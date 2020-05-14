package mockusecases

import (
	"context"
	"errors"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

const (
	//HandlerErrorActivate err for mock
	HandlerErrorActivate = "needErrorActivate"
	//HandlerErrorDeactivate err for mock
	HandlerErrorDeactivate = "needErrorDeactivate"
	//HandlerErrorUpdate err for mock
	HandlerErrorUpdate = "needErrorUpdate"
	//HandlerErrorIsApplicable err for mock
	HandlerErrorIsApplicable = "needErrorIsApplicable"
	//HandlerErrorGetTask err for mock get task
	HandlerErrorGetTask = "needErrorGetTask"
	//HandlerErrorPreExecution err for mock pre exec
	HandlerErrorPreExecution = "needErrorGetTaskPreExecution"
	//HandlerErrorPostExecution err for mock post exec
	HandlerErrorPostExecution = "needErrorGetTaskPostExecution"
)

// MockTriggerHandler implements mock for handler
type MockTriggerHandler struct{}

// Activate ..
func (m MockTriggerHandler) Activate(ctx context.Context, exactType string, tasks []models.Task) error {
	if len(tasks) == 0 {
		return nil
	}
	if tasks[0].Parameters == HandlerErrorActivate {
		return errors.New("err")
	}
	return nil
}

// Update ..
func (m MockTriggerHandler) Update(ctx context.Context, exactType string, tasks []models.Task) error {
	if tasks[0].Parameters == HandlerErrorUpdate {
		return errors.New("err")
	}
	return nil
}

//Deactivate ..
func (m MockTriggerHandler) Deactivate(ctx context.Context, exactType string, tasks []models.Task) error {
	if len(tasks) == 0 {
		return nil
	}
	if tasks[0].Parameters == HandlerErrorDeactivate {
		return errors.New("err")
	}
	return nil
}

//IsApplicable ..
func (m MockTriggerHandler) IsApplicable(ctx context.Context, task models.Task, payload tasking.TriggerExecutionPayload) bool {
	if task.Parameters == HandlerErrorIsApplicable {
		return false
	}
	return true
}

//GetTask ..
func (m MockTriggerHandler) GetTask(ctx context.Context, taskID gocql.UUID) (models.Task, error) {
	return models.Task{}, nil
}

//PreExecution ..
func (m MockTriggerHandler) PreExecution(task models.Task) error {
	if task.Parameters == HandlerErrorPreExecution {
		return errors.New("err")
	}
	return nil
}

//PostExecution ..
func (m MockTriggerHandler) PostExecution(task models.Task) error {
	if task.Parameters == HandlerErrorPostExecution {
		return errors.New("err")
	}
	return nil
}
