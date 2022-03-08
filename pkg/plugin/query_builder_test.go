package plugin

import (
	"strings"
	"testing"
)

func Test_buildQuery(t *testing.T) {
	type args struct {
		d  *YouboraDataSource
		qm *QueryModel
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Simple Views",
			args: args{
				d: &YouboraDataSource{
					apikey:  "mytestapikey",
					baseurl: "https://api.youbora.com",
					account: "mytestaccount",
				},
				qm: &QueryModel{
					FromDate:      "1646220720",
					Granularity:   "minute",
					Metrics:       []string{"views"},
					StreamingType: []string{"ALL"},
				},
			},
			want: "https://api.youbora.com/mytestaccount/data?fromDate=1646220720&granularity=minute&metrics=views&timezone=GMT&type=ALL",
		},
		{
			name: "Filtered",
			args: args{
				d: &YouboraDataSource{
					apikey:  "mytestapikey",
					baseurl: "https://api.youbora.com",
					account: "mytestaccount",
				},
				qm: &QueryModel{
					FromDate:      "1646220720",
					Granularity:   "minute",
					Metrics:       []string{"views"},
					StreamingType: []string{"ALL"},
					FilterBy:      "country",
					FilterValue:   "Brazil",
				},
			},
			want: "https://api.youbora.com/mytestaccount/data?filter=%5B%7B%22name%22%3A%22country%3DBrazil%22%2C%22rules%22%3A%7B%22country%22%3A%5B%22Brazil%22%5D%7D%7D%5D&fromDate=1646220720&granularity=minute&metrics=views&timezone=GMT&type=ALL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildQuery(tt.args.d, tt.args.qm)
			i := strings.Index(got, "&dateToken")
			got = got[:i]
			if got != tt.want {
				t.Errorf("buildQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
