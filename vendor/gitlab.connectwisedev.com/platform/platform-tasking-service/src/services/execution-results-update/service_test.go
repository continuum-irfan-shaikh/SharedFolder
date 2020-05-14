package executionResultsUpdate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"golang.org/x/net/context"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/asset"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-repository"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"
)

const defaultMsg = `failed on unexpected value of result "%v"`

var (
	validManagedEndpoint        = "58a1af2f-6579-4aec-b45d-5dfde879ef01"
	validManagedEndpointUUID, _ = gocql.ParseUUID(validManagedEndpoint)
	validTaskInstanceID         = "58a1af2f-6579-4aec-b45d-000000000001"
	validTaskInstanceUUID, _    = gocql.ParseUUID(validTaskInstanceID)
	validOriginID               = "804823ce-b84d-11e9-834d-e4e74935e7a7"
	validOriginUUID, _          = gocql.ParseUUID(validOriginID)
	partnerID                   = "0"
	someTime                    = time.Now().UTC()
	userHasNOCAccess            = true
	updateTaskURL               = "/" + partnerID + "/task-execution-results/task-instances/" + validTaskInstanceID
	validBodySuccesesStatus     = tasking.ExecutionResult{
		EndpointID:       validManagedEndpoint,
		UpdateTime:       someTime,
		CompletionStatus: "Success",
	}
	validBodyFailedStatus = tasking.ExecutionResult{
		EndpointID:       validManagedEndpoint,
		UpdateTime:       someTime,
		CompletionStatus: "Failed",
	}
	validBodyFailedByTimeoutStatus = tasking.ExecutionResult{
		EndpointID:       validManagedEndpoint,
		UpdateTime:       someTime,
		CompletionStatus: "Failed",
		ErrorDetails:     failedByTimeout,
	}
)

func init() {
	logger.Load(config.Config.Log)
	translation.Load()
}

