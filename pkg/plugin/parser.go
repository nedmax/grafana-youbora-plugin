package plugin

import (
	"time"
)

func ParseYouboraResponse(yr *YouboraResponse) ([]time.Time, []float64, error) {
	n := len(yr.Data[0].Metrics[0].Values[0].Points)
	x := make([]time.Time, n)
	y := make([]float64, n)

	for i := 0; i < n; i++ {
		x[i] = time.Unix(int64(yr.Data[0].Metrics[0].Values[0].Points[i][0]/1000), 0)
		y[i] = yr.Data[0].Metrics[0].Values[0].Points[i][1]
	}

	return x, y, nil
}
