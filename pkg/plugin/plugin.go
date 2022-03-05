package plugin

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
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

var httpClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			Renegotiation: tls.RenegotiateFreelyAsClient,
		},
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 3 * time.Second,
			DualStack: true,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
	},
	Timeout: time.Duration(time.Second * 5),
}

// NewYouboraDataSource creates a new datasource instance.
func NewYouboraDataSource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var secureData = settings.DecryptedSecureJSONData
	var jsondata JsonData

	if err := json.Unmarshal(settings.JSONData, &jsondata); err != nil {
		log.DefaultLogger.Warn("Error getting API key", "err", err)
	}

	return &YouboraDataSource{
		apikey:  secureData["apikey"],
		baseurl: "https://api.youbora.com",
		account: jsondata.Account,
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
	log.DefaultLogger.Info("QueryData called", "request", req)

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
	log.DefaultLogger.Info("CheckHealth called", "request", req)
	var status = backend.HealthStatusOk
	var message = "Data source is working"
	var qm = &QueryModel{
		FromDate: "last5minutes",
		Metrics:  []string{"views"},
	}

	_, err := d.doRequest(ctx, qm)

	if err != nil {
		status = backend.HealthStatusError
		message = "error getting API version"
		return nil, err
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

func (d *YouboraDataSource) doRequest(ctx context.Context, qm *QueryModel) (body []byte, err error) {

	const ttl = 20 * 60 * 1000 // 20 minutes in milliseconds
	expirationTime := (time.Now().UnixNano() / int64(time.Millisecond)) + ttl

	basePath := fmt.Sprintf("/%s/data", d.account)
	orderedParams := fmt.Sprintf(
		"fromDate=%s&metrics=%s&type=%s&timezone=GMT&granularity=%s",
		qm.FromDate,
		strings.Join(qm.Metrics, ","),
		strings.Join(qm.StreamingType, ","),
		qm.Granularity,
	)
	if qm.ToDate != "" {
		orderedParams = orderedParams + fmt.Sprintf("&toDate=%s", qm.ToDate)
	}
	if qm.GroupBy != "" {
		orderedParams = orderedParams + fmt.Sprintf("&groupBy=%s", qm.GroupBy)
	}

	baseParams := fmt.Sprintf("dateToken=%d&%s", expirationTime, orderedParams)
	baseToken := fmt.Sprintf("%s?%s", basePath, baseParams)
	token := md5.Sum([]byte(baseToken + d.apikey))

	url := fmt.Sprintf("%s%s?%s&token=%s", d.baseurl, basePath, baseParams, hex.EncodeToString(token[:]))
	log.DefaultLogger.Debug("calling API", "url", url)

	resp, err := httpClient.Get(url)
	if err != nil {
		return body, err
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)

	return body, nil
}
