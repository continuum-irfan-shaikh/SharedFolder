package remoteAccess

import "time"

//EndpointConfigurationReportReqBody represents http POST request body structure
type EndpointConfigurationReportReqBody struct {
	Vendor        string   `json:"vendor"`
	ClientID      string   `json:"clientId"`
	SiteID        string   `json:"siteId"`
	EndpointsList []string `json:"endpointsList"`
}

//EndpointConfigurationReportNode represents Details for Endpoint Configuration Report Record
type EndpointConfigurationReportNode struct {
	PartnerID       string    `json:"partnerId"`
	EndpointID      string    `json:"endpointId"`
	ClientID        string    `json:"clientId"`
	SiteID          string    `json:"siteId"`
	RegID           string    `json:"regId"`
	AgentInstalled  int       `json:"agentInstalled"`
	DCDtime         time.Time `json:"dcdtime"`
	EndpointType    string    `json:"endpointType"`
	OneClickEnabled bool      `json:"oneClickEnabled"`
}

//EndpointConfigurationReportData represents Details for Endpoint Configuration Report
type EndpointConfigurationReportData struct {
	OutData []EndpointConfigurationReportNode `json:"outData"`
}
