package trigger

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	managedEndpoints "gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/managed-endpoints"

	"gopkg.in/jarcoal/httpmock.v1"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/triggers"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository"
)

const (
	validUUIDstr  = "58a1af2f-6579-4aec-b45d-5dfde879ef01"
	validUUID2str = "58a1af2f-6579-4aec-b45d-5dfde879ef21"
)

func init() {
	config.Load()
	logger.Load(config.Config.Log)
}

var (
	statusesMap   = make(map[gocql.UUID]statuses.TaskInstanceStatus)
	validUUID, _  = gocql.ParseUUID(validUUIDstr)
	validUUID2, _ = gocql.ParseUUID(validUUID2str)
)

func TestService_groupTaskInstances(t *testing.T) {
	timeNow := time.Now().UTC()

	type test struct {
		name           string
		tasksInst      []models.TaskInstance
		tasks          []models.Task
		isNeedHTTPMock bool
		mock           func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine, persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo, models.TaskInstancePersistence)
	}

	ctx := context.Background()
	validTasks := []models.Task{
		{
			ID:        validUUID,
			PartnerID: "1",
			Schedule: apiModels.Schedule{
				Regularity:   apiModels.Trigger,
				EndRunTime:   timeNow,
				TriggerTypes: []string{triggers.MockGeneric},
			},
			State:              statuses.TaskStateActive,
			ManagedEndpointID:  validUUID2,
			RunTimeUTC:         timeNow,
			LastTaskInstanceID: validUUID,
			Targets: models.Target{
				IDs:  []string{validUUID2str},
				Type: models.Site,
			},
		},
	}

	updatedTasks := []models.Task{
		{
			ID:        validUUID,
			PartnerID: "1",
			Schedule: apiModels.Schedule{
				Regularity:   apiModels.Trigger,
				EndRunTime:   timeNow,
				TriggerTypes: []string{triggers.MockGeneric},
			},
			State:              statuses.TaskStateActive,
			ManagedEndpointID:  validUUID2,
			RunTimeUTC:         timeNow,
			LastTaskInstanceID: validUUID,
			Targets: models.Target{
				IDs:  []string{validUUID2str},
				Type: models.Site,
			},
		},
	}

	validTaskInstances := []models.TaskInstance{
		{
			PartnerID:     "1",
			ID:            str2uuid("01000000-0000-0000-0000-000000000001"),
			TaskID:        validUUID,
			OriginID:      gocql.UUID{},
			StartedAt:     time.Now(),
			LastRunTime:   time.Now(),
			Statuses:      statusesMap,
			OverallStatus: statuses.TaskInstanceScheduled,
			FailureCount:  0,
			SuccessCount:  0,
			TriggeredBy:   "",
		},
	}

	triggers := []entities.ActiveTrigger{
		{
			TaskID:    validUUID,
			Type:      validTasks[0].Schedule.TriggerTypes[0],
			PartnerID: validTasks[0].PartnerID,
		},
	}

	logger.Load(config.Config.Log)
	tests := []test{
		{
			name: "testCase 1 - zero tasks",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo, models.TaskInstancePersistence) {

				err := errors.New("ActiveTriggersReopening: can't get active triggers. err - error")

				tr := mocks.NewMockTriggersRepo(ctl)
				tr.EXPECT().GetAll().Return(triggers, err)

				return nil, nil, nil, nil, nil, nil, tr, nil
			},
			tasks: validTasks,
		},
		{
			name: "testCase 2 - can't getting tasks",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo, models.TaskInstancePersistence) {

				var taskIDS []gocql.UUID
				taskIDS = append(taskIDS, tasks[0].ID)
				err := errors.New("err")

				tr := mocks.NewMockTriggersRepo(ctl)
				tasksRepo := mocks.NewMockTaskPersistence(ctl)
				tr.EXPECT().GetAll().Return(triggers, nil)
				tasksRepo.EXPECT().GetByIDs(gomock.Any(), nil, "1", false, validUUID).Return(tasks, err)

				return nil, tasksRepo, nil, nil, nil, nil, tr, nil
			},
			tasks: validTasks,
		},
		{
			name:           "testCase 3 - can't save tasks in DB",
			isNeedHTTPMock: true,
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo, models.TaskInstancePersistence) {

				err := errors.New("err")
				var (
					taskIDS []gocql.UUID
				)
				taskIDS = append(taskIDS, tasks[0].ID)

				tr := mocks.NewMockTriggersRepo(ctl)
				tasksRepo := mocks.NewMockTaskPersistence(ctl)
				sites := mocks.NewMockUserSitesPersistence(ctl)
				ti := mocks.NewMockTaskInstancePersistence(ctl)

				tr.EXPECT().GetAll().Return(triggers, nil)
				tasksRepo.EXPECT().GetByIDs(gomock.Any(), nil, "1", false, validUUID).Return(tasks, nil)
				tasksRepo.EXPECT().InsertOrUpdate(gomock.Any(), updatedTasks).Return(err)
				return nil, tasksRepo, nil, nil, nil, sites, tr, ti
			},
			tasks: validTasks,
		},
		{
			name:           "testCase 4 - can't save tasks in DB",
			isNeedHTTPMock: true,
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo, models.TaskInstancePersistence) {

				err := errors.New("err")
				var (
					taskIDS             []gocql.UUID
					managedEndpoindsMap = make(map[gocql.UUID]struct{})
				)
				taskIDS = append(taskIDS, tasks[0].ID)
				managedEndpoindsMap[validUUID] = struct{}{}
				inst := []gocql.UUID{tasks[0].LastTaskInstanceID}

				tr := mocks.NewMockTriggersRepo(ctl)
				tasksRepo := mocks.NewMockTaskPersistence(ctl)
				ti := mocks.NewMockTaskInstancePersistence(ctl)

				tr.EXPECT().GetAll().Return(triggers, nil)
				tasksRepo.EXPECT().GetByIDs(gomock.Any(), nil, "1", false, validUUID).Return(tasks, nil)
				tasksRepo.EXPECT().InsertOrUpdate(gomock.Any(), updatedTasks).Return(nil)
				ti.EXPECT().GetByIDs(gomock.Any(), inst).Return(validTaskInstances, err)
				return nil, tasksRepo, nil, nil, nil, nil, tr, ti
			},
			tasks: validTasks,
		},
		{
			name:           "testCase 4 - valid",
			isNeedHTTPMock: true,
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo, models.TaskInstancePersistence) {
				var (
					taskIDS             []gocql.UUID
					managedEndpointsMap = make(map[gocql.UUID]struct{})
				)

				taskIDS = append(taskIDS, tasks[0].ID)
				managedEndpointsMap[validUUID] = struct{}{}
				inst := []gocql.UUID{tasks[0].LastTaskInstanceID}

				tr := mocks.NewMockTriggersRepo(ctl)
				tasksRepo := mocks.NewMockTaskPersistence(ctl)
				ti := mocks.NewMockTaskInstancePersistence(ctl)

				tr.EXPECT().GetAll().Return(triggers, nil)
				tasksRepo.EXPECT().GetByIDs(gomock.Any(), nil, "1", false, taskIDS).Return(tasks, nil)
				tasksRepo.EXPECT().InsertOrUpdate(gomock.Any(), updatedTasks).Return(nil)
				ti.EXPECT().GetByIDs(gomock.Any(), inst).Return(validTaskInstances, nil)
				ti.EXPECT().GetNearestInstanceAfter(validTaskInstances[0].TaskID, validTaskInstances[0].StartedAt).Return(validTaskInstances[0], gocql.ErrNotFound)
				ti.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
				return nil, tasksRepo, nil, nil, nil, nil, tr, ti

			},
			tasks: validTasks,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpClient := http.DefaultClient
			if tc.isNeedHTTPMock {
				mes := []managedEndpoints.ManagedEndpoint{{
					ID:         validUUID,
					SiteID:     "5009",
					EndpointID: validUUID2,
				}}
				mesPayload, err := json.Marshal(&mes)
				if err != nil {
					t.Fatalf("cant marshall json %v", err)
				}
				httpmock.RegisterResponder("GET", "http://127.0.0.1:8084/asset/v1/partner/1/sites/58a1af2f-6579-4aec-b45d-5dfde879ef21/summary", httpmock.NewBytesResponder(200, mesPayload))
				httpmock.RegisterResponder("GET", "http://127.0.0.1:8084/asset/v1/partner/1/sites/58a1af2f-6579-4aec-b45d-5dfde879ef01/summary", httpmock.NewBytesResponder(200, mesPayload))
			}
			statusesMap[validUUID] = statuses.TaskInstanceScheduled
			_, tasksRepo, ae, cache, dg, sites, tr, ti := tc.mock(tc.tasks, mockCtrl)
			s := Service{
				taskRepo:         tasksRepo,
				log:              logger.Log,
				automationEngine: ae,
				cache:            cache,
				dynamicGroups:    dg,
				sitesRepo:        sites,
				triggerRepo:      tr,
				httpClient:       httpClient,
				taskInstanceRepo: ti,
			}
			s.ActiveTriggersReopening(ctx)
		})
	}
}

func str2uuid(stringUUID string) gocql.UUID {
	uuid, err := gocql.ParseUUID(stringUUID)
	if err != nil {
 		return gocql.UUID{}
 	}
	return uuid
}
