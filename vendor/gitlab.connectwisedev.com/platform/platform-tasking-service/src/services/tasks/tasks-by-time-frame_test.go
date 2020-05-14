package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	mockusecases "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mock-usecases"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-repository"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/model-mocks"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const (
	partnerID     = "1"
	validUUIDstr  = "58a1af2f-6579-4aec-b45d-5dfde879ef01"
	validUUID2str = "58a1af2f-6579-4aec-b45d-5dfde879ef21"
	validUUID3str = "58a1af2f-6579-4aec-b45d-5dfde879ef31"
	validUUID4str = "58a1af2f-6579-4aec-b45d-5dfde879ef34"
)

var (
	to                       = time.Now().UTC()
	from                     = to.Add(-time.Hour)
	goodTimeFrame            = fmt.Sprintf(`?from=%v&to=%v`, from.Format(time.RFC3339Nano), to.Format(time.RFC3339Nano))
	badTimeFrameToIsZero     = fmt.Sprintf(`?from=%v&to=%v`, from.Format(time.RFC3339Nano), 0)
	badTimeFrameToBeforeFrom = fmt.Sprintf(`?from=%v&to=%v`, to.Format(time.RFC3339Nano), from.Format(time.RFC3339Nano))
	validUUID, _             = gocql.ParseUUID(validUUIDstr)
	validUUID2, _            = gocql.ParseUUID(validUUID2str)
	validUUID3, _            = gocql.ParseUUID(validUUID3str)
	validUUID4, _            = gocql.ParseUUID(validUUID4str)
)

