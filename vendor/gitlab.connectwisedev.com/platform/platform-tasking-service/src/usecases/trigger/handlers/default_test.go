package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

func TestDefaultTriggerHandler_GetTask(t *testing.T) {
	var (
		ctrl                = gomock.NewController(t)
		taskPersistanceMock = mocks.NewMockTaskPersistence(ctrl)
		handler             = DefaultTriggerHandler{
			taskRepo:     taskPersistanceMock,
			triggersRepo: nil,
			log:          logger.Log,
		}

		partnerID     = "pid"
		taskID, _     = gocql.RandomUUID()
		endpointID, _ = gocql.RandomUUID()
		ctx           = context.WithValue(
			context.WithValue(
				context.Background(), config.PartnerIDKeyCTX, partnerID),
			config.EndpointIDKeyCTX, endpointID)
	)

	t.Run("positive", func(t *testing.T) {
		internalTasks := []models.Task{{}, {}}
		taskPersistanceMock.EXPECT().GetByIDAndManagedEndpoints(ctx, partnerID, taskID, endpointID).Return(internalTasks, nil)
		_, err := handler.GetTask(ctx, taskID)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		internalTasks := []models.Task{}
		taskPersistanceMock.EXPECT().GetByIDAndManagedEndpoints(ctx, partnerID, taskID, endpointID).Return(internalTasks, nil)
		_, err := handler.GetTask(ctx, taskID)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		internalTasks := []models.Task{}
		taskPersistanceMock.EXPECT().GetByIDAndManagedEndpoints(ctx, partnerID, taskID, endpointID).Return(internalTasks, errors.New(""))
		_, err := handler.GetTask(ctx, taskID)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})

	t.Run("negative_3", func(t *testing.T) {
		internalTasks := []models.Task{}
		taskPersistanceMock.EXPECT().GetByIDAndManagedEndpoints(ctx, partnerID, taskID, endpointID).Return(internalTasks, models.TaskNotFoundError{})
		_, err := handler.GetTask(ctx, taskID)
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})
}

func TestName(t *testing.T) {
	tr := NewDefaultTrigger(nil, nil, nil)
	tasks := []models.Task{}
	trType := ""
	ctx := context.TODO()
	tr.Activate(ctx, trType, tasks)
	tr.Deactivate(ctx, trType, tasks)
	tr.Update(ctx, trType, tasks)
	tr.IsApplicable(context.Background(), models.Task{}, tasking.TriggerExecutionPayload{})
	tr.PostExecution(models.Task{})
	tr.PreExecution(models.Task{})
}
