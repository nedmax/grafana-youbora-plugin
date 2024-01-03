package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

// Make sure YouboraDataSource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler, backend.StreamHandler interfaces. Plugin should not
// implement all these interfaces - only those which are required for a particular task.
// For example if plugin does not need streaming functionality then you are free to remove
// methods that implement backend.StreamHandler. Implementing instancemgmt.InstanceDisposer
// is useful to clean up resources used by previous datasource instance when a new datasource
// instance created upon datasource settings changed.
var (
	_ backend.QueryDataHandler      = (*YouboraDataSource)(nil)
	_ backend.CheckHealthHandler    = (*YouboraDataSource)(nil)
	_ instancemgmt.InstanceDisposer = (*YouboraDataSource)(nil)
)

// NewYouboraDataSource creates a new datasource instance.
func NewYouboraDataSource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var secureData = settings.DecryptedSecureJSONData
	var jsondata JsonData

	client, err := httpclient.New(httpclient.Options{
		Timeouts: &httpclient.TimeoutOptions{
			Timeout: 5 * time.Second,
		},
	})
	if err != nil {
		log.DefaultLogger.Error("failed to create HTTP client.", "error", err)
	}

	if err := json.Unmarshal(settings.JSONData, &jsondata); err != nil {
		log.DefaultLogger.Error("Error getting API key.", "error", err)
		return nil, err
	}

	return &YouboraDataSource{
		apikey:     secureData["apikey"],
		baseurl:    "https://api.youbora.com",
		account:    jsondata.Account,
		httpclient: client,
	}, nil
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewYouboraDataSource factory function.
func (d *YouboraDataSource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *YouboraDataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

func (d *YouboraDataSource) query(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}
	qm := QueryModel{}
	yr := YouboraResponse{}

	// Unmarshal the JSON into our queryModel.
	response.Error = ParseQuery(query, &qm)
	if response.Error != nil {
		return response
	}

	// query Youbora API.
	body, err := d.doRequest(ctx, &qm)
	if err != nil {
		response.Error = err
		return response
	}
	if err := json.Unmarshal(body, &yr); err != nil {
		response.Error = err
		return response
	}

	// create data frame response.
	frames, err := ParseYouboraResponse(&yr)
	if err != nil {
		response.Error = err
		return response
	}

	// add the frames to the response.
	response.Frames = frames

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *YouboraDataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	yr := YouboraResponse{}
	var message = "Data source is working"
	var qm = &QueryModel{
		FromDate:      "last5minutes",
		Metrics:       []string{"views"},
		Granularity:   "minute",
		StreamingType: []string{"ALL"},
	}

	body, err := d.doRequest(ctx, qm)
	if err != nil {
		log.DefaultLogger.Error("doRequest", "error", err)
		status = backend.HealthStatusError
		message = "error getting API data"
		return &backend.CheckHealthResult{
			Status:  status,
			Message: message,
		}, err
	}
	if err := json.Unmarshal(body, &yr); err != nil {
		log.DefaultLogger.Error("doRequest", "body", body)
		status = backend.HealthStatusError
		message = "error parsing API response"
		return &backend.CheckHealthResult{
			Status:  status,
			Message: message,
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

func (d *YouboraDataSource) doRequest(ctx context.Context, qm *QueryModel) (body []byte, err error) {
	url := buildQuery(d, qm)

	rsp, err := d.httpclient.Get(url)
	log.DefaultLogger.Debug("DEBUG URL", "url", url)

	if err != nil {
		log.DefaultLogger.Error("failed executing GET.", "error", err)

		return body, err
	}
	defer rsp.Body.Close()

	if body, err = io.ReadAll(rsp.Body); err != nil {
		return nil, err
	}

	if rsp.StatusCode > 399 {
		var parsed map[string]interface{}
		err := json.Unmarshal(body, &parsed)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTP response %v (invalid response json)", rsp.Status)
		}
		return nil, fmt.Errorf("invalid API request: %v", parsed)
	}

	return body, err
}
