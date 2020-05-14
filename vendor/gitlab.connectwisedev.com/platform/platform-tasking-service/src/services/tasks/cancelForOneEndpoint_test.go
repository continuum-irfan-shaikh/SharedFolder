package tasks

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

func TestCancelNearestExecutionForOneME(t *testing.T) {
	var (
		timeNow                = time.Now().UTC().Truncate(time.Minute)
		targetsStrSlice        = []string{validUUIDstr}
		targetsManagedEndpoint = models.Target{IDs: targetsStrSlice, Type: models.ManagedEndpoint}
		scheduledStatusesMap   = make(map[gocql.UUID]statuses.TaskInstanceStatus)
		pendingStatusesMap     = make(map[gocql.UUID]statuses.TaskInstanceStatus)
		scheduledTIMap         = make(map[gocql.UUID]statuses.TaskInstanceStatus)
		scheduledTI            = []models.TaskInstance{{ID: validTaskInstanceUUID, Statuses: scheduledTIMap, OverallStatus: statuses.TaskInstanceScheduled}}
		scheduledTIWithPending = []models.TaskInstance{{ID: validTaskInstanceUUID, Statuses: scheduledTIMap, OverallStatus: statuses.TaskInstanceScheduled}, {Statuses: pendingStatusesMap}}

		oneTimeTask = []models.Task{
			{
				Name:               "oneTimeTask",
				RunTimeUTC:         timeNow.Truncate(time.Minute),
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           apiModels.Schedule{Regularity: apiModels.OneTime},
				ManagedEndpointID:  validUUID,
			},
		}

		scheduledTaskInstances = []models.TaskInstance{{
			ID:            validTaskInstanceUUID,
			Statuses:      scheduledStatusesMap,
			OverallStatus: statuses.TaskInstanceScheduled},
		}

		scheduleRecurrent = apiModels.Schedule{
			Location:   "UTC",
			Regularity: apiModels.Recurrent,
			Repeat: apiModels.Repeat{
				Frequency: apiModels.Hourly,
				Every:     1,
			},
			EndRunTime: timeNow.Add(threeHoursDuration),
		}

		recurrentTask = []models.Task{
			{
				Name:              "recurrentTask",
				RunTimeUTC:        timeNow,
				PartnerID:         partnerID,
				ID:                validUUID,
				Targets:           targetsManagedEndpoint,
				Schedule:          scheduleRecurrent,
				ManagedEndpointID: validUUID,
			},
		}
		expectedRecurrent = []models.Task{
			{
				Name:              "expectedRecurrentTask",
				RunTimeUTC:        timeNow,
				PartnerID:         partnerID,
				ID:                validUUID,
				State:             statuses.TaskStateActive,
				Targets:           targetsManagedEndpoint,
				Schedule:          scheduleRecurrent,
				ManagedEndpointID: validUUID,
			},
		}
	)

	scheduledTIMap[managedEndpointUUID] = statuses.TaskInstanceScheduled
	scheduledStatusesMap[managedEndpointUUID] = statuses.TaskInstanceScheduled
	pendingStatusesMap[validUUID] = statuses.TaskInstancePending

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
			URL:                  "/" + partnerID + "/tasks/" + emptyUUIDstr + "/managed-endpoint/" + managedEndpointID + "/cancel",
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorTaskIDHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name:                 "testCase 2 - bad taskID",
			URL:                  "/" + partnerID + "/tasks/badTaskID/managed-endpoint/" + managedEndpointID + "/cancel",
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorTaskIDHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name:                 "testCase 3 - bad endpointID",
			URL:                  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoint/badID/cancel",
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorEndpointIDHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name:                 "testCase 4 - empty endpointID",
			URL:                  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoint/" + emptyUUIDstr + "/cancel",
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorEndpointIDHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name:                 "testCase 5 - GetByIDs get DB error",
			URL:                  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoint/" + managedEndpointID + "/cancel",
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
			expectedCode:         http.StatusInternalServerError,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, taskID).
					Return([]models.Task{}, errors.New("error from DB"))
				return tp
			},
		},
		{
			name:                 "testCase 6 - GetByIDs returned zero tasks",
			URL:                  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoint/" + managedEndpointID + "/cancel",
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
			name:   "testCase 7 - notNOC partner wants to cancel NOC task",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoint/" + managedEndpointID + "/cancel",
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
			name:                 "testCase 8 - task has pending instance -> can't cancel",
			URL:                  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoint/" + managedEndpointID + "/cancel",
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
			expectedCode:         http.StatusBadRequest,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, taskID).
					Return(recurrentTask, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(scheduledTIWithPending, nil)
				return ti
			},
		},
		{
			name:                 "testCase 9 - OneTimeTask, GetByID for taskInstance error",
			URL:                  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoint/" + managedEndpointID + "/cancel",
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
					Return([]models.TaskInstance{}, errors.New("err from DB"))
				return ti
			},
		},
		{
			name:                 "testCase 10 - OneTimeTask, TaskInstance GetByIDs returned zero instances",
			URL:                  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoint/" + managedEndpointID + "/cancel",
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
					Return([]models.TaskInstance{}, nil)
				return ti
			},
		},
		{
			name:                 "testCase 12 - OneTimeTask, TaskInstance Insert err",
			URL:                  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoint/" + managedEndpointID + "/cancel",
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
			name:         "testCase 13 - Recurrent canceled ok",
			URL:          "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoint/" + managedEndpointID + "/cancel",
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
					Return(scheduledTI, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedTask: expectedRecurrent,
		},
		{
			name:         "testCase 14 - Recurrent canceled error while updating fields",
			URL:          "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoint/" + managedEndpointID + "/cancel",
			method:       http.MethodPut,
			expectedCode: http.StatusBadRequest,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, taskID).
					Return(recurrentTask, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(scheduledTI, nil)
				return ti
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {

			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, nil, nil, tc.taskInstanceMock, tc.taskPersistenceMock, nil,
				nil, nil, nil, tc.userMock, nil, nil, nil, nil, nil)

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
				var gotTask models.Task

				err := json.Unmarshal(w.Body.Bytes(), &gotTask)
				if err != nil {
					t.Errorf("Error while unmarshalling body, %v", err)
				}

				if gotTask.RunTimeUTC != expectedRecurrent[0].RunTimeUTC {
					t.Errorf("Want RunTimeUTC %v but got %v", expectedRecurrent[0].RunTimeUTC, gotTask.RunTimeUTC)
				}

			}
		})
	}
}
