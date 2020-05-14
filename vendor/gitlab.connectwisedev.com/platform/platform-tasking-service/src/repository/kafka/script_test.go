package kafka

import (
	"testing"

	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
)

func TestName(t *testing.T) {
	mock := NewExecutionResultMock(false)

	service := NewExecutionResult("", mock)
	service.Publish(tasking.ExecutionResultKafkaMessage{})

	mock = NewExecutionResultMock(true)

	service = NewExecutionResult("", mock)
	service.Publish(tasking.ExecutionResultKafkaMessage{})
}
