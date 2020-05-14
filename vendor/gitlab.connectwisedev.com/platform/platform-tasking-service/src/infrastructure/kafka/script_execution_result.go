package kafka

import "gitlab.connectwisedev.com/platform/platform-common-lib/src/messaging"

// NewExecutionResult represents new messaging exec results service
func NewExecutionResult(config messaging.Config) messaging.Service {
	return messaging.NewUniqueService(config)
}
