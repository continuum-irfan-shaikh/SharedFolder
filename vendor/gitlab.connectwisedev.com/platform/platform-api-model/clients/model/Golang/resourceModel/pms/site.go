package pms

// Site represent site response structure for PMS
type Site struct {
	ResType           int      `json:"ResType,omitempty"`
	ResDescription    string   `json:"ResDescription,omitempty"`
	EventOriginatedAt string   `json:"EventOriginatedAt,omitempty"`
	Data              SiteData `json:"Data,omitempty"`
}

// SiteData represent site response structure for PMS
// each attribute map to mstsite table column of database
type SiteData struct {
	Site_CallerGreeting string `json:"Site_CallerGreeting,omitempty"`
	Client_Id           int    `json:"Client_Id,omitempty"`
	Site_MainPhoneno    string `json:"Site_MainPhoneno,omitempty"`
	Site_CreationDate   string `json:"Site_CreationDate,omitempty"`
	Site_DisabledOn     string `json:"Site_DisabledOn,omitempty"`
	Site_Status         bool   `json:"Site_Status,omitempty"`
	Partner_Id          int    `json:"Partner_Id,omitempty"`
	Site_Proxy          bool   `json:"Site_Proxy,omitempty"`
	Site_Address        string `json:"Site_Address,omitempty"`
	Site_Address2       string `json:"Site_Address2,omitempty"`
	Site_City           string `json:"Site_City,omitempty"`
	Site_Country        string `json:"Site_Country,omitempty"`
	Site_Id             int    `json:"Site_Id,omitempty"`
	Site_Name           string `json:"Site_Name,omitempty"`
	Site_Postalcode     string `json:"Site_Postalcode,omitempty"`
	Site_State          string `json:"Site_State,omitempty"`
	Site_Code           string `json:"Site_Code,omitempty"`
	Status              int    `json:"Status,omitempty"`
	TimeStamp           string `json:"TimeStamp,omitempty"`
	Site_Timezone       string `json:"Site_Timezone,omitempty"`
	Site_Proxy_UserName string `json:"Site_Proxy_UserName,omitempty"`
	Operation           int    `json:"Operation,omitempty"`
}
