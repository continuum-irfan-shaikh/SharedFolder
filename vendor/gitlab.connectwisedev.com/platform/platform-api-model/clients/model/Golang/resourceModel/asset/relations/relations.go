package relations

import "time"

// ParentObject ...
type ParentObject struct {
	PartnerID       string `json:"partnerId,omitempty"`
	ClientID        string `json:"clientId,omitempty"`
	SiteID          string `json:"siteId,omitempty"`
	EndpointID      string `json:"endpointId,omitempty"`
	ProbeEndpointID string `json:"probeEndpointId,omitempty"`

	UUID string `json:"uuid,omitempty"`
	GUID string `json:"guid,omitempty"`

	HostName string `json:"hostname,omitempty"`
	Version  string `json:"version,omitempty"`
	Type     string `json:"type,omitempty"`
	State    string `json:"state,omitempty"`

	HasAgent    bool          `json:"hasAgent,omitempty"`
	IsMonitored bool          `json:"isMonitored,omitempty"`
	Networks    []NetworkType `json:"networks,omitempty"`

	EventType         string    `json:"eventType,omitempty"`
	DCtimeStampUTC    time.Time `json:"dctimeStampUTC"`
	EventTimeStampUTC time.Time `json:"eventTimeStampUTC"`
}

// ChildObject - Child Object
type ChildObject struct {
	PartnerID  string `json:"partnerId,omitempty"`
	ClientID   string `json:"clientId,omitempty"`
	SiteID     string `json:"siteId,omitempty"`
	EndpointID string `json:"endpointId,omitempty"`

	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	OS       string `json:"os,omitempty"`
	HostName string `json:"hostname,omitempty"`
	State    string `json:"state,omitempty"`

	UUID        string `json:"uuid,omitempty"`
	HasAgent    bool   `json:"hasAgent,omitempty"`
	IsMonitored bool   `json:"isMonitored,omitempty"`

	Networks []NetworkType `json:"networks,omitempty"`
}

// Relations ...
type Relations struct {
	*ParentObject
	Parent *Relations
}

// AssetRelations ...
type AssetRelations struct {
	*ChildObject
	Parent *Relations
}

// ChildDetails - Child Details
type ChildDetails struct {
	*ParentObject
	Childs []ChildRelations `json:"childs,omitempty"`
}

// ChildRelations - Child Relations
type ChildRelations struct {
	*ChildObject
	Childs []ChildRelations `json:"childs,omitempty"`
}
