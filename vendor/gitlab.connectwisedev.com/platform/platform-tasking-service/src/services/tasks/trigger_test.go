package tasks

import (
	"fmt"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/trigger"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const defaultMsg = `failed on unexpected value of result "%v"`

var (
	now           = time.Now()
	notValidStart = now.AddDate(0, 0, -5).Add(2 * time.Hour)
	notValidEnd   = now.AddDate(0, 0, -5).Add(4 * time.Hour)
)

func TestTaskService_ExecuteTrigger(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrl *gomock.Controller

	type payload struct {
		req             func() *http.Request
		triggerUC       func() trigger.Usecase
		tiPersistence   func() models.TaskInstancePersistence
		taskPersistence func() models.TaskPersistence
		templateCache   func() models.TemplateCache
	}

	type expected struct {
		body       string
		statusCode int
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name: "error while extracting struct from request",
			expected: expected{
				statusCode: http.StatusBadRequest,
				body:       `{"message":"Can not decode input data","errorCode":"error_cant_decode_input_data"}`,
			},
			payload: payload{
				req: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "", strings.NewReader(``))
					Ω(err).To(BeNil(), "error with NewRequest")
					return r
				},
				triggerUC: func() trigger.Usecase {
					triggerUCMock := mocks.NewMockUsecase(mockCtrl)
					return triggerUCMock
				},
				tiPersistence: func() models.TaskInstancePersistence {
					tiPersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrl)
					return tiPersistenceMock
				},
				taskPersistence: func() models.TaskPersistence {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrl)
					return taskPersistenceMock
				},
				templateCache: func() models.TemplateCache {
					templateCacheMock := mocks.NewMockTemplateCache(mockCtrl)
					return templateCacheMock
				},
			},
		},
		{
			name: "error while parsing UUID",
			expected: expected{
				statusCode: http.StatusBadRequest,
				body:       `{"message":"Can not decode input data","errorCode":"error_cant_decode_input_data"}`,
			},
			payload: payload{
				req: func() *http.Request {
					r, err := http.NewRequest(
						http.MethodGet,
						"",
						strings.NewReader("{\"dynamicGroupId\": \"groupID\"}"),
					)
					Ω(err).To(BeNil(), "error with NewRequest")
					return r
				},
				triggerUC: func() trigger.Usecase {
					triggerUCMock := mocks.NewMockUsecase(mockCtrl)
					return triggerUCMock
				},
				tiPersistence: func() models.TaskInstancePersistence {
					tiPersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrl)
					return tiPersistenceMock
				},
				taskPersistence: func() models.TaskPersistence {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrl)
					return taskPersistenceMock
				},
				templateCache: func() models.TemplateCache {
					templateCacheMock := mocks.NewMockTemplateCache(mockCtrl)
					return templateCacheMock
				},
			},
		},
		{
			name: "error while get active triggers",
			expected: expected{
				statusCode: http.StatusInternalServerError,
				body:       `{"message":"Can't get active triggers","errorCode":"error_cant_get_active_triggers"}`,
			},
			payload: payload{
				req: func() *http.Request {
					r, err := http.NewRequest(
						http.MethodGet,
						"",
						strings.NewReader("{\"dynamicGroupId\": \"groupID\"}"),
					)
					Ω(err).To(BeNil(), "error with NewRequest")
					r = mux.SetURLVars(r,
						map[string]string{
							"partnerID":   "partner_id",
							"endpointID":  "58a1af2f-6579-4aec-b45d-5dfde879ef01",
							"triggerType": "mockedTrigger",
						},
					)
					return r
				},
				triggerUC: func() trigger.Usecase {
					triggerUCMock := mocks.NewMockUsecase(mockCtrl)
					triggerUCMock.EXPECT().GetActiveTriggers(gomock.Any()).
						Return(nil, errors.New("some_err")).Times(1)
					return triggerUCMock
				},
				tiPersistence: func() models.TaskInstancePersistence {
					tiPersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrl)
					return tiPersistenceMock
				},
				taskPersistence: func() models.TaskPersistence {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrl)
					return taskPersistenceMock
				},
				templateCache: func() models.TemplateCache {
					templateCacheMock := mocks.NewMockTemplateCache(mockCtrl)
					return templateCacheMock
				},
			},
		},
		{
			name: "success with triggers with not valid time frame",
			expected: expected{
				statusCode: http.StatusCreated,
				body:       `{"message":"Updated"}`,
			},
			payload: payload{
				req: func() *http.Request {
					r, err := http.NewRequest(
						http.MethodGet,
						"",
						strings.NewReader("{\"dynamicGroupId\": \"groupID\"}"),
					)
					Ω(err).To(BeNil(), "error with NewRequest")
					r = mux.SetURLVars(r,
						map[string]string{
							"partnerID":   "partner_id",
							"endpointID":  "58a1af2f-6579-4aec-b45d-5dfde879ef01",
							"triggerType": "mockedTrigger",
						},
					)
					return r
				},
				triggerUC: func() trigger.Usecase {
					triggerUCMock := mocks.NewMockUsecase(mockCtrl)
					triggerUCMock.EXPECT().GetActiveTriggers(gomock.Any()).
						Return([]entities.ActiveTrigger{
							{
								Type:           models.LogoffTriggerType,
								PartnerID:      "partner_id",
								StartTimeFrame: notValidStart,
								EndTimeFrame:   notValidEnd,
							},
						}, nil).Times(1)
					return triggerUCMock
				},
				tiPersistence: func() models.TaskInstancePersistence {
					tiPersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrl)
					return tiPersistenceMock
				},
				taskPersistence: func() models.TaskPersistence {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrl)
					return taskPersistenceMock
				},
				templateCache: func() models.TemplateCache {
					templateCacheMock := mocks.NewMockTemplateCache(mockCtrl)
					return templateCacheMock
				},
			},
		},
	}

	for _, test := range tc {
		mockCtrl = gomock.NewController(t)
		ts := TaskService{
			trigger:                 test.triggerUC(),
			taskInstancePersistence: test.payload.tiPersistence(),
			taskPersistence:         test.payload.taskPersistence(),
			templateCache:           test.payload.templateCache(),
		}

		w := httptest.NewRecorder()
		ts.ExecuteTrigger(w, test.payload.req())
		Ω(w.Result().StatusCode).To(Equal(test.expected.statusCode), fmt.Sprintf(defaultMsg, test.name))
		Ω(w.Body.String()).To(Equal(test.expected.body), fmt.Sprintf(defaultMsg, test.name))

		mockCtrl.Finish()
	}
}
