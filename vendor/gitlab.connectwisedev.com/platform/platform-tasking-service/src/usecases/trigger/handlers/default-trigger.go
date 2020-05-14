package handlers

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	m "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository"
)

const defaultTriggerUseCaseLogPrefix = "usecases/trigger/handler/default-trigger.go: "

// DefaultTriggerHandler represent default handler for triggers that don't need external connections to activate trigger
type DefaultTriggerHandler struct {
	taskRepo     m.TaskPersistence
	triggersRepo repository.TriggersRepo
	log          logger.Logger
}

// NewDefaultTrigger returns DefaultTriggerHandler
func NewDefaultTrigger(tRepo m.TaskPersistence, tr repository.TriggersRepo, log logger.Logger) *DefaultTriggerHandler {
	return &DefaultTriggerHandler{
		taskRepo:     tRepo,
		triggersRepo: tr,
		log:          log,
	}
}

// Activate implements standard trigger handlers behaviour
func (tr *DefaultTriggerHandler) Activate(_ context.Context, _ string, _ []m.Task) error {
	return nil
}

// Deactivate implements standard trigger handlers behaviour
func (tr *DefaultTriggerHandler) Deactivate(_ context.Context,_ string, _ []m.Task) error {
	return nil
}

// Update implements standard trigger handlers behaviour
func (tr *DefaultTriggerHandler) Update(_ context.Context,_ string, _ []m.Task) error {
	return nil
}

// GetTask returns task to process
func (tr *DefaultTriggerHandler) GetTask(ctx context.Context, taskID gocql.UUID) (task m.Task, err error) {
	var (
		partnerID  = ctx.Value(config.PartnerIDKeyCTX).(string)
		endpointID = ctx.Value(config.EndpointIDKeyCTX).(gocql.UUID)
	)

	internalTasks, err := tr.taskRepo.GetByIDAndManagedEndpoints(ctx, partnerID, taskID, endpointID)
	if err != nil {
		switch err.(type) {
		case m.TaskNotFoundError:
			return task, fmt.Errorf("no tasks for endpoint %v and taskId %v", endpointID, taskID)
		default:
			tr.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, defaultTriggerUseCaseLogPrefix+canNotGetInternalTasksErrMsg, taskID, endpointID, err)
			return task, err
		}
	}

	if len(internalTasks) == 0 {
		return task, fmt.Errorf("no tasks for endpoint %v and taskId %v", endpointID, taskID)
	}

	return internalTasks[0], nil
}

// IsApplicable noting to check on default trigger
func (tr *DefaultTriggerHandler) IsApplicable(_ context.Context, _ m.Task, _ apiModels.TriggerExecutionPayload) bool {
	return true // no special behaviour for default one
}

// PreExecution ..
func (tr *DefaultTriggerHandler) PreExecution(_ m.Task) error {
	return nil // no special behaviour for default one
}

// PostExecution ..
func (tr *DefaultTriggerHandler) PostExecution(_ m.Task) error {
	return nil // no special behaviour for default one
}
