package handlers

import (
	"fmt"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/asset"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/triggers"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository"
	"golang.org/x/net/context"
)

func TestLogoffTrigger_Activate(t *testing.T) {
	NewLogOnOffTrigger(nil, nil, nil, nil, nil, nil, nil)
	dg := NewDynamicBasedTrigger(nil, nil, nil, nil, nil)
	dg.PostExecution(models.Task{})
	dg.PreExecution(models.Task{})
	dg.Update(context.TODO(), "", []models.Task{})

	mockCtrl := gomock.NewController(t)
	id := gocql.TimeUUID()
	asset.Load()

	type args struct {
		tasks []models.Task
	}

	type test struct {
		name string
		args args
		mock func(tasks []models.Task) (integration.AgentConfig, models.TaskPersistence, repository.TriggersRepo, logger.Logger, models.UserSitesPersistence, repository.Profiles, integration.Asset)
		init func(test)
	}

	validTasks := []models.Task{
		{
			ID:                id,
			PartnerID:         "400500",
			Schedule:          apiModels.Schedule{Regularity: apiModels.Trigger, EndRunTime: time.Now(), TriggerTypes: []string{triggers.LoginTrigger}},
			ManagedEndpointID: id,
			RunTimeUTC:        time.Now(),
			Targets: models.Target{
				IDs:  []string{id.String()},
				Type: models.ManagedEndpoint,
			},
		},
	}

	logger.Load(config.Config.Log)
	tests := []test{

		{
			name: "testCase 1 - good",
			mock: func(tasks []models.Task) (integration.AgentConfig, models.TaskPersistence, repository.TriggersRepo, logger.Logger, models.UserSitesPersistence, repository.Profiles, integration.Asset) {
				agent := mock.NewMockAgentConfig(mockCtrl)
				as := mock.NewMockAsset(mockCtrl)
				profiles := mock.NewMockProfiles(mockCtrl)

				as.EXPECT().GetSiteIDByEndpointID(gomock.Any(), tasks[0].PartnerID, id).Return("", "", nil)
				agent.EXPECT().Activate(gomock.Any(), gomock.Any(), gomock.Any(), tasks[0].PartnerID).Return(id, nil)
				profiles.EXPECT().Insert(tasks[0].ID, id)
				return agent, nil, nil, nil, nil, profiles, as
			},
			args: args{
				tasks: validTasks,
			},
		},

		{
			name: "testCase 2 - can't insert to profile",
			mock: func(tasks []models.Task) (integration.AgentConfig, models.TaskPersistence, repository.TriggersRepo, logger.Logger, models.UserSitesPersistence, repository.Profiles, integration.Asset) {
				agent := mock.NewMockAgentConfig(mockCtrl)
				as := mock.NewMockAsset(mockCtrl)
				profiles := mock.NewMockProfiles(mockCtrl)
				err := fmt.Errorf("err")

				as.EXPECT().GetSiteIDByEndpointID(gomock.Any(), tasks[0].PartnerID, id).Return("", "", nil)
				agent.EXPECT().Activate(gomock.Any(), gomock.Any(), gomock.Any(), tasks[0].PartnerID).Return(id, nil)
				profiles.EXPECT().Insert(tasks[0].ID, id).Return(err)
				return agent, nil, nil, nil, nil, profiles, as
			},
			args: args{
				tasks: validTasks,
			},
		},

		{
			name: "testCase 3 - can't activate",
			mock: func(tasks []models.Task) (integration.AgentConfig, models.TaskPersistence, repository.TriggersRepo, logger.Logger, models.UserSitesPersistence, repository.Profiles, integration.Asset) {
				agent := mock.NewMockAgentConfig(mockCtrl)
				as := mock.NewMockAsset(mockCtrl)
				err := fmt.Errorf("err")

				as.EXPECT().GetSiteIDByEndpointID(gomock.Any(), tasks[0].PartnerID, id).Return("", "", nil)
				agent.EXPECT().Activate(gomock.Any(), gomock.Any(), gomock.Any(), tasks[0].PartnerID).Return(id, err)
				return agent, nil, nil, nil, nil, nil, as
			},
			args: args{
				tasks: validTasks,
			},
		},

		{
			name: "testCase 4 - can't get info from asset",
			mock: func(tasks []models.Task) (integration.AgentConfig, models.TaskPersistence, repository.TriggersRepo, logger.Logger, models.UserSitesPersistence, repository.Profiles, integration.Asset) {
				as := mock.NewMockAsset(mockCtrl)
				err := fmt.Errorf("err")
				agent := mock.NewMockAgentConfig(mockCtrl)

				as.EXPECT().GetSiteIDByEndpointID(gomock.Any(), tasks[0].PartnerID, id).Return("", "", err)
				agent.EXPECT().Activate(gomock.Any(), gomock.Any(), gomock.Any(), tasks[0].PartnerID).Return(id, err)
				return agent, nil, nil, nil, nil, nil, as
			},
			args: args{
				tasks: validTasks,
			},
		},
	}

	for _, tc := range tests {
		agent, task, repoTrig, _, userSites, profile, asset := tc.mock(tc.args.tasks)
		triggerLogOnOff := logOnOffTrigger{
			DefaultTriggerHandler: &DefaultTriggerHandler{
				taskRepo:     task,
				triggersRepo: repoTrig,
				log:          logger.Log,
			},
			agentConf:     agent,
			userSites:     userSites,
			profilesRepo:  profile,
			assetsService: asset,
		}
		triggerLogOnOff.Activate(context.TODO(), triggers.LoginTrigger, tc.args.tasks)
	}
}

