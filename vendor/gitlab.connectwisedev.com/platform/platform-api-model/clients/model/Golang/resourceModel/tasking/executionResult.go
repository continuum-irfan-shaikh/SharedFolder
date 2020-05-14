package tasking

import (
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	"github.com/gocql/gocql"

	"time"
)

// ExecutionResult represents result details returned by Origin MS to Tasking MS
// Tasking POST API route /partners/{partnerID}/task-execution-results/task-instances/{taskInstanceID}
type ExecutionResult struct {
	// Possible values: "Success", "Failed", "Some Failures"
	CompletionStatus string    `json:"completionStatus"  valid:"required"`
	EndpointID       string    `json:"endpointId"        valid:"required,uuid"`
	UpdateTime       time.Time `json:"updateTime"        valid:"required"`
	ErrorDetails     string    `json:"errorDetails"      valid:"-"`
	ResultDetails    string    `json:"resultDetails"     valid:"-"`
}

// ExecutionResultKafkaMessage structure represents the Kafka message with Script execution results
type ExecutionResultKafkaMessage struct {
	agent.BrokerEnvelope
	Message ScriptPluginReturnMessage `json:"message"`
}

// TaskResult is Task result for result webhook
type TaskResult struct {
	ID            gocql.UUID `json:"id"`
	ResultMessage string     `json:"result_message"`
	StdOut        string     `json:"std_out"`
	StdErr        string     `json:"std_err"`
	Success       bool       `json:"success"`
	ResultWebhook string     `json:"result_webhook"`
}