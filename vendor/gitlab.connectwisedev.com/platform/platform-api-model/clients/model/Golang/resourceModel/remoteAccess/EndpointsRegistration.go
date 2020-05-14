package remoteAccess

//EndpointsRegistration represents Details for Registered Endpoints in Site
type EndpointsRegistration struct {
	PartnerID int64              `json:"partnerId"`
	ClientID  int64              `json:"clientId"`
	SiteID    int64              `json:"siteId"`
	OutData   []RegistrationData `json:"outData"`
}

//RegisterEndpointsList represents list of Endpoints to Register
type RegisterEndpointsList struct {
	EndpointID   string `json:"endpointId"`
	ResourceType string `json:"resourceType"`
	RegisterFlag bool   `json:"registerFlag"`
	OneClickFlag bool   `json:"oneClickFlag"`
}

//EndpointsRegistrationReqBody represents http POST request body structure
type EndpointsRegistrationReqBody struct {
	Vendor        string                  `json:"vendor"`
	ClientID      string                  `json:"clientId"`
	SiteID        string                  `json:"siteId"`
	EndpointsList []RegisterEndpointsList `json:"endpointsList"`
}
