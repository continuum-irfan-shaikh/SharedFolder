package models_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
)

const defaultMsg = `failed on unexpected value of result "%v"`

var endpointID, _ = gocql.ParseUUID("b49082f1-6707-4732-ad8b-50d62b15c863")
var taskId, _ = gocql.ParseUUID("55a98dda-f855-4197-a4b3-5bc72c70fd65")
var taskInstID1, _ = gocql.ParseUUID("554258b7-1aba-423b-b0d6-e450f0db4780")
var lastTaskInstID1, _ = gocql.ParseUUID("fade2a3b-947e-4c2b-93f4-0330b5ad274a")
var lastTaskInstID2, _ = gocql.ParseUUID("a938d7b9-0252-4bf2-9af7-d093046e3ea9")
var originID, _ = gocql.ParseUUID("35ba90c3-be11-448c-b5c5-814e9fc337a6")

func init()  {
	logger.Load(config.Config.Log)
}

func TestExecutionResultViewRepoCassandra_Get(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrlr *gomock.Controller

	type payload struct {
		isNoc                   bool
		taskPersistence         func()
		taskInstancePersistence func()
		execResultPersistence   func()
		templateCache           func()
	}

	type expected struct {
		results []*models.ExecutionResultView
		err     error
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name: "Success_#1",
			payload: payload{
				isNoc: false,
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByPartnerAndManagedEndpointID(context.Background(), "partnerID", endpointID, 2).
						Return([]models.Task{
							{
								LastTaskInstanceID: lastTaskInstID1,
							},
							{
								LastTaskInstanceID: lastTaskInstID2,
							},
						}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				execResultPersistence: func() {
					execResultsPersistenceMock := mocks.NewMockExecutionResultPersistence(mockCtrlr)
					execResultsPersistenceMock.EXPECT().GetByTargetAndTaskInstanceIDs(endpointID, gomock.Any()).
						Return([]models.ExecutionResult{
							{},
						}, nil).Times(1)
					models.ExecutionResultPersistenceInstance = execResultsPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					taskInstancePersistenceMock.EXPECT().GetByIDs(context.Background(), gomock.Any()).
						Return([]models.TaskInstance{
							{},
						}, nil).Times(1)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				templateCache: func() {
					models.TemplateCacheInstance = mocks.NewMockTemplateCache(mockCtrlr)
				},
			},
			expected: expected{
				results: []*models.ExecutionResultView{},
			},
		},
		{
			name: "Success_#2",
			payload: payload{
				isNoc: false,
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByPartnerAndManagedEndpointID(context.Background(), "partnerID", endpointID, 2).
						Return([]models.Task{
							{
								LastTaskInstanceID: lastTaskInstID1,
								IsRequireNOCAccess: true,
							},
							{
								LastTaskInstanceID: lastTaskInstID2,
								IsRequireNOCAccess: true,
							},
						}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				execResultPersistence: func() {
					execResultsPersistenceMock := mocks.NewMockExecutionResultPersistence(mockCtrlr)
					execResultsPersistenceMock.EXPECT().GetByTargetAndTaskInstanceIDs(endpointID, gomock.Any()).
						Return([]models.ExecutionResult{
							{},
						}, nil).Times(1)
					models.ExecutionResultPersistenceInstance = execResultsPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					taskInstancePersistenceMock.EXPECT().GetByIDs(context.Background(), gomock.Any()).
						Return([]models.TaskInstance{
							{},
						}, nil).Times(1)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				templateCache: func() {
					models.TemplateCacheInstance = mocks.NewMockTemplateCache(mockCtrlr)
				},
			},
			expected: expected{
				results: []*models.ExecutionResultView{},
			},
		},
		{
			name: "Success_#3",
			payload: payload{
				isNoc: false,
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByPartnerAndManagedEndpointID(context.Background(), "partnerID", endpointID, 2).
						Return([]models.Task{
							{
								LastTaskInstanceID: lastTaskInstID1,
								State:              statuses.TaskStateActive,
								OriginID:           originID,
								Schedule: apiModels.Schedule{
									Regularity: apiModels.RunNow,
								},
							},
						}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				execResultPersistence: func() {
					execResultsPersistenceMock := mocks.NewMockExecutionResultPersistence(mockCtrlr)
					execResultsPersistenceMock.EXPECT().GetByTargetAndTaskInstanceIDs(endpointID, gomock.Any()).
						Return([]models.ExecutionResult{
							{
								TaskInstanceID:  lastTaskInstID1,
								ExecutionStatus: statuses.TaskInstanceScheduled,
							},
						}, nil).Times(1)
					models.ExecutionResultPersistenceInstance = execResultsPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					taskInstancePersistenceMock.EXPECT().GetByIDs(context.Background(), gomock.Any()).
						Return([]models.TaskInstance{
							{
								ID:       lastTaskInstID1,
								Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{endpointID: statuses.TaskInstanceRunning},
							},
						}, nil).Times(1)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				templateCache: func() {
					templateCacheMock := mocks.NewMockTemplateCache(mockCtrlr)
					templateCacheMock.EXPECT().GetByOriginID(context.Background(), "partnerID", originID, true).
						Return(models.TemplateDetails{}, nil).Times(1)
					models.TemplateCacheInstance = templateCacheMock
				},
			},
			expected: expected{
				results: []*models.ExecutionResultView{
					{
						ManagedEndpointID: endpointID,
						ExecutionID:       lastTaskInstID1,
						OriginID:          originID,
						CanBePostponed:    true,
						CanBeCanceled:     true,
						DeviceCount:       1,
						Regularity:        apiModels.RunNow,
						LastRunStatus:     statuses.TaskInstanceScheduled,
						Status:            statuses.TaskStateActive,
					},
				},
			},
		},
		{
			name: "Success_#4",
			payload: payload{
				isNoc: false,
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByPartnerAndManagedEndpointID(context.Background(), "partnerID", endpointID, 2).
						Return([]models.Task{
							{
								LastTaskInstanceID: lastTaskInstID1,
								State:              statuses.TaskStateActive,
								OriginID:           originID,
								Schedule: apiModels.Schedule{
									Regularity: apiModels.RunNow,
								},
							},
						}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				execResultPersistence: func() {
					execResultsPersistenceMock := mocks.NewMockExecutionResultPersistence(mockCtrlr)
					execResultsPersistenceMock.EXPECT().GetByTargetAndTaskInstanceIDs(endpointID, gomock.Any()).
						Return([]models.ExecutionResult{
							{
								TaskInstanceID:  lastTaskInstID1,
								ExecutionStatus: statuses.TaskInstanceFailed,
							},
						}, nil).Times(1)
					models.ExecutionResultPersistenceInstance = execResultsPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					taskInstancePersistenceMock.EXPECT().GetByIDs(context.Background(), gomock.Any()).
						Return([]models.TaskInstance{
							{
								ID:       lastTaskInstID1,
								Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{endpointID: statuses.TaskInstanceRunning},
							},
						}, nil).Times(1)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				templateCache: func() {
					templateCacheMock := mocks.NewMockTemplateCache(mockCtrlr)
					templateCacheMock.EXPECT().GetByOriginID(context.Background(), "partnerID", originID, true).
						Return(models.TemplateDetails{
							FailureMessage: "Fail",
						}, nil).Times(1)
					models.TemplateCacheInstance = templateCacheMock
				},
			},
			expected: expected{
				results: []*models.ExecutionResultView{
					{
						ManagedEndpointID: endpointID,
						ExecutionID:       lastTaskInstID1,
						OriginID:          originID,
						DeviceCount:       1,
						Regularity:        apiModels.RunNow,
						ResultMessage:     "Fail",
						LastRunStatus:     statuses.TaskInstanceFailed,
						Status:            statuses.TaskStateActive,
					},
				},
			},
		},
		{
			name: "Success_#5",
			payload: payload{
				isNoc: false,
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByPartnerAndManagedEndpointID(context.Background(), "partnerID", endpointID, 2).
						Return([]models.Task{
							{
								LastTaskInstanceID: lastTaskInstID1,
								State:              statuses.TaskStateActive,
								OriginID:           originID,
								Schedule: apiModels.Schedule{
									Regularity: apiModels.RunNow,
								},
							},
						}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				execResultPersistence: func() {
					execResultsPersistenceMock := mocks.NewMockExecutionResultPersistence(mockCtrlr)
					execResultsPersistenceMock.EXPECT().GetByTargetAndTaskInstanceIDs(endpointID, gomock.Any()).
						Return([]models.ExecutionResult{
							{
								TaskInstanceID:  lastTaskInstID1,
								ExecutionStatus: statuses.TaskInstanceSuccess,
							},
						}, nil).Times(1)
					models.ExecutionResultPersistenceInstance = execResultsPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					taskInstancePersistenceMock.EXPECT().GetByIDs(context.Background(), gomock.Any()).
						Return([]models.TaskInstance{
							{
								ID:       lastTaskInstID1,
								Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{endpointID: statuses.TaskInstanceRunning},
							},
						}, nil).Times(1)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				templateCache: func() {
					templateCacheMock := mocks.NewMockTemplateCache(mockCtrlr)
					templateCacheMock.EXPECT().GetByOriginID(context.Background(), "partnerID", originID, true).
						Return(models.TemplateDetails{
							SuccessMessage: "Success",
						}, nil).Times(1)
					models.TemplateCacheInstance = templateCacheMock
				},
			},
			expected: expected{
				results: []*models.ExecutionResultView{
					{
						ManagedEndpointID: endpointID,
						ExecutionID:       lastTaskInstID1,
						OriginID:          originID,
						DeviceCount:       1,
						Regularity:        apiModels.RunNow,
						ResultMessage:     "Success",
						LastRunStatus:     statuses.TaskInstanceSuccess,
						Status:            statuses.TaskStateActive,
					},
				},
			},
		},
		{
			name: "Success_#6",
			payload: payload{
				isNoc: false,
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByPartnerAndManagedEndpointID(context.Background(), "partnerID", endpointID, 2).
						Return([]models.Task{
							{
								LastTaskInstanceID: lastTaskInstID1,
								State:              statuses.TaskStateActive,
								OriginID:           originID,
								Schedule: apiModels.Schedule{
									Regularity: apiModels.RunNow,
								},
							},
						}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				execResultPersistence: func() {
					execResultsPersistenceMock := mocks.NewMockExecutionResultPersistence(mockCtrlr)
					execResultsPersistenceMock.EXPECT().GetByTargetAndTaskInstanceIDs(endpointID, gomock.Any()).
						Return([]models.ExecutionResult{
							{
								TaskInstanceID:  lastTaskInstID1,
								ExecutionStatus: statuses.TaskInstanceSuccess,
							},
						}, nil).Times(1)
					models.ExecutionResultPersistenceInstance = execResultsPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					taskInstancePersistenceMock.EXPECT().GetByIDs(context.Background(), gomock.Any()).
						Return([]models.TaskInstance{
							{
								ID:       lastTaskInstID1,
								Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{endpointID: statuses.TaskInstanceRunning},
							},
						}, nil).Times(1)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				templateCache: func() {
					templateCacheMock := mocks.NewMockTemplateCache(mockCtrlr)
					templateCacheMock.EXPECT().GetByOriginID(context.Background(), "partnerID", originID, true).
						Return(models.TemplateDetails{}, errors.New("some_err")).Times(1)
					models.TemplateCacheInstance = templateCacheMock
				},
			},
			expected: expected{
				results: []*models.ExecutionResultView{
					{
						ManagedEndpointID: endpointID,
						ExecutionID:       lastTaskInstID1,
						OriginID:          originID,
						DeviceCount:       1,
						Regularity:        apiModels.RunNow,
						LastRunStatus:     statuses.TaskInstanceSuccess,
						Status:            statuses.TaskStateActive,
					},
				},
			},
		},
		{
			name: "error during getting execution results and task instances",
			payload: payload{
				isNoc: false,
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByPartnerAndManagedEndpointID(context.Background(), "partnerID", endpointID, 2).
						Return([]models.Task{
							{
								LastTaskInstanceID: lastTaskInstID1,
							},
							{
								LastTaskInstanceID: lastTaskInstID2,
							},
						}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				execResultPersistence: func() {
					execResultsPersistenceMock := mocks.NewMockExecutionResultPersistence(mockCtrlr)
					execResultsPersistenceMock.EXPECT().GetByTargetAndTaskInstanceIDs(endpointID, gomock.Any()).
						Return(nil, errors.New("some_err")).Times(1)
					models.ExecutionResultPersistenceInstance = execResultsPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					taskInstancePersistenceMock.EXPECT().GetByIDs(context.Background(), gomock.Any()).
						Return(nil, errors.New("some_err")).Times(1)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				templateCache: func() {
					models.TemplateCacheInstance = mocks.NewMockTemplateCache(mockCtrlr)
				},
			},
			expected: expected{
				results: []*models.ExecutionResultView{},
			},
		},
		{
			name: "error while getting list of tasks",
			payload: payload{
				isNoc: false,
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByPartnerAndManagedEndpointID(context.Background(), "partnerID", endpointID, 2).
						Return(nil, errors.New("some_err")).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				execResultPersistence: func() {
					execResultsPersistenceMock := mocks.NewMockExecutionResultPersistence(mockCtrlr)
					models.ExecutionResultPersistenceInstance = execResultsPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				templateCache: func() {
					models.TemplateCacheInstance = mocks.NewMockTemplateCache(mockCtrlr)
				},
			},
			expected: expected{
				results: []*models.ExecutionResultView{},
				err:     errors.New("some_err"),
			},
		},
	}

	for _, test := range tc {
		mockCtrlr = gomock.NewController(t)
		test.payload.taskPersistence()
		test.payload.execResultPersistence()
		test.payload.taskInstancePersistence()
		test.payload.templateCache()

		e := models.ExecutionResultViewRepoCassandra{}

		v, err := e.Get(context.Background(), "partnerID", endpointID, 2, test.payload.isNoc)
		mockCtrlr.Finish()

		if test.expected.err == nil {
			Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
			Ω(v).To(ConsistOf(test.expected.results), fmt.Sprintf(defaultMsg, test.name))
			continue
		}

		Ω(err.Error()).To(Equal(test.expected.err.Error()), fmt.Sprintf(defaultMsg, test.name))
		Ω(v).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
	}
}

func TestExecutionResultViewRepoCassandra_History(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrlr *gomock.Controller

	type payload struct {
		isNoc                   bool
		taskPersistence         func()
		taskInstancePersistence func()
		execResultPersistence   func()
		templateCache           func()
	}

	type expected struct {
		results []*models.ExecutionResultView
		err     error
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name: "error with getting tasks",
			payload: payload{
				isNoc: false,
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByIDAndManagedEndpoints(context.Background(), "partnerID", taskId, endpointID).
						Return(nil, errors.New("some_err")).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				execResultPersistence: func() {
					execResultsPersistenceMock := mocks.NewMockExecutionResultPersistence(mockCtrlr)
					models.ExecutionResultPersistenceInstance = execResultsPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				templateCache: func() {
					models.TemplateCacheInstance = mocks.NewMockTemplateCache(mockCtrlr)
				},
			},
			expected: expected{
				results: []*models.ExecutionResultView{},
				err:     errors.New("some_err"),
			},
		},
		{
			name: "error with nil tasks",
			payload: payload{
				isNoc: false,
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByIDAndManagedEndpoints(context.Background(), "partnerID", taskId, endpointID).
						Return([]models.Task{}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				execResultPersistence: func() {
					execResultsPersistenceMock := mocks.NewMockExecutionResultPersistence(mockCtrlr)
					models.ExecutionResultPersistenceInstance = execResultsPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				templateCache: func() {
					models.TemplateCacheInstance = mocks.NewMockTemplateCache(mockCtrlr)
				},
			},
			expected: expected{
				results: []*models.ExecutionResultView{},
				err:     errors.New("task not found, 55a98dda-f855-4197-a4b3-5bc72c70fd65 id"),
			},
		},
		{
			name: "Success_#2",
			payload: payload{
				isNoc: false,
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByIDAndManagedEndpoints(context.Background(), "partnerID", taskId, endpointID).
						Return([]models.Task{
							{
								IsRequireNOCAccess: true,
							},
						}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				execResultPersistence: func() {
					execResultsPersistenceMock := mocks.NewMockExecutionResultPersistence(mockCtrlr)
					models.ExecutionResultPersistenceInstance = execResultsPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				templateCache: func() {
					models.TemplateCacheInstance = mocks.NewMockTemplateCache(mockCtrlr)
				},
			},
			expected: expected{
				results: []*models.ExecutionResultView{},
			},
		},
		{
			name: "error with getting task instances",
			payload: payload{
				isNoc: false,
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByIDAndManagedEndpoints(context.Background(), "partnerID", taskId, endpointID).
						Return([]models.Task{
							{},
						}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				execResultPersistence: func() {
					execResultsPersistenceMock := mocks.NewMockExecutionResultPersistence(mockCtrlr)
					models.ExecutionResultPersistenceInstance = execResultsPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					taskInstancePersistenceMock.EXPECT().GetByTaskID(context.Background(), taskId).
						Return(nil, errors.New("some_err")).Times(1)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				templateCache: func() {
					models.TemplateCacheInstance = mocks.NewMockTemplateCache(mockCtrlr)
				},
			},
			expected: expected{
				results: []*models.ExecutionResultView{},
				err:     errors.New("some_err"),
			},
		},
		{
			name: "error with getting execution results",
			payload: payload{
				isNoc: false,
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByIDAndManagedEndpoints(context.Background(), "partnerID", taskId, endpointID).
						Return([]models.Task{
							{},
						}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				execResultPersistence: func() {
					execResultsPersistenceMock := mocks.NewMockExecutionResultPersistence(mockCtrlr)
					execResultsPersistenceMock.EXPECT().GetByTargetAndTaskInstanceIDs(endpointID, gomock.Any()).
						Return(nil, errors.New("some_err")).Times(1)
					models.ExecutionResultPersistenceInstance = execResultsPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					taskInstancePersistenceMock.EXPECT().GetByTaskID(context.Background(), taskId).
						Return([]models.TaskInstance{
							{
								ID:       lastTaskInstID1,
								Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{endpointID: statuses.TaskInstanceRunning},
							},
							{
								ID: lastTaskInstID2,
							},
						}, nil).Times(1)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				templateCache: func() {
					models.TemplateCacheInstance = mocks.NewMockTemplateCache(mockCtrlr)
				},
			},
			expected: expected{
				results: []*models.ExecutionResultView{},
				err: errors.New("ExecutionResultViewRepoCassandra.History: Error while trying to retrive " +
					"Script Execution Results by Task Instance ID and Managed Endpoint ID 55a98dda-f855-4197-a4b3-5bc72c70fd65. " +
					"Err: some_err"),
			},
		},
		{
			name: "Success_#3",
			payload: payload{
				isNoc: false,
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByIDAndManagedEndpoints(context.Background(), "partnerID", taskId, endpointID).
						Return([]models.Task{
							{
								OriginID: originID,
							},
						}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				execResultPersistence: func() {
					execResultsPersistenceMock := mocks.NewMockExecutionResultPersistence(mockCtrlr)
					execResultsPersistenceMock.EXPECT().GetByTargetAndTaskInstanceIDs(endpointID, gomock.Any()).
						Return([]models.ExecutionResult{
							{
								ManagedEndpointID: endpointID,
							},
						}, nil).Times(1)
					models.ExecutionResultPersistenceInstance = execResultsPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					taskInstancePersistenceMock.EXPECT().GetByTaskID(context.Background(), taskId).
						Return([]models.TaskInstance{
							{
								ID: lastTaskInstID1,
							},
						}, nil).Times(1)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				templateCache: func() {
					templateCacheMock := mocks.NewMockTemplateCache(mockCtrlr)
					templateCacheMock.EXPECT().GetByOriginID(context.Background(), "partnerID", originID, true).
						Return(models.TemplateDetails{}, nil).Times(1)
					models.TemplateCacheInstance = templateCacheMock
				},
			},
			expected: expected{
				results: []*models.ExecutionResultView{
					{
						ManagedEndpointID: endpointID,
						OriginID:          originID,
					},
				},
			},
		},
	}

	for _, test := range tc {
		mockCtrlr = gomock.NewController(t)
		test.payload.taskPersistence()
		test.payload.execResultPersistence()
		test.payload.taskInstancePersistence()
		test.payload.templateCache()

		e := models.ExecutionResultViewRepoCassandra{}

		history, err := e.History(context.Background(), "partnerID", taskId, endpointID, 1, test.payload.isNoc)
		mockCtrlr.Finish()

		if test.expected.err == nil {
			Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
			Ω(history).To(ConsistOf(test.expected.results), fmt.Sprintf(defaultMsg, test.name))
			continue
		}

		Ω(err.Error()).To(Equal(test.expected.err.Error()), fmt.Sprintf(defaultMsg, test.name))
		Ω(history).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
	}
}
