package kafka

import (
	"encoding/json"
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/messaging"
)

// KafkaConn  kafka repo for exec results
type KafkaConn interface {
	Publish(env *messaging.Envelope) error
}

// NewExecutionResult returns new ExecutionResult
func NewExecutionResult(topic string, kafkaConn KafkaConn) *ExecutionResult {
	return &ExecutionResult{
		topic:     topic,
		kafkaConn: kafkaConn,
	}
}

// ExecutionResult is a kafka repo client
type ExecutionResult struct {
	topic     string
	kafkaConn KafkaConn
}

// Publish publishes exec result kafka message
func (s *ExecutionResult) Publish(msg tasking.ExecutionResultKafkaMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("can't marshal message: %v", msg)
	}

	envelope := messaging.Envelope{
		Topic:   s.topic,
		Message: string(data),
	}

	if err = s.kafkaConn.Publish(&envelope); err != nil {
		return fmt.Errorf("can't publish message: %v", envelope)
	}
	return nil
}
