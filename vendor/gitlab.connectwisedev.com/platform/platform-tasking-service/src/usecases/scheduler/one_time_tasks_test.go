package scheduler

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mockLoggerTasking"
	modelMocks "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/model-mocks"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
)

func init() {
	logger.Load(config.Config.Log)
}

func TestOneTimeTasksBuilder_Build(t *testing.T) {
	RegisterTestingT(t)
	expected := &OneTimeTasks{}
	builder := OneTimeTasksBuilder{}
	actual := builder.Build()
	Ω(actual).To(Equal(expected), fmt.Sprintf(defaultMsg, expected))
}

func TestOneTimeTasks_shouldBeRun(t *testing.T) {
	RegisterTestingT(t)

	type payload struct {
		time time.Time
		ti   models.TaskInstance
		task models.Task
	}

	tc := []struct {
		name    string
		payload func() payload
	}{
		{
			name: "test-case1: can't find task instance status",
			payload: func() (p payload) {
				return p
			},
		},
		{
			name: "test-case2: status determines that the task should not be running",
			payload: func() (p payload) {
				uuid := gocql.TimeUUID()
				p = payload{
					ti: models.TaskInstance{
						Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{
							uuid: 8,
						},
					},
					task: models.Task{
						ManagedEndpointID: uuid,
					},
				}
				return p
			},
		},
		{
			name: "test-case3: task should be run",
			payload: func() (p payload) {
				uuid := gocql.TimeUUID()
				p = payload{
					ti: models.TaskInstance{
						Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{
							uuid: 1,
						},
					},
					task: models.Task{
						ManagedEndpointID: uuid,
						Schedule: apiModels.Schedule{
							EndRunTime: time.Now().Add(1 * time.Minute),
						},
					},
				}
				return p
			},
		},
		{
			name: "test-case4: end run time after current time",
			payload: func() (p payload) {
				uuid := gocql.TimeUUID()
				p = payload{
					ti: models.TaskInstance{
						Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{
							uuid: 1,
						},
					},
					task: models.Task{
						ManagedEndpointID: uuid,
						Schedule: apiModels.Schedule{
							EndRunTime: time.Now().Add(-1 * time.Minute),
						},
					},
					time: time.Now(),
				}
				return p
			},
		},
		{
			name: "test-case5: end run time is zero",
			payload: func() (p payload) {
				uuid := gocql.TimeUUID()
				p = payload{
					ti: models.TaskInstance{
						Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{
							uuid: 1,
						},
					},
					task: models.Task{
						ManagedEndpointID: uuid,
					},
				}
				return p
			},
		},
	}

	for _, test := range tc {
		payload := test.payload()
		oneTimeTasks := OneTimeTasks{}

		oneTimeTasks.shouldBeRun(payload.time, payload.ti, payload.task)
	}
}

func TestOneTimeTasks_sendExecutionResultsErr(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)
	tasks := []models.Task{{}, {}}
	currentTime := time.Now()

	// Mocks
	executionResultMock := NewMockExecutionResultRepo(mockCtrl)

	executionResult := tasking.ExecutionResultKafkaMessage{
		Message: tasking.ScriptPluginReturnMessage{
			ExecutionID:  tasks[0].LastTaskInstanceID.String(),
			TimestampUTC: currentTime,
			Status:       statuses.TaskInstanceFailedText,
			Stderr:       "failed by time out",
		},
		BrokerEnvelope: agent.BrokerEnvelope{
			EndpointID: tasks[0].ManagedEndpointID.String(),
			PartnerID:  tasks[0].PartnerID,
		},
	}

	// Expect execution
	err := errors.New("some error")
	executionResultMock.EXPECT().Publish(executionResult).Return(err).Times(1)
	executionResultMock.EXPECT().Publish(executionResult).Return(err).Times(1)

	oneTimeTasks := OneTimeTasks{
		executionResultRepo: executionResultMock,
		log:                 logger.Log,
	}
	oneTimeTasks.sendExecutionResultsErr(context.TODO(), currentTime, tasks...)
}

