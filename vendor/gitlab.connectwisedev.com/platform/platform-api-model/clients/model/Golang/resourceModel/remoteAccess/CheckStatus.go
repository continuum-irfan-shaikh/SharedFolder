package remoteAccess

//CheckStatus represents message format for checkStatus endpoint
type CheckStatus struct {
	PartnerID  int64  `json:"partnerId"`
	EndpointID string `json:"endpointId"`
	RegID      string `json:"regId"`
	OnlineFlag bool   `json:"onlineFlag"`
}

//CheckStatusReqBody represents http POST request body structure
type CheckStatusReqBody struct {
	Vendor             string `json:"vendor"`
	RemoteAccessHostID string `json:"remoteAccessHostId"`
	RegID              string `json:"regId"`
	EndpointID         string `json:"endpointId"`
}
