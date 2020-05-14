package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	api "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

func TestDynamicBasedTrigger_IsApplicable(t *testing.T) {
	var (
		ctrl          *gomock.Controller
		mockSitesRepo *mocks.MockUserSitesPersistence
		trigger       DynamicBasedTrigger

		partnerID = "pid"
		siteID    = "10"
		ctx       = context.WithValue(
			context.WithValue(
				context.Background(), config.PartnerIDKeyCTX, partnerID),
			config.SiteIDKeyCTX, siteID)
	)

	before := func() {
		ctrl = gomock.NewController(t)
		logger.Load(config.Config.Log)
		mockSitesRepo = mocks.NewMockUserSitesPersistence(ctrl)
		trigger = DynamicBasedTrigger{
			log:       logger.Log,
			cache:     nil,
			sitesRepo: mockSitesRepo,
		}
	}

	dynamicGroupID, err := gocql.RandomUUID()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Run("positive", func(t *testing.T) {
		before()
		payload := api.TriggerExecutionPayload{
			DynamicGroupID: dynamicGroupID.String(),
		}

		task := models.Task{
			TargetsByType: models.TargetsByType{
				models.DynamicGroup: []string{dynamicGroupID.String()},
			},
			CreatedBy: partnerID,
		}

		usites := entities.UserSites{
			PartnerID: partnerID,
			UserID:    "uid",
			SiteIDs:   []int64{10},
		}

		mockSitesRepo.EXPECT().Sites(gomock.Any(), partnerID, partnerID).Return(usites, nil)

		b := trigger.IsApplicable(ctx, task, payload)
		ctrl.Finish()
		if !b {
			t.Fatalf("expected result to be true")
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		before()
		payload := api.TriggerExecutionPayload{
			DynamicGroupID: "invalid uuid",
		}

		task := models.Task{}

		b := trigger.IsApplicable(ctx, task, payload)
		ctrl.Finish()
		if b {
			t.Fatalf("expected result to be false")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		before()
		payload := api.TriggerExecutionPayload{
			DynamicGroupID: dynamicGroupID.String(),
		}

		task := models.Task{
			TargetsByType: models.TargetsByType{
				models.ManagedEndpoint: []string{dynamicGroupID.String()},
			},
			CreatedBy: partnerID,
		}

		b := trigger.IsApplicable(ctx, task, payload)
		ctrl.Finish()
		if b {
			t.Fatalf("expected result to be false")
		}
	})

	t.Run("negative_3", func(t *testing.T) {
		before()
		payload := api.TriggerExecutionPayload{
			DynamicGroupID: dynamicGroupID.String(),
		}

		task := models.Task{
			TargetsByType: models.TargetsByType{
				models.DynamicGroup: []string{dynamicGroupID.String()},
			},
			CreatedBy: partnerID,
		}

		usites := entities.UserSites{
			PartnerID: partnerID,
			UserID:    "uid",
			SiteIDs:   []int64{10},
		}

		mockSitesRepo.EXPECT().Sites(gomock.Any(), partnerID, partnerID).Return(usites, errors.New(""))

		b := trigger.IsApplicable(ctx, task, payload)
		ctrl.Finish()
		if b {
			t.Fatalf("expected result to be false")
		}
	})

	t.Run("negative_4", func(t *testing.T) {
		before()
		payload := api.TriggerExecutionPayload{
			DynamicGroupID: dynamicGroupID.String(),
		}

		task := models.Task{
			TargetsByType: models.TargetsByType{
				models.DynamicGroup: []string{dynamicGroupID.String()},
			},
			CreatedBy: partnerID,
		}

		usites := entities.UserSites{
			PartnerID: partnerID,
			UserID:    "uid",
			SiteIDs:   []int64{10},
		}

		ctx := context.WithValue(
			context.WithValue(
				context.Background(), config.PartnerIDKeyCTX, partnerID),
			config.SiteIDKeyCTX, "invalid integer")

		mockSitesRepo.EXPECT().Sites(gomock.Any(), partnerID, partnerID).Return(usites, nil)

		b := trigger.IsApplicable(ctx, task, payload)
		ctrl.Finish()
		if b {
			t.Fatalf("expected result to be false")
		}
	})

	t.Run("negative_5", func(t *testing.T) {
		before()
		payload := api.TriggerExecutionPayload{
			DynamicGroupID: dynamicGroupID.String(),
		}

		task := models.Task{
			TargetsByType: models.TargetsByType{
				models.DynamicGroup: []string{dynamicGroupID.String()},
			},
			CreatedBy: partnerID,
		}

		usites := entities.UserSites{
			PartnerID: partnerID,
			UserID:    "uid",
			SiteIDs:   []int64{20},
		}

		mockSitesRepo.EXPECT().Sites(gomock.Any(), partnerID, partnerID).Return(usites, nil)

		b := trigger.IsApplicable(ctx, task, payload)
		ctrl.Finish()
		if b {
			t.Fatalf("expected result to be false")
		}
	})
}

func TestDynamicBasedTrigger_GetTask(t *testing.T) {
	var (
		ctrl                *gomock.Controller
		mockDG              *mocks.MockDynamicGroups
		taskPersistanceMock *mocks.MockTaskPersistence
		trigger             DynamicBasedTrigger

		partnerID     = "pid"
		endpointID, _ = gocql.RandomUUID()
		taskID, _     = gocql.RandomUUID()
		ctx           = context.WithValue(
			context.WithValue(
				context.Background(), config.PartnerIDKeyCTX, partnerID),
			config.EndpointIDKeyCTX, endpointID)
	)

	before := func() {
		ctrl = gomock.NewController(t)
		mockDG = mocks.NewMockDynamicGroups(ctrl)
		logger.Load(config.Config.Log)
		taskPersistanceMock = mocks.NewMockTaskPersistence(ctrl)
		trigger = DynamicBasedTrigger{
			dgClient:  mockDG,
			tasksRepo: taskPersistanceMock,
			log:       logger.Log,
			cache:     nil,
			sitesRepo: nil,
		}
	}

	t.Run("positive", func(t *testing.T) {
		before()
		tasks := []models.Task{{
			TargetsByType: models.TargetsByType{
				models.DynamicGroup: nil,
			},
		}}

		taskPersistanceMock.EXPECT().GetByIDs(ctx, gomock.Any(), partnerID, true, taskID).Return(tasks, nil)
		_, err := trigger.GetTask(ctx, taskID)
		ctrl.Finish()
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		before()
		taskPersistanceMock.EXPECT().GetByIDs(ctx, gomock.Any(), partnerID, true, taskID).Return(nil, errors.New(""))
		_, err := trigger.GetTask(ctx, taskID)
		ctrl.Finish()
		if err == nil {
			t.Fatalf("error cannot be empty")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		before()
		tasks := []models.Task{}
		taskPersistanceMock.EXPECT().GetByIDs(ctx, gomock.Any(), partnerID, true, taskID).Return(tasks, nil)
		_, err := trigger.GetTask(ctx, taskID)
		ctrl.Finish()
		if err == nil {
			t.Fatalf("error cannot be empty")
		}
	})
}

func TestDynamicBasedTrigger_Deactivate(t *testing.T) {
	var (
		ctrl    *gomock.Controller
		mockDG  *mocks.MockDynamicGroups
		trigger DynamicBasedTrigger

		triggerType = "tt"
	)

	dynamicGroupID, err := gocql.RandomUUID()
	if err != nil {
		t.Fatalf(err.Error())
	}

	before := func() {
		ctrl = gomock.NewController(t)
		mockDG = mocks.NewMockDynamicGroups(ctrl)
		logger.Load(config.Config.Log)
		trigger = DynamicBasedTrigger{
			dgClient: mockDG,
			log:      logger.Log,
		}
	}

	t.Run("positive", func(t *testing.T) {
		before()
		tasks := []models.Task{{
			TargetsByType: models.TargetsByType{
				models.DynamicGroup: []string{dynamicGroupID.String()},
			},
		}}

		mockDG.EXPECT().StopGroupsMonitoring(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		err := trigger.Deactivate(context.TODO(), triggerType, tasks)
		ctrl.Finish()
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		before()
		tasks := []models.Task{{
			TargetsByType: models.TargetsByType{
				models.ManagedEndpoint: nil, // not dynamic group
			},
		}}

		err := trigger.Deactivate(context.TODO(), triggerType, tasks)
		ctrl.Finish()
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		before()
		tasks := []models.Task{{
			TargetsByType: models.TargetsByType{
				models.DynamicGroup: []string{dynamicGroupID.String()},
			},
		}}

		mockDG.EXPECT().StopGroupsMonitoring(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New(""))
		err := trigger.Deactivate(context.TODO(), triggerType, tasks)
		ctrl.Finish()
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})

}

func TestDynamicBasedTrigger_Activate(t *testing.T) {
	var (
		ctrl    *gomock.Controller
		mockDG  *mocks.MockDynamicGroups
		trigger DynamicBasedTrigger

		triggerType = "tt"
	)

	dynamicGroupID, err := gocql.RandomUUID()
	if err != nil {
		t.Fatalf(err.Error())
	}
	logger.Load(config.Config.Log)

	before := func() {
		ctrl = gomock.NewController(t)
		mockDG = mocks.NewMockDynamicGroups(ctrl)
		trigger = DynamicBasedTrigger{
			dgClient: mockDG,
			log:      logger.Log,
		}
	}

	t.Run("positive", func(t *testing.T) {
		before()
		tasks := []models.Task{{
			TargetsByType: models.TargetsByType{
				models.DynamicGroup: []string{dynamicGroupID.String()},
			},
		}}

		mockDG.EXPECT().StartMonitoringGroups(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		err := trigger.Activate(context.TODO(), triggerType, tasks)
		ctrl.Finish()
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		before()
		tasks := []models.Task{{
			TargetsByType: models.TargetsByType{
				models.ManagedEndpoint: nil, // not dynamic group
			},
		}}

		err := trigger.Activate(context.TODO(), triggerType, tasks)
		ctrl.Finish()
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		before()
		tasks := []models.Task{{
			TargetsByType: models.TargetsByType{
				models.DynamicGroup: []string{dynamicGroupID.String()},
			},
		}}

		mockDG.EXPECT().StartMonitoringGroups(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New(""))
		err := trigger.Activate(context.TODO(), triggerType, tasks)
		ctrl.Finish()
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})
}