func TestOneTimeTasks_saveTimeExpiration(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()
	nowUUID := gocql.TimeUUID()
	type fields struct {
		executionExpirationRepo ExecutionExpirationRepo
		cacheRepo               CacheRepo
		log                     logger.Logger
	}
	type args struct {
		task models.Task
		ti   models.TaskInstance
	}

	type test struct {
		name   string
		fields fields
		args   args
		init   func(*test)
	}
	tests := []test{
		{
			name: "test-case2: can't insert to execution repo",
			args: args{
				task: models.Task{
					Parameters: "{{",
					OriginID:   nowUUID,
					Type:       models.TaskTypeScript,
				},
				ti: models.TaskInstance{
					Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{
						nowUUID: statuses.TaskInstanceRunning,
					},
				},
			},
			init: func(test *test) {
				data := models.TemplateDetails{
					ExpectedExecutionTimeSec: 5,
				}
				err := errors.New("some error")
				//initialize interfaces
				cache := NewMockCacheRepo(mockCtrl)
				executionExpirationRepo := NewMockExecutionExpirationRepo(mockCtrl)

				cache.EXPECT().GetByOriginID(ctx, "", nowUUID, true).Return(data, nil).Times(1)
				executionExpirationRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(err).Times(1)
				//set fields
				test.fields.cacheRepo = cache
				test.fields.log = logger.Log
				test.fields.executionExpirationRepo = executionExpirationRepo
			},
		},
	}
	for _, tt := range tests {
		tt.init(&tt)
		t.Run(tt.name, func(t *testing.T) {
			o := &OneTimeTasks{
				executionExpirationRepo: tt.fields.executionExpirationRepo,
				cacheRepo:               tt.fields.cacheRepo,
				log:                     logger.Log,
			}
			o.saveTimeExpiration(context.TODO(), tt.args.ti, 5)
		})
	}
}

func TestOneTimeTasks_deactivateTasks(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()
	currentTime := time.Now()
	me := gocql.TimeUUID()
	tasks := []models.Task{{ManagedEndpointID: me}}

	executionResult := tasking.ExecutionResultKafkaMessage{
		Message: tasking.ScriptPluginReturnMessage{
			ExecutionID:  tasks[0].LastTaskInstanceID.String(),
			TimestampUTC: currentTime,
			Status:       statuses.TaskInstanceFailedText,
			Stderr:       "failed by time out",
		},
		BrokerEnvelope: agent.BrokerEnvelope{
			EndpointID: tasks[0].ManagedEndpointID.String(),
			PartnerID:  tasks[0].PartnerID,
		},
	}

	//Execution result repo
	erRepo := NewMockExecutionResultRepo(mockCtrl)
	taskRepo := NewMockTaskRepo(mockCtrl)
	err := errors.New("some error")
	taskRepo.EXPECT().UpdateSchedulerFields(ctx, gomock.Any()).Return(err).Times(len(tasks))
	erRepo.EXPECT().Publish(executionResult).Times(1)

	oneTimeTasks := OneTimeTasks{
		taskRepo:            taskRepo,
		log:                 logger.Log,
		executionResultRepo: erRepo,
	}

	oneTimeTasks.deactivateTasks(tasks...)
}

func TestOneTimeTasks_recalculateRunTime(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()
	tasks := []models.Task{{}, {}}
	conf := config.Configuration{RecalculateTime: 2}

	// Mocks
	taskMock := NewMockTaskRepo(mockCtrl)

	// Add time to recalculate
	addTime := time.Duration(conf.RecalculateTime) * time.Second
	tasks[0].RunTimeUTC.Add(addTime)
	tasks[1].RunTimeUTC.Add(addTime)

	// Expect execution
	err := errors.New("some error")
	taskMock.EXPECT().UpdateSchedulerFields(ctx, tasks[0]).Return(err).Times(1)
	taskMock.EXPECT().UpdateSchedulerFields(ctx, tasks[1]).Return(err).Times(1)

	oneTimeTasks := OneTimeTasks{
		taskRepo: taskMock,
		log:      logger.Log,
	}
	oneTimeTasks.recalculateRunTime(tasks...)
}

