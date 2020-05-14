package agent

import "time"

//ErrorMessage is the struct definition of Error message structure
type ErrorMessage struct {
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Version      string            `json:"version"`
	TimeUUID     string            `json:"timeUUID"`
	TimestampUTC time.Time         `json:"timestampUTC"`
	Path         string            `json:"path"`
	ErrorTrace   string            `json:"errorTrace"`
	StatusCode   int               `json:"statusCode"`
	ErrorCode    string            `json:"errorCode"`
	ErrorData    map[string]string `json:"errorData"`
}

const (
	//ErrBlankMessageRecieved error Blank Message Recieved
	ErrBlankMessageRecieved = "ErrBlankMessageRecieved"

	//ErrAgentMappingNotFound when HB received for unavailability of Mapping for Agents
	ErrAgentMappingNotFound = "ErrAgentMappingNotFound"

	//ErrCircuitOpen thrown when circuit breaker trips
	ErrCircuitOpen = "ErrCircuitOpen"

	//ErrMaxConcurrency thrown when max concurrency reached
	ErrMaxConcurrency = "ErrMaxConcurrency"

	//ErrTimeout thrown when broker/Cassandra command times out
	ErrTimeout = "ErrTimeout"

	//ErrNoError indicates no error
	ErrNoError = "ErrNoError"

	//OfflineMessageSeparator separates different messages in a batch
	OfflineMessageSeparator = "Me$@"
)
