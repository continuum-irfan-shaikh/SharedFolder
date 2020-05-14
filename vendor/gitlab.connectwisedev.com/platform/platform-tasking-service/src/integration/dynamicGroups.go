package integration

import (
	"context"

	"github.com/gocql/gocql"
)

//go:generate mockgen -destination=../mocks/mocks-gomock/dynamicGroupRepo_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration DynamicGroups

// DynamicGroups represents repository to work with Dynamic Groups MS
type DynamicGroups interface {
	GetEndpointsByGroupIDs(ctx context.Context, targetIDs []string, createdBy, partnerID string, hasNOCAccess bool) ([]gocql.UUID, error)

	StartMonitoringGroups(ctx context.Context,partnerID string, groupIDs []string, taskID gocql.UUID) error
	StopGroupsMonitoring(ctx context.Context, partnerID string, groupIDs []string, taskID gocql.UUID) error
}
