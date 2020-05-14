package remoteAccess

import "time"

//SessionAuditReportReqBody represents http POST request body structure
type SessionAuditReportReqBody struct {
	Vendor        string   `json:"vendor"`
	EndpointsList []string `json:"endpointsList"`
	SessionID     string   `json:"sessionId"`
}

//SessionAuditReportNode represents Details for Session Audit Report Record
type SessionAuditReportNode struct {
	PartnerID          string    `json:"partnerId"`
	EndpointID         string    `json:"endpointId"`
	SessionID          string    `json:"sessionId"`
	RegID              string    `json:"regId"`
	ITSUserID          string    `json:"itsUserId"`
	Reason             string    `json:"reason"`
	RemoteAccessHostID string    `json:"remoteAccessHostId"`
	Dcdtime            time.Time `json:"dcdtime"`
	Referer            string    `json:"referer"`
	RUsername          string    `json:"rUsername"`
	NOCUserID          string    `json:"nocUserId"`
	PortalType         string    `json:"portalType"`
	TempURL            string    `json:"tempURL"`
	ResourceType       string    `json:"resourceType"`
}

//SessionAuditReportData represents Details for Session Audit Report
type SessionAuditReportData struct {
	OutData []SessionAuditReportNode `json:"outData"`
}
