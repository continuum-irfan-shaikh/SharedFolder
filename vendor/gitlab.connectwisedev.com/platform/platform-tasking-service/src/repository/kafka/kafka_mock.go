package kafka

import (
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/messaging"
)

// NewExecutionResultMock returns new ExecutionResult
func NewExecutionResultMock(withErr bool) *ExecutionResultMock {
	return &ExecutionResultMock{
		withErr: withErr,
	}
}

// ExecutionResultMock is a kafka repo client
type ExecutionResultMock struct {
	withErr bool
}

// Publish publishes exec result kafka message
func (s *ExecutionResultMock) Publish(env *messaging.Envelope) error {
	if s.withErr {
		return fmt.Errorf("err")
	}
	return nil
}
