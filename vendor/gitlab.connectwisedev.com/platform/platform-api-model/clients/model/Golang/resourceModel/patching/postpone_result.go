package patching

import (
	"time"

	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
)

// PostponeResult postpone result from plugin listened by kafka
type PostponeResult struct {
	agent.BrokerEnvelope
	Message PostponeResultResponseMessage `json:"message"`
}

// DownloadResultResponseMessage download response message
type PostponeResultResponseMessage struct {
	TimestampUTC time.Time `json:"timestampUTC" description:"UTC time"`
	Metadata     string    `json:"metadata" description:"Postpone metadata with json schema"`
}