func TestLogoffTrigger_Deactivate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	id := gocql.TimeUUID()
	asset.Load()
	logger.Load(config.Config.Log)
	type args struct {
		tasks []models.Task
	}

	type test struct {
		name string
		args args
		mock func(tasks []models.Task) (integration.AgentConfig, models.TaskPersistence, repository.TriggersRepo, logger.Logger, models.UserSitesPersistence, repository.Profiles, integration.Asset)
		init func(test)
	}

	validTasks := []models.Task{
		{
			ID:                id,
			PartnerID:         "400500",
			Schedule:          apiModels.Schedule{Regularity: apiModels.Trigger, EndRunTime: time.Now(), TriggerTypes: []string{triggers.LogoutTrigger}},
			ManagedEndpointID: id,
			RunTimeUTC:        time.Now(),
			Targets: models.Target{
				IDs:  []string{id.String()},
				Type: models.ManagedEndpoint,
			},
		},
	}

	tests := []test{

		{
			name: "testCase 1 - good",
			mock: func(tasks []models.Task) (integration.AgentConfig, models.TaskPersistence, repository.TriggersRepo, logger.Logger, models.UserSitesPersistence, repository.Profiles, integration.Asset) {
				agent := mock.NewMockAgentConfig(mockCtrl)
				profiles := mock.NewMockProfiles(mockCtrl)

				profiles.EXPECT().GetByTaskID(id).Return(id, nil)
				agent.EXPECT().Deactivate(gomock.Any(), id, tasks[0].PartnerID).Return(nil)
				profiles.EXPECT().Delete(tasks[0].ID)
				return agent, nil, nil, nil, nil, profiles, nil
			},
			args: args{
				tasks: validTasks,
			},
		},

		{
			name: "testCase 2 - can't delete from profiles",
			mock: func(tasks []models.Task) (integration.AgentConfig, models.TaskPersistence, repository.TriggersRepo, logger.Logger, models.UserSitesPersistence, repository.Profiles, integration.Asset) {
				agent := mock.NewMockAgentConfig(mockCtrl)
				profiles := mock.NewMockProfiles(mockCtrl)
				err := fmt.Errorf("err")

				profiles.EXPECT().GetByTaskID(id).Return(id, nil)
				agent.EXPECT().Deactivate(gomock.Any(), id, tasks[0].PartnerID).Return(nil)
				profiles.EXPECT().Delete(tasks[0].ID).Return(err)

				return agent, nil, nil, nil, nil, profiles, nil
			},
			args: args{
				tasks: validTasks,
			},
		},

		{
			name: "testCase 3 - can't deactivate",
			mock: func(tasks []models.Task) (integration.AgentConfig, models.TaskPersistence, repository.TriggersRepo, logger.Logger, models.UserSitesPersistence, repository.Profiles, integration.Asset) {
				agent := mock.NewMockAgentConfig(mockCtrl)
				profiles := mock.NewMockProfiles(mockCtrl)
				err := fmt.Errorf("err")

				profiles.EXPECT().GetByTaskID(id).Return(id, nil)
				agent.EXPECT().Deactivate(gomock.Any(), id, tasks[0].PartnerID).Return(err)
				return agent, nil, nil, nil, nil, profiles, nil
			},
			args: args{
				tasks: validTasks,
			},
		},
		{
			name: "testCase 4 - can't get info from profiles",
			mock: func(tasks []models.Task) (integration.AgentConfig, models.TaskPersistence, repository.TriggersRepo, logger.Logger, models.UserSitesPersistence, repository.Profiles, integration.Asset) {
				profiles := mock.NewMockProfiles(mockCtrl)
				err := fmt.Errorf("err")

				profiles.EXPECT().GetByTaskID(id).Return(id, err)
				return nil, nil, nil, nil, nil, profiles, nil
			},
			args: args{
				tasks: validTasks,
			},
		},
	}

	for _, tc := range tests {
		agent, task, repoTrig, _, userSites, profile, asset := tc.mock(tc.args.tasks)
		triggerLogOnOff := logOnOffTrigger{
			DefaultTriggerHandler: &DefaultTriggerHandler{
				taskRepo:     task,
				triggersRepo: repoTrig,
				log:          logger.Log,
			},
			agentConf:     agent,
			userSites:     userSites,
			profilesRepo:  profile,
			assetsService: asset,
		}
		triggerLogOnOff.Deactivate(context.TODO(), triggers.LogoutTrigger, tc.args.tasks)
	}
}
