package trigger

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	mockUC "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mock-usecases"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/triggers"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/trigger/handlers"
)

func TestNew(t *testing.T) {
	RegisterTestingT(t)
	expected := &Service{}
	actual := New(nil, nil, models.DataBaseConnectors{}, repository.DatabaseRepositories{}, nil, integration.ExternalClients{})
	Ω(actual).To(Equal(expected), fmt.Sprintf("failed on unexpected value of result %v", expected))
}

func TestGetTriggerHandler(t *testing.T) {
	s := New(nil, nil, models.DataBaseConnectors{}, repository.DatabaseRepositories{}, nil, integration.ExternalClients{})
	handler := s.getTriggerHandler(triggers.DynamicGroupExitTrigger)
	switch handler.(type) {
	case *handlers.DynamicBasedTrigger:
		return
	default:
		t.Errorf("got wrong handler - %v, expected - DynamicBasedTrigger handler", reflect.TypeOf(handler))
	}

	handler = s.getTriggerHandler(triggers.StartupTrigger)
	switch handler.(type) {
	case *handlers.DefaultTriggerHandler:
		return
	default:
		t.Errorf("got wrong handler - %v, expected - DefaultTriggerHandler", reflect.TypeOf(handler))
	}
}

func TestService_Activate(t *testing.T) {
	logger.Load(config.Config.Log)
	id := gocql.TimeUUID()
	type test struct {
		name  string
		tasks []models.Task
		mock  func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine, persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo)
	}

	validTasks := []models.Task{
		{
			ID:        id,
			PartnerID: "1",
			Schedule: apiModels.Schedule{
				Regularity:    apiModels.Trigger,
				EndRunTime:    time.Now(),
				TriggerTypes:  []string{triggers.MockGeneric},
				TriggerFrames: []apiModels.TriggerFrame{{TriggerType: triggers.MockGeneric}},
			},
			ManagedEndpointID: id,
			RunTimeUTC:        time.Now(),
			Targets: models.Target{
				IDs:  []string{id.String()},
				Type: models.DynamicGroup,
			},
		},
	}

	tests := []test{
		{
			name: "testCase 1 - zero tasks",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo) {
				return nil, nil, nil, nil, nil, nil, nil
			},
			tasks: []models.Task{},
		},
		{
			name: "testCase 2 - handler activate returned error",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo) {
				return nil, nil, nil, nil, nil, nil, nil
			},
			tasks: []models.Task{{
				ID:        id,
				PartnerID: "1",
				Schedule: apiModels.Schedule{
					Regularity:    apiModels.Trigger,
					EndRunTime:    time.Now(),
					TriggerTypes:  []string{triggers.MockGeneric},
					TriggerFrames: []apiModels.TriggerFrame{{TriggerType: triggers.MockGeneric}},
				},
				ManagedEndpointID: id,
				Parameters:        mockUC.HandlerErrorActivate,
				RunTimeUTC:        time.Now(),
				Targets: models.Target{
					IDs:  []string{id.String()},
					Type: models.DynamicGroup,
				},
			},
			},
		},
		{
			name: "testCase 3 - trigger repo insert error",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo) {
				err := errors.New("err")

				tr := mocks.NewMockTriggersRepo(ctl)
				trigger := entities.ActiveTrigger{
					TaskID:    tasks[0].ID,
					Type:      tasks[0].Schedule.TriggerTypes[0],
					PartnerID: tasks[0].PartnerID,
				}
				tr.EXPECT().Insert(trigger).Return(err)
				return nil, nil, nil, nil, nil, nil, tr
			},
			tasks: validTasks,
		},
		{
			name: "testCase 4 - trigger counter returned error",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo) {
				err := errors.New("err")

				tr := mocks.NewMockTriggersRepo(ctl)
				trigger := entities.ActiveTrigger{
					TaskID:    tasks[0].ID,
					Type:      tasks[0].Schedule.TriggerTypes[0],
					PartnerID: tasks[0].PartnerID,
				}

				tr.EXPECT().Insert(trigger)
				tr.EXPECT().GetTriggerCounterByType(tasks[0].Schedule.TriggerTypes[0]).Return(entities.TriggerCounter{}, err)
				return nil, nil, nil, nil, nil, nil, tr
			},
			tasks: validTasks,
		},
		{
			name: "testCase 5 - trigger counter is zero, definition get err",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo) {
				err := errors.New("err")

				tr := mocks.NewMockTriggersRepo(ctl)
				trigger := entities.ActiveTrigger{
					TaskID:    tasks[0].ID,
					Type:      tasks[0].Schedule.TriggerTypes[0],
					PartnerID: tasks[0].PartnerID,
				}

				tr.EXPECT().Insert(trigger)
				tr.EXPECT().GetTriggerCounterByType(tasks[0].Schedule.TriggerTypes[0]).Return(entities.TriggerCounter{}, nil)
				tr.EXPECT().GetDefinition(tasks[0].Schedule.TriggerTypes[0]).Return(entities.TriggerDefinition{}, err)
				return nil, nil, nil, nil, nil, nil, tr
			},
			tasks: validTasks,
		},
		{
			name: "testCase 6 - trigger counter is zero, update AE policy err",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo) {
				err := errors.New("err")

				tr := mocks.NewMockTriggersRepo(ctl)
				trigger := entities.ActiveTrigger{
					TaskID:    tasks[0].ID,
					Type:      tasks[0].Schedule.TriggerTypes[0],
					PartnerID: tasks[0].PartnerID,
				}

				tr.EXPECT().Insert(trigger)
				tr.EXPECT().GetTriggerCounterByType(tasks[0].Schedule.TriggerTypes[0]).Return(entities.TriggerCounter{}, nil)
				tr.EXPECT().GetDefinition(tasks[0].Schedule.TriggerTypes[0])

				ae := mocks.NewMockAutomationEngine(ctl)
				ae.EXPECT().UpdateRemotePolicies(gomock.Any(), gomock.Any()).Return("", err)
				return nil, nil, ae, nil, nil, nil, tr
			},
			tasks: validTasks,
		},
		{
			name: "testCase 7 - trigger counter is zero, increase counter err but still proceed ",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo) {
				err := errors.New("err")

				tr := mocks.NewMockTriggersRepo(ctl)
				trigger := entities.ActiveTrigger{
					TaskID:    tasks[0].ID,
					Type:      tasks[0].Schedule.TriggerTypes[0],
					PartnerID: tasks[0].PartnerID,
				}

				tr.EXPECT().Insert(trigger)
				tr.EXPECT().GetTriggerCounterByType(tasks[0].Schedule.TriggerTypes[0]).Return(entities.TriggerCounter{}, nil)
				tr.EXPECT().GetDefinition(tasks[0].Schedule.TriggerTypes[0])

				ae := mocks.NewMockAutomationEngine(ctl)
				ae.EXPECT().UpdateRemotePolicies(gomock.Any(), gomock.Any()).Return("policyID", nil)

				tr.EXPECT().IncreaseCounter(entities.TriggerCounter{PolicyID: "policyID", TriggerID: triggers.MockGeneric}).Return(err)
				tr.EXPECT().GetAllByType(gomock.Any(), gomock.Any(), gomock.Any(), false).Return([]entities.ActiveTrigger{}, err)

				tasksRepo := mocks.NewMockTaskPersistence(ctl)

				var expectedTasks []models.Task
				for _, task := range tasks {
					if !task.Schedule.EndRunTime.IsZero() { // we don't need to reschedule 'never' trigger
						task.RunTimeUTC = task.Schedule.EndRunTime.UTC()
						expectedTasks = append(expectedTasks, task)
						continue
					}
					task.RunTimeUTC = task.RunTimeUTC.Add(hundredYears) // never task will never end
					expectedTasks = append(expectedTasks, task)
				}
				tasksRepo.EXPECT().UpdateSchedulerFields(gomock.Any(), expectedTasks).Return(err)

				return nil, tasksRepo, ae, nil, nil, nil, tr
			},
			tasks: validTasks,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			_, tasksRepo, ae, cache, dg, sites, tr := tc.mock(tc.tasks, mockCtrl)
			s := Service{
				taskRepo:         tasksRepo,
				log:              logger.Log,
				automationEngine: ae,
				cache:            cache,
				dynamicGroups:    dg,
				sitesRepo:        sites,
				triggerRepo:      tr,
			}
			s.Activate(context.TODO(), tc.tasks)
		})
	}
}

