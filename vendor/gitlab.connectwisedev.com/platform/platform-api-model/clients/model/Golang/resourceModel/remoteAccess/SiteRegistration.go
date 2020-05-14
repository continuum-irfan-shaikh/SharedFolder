package remoteAccess

//SiteRegistration represents Details for Registered Sites
type SiteRegistration struct {
	PartnerID  int64  `json:"partnerId"`
	ClientID   int64  `json:"clientId"`
	SiteID     int64  `json:"siteId"`
	StatusCode int    `json:"statusCode"`
	Details    string `json:"details"`
}

//SiteRegistrationReqBody represents http POST request body structure
type SiteRegistrationReqBody struct {
	Vendor       string `json:"vendor"`
	ClientID     string `json:"clientId"`
	SiteID       string `json:"siteId"`
	ResourceType string `json:"resourceType"`
	RegisterFlag bool   `json:"registerFlag"`
	OneClickFlag bool   `json:"oneClickFlag"`
}
