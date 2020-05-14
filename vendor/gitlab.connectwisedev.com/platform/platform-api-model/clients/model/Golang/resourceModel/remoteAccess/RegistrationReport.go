package remoteAccess

//RegistrationReportReqBody represents http POST request body structure
type RegistrationReportReqBody struct {
	Vendor        string   `json:"vendor"`
	ClientID      string   `json:"clientId"`
	SiteID        string   `json:"siteId"`
	EndpointsList []string `json:"endpointsList"`
}

//RegistrationReportNode represents Details for Registration Report Record
type RegistrationReportNode struct {
	PartnerID        string `json:"partnerId"`
	EndpointID       string `json:"endpointId"`
	ClientID         string `json:"clientId"`
	SiteID           string `json:"siteId"`
	RegID            string `json:"regId"`
	CurrentPort      int    `json:"currentPort"`
	InstallBuild     string `json:"installBuild"`
	LicenseType      int    `json:"licenseType"`
	ProductBuild     int    `json:"productBuild"`
	ProductType      int    `json:"productType"`
	ProductVersion   string `json:"productVersion"`
	InstalledVersion string `json:"installedVersion"`
	LicenseInfo      string `json:"licenseInfo"`
	EndpointType     string `json:"endpointType"`
	OSLmi            int    `json:"osLmi"`
	OSSpec           int    `json:"osSpec"`
	TimezoneIndex    int    `json:"timezoneIndex"`
}

//RegistrationReportData represents Details for Registration Report
type RegistrationReportData struct {
	OutData []RegistrationReportNode `json:"outData"`
}