func Test_UpdateTaskAndTaskInstanceStatuses(t *testing.T) {
	var (
		taskInstancePersMockNoErr = func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
			ti.EXPECT().
				GetMinimalInstanceByID(gomock.Any(), validTaskInstanceUUID).Return(models.TaskInstance{
				OriginID: validManagedEndpointUUID,
				ID:       validTaskInstanceUUID,
				TaskID:   validTaskInstanceUUID,
			}, nil)
			return ti
		}
	)

	testCases := []struct {
		name                           string
		executionResultPersistenceMock mock.ExecutionResultPersistenceConf
		taskInstancePersistenceMock    mock.TaskInstancePersistenceConf
		taskSummaryPersistenceMock     mock.TaskSummaryPersistenceConf
		taskPersistenceMock            mock.TaskPersistenceConf
		templateCacheMock              mock.TemplateCacheConf
		URL                            string
		method                         string
		expectedCode                   int
		expectedErrorMessage           string
		bodyToSend                     tasking.ExecutionResult
	}{
		{
			name:                 "testCase 1 - bad TaskInstance ID",
			URL:                  "/" + partnerID + "/task-execution-results/task-instances/badTaskInstanceID",
			method:               http.MethodPost,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTaskInstanceIDHasBadFormat,
		},
		{
			name: "testCase 2 - bad ExecutionResult body, no CompletionStatus",
			URL:  updateTaskURL,
			bodyToSend: tasking.ExecutionResult{
				EndpointID: validManagedEndpoint,
				UpdateTime: someTime,
			},
			method:               http.MethodPost,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantDecodeInputData,
		},
		{
			name: "testCase 3 - bad EndpointID in body, cannot parse struct body",
			URL:  updateTaskURL,
			bodyToSend: tasking.ExecutionResult{
				EndpointID:       "58a1af2f-6579-4~~c-b45d-{dfde879ef/1",
				UpdateTime:       someTime,
				CompletionStatus: "Successes",
			},
			method:               http.MethodPost,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantDecodeInputData,
		},
		{
			name: "testCase 4 - cannot get Task Instances",
			URL:  updateTaskURL,
			taskInstancePersistenceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetMinimalInstanceByID(gomock.Any(), validTaskInstanceUUID).
					Return(models.TaskInstance{}, errors.New("err"))
				return ti
			},
			bodyToSend: tasking.ExecutionResult{
				EndpointID:       validManagedEndpoint,
				UpdateTime:       someTime,
				CompletionStatus: "bad",
			},
			method:               http.MethodPost,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantGetTaskInstances,
		},
		{
			name: "testCase 5 - len of Task Instances is 0",
			URL:  updateTaskURL,
			taskInstancePersistenceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetMinimalInstanceByID(gomock.Any(), validTaskInstanceUUID).
					Return(models.TaskInstance{}, nil)
				return ti
			},
			bodyToSend:           validBodySuccesesStatus,
			method:               http.MethodPost,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantGetTaskInstances,
		},
		{
			name: "testCase 6 - bad completionStatus",
			URL:  updateTaskURL,
			bodyToSend: tasking.ExecutionResult{
				EndpointID:       validManagedEndpoint,
				UpdateTime:       someTime,
				CompletionStatus: `bad`,
			},
			taskInstancePersistenceMock: taskInstancePersMockNoErr,
			method:                      http.MethodPost,
			expectedCode:                http.StatusCreated,
		},
		{
			name:                        "testCase 7 - error from Server database, cannot upsert results, body with Successes Status",
			URL:                         updateTaskURL,
			bodyToSend:                  validBodySuccesesStatus,
			taskInstancePersistenceMock: taskInstancePersMockNoErr,
			executionResultPersistenceMock: func(ep *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				ep.EXPECT().
					Upsert(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("error"))
				return ep
			},
			method:               http.MethodPost,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantInsertData,
		},
		{
			name:                        "testCase 8 - error from Server database, cannot upsert results, body with Failed Status",
			URL:                         updateTaskURL,
			bodyToSend:                  validBodyFailedStatus,
			taskInstancePersistenceMock: taskInstancePersMockNoErr,
			executionResultPersistenceMock: func(ep *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				ep.EXPECT().
					Upsert(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("error"))
				return ep
			},
			method:               http.MethodPost,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantInsertData,
		},
		{
			name:       "testCase 9 - error from Server database, cannot insert into taskInstance",
			URL:        updateTaskURL,
			bodyToSend: validBodyFailedStatus,
			taskInstancePersistenceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetMinimalInstanceByID(gomock.Any(), validTaskInstanceUUID).Return(models.TaskInstance{
					OriginID: validManagedEndpointUUID,
					ID:       validTaskInstanceUUID,
					TaskID:   validTaskInstanceUUID,
				}, nil)
				ti.EXPECT().
					UpdateStatuses(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return ti
			},
			executionResultPersistenceMock: func(ep *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				ep.EXPECT().
					Upsert(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
				return ep
			},
			method:               http.MethodPost,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantUpdateTaskInstances,
		},
		{
			name:       "testCas 12 - Updated",
			URL:        updateTaskURL,
			bodyToSend: validBodyFailedStatus,
			taskInstancePersistenceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetMinimalInstanceByID(gomock.Any(), validTaskInstanceUUID).Return(models.TaskInstance{
					OriginID: validManagedEndpointUUID,
					ID:       validTaskInstanceUUID,
					TaskID:   validTaskInstanceUUID,
				}, nil)
				ti.EXPECT().
					UpdateStatuses(gomock.Any(), gomock.Any())
				return ti
			},
			executionResultPersistenceMock: func(ep *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				ep.EXPECT().
					Upsert(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
				return ep
			},
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetExecutionResultTaskData(partnerID, gomock.Any(), gomock.Any()).
					Return(models.ExecutionResultTaskData{Name: "name", CreatedBy: "someBody"}, nil)
				return tp
			},
			method:       http.MethodPost,
			expectedCode: http.StatusCreated,
		},
		{
			name:       "testCas 12.1 - Updated for timeout error",
			URL:        updateTaskURL,
			bodyToSend: validBodyFailedByTimeoutStatus,
			taskInstancePersistenceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByIDs(gomock.Any(), validTaskInstanceUUID).Return([]models.TaskInstance{
					{
						OriginID: validManagedEndpointUUID,
						ID:       validTaskInstanceUUID,
						TaskID:   validTaskInstanceUUID,
					}, {
						OriginID: validManagedEndpointUUID,
						ID:       validTaskInstanceUUID,
						TaskID:   validTaskInstanceUUID,
					}}, nil)
				ti.EXPECT().
					UpdateStatuses(gomock.Any(), gomock.Any())
				return ti
			},
			executionResultPersistenceMock: func(ep *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				ep.EXPECT().
					Upsert(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
				return ep
			},
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetExecutionResultTaskData(partnerID, gomock.Any(), gomock.Any()).
					Return(models.ExecutionResultTaskData{Name: "name", CreatedBy: "someBody"}, nil)
				return tp
			},
			method:       http.MethodPost,
			expectedCode: http.StatusCreated,
		},
		{
			name:       "testCas 12.5 - Nothing to update because TI device status - failed",
			URL:        updateTaskURL,
			bodyToSend: validBodyFailedStatus,
			taskInstancePersistenceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetMinimalInstanceByID(gomock.Any(), validTaskInstanceUUID).Return(models.TaskInstance{
					OriginID: validManagedEndpointUUID,
					ID:       validTaskInstanceUUID,
					TaskID:   validTaskInstanceUUID,
					Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{
						validManagedEndpointUUID: statuses.TaskInstanceFailed,
					},
				}, nil)
				return ti
			},
			method:       http.MethodPost,
			expectedCode: http.StatusCreated,
		},
		{
			name:       "testCas 13 - Updated but with cache error",
			URL:        updateTaskURL,
			bodyToSend: validBodyFailedStatus,
			taskInstancePersistenceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetMinimalInstanceByID(gomock.Any(), validTaskInstanceUUID).Return(models.TaskInstance{
					OriginID: validManagedEndpointUUID,
					ID:       validTaskInstanceUUID,
					TaskID:   validTaskInstanceUUID,
				}, nil)
				ti.EXPECT().
					UpdateStatuses(gomock.Any(), gomock.Any())
				return ti
			},
			executionResultPersistenceMock: func(ep *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				ep.EXPECT().
					Upsert(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
				return ep
			},
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetExecutionResultTaskData(partnerID, gomock.Any(), gomock.Any()).
					Return(models.ExecutionResultTaskData{
						Name:          "name3",
						ResultWebHook: "notEmpty",
					}, nil)
				return tp
			},
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validManagedEndpointUUID, userHasNOCAccess).
					Return(models.TemplateDetails{}, errors.New("")).Times(1)
				return tc
			},
			method:       http.MethodPost,
			expectedCode: http.StatusCreated,
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedService(
				mockController,
				tc.executionResultPersistenceMock,
				tc.taskInstancePersistenceMock,
				tc.taskSummaryPersistenceMock,
				tc.taskPersistenceMock,
				tc.templateCacheMock,
			)
			service = getMockedServiceWithHistory(mockController, service)
			router := getExecResultsUpdateRouter(service)

			w := httptest.NewRecorder()

			payload := []tasking.ExecutionResult{tc.bodyToSend}
			body, err := json.Marshal(payload)
			if err != nil {
				t.Errorf("Cannot parse body err=%v", err)
			}

			r := httptest.NewRequest(tc.method, tc.URL, bytes.NewReader(body))
			router.ServeHTTP(w, r)
			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
		})
	}
}

