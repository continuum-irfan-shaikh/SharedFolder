package remoteAccess

import "time"

//EnableManifestReportReqBody represents http POST request body structure
type EnableManifestReportReqBody struct {
	Vendor        string   `json:"vendor"`
	ClientID      string   `json:"clientId"`
	SiteID        string   `json:"siteId"`
	EndpointsList []string `json:"endpointsList"`
}

//EnableManifestReportNode represents Details for Enable Manifest Report Record
type EnableManifestReportNode struct {
	PartnerID       string    `json:"partnerId"`
	EndpointID      string    `json:"endpointId"`
	ClientID        string    `json:"clientId"`
	SiteID          string    `json:"siteId"`
	ManifestVersion string    `json:"manifestVersion"`
	RetryCount      int       `json:"retryCount"`
	DCDtime         time.Time `json:"dcdtime"`
}

//EnableManifestReportData represents Details for Enable Manifest Report
type EnableManifestReportData struct {
	OutData []EnableManifestReportNode `json:"outData"`
}

//EnableManifestResetReqBody represents http DELETE request body structure
type EnableManifestResetReqBody struct {
	Vendor        string   `json:"vendor"`
	EndpointsList []string `json:"endpointsList"`
}
