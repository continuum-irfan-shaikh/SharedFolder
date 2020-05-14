package tasks

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"github.com/golang/mock/gomock"
)

func TestName(t *testing.T) {
	validCount := 999

	testCases := []struct {
		name                 string
		method               string
		URL                  string
		expectedCode         int
		expectedErrorMessage string
		expectedBody         interface{}
		mockTaskInstance     mock.TaskInstancePersistenceConf
	}{
		{
			name:                 "testCase 0 - invalid taskID",
			URL:                  "/" + partnerID + "/tasks/badUUID/task-instances/count",
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTaskIDHasBadFormat,
		},
		{
			name: "testCase 1 - Cant get count from database",
			URL:  "/" + partnerID + "/tasks/" + validUUIDstr + "/task-instances/count",
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT(). //context + taskid
						GetInstancesCountByTaskID(gomock.Any(), validUUID).
						Return(0, errors.New("db error"))
				return ti
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskInstanceCountByTaskID,
		},
		{
			name: "testCase 2 - Successes",
			URL:  "/" + partnerID + "/tasks/" + validUUIDstr + "/task-instances/count",
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT(). //context + taskid
						GetInstancesCountByTaskID(gomock.Any(), validUUID).
						Return(validCount, nil)
				return ti
			},
			expectedCode: http.StatusOK,
			expectedBody: InstancesCount{
				InstancesCount: validCount,
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, nil, nil, tc.mockTaskInstance, nil, nil, nil,
				nil, nil, userServiceMockNoErr, nil, nil, nil, nil, nil)
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

			if tc.expectedBody != nil {
				var gotBody InstancesCount
				err := json.Unmarshal(w.Body.Bytes(), &gotBody)
				if err != nil {
					t.Errorf("Cant unmarshall received body %v", err)
				}
				expectedBody := tc.expectedBody.(InstancesCount)

				if gotBody.InstancesCount != expectedBody.InstancesCount {
					t.Errorf("Want % but got %v", expectedBody.InstancesCount, gotBody.InstancesCount)
				}
			}
		})
	}
}