func getMockedServiceWithHistory(mc *gomock.Controller, s ExecutionResultUpdateService) ExecutionResultUpdateService {
	taskExecHistoryRepoMock := mockrepositories.NewMockTaskExecutionHistoryRepo(mc)
	taskExecHistoryRepoMock.EXPECT().Insert(gomock.Any()).Return(nil).MaxTimes(1)

	assetRepoMock := mock.NewMockAsset(mc)
	assetRepoMock.EXPECT().GetSiteIDByEndpointID(gomock.Any(), partnerID, validManagedEndpointUUID).Return("siteID", "", nil).MaxTimes(1)
	assetRepoMock.EXPECT().GetMachineNameByEndpointID(gomock.Any(), partnerID, validManagedEndpointUUID).Return("someName", nil).MaxTimes(1)
	asset.ServiceInstance = assetRepoMock

	s.assetService = assetRepoMock
	s.taskExecHistory = taskExecHistoryRepoMock
	return s
}

// getMockedService is a function that returns configured mocked TaskDefinitionService
func getMockedService(mc *gomock.Controller,
	ep mock.ExecutionResultPersistenceConf,
	ti mock.TaskInstancePersistenceConf,
	ts mock.TaskSummaryPersistenceConf,
	tp mock.TaskPersistenceConf,
	tc mock.TemplateCacheConf,
) ExecutionResultUpdateService {

	execResultMock := mock.NewMockExecutionResultPersistence(mc)
	if ep != nil {
		execResultMock = ep(execResultMock)
	}
	taskInstancePersMock := mock.NewMockTaskInstancePersistence(mc)
	if ti != nil {
		taskInstancePersMock = ti(taskInstancePersMock)
		models.TaskInstancePersistenceInstance = taskInstancePersMock
	}
	taskSummaryMock := mock.NewMockTaskSummaryPersistence(mc)
	if ts != nil {
		taskSummaryMock = ts(taskSummaryMock)
		models.TaskSummaryPersistenceInstance = taskSummaryMock
	}
	taskPersistenceMock := mock.NewMockTaskPersistence(mc)
	if tp != nil {
		taskPersistenceMock = tp(taskPersistenceMock)
		models.TaskPersistenceInstance = taskPersistenceMock
	}
	templateCacheMock := mock.NewMockTemplateCache(mc)
	if tc != nil {
		templateCacheMock = tc(templateCacheMock)
		models.TemplateCacheInstance = templateCacheMock
	}

	return NewExecutionResultUpdateService(execResultMock, taskPersistenceMock, taskInstancePersMock, nil, nil)
}