func TestLastTasks(t *testing.T) {
	targetsStrSlice := []string{validUUIDstr}
	targets := models.Target{IDs: targetsStrSlice}

	validTasks := []models.Task{
		{Name: "task1", ID: validUUID, ExternalTask: true, IsRequireNOCAccess: false, Targets: targets},
		{Name: "task2", ID: validUUID2, IsRequireNOCAccess: true, Targets: targets},
		{Name: "task3", ID: validUUID3, IsRequireNOCAccess: false, Targets: targets},
		{Name: "task4", ID: validUUID4, IsRequireNOCAccess: false, Targets: targets},
	}

	validTasksForScheduledFrame := []models.Task{
		{Name: "task1", ID: validUUID, Targets: targets},
		{Name: "task2", ID: validUUID2, Targets: targets, State: statuses.TaskStateInactive},
	}

	validTaskInstances := []models.TaskInstance{
		{TaskID: validUUID},
		{TaskID: validUUID2},
		{TaskID: validUUID3},
		{TaskID: validUUID4},
	}

	validTaskInstanceStatusesSuccessesLastTasks := make(map[gocql.UUID]statuses.TaskInstanceStatus)
	validTaskInstanceStatusesScheduledLastTasks := make(map[gocql.UUID]statuses.TaskInstanceStatus)
	validTaskInstanceStatusesSuccessAndDisabled := make(map[gocql.UUID]statuses.TaskInstanceStatus)
	validTaskInstanceStatusesSuccessesLastTasks[validUUID] = statuses.TaskInstanceSuccess
	validTaskInstanceStatusesScheduledLastTasks[validUUID4] = statuses.TaskInstanceScheduled

	validTaskInstanceStatusesSuccessAndDisabled[validUUID] = statuses.TaskInstanceSuccess
	validTaskInstanceStatusesSuccessAndDisabled[validUUID4] = statuses.TaskInstanceDisabled

	scheduledStatusStr, err := statuses.TaskInstanceStatusText(statuses.TaskInstanceScheduled)
	if err != nil {
		t.Errorf("cannot parse status")
	}

	testCases := []struct {
		name                 string
		URL                  string
		taskSummaryMock      mock.TaskSummaryPersistenceConf
		taskInstanceMock     mock.TaskInstancePersistenceConf
		taskPersistenceMock  mock.TaskPersistenceConf
		userMock             user.Service
		method               string
		expectedErrorMessage string
		expectedBody         interface{}
		expectedCode         int
	}{
		{
			name:                 "testCase 0 - Bad time, 'from'  & 'to' is zero",
			URL:                  "/" + partnerID + "/tasks/data/recent",
			method:               http.MethodGet,
			expectedErrorMessage: errorcode.ErrorTimeFrameHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name:                 "testCase 1 - Bad time, 'to' is zero",
			URL:                  "/" + partnerID + "/tasks/data/recent" + badTimeFrameToIsZero,
			method:               http.MethodGet,
			expectedErrorMessage: errorcode.ErrorTimeFrameHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name:                 "testCase 2 - Bad time, from and to are swapped",
			URL:                  "/" + partnerID + "/tasks/data/recent" + badTimeFrameToBeforeFrom,
			method:               http.MethodGet,
			expectedErrorMessage: errorcode.ErrorTimeFrameHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name: "testCase 3 - Server error, cannot get TaskInstances GetByStartedAtAfter err",
			URL:  "/" + partnerID + "/tasks/data/recent" + goodTimeFrame,
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByStartedAtAfter(gomock.Any(), partnerID, from, to).
					Return([]models.TaskInstance{}, errors.New("some error"))
				return ti
			},
			method:               http.MethodGet,
			expectedErrorMessage: errorcode.ErrorCantGetTaskInstances,
			expectedCode:         http.StatusInternalServerError,
		},
		{
			name: "testCase 4 - Server error, cannot get Tasks by TaskID taskPersistenceErr",
			URL:  "/" + partnerID + "/tasks/data/recent" + goodTimeFrame,
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByStartedAtAfter(gomock.Any(), partnerID, from, to).
					Return(validTaskInstances, nil)
				return ti
			},
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, gomock.Any()).
					Return(nil, errors.New("taskPersistenceErr"))
				return tp
			},
			method:               http.MethodGet,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
			expectedCode:         http.StatusInternalServerError,
		},
		{
			// same as previous but user doesnt have noc access, so we skip NOC task
			name: "testCase 6 - user doesn't have NOC access",
			URL:  "/" + partnerID + "/tasks/data/recent" + goodTimeFrame,
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByStartedAtAfter(gomock.Any(), partnerID, from, to).
					Return([]models.TaskInstance{
						{TaskID: validUUID},
						{TaskID: validUUID2, Statuses: validTaskInstanceStatusesSuccessesLastTasks},
						{TaskID: validUUID3},
						{TaskID: validUUID4, Statuses: validTaskInstanceStatusesScheduledLastTasks},
					}, nil)
				return ti
			},
			userMock: user.NewMock("", partnerID, "", "", false),
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), gomock.Any(), partnerID, true, gomock.Any()).
					Return(validTasks, nil)
				return tp
			},
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			// expecting empty body
		},
		{
			// skip task2, because its inactive
			name:     "testCase 7 - future time frame (scheduled filter activated), getByPartnerErr",
			URL:      "/" + partnerID + "/tasks/data/recent" + goodTimeFrame + "&isScheduled=true",
			userMock: user.NewMock("", partnerID, "", "", false),
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByPartnerAndTime(gomock.Any(), partnerID, gomock.Any()).
					Return([]models.Task{}, errors.New("err"))
				return tp
			},
			method:               http.MethodGet,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
		},
		{
			// skip task2, because its inactive
			name: "testCase 8 - future time frame (scheduled filter activated), get Task Instances Errf",
			URL:  "/" + partnerID + "/tasks/data/recent" + goodTimeFrame + "&isScheduled=true",
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), gomock.Any()).
					Return([]models.TaskInstance{}, errors.New("err"))
				return ti
			},
			userMock: user.NewMock("", partnerID, "", "", false),
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByPartnerAndTime(gomock.Any(), partnerID, gomock.Any()).
					Return(validTasksForScheduledFrame, nil)
				return tp
			},
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
		},
		{
			// skip task2, because its inactive
			name: "testCase 9 - future time frame (scheduled filter activated)",
			URL:  "/" + partnerID + "/tasks/data/recent" + goodTimeFrame + "&isScheduled=true",
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), gomock.Any()).
					Return([]models.TaskInstance{
						{TaskID: validUUID, Statuses: validTaskInstanceStatusesScheduledLastTasks},
						{TaskID: validUUID2, Statuses: validTaskInstanceStatusesScheduledLastTasks},
					}, nil)
				return ti
			},
			userMock: user.NewMock("", partnerID, "", "", false),
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByPartnerAndTime(gomock.Any(), partnerID, gomock.Any()).
					Return(validTasksForScheduledFrame, nil)
				return tp
			},
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: []models.TaskDetailsWithStatuses{{
				Task:           validTasksForScheduledFrame[0],
				TaskInstance:   models.TaskInstance{TaskID: validUUID, Statuses: validTaskInstanceStatusesScheduledLastTasks},
				Statuses:       map[string]int{scheduledStatusStr: 1},
				CanBePostponed: true,
				CanBeCanceled:  true,
			}},
		},
		{
			// skip task2, because its inactive
			name: "testCase 10 - future time frame, skipping disabled oneTime",
			URL:  "/" + partnerID + "/tasks/data/recent" + goodTimeFrame + "&isScheduled=true",
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), gomock.Any()).
					Return([]models.TaskInstance{
						{TaskID: validUUID, Statuses: validTaskInstanceStatusesSuccessAndDisabled},
					}, nil)
				return ti
			},
			userMock: user.NewMock("", partnerID, "", "", false),
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByPartnerAndTime(gomock.Any(), partnerID, gomock.Any()).
					Return(validTasksForScheduledFrame, nil)
				return tp
			},
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			// expecting nil
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, nil,
				tc.taskSummaryMock,
				tc.taskInstanceMock,
				tc.taskPersistenceMock,
				nil,
				nil,
				nil,
				nil,
				tc.userMock,
				nil, nil, nil, nil, nil)

			router := getFramesRouter(service)

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
			expectedBodyBytes, err := json.Marshal(tc.expectedBody)
			if err != nil {
				t.Errorf("Cannot parse expected body%v", tc.name)
			}
			if !bytes.Equal(expectedBodyBytes, w.Body.Bytes()) && tc.expectedBody != nil {
				t.Errorf("Want body %v but got  \n%v", string(expectedBodyBytes), w.Body.String())
			}
		})
	}
}

