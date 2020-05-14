package taskExecutionResults

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	appLoader "gitlab.connectwisedev.com/platform/platform-tasking-service/src/app-loader"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

var (
	partnerID                   = "1"
	validManagedEndpoint        = "58a1af2f-6579-4aec-b45d-5dfde879ef01"
	validManagedEndpointUUID, _ = gocql.ParseUUID(validManagedEndpoint)
	validTaskingID              = "58a1af2f-6579-4aec-b45d-5dfde879ef01"
	validTaskingUUID, _         = gocql.ParseUUID(validTaskingID)
	validTaskInstanceID         = "58a1af2f-6579-4aec-b45d-000000000001"
	validTaskInstanceUUID, _    = gocql.ParseUUID(validTaskInstanceID)
	userServiceNoError          = user.NewMock("", "", "", "", false)
)

func init() {
	appLoader.LoadApplicationServices(true)
}

func TestTaskResultsServiceGet(t *testing.T) {
	executionResViewSlice := make([]*models.ExecutionResultView, 0)
	executionResViewSlice = append(executionResViewSlice, &models.ExecutionResultView{
		PartnerID:         partnerID,
		TaskID:            validTaskingUUID,
		TaskName:          "Task",
		ManagedEndpointID: validManagedEndpointUUID,
		LastRunStatus:     statuses.TaskInstanceSuccess,
	})

	testCases := []struct {
		name                 string
		execResViewMock      mock.ExecutionResultViewPersistenceConf
		userServiceMock      user.Service
		URL                  string
		method               string
		expectedErrorMessage string
		expectedCode         int
		expectedBody         []models.ExecutionResultView
	}{
		{
			name:                 "testCase 1 - Bad Request, bad Managed EndpointID",
			URL:                  "/" + partnerID + "/task-execution-results/managed-endpoints/badManagedEndpointID",
			method:               http.MethodGet,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorEndpointIDHasBadFormat,
		},
		{
			name:                 "testCase 2 -BadRequest badCount",
			URL:                  "/" + partnerID + "/task-execution-results/managed-endpoints/" + validManagedEndpoint + "?count=badCount",
			method:               http.MethodGet,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCountVarHasBadFormat,
		},
		{
			name: "testCase 4 - ServerError while getting exec results",
			URL:  "/" + partnerID + "/task-execution-results/managed-endpoints/" + validManagedEndpoint + "?count=2",
			execResViewMock: func(ev *mock.MockExecutionResultViewPersistence) *mock.MockExecutionResultViewPersistence {
				ev.EXPECT().
					Get(gomock.Any(), partnerID, validManagedEndpointUUID, 2, false).
					Return(nil, errors.New("err"))
				return ev
			},
			userServiceMock:      userServiceNoError,
			method:               http.MethodGet,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskExecutionResultsForManagedEndpoint,
		},
		{
			name: "testCase 5 - success",
			URL:  "/" + partnerID + "/task-execution-results/managed-endpoints/" + validManagedEndpoint + "?count=999",
			execResViewMock: func(ev *mock.MockExecutionResultViewPersistence) *mock.MockExecutionResultViewPersistence {
				ev.EXPECT().
					Get(gomock.Any(), partnerID, validManagedEndpointUUID, 999, false).
					Return(executionResViewSlice, nil)
				return ev
			},
			userServiceMock: userServiceNoError,
			method:          http.MethodGet,
			expectedCode:    http.StatusOK,
			expectedBody: []models.ExecutionResultView{
				{
					PartnerID:         partnerID,
					ManagedEndpointID: validManagedEndpointUUID,
					TaskID:            validTaskingUUID,
					TaskName:          "Task",
					LastRunStatus:     statuses.TaskInstanceSuccess,
				},
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedService(mockController, nil, tc.execResViewMock, tc.userServiceMock)
			router := getExecResultsRouter(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, nil)
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
			// expected body checking
			if tc.expectedErrorMessage == "" {
				expectedBody, err := json.Marshal(tc.expectedBody)
				if err != nil {
					t.Errorf("Cannot parse expected body")
				}

				if !bytes.Equal(expectedBody, w.Body.Bytes()) {
					t.Errorf("Wanted %v but got %v", string(expectedBody), w.Body.String())
				}
			}
		})
	}
}

