package repository

import (
	"context"

	e "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"github.com/gocql/gocql"
)

//go:generate mockgen -destination=../mocks/mocks-gomock/triggers_mock.go  -package=mocks -source=./triggers.go

// TriggersRepo is a repo to work with triggers specific tables
type TriggersRepo interface {
	Insert(tr e.ActiveTrigger) error
	Delete(tr e.ActiveTrigger) error

	GetAll() ([]e.ActiveTrigger, error)
	GetAllByType(ctx context.Context, typeTrigger string, partnerID string, fromCache bool) ([]e.ActiveTrigger, error)
	GetAllByTaskID(partnerID string, taskID gocql.UUID) ([]e.ActiveTrigger, error)

	InsertDefinitions(defs []e.TriggerDefinition) error
	TruncateDefinitions() error

	GetDefinition(triggerType string) (e.TriggerDefinition, error)
	GetAllDefinitions() ([]e.TriggerDefinition, error)
	GetAllDefinitionsNamesAndIDs() ([]e.TriggerDefinition, error)

	GetTriggerCounterByType(triggerType string) (e.TriggerCounter, error)
	IncreaseCounter(counter e.TriggerCounter) error
	DecreaseCounter(counter e.TriggerCounter) error
}