func TestService_Deactivate(t *testing.T) {
	id := gocql.TimeUUID()
	currentTime := time.Now()
	logger.Load(config.Config.Log)
	type test struct {
		name  string
		tasks []models.Task
		mock  func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine, persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo)
	}
	validTasks := []models.Task{
		{
			ID:        id,
			PartnerID: "1",
			Schedule: apiModels.Schedule{
				Regularity:   apiModels.Trigger,
				EndRunTime:   currentTime,
				TriggerTypes: []string{triggers.MockGeneric},
			},
			ManagedEndpointID: id,
			RunTimeUTC:        currentTime,
			Targets: models.Target{
				IDs:  []string{id.String()},
				Type: models.DynamicGroup,
			},
		},
	}

	tests := []test{
		{
			name: "testCase 1 - zero tasks",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo) {
				return nil, nil, nil, nil, nil, nil, nil
			},
			tasks: []models.Task{},
		},
		{
			name: "testCase 2 - handler activate returned error",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo) {
				return nil, nil, nil, nil, nil, nil, nil
			},
			tasks: []models.Task{{
				ID:                id,
				PartnerID:         "1",
				Schedule:          apiModels.Schedule{Regularity: apiModels.Trigger, EndRunTime: time.Now(), TriggerTypes: []string{triggers.MockGeneric}},
				ManagedEndpointID: id,
				Parameters:        mockUC.HandlerErrorDeactivate,
				RunTimeUTC:        time.Now(),
				Targets: models.Target{
					IDs:  []string{id.String()},
					Type: models.DynamicGroup,
				},
			},
			},
		},
		{
			name: "testCase 3 - trigger repo insert error",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo) {
				err := errors.New("err")

				tr := mocks.NewMockTriggersRepo(ctl)
				trigger := entities.ActiveTrigger{
					TaskID:    tasks[0].ID,
					Type:      tasks[0].Schedule.TriggerTypes[0],
					PartnerID: tasks[0].PartnerID,
				}
				tr.EXPECT().Delete(trigger).Return(err)
				return nil, nil, nil, nil, nil, nil, tr
			},
			tasks: validTasks,
		},
		{
			name: "testCase 4 - trigger counter returned error",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo) {
				err := errors.New("err")

				tr := mocks.NewMockTriggersRepo(ctl)
				trigger := entities.ActiveTrigger{
					TaskID:    tasks[0].ID,
					Type:      tasks[0].Schedule.TriggerTypes[0],
					PartnerID: tasks[0].PartnerID,
				}

				tr.EXPECT().Delete(trigger)
				tr.EXPECT().GetTriggerCounterByType(tasks[0].Schedule.TriggerTypes[0]).Return(entities.TriggerCounter{}, err)

				tr.EXPECT().GetAllByType(gomock.Any(), gomock.Any(), gomock.Any(), false).Return([]entities.ActiveTrigger{}, err)

				tasksRepo := mocks.NewMockTaskPersistence(ctl)
				tasksRepo.EXPECT().UpdateSchedulerFields(gomock.Any(), gomock.Any())
				return nil, tasksRepo, nil, nil, nil, nil, tr
			},
			tasks: validTasks,
		},
		{
			name: "testCase 5 - remove policy err",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo) {
				err := errors.New("err")

				tr := mocks.NewMockTriggersRepo(ctl)
				trigger := entities.ActiveTrigger{
					TaskID:    tasks[0].ID,
					Type:      tasks[0].Schedule.TriggerTypes[0],
					PartnerID: tasks[0].PartnerID,
				}

				counter := entities.TriggerCounter{Count: 1, PolicyID: "id"}
				tr.EXPECT().Delete(trigger)
				tr.EXPECT().GetTriggerCounterByType(tasks[0].Schedule.TriggerTypes[0]).Return(counter, nil)
				tr.EXPECT().GetAllByType(gomock.Any(), gomock.Any(), gomock.Any(), false).Return([]entities.ActiveTrigger{}, err)

				ae := mocks.NewMockAutomationEngine(ctl)
				tr.EXPECT().GetDefinition(gomock.Any())
				ae.EXPECT().RemovePolicy(gomock.Any(), gomock.Any()).Return(err)

				tasksRepo := mocks.NewMockTaskPersistence(ctl)
				tasksRepo.EXPECT().UpdateSchedulerFields(gomock.Any(), gomock.Any())
				return nil, tasksRepo, ae, nil, nil, nil, tr
			},
			tasks: validTasks,
		},
		{
			name: "testCase 6 - decrease counter err, update tasks err",
			mock: func(tasks []models.Task, ctl *gomock.Controller) (logger.Logger, models.TaskPersistence, integration.AutomationEngine,
				persistency.Cache, integration.DynamicGroups, models.UserSitesPersistence, repository.TriggersRepo) {
				err := errors.New("err")

				tr := mocks.NewMockTriggersRepo(ctl)
				trigger := entities.ActiveTrigger{
					TaskID:    tasks[0].ID,
					Type:      tasks[0].Schedule.TriggerTypes[0],
					PartnerID: tasks[0].PartnerID,
				}

				counter := entities.TriggerCounter{Count: 1, PolicyID: "id"}
				tr.EXPECT().Delete(trigger)
				tr.EXPECT().GetTriggerCounterByType(tasks[0].Schedule.TriggerTypes[0]).Return(counter, nil)
				tr.EXPECT().GetAllByType(gomock.Any(), gomock.Any(), gomock.Any(), false).Return([]entities.ActiveTrigger{}, err)

				ae := mocks.NewMockAutomationEngine(ctl)
				tr.EXPECT().GetDefinition(gomock.Any())
				ae.EXPECT().RemovePolicy(gomock.Any(), gomock.Any())

				tr.EXPECT().DecreaseCounter(counter).Return(err)

				tasksRepo := mocks.NewMockTaskPersistence(ctl)

				var expectedTasks []models.Task
				for _, task := range tasks {
					task.State = statuses.TaskStateInactive
					expectedTasks = append(expectedTasks, task)
				}
				tasksRepo.EXPECT().UpdateSchedulerFields(gomock.Any(), expectedTasks).Return(err)
				return nil, tasksRepo, ae, nil, nil, nil, tr
			},
			tasks: validTasks,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			_, tasksRepo, ae, cache, dg, sites, tr := tc.mock(tc.tasks, mockCtrl)
			s := Service{
				taskRepo:         tasksRepo,
				log:              logger.Log,
				automationEngine: ae,
				cache:            cache,
				dynamicGroups:    dg,
				sitesRepo:        sites,
				triggerRepo:      tr,
			}
			s.Deactivate(context.TODO(), tc.tasks)
		})
	}
}