// getMockedTaskService is a function that returns configured mocked TaskDefinitionService
func getMockedTaskService(mockController *gomock.Controller,
	td mock.TaskDefinitionConf,
	ts mock.TaskSummaryPersistenceConf,
	ti mock.TaskInstancePersistenceConf,
	tp mock.TaskPersistenceConf,
	tc mock.TemplateCacheConf,
	er mock.ExecutionResultPersistenceConf,
	ee mock.ExecutionExpirationPersistenceConf,
	us mock.UserSitesConf,
	userMock user.Service,
	u mock.UserUC,
	triggerUC mock.TriggerUCConf,
	targetsRepoConf mock.TargetRepoConf,
	dg mock.DGRepoConf,
	sr mock.SitesRepoConf) TaskService {

	taskPersistenceMock := mock.NewMockTaskPersistence(mockController)
	if tp != nil {
		taskPersistenceMock = tp(taskPersistenceMock)
	}

	tdMock := mock.NewMockTaskDefinitionPersistence(mockController)
	if td != nil {
		tdMock = td(tdMock)
	}

	targetsRepoMock := mock.NewMockTargetsRepo(mockController)
	if targetsRepoConf != nil {
		targetsRepoMock = targetsRepoConf(targetsRepoMock)
	}

	taskSummaryPersistenceMock := mock.NewMockTaskSummaryPersistence(mockController)
	if ts != nil {
		taskSummaryPersistenceMock = ts(taskSummaryPersistenceMock)
	}

	taskInstanceMock := mock.NewMockTaskInstancePersistence(mockController)
	if ti != nil {
		taskInstanceMock = ti(taskInstanceMock)
	}

	templateCacheMock := mock.NewMockTemplateCache(mockController)
	if tc != nil {
		templateCacheMock = tc(templateCacheMock)
	}

	execResultsMock := mock.NewMockExecutionResultPersistence(mockController)
	if er != nil {
		execResultsMock = er(execResultsMock)
	}

	executionExpiration := mock.NewMockExecutionExpirationPersistence(mockController)
	if ee != nil {
		executionExpiration = ee(executionExpiration)
	}

	userSitesMock := mock.NewMockUserSitesPersistence(mockController)
	if us != nil {
		userSitesMock = us(userSitesMock)
	}

	userUCMock := mock.NewMockUserUC(mockController)
	if u != nil {
		userUCMock = u(userUCMock)
	}

	userServiceMock := userMock
	if userServiceMock == nil {
		userServiceMock = userServiceMockNoErr
	}

	trUCMock := mock.NewMockUsecase(mockController)
	if triggerUC != nil {
		trUCMock = triggerUC(trUCMock)
	}

	dgRepo := mock.NewMockDynamicGroupRepo(mockController)
	if dg != nil {
		dgRepo = dg(dgRepo)
	}

	siteRepo := mock.NewMockSiteRepo(mockController)
	if sr != nil {
		siteRepo = sr(siteRepo)
	}

	executionResultUpdateService := mockusecases.NewMockExecutionResultUpdateUC(mockController)
	executionResultUpdateService.EXPECT().ProcessExecutionResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	// setting up the others mock's that are used in this handler but we don't need to test them
	return New(
		tdMock,
		taskPersistenceMock,
		taskInstanceMock,
		templateCacheMock,
		execResultsMock,
		taskSummaryPersistenceMock,
		executionExpiration,
		userSitesMock,
		userServiceMock,
		DefaultMock{},
		siteRepo,
		dgRepo,
		modelMocks.TaskCounterDefaultMock{},
		nil,
		MessagingMock{},
		http.DefaultClient,
		userUCMock,
		trUCMock, //  add a mocked service
		nil,
		targetsRepoMock,
		executionResultUpdateService,
		&mockrepositories.EncryptionServiceMock{},
		&mockrepositories.AgentEncryptionServiceMock{},
	)
}

func getFramesRouter(service TaskService) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/{partnerID}/tasks/data/recent", http.HandlerFunc(service.LastTasks)).Methods(http.MethodGet)
	return
}

// DefaultMock ..
type DefaultMock struct{}

// GetLocationByEndpointID ..
func (DefaultMock) GetLocationByEndpointID(_ context.Context, partnerID string, endpointID gocql.UUID) (location *time.Location, err error) {
	return nil, nil
}

// GetSiteIDByEndpointID ..
func (DefaultMock) GetSiteIDByEndpointID(_ context.Context, partnerID string, endpointID gocql.UUID) (siteID, clientID string, err error) {
	return "", "", nil
}

// GetMachineNameByEndpointID ..
func (DefaultMock) GetMachineNameByEndpointID(_ context.Context, partnerID string, endpointID gocql.UUID) (string, error) {
	return "", nil
}

// GetResourceTypeByEndpointID  ..
func (DefaultMock) GetResourceTypeByEndpointID(_ context.Context, partnerID string, endpointID gocql.UUID) (integration.ResourceType, error) {
	return "", nil
}

// MessagingMock ..
type MessagingMock struct{}

// Push ..
func (MessagingMock) Push(message interface{}, msgType string) error { return nil }