func TestTaskResultsServiceHistory(t *testing.T) {
	executionResViewSlice := make([]*models.ExecutionResultView, 0)
	executionResViewSlice = append(executionResViewSlice, &models.ExecutionResultView{
		PartnerID:         partnerID,
		TaskID:            validTaskingUUID,
		TaskName:          "Task",
		ManagedEndpointID: validManagedEndpointUUID,
	})

	testCases := []struct {
		name                 string
		execResViewMock      mock.ExecutionResultViewPersistenceConf
		userServiceMock      user.Service
		URL                  string
		method               string
		expectedErrorMessage string
		expectedCode         int
		expectedBody         []models.ExecutionResultView
	}{
		{
			name:                 "testCase 1 - BadRequest  badTaskingID",
			URL:                  "/" + partnerID + "/task-execution-results/tasks/badTaskingID/managed-endpoints/badManagedEndpointID/history",
			method:               http.MethodGet,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTaskIDHasBadFormat,
		},
		{
			name:                 "testCase 2 - BadRequest  badManagedEndpointID",
			URL:                  "/" + partnerID + "/task-execution-results/tasks/" + validTaskingID + "/managed-endpoints/badManagedEndpointID/history",
			method:               http.MethodGet,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorEndpointIDHasBadFormat,
		},
		{
			name:                 "testCase 3 - BadRequest  bad count format",
			URL:                  "/" + partnerID + "/task-execution-results/tasks/" + validTaskingID + "/managed-endpoints/" + validManagedEndpoint + "/history?count=bad",
			method:               http.MethodGet,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCountVarHasBadFormat,
		},
		{
			name:            "testCase 5 - ServerError while getting history,server err, count=999",
			URL:             "/" + partnerID + "/task-execution-results/tasks/" + validTaskingID + "/managed-endpoints/" + validManagedEndpoint + "/history?count=999",
			userServiceMock: userServiceNoError,
			execResViewMock: func(ev *mock.MockExecutionResultViewPersistence) *mock.MockExecutionResultViewPersistence {
				ev.EXPECT().
					History(gomock.Any(), partnerID, validTaskingUUID, validManagedEndpointUUID, 999, false).
					Return(nil, errors.New("err"))
				return ev
			},
			method:               http.MethodGet,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskExecutionResults,
		},
		{
			name:            "testCase 6 - error while getting history,TaskNotFoundError, count=999",
			URL:             "/" + partnerID + "/task-execution-results/tasks/" + validTaskingID + "/managed-endpoints/" + validManagedEndpoint + "/history?count=999",
			userServiceMock: userServiceNoError,
			execResViewMock: func(ev *mock.MockExecutionResultViewPersistence) *mock.MockExecutionResultViewPersistence {
				ev.EXPECT().
					History(gomock.Any(), partnerID, validTaskingUUID, validManagedEndpointUUID, 999, false).
					Return(nil, models.TaskNotFoundError{})
				return ev
			},
			method:               http.MethodGet,
			expectedCode:         http.StatusNotFound,
			expectedErrorMessage: errorcode.ErrorCantGetTaskExecutionResults,
		},
		{
			name:            "testCase 7 - successes, no count specified",
			URL:             "/" + partnerID + "/task-execution-results/tasks/" + validTaskingID + "/managed-endpoints/" + validManagedEndpoint + "/history",
			userServiceMock: userServiceNoError,
			execResViewMock: func(ev *mock.MockExecutionResultViewPersistence) *mock.MockExecutionResultViewPersistence {
				ev.EXPECT().
					History(gomock.Any(), partnerID, validTaskingUUID, validManagedEndpointUUID, defaultHistoryCount, false).
					Return(executionResViewSlice, nil)
				return ev
			},
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: []models.ExecutionResultView{
				{PartnerID: partnerID, ManagedEndpointID: validManagedEndpointUUID, TaskID: validTaskingUUID, TaskName: "Task"},
			},
		},
		{
			name:            "testCase 8 - successes, with count specified",
			URL:             "/" + partnerID + "/task-execution-results/tasks/" + validTaskingID + "/managed-endpoints/" + validManagedEndpoint + "/history?count=999",
			userServiceMock: userServiceNoError,
			execResViewMock: func(ev *mock.MockExecutionResultViewPersistence) *mock.MockExecutionResultViewPersistence {
				ev.EXPECT().
					History(gomock.Any(), partnerID, validTaskingUUID, validManagedEndpointUUID, 999, false).
					Return(executionResViewSlice, nil)
				return ev
			},
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: []models.ExecutionResultView{
				{PartnerID: partnerID, ManagedEndpointID: validManagedEndpointUUID, TaskID: validTaskingUUID, TaskName: "Task"},
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedService(mockController, nil, tc.execResViewMock, tc.userServiceMock)
			router := getExecResultsRouter(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, nil)
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
			// expected body checking
			if tc.expectedErrorMessage == "" {
				expectedBody, err := json.Marshal(tc.expectedBody)
				if err != nil {
					t.Errorf("Cannot parse expected body")
				}

				if !bytes.Equal(expectedBody, w.Body.Bytes()) {
					t.Errorf("Wanted %v but got %v", string(expectedBody), w.Body.String())
				}
			}
		})
	}
}