func TestValidateCtx(t *testing.T) {
	s := New(nil, nil, models.DataBaseConnectors{}, repository.DatabaseRepositories{}, nil, integration.ExternalClients{})

	ctx := context.WithValue(context.TODO(), config.TriggerTypeIDKeyCTX, 1) // must be string)
	err := s.validateContext(ctx)
	if err == nil {
		t.Errorf("expected err, because triggerType must be string")
	}
	ctx = context.WithValue(ctx, config.TriggerTypeIDKeyCTX, "trigger") // must be string)

	ctx = context.WithValue(ctx, config.PartnerIDKeyCTX, 1)
	err = s.validateContext(ctx)
	if err == nil {
		t.Errorf("expected err, because PartnerID must be string")
	}
	ctx = context.WithValue(ctx, config.PartnerIDKeyCTX, "partnerID")

	ctx = context.WithValue(ctx, config.EndpointIDKeyCTX, 1)
	err = s.validateContext(ctx)
	if err == nil {
		t.Errorf("expected err, because EndpointID must be gocql.UUID")
	}
	ctx = context.WithValue(ctx, config.EndpointIDKeyCTX, gocql.TimeUUID())

	ctx = context.WithValue(ctx, config.SiteIDKeyCTX, 1)
	err = s.validateContext(ctx)
	if err == nil {
		t.Errorf("expected err, because siteID must be string")
	}
	ctx = context.WithValue(ctx, config.SiteIDKeyCTX, "siteID")

	ctx = context.WithValue(ctx, config.ClientIDKeyCTX, 1)
	err = s.validateContext(ctx)
	if err == nil {
		t.Errorf("expected err, because ClientID must be string")
	}

	ctx = context.WithValue(ctx, config.ClientIDKeyCTX, "ClientID")
	err = s.validateContext(ctx)
	if err != nil {
		t.Errorf("expected nil but gor %v", err)
	}
}

