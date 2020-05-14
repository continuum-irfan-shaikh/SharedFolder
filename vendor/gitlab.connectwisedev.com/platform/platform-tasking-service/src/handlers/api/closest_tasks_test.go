package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	mockusecases "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mock-usecases"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

const (
	defaultMsg = `failed on unexpected value of result "%v"`
)

func TestClosestTasks_ServeHTTP(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrl *gomock.Controller

	type expected struct {
		code     int
		bodyFile string
	}

	type payload struct {
		uc     func() TasksInteractor
		method string
		url    string
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "get tasks error",
			expected: expected{
				code:     http.StatusInternalServerError,
				bodyFile: "./testdata/get_tasks_error.json",
			},
			payload: payload{
				method: http.MethodPost,
				url:    "/tasking/v1/partners/50016646/tasks/data/closest",
				uc: func() TasksInteractor {
					taskInteractorMock := mockusecases.NewMockTasksInteractor(mockCtrl)
					taskInteractorMock.EXPECT().GetClosestTasks(gomock.Any(), gomock.Any()).
						Return(entities.EndpointsClosestTasks{}, errors.New("some_err")).Times(1)
					return taskInteractorMock
				},
			},
		}, {
			name: "success",
			expected: expected{
				code:     http.StatusOK,
				bodyFile: "./testdata/success.json",
			},
			payload: payload{
				method: http.MethodPost,
				url:    "/tasking/v1/partners/50016646/tasks/data/closest",
				uc: func() TasksInteractor {
					taskInteractorMock := mockusecases.NewMockTasksInteractor(mockCtrl)
					taskInteractorMock.EXPECT().GetClosestTasks(gomock.Any(), gomock.Any()).
						Return(entities.EndpointsClosestTasks{"taskID": entities.ClosestTasks{
							Previous: &entities.ClosestTask{
								ID:      "firstID",
								Name:    "firstName",
								RunDate: 56897412,
								Status:  "success",
							},
							Next: &entities.ClosestTask{
								ID:      "secondID",
								Name:    "secondName",
								RunDate: 66897412,
							},
						}}, nil).Times(1)
					return taskInteractorMock
				},
			},
		},
	}

	logger.Load(config.Config.Log)
	for _, test := range tc {
		mockCtrl = gomock.NewController(t)
		loadErr := translation.Load()
		Ω(loadErr).To(BeNil(), fmt.Sprintf(defaultMsg, loadErr))
		req, err := http.NewRequest(test.payload.method, test.payload.url, strings.NewReader("[]"))
		Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, err))
		req.Header.Add("Accept-Language", "en-US")

		mockErr := translation.MockTranslations()
		Ω(mockErr).To(BeNil(), fmt.Sprintf(defaultMsg, mockErr))

		rw := httptest.NewRecorder()
		handler := NewClosestTasks(test.payload.uc(), logger.Log)

		handler.ServeHTTP(rw, req)
		mockCtrl.Finish()

		body, err := readTestData(test.expected.bodyFile)
		Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, err))

		Ω(rw.Body.Bytes()).To(MatchJSON(body), fmt.Sprintf(defaultMsg, test.name))
		Ω(rw.Code).To(Equal(test.expected.code), fmt.Sprintf(defaultMsg, rw.Code))
	}
}

func readTestData(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
