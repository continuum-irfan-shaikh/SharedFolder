package integration

import (
	"context"

	e "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
)

//go:generate mockgen -destination=../mocks/mocks-gomock/automationEngine_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration AutomationEngine

// AutomationEngine represents interface to work with AE and policies
type AutomationEngine interface {
	UpdateRemotePolicies(ctx context.Context, data []e.TriggerDefinition) (string, error)
	RemovePolicy(ctx context.Context, systemIdentifier map[string]interface{}) (err error)
}