func TestService_GetTask(t *testing.T) {
	s := New(nil, nil, models.DataBaseConnectors{}, repository.DatabaseRepositories{}, nil, integration.ExternalClients{})
	if _, err := s.GetTask(context.TODO(), gocql.TimeUUID()); err == nil {
		t.Errorf("expected err becaues of empty context but got nil")
	}

	ctx := getTestContext()
	if _, err := s.GetTask(ctx, gocql.TimeUUID()); err != nil {
		t.Errorf("expected nil but got %v", err)
	}
}

func TestService_IsApplicable(t *testing.T) {
	s := New(nil, nil, models.DataBaseConnectors{}, repository.DatabaseRepositories{}, nil, integration.ExternalClients{})
	if applicable := s.IsApplicable(context.TODO(), models.Task{}, apiModels.TriggerExecutionPayload{}); applicable {
		t.Errorf("expected false but got true")
	}

	ctx := getTestContext()
	if applicable := s.IsApplicable(ctx, models.Task{}, apiModels.TriggerExecutionPayload{}); !applicable {
		t.Errorf("expected true but got false")
	}
}

func TestService_PostExecution(t *testing.T) {
	s := New(nil, nil, models.DataBaseConnectors{}, repository.DatabaseRepositories{}, nil, integration.ExternalClients{})
	if err := s.PostExecution(triggers.MockGeneric, models.Task{}); err != nil {
		t.Errorf("expected false but got true")
	}
}

func TestService_PreExecution(t *testing.T) {
	s := New(nil, nil, models.DataBaseConnectors{}, repository.DatabaseRepositories{}, nil, integration.ExternalClients{})
	if err := s.PreExecution(triggers.MockGeneric, models.Task{}); err != nil {
		t.Errorf("expected false but got true")
	}
}

func getTestContext() context.Context {
	ctx := context.TODO()
	ctx = context.WithValue(ctx, config.PartnerIDKeyCTX, "partnerID")
	ctx = context.WithValue(ctx, config.SiteIDKeyCTX, "siteID")
	ctx = context.WithValue(ctx, config.ClientIDKeyCTX, "clientID")
	ctx = context.WithValue(ctx, config.EndpointIDKeyCTX, gocql.TimeUUID())
	ctx = context.WithValue(ctx, config.TriggerTypeIDKeyCTX, triggers.MockGeneric)
	return ctx
}
