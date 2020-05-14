package health

import "time"

//Action is the recommended action by AgentCore for watchdog
type Action string

//PluginName ...
type PluginName string

const (
	//Restart for restarting ITSPlatform service
	Restart Action = "Restart"
)

//AgentHealth returns health of the agent
type AgentHealth struct {
	TimestampUTC        time.Time                         `json:"timestampUTC"`
	AgentCore           AgentCore                         `json:"agentCore"`
	ShortRunningPlugins map[PluginName]ShortRunningPlugin `json:"shortRunningPlugins"`
	LongRunningPlugins  map[PluginName]LongRunningPlugin  `json:"longRunningPlugins"`
	ServiceDetail       ServiceDetail                     `json:"serviceDetail"`
	RecommendedAction   Action                            `json:"recommendedAction"`
}

//AgentCore ...
type AgentCore struct {
	TimestampUTC           time.Time              `json:"timestampUTC"`
	State                  State                  `json:"state"`
	HeartbeatStatus        HeartbeatStatus        `json:"heartbeatStatus"`
	PendingMailboxMessages PendingMailboxMessages `json:"pendingMailboxMessages"`
	OfflineMessages        OfflineMessages        `json:"offlineMessages"`
}

//State is a struct to hold Agent State
type State struct {
	TimestampUTC    time.Time `json:"timestampUTC"`
	Online          bool      `json:"online"`
	ActiveBroker    bool      `json:"activeBroker"`
	ActiveHeartbeat bool      `json:"activeHeartbeat"`
}

//HeartbeatStatus is a struct to show detailed HB stats
type HeartbeatStatus struct {
	TimestampUTC   time.Time `json:"timestampUTC"`
	TotalAttempts  int64     `json:"totalAttempts"`
	Success        int64     `json:"success"`
	NetworkFailure int64     `json:"networkFailure"`
	OtherFailure   int64     `json:"otherFailure"`
}

//PendingMailboxMessages is a struct to hold Agent State
type PendingMailboxMessages struct {
	TimestampUTC time.Time `json:"timestampUTC"`
	Count        int64     `json:"count"`
}

//OfflineMessages is a struct to hold Agent State
type OfflineMessages struct {
	TimestampUTC time.Time `json:"timestampUTC"`
	Count        int64     `json:"count"`
}

//ShortRunningPlugin ...
type ShortRunningPlugin struct {
	TotalInvocation      int64 `json:"totalInvocation"`
	SuccessfulInvocation int64 `json:"successfulInvocation"`
	FailedInvocation     int64 `json:"failedInvocation"`
	OtherFailures        int64 `json:"otherFailures"`
}

//LongRunningPlugin ...
type LongRunningPlugin struct{}

//ServiceDetail ...
type ServiceDetail struct {
	TimestampUTC            time.Time `json:"timestampUTC"`
	NumberOfServiceRestarts int64     `json:"numberOfServiceRestarts"`
}
