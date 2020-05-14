package entities

import (
	"strings"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/triggers"
)

// Policy is a struct that represents policy config. Needed for AE MS
type Policy struct {
	ID               string                 `json:"policyid"`
	Description      string                 `json:"description"`
	Type             string                 `json:"type"` // for AE must have 'policy' value
	Topic            string                 `json:"topic"`
	SystemIdentifier map[string]interface{} `json:"systemidentifier"` // message identifier
	Name             string                 `json:"name"`
	Version          string                 `json:"version"` // 1
	Variables        map[string]PolicyVars  `json:"variables"`
	Actions          []Action               `json:"actions"`
}

// PolicyVars represents struct with kafka message values
type PolicyVars struct {
	Type     string `json:"type"` // value is 'message'
	Key      string `json:"key"`
	DataType string `json:"dataType"` //string
}

// Action represents action section of a policy file
type Action struct {
	TriggerID     string
	Name          string            `json:"Name"`
	Mode          string            `json:"Mode"`
	Protocol      string            `json:"Protocol"` // protocol to send webhook request
	Endpoint      string            `json:"Endpoint"` // URL of the webhook service
	Method        string            `json:"Method"`
	Context       string            `json:"Context"`     // prefix of the called API
	EndResource   string            `json:"EndResource"` // postfix of an API
	PathVariables []KeyValue        `json:"PathVariables"`
	Payload       map[string]string `json:"payload"`
	Headers       Headers
}

// KeyValue is just a simple key-value struct to store string key-value. Made to store the order
type KeyValue struct {
	Key   string
	Value string
}

// Headers represents API headers
type Headers struct {
	Uid   string `json:"uid"`
	Realm string `json:"realm"`
}

// IsAlertingPolicy checks if this is alerting policy type or not
func (p *Policy) IsAlertingPolicy() bool {
	return strings.Contains(p.ID, triggers.AlertTypePrefix)
}
