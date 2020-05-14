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
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mockLoggerTasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

func init() {
	logger.Load(config.Config.Log)
}

func TestNewFirstCheckIn(t *testing.T) {
	var (
		ctrl          = gomock.NewController(t)
		mockLogger    = mockLoggerTasking.NewMockLogger(ctrl)
		mockSitesRepo = mocks.NewMockUserSitesPersistence(ctrl)
		mockTaskRepo  = mock.NewMockTaskPersistence(ctrl)
	)

	NewFirstCheckIn(mockTaskRepo, mockLogger, mockSitesRepo)
}

func TestFirstCheckInHandler_GetTask(t *testing.T) {
	var (
		ctrl         *gomock.Controller
		mockTaskRepo *mock.MockTaskPersistence
		handler      FirstCheckInHandler

		taskID, _ = gocql.RandomUUID()
		partnerID = "pid"
		ctx       = context.WithValue(context.Background(), config.PartnerIDKeyCTX, partnerID)
	)

	before := func() {
		ctrl = gomock.NewController(t)
		mockTaskRepo = mock.NewMockTaskPersistence(ctrl)
		handler = FirstCheckInHandler{
			DefaultTriggerHandler: &DefaultTriggerHandler{
				taskRepo: mockTaskRepo,
				log:      logger.Log,
			},
		}
	}

	t.Run("positive", func(t *testing.T) {
		before()
		internalTasks := []models.Task{{}, {}}
		mockTaskRepo.EXPECT().GetByIDs(ctx, nil, partnerID, false, taskID).Return(internalTasks, nil)

		_, err := handler.GetTask(ctx, taskID)
		ctrl.Finish()
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		before()
		internalTasks := []models.Task{{}, {}}
		mockTaskRepo.EXPECT().GetByIDs(ctx, nil, partnerID, false, taskID).Return(internalTasks, errors.New(""))

		_, err := handler.GetTask(ctx, taskID)
		ctrl.Finish()
		if err == nil {
			t.Fatalf("error shouldn't be nil")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		before()
		internalTasks := []models.Task{{}, {}}
		mockTaskRepo.EXPECT().GetByIDs(ctx, nil, partnerID, false, taskID).Return(internalTasks, models.TaskNotFoundError{})

		_, err := handler.GetTask(ctx, taskID)
		ctrl.Finish()
		if err == nil {
			t.Fatalf("error shouldn't be nil")
		}
	})

	t.Run("negative_3", func(t *testing.T) {
		before()
		internalTasks := []models.Task{}
		mockTaskRepo.EXPECT().GetByIDs(ctx, nil, partnerID, false, taskID).Return(internalTasks, nil)

		_, err := handler.GetTask(ctx, taskID)
		ctrl.Finish()
		if err == nil {
			t.Fatalf("error shouldn't be nil")
		}
	})
}

func TestFirstCheckInHandler_IsApplicable(t *testing.T) {
	var (
		ctrl       *gomock.Controller
		handler    FirstCheckInHandler
		siteID = "sid"
		ctx    = context.WithValue(context.Background(), config.SiteIDKeyCTX, siteID)
	)

	before := func() {
		ctrl = gomock.NewController(t)
		handler = FirstCheckInHandler{
			DefaultTriggerHandler: &DefaultTriggerHandler{
				log: logger.Log,
			},
		}
	}

	t.Run("positive", func(t *testing.T) {
		before()
		task := models.Task{
			TargetsByType: models.TargetsByType{
				models.Site: []string{siteID},
			},
		}

		b := handler.IsApplicable(ctx, task, tasking.TriggerExecutionPayload{})
		ctrl.Finish()
		if !b {
			t.Fatalf("result should be true")
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		before()

		b := handler.IsApplicable(context.WithValue(context.Background(), config.SiteIDKeyCTX, 100), models.Task{}, tasking.TriggerExecutionPayload{})
		ctrl.Finish()
		if b {
			t.Fatalf("result should be false")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		before()
		task := models.Task{
			Targets: models.Target{
				IDs:  []string{""},
				Type: 0,
			},
		}

		b := handler.IsApplicable(ctx, task, tasking.TriggerExecutionPayload{})
		ctrl.Finish()
		if b {
			t.Fatalf("result should be false")
		}
	})
}
