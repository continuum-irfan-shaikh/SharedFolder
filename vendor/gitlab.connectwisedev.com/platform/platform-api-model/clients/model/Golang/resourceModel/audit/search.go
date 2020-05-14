package audit

import "time"

// Search - Search is an data structure used for searching from Audit records
type Search struct {
	First       int       `json:"first,omitempty"`
	Size        int       `json:"size,omitempty"`
	FreeText    string    `json:"freeText,omitempty"`
	PartnerID   string    `json:"partnerID,omitempty"`
	ClientID    string    `json:"clientID,omitempty"`
	SiteID      string    `json:"siteID,omitempty"`
	EndpointID  string    `json:"endpointID,omitempty"`
	ObjectName  string    `json:"objectName,omitempty"`
	EventType   string    `json:"eventType,omitempty"`
	ServiceName string    `json:"serviceName,omitempty"`
	Timeframe   Timeframe `json:"timeframe,omitempty"`
}

// Timeframe - Timeframe used for searching from Audit recodrs
type Timeframe struct {
	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`
}
