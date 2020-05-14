package pms

// Client represent client response structure for PMS
type Client struct {
	ResType           int        `json:"ResType,omitempty"`
	ResDescription    string     `json:"ResDescription,omitempty"`
	EventOriginatedAt string     `json:"EventOriginatedAt,omitempty"`
	Data              ClientData `json:"Data,omitempty"`
}

// ClientData represent client response structure for PMS
// each attribute map to mstsite table as in legacy
// system client concept does not exist column of database
type ClientData struct {
	Client_Id             int    `json:"Client_Id,omitempty"`
	Client_CallerGreeting string `json:"Client_CallerGreeting,omitempty"`
	Client_PrimaryPhoneNo string `json:"Client_PrimaryPhoneNo,omitempty"`
	Client_CreatedDate    string `json:"Client_CreatedDate,omitempty"`
	Client_DisabledOn     string `json:"Client_DisabledOn,omitempty"`
	Client_Status         bool   `json:"Client_Status,omitempty"`
	Partner_Id            int    `json:"Partner_Id,omitempty"`
	Client_Proxy          bool   `json:"Client_Proxy,omitempty"`
	Client_Address        string `json:"Client_Address,omitempty"`
	Client_Address2       string `json:"Client_Address2,omitempty"`
	Client_City           string `json:"Client_City,omitempty"`
	Client_Country        string `json:"Client_Country,omitempty"`
	Client_Name           string `json:"Client_Name,omitempty"`
	Client_PostalCode     string `json:"Client_PostalCode,omitempty"`
	Client_State          string `json:"Client_State,omitempty"`
	Client_Code           string `json:"Client_Code,omitempty"`
	Status                int    `json:"Status,omitempty"`
	TimeStamp             string `json:"TimeStamp,omitempty"`
	Client_Timezone       string `json:"Client_Timezone,omitempty"`
	Client_Proxy_UserName string `json:"Client_Proxy_UserName,omitempty"`
	Operation             int    `json:"Operation,omitempty"`
}
