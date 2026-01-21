package plugin

import (
	"fmt"
	"net/url"
	"strings"
)

func buildQuery(d *YouboraDataSource, qm *QueryModel) string {
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

	baseParams := params.Encode()
	return fmt.Sprintf("%s%s?%s", d.baseurl, basePath, baseParams)
}