func getExecResultsUpdateRouter(service ExecutionResultUpdateService) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/{partnerID}/task-execution-results/task-instances/{taskInstanceID}", http.HandlerFunc(service.UpdateTaskAndTaskInstanceStatuses)).Methods(http.MethodPost)
	return
}

func TestCreateTaskExecHistory(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrl *gomock.Controller

	type payload struct {
		pData           *processData
		result          models.ExecutionResult
		asset           func() integration.Asset
		taskExecHistory func() TaskExecutionHistoryRepo
	}

	type expected struct {
		err string
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name: "Not enough process execution results",
			expected: expected{
				err: "not enough data to process task execution history",
			},
			payload: payload{
				pData:  nil,
				result: models.ExecutionResult{},
				asset: func() integration.Asset {
					assetMock := mock.NewMockAsset(mockCtrl)
					return assetMock
				},
				taskExecHistory: func() TaskExecutionHistoryRepo {
					taskExecHistoryMock := mockrepositories.NewMockTaskExecutionHistoryRepo(mockCtrl)
					return taskExecHistoryMock
				},
			},
		},
		{
			name: "error with getting siteID from asset service",
			expected: expected{
				err: `error while getting siteID from asset service: some_err`,
			},
			payload: payload{
				pData: &processData{
					PartnerID:     "partnerID",
					ResultWebhook: partnerID,
				},
				result: models.ExecutionResult{
					ManagedEndpointID: validManagedEndpointUUID,
				},
				asset: func() integration.Asset {
					assetMock := mock.NewMockAsset(mockCtrl)
					assetMock.EXPECT().GetSiteIDByEndpointID(gomock.Any(), "partnerID", validManagedEndpointUUID).
						Return("", "", errors.New("some_err")).Times(1)
					return assetMock
				},
				taskExecHistory: func() TaskExecutionHistoryRepo {
					taskExecHistoryMock := mockrepositories.NewMockTaskExecutionHistoryRepo(mockCtrl)
					return taskExecHistoryMock
				},
			},
		},
		{
			name: "error with getting machine name from asset service",
			expected: expected{
				err: `error while getting machineName from asset service: some_err`,
			},
			payload: payload{
				pData: &processData{
					PartnerID:     "partnerID",
					ResultWebhook: partnerID,
				},
				result: models.ExecutionResult{
					ManagedEndpointID: validManagedEndpointUUID,
				},
				asset: func() integration.Asset {
					assetMock := mock.NewMockAsset(mockCtrl)
					assetMock.EXPECT().GetSiteIDByEndpointID(gomock.Any(), "partnerID", validManagedEndpointUUID).
						Return("50000031", "", nil).Times(1)
					assetMock.EXPECT().GetMachineNameByEndpointID(gomock.Any(), "partnerID", validManagedEndpointUUID).
						Return("", errors.New("some_err")).Times(1)
					return assetMock
				},
				taskExecHistory: func() TaskExecutionHistoryRepo {
					taskExecHistoryMock := mockrepositories.NewMockTaskExecutionHistoryRepo(mockCtrl)
					return taskExecHistoryMock
				},
			},
		},
		{
			name: "error with inserting data to DB",
			expected: expected{
				err: `error while inserting history: some_err`,
			},
			payload: payload{
				pData: &processData{
					PartnerID: "partnerID",
					TaskInstance: models.TaskInstance{
						StartedAt: parseTime("2019-08-06T16:10:41+00:00"),
						OriginID:  validOriginUUID,
					},
					Name:      "taskName",
					CreatedBy: "someUser",
				},
				result: models.ExecutionResult{
					ManagedEndpointID: validManagedEndpointUUID,
					ExecutionStatus:   statuses.TaskInstanceStatus(2),
					UpdatedAt:         parseTime("2019-08-15T16:10:41+00:00"),
				},
				asset: func() integration.Asset {
					assetMock := mock.NewMockAsset(mockCtrl)
					assetMock.EXPECT().GetSiteIDByEndpointID(gomock.Any(), "partnerID", validManagedEndpointUUID).
						Return("50000031", "", nil).Times(1)
					assetMock.EXPECT().GetMachineNameByEndpointID(gomock.Any(), "partnerID", validManagedEndpointUUID).
						Return("machineName", nil).Times(1)
					return assetMock
				},
				taskExecHistory: func() TaskExecutionHistoryRepo {
					taskExecHistoryMock := mockrepositories.NewMockTaskExecutionHistoryRepo(mockCtrl)

					t := entities.TaskExecHistory{
						ExecYear:      "2019",
						ExecMonth:     "08",
						ExecDate:      "2019-08-06",
						ExecTime:      parseTime("2019-08-06T16:10:41+00:00"),
						EndpointID:    "58a1af2f-6579-4aec-b45d-5dfde879ef01",
						ScriptName:    "taskName",
						ScriptID:      "804823ce-b84d-11e9-834d-e4e74935e7a7",
						CompletedTime: parseTime("2019-08-15T16:10:41+00:00"),
						ExecStatus:    "Success",
						PartnerID:     "partnerID",
						SiteID:        "50000031",
						MachineName:   "machineName",
						ExecBy:        "someUser",
					}
					taskExecHistoryMock.EXPECT().Insert(t).
						Return(errors.New("some_err")).Times(1)

					return taskExecHistoryMock
				},
			},
		},
		{
			name: "success",
			expected: expected{
				err: "",
			},
			payload: payload{
				pData: &processData{
					PartnerID: "partnerID",
					TaskInstance: models.TaskInstance{
						StartedAt: parseTime("2019-08-06T16:10:41+00:00"),
						OriginID:  validOriginUUID,
					},
					Name:      "taskName",
					CreatedBy: "someUser",
				},
				result: models.ExecutionResult{
					ManagedEndpointID: validManagedEndpointUUID,
					ExecutionStatus:   statuses.TaskInstanceStatus(2),
					UpdatedAt:         parseTime("2019-08-15T16:10:41+00:00"),
				},
				asset: func() integration.Asset {
					assetMock := mock.NewMockAsset(mockCtrl)
					assetMock.EXPECT().GetSiteIDByEndpointID(gomock.Any(), "partnerID", validManagedEndpointUUID).
						Return("50000031", "", nil).Times(1)
					assetMock.EXPECT().GetMachineNameByEndpointID(gomock.Any(), "partnerID", validManagedEndpointUUID).
						Return("machineName", nil).Times(1)
					return assetMock
				},
				taskExecHistory: func() TaskExecutionHistoryRepo {
					taskExecHistoryMock := mockrepositories.NewMockTaskExecutionHistoryRepo(mockCtrl)

					t := entities.TaskExecHistory{
						ExecYear:      "2019",
						ExecMonth:     "08",
						ExecDate:      "2019-08-06",
						ExecTime:      parseTime("2019-08-06T16:10:41+00:00"),
						EndpointID:    "58a1af2f-6579-4aec-b45d-5dfde879ef01",
						ScriptName:    "taskName",
						ScriptID:      "804823ce-b84d-11e9-834d-e4e74935e7a7",
						CompletedTime: parseTime("2019-08-15T16:10:41+00:00"),
						ExecStatus:    "Success",
						PartnerID:     "partnerID",
						SiteID:        "50000031",
						MachineName:   "machineName",
						ExecBy:        "someUser",
					}
					taskExecHistoryMock.EXPECT().Insert(t).
						Return(nil).Times(1)

					return taskExecHistoryMock
				},
			},
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			mockCtrl = gomock.NewController(t)
			s := NewExecutionResultUpdateService(nil, nil, nil, test.payload.taskExecHistory(), test.payload.asset())

			err := s.createTaskExecHistory(context.Background(), test.payload.pData, test.payload.result)
			mockCtrl.Finish()

			if test.expected.err == "" {
				Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
				return
			}
			Ω(err.Error()).To(Equal(test.expected.err), fmt.Sprintf(defaultMsg, test.name))
		})

	}
}

func parseTime(v string) time.Time {
	t, _ := time.Parse(time.RFC3339, v)
	return t
}
