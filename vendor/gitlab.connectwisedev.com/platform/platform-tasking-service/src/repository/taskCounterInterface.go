package repository

import (
	"context"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
)

// TaskCounter interface to perform actions with task_counters table
type TaskCounter interface {
	GetCounters(ctx context.Context, partnerID string, endpointID gocql.UUID) ([]models.TaskCount, error)
	IncreaseCounter(partnerID string, counters []models.TaskCount, isExternal bool) error
	DecreaseCounter(partnerID string, counters []models.TaskCount, isExternal bool) error

	GetAllPartners(ctx context.Context) (partnerIDs map[string]struct{}, err error)
}
