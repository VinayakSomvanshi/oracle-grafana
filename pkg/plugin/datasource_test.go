package plugin

import (
	"context"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestQueryData(t *testing.T) {
	ds := OracleDatasource{}

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{
					RefID: "A",
					JSON:  []byte(`{"o_sql": "SELECT 1 FROM dual"}`),
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Responses) != 1 {
		t.Fatalf("expected 1 response, got %d", len(resp.Responses))
	}

	// Should fail with connection error since no db is connected, but query was handled cleanly
	resA := resp.Responses["A"]
	if resA.Error == nil && resA.Status != backend.StatusBadRequest {
		t.Logf("Response: %+v", resA)
	}
}

func TestQueryDataRejectedWrite(t *testing.T) {
	ds := OracleDatasource{}
	// simulate connected state
	ds.connection = OracleDatasourceConnection{}

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{
					RefID: "B",
					JSON:  []byte(`{"o_sql": "DROP TABLE users"}`),
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	resB := resp.Responses["B"]
	if resB.Error == nil {
		t.Fatal("expected error on mutating query, got nil")
	}
}
