package plugin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"
)

func buildQuery(d *YouboraDataSource, qm *QueryModel) string {
	// 20 minutes in milliseconds
	const ttl = 20 * 60 * 1000
	expirationTime := (time.Now().UnixNano() / int64(time.Millisecond)) + ttl

	basePath := fmt.Sprintf("/%s/data", d.account)
	params := url.Values{}
	params.Add("fromDate", qm.FromDate)
	if qm.ToDate != "" {
		params.Add("toDate", qm.ToDate)
	}
	params.Add("granularity", qm.Granularity)
	params.Add("timezone", "GMT")
	params.Add("metrics", strings.Join(qm.Metrics, ","))
	params.Add("type", strings.Join(qm.StreamingType, ","))
	if qm.GroupBy != "" {
		params.Add("groupBy", qm.GroupBy)
	}
	if qm.FilterBy != "" && qm.FilterValue != "" {
		params.Add("filter", fmt.Sprintf(
			"[{\"name\":\"%s=%s\",\"rules\":{\"%s\":[\"%s\"]}}]",
			qm.FilterBy,
			qm.FilterValue,
			qm.FilterBy,
			qm.FilterValue,
		))
	}

	baseParams := fmt.Sprintf("%s&dateToken=%d", params.Encode(), expirationTime)
	baseToken := fmt.Sprintf("%s?%s", basePath, baseParams)
	token := md5.Sum([]byte(baseToken + d.apikey))

	url := fmt.Sprintf("%s%s?%s&token=%s", d.baseurl, basePath, baseParams, hex.EncodeToString(token[:]))
	return url
}
