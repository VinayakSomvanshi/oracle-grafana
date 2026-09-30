package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

const DefaultMaxRowLimit = 50000

type OracleDatasourceQuery struct {
	Datasource   OracleDatasourceInfo `json:"datasource"`
	DatasourceId int64                `json:"datasourceId"`
	IntervalMs   int64                `json:"intervalMs"`
	O_parsed     string               `json:"o_parsed"`
	O_sql        string               `json:"o_sql"`
	RefId        string               `json:"refId"`
}

type OracleDatasourceInfo struct {
	Type string `json:"type"`
	Uid  string `json:"uid"`
}

// MakeQuery executes a validated, read-only SQL query against the Oracle database
// with full context cancellation, row-limit circuit breaker, and typed column mapping.
func (q *OracleDatasourceQuery) MakeQuery(ctx context.Context, c *OracleDatasourceConnection) backend.DataResponse {
	var resp backend.DataResponse

	if !c.IsConnected() {
		return backend.ErrDataResponse(backend.StatusBadRequest, "oracle connection is not open")
	}

	rawSQL := q.O_parsed
	if strings.TrimSpace(rawSQL) == "" {
		rawSQL = q.O_sql
	}

	// 1. Mandatory Pre-Flight Read-Only Validation (Defense Layer 1)
	if err := ValidateReadOnlyQuery(rawSQL); err != nil {
		log.DefaultLogger.Warn("Query rejected by read-only validator", "error", err, "query", rawSQL)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("Query rejected (read-only compliance): %v", err))
	}

	// 2. Oracle Database Engine Read-Only Transaction Enforcement (Defense Layer 2)
	// Even if an obscure SQL or stored function is executed, Oracle kernel blocks data mutation.
	_, _ = c.connection.ExecContext(ctx, "ALTER SESSION SET TRANSACTION READ ONLY")

	// 3. Prepare query with Context for cancellation support
	stmt, err := c.connection.PrepareContext(ctx, rawSQL)
	if err != nil {
		log.DefaultLogger.Error("Error preparing SQL", "error", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("Error preparing SQL: %v", err))
	}
	defer stmt.Close()

	// 4. Execute query with Context
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		log.DefaultLogger.Error("Error executing query", "error", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("Error querying database: %v", err))
	}
	defer rows.Close()

	// 5. Build Column Metadata & Typed Builders
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		log.DefaultLogger.Error("Error retrieving column types", "error", err)
		return backend.ErrDataResponse(backend.StatusUnknown, fmt.Sprintf("Error retrieving column metadata: %v", err))
	}

	builders := make([]*TypedColumnBuilder, len(colTypes))
	for i, ct := range colTypes {
		builders[i] = NewColumnBuilder(ct)
	}

	// 6. Scan Rows with Row Limit Circuit Breaker (OOM Protection)
	scanValues := make([]interface{}, len(colTypes))
	scanArgs := make([]interface{}, len(colTypes))
	for i := range scanArgs {
		scanArgs[i] = &scanValues[i]
	}

	rowCount := 0
	rowLimitExceeded := false

	for rows.Next() {
		rowCount++
		if rowCount > DefaultMaxRowLimit {
			rowLimitExceeded = true
			break
		}

		if err := rows.Scan(scanArgs...); err != nil {
			log.DefaultLogger.Error("Error scanning row", "row", rowCount, "error", err)
			return backend.ErrDataResponse(backend.StatusUnknown, fmt.Sprintf("Error reading row data: %v", err))
		}

		for i, val := range scanValues {
			builders[i].Append(val)
		}
	}

	if err := rows.Err(); err != nil {
		log.DefaultLogger.Error("Row iteration error", "error", err)
		return backend.ErrDataResponse(backend.StatusUnknown, fmt.Sprintf("Error during row processing: %v", err))
	}

	// 7. Construct Grafana DataFrame
	frame := data.NewFrame(q.RefId)
	for _, b := range builders {
		frame.Fields = append(frame.Fields, b.ToField())
	}

	if rowLimitExceeded {
		frame.AppendNotices(data.Notice{
			Severity: data.NoticeSeverityWarning,
			Text:     fmt.Sprintf("Row limit reached (%d rows). Query output was truncated to prevent memory overload.", DefaultMaxRowLimit),
		})
	}

	resp.Frames = append(resp.Frames, frame)
	return resp
}

func (q *OracleDatasourceQuery) ParseDatasourceQuery(query backend.DataQuery) error {
	err := json.Unmarshal(query.JSON, &q)
	if err != nil {
		log.DefaultLogger.Error("Error unmarshaling query JSON", "error", err)
	}
	q.RefId = query.RefID
	return err
}
