package err

import "time"

//StartupFailure is the struct definition of Error message sent by Agent Manager
type StartupFailure struct {
	AgentFailTSUTC  time.Time `json:"agentFailTsUTC"`
	AgentFailRemark string    `json:"agentFailRemark"`
	FailureSource   string    `json:"failureSource"`
}

//Message is the struct definition of Error message structure
type Message struct {
	EndpointID     string              `json:"endpointID,omitempty"`
	AgentID        string              `json:"agentID,omitempty"`
	PartnerID      string              `json:"partnerID,omitempty"`
	SiteID         string              `json:"siteID,omitempty"`
	ClientID       string              `json:"clientID,omitempty"`
	LegacyRegID    string              `json:"legacyRegID,omitempty"`
	TimestampUTC   time.Time           `json:"timestampUTC,omitempty"`
	DCTimestampUTC time.Time           `json:"dcTimestampUTC,omitempty"`
	Source         string              `json:"source,omitempty"`
	Type           string              `json:"type,omitempty"`
	SubType        string              `json:"subType,omitempty"`
	ErrorCode      string              `json:"errorCode,omitempty"`
	ErrorTrace     string              `json:"errorTrace,omitempty"`
	Metadata       map[string][]string `json:"metadata,omitempty"`
	TransactionID  string              `json:"transactionID,omitempty"`
}

const (
	//StatusCode is a Key for Status Code
	StatusCode = "StatusCode"

	//AutoUpdate is a type used for Autoupdate messages
	AutoUpdate = "Auto Update"

	//AgentOffline is a type used for Agent Offline messages
	AgentOffline = "Agent Offline"

	//AgentStartupFailure is a type used for Agent Startup Failure messages
	AgentStartupFailure = "Agent Startup Failure"

	//AgentPluginFailure is a type used for Agent Plugin Failure messages
	AgentPluginFailure = "Agent Plugin Failure"

	//AgentCoreRecovery is a type used for Agent Recovery messages
	AgentCoreRecovery = "Agent-Core-Recovery"

	//JunoManagerReporting is a type used for Juno manager reporting  messages
	JunoManagerReporting = "Juno-Manager-Reporting"

	//AgentCoreRecoveryType is a type used for Agent Recovery Type
	AgentCoreRecoveryType = "Recovery"

	//ReportingType is a type used for Juno manager reporting  messages
	ReportingType = "Reporting"

	//DefaultPartnerID is a type used for DefaultPartnerID
	DefaultPartnerID = "NotExists"

	//Failed to execute the run as script
	RunAsScriptFailure = "Run-As-Script-Failure"

	//Failed to execute the gateway pluin
	GatewayPluginFailure = "Gateway-Plugin-Failure"

	//Failed to download packages from Gateway
	GatewayDownloadFaliure = "Gateway-Download-Failure"

	//PluginValidationError used to indicate a type to be used in case of plugin validation failure scenario
	PluginValidationError = "Plugin-Validation-Error"
)

//AgentFailureSource ...
var AgentFailureSource = map[string]bool{
	AgentCoreRecovery:    true,
	JunoManagerReporting: true,
}

//AgentFailureType ...
var AgentFailureType = map[string]bool{
	AgentCoreRecoveryType: true,
	ReportingType:         true,
}