func TestGetTaskExecutionStdoutLogs_GetTaskExecutionStderrLogs(t *testing.T) {

	testCases := []struct {
		name                   string
		userSitesMock          user.Service
		execResPersistenceMock mock.ExecutionResultPersistenceConf
		URL                    string
		expectedErrorMessage   string
		method                 string
		expectedCode           int
		expectedResult         string
	}{
		{
			name:                 "testCase 1 - BadRequest  badManagedEndpointID in URL",
			URL:                  "/" + partnerID + "/task-execution-results/managed-endpoints/badManagedEndpoint/task-instances/" + validTaskInstanceID + "/logs/stdout",
			method:               http.MethodGet,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorEndpointIDHasBadFormat,
		},
		{
			name:                 "testCase 2 - BadRequest  badTaskInstanceID in URL",
			URL:                  "/" + partnerID + "/task-execution-results/managed-endpoints/" + validManagedEndpoint + "/task-instances/badTaskInstanceID/logs/stdout",
			method:               http.MethodGet,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTaskInstanceIDHasBadFormat,
		},
		{
			name: "testCase 3 - server error from GetByTargetAndTaskInstanceIDs",
			URL:  "/" + partnerID + "/task-execution-results/managed-endpoints/" + validManagedEndpoint + "/task-instances/" + validTaskInstanceID + "/logs/stderr",
			execResPersistenceMock: func(ep *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				ep.EXPECT().
					GetByTargetAndTaskInstanceIDs(validManagedEndpointUUID, validTaskInstanceUUID).
					Return(nil, errors.New(""))
				return ep
			},
			method:               http.MethodGet,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantDecodeInputData,
		},
		{
			name: "testCase 4 - Successes getting stderr",
			URL:  "/" + partnerID + "/task-execution-results/managed-endpoints/" + validManagedEndpoint + "/task-instances/" + validTaskInstanceID + "/logs/stderr",
			execResPersistenceMock: func(ep *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				ep.EXPECT().
					GetByTargetAndTaskInstanceIDs(validManagedEndpointUUID, validTaskInstanceUUID).
					Return([]models.ExecutionResult{{StdErr: "stdErr"}}, nil)
				return ep
			},
			method:         http.MethodGet,
			expectedCode:   http.StatusOK,
			expectedResult: "\"stdErr\"",
		},
		{
			name: "testCase 5 - Successes getting stdout",
			URL:  "/" + partnerID + "/task-execution-results/managed-endpoints/" + validManagedEndpoint + "/task-instances/" + validTaskInstanceID + "/logs/stdout",
			execResPersistenceMock: func(ep *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				ep.EXPECT().
					GetByTargetAndTaskInstanceIDs(validManagedEndpointUUID, validTaskInstanceUUID).
					Return([]models.ExecutionResult{{StdOut: "stdOut"}}, nil)
				return ep
			},
			method:         http.MethodGet,
			expectedCode:   http.StatusOK,
			expectedResult: "\"stdOut\"",
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedService(mockController, tc.execResPersistenceMock, nil, tc.userSitesMock)
			router := getExecResultsRouter(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, nil)
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}

			if tc.expectedResult != w.Body.String() && tc.expectedResult != "" {
				t.Errorf("Wanted %v but got %v", tc.expectedResult, w.Body.String())
			}
		})
	}
}

func getExecResultsRouter(service TaskResultsService) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/{partnerID}/task-execution-results/managed-endpoints/{managedEndpointID}", http.HandlerFunc(service.Get)).Methods(http.MethodGet)
	router.HandleFunc("/{partnerID}/task-execution-results/tasks/{taskID}/managed-endpoints/{managedEndpointID}/history", http.HandlerFunc(service.History)).Methods(http.MethodGet)
	router.HandleFunc("/{partnerID}/task-execution-results/managed-endpoints/{managedEndpointID}/task-instances/{taskInstanceID}/logs/stdout", http.HandlerFunc(service.GetTaskExecutionStdoutLogs)).Methods(http.MethodGet)
	router.HandleFunc("/{partnerID}/task-execution-results/managed-endpoints/{managedEndpointID}/task-instances/{taskInstanceID}/logs/stderr", http.HandlerFunc(service.GetTaskExecutionStderrLogs)).Methods(http.MethodGet)
	return
}

// getMockedService is a function that returns configured mocked TaskDefinitionService
func getMockedService(mockController *gomock.Controller,
	ep mock.ExecutionResultPersistenceConf,
	ev mock.ExecutionResultViewPersistenceConf,
	us user.Service) TaskResultsService {

	execResultMock := mock.NewMockExecutionResultPersistence(mockController)
	execResultViewMock := mock.NewMockExecutionResultViewPersistence(mockController)
	if ep != nil {
		execResultMock = ep(execResultMock)
	}
	if ev != nil {
		execResultViewMock = ev(execResultViewMock)
	}
	return NewTaskResultsService(execResultViewMock, execResultMock, us, http.DefaultClient)
}
