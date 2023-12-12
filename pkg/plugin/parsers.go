package plugin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func ParseYouboraResponse(yr *YouboraResponse) (frames data.Frames, err error) {

	for i := 0; i < len(yr.Data); i++ {
		var frameName string
		if yr.Data[i].Name != "" {
			frameName = fmt.Sprintf("%s | %s", yr.Data[i].Name, yr.Data[i].Type)
		} else {
			frameName = yr.Data[i].Type
		}
		frame := data.NewFrame(frameName)

		for k := 0; k < len(yr.Data[i].Metrics); k++ {
			fieldName := fmt.Sprintf("(%s)", yr.Data[i].Metrics[k].Label)
			n := len(yr.Data[i].Metrics[k].Values[0].Points)
			x := make([]time.Time, n)
			y := make([]float64, n)

			for j := 0; j < n; j++ {
				x[j] = time.Unix(int64(yr.Data[i].Metrics[k].Values[0].Points[j][0]/1000), 0)
				y[j] = yr.Data[i].Metrics[k].Values[0].Points[j][1]
			}
			// add only with time field per data frame.
			if k == 0 {
				frame.Fields = append(frame.Fields, data.NewField("time", nil, x))
			}
			frame.Fields = append(frame.Fields, data.NewField(fieldName, nil, y))
		}

		frames = append(frames, frame)
	}

	return frames, nil
}

func ParseQuery(dq backend.DataQuery, qm *QueryModel) error {

	err := json.Unmarshal(dq.JSON, &qm)
	log.DefaultLogger.Debug("QUERY", "data", dq.JSON)

	// setup time range.
	qm.FromDate = fmt.Sprintf("%d", dq.TimeRange.From.UnixNano()/int64(time.Millisecond))
	qm.ToDate = fmt.Sprintf("%d", dq.TimeRange.To.UnixNano()/int64(time.Millisecond))

	// setup granularity.
	qm.Granularity = "minute"
	if dq.Interval.Minutes() > 120 {
		qm.Granularity = "hour"
	}
	if dq.Interval.Hours() > 48 {
		qm.Granularity = "day"
	}

	return err
}
