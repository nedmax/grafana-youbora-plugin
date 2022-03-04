package plugin

type YouboraResponse struct {
	Messages []Message `json:"messages,omitempty"`
	Data     []Datum   `json:"data,omitempty"`
	Metadata Metadata  `json:"metadata,omitempty"`
}
type Message struct {
	Class     string `json:"class,omitempty"`
	Text      string `json:"message,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

type Datum struct {
	Name    string    `json:"name,omitempty"`
	Type    string    `json:"type"`
	Date    []float64 `json:"date"`
	Metrics []Metric  `json:"metrics,omitempty"`
}

type Metadata struct {
	Account   string  `json:"account,omitempty"`
	Timestamp float64 `json:"timestamp,omitempty"`
	Refresh   int16   `json:"refresh,omitempty"`
}

type Metric struct {
	Values     []Value `json:"values,omitempty"`
	Code       string  `json:"code,omitempty"`
	Label      string  `json:"label,omitempty"`
	Magnitudes struct {
		X string `json:"x,omitempty"`
		Y string `json:"y,omitempty"`
	} `json:"magnitudes,omitempty"`
}

type Value struct {
	Points [][2]float64 `json:"data,omitempty"`
}

type YouboraDataSource struct {
	apikey  string
	baseurl string
	account string
}

type JsonData struct {
	Account string `json:"account"`
}

type QueryModel struct {
	FromDate      string
	ToDate        string
	Filter        string   `json:"filter,omitempty"`
	StreamingType []string `json:"streamingType"`
	Metrics       []string `json:"metrics"`
	Granularity   string
}
