package cloud

//ResourceMetrics stores metrics data
type ResourceMetrics struct {
	ID         string        `json:"id"`
	Name       MetricsName   `json:"name"`
	Unit       string        `json:"unit"`
	Timeseries []MetricsData `json:"timeseries"`
}

//MetricsName stores information of metrics name
type MetricsName struct {
	Value          string `json:"value"`
	LocalizedValue string `json:"localizedValue"`
}

//MetricsData stores array of timeseries data
type MetricsData struct {
	Data []TimeSeriesData `json:"data"`
}

//TimeSeriesData stores values of time series data
type TimeSeriesData struct {
	TimeStamp string  `json:"timeStamp,omitempty"`
	Total     float64 `json:"total,omitempty"`
	Count     float64 `json:"count,omitempty"`
	Average   float64 `json:"average,omitempty"`
	Minimum   float64 `json:"minimum,omitempty"`
	Maximum   float64 `json:"maximum,omitempty"`
}
