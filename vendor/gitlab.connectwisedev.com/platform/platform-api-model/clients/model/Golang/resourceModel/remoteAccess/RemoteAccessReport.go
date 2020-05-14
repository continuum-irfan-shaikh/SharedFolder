package remoteAccess

import "time"

//RemoteAccessReportReqBody represents http POST request body structure
type RemoteAccessReportReqBody struct {
	Vendor    string `json:"vendor"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

//RemoteAccessReportNode represents Details for Remote Access Report Record
type RemoteAccessReportNode struct {
	PartnerID          string    `json:"partnerId"`
	EndpointID         string    `json:"endpointId"`
	SessionID          string    `json:"sessionId"`
	RemoteAccessHostID string    `json:"remoteAccessHostId"`
	StartTime          time.Time `json:"startTime"`
	EndTime            time.Time `json:"endTime"`
	UserIP             string    `json:"userIP"`
	Dcdtime            time.Time `json:"dcdtime"`
	ResourceType       string    `json:"resourceType"`
	RegID              string    `json:"regId"`
	ITSUserID          string    `json:"itsUserId"`
	Reason             string    `json:"reason"`
	Referer            string    `json:"referer"`
	RUsername          string    `json:"rUsername"`
	NOCUserID          string    `json:"nocUserId"`
	PortalType         string    `json:"portalType"`
	TempURL            string    `json:"tempURL"`
}

//RemoteAccessReportData represents Details for Remote Access Report
type RemoteAccessReportData struct {
	OutData []RemoteAccessReportNode `json:"outData"`
}
