package plugin_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-starter-datasource-backend/pkg/plugin"
)

// This is where the tests for the datasource backend live.
func TestQueryData(t *testing.T) {
	ds := plugin.YouboraDataSource{}

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A"},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Responses) != 1 {
		t.Fatal("QueryData must return a response")
	}
}

func TestReadResponseData(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/simple_views.json")
	if err != nil {
		t.Fatal(err)
	}

	var result plugin.YouboraResponse
	err = json.Unmarshal(b, &result)

	if err != nil {
		t.Error(err)
	}

	if len(result.Messages) > 0 {
		t.Errorf("%v: output doesn't match expected result", result)
	}

	if result.Metadata.Account != "globo" {
		t.Errorf("%v: output doesn't match expected result", result)
	}

	if result.Data[0].Date[0] != 1646152500000 {
		t.Errorf("%v: output doesn't match expected result", result)
	}

	if result.Data[0].Metrics[0].Values[0].Points[0].Value != 12554369 {
		t.Errorf("%v: output doesn't match expected result", result)
	}

}
