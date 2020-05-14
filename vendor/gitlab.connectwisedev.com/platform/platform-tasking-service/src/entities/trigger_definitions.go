package entities

// TriggerDefinition defines struct that describes trigger config definition
// must be synced with https://github.com/ContinuumLLC/rmm-task-triggers/blob/master/schemas/trigger_schema.json
type TriggerDefinition struct {
	ID              string       `json:"id"`
	Internal        bool         `json:"internal"`
	DisplayName     string       `json:"displayName"`
	Description     string       `json:"description"`
	TriggerCategory string       `json:"triggerCategory"`
	EventDetails    EventDetails `json:"eventDetails"`
	SourceSystem    string       `json:"sourceSystem,omitempty"` // field that is used to delete policies from AE
}

// EventDetails is a struct that represents event details of a trigger with info about kafka topics and ide
type EventDetails struct {
	Topic               string                 `json:"topic"`
	MessageIdentifier   map[string]interface{} `json:"messageIdentifier"`
	EndpointIdentifiers EndpointIdentifiers    `json:"endpointIdentifiers"`
	PayloadIdentifiers  map[string]Key         `json:"payloadIdentifiers,omitempty"`
}

// EndpointIdentifiers represents identifier that is used in API endpoint for tasking webhook
type EndpointIdentifiers struct {
	PartnerID  Key `json:"partnerID"`
	ClientID   Key `json:"clientID"`
	SiteID     Key `json:"siteID"`
	EndpointID Key `json:"endpointID"`
}

// Key is a simple struct that contains value of a key
type Key struct {
	Value       string `json:"key"`
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
}
