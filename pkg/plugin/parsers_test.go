package plugin_test

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-starter-datasource-backend/pkg/plugin"
)

func TestParseQuery(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/query.json")
	if err != nil {
		t.Fatal(err)
	}
	var qm plugin.QueryModel
	var query backend.DataQuery
	query.JSON = b
	err = plugin.ParseQuery(query, &qm)

	if err != nil {
		t.Errorf("%v: error parsing JSON", string(b))
	}

	if strings.Join(qm.StreamingType, ",") != "ALL,VOD" {
		t.Errorf("%v: output doesn't match expected result", qm.StreamingType)
	}

	if strings.Join(qm.Metrics, ",") != "views,concurrent" {
		t.Errorf("%v: output doesn't match expected result", qm.Metrics)
	}

}

func TestParseYouboraResponse(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/simple_views.json")
	if err != nil {
		t.Fatal(err)
	}

	var result plugin.YouboraResponse
	err = json.Unmarshal(b, &result)

	if err != nil {
		t.Errorf("%v: error parsing JSON", string(b))
	}

	if len(result.Messages) > 0 {
		t.Errorf("%v: output doesn't match expected result", result)
	}

	if result.Metadata.Account != "globo" {
		t.Errorf("%v: output doesn't match expected result", result)
	}

	if result.Data[0].Date[0] != 1646220720000 {
		t.Errorf("%v: output doesn't match expected result", result)
	}

	frames, err := plugin.ParseYouboraResponse(&result)
	if err != nil {
		t.Fatalf("%v: error parsing Youbora response", result)
	}

	if frames[0].Fields[0].At(0) != time.Unix(1646220720, 0) {
		t.Errorf("%v != %v: output doesn't match expected result", frames[0].Fields[0].At(0), time.Unix(1646220720, 0))
	}

	if frames[0].Fields[1].At(0) != float64(15655) {
		t.Errorf("%v != %v: output doesn't match expected result", frames[0].Fields[1].At(0), 15655)
	}

}

func TestDoubleData(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/double_data.json")
	if err != nil {
		t.Fatal(err)
	}

	var result plugin.YouboraResponse
	err = json.Unmarshal(b, &result)
	if err != nil {
		t.Errorf("%v: error parsing JSON", string(b))
	}

	frames, err := plugin.ParseYouboraResponse(&result)
	if err != nil {
		t.Fatalf("%v: error parsing Youbora response", result)
	}

	if frames[0].Fields[0].At(0) != time.Unix(1646232180, 0) {
		t.Errorf("%v != %v: output doesn't match expected result", frames[0].Fields[0].At(0), time.Unix(1646232180, 0))
	}

	if frames[0].Fields[1].At(0) != float64(27585) {
		t.Errorf("%v != %v: output doesn't match expected result", frames[0].Fields[1].At(0), 27585)
	}

	if frames[0].Fields[1].Name != "(Plays)" {
		t.Errorf("%v: output doesn't match expected result", frames[0].Fields[1].Name)
	}

	if frames[0].Fields[2].Name != "(Concurrent Plays)" {
		t.Errorf("%v: output doesn't match expected result", frames[0].Fields[2].Name)
	}

	if frames[0].Name != "ALL" {
		t.Errorf("%v: output doesn't match expected result", frames[0].Name)
	}

}

func TestWithDimension(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/with_dimension.json")
	if err != nil {
		t.Fatal(err)
	}

	var result plugin.YouboraResponse
	err = json.Unmarshal(b, &result)
	if err != nil {
		t.Errorf("%v: error parsing JSON", string(b))
	}

	frames, err := plugin.ParseYouboraResponse(&result)
	if err != nil {
		t.Fatalf("%v: error parsing Youbora response", result)
	}

	if frames[0].Fields[0].At(0) != time.Unix(1646489220, 0) {
		t.Errorf("%v != %v: output doesn't match expected result", frames[0].Fields[0].At(0), time.Unix(1646489220, 0))
	}

	if frames[0].Fields[1].At(0) != float64(23196) {
		t.Errorf("%v != %v: output doesn't match expected result", frames[0].Fields[1].At(0), 23196)
	}

	if frames[0].Name != "Globo | ALL" {
		t.Errorf("%v: output doesn't match expected result", frames[0].Name)
	}

}
