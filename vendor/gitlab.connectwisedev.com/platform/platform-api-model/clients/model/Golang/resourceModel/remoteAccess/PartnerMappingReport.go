package remoteAccess

//PartnerEndpointMappingReportReqBody represents http POST request body structure
type PartnerEndpointMappingReportReqBody struct {
	Vendor        string   `json:"vendor"`
	ClientID      string   `json:"clientId"`
	SiteID        string   `json:"siteId"`
	EndpointsList []string `json:"endpointsList"`
}

//PartnerRegidMappingReportReqBody represents http POST request body structure
type PartnerRegidMappingReportReqBody struct {
	Vendor    string   `json:"vendor"`
	ClientID  string   `json:"clientId"`
	SiteID    string   `json:"siteId"`
	RegIDList []string `json:"regIdList"`
}

//PartnerHostidMappingReportReqBody represents http POST request body structure
type PartnerHostidMappingReportReqBody struct {
	Vendor     string   `json:"vendor"`
	ClientID   string   `json:"clientId"`
	SiteID     string   `json:"siteId"`
	HostIDList []string `json:"hostIdList"`
}

//PartnerMappingReportNode represents Details for Partner Mapping Report Record
type PartnerMappingReportNode struct {
	PartnerID          string `json:"partnerId"`
	EndpointID         string `json:"endpointId"`
	ClientID           string `json:"clientId"`
	SiteID             string `json:"siteId"`
	RegID              string `json:"regId"`
	RemoteAccessHostID string `json:"remoteAccessHostId"`
	EndpointType       string `json:"endpointType"`
}

//PartnerMappingReportData represents Details for Partner Mapping Report
type PartnerMappingReportData struct {
	OutData []PartnerMappingReportNode `json:"outData"`
}
