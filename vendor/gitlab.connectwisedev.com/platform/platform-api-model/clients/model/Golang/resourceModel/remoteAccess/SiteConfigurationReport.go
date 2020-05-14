package remoteAccess

import "time"

//SiteConfigurationReportReqBody represents http POST request body structure
type SiteConfigurationReportReqBody struct {
	Vendor   string `json:"vendor"`
	ClientID string `json:"clientId"`
	SiteID   string `json:"siteId"`
}

//SiteConfigurationReportNode represents Details for Site Configuration Report Record
type SiteConfigurationReportNode struct {
	PartnerID       string    `json:"partnerId"`
	ClientID        string    `json:"clientId"`
	SiteID          string    `json:"siteId"`
	EndpointType    string    `json:"endpointType"`
	OneClickEnabled bool      `json:"oneClickEnabled"`
	Status          int       `json:"status"`
	DCDtime         time.Time `json:"dcdtime"`
}

//SiteConfigurationReportData represents Details for Site Configuration Report
type SiteConfigurationReportData struct {
	OutData []SiteConfigurationReportNode `json:"outData"`
}
