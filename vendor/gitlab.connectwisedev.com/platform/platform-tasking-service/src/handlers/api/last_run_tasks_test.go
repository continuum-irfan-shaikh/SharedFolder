package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	mockusecases "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mock-usecases"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"
)

const (
	tasksHistoryURI = "/{partnerID}/tasks/data/history"
	partnerID       = "1"
)

func TestLastRunTasks_ServeHTTP_Get(t *testing.T) {
	if err := translation.Load(); err != nil {
		t.Fatal(err)
	}
	var (
		to                       = time.Now().UTC()
		from                     = to.Add(-time.Hour)
		goodTimeFrame            = fmt.Sprintf(`?from=%v&to=%v`, from.Format(time.RFC3339Nano), to.Format(time.RFC3339Nano))
		badTimeFrameToBeforeFrom = fmt.Sprintf(`?from=%v&to=%v`, to.Format(time.RFC3339Nano), from.Format(time.RFC3339Nano))
		badTimeFrameToIsZero     = fmt.Sprintf(`?from=%v&to=%v`, from.Format(time.RFC3339Nano), 0)
		badTimeFrameFromIsZero   = fmt.Sprintf(`?from=%v&to=%v`, 0, to.Format(time.RFC3339Nano))
	)
	logger.Load(config.Config.Log)
	testCases := []struct {
		name         string
		mock         func(ctl *gomock.Controller) (logger.Logger, TasksInteractor)
		expectedCode int
		URL          string
	}{
		{
			name: "testCase 1 - uc returned err",
			mock: func(ctl *gomock.Controller) (logger.Logger, TasksInteractor) {
				ti := mockusecases.NewMockTasksInteractor(ctl)
				ti.EXPECT().GetTasksHistory(gomock.Any(), from, to).Return([]entities.ScheduledTasks{}, errors.New("err"))
				return nil, ti
			},
			URL:          "/" + partnerID + "/tasks/data/history" + goodTimeFrame,
			expectedCode: http.StatusBadRequest,
		}, {
			name: "testCase 2 - uc returned bad request err (invalid time frame)",
			mock: func(ctl *gomock.Controller) (logger.Logger, TasksInteractor) {
				return nil, nil
			},
			URL:          "/" + partnerID + "/tasks/data/history" + badTimeFrameToBeforeFrom,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "testCase 3 - uc returned bad request err (bad to time frame)",
			mock: func(ctl *gomock.Controller) (logger.Logger, TasksInteractor) {
				return nil, nil
			},
			URL:          "/" + partnerID + "/tasks/data/history" + badTimeFrameToIsZero,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "testCase 4 - uc returned bad request err (bad from time frame)",
			mock: func(ctl *gomock.Controller) (logger.Logger, TasksInteractor) {
				return nil, nil
			},
			URL:          "/" + partnerID + "/tasks/data/history" + badTimeFrameFromIsZero,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "testCase 6 - uc returned internal err",
			mock: func(ctl *gomock.Controller) (logger.Logger, TasksInteractor) {
				ti := mockusecases.NewMockTasksInteractor(ctl)
				ti.EXPECT().GetTasksHistory(gomock.Any(), from, to).Return([]entities.ScheduledTasks{}, errorcode.NewInternalServerErr("", ""))

				return nil, ti
			},
			URL:          "/" + partnerID + "/tasks/data/history" + goodTimeFrame,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			_, ti := tc.mock(mc)
			s := *NewLastRunTasksApi(ti, logger.Log)
			router := getHistoryRouter(s)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tc.URL, nil)
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
		})
	}
}

func getHistoryRouter(service LastRunTasksApi) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc(tasksHistoryURI, service.ServeHTTP).Methods(http.MethodGet)
	return
}