func TestOneTimeTasks_updateTaskInstance(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)
	timeUUID := gocql.TimeUUID()
	timeUUID2 := gocql.TimeUUID()

	type payload struct {
		time             time.Time
		ttl              int
		taskInstanceRepo TaskInstanceRepo
		log              logger.Logger
	}

	tc := []struct {
		name    string
		payload func() payload
	}{
		{
			name: "test-case1: can't update taskInstance",
			payload: func() payload {
				taskInstanceMock := NewMockTaskInstanceRepo(mockCtrl)
				logMock := mockLoggerTasking.NewMockLogger(mockCtrl)
				err := errors.New("some error")
				t := time.Now()
				ttl := 0
				taskInstance := models.TaskInstance{
					LastRunTime: t,
					Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{
						timeUUID:  statuses.TaskInstancePending,
						timeUUID2: statuses.TaskInstanceStopped,
					},
				}
				taskInstanceMock.EXPECT().Insert(taskInstance, ttl).Return(err)
				p := payload{
					time:             t,
					ttl:              ttl,
					taskInstanceRepo: taskInstanceMock,
					log:              logMock,
				}
				return p
			},
		},
		{
			name: "test-case2: successful update taskInstance",
			payload: func() payload {
				taskInstanceMock := NewMockTaskInstanceRepo(mockCtrl)
				t := time.Now()
				ttl := 0
				taskInstance := models.TaskInstance{
					LastRunTime: t,
					Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{
						timeUUID:  statuses.TaskInstancePending,
						timeUUID2: statuses.TaskInstanceStopped,
					},
				}
				taskInstanceMock.EXPECT().Insert(taskInstance, ttl).Return(nil)
				p := payload{
					time:             t,
					ttl:              ttl,
					taskInstanceRepo: taskInstanceMock,
				}
				return p
			},
		},
	}

	for _, test := range tc {
		payload := test.payload()
		oneTimeTasks := OneTimeTasks{
			taskInstanceRepo: payload.taskInstanceRepo,
			log:              logger.Log,
		}
		ti := models.TaskInstance{
			Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{
				timeUUID:  statuses.TaskInstanceScheduled,
				timeUUID2: statuses.TaskInstanceStopped,
			},
		}
		oneTimeTasks.updateTaskInstance(context.TODO(), payload.time, ti)
	}
}

