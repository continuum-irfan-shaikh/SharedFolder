package entities

import (
	"time"

	"github.com/gocql/gocql"
)

// SysEventsPayload is a payload entity that need sys-event plugin to activate triggers
type SysEventsPayload struct {
	PartnerID   string       `json:"-"` // this field is ignored and we don't need it on plugin side. needs to transfer the data
	Name        string       `json:"name"`
	PluginPath  []string     `json:"pluginPaths"`
	TriggerType string       `json:"type"`
	Shutdown    *TriggerData `json:"shutdown,omitempty"` // only one of this must be at the same time in one payload
	Logoff      *TriggerData `json:"logoff,omitempty"`   // only one of this must be at the same time in one payload
}

// TriggerData represents user tasking event data
type TriggerData struct {
	Payload      TaskMessage `json:"payload"`
	DelayTimeSec int         `json:"delayTimeSec"`
}

// TaskMessage struct describes message for Execution TaskMessage model
type TaskMessage struct {
	Task       string `json:"task"`
	ExecuteNow string `json:"executeNow"` // must be "true"
	Schedule   string `json:"schedule"`
	TaskInput  string `json:"taskInput"` // must be in base64 and comes from scripting
}

// AgentTargets struct that says where the changes will be deployed
type AgentTargets struct {
	PartnerID  string `json:"partnerID"`
	SiteID     string `json:"siteID"`
	ClientID   string `json:"clientID"`
	EndpointID string `json:"endpointID"`
}

// AgentActivateResp represents the response of Agent-config activate profile
type AgentActivateResp struct {
	ProfileID gocql.UUID `json:"id"`
}

// ActiveTrigger represents active triggers type
type ActiveTrigger struct {
	Type           string
	PartnerID      string
	TaskID         gocql.UUID
	StartTimeFrame time.Time
	EndTimeFrame   time.Time
}

// TriggerCounter entity that stores active triggers count
type TriggerCounter struct {
	TriggerID string
	PolicyID  string
	Count     int
}

// Rule is event log collector description
type Rule struct {
	TriggerID    string           `json:"-"`
	Collector    string           `json:"collector,omitempty"`
	Description  string           `json:"description,omitempty"`
	EventDetails RuleEventDetails `json:"eventDetails,omitempty"`
}

// RuleEventDetails is settings for rules
type RuleEventDetails struct {
	Facility       int      `json:"facility,omitempty"`
	Source         string   `json:"source,omitempty"`
	Channel        string   `json:"channel,omitempty"`
	EventIDs       []string `json:"eventIDs,omitempty"`
	Severity       []int    `json:"severity,omitempty"`
	FetchTime      string   `json:"fetchTime,omitempty"`
	CaptureXMLView bool     `json:"captureXMLView,omitempty"`
}

// AgentConfigPayload  is a payload entity that need Agent-config service to change sys-events configuration file
type AgentConfigPayload struct {
	Tag           string          `json:"tag"`
	Description   string          `json:"description"`
	Configuration []Configuration `json:"configurations"`
	Targets       []AgentTargets  `json:"targets"`
}

// Configuration agent config needed field that says what and where will be changed
type Configuration struct {
	PackageName string  `json:"packageName"`
	FileName    string  `json:"fileName"`
	Patch       []Patch `json:"patch"`
}

// Patch ...
type Patch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value Rule   `json:"value"`
}
