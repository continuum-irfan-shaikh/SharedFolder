package remoteAccess

import "time"

//InstallationStatusReportReqBody represents http POST request body structure
type InstallationStatusReportReqBody struct {
	Vendor        string   `json:"vendor"`
	ClientID      string   `json:"clientId"`
	SiteID        string   `json:"siteId"`
	EndpointsList []string `json:"endpointsList"`
	CompactFlag   bool     `json:"compactFlag"`
}

//InstallationStatusReportNode represents Details for Installation Status Report Record
type InstallationStatusReportNode struct {
	PartnerID          string     `json:"partnerId,omitempty"`
	EndpointID         string     `json:"endpointId"`
	ClientID           string     `json:"clientId,omitempty"`
	SiteID             string     `json:"siteId,omitempty"`
	RegID              string     `json:"regId,omitempty"`
	RemoteAccessHostID string     `json:"remoteAccessHostId,omitempty"`
	EndpointType       string     `json:"endpointType,omitempty"`
	InstallationStatus string     `json:"installationStatus"`
	ErrorCode          string     `json:"errorCode,omitempty"`
	VerificationFlag   bool       `json:"verificationFlag"`
	DCDtime            *time.Time `json:"dcdtime,omitempty"`
	ControlSessionID    string    `json:"controlsessionid,omitempty"`
}

//InstallationStatusReportData represents Details for Installation Status Report
type InstallationStatusReportData struct {
	OutData []InstallationStatusReportNode `json:"outData"`
}