func TestOneTimeTasks_executeTasks(t *testing.T) {
	ctx := context.Background()
	ti := models.TaskInstance{}
	currentTime := time.Now()

	type payload struct {
		taskExecutionRepo       TaskExecutionRepo
		executionResultRepo     ExecutionResultRepo
		executionExpirationRepo ExecutionExpirationRepo
		taskRepo                TaskRepo
		log                     logger.Logger
		tasks                   []models.Task
		cache                   CacheRepo
	}

	tc := []struct {
		name    string
		payload func(mockCtrl *gomock.Controller) payload
	}{
		{
			name: "test-case4: successful result",
			payload: func(mockCtrl *gomock.Controller) payload {
				tasks := []models.Task{{RunTimeUTC: currentTime}}
				eps := []apiModels.ManagedEndpoint{
					{
						ID:          tasks[0].ManagedEndpointID.String(),
						NextRunTime: tasks[0].RunTimeUTC.UTC(),
					},
				}
				pattern := "%s/partners/%s/task-execution-results/task-instances/%s"
				webHookURL := fmt.Sprintf(pattern, "", tasks[0].PartnerID, ti.ID)
				executionPayload := apiModels.ExecutionPayload{
					ExecutionID:              ti.ID.String(),
					OriginID:                 ti.OriginID.String(),
					ManagedEndpoints:         eps,
					Parameters:               tasks[0].Parameters,
					TaskID:                   tasks[0].ID,
					WebhookURL:               webHookURL,
					Credentials:              tasks[0].Credentials,
					ExpectedExecutionTimeSec: 300,
				}

				err := errors.New("some error")

				//Mock repos
				logMock := mockLoggerTasking.NewMockLogger(mockCtrl)
				taskExecutionRepo := NewMockTaskExecutionRepo(mockCtrl)
				cache := NewMockCacheRepo(mockCtrl)
				executionExpirationRepo := NewMockExecutionExpirationRepo(mockCtrl)

				// Expects
				taskExecutionRepo.EXPECT().ExecuteTasks(gomock.Any(), executionPayload, tasks[0].PartnerID, tasks[0].Type)
				executionExpirationRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(err)
				cache.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any()).Return(300).AnyTimes()

				return payload{
					taskExecutionRepo:       taskExecutionRepo,
					executionExpirationRepo: executionExpirationRepo,
					log:                     logMock,
					tasks:                   tasks,
					cache:                   cache,
				}
			},
		},
		{
			name: "test-case5: can't execute task",
			payload: func(mockCtrl *gomock.Controller) payload {
				tasks := []models.Task{{RunTimeUTC: currentTime}}
				eps := []apiModels.ManagedEndpoint{
					{
						ID:          tasks[0].ManagedEndpointID.String(),
						NextRunTime: tasks[0].RunTimeUTC.UTC(),
					},
				}
				pattern := "%s/partners/%s/task-execution-results/task-instances/%s"
				webHookURL := fmt.Sprintf(pattern, "", tasks[0].PartnerID, ti.ID)
				executionPayload := apiModels.ExecutionPayload{
					ExecutionID:              ti.ID.String(),
					OriginID:                 ti.OriginID.String(),
					ManagedEndpoints:         eps,
					Parameters:               tasks[0].Parameters, //Common parameters for all tasks inside task instance
					TaskID:                   tasks[0].ID,
					WebhookURL:               webHookURL,
					Credentials:              tasks[0].Credentials,
					ExpectedExecutionTimeSec: 300,
				}

				err := errors.New("some error")

				//Mock repos
				taskExecutionRepo := NewMockTaskExecutionRepo(mockCtrl)
				taskRepo := NewMockTaskRepo(mockCtrl)
				cache := NewMockCacheRepo(mockCtrl)
				executionExpirationRepo := NewMockExecutionExpirationRepo(mockCtrl)

				// Expects
				taskExecutionRepo.EXPECT().ExecuteTasks(gomock.Any(), executionPayload, tasks[0].PartnerID, tasks[0].Type).Return(err)
				taskRepo.EXPECT().UpdateSchedulerFields(ctx, tasks[0]).Return(nil)
				cache.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any()).Return(300).AnyTimes()
				executionExpirationRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
				return payload{
					taskRepo:                taskRepo,
					taskExecutionRepo:       taskExecutionRepo,
					executionExpirationRepo: executionExpirationRepo,
					cache:                   cache,
					log:                     logger.Log,
					tasks:                   tasks,
				}
			},
		},
		{
			name: "test-case5.1: can't execute task - out of schedule",
			payload: func(mockCtrl *gomock.Controller) payload {
				tasks := []models.Task{{RunTimeUTC: currentTime.Add(time.Minute * -60)}}

				//Mock repos
				taskExecutionRepo := NewMockTaskExecutionRepo(mockCtrl)
				taskRepo := NewMockTaskRepo(mockCtrl)
				cache := NewMockCacheRepo(mockCtrl)

				cache.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any()).Return(300).AnyTimes()
				return payload{
					taskRepo:          taskRepo,
					taskExecutionRepo: taskExecutionRepo,
					cache:             cache,
					log:               logger.Log,
					tasks:             tasks,
				}
			},
		},
		{
			name: "test-case6: can't execute task, end run time is equal to run time ",
			payload: func(mockCtrl *gomock.Controller) payload {
				tasks := []models.Task{{
					RunTimeUTC: currentTime,
					Schedule: apiModels.Schedule{
						EndRunTime: currentTime,
					},
				}}
				eps := []apiModels.ManagedEndpoint{
					{
						ID:          tasks[0].ManagedEndpointID.String(),
						NextRunTime: tasks[0].RunTimeUTC.UTC(),
					},
				}
				pattern := "%s/partners/%s/task-execution-results/task-instances/%s"
				webHookURL := fmt.Sprintf(pattern, "", tasks[0].PartnerID, ti.ID)
				executionPayload := apiModels.ExecutionPayload{
					ExecutionID:              ti.ID.String(),
					OriginID:                 ti.OriginID.String(),
					ManagedEndpoints:         eps,
					Parameters:               tasks[0].Parameters,
					TaskID:                   tasks[0].ID,
					WebhookURL:               webHookURL,
					Credentials:              tasks[0].Credentials,
					ExpectedExecutionTimeSec: 300,
				}
				err := errors.New("some error")

				//Mock repos
				taskExecutionRepo := NewMockTaskExecutionRepo(mockCtrl)
				executionResultRepo := NewMockExecutionResultRepo(mockCtrl)
				cache := NewMockCacheRepo(mockCtrl)

				// Expects
				taskExecutionRepo.EXPECT().ExecuteTasks(gomock.Any(), executionPayload, tasks[0].PartnerID, tasks[0].Type).Return(err)
				executionResultRepo.EXPECT().Publish(gomock.Any()).Return(nil)
				cache.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any()).Return(300).AnyTimes()
				return payload{
					taskExecutionRepo:   taskExecutionRepo,
					executionResultRepo: executionResultRepo,
					log:                 logger.Log,
					tasks:               tasks,
					cache:               cache,
				}
			},
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			payload := test.payload(mockCtrl)
			oneTimeTasks := OneTimeTasks{
				taskRepo:                payload.taskRepo,
				taskExecutionRepo:       payload.taskExecutionRepo,
				executionResultRepo:     payload.executionResultRepo,
				executionExpirationRepo: payload.executionExpirationRepo,
				log:                     logger.Log,
				cacheRepo:               payload.cache,
			}
			oneTimeTasks.executeTasks(context.TODO(), currentTime, ti, payload.tasks...)
		})
	}
}

