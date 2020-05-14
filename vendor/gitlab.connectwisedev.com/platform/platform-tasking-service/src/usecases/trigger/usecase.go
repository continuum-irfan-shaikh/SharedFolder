package trigger

import (
	"context"

	api "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	e "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	m "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
)

//go:generate mockgen -destination=../../mocks/mocks-gomock/triggerUC_mock.go  -package=mocks -source=./usecase.go

// Usecase - represents trigger usecases
type Usecase interface {
	Activate(ctx context.Context, tasks []m.Task) error
	Deactivate(ctx context.Context, tasks []m.Task) error

	PreExecution(triggerType string, task m.Task) error
	PostExecution(triggerType string, task m.Task) error

	IsApplicable(ctx context.Context, task m.Task, payload api.TriggerExecutionPayload) bool

	GetActiveTriggers(ctx context.Context) ([]e.ActiveTrigger, error)
	GetActiveTriggersByTaskID(partnerID string, taskID gocql.UUID) ([]e.ActiveTrigger, error)
	GetTask(ctx context.Context, taskID gocql.UUID) (m.Task, error)

	DeleteActiveTriggers(triggers []e.ActiveTrigger) error

	ActiveTriggersReopening(ctx context.Context)
}

// Handler interface for handlers of different types of trigger
type Handler interface {
	Activate(ctx context.Context, exactType string, tasks []m.Task) error
	Deactivate(ctx context.Context, exactType string, tasks []m.Task) error
	Update(ctx context.Context,exactType string, tasks []m.Task) error

	IsApplicable(ctx context.Context, task m.Task, payload api.TriggerExecutionPayload) bool

	GetTask(ctx context.Context, taskID gocql.UUID) (m.Task, error)
	PreExecution(task m.Task) error
	PostExecution(task m.Task) error
}

// DefinitionUseCase represents interface to import externally defined triggers
type DefinitionUseCase interface {
	ImportExternalTriggers(defs []e.TriggerDefinition) error
	GetTriggerTypes() (triggerTypes []e.TriggerDefinition, err error)
}
