package integration

import (
	"context"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
)

//go:generate mockgen -destination=../mocks/mocks-gomock/agentConfigRepo_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration AgentConfig

// AgentConfig is a repository that represents activating and deactivating triggers using external services
type AgentConfig interface {
	Activate(ctx context.Context, content entities.Rule, managedEndpointsIDs map[string]entities.Endpoints, partnerID string) (profileID gocql.UUID, err error)
	Deactivate(ctx context.Context, profileID gocql.UUID, partnerID string) error
	Update(ctx context.Context, content entities.Rule, managedEndpointsIDs map[string]entities.Endpoints, partnerID string, profileID gocql.UUID) (err error)
}
