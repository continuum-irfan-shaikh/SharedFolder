package remoteAccess

//Connection represents message format for connection endpoint
type Connection struct {
	GatewayURL         string `json:"gatewayUrl"`
	SessionTicketID    int64  `json:"sessionTicketId"`
	RemoteAccessHostID int64  `json:"remoteAccessHostId"`
	PartnerID          int64  `json:"partnerId"`
	UserID             int64  `json:"userId"`
	ClientID           string `json:"clientId"`
	SiteID             string `json:"siteId"`
	EndpointID         string `json:"endpointId"`
	RegID              string `json:"regId"`
	MessageControl     string `json:"message"`
}

//ConnectionReqBody represents http POST request body structure
type ConnectionReqBody struct {
	Vendor                      string `json:"vendor"`
	ResourceType                string `json:"resourceType"`
	OneClickFlag                bool   `json:"oneClickFlag"`
	RemoteAccessHostID          string `json:"remoteAccessHostId"`
	RegID                       string `json:"regId"`
	EndpointID                  string `json:"endpointId"`
	Reason                      string `json:"reason"`
	Referer                     string `json:"referer"`
	RUserName                   string `json:"rUserName"`
	NOCUserID                   string `json:"nocUserID"`
	PortalType                  string `json:"portalType"`
	IsNOCFlag                   bool   `json:"isNOCFlag"`
	DirectLink                  bool   `json:"directLink"`
	FailWhenRaSessionInProgress bool   `json:"failWhenRaSessionInProgress"`
	QuickLaunch                 string `json:"quickLaunch"`
	SessionID                   string `json:"sessionID"`
	UserDisplayName             string `json:"userDisplayName"`
}
