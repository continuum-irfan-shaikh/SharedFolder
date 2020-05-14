package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	mockusecases "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mock-usecases"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"

	. "github.com/onsi/gomega"
)

const scheduledURI = "/{partnerID}/tasks/data/scheduled"

func TestScheduledTasks_ServeHTTP_Get(t *testing.T) {
	if err := translation.Load(); err != nil {
		t.Fatal(err)
	}
	logger.Load(config.Config.Log)
	testCases := []struct {
		name         string
		mock         func(ctl *gomock.Controller) (logger.Logger, TasksInteractor)
		expectedCode int
	}{
		{
			name: "testCase 1 - uc returned err",
			mock: func(ctl *gomock.Controller) (logger.Logger, TasksInteractor) {
				ti := mockusecases.NewMockTasksInteractor(ctl)
				ti.EXPECT().GetScheduledTasks(gomock.Any()).Return([]entities.ScheduledTasks{}, errors.New("err"))
				return nil, ti
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "testCase 2 - uc returned bad request err",
			mock: func(ctl *gomock.Controller) (logger.Logger, TasksInteractor) {

				ti := mockusecases.NewMockTasksInteractor(ctl)
				ti.EXPECT().GetScheduledTasks(gomock.Any()).Return([]entities.ScheduledTasks{}, errorcode.NewBadRequestErr("", ""))
				return nil, ti
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "testCase 3 - uc returned internal err",
			mock: func(ctl *gomock.Controller) (logger.Logger, TasksInteractor) {

				ti := mockusecases.NewMockTasksInteractor(ctl)
				ti.EXPECT().GetScheduledTasks(gomock.Any()).Return([]entities.ScheduledTasks{}, errorcode.NewInternalServerErr("", ""))
				return nil, ti
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "testCase 4 - ok",
			mock: func(ctl *gomock.Controller) (logger.Logger, TasksInteractor) {
				ti := mockusecases.NewMockTasksInteractor(ctl)
				ti.EXPECT().GetScheduledTasks(gomock.Any()).Return([]entities.ScheduledTasks{}, nil)
				return nil, ti
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			_, ti := tc.mock(mc)
			s := *NewScheduledTasksApi(ti, logger.Log)
			router := getScheduledRouter(s)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, scheduledURI, nil)
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
		})
	}
}

func getScheduledRouter(service ScheduledTasks) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc(scheduledURI, service.ServeHTTP).Methods(http.MethodGet)
	return
}

func TestClosestTasks_ServeHTTP_Delete(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrl *gomock.Controller

	type expected struct {
		code int
	}

	type payload struct {
		uc     func() TasksInteractor
		method string
		body   string
		url    string
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "success",
			expected: expected{
				code: http.StatusNoContent,
			},
			payload: payload{
				method: http.MethodDelete,
				url:    "/tasking/v1/partners/50016646/tasks/data/scheduled",
				body:   `{"ids":["111"]}`,
				uc: func() TasksInteractor {
					taskInteractorMock := mockusecases.NewMockTasksInteractor(mockCtrl)

					taskInteractorMock.EXPECT().DeleteScheduledTasks(gomock.Any(), gomock.Any()).
						Return(nil).
						Times(1)

					return taskInteractorMock
				},
			},
		},
		{
			name: "failed to parse request",
			expected: expected{
				code: http.StatusBadRequest,
			},
			payload: payload{
				method: http.MethodDelete,
				url:    "/tasking/v1/partners/50016646/tasks/data/scheduled",
				body:   `111`,
				uc: func() TasksInteractor {
					taskInteractorMock := mockusecases.NewMockTasksInteractor(mockCtrl)
					return taskInteractorMock
				},
			},
		},
		{
			name: "failed to delete tasks",
			expected: expected{
				code: http.StatusInternalServerError,
			},
			payload: payload{
				method: http.MethodDelete,
				url:    "/tasking/v1/partners/50016646/tasks/data/scheduled",
				body:   `{"ids":["111"]}`,
				uc: func() TasksInteractor {
					taskInteractorMock := mockusecases.NewMockTasksInteractor(mockCtrl)

					taskInteractorMock.EXPECT().DeleteScheduledTasks(gomock.Any(), gomock.Any()).
						Return(errors.New("fail")).
						Times(1)

					return taskInteractorMock
				},
			},
		},
		{
			name: "failed because of invalid method",
			expected: expected{
				code: http.StatusMethodNotAllowed,
			},
			payload: payload{
				method: http.MethodPut,
				url:    "/tasking/v1/partners/50016646/tasks/data/scheduled",
				body:   `{"ids":["111"]}`,
				uc: func() TasksInteractor {
					taskInteractorMock := mockusecases.NewMockTasksInteractor(mockCtrl)
					return taskInteractorMock
				},
			},
		},
	}

	logger.Load(config.Config.Log)
	for _, test := range tc {
		mockCtrl = gomock.NewController(t)
		req, err := http.NewRequest(test.payload.method, test.payload.url, bytes.NewBufferString(test.payload.body))
		Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, err))
		req.Header.Add("Accept-Language", "en-US")

		mockErr := translation.MockTranslations()
		Ω(mockErr).To(BeNil(), fmt.Sprintf(defaultMsg, mockErr))

		rw := httptest.NewRecorder()

		handler := NewScheduledTasksApi(test.payload.uc(), logger.Log)

		handler.ServeHTTP(rw, req)
		mockCtrl.Finish()

		Ω(rw.Code).To(Equal(test.expected.code), fmt.Sprintf(defaultMsg, rw.Code))
	}
}