func TestOneTimeTasks_Process(t *testing.T) {
	currentTime := time.Now().UTC().Truncate(time.Minute)
	ctx := context.Background()

	tasks := []models.Task{{}}

	tc := []struct {
		name     string
		payload  func(mockCtrl *gomock.Controller) (TaskInstanceRepo, TaskRepo, logger.Logger)
		expected string
	}{
		{
			name: "test-case1: error with getting taskInstance",
			payload: func(mockCtrl *gomock.Controller) (TaskInstanceRepo, TaskRepo, logger.Logger) {
				taskInstanceRepo := NewMockTaskInstanceRepo(mockCtrl)
				logMock := mockLoggerTasking.NewMockLogger(mockCtrl)
				taskRepo := NewMockTaskRepo(mockCtrl)
				err := errors.New("some error")
				taskInstanceRepo.EXPECT().GetInstance(gomock.Any()).Return(models.TaskInstance{}, err)
				taskRepo.EXPECT().UpdateSchedulerFields(ctx, gomock.Any()).Return(nil).Times(len(tasks))
				return taskInstanceRepo, taskRepo, logMock
			},
			expected: "some error",
		},
		{
			name: "test-case1: successful execution",
			payload: func(mockCtrl *gomock.Controller) (TaskInstanceRepo, TaskRepo, logger.Logger) {
				taskInstanceRepo := NewMockTaskInstanceRepo(mockCtrl)
				taskRepo := NewMockTaskRepo(mockCtrl)
				logMock := mockLoggerTasking.NewMockLogger(mockCtrl)
				ti := models.TaskInstance{
					Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{
						tasks[0].ManagedEndpointID: 0,
					},
				}
				taskInstanceRepo.EXPECT().GetInstance(gomock.Any()).Return(ti, nil)
				taskInstanceRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				taskRepo.EXPECT().UpdateSchedulerFields(ctx, gomock.Any()).Return(nil).AnyTimes()
				return taskInstanceRepo, taskRepo, logMock
			},
			expected: "some error",
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			m := gomock.NewController(t)
			defer m.Finish()

			taskInstanceRepo, taskRepo, _ := test.payload(m)
			oneTimeTasks := OneTimeTasks{
				taskInstanceRepo: taskInstanceRepo,
				taskRepo:         taskRepo,
				log:              logger.Log,
				cacheRepo:        modelMocks.NewTemplateCacheMock(false),
			}

			oneTimeTasks.Process(context.Background(), currentTime, tasks)
		})
	}
	time.Sleep(time.Second)
}
