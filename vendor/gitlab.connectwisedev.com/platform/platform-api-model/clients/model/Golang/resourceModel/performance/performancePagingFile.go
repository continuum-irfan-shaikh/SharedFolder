package performance

import "time"

//PagingFiles Overall paging file performance data
type PagingFiles struct {
	EndpointID    string       `json:"endpointID"`
	PartnerID     string       `json:"partnerID"`
	ClientID      string       `json:"clientID"`
	SiteID        string       `json:"siteID"`
	Name          string       `json:"name"`
	CreateTimeUTC time.Time    `json:"createTimeUTC"`
	CreatedBy     string       `json:"createdBy"`
	Metric        PagingMetric `json:"metric" cql:"metric"`
	PagingFiles   []PagingFile `json:"pagingFiles" cql:"paging_files"`
}

//PagingFile Single paging file perf data
type PagingFile struct {
	Name   string       `json:"name" cql:"name"`
	Metric PagingMetric `json:"metric" cql:"metric"`
}

//PagingMetric paging file usage details in percentage
type PagingMetric struct {
	PercentUsage     float32 `json:"percentUsage" cql:"percent_usage"`
	PercentPeakUsage float32 `json:"percentPeakUsage" cql:"percent_peakusage"`
}
