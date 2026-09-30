package plugin

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"

	_ "github.com/sijms/go-ora/v2"
)

var (
	_ backend.QueryDataHandler      = (*OracleDatasource)(nil)
	_ backend.CheckHealthHandler    = (*OracleDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*OracleDatasource)(nil)
)

type OracleDatasource struct {
	connection OracleDatasourceConnection
	name       string
	settings   OracleDatasourceSettings
}

func NewDatasource(_ context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	datasourceSettings := ParseDatasourceSettings(settings.JSONData, settings.DecryptedSecureJSONData)
	log.DefaultLogger.Debug("Creating new Oracle datasource instance", "name", settings.Name)
	return &OracleDatasource{
		connection: OracleDatasourceConnection{},
		name:       settings.Name,
		settings:   datasourceSettings,
	}, nil
}

func (d *OracleDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	d.name = req.PluginContext.DataSourceInstanceSettings.Name
	d.settings = ParseDatasourceSettings(req.PluginContext.DataSourceInstanceSettings.JSONData, req.PluginContext.DataSourceInstanceSettings.DecryptedSecureJSONData)

	err := d.connection.Reconnect(&d.settings)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("Oracle connection error: %v", err),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Oracle datasource successfully connected (Read-Only Mode Active)",
	}, nil
}

func (d *OracleDatasource) Dispose() {
	if err := d.connection.Disconnect(); err != nil {
		log.DefaultLogger.Error("Error closing Oracle connection during Dispose", "error", err)
	}
}

func (d *OracleDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	var connErr error
	if !d.connection.IsConnected() {
		connErr = d.connection.Connect(&d.settings)
	}

	for _, q := range req.Queries {
		if connErr != nil {
			response.Responses[q.RefID] = backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("Error connecting to Oracle: %v", connErr))
			continue
		}

		queryObj := OracleDatasourceQuery{}
		if err := queryObj.ParseDatasourceQuery(q); err != nil {
			response.Responses[q.RefID] = backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("Error parsing query: %v", err))
			continue
		}

		response.Responses[q.RefID] = queryObj.MakeQuery(ctx, &d.connection)
	}

	return response, nil
}
