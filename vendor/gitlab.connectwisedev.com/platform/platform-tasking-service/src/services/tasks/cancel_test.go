package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
)

func TestCancelNearestExecution(t *testing.T) {
	var (
		timeNow                = time.Now().UTC().Truncate(time.Minute)
		targetsStrSlice        = []string{validUUIDstr}
		targetsManagedEndpoint = models.Target{IDs: targetsStrSlice, Type: models.ManagedEndpoint}
		scheduledStatusesMap   = make(map[gocql.UUID]statuses.TaskInstanceStatus)
		pendingStatusesMap     = make(map[gocql.UUID]statuses.TaskInstanceStatus)
		canceledStatusesMap    = make(map[gocql.UUID]statuses.TaskInstanceStatus)

		scheduleRecurrent = apiModels.Schedule{
			Location:   "UTC",
			Regularity: apiModels.Recurrent,
			Repeat: apiModels.Repeat{
				Frequency: apiModels.Hourly,
				Every:     1,
			},
			EndRunTime: timeNow.Add(threeHoursDuration),
		}
		//------------------------------------------------test cases

		recurrentTask = []models.Task{{
			Name:              "task1",
			RunTimeUTC:        timeNow,
			PartnerID:         partnerID,
			ID:                validUUID,
			Targets:           targetsManagedEndpoint,
			Schedule:          scheduleRecurrent,
			ManagedEndpointID: validUUID},
		}
		expectedRecurrent = []models.Task{{
			Name:              "task1",
			RunTimeUTC:        timeNow,
			PartnerID:         partnerID,
			ID:                validUUID,
			State:             statuses.TaskStateActive,
			Targets:           targetsManagedEndpoint,
			Schedule:          scheduleRecurrent,
			ManagedEndpointID: validUUID},
		}

		oneTimeTask = []models.Task{{
			Name:               "task1",
			RunTimeUTC:         timeNow.Truncate(time.Minute),
			PartnerID:          partnerID,
			ID:                 validUUID,
			ExternalTask:       true,
			IsRequireNOCAccess: false,
			Targets:            targetsManagedEndpoint,
			Schedule:           apiModels.Schedule{Regularity: apiModels.OneTime},
			ManagedEndpointID:  validUUID},
		}
		//--------------------------------------- test cases
		scheduledTaskInstances = []models.TaskInstance{{
			ID:            validTaskInstanceUUID,
			Statuses:      scheduledStatusesMap,
			OverallStatus: statuses.TaskInstanceScheduled},
		}
		scheduledTaskInstancesWithPending = []models.TaskInstance{{
			ID:            validTaskInstanceUUID,
			Statuses:      scheduledStatusesMap,
			OverallStatus: statuses.TaskInstanceScheduled},
			{
				Statuses: pendingStatusesMap,
			},
		}
		pendingTaskInstances = []models.TaskInstance{{
			ID:            validTaskInstanceUUID,
			Statuses:      pendingStatusesMap,
			OverallStatus: statuses.TaskInstanceScheduled},
		}
		canceledTaskInstances = []models.TaskInstance{{
			ID:            validTaskInstanceUUID,
			Statuses:      canceledStatusesMap,
			OverallStatus: statuses.TaskInstanceScheduled},
		}
	)

	testCases := []struct {
		name                 string
		URL                  string
		taskInstanceMock     mock.TaskInstancePersistenceConf
		taskPersistenceMock  mock.TaskPersistenceConf
		userMock             user.Service
		method               string
		expectedErrorMessage string
		expectedCode         int
		expectedTask         interface{}
	}{
		{
			name:                 "testCase 1 - empty taskID",
			URL:                  fmt.Sprintf("/%v/tasks/%v/cancel", partnerID, emptyUUIDstr),
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorTaskIDHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name:                 "testCase 2 - bad taskID",
			URL:                  fmt.Sprintf("/%v/tasks/%v/cancel", partnerID, "badTaskID"),
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorTaskIDHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name:                 "testCase 3 - GetByIDs db error",
			URL:                  fmt.Sprintf("/%v/tasks/%v/cancel", partnerID, taskIDstr),
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
			expectedCode:         http.StatusInternalServerError,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, taskID).
					Return([]models.Task{}, errors.New("db error"))
				return tp
			},
		},
		{
			name:                 "testCase 4 - GetByIDs returned zero tasks",
			URL:                  fmt.Sprintf("/%v/tasks/%v/cancel", partnerID, taskIDstr),
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
			expectedCode:         http.StatusNotFound,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, taskID).
					Return([]models.Task{}, nil)
				return tp
			},
		},
		{
			name:   "testCase 5 - NOC task for not NOC, can't postpone",
			URL:    fmt.Sprintf("/%v/tasks/%v/cancel", partnerID, taskIDstr),
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, taskID).
					Return([]models.Task{{IsRequireNOCAccess: true}}, nil)
				return tp
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
		},
		{
			name:                 "testCase 6 - OneTimeTask, TaskInstance get error",
			URL:                  fmt.Sprintf("/%v/tasks/%v/cancel", partnerID, taskIDstr),
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantGetTaskInstances,
			expectedCode:         http.StatusInternalServerError,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, taskID).
					Return(oneTimeTask, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return([]models.TaskInstance{}, errors.New("err"))
				return ti
			},
		},
		{
			name:                 "testCase 7 - OneTimeTask, TaskInstance GetTopInstancesByTaskID returned zero instances",
			URL:                  fmt.Sprintf("/%v/tasks/%v/cancel", partnerID, taskIDstr),
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantGetTaskInstances,
			expectedCode:         http.StatusNotFound,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, taskID).
					Return(oneTimeTask, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return([]models.TaskInstance{}, nil)
				return ti
			},
		},
		{
			name:                 "testCase 8 - OneTimeTask, TaskInstance GetTopInstancesByTaskID returned instance with pending device",
			URL:                  fmt.Sprintf("/%v/tasks/%v/cancel", partnerID, taskIDstr),
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
			expectedCode:         http.StatusBadRequest,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, taskID).
					Return(oneTimeTask, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(pendingTaskInstances, nil)
				return ti
			},
		},
		{
			name:                 "testCase 8.5 - OneTimeTask, TaskInstance GetTopInstancesByTaskID returned instance with pending device",
			URL:                  fmt.Sprintf("/%v/tasks/%v/cancel", partnerID, taskIDstr),
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
			expectedCode:         http.StatusBadRequest,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, taskID).
					Return(oneTimeTask, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(canceledTaskInstances, nil)
				return ti
			},
		},
		{
			name:                 "testCase 9 - OneTimeTask, TaskInstance Insert err",
			URL:                  fmt.Sprintf("/%v/tasks/%v/cancel", partnerID, taskIDstr),
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantUpdateTaskInstances,
			expectedCode:         http.StatusInternalServerError,
			// 1. TP - GetByIDs, 2. TI - GetByIDs, 3. TI - Insert
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, taskID).
					Return(oneTimeTask, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(scheduledTaskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return ti
			},
		},
		{
			name:                 "testCase 10 - OneTimeTask TaskInsert Err",
			URL:                  fmt.Sprintf("/%v/tasks/%v/cancel", partnerID, taskIDstr),
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantSaveTaskToDB,
			expectedCode:         http.StatusInternalServerError,
			// 1. TP - GetByIDs, 2. TI - GetByIDs, 3. TI - Insert 4. TP - insert
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, taskID).
					Return(oneTimeTask, nil)
				tp.EXPECT().
					UpdateModifiedFieldsByMEs(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(scheduledTaskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
		},
		{
			name:         "testCase 11 - Recurrent canceled ok",
			URL:          fmt.Sprintf("/%v/tasks/%v/cancel", partnerID, taskIDstr),
			method:       http.MethodPut,
			expectedCode: http.StatusOK,
			// 1. TP - GetByIDs, 2. TI - GetByIDs, 3. TI - Insert 4. TP - insert
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, taskID).
					Return(recurrentTask, nil)
				tp.EXPECT().
					UpdateModifiedFieldsByMEs(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(scheduledTaskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedTask: expectedRecurrent,
		},
		{
			name:                 "testCase 12 - Recurrent with pending instance",
			URL:                  fmt.Sprintf("/%v/tasks/%v/cancel", partnerID, taskIDstr),
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
			expectedCode:         http.StatusBadRequest,
			// 1. TP - GetByIDs, 2. TI - GetByIDs, 3. TI - Insert 4. TP - insert
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, taskID).
					Return(recurrentTask, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(scheduledTaskInstancesWithPending, nil)
				return ti
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		scheduledStatusesMap[validUUID] = statuses.TaskInstanceScheduled
		pendingStatusesMap[validUUID] = statuses.TaskInstancePending
		canceledStatusesMap[validUUID] = statuses.TaskInstanceCanceled

		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, nil, nil, tc.taskInstanceMock, tc.taskPersistenceMock,
				nil, nil, nil, nil, tc.userMock, nil, nil, nil, nil, nil)

			router := getTaskServiceRouter(service)
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

			if tc.expectedTask != nil {
				var gotTasks models.Task
				expectedTasks := tc.expectedTask.([]models.Task)

				err := json.Unmarshal(w.Body.Bytes(), &gotTasks)
				if err != nil {
					t.Errorf("Error while unmarshalling body, %v", err)
				}

				if gotTasks.RunTimeUTC != expectedTasks[0].RunTimeUTC {
					t.Errorf("Want RunTimeUTC %v but got %v", expectedTasks[0].RunTimeUTC, gotTasks.RunTimeUTC)
				}
			}
		})
	}
}
