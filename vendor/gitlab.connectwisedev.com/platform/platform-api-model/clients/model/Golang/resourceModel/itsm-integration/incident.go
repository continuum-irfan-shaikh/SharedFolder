package itsm_integration

import "time"

// Incident represents incident object for Cherwell
type Incident struct {
	ID string `json:"id,omitempty"`

	// Incident title
	Title string `json:"title,omitempty"`

	// Incident description
	Description string `json:"description,omitempty"`

	// Customer ID
	CustomerID string `json:"customerId,omitempty"`

	// Partner ID
	PartnerID string `json:"partnerId,omitempty"`

	// Incident priority (1-5)
	Priority int `json:"priority,omitempty"`

	// Status is incident's status
	Status string `json:"status,omitempty"`

	// Incident creation date in RFC3339 format
	CreatedAt string `json:"createdAt,omitempty"`

	// Assignee
	Assignee string `json:"assignee,omitempty"`

	// SiteID
	SiteID string `json:"siteId,omitempty"`

	// ClientID
	ClientID string `json:"clientId,omitempty"`

	// EndpointID - configItemRecID field in CSM
	EndpointID string `json:"endpointId,omitempty"`

	LegacyEndpointID string `json:"legacyRegId,omitempty"`

	Source string `json:"source,omitempty"`

	ConditionID string `json:"conditionId,omitempty"`

	IncidentID string `json:"incidentId,omitempty"`

	IncidentType string `json:"incidentType,omitempty"`

	AlertDetails string `json:"alertDetails,omitempty"`

	// CustomerEmail - CustomerEmail field in CSM
	CustomerEmail string `json:"customerEmail,omitempty"`

	PSAContactID string `json:"PSAContactID,omitempty"`

	Category string `json:"category,omitempty"`

	WorkGroup string `json:"workGroup,omitempty"`

	Service string `json:"service,omitempty"`

	Origin string `json:"origin,omitempty"`

	PSAIncidentID string `json:"PSAIncidentID,omitempty"`

	PSAType string `json:"PSAType,omitempty"`
}

// GetCreationDate returns incident creation Time object
func (i *Incident) GetCreationDate() (time.Time, error) {
	return time.Parse(time.RFC3339, i.CreatedAt)
}
