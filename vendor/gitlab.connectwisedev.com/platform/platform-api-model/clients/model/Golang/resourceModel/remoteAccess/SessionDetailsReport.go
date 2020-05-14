package remoteAccess

import "time"

//SessionDetailsReportReqBody represents http POST request body structure
type SessionDetailsReportReqBody struct {
	Vendor        string   `json:"vendor"`
	EndpointsList []string `json:"endpointsList"`
	SessionID     string   `json:"sessionId"`
}

//SessionDetailsReportNode represents Details for Session Details Report Record
type SessionDetailsReportNode struct {
	PartnerID          string    `json:"partnerId"`
	EndpointID         string    `json:"endpointId"`
	SessionID          string    `json:"sessionId"`
	RemoteAccessHostID string    `json:"remoteAccessHostId"`
	StartTime          time.Time `json:"startTime"`
	EndTime            time.Time `json:"endTime"`
	UserIP             string    `json:"userIP"`
	Dcdtime            time.Time `json:"dcdtime"`
	ResourceType       string    `json:"resourceType"`
}

//SessionDetailsReportData represents Details for Session Details Report
type SessionDetailsReportData struct {
	OutData []SessionDetailsReportNode `json:"outData"`
}
