package ticket

// CherwellIncident represent Cherwell incident object
type CherwellIncident struct {
	NOCID          int    `json:"NOCID,omitempty"`
	NOCType        string `json:"NOCType,omitempty"`
	PartnerID      int    `json:"PartnerID,omitempty"`
	ClientID       int    `json:"ClientID,omitempty"`
	SiteID         int    `json:"SiteID,omitempty"`
	IncidentStatus string `json:"IncidentStatus,omitempty"`
	LegacyStatus   string `json:"LegacyStatus,omitempty"`
	AssignedTo     string `json:"AssignedTo,omitempty"`
	Title          string `json:"Title,omitempty"`
	Description    string `json:"Description,omitempty"`
	Source         string `json:"Source,omitempty"`
	PSATicketID    string `json:"PSATicketID,omitempty"`
	PSAType        string `json:"PSAType,omitempty"`
	Priority       string `json:"Priority,omitempty"`
	IncidentType   string `json:"IncidentType,omitempty"`
	ServiceCatalog string `json:"ServiceCatalog,omitempty"`
	CreatedBy      string `json:"CreatedBy,omitempty"`
	CreatedOn      string `json:"CreatedOn,omitempty"`
	IncidentID     string `json:"IncidentID"`
	Origin         string `json:"Origin,omitempty"`
	UpdatedOn      string `json:"UpdatedOn,omitempty"`
	ResourceID     int    `json:"ResourceID,omitempty"`
	ConditionID    int    `json:"ConditionID,omitempty"`
	PartnerCode    string `json:"PartnerCode,omitempty"`
	SiteCode       string `json:"SiteCode,omitempty"`
	Version        string `json:"Version,omitempty"`
}
