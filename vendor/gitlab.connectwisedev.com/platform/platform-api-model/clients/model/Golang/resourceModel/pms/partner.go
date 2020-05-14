package pms

// Partner represent partner response structure for PMS
type Partner struct {
	ResType           int         `json:"ResType,omitempty"`
	ResDescription    string      `json:"ResDescription,omitempty"`
	EventOriginatedAt string      `json:"EventOriginatedAt,omitempty"`
	Data              PartnerData `json:"Data,omitempty"`
}

// PartnerData represent partner response structure for PMS
// each attribute map to mstmember table column of database
type PartnerData struct {
	Partner_ActivatedOn             string `json:"Partner_ActivatedOn,omitempty"`
	Partner_Address                 string `json:"Partner_Address,omitempty"`
	Partner_City                    string `json:"Partner_City,omitempty"`
	Partner_Country                 string `json:"Partner_Country,omitempty"`
	Partner_DisabledOn              string `json:"Partner_DisabledOn,omitempty"`
	Partner_StartDate               string `json:"Partner_StartDate,omitempty"`
	Partner_EmailId                 string `json:"Partner_EmailId,omitempty"`
	Partner_Freezed_EffectiveDate   string `json:"Partner_Freezed_EffectiveDate,omitempty"`
	Partner_Freezed_Reason          string `json:"Partner_Freezed_Reason,omitempty"`
	Partner_HDRegion                string `json:"Partner_HDRegion,omitempty"`
	Partner_Freezed                 bool   `json:"Partner_Freezed,omitempty"`
	Partner_Currently_Active        bool   `json:"Partner_Currently_Active,omitempty"`
	Partner_Code                    string `json:"Partner_Code,omitempty"`
	Partner_Id                      int    `json:"Partner_Id,omitempty"`
	Partner_Name                    string `json:"Partner_Name,omitempty"`
	Partner_Mobileno                string `json:"Partner_Mobileno,omitempty"`
	Partner_SalesforceId            string `json:"Partner_SalesforceId,omitempty"`
	Partner_State                   string `json:"Partner_State,omitempty"`
	Status                          int    `json:"Status,omitempty"`
	Partner_TelNo                   string `json:"Partner_TelNo,omitempty"`
	TimeStamp                       string `json:"TimeStamp,omitempty"`
	Partner_UnFreezed_EffectiveDate string `json:"Partner_UnFreezed_EffectiveDate,omitempty"`
	Partner_ZipCode                 string `json:"Partner_ZipCode,omitempty"`
	Operation                       int    `json:"Operation,omitempty"`
	Partner_Currency                string `json:"Partner_Currency,omitempty"`
}
