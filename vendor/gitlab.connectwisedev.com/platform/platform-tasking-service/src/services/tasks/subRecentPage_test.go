package tasks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const timeFormat = "2009-11-10 23:00:00 +0000 UTC"

func TestSubRecent(t *testing.T) {

	var (
		currentTime                  = time.Now()
		currentTimePlusMinute        = currentTime.Add(time.Minute)
		currentTimePlusTwentyMinutes = currentTime.Add(time.Minute * 20)

		targetsStrSlice = []string{validUUIDstr}
		targets         = models.Target{IDs: targetsStrSlice}

		validTasksSubRecent = []models.Task{{
			Name:               "task1",
			ID:                 validUUID,
			IsRequireNOCAccess: false,
			Targets:            targets,
			OriginID:           validUUID,
			Type:               models.TaskTypeScript,
			RunTimeUTC:         currentTimePlusTwentyMinutes,
			ManagedEndpointID:  validUUID,
		}, {
			Name:               "task1",
			ID:                 validUUID,
			IsRequireNOCAccess: true,
			Targets:            targets,
			RunTimeUTC:         currentTimePlusMinute, // nearest RunTime
			ManagedEndpointID:  validUUID2,
			Type:               models.TaskTypeScript,
		}}

		validTasksSubRecentNotScript = []models.Task{
			{Name: "task1", ID: validUUID, ExternalTask: true, IsRequireNOCAccess: false, Targets: targets, OriginID: validUUID},
			{Name: "task2", ID: validUUID, IsRequireNOCAccess: true, Targets: targets},
		}

		validTaskInstancesSubRecent = []models.TaskInstance{
			{ID: validUUID, TaskID: validUUID, Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{validUUID: statuses.TaskInstanceSuccess}},
			{ID: validUUID2, TaskID: validUUID, Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{validUUID2: statuses.TaskInstanceFailed}},
			{ID: validUUID3, TaskID: validUUID},
			{ID: validUUID4, TaskID: validUUID, Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{validUUID4: statuses.TaskInstanceScheduled}},
		}
	)

	userMock := user.NewMock("", partnerID, "", "", true)

	testCases := []struct {
		name                       string
		URL                        string
		taskSummaryMock            mock.TaskSummaryPersistenceConf
		taskInstanceMock           mock.TaskInstancePersistenceConf
		taskPersistenceMock        mock.TaskPersistenceConf
		templateCacheMock          mock.TemplateCacheConf
		resultMock                 mock.ExecutionResultPersistenceConf
		method                     string
		expectedCode               int
		expectedErrorMessage       string
		isNeedToCheckBody          bool
		expectedNearestNextRunTime time.Time
		expectedTaskID             gocql.UUID
	}{
		{
			name:                 "testCase 0 - Cannot parse UUID in request",
			URL:                  "/" + partnerID + "/tasks/BADUUID/data",
			method:               http.MethodGet,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantDecodeInputData,
		},
		{
			name:   "testCase 1 - Error from cassandra while getting task by id",
			URL:    "/" + partnerID + "/tasks/" + validUUIDstr + "/data",
			method: http.MethodGet,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, validUUID).
					Return([]models.Task{}, errors.New("err"))
				return tp
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
		},
		{
			name:   "testCase 2 - Cassandra returned 0 tasks",
			URL:    "/" + partnerID + "/tasks/" + validUUIDstr + "/data",
			method: http.MethodGet,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, validUUID).
					Return([]models.Task{}, nil)
				return tp
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:   "testCase 3 - Error while getting taskInstances, also task type is not script",
			URL:    "/" + partnerID + "/tasks/" + validUUIDstr + "/data",
			method: http.MethodGet,
			// first we get tasks by ids, then in goroutine trying to get TIs but get error, also in in templateCache goroutine
			// see that task type is not script and return taskTemplate as failed
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, validUUID).
					Return(validTasksSubRecentNotScript, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByTaskID(gomock.Any(), validUUID).
					Return([]models.TaskInstance{}, errors.New("err"))
				return ti
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskInstances,
		},
		{
			name:   "testCase 4 - Error while getting templateDetails",
			URL:    "/" + partnerID + "/tasks/" + validUUIDstr + "/data",
			method: http.MethodGet,
			// first we get tasks by ids, then in goroutine trying to get TIs but get error, also in templateCache goroutine
			// see that task type is script an then get templateDetaills but get err
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, validUUID).
					Return(validTasksSubRecent, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByTaskID(gomock.Any(), validUUID).
					Return(validTaskInstancesSubRecent, nil)
				return ti
			},
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validTasksSubRecent[0].OriginID, true).
					Return(models.TemplateDetails{}, errors.New("err"))
				return tc
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskDefinitionTemplate,
		},
		{
			name:   "testCase 5 - Error while getting executionResults, with empty data",
			URL:    "/" + partnerID + "/tasks/" + validUUIDstr + "/data",
			method: http.MethodGet,
			// first we get tasks by ids, then in goroutine trying to get TIs but get error, also in templateCache goroutine
			// see that task type is script an then get templateDetail's, then get executionResults but get err
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, validUUID).
					Return(validTasksSubRecent, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByTaskID(gomock.Any(), validUUID).
					Return(validTaskInstancesSubRecent, nil)
				return ti
			},
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validTasksSubRecent[0].OriginID, true).
					Return(models.TemplateDetails{}, nil)
				return tc
			},
			resultMock: func(er *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				er.EXPECT().
					GetByTaskInstanceIDs(gomock.Any()).
					Return([]models.ExecutionResult{}, errors.New("err"))
				return er
			},
			expectedCode:               http.StatusOK,
			isNeedToCheckBody:          true,
			expectedNearestNextRunTime: validTasksSubRecent[1].RunTimeUTC,
			expectedTaskID:             validTasksSubRecent[0].ID,
		},
		{
			name:   "testCase 6 - Ok",
			URL:    "/" + partnerID + "/tasks/" + validUUIDstr + "/data",
			method: http.MethodGet,
			// first we get tasks by ids, then in goroutine trying to get TIs but get error, also in templateCache goroutine
			// see that task type is script an then get templateDetail's then get execution results
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, validUUID).
					Return(validTasksSubRecent, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByTaskID(gomock.Any(), validUUID).
					Return(validTaskInstancesSubRecent, nil)
				return ti
			},
			templateCacheMock: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, validTasksSubRecent[0].OriginID, true).
					Return(models.TemplateDetails{}, nil)
				return tc
			},
			resultMock: func(er *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				er.EXPECT().
					GetByTaskInstanceIDs(gomock.Any()).
					Return([]models.ExecutionResult{{ManagedEndpointID: validUUID, TaskInstanceID: validUUID}}, nil)
				return er
			},
			expectedCode:               http.StatusOK,
			isNeedToCheckBody:          true,
			expectedNearestNextRunTime: validTasksSubRecent[1].RunTimeUTC,
			expectedTaskID:             validTasksSubRecent[0].ID,
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, nil,
				tc.taskSummaryMock,
				tc.taskInstanceMock,
				tc.taskPersistenceMock,
				tc.templateCacheMock,
				tc.resultMock,
				nil,
				nil,
				userMock,
				nil,
				nil,
				nil, nil, nil)

			router := getSubRecentRouter(service)

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

			if tc.isNeedToCheckBody {
				var receivedBody models.TaskSummaryDetails

				err := json.Unmarshal(w.Body.Bytes(), &receivedBody)
				if err != nil {
					t.Fatalf("Cant unmarshall received body")
				}

				if receivedBody.TaskSummary.NearestNextRunTime.Format(timeFormat) != tc.expectedNearestNextRunTime.Format(timeFormat) {
					t.Fatalf("Want %v but got %v", tc.expectedNearestNextRunTime, receivedBody.TaskSummary.NearestNextRunTime)
				}

				if receivedBody.TaskSummary.TaskID != tc.expectedTaskID {
					t.Fatalf("Want %v but got %v", tc.expectedTaskID, receivedBody.TaskSummary.TaskID)
				}
			}
		})
	}
}

func getSubRecentRouter(service TaskService) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/{partnerID}/tasks/{taskID}/data", http.HandlerFunc(service.SubRecent)).Methods(http.MethodGet)
	return
}
