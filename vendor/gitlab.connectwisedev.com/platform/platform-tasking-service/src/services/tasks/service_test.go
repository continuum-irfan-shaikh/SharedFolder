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

	"github.com/coocood/freecache"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	accessControl "gitlab.connectwisedev.com/platform/platform-tasking-service/src/access-control"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/asset"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	mockusecases "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mock-usecases"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	modelMocks "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/model-mocks"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"
	"gopkg.in/jarcoal/httpmock.v1"
)

// siteData describes data which is received from SitesMS
type siteData struct {
	SiteList []site `json:"siteDetailList"`
}

// Site describes site data
type site struct {
	ID int64 `json:"siteId"`
}

var service TaskService

const (
	validTaskInstanceUUIDstr = "58a1af2f-6579-4aec-b45d-000000000001"
	taskIDstr                = "58a1af2f-6579-4aec-b45d-5dfde879ef01"
	managedEndpointID        = "a0e831e9-145d-48ee-b9c0-a44c63062e37"
	emptyUUIDstr             = "00000000-0000-0000-0000-000000000000"
	uid                      = "admin"
	sequenceType             = "sequence"
	uidHeader                = "admin"
	dgDelimiter              = "%20OR%20"
	siteDelimiter            = "%2C"
)

var (
	testUUID                 = gocql.TimeUUID()
	validTaskInstanceUUID, _ = gocql.ParseUUID(validTaskInstanceUUIDstr)
	taskID, _                = gocql.ParseUUID(taskIDstr)
	managedEndpointUUID, _   = gocql.ParseUUID(managedEndpointID)

	targetDG = models.TargetsByType{
		models.DynamicGroup: []string{testUUID.String()},
	}

	targetSitesFixedUUID = models.TargetsByType{
		models.Site: []string{"1"},
	}

	targetManagedEndpoint = models.TargetsByType{
		models.ManagedEndpoint: []string{testUUID.String()},
	}

	userServiceMockNoErr = user.NewMock("", partnerID, uidHeader, "ValidToken", false)
	httpClient           = http.DefaultClient
)

func init() {
	config.Load()
	logger.Load(config.Config.Log)
	asset.Load()
	translation.Load()
	accessControl.Load(httpClient)
	models.TaskInstancePersistenceInstance = modelMocks.TaskInstanceCustomizableMock{
		InsertF: func(context.Context, models.TaskInstance) error { return nil },
	}
	models.TemplatesCache = freecache.NewCache(config.Config.TDTCacheSettings.Size)
}

func TestEdit(t *testing.T) {
	var (
		validRunNowTask = models.Task{
			TargetsByType: targetManagedEndpoint,
			OriginID:      taskID,
			Type:          config.ScriptTaskType,
			Parameters:    `{"num":1}`,
			Schedule:      apiModels.Schedule{Regularity: apiModels.RunNow},
		}
		validRunNowTaskWithObject = models.Task{
			TargetsByType:    targetManagedEndpoint,
			OriginID:         taskID,
			Type:             config.ScriptTaskType,
			ParametersObject: map[string]interface{}{"param": false},
			Schedule:         apiModels.Schedule{Regularity: apiModels.RunNow},
		}
		validRunNowTaskEmptyParams = models.Task{
			TargetsByType: targetManagedEndpoint,
			OriginID:      taskID,
			Type:          config.ScriptTaskType,
			Schedule:      apiModels.Schedule{Regularity: apiModels.RunNow},
		}
		validRunNowTaskBadParams = models.Task{
			TargetsByType: targetManagedEndpoint,
			OriginID:      taskID,
			Type:          config.ScriptTaskType,
			Parameters:    `WRONG`,
			Schedule:      apiModels.Schedule{Regularity: apiModels.RunNow},
		}

		validRunNowTaskDG = models.Task{
			TargetsByType: targetDG,
			OriginID:      taskID,
			Type:          config.ScriptTaskType,
			Schedule:      apiModels.Schedule{Regularity: apiModels.RunNow},
		}

		sitesSlice = []int64{1, 2}
		sitesData  = siteData{[]site{{1}, {2}}}
		siteIDsURL = fmt.Sprintf("%s/partner/%s/sites", config.Config.SitesMsURL, partnerID)
	)

	testCases := []struct {
		name                    string
		URL                     string
		method                  string
		isNeedUIDHeader         bool
		isNeedSiteIDs           bool
		mockTargetsRepo         mock.TargetRepoConf
		mockUserSites           mock.UserSitesConf
		mockTemplateCache       mock.TemplateCacheConf
		mockTaskPersistence     mock.TaskPersistenceConf
		mockTaskInstance        mock.TaskInstancePersistenceConf
		mockExecutionExpiration mock.ExecutionExpirationPersistenceConf
		mockTaskDef             mock.TaskDefinitionConf
		mockUserService         user.Service
		mockTaskCounter         mock.TaskCounterConf
		bodyToSend              models.Task
		expectedCode            int
		expectedErrorMessage    string
		expectedBody            models.Task
	}{
		{
			name:                 "testCase 1 - cannot parse taskID",
			URL:                  "/" + partnerID + "/tasks/bad",
			method:               http.MethodPost,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTaskIDHasBadFormat,
		},
		{
			name:       "testCase 2 - Error while getting task by id",
			URL:        "/" + partnerID + "/tasks/" + taskIDstr,
			method:     http.MethodPost,
			bodyToSend: validRunNowTaskWithObject,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{}, errors.New("err"))
				return tp
			},
			isNeedUIDHeader:      true,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
			expectedCode:         http.StatusInternalServerError,
		},
		{
			name:       "testCase 3 - Error while getting task by id",
			URL:        "/" + partnerID + "/tasks/" + taskIDstr,
			method:     http.MethodPost,
			bodyToSend: validRunNowTask,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{}, errors.New("err"))
				return tp
			},
			isNeedUIDHeader:      true,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
			expectedCode:         http.StatusInternalServerError,
		},
		{
			name:       "testCase 4 - Error while getting task by id, not found",
			URL:        "/" + partnerID + "/tasks/" + taskIDstr,
			method:     http.MethodPost,
			bodyToSend: validRunNowTask,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{}, gocql.ErrNotFound)
				return tp
			},
			isNeedUIDHeader:      true,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
			expectedCode:         http.StatusNotFound,
		},
		{
			name:       "testCase 5 - Error while inserting task",
			URL:        "/" + partnerID + "/tasks/" + taskIDstr,
			method:     http.MethodPost,
			bodyToSend: validRunNowTask,
			// Here GetByIDs returns 2 tasks with active statuses and then we're expecting
			// this tasks with inactive statuses to be inserted but with Cassandra error while inserting
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{
						PartnerID: partnerID,
						Name:      "Task",
						ID:        taskID,
						State:     statuses.TaskStateActive},
						{Name: "task2", State: statuses.TaskStateActive}},
						nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return tp
			},
			isNeedUIDHeader:      true,
			expectedErrorMessage: errorcode.ErrorCantSaveTaskToDB,
			expectedCode:         http.StatusInternalServerError,
		},
		{
			name:       "testCase 6 - Error getting siteIDs, no responder",
			URL:        "/" + partnerID + "/tasks/" + taskIDstr,
			method:     http.MethodPost,
			bodyToSend: validRunNowTaskDG,
			// Here GetByIDs returns 2 tasks with active statuses and then we're expecting
			// this tasks with inactive statuses to be inserted but with Cassandra error while inserting
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{validRunNowTask}, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			isNeedUIDHeader:      true,
			expectedErrorMessage: errorcode.ErrorCantCreateNewTask,
			expectedCode:         http.StatusInternalServerError,
		},
		{
			name:       "testCase 7 - Cant insert sites",
			URL:        "/" + partnerID + "/tasks/" + taskIDstr,
			method:     http.MethodPost,
			bodyToSend: validRunNowTaskDG,
			// Here GetByIDs returns 2 we're expecting
			// this tasks to be inserted, then we get sideIDs and trying to insert sites - error
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{}}, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uid, sitesSlice).
					Return(errors.New("err"))
				return us
			},
			isNeedUIDHeader:      true,
			isNeedSiteIDs:        true,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantCreateNewTask,
		},
		{
			name:       "testCase 8 - Cannot get templates for partner TemplateNotFoundError",
			URL:        "/" + partnerID + "/tasks/" + taskIDstr,
			method:     http.MethodPost,
			bodyToSend: validRunNowTask,
			// Here GetByIDs returns a task with and then we're inserting inactive task
			// then we're trying to ge templates but get notFoundError
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{Name: "name", Type: validRunNowTask.Type, PartnerID: partnerID, OriginID: taskID}}, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), gomock.Any(), gomock.Any(), sitesSlice).
					Return(nil)
				return us
			},
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, taskID, false).
					Return(models.TemplateDetails{}, models.TemplateNotFoundError{})
				return tc
			},
			isNeedUIDHeader:      true,
			isNeedSiteIDs:        true,
			expectedCode:         http.StatusNotFound,
			expectedErrorMessage: errorcode.ErrorCantGetTaskDefinitionTemplate,
		},
		{
			name:       "testCase 9 - Cannot get templates for partner, db error",
			URL:        "/" + partnerID + "/tasks/" + taskIDstr,
			method:     http.MethodPost,
			bodyToSend: validRunNowTask,
			// Here GetByIDs returns a task with and then we're inserting inactive task
			// then we're trying to ge templates but get Cassandra error
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{Name: "name", Type: validRunNowTask.Type, PartnerID: partnerID, OriginID: taskID}}, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, taskID, false).
					Return(models.TemplateDetails{}, errors.New("err"))
				return tc
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), gomock.Any(), gomock.Any(), sitesSlice).
					Return(nil)
				return us
			},
			isNeedUIDHeader:      true,
			isNeedSiteIDs:        true,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTemplatesForExecutionMS,
		},
		{
			name:       "testCase 10 - Cannot umarshall params",
			URL:        "/" + partnerID + "/tasks/" + taskIDstr,
			method:     http.MethodPost,
			bodyToSend: validRunNowTaskBadParams,
			// Here GetByIDs returns a task with no task type and then we're inserting inactive task
			// but we cant because there is no task type ValidateParametersField
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{Name: "name", PartnerID: partnerID, OriginID: taskID}}, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), gomock.Any(), gomock.Any(), sitesSlice).
					Return(nil)
				return us
			},
			isNeedUIDHeader:      true,
			isNeedSiteIDs:        true,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantDecodeInputData,
		},
		{
			name:       "testCase 11 - Cannot insert taskSites",
			URL:        "/" + partnerID + "/tasks/" + taskIDstr,
			method:     http.MethodPost,
			bodyToSend: validRunNowTask,
			// Here GetByIDs returns a task with no task type and then we're inserting inactive task
			// but we cant because there is no task type ValidateParametersField
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{Name: "name", PartnerID: partnerID, OriginID: taskID, Parameters: `{"num":"1"}`}}, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any())
				return tp
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), gomock.Any(), gomock.Any(), sitesSlice).
					Return(nil)
				us.EXPECT().InsertSitesByTaskID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("err"))
				return us
			},
			isNeedUIDHeader:      true,
			isNeedSiteIDs:        true,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantCreateNewTask,
		},
		{
			name:       "testCase 12 - Error while saving and processing task, InsertOrUpdate error",
			URL:        "/" + partnerID + "/tasks/" + taskIDstr,
			method:     http.MethodPost,
			bodyToSend: validRunNowTask,
			//Here GetByIDs returns a task with no task type and then we're inserting inactive task
			//then we get template cache, inserting new TaskInstance but cant insert edited task, cassandra error
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{Name: "name", PartnerID: partnerID, OriginID: taskID, Parameters: `{"num":"1"}`}}, nil)
				//here we can see, that we're inserting task 2 times: 1st - to make it inactive, and then to actually edit it
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any())
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return tp
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), gomock.Any(), gomock.Any(), sitesSlice).
					Return(nil)
				us.EXPECT().InsertSitesByTaskID(gomock.Any(), gomock.Any(), gomock.Any())
				return us
			},
			isNeedUIDHeader:      true,
			isNeedSiteIDs:        true,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantSaveTaskToDB,
		},
		{
			name:       "testCase 13 - Error while saving and processing task, InsertExecutionExpiration error",
			URL:        "/" + partnerID + "/tasks/" + taskIDstr,
			method:     http.MethodPost,
			bodyToSend: validRunNowTaskEmptyParams,
			//Here GetByIDs returns a task with no task type and then we're inserting inactive task
			//then we get template cache, inserting new TaskInstance inserting edited task, Get template detailt again (!)
			//inserting new task Instance
			//save execution expiration but get an error while inserting them
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{Name: "name", PartnerID: partnerID, OriginID: taskID}}, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), gomock.Any(), gomock.Any(), sitesSlice).
					Return(nil)
				us.EXPECT().
					InsertSitesByTaskID(gomock.Any(), gomock.Any(), gomock.Any())
				return us
			},
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {

				tc.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any())
				return tc
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).Times(1).
					Return(nil)
				return ti
			},
			mockExecutionExpiration: func(ee *mock.MockExecutionExpirationPersistence) *mock.MockExecutionExpirationPersistence {
				ee.EXPECT().
					InsertExecutionExpiration(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return ee
			},
			isNeedUIDHeader:      true,
			isNeedSiteIDs:        true,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantPrepareTaskForSendingOnExecution,
		},
		{
			name:       "testCase 14 - Successes but with counter Increase err",
			URL:        "/" + partnerID + "/tasks/" + taskIDstr,
			method:     http.MethodPost,
			bodyToSend: validRunNowTaskEmptyParams,
			//Here GetByIDs returns a task with no task type and then we're inserting inactive task
			//then we get template cache, inserting new TaskInstance inserting edited task, Get template details again (!)
			//inserting new task Instance, save execution expiration
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{Name: "name", PartnerID: partnerID, Type: models.TaskTypeScript, OriginID: taskID}}, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), gomock.Any(), gomock.Any(), sitesSlice).
					Return(nil)
				us.EXPECT().
					InsertSitesByTaskID(gomock.Any(), gomock.Any(), gomock.Any())
				return us
			},
			mockTargetsRepo: func(us *mock.MockTargetsRepo) *mock.MockTargetsRepo {
				us.EXPECT().Insert(partnerID, gomock.Any(), gomock.Any())
				return us
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).Times(1).
					Return(nil)
				return ti
			},
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, taskID, false).
					Return(models.TemplateDetails{}, nil)
				tc.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any())
				return tc
			},
			mockExecutionExpiration: func(ee *mock.MockExecutionExpirationPersistence) *mock.MockExecutionExpirationPersistence {
				ee.EXPECT().
					InsertExecutionExpiration(gomock.Any(), gomock.Any()).
					Return(nil)
				return ee
			},
			mockTaskCounter: func(tc *mock.MockTaskCounterPersistence) *mock.MockTaskCounterPersistence {
				tc.EXPECT().
					IncreaseCounter(partnerID, gomock.Any(), false).
					Return(errors.New("err"))
				tc.EXPECT().
					DecreaseCounter(partnerID, gomock.Any(), false).
					Return(errors.New("err"))
				return tc
			},
			isNeedUIDHeader: true,
			isNeedSiteIDs:   true,
			expectedCode:    http.StatusOK,
			expectedBody: models.Task{
				OriginID:  validRunNowTask.OriginID,
				Type:      validRunNowTask.Type,
				PartnerID: partnerID,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, tc.mockTaskDef,
				nil,
				tc.mockTaskInstance,
				tc.mockTaskPersistence,
				tc.mockTemplateCache,
				nil,
				tc.mockExecutionExpiration,
				tc.mockUserSites,
				tc.mockUserService,
				nil, nil, tc.mockTargetsRepo,
				nil, nil)
			router := getTaskServiceRouter(service)

			body, err := json.Marshal(tc.bodyToSend)
			if err != nil {
				t.Errorf("Can't encode body")
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, bytes.NewReader(body))
			if tc.isNeedUIDHeader {
				r.Header.Add(common.InitiatedByHeader, uidHeader)
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			if tc.isNeedSiteIDs {
				customHTTPMock(siteIDsURL, http.MethodGet, sitesData, http.StatusOK)
			}

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
				var gotBody models.TaskDefinitionDetails
				err = json.Unmarshal(w.Body.Bytes(), &gotBody)
				if err != nil {
					t.Errorf("Cannot parse expected body")
				}

				if gotBody.OriginID != tc.expectedBody.OriginID {
					t.Errorf("Wanted %v but got %v", tc.expectedBody.OriginID, gotBody.OriginID)
				}
				if gotBody.PartnerID != tc.expectedBody.PartnerID {
					t.Errorf("Wanted %v but got %v", tc.expectedBody.PartnerID, gotBody.PartnerID)
				}
				if gotBody.Type != tc.expectedBody.Type {
					t.Errorf("Wanted %v but got %v", tc.expectedBody.Type, gotBody.Type)
				}
			}
		})
	}
}

func TestStopTaskInstanceExecutions(t *testing.T) {

	type stopExecutionPayload struct {
		InstanceIDs []gocql.UUID `json:"instanceIDs"`
	}

	var (
		validInstancesPayload            = stopExecutionPayload{InstanceIDs: []gocql.UUID{validTaskInstanceUUID}}
		validTaskInstanceStatusesPending = make(map[gocql.UUID]statuses.TaskInstanceStatus)
		validTaskIDs                     = make([]gocql.UUID, 0)
	)
	validTaskIDs = append(validTaskIDs, validUUID)

	testCases := []struct {
		name                 string
		URL                  string
		method               string
		mockTaskInstance     mock.TaskInstancePersistenceConf
		mockTaskPersistence  mock.TaskPersistenceConf
		bodyToSend           stopExecutionPayload
		isWantInvalidBody    bool
		expectedCode         int
		expectedErrorMessage string
	}{
		{
			name:                 "testCase 1 - Cannot parse payload",
			URL:                  "/" + partnerID + "/tasks/task-instances/stop",
			method:               http.MethodPut,
			isWantInvalidBody:    true,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantDecodeInputData,
		},
		{
			name:       "testCase 2 - Cannot parse payload",
			URL:        "/" + partnerID + "/tasks/task-instances/stop",
			method:     http.MethodPut,
			bodyToSend: validInstancesPayload,
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByIDs(gomock.Any(), validInstancesPayload.InstanceIDs[0]).
					Return([]models.TaskInstance{}, errors.New("err"))
				return ti
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskInstances,
		},
		{
			name:       "testCase 3 - got zero task instances",
			URL:        "/" + partnerID + "/tasks/task-instances/stop",
			method:     http.MethodPut,
			bodyToSend: validInstancesPayload,
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByIDs(gomock.Any(), validInstancesPayload.InstanceIDs[0]).
					Return([]models.TaskInstance{}, nil)
				return ti
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantGetTaskInstances,
		},
		{
			name:       "testCase 4 - Marked taskInstance with STOP, but then Insert returned ERROR",
			URL:        "/" + partnerID + "/tasks/task-instances/stop",
			method:     http.MethodPut,
			bodyToSend: validInstancesPayload,
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					// Here we're return taskInstance with scheduled status for a target and then TI
					GetByIDs(gomock.Any(), validInstancesPayload.InstanceIDs[0]).
					Return([]models.TaskInstance{
						{
							TaskID:   validUUID,
							ID:       validTaskInstanceUUID,
							Statuses: validTaskInstanceStatusesPending,
						},
					}, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
				return ti
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantUpdateTaskInstances,
		},
		{
			name:       "testCase 5 - Success",
			URL:        "/" + partnerID + "/tasks/task-instances/stop",
			method:     http.MethodPut,
			bodyToSend: validInstancesPayload,
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					// Here we're return taskInstance with scheduled status for a target and then TI
					GetByIDs(gomock.Any(), validInstancesPayload.InstanceIDs[0]).
					Return([]models.TaskInstance{
						{
							TaskID:   validUUID,
							ID:       validTaskInstanceUUID,
							Statuses: validTaskInstanceStatusesPending,
						},
					}, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// because of pointers
			validTaskInstanceStatusesPending[validTaskInstanceUUID] = statuses.TaskInstancePending

			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, nil, nil, tc.mockTaskInstance, tc.mockTaskPersistence, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil)
			router := getTaskServiceRouter(service)

			body, err := json.Marshal(tc.bodyToSend)
			if err != nil {
				t.Errorf("Can't encode body")
			}

			// making body invalid adding wrong symbol
			if tc.isWantInvalidBody {
				body = []byte(`{`)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, bytes.NewReader(body))

			router.ServeHTTP(w, r)
			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
		})
	}
}

func TestCreate(t *testing.T) {

	config.Config.TaskTypes[sequenceType] = "http://sequence.elb"
	sitesURL := fmt.Sprintf("%s/partner/%s/sites", config.Config.SitesMsURL, partnerID)

	validRunNowTaskDG := models.Task{
		TargetsByType: targetDG,
		OriginID:      taskID,
		Type:          config.ScriptTaskType,
		Schedule:      apiModels.Schedule{Regularity: apiModels.RunNow},
	}

	validRunNowTask := models.Task{
		TargetsByType: targetManagedEndpoint,
		OriginID:      taskID,
		Type:          config.ScriptTaskType,
		Schedule:      apiModels.Schedule{Regularity: apiModels.RunNow},
	}

	validRunNowTaskSequenceType := models.Task{
		TargetsByType: targetManagedEndpoint,
		OriginID:      taskID,
		Type:          sequenceType,
		Schedule:      apiModels.Schedule{Regularity: apiModels.RunNow, EndRunTime: time.Now().AddDate(100, 0, 0)},
	}

	validRunNowTaskSequenceTypeOldTarget := models.Task{
		Targets: models.Target{
			IDs:  targetManagedEndpoint[models.ManagedEndpoint],
			Type: models.ManagedEndpoint,
		},
		OriginID: taskID,
		Type:     sequenceType,
		Schedule: apiModels.Schedule{Regularity: apiModels.RunNow, EndRunTime: time.Now().AddDate(100, 0, 0)},
	}

	sitesSlice := []int64{1, 2}
	sitesData := siteData{[]site{{1}, {2}}}
	ctx := context.Background()
	userCTX := context.WithValue(ctx, config.UserKeyCTX, entities.User{PartnerID: partnerID, UID: uidHeader, Token: "ValidToken"})

	testCases := []struct {
		name                             string
		URL                              string
		method                           string
		isNeedUIDHeader                  bool
		isNeedSitesData                  bool
		isNeedManagedEndpointsData       bool
		isNeedEmptyCtx                   bool
		managedEndpointHTTPMockResponder func(req *http.Request) (*http.Response, error)
		userUC                           mock.UserUC
		mockUserService                  user.Service
		mockUserSites                    mock.UserSitesConf
		mockTargets                      mock.TargetRepoConf
		mockTemplateCache                mock.TemplateCacheConf
		mockTaskPersistence              mock.TaskPersistenceConf
		mockTaskCounter                  mock.TaskCounterPersistenceConf
		mockTaskInstance                 mock.TaskInstancePersistenceConf
		mockExecutionExpiration          mock.ExecutionExpirationPersistenceConf
		mockDGRepo                       mock.DGRepoConf
		mockSiteRepo                     mock.SitesRepoConf
		bodyToSend                       models.Task
		expectedCode                     int
		expectedErrorMessage             string
	}{
		{
			name:   "testCase 0 - Bad task format, Recurrent task with no details",
			URL:    "/" + partnerID + "/tasks",
			method: http.MethodPost,
			bodyToSend: models.Task{
				TargetsByType: targetManagedEndpoint,
				OriginID:      taskID,
				Type:          config.ScriptTaskType,
				Schedule:      apiModels.Schedule{Regularity: apiModels.Recurrent},
			},
			isNeedUIDHeader:      true,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantDecodeInputData,
		},
		{
			name:            "testCase 1 - can't get user",
			URL:             "/" + partnerID + "/tasks",
			method:          http.MethodPost,
			bodyToSend:      validRunNowTaskDG,
			isNeedUIDHeader: true,
			expectedCode:    http.StatusInternalServerError,
		},
		{
			name:            "testCase 2 - Internal error, can't get sites data",
			URL:             "/" + partnerID + "/tasks",
			method:          http.MethodPost,
			bodyToSend:      validRunNowTaskDG,
			isNeedUIDHeader: true,
			isNeedEmptyCtx:  true,
			expectedCode:    http.StatusBadRequest,
		},
		{
			name:            "testCase 3 - error while inserting user sites",
			URL:             "/" + partnerID + "/tasks",
			method:          http.MethodPost,
			bodyToSend:      validRunNowTaskDG,
			isNeedUIDHeader: true,
			isNeedSitesData: true,
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uidHeader, sitesSlice).
					Return(errors.New("some error"))
				return us
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantCreateNewTask,
		},
		{
			name:   "testCase 4 - Bad parameters field for sequence type",
			URL:    "/" + partnerID + "/tasks",
			method: http.MethodPost,
			bodyToSend: models.Task{
				TargetsByType: targetDG,
				OriginID:      taskID,
				Type:          sequenceType,
				Parameters:    `bad parameters`,
				Schedule:      apiModels.Schedule{Regularity: apiModels.RunNow},
			},
			// we use userSites here because we're sending DG
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uidHeader, sitesSlice).
					Return(nil)
				return us
			},
			isNeedUIDHeader:      true,
			isNeedSitesData:      true,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantDecodeInputData,
		},
		{
			name:       "testCase 6 - Cannot get templates for partner",
			URL:        "/" + partnerID + "/tasks",
			method:     http.MethodPost,
			bodyToSend: validRunNowTask,
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, taskID, false).
					Return(models.TemplateDetails{}, errors.New("err"))
				return tc
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uidHeader, sitesSlice).
					Return(nil)
				return us
			},
			isNeedUIDHeader:      true,
			isNeedSitesData:      true,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTemplatesForExecutionMS,
		},
		{
			name:       "testCase 7 - Cannot get templates for partner TemplateNotFoundError",
			URL:        "/" + partnerID + "/tasks",
			method:     http.MethodPost,
			bodyToSend: validRunNowTask,
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, taskID, false).
					Return(models.TemplateDetails{}, models.TemplateNotFoundError{OriginID: taskID, PartnerID: partnerID})
				return tc
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uidHeader, sitesSlice).
					Return(nil)
				return us
			},
			isNeedUIDHeader:      true,
			isNeedSitesData:      true,
			expectedCode:         http.StatusNotFound,
			expectedErrorMessage: errorcode.ErrorCantGetTaskDefinitionTemplate,
		},
		{
			name:   "testCase 8 - Cannot validate parameters field, template json schema is empty but parameters don't",
			URL:    "/" + partnerID + "/tasks",
			method: http.MethodPost,
			bodyToSend: models.Task{
				TargetsByType: targetManagedEndpoint,
				OriginID:      taskID,
				Type:          config.ScriptTaskType,
				Parameters:    `{}`,
				Schedule:      apiModels.Schedule{Regularity: apiModels.RunNow},
			},
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, taskID, false).
					Return(models.TemplateDetails{PartnerID: partnerID, Name: "Name"}, nil)
				return tc
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uidHeader, sitesSlice).
					Return(nil)
				return us
			},
			isNeedUIDHeader:      true,
			isNeedSitesData:      true,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantDecodeInputData,
		},
		{
			name:   "testCase 9 - Can't get managed endpoints for the targets, error in DG repo",
			URL:    "/" + partnerID + "/tasks",
			method: http.MethodPost,
			bodyToSend: models.Task{
				TargetsByType: targetDG,
				OriginID:      taskID,
				Type:          sequenceType,
				Schedule:      apiModels.Schedule{Regularity: apiModels.RunNow},
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uidHeader, sitesSlice).
					Return(nil)
				return us
			},
			mockDGRepo: func(dg *mock.MockDynamicGroupRepo) *mock.MockDynamicGroupRepo {
				dg.EXPECT().
					GetEndpointsByGroupIDs(gomock.Any(), gomock.Any(), uidHeader, partnerID, false).
					Return([]gocql.UUID{}, errors.New("some error"))
				return dg
			},
			isNeedUIDHeader:      true,
			isNeedSitesData:      true,
			expectedErrorMessage: errorcode.ErrorCantOpenTargets,
			expectedCode:         http.StatusInternalServerError,
		},
		{
			name:   "testCase 11 - Can't get managed endpoints for the targets, error in Sites repo",
			URL:    "/" + partnerID + "/tasks",
			method: http.MethodPost,
			bodyToSend: models.Task{
				TargetsByType: targetSitesFixedUUID,
				OriginID:      taskID,
				Type:          sequenceType,
				Schedule:      apiModels.Schedule{Regularity: apiModels.RunNow},
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uidHeader, sitesSlice).
					Return(nil)
				return us
			},
			mockSiteRepo: func(sr *mock.MockSiteRepo) *mock.MockSiteRepo {
				sr.EXPECT().GetEndpointsBySiteIDs(gomock.Any(),partnerID, gomock.Any()).
					Return([]gocql.UUID{}, errors.New("some error"))
				return sr
			},
			isNeedUIDHeader:            true,
			isNeedSitesData:            true,
			isNeedManagedEndpointsData: true,
			expectedErrorMessage:       errorcode.ErrorCantOpenTargets,
			expectedCode:               http.StatusInternalServerError,
		},
		{
			name:   "testCase 12 - Can't get managed endpoints for the targets, empty endpoints from Sites repo",
			URL:    "/" + partnerID + "/tasks",
			method: http.MethodPost,
			bodyToSend: models.Task{
				TargetsByType: targetSitesFixedUUID,
				OriginID:      taskID,
				Type:          sequenceType,
				Schedule:      apiModels.Schedule{Regularity: apiModels.RunNow},
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uidHeader, sitesSlice).
					Return(nil)
				return us
			},
			mockSiteRepo: func(sr *mock.MockSiteRepo) *mock.MockSiteRepo {
				sr.EXPECT().GetEndpointsBySiteIDs(gomock.Any(), partnerID, gomock.Any()).
					Return([]gocql.UUID{}, nil)
				return sr
			},
			isNeedUIDHeader:            true,
			isNeedSitesData:            true,
			isNeedManagedEndpointsData: true,
			expectedErrorMessage:       errorcode.ErrorNoEndpointsForTargets,
			expectedCode:               http.StatusInternalServerError,
		},
		{
			name:       "testCase 14 - Can't create a task, cassandra error",
			URL:        "/" + partnerID + "/tasks",
			method:     http.MethodPost,
			bodyToSend: validRunNowTask,
			// Here we're getting templates and then inserting a TaskInstance, inserting a new task but getting an error
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().
					GetByOriginID(gomock.Any(), partnerID, taskID, false)
				return tc
			},
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(errors.New("Cassandra err"))
				return tp
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uidHeader, sitesSlice).
					Return(nil)
				us.EXPECT().
					InsertSitesByTaskID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return us
			},
			isNeedUIDHeader:      true,
			isNeedSitesData:      true,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantSaveTaskToDB,
		},
		{
			name:       "testCase 15 - Can't send task, cassandra err while inserting execution expirations",
			URL:        "/" + partnerID + "/tasks",
			method:     http.MethodPost,
			bodyToSend: validRunNowTaskSequenceType,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().Insert(gomock.Any(), gomock.Any())
				return ti
			},
			mockExecutionExpiration: func(ee *mock.MockExecutionExpirationPersistence) *mock.MockExecutionExpirationPersistence {
				ee.EXPECT().
					InsertExecutionExpiration(gomock.Any(), gomock.Any()).
					Return(errors.New("Cassandra err"))
				return ee
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uidHeader, sitesSlice).
					Return(nil)
				us.EXPECT().
					InsertSitesByTaskID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return us
			},
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any())
				return tc
			},
			isNeedUIDHeader:      true,
			isNeedSitesData:      true,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantPrepareTaskForSendingOnExecution,
		},
		{
			name:       "testCase 16 - Can't prepare task for sending to Scripting MS, cassandra err while inserting task instance",
			URL:        "/" + partnerID + "/tasks",
			method:     http.MethodPost,
			bodyToSend: validRunNowTaskSequenceType,
			// Here we're inserting a TaskInstance, inserting a new task, inserting taskInstance for SaveAndProcess and get err
			// in second goroutine in saveExecutionExpiration we get Templates and inserting ExecExpiration
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return ti
			},
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			mockExecutionExpiration: func(ee *mock.MockExecutionExpirationPersistence) *mock.MockExecutionExpirationPersistence {
				ee.EXPECT().
					InsertExecutionExpiration(gomock.Any(), gomock.Any()).
					Return(nil)
				return ee
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uidHeader, sitesSlice).
					Return(nil)
				us.EXPECT().
					InsertSitesByTaskID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return us
			},
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any())
				return tc
			},
			isNeedUIDHeader:      true,
			isNeedSitesData:      true,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantPrepareTaskForSendingOnExecution,
		},
		{
			name:       "testCase 17 - can't insert task sites",
			URL:        "/" + partnerID + "/tasks",
			method:     http.MethodPost,
			bodyToSend: validRunNowTaskSequenceType,
			// Here we're inserting a TaskInstance, inserting a new task, inserting taskInstance for SaveAndProcess and get err
			// in second goroutine in saveExecutionExpiration we get Templates and inserting ExecExpiration
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uidHeader, sitesSlice).
					Return(nil)
				us.EXPECT().
					InsertSitesByTaskID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return us
			},
			isNeedUIDHeader: true,
			isNeedSitesData: true,
			expectedCode:    http.StatusInternalServerError,
		},
		{
			name:       "testCase 18 - Created",
			URL:        "/" + partnerID + "/tasks",
			method:     http.MethodPost,
			bodyToSend: validRunNowTaskSequenceType,
			// Here we're inserting a TaskInstance, inserting a new task
			// inserting InsertExecutionExpiration but get templates error
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil).Times(1)
				return ti
			},
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			mockExecutionExpiration: func(ee *mock.MockExecutionExpirationPersistence) *mock.MockExecutionExpirationPersistence {
				ee.EXPECT().
					InsertExecutionExpiration(gomock.Any(), gomock.Any()).
					Return(nil)
				return ee
			},
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any())
				return tc
			},
			mockTaskCounter: func(tc *mock.MockTaskCounterPersistence) *mock.MockTaskCounterPersistence {
				tc.EXPECT().
					IncreaseCounter(partnerID, gomock.Any(), false).
					Return(nil)
				return tc
			},
			mockTargets: func(us *mock.MockTargetsRepo) *mock.MockTargetsRepo {
				us.EXPECT().Insert(partnerID, gomock.Any(), gomock.Any())
				return us
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uidHeader, sitesSlice).
					Return(nil)
				us.EXPECT().
					InsertSitesByTaskID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return us
			},
			userUC: func(u *mock.MockUserUC) *mock.MockUserUC {
				u.EXPECT().EndpointsFromAsset(gomock.Any(), gomock.Any()).Return([]entities.Endpoints{}, nil).AnyTimes()
				u.EXPECT().SaveEndpoints(gomock.Any(), gomock.Any()).AnyTimes()
				return u
			},
			isNeedUIDHeader: true,
			isNeedSitesData: true,
			expectedCode:    http.StatusCreated,
			// cant check returned body here because there is a time.Now() function used there
		},
		{
			name:       "testCase 19 - Created by mono-target request",
			URL:        "/" + partnerID + "/tasks",
			method:     http.MethodPost,
			bodyToSend: validRunNowTaskSequenceTypeOldTarget,
			// Here we're inserting a TaskInstance, inserting a new task
			// inserting InsertExecutionExpiration but get templates error
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil).Times(1)
				return ti
			},
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			mockExecutionExpiration: func(ee *mock.MockExecutionExpirationPersistence) *mock.MockExecutionExpirationPersistence {
				ee.EXPECT().
					InsertExecutionExpiration(gomock.Any(), gomock.Any()).
					Return(nil)
				return ee
			},
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any())
				return tc
			},
			mockTaskCounter: func(tc *mock.MockTaskCounterPersistence) *mock.MockTaskCounterPersistence {
				tc.EXPECT().
					IncreaseCounter(partnerID, gomock.Any(), false).
					Return(nil)
				return tc
			},
			mockTargets: func(us *mock.MockTargetsRepo) *mock.MockTargetsRepo {
				us.EXPECT().Insert(partnerID, gomock.Any(), gomock.Any())
				return us
			},
			mockUserSites: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					InsertUserSites(gomock.Any(), partnerID, uidHeader, sitesSlice).
					Return(nil)
				us.EXPECT().
					InsertSitesByTaskID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return us
			},
			userUC: func(u *mock.MockUserUC) *mock.MockUserUC {
				u.EXPECT().EndpointsFromAsset(gomock.Any(), gomock.Any()).Return([]entities.Endpoints{}, nil).AnyTimes()
				u.EXPECT().SaveEndpoints(gomock.Any(), gomock.Any()).AnyTimes()
				return u
			},
			isNeedUIDHeader: true,
			isNeedSitesData: true,
			expectedCode:    http.StatusCreated,
			// cant check returned body here because there is a time.Now() function used there
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, nil,
				nil,
				tc.mockTaskInstance,
				tc.mockTaskPersistence,
				tc.mockTemplateCache,
				nil,
				tc.mockExecutionExpiration,
				tc.mockUserSites,
				tc.mockUserService,
				tc.userUC, nil, tc.mockTargets, tc.mockDGRepo, tc.mockSiteRepo)
			router := getTaskServiceRouter(service)

			body, err := json.Marshal(tc.bodyToSend)
			if err != nil {
				t.Errorf("Can't encode body")
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, bytes.NewReader(body))

			if !tc.isNeedEmptyCtx {
				r = r.WithContext(userCTX)
			}

			if tc.isNeedUIDHeader {
				r.Header.Add(common.InitiatedByHeader, uidHeader)
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			if tc.isNeedSitesData {
				customHTTPMock(sitesURL, http.MethodGet, sitesData, http.StatusOK)
			}

			router.ServeHTTP(w, r)
			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
		})
	}
}

func TestSendTaskOnExecution(t *testing.T) {
	config.Config.TaskTypes["sequence"] = "sequence"
	logger.Load(config.Config.Log)
	ctx := context.Background()

	validTasks := []models.Task{{
		PartnerID: partnerID,
		Type:      "sequence",
		OriginID:  validTaskInstanceUUID,
	}}

	statusesMap := make(map[gocql.UUID]statuses.TaskInstanceStatus)
	statusesMap[validTaskInstanceUUID] = statuses.TaskInstanceRunning

	validTaskInstance := models.TaskInstance{PartnerID: partnerID,
		TaskID:        taskID,
		OverallStatus: statuses.TaskInstanceRunning,
		Statuses:      statusesMap,
	}

	executionExpirationMockNoErr := func(ee *mock.MockExecutionExpirationPersistence) *mock.MockExecutionExpirationPersistence {
		ee.EXPECT().InsertExecutionExpiration(gomock.Any(), gomock.Any())
		return ee
	}

	executionExpirationMockErr := func(ee *mock.MockExecutionExpirationPersistence) *mock.MockExecutionExpirationPersistence {
		ee.EXPECT().
			InsertExecutionExpiration(gomock.Any(), gomock.Any()).
			Return(errors.New("db err"))
		return ee
	}

	testCases := []struct {
		name               string
		tasksToSend        []models.Task
		taskInstanceToSend models.TaskInstance
		mockTemplateCache  mock.TemplateCacheConf
		mockTaskInstance   mock.TaskInstancePersistenceConf
		mockExecExpiration mock.ExecutionExpirationPersistenceConf
		err                error
	}{
		{
			name: "testCase 0 - len of task is less than 1",
			err:  fmt.Errorf("empty tasks"),
		},
		{
			name: "testCase 1 - got error from saveExecutionExpiration because taskType == script",
			tasksToSend: []models.Task{{
				PartnerID: partnerID,
				Type:      "script",
				OriginID:  validTaskInstanceUUID,
			}},
			taskInstanceToSend: validTaskInstance,
			// We're getting templates, inserting execution expiration, but get error  while inserting Task Instance
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any())
				return tc
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			mockExecExpiration: executionExpirationMockErr,
			err:                fmt.Errorf("db err"),
		},
		{
			name:               "testCase 2 - got error while inserting exec expiration",
			tasksToSend:        validTasks,
			taskInstanceToSend: validTaskInstance,
			// We're getting templates, inserting execution expiration, but get error  while inserting Task Instance
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any())
				return tc
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			mockExecExpiration: executionExpirationMockErr,
			err:                fmt.Errorf("db err"),
		},
		{
			name:               "testCase 3 - got error while inserting task instance",
			tasksToSend:        validTasks,
			taskInstanceToSend: validTaskInstance,
			// We're getting templates, inserting execution expiration, but get error  while inserting Task Instance
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any())
				return tc
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(errors.New("db err"))
				return ti
			},
			mockExecExpiration: executionExpirationMockNoErr,
			err:                fmt.Errorf("db err"),
		},
		{
			name: "testCase 4 - got error from getExecutionURL",
			tasksToSend: []models.Task{{
				PartnerID: partnerID,
				Type:      "invalid",
				OriginID:  validTaskInstanceUUID,
			}},
			taskInstanceToSend: validTaskInstance,
			// We're getting templates, inserting execution expiration and inserting Task Instance
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any())
				return tc
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			mockExecExpiration: executionExpirationMockNoErr,
			err:                fmt.Errorf("unknown execution MS type: invalid"),
		},
		{
			name:               "testCase 5 - No err",
			tasksToSend:        validTasks,
			taskInstanceToSend: validTaskInstance,
			// We're getting templates, inserting execution expiration and inserting Task Instance
			mockTemplateCache: func(tc *mock.MockTemplateCache) *mock.MockTemplateCache {
				tc.EXPECT().CalculateExpectedExecutionTimeSec(gomock.Any(), gomock.Any())
				return tc
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			mockExecExpiration: executionExpirationMockNoErr,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			taskService := getMockedTaskService(mockController, nil,
				nil,
				tc.mockTaskInstance,
				nil,
				tc.mockTemplateCache,
				nil,
				tc.mockExecExpiration,
				nil,
				nil,
				nil, nil, nil, nil, nil)

			gotErr := taskService.SendTaskOnExecution(ctx, tc.tasksToSend, tc.taskInstanceToSend)
			if tc.err != nil && gotErr.Error() != tc.err.Error() {
				t.Errorf("Wanted %v but got %v", tc.err, gotErr)
			}
		})
	}
}

func TestGetByPartnerAndManagedEndpointID(t *testing.T) {
	managedEndpoints := apiModels.ManagedEndpoint{ID: taskIDstr}
	expectedBody := []models.Task{
		{
			ID:               taskID,
			OriginID:         taskID,
			State:            3,
			ManagedEndpoints: []models.ManagedEndpointDetailed{{ManagedEndpoint: managedEndpoints, State: 0}},
		},
	}

	expectedBodyJSON, err := json.Marshal(expectedBody)
	if err != nil {
		t.Errorf("Can't marshal custom script: %v", expectedBody)
		return
	}

	testCases := []struct {
		name                 string
		method               string
		URL                  string
		expectedBody         interface{}
		expectedErrorMessage string
		expectedCode         int
		userServiceMock      user.Service
		taskPersistenceMock  mock.TaskPersistenceConf
	}{
		{
			name:                 "testCase 1 - bad managedEndpointID",
			URL:                  "/" + partnerID + "/tasks/managed-endpoints/badManagedEndpointID",
			method:               http.MethodGet,
			expectedErrorMessage: errorcode.ErrorEndpointIDHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name:         "testCase 2 - Good",
			URL:          "/" + partnerID + "/tasks/managed-endpoints/" + taskID.String(),
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: string(expectedBodyJSON),
			// First we're Get tasks ByPartnerAndManagedEndpointID with no err, and then get lists of internalTasks by IDs
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByPartnerAndManagedEndpointID(gomock.Any(), partnerID, taskID, common.UnlimitedCount).
					Return([]models.Task{}, nil)
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, gomock.Any()).
					Return([]models.Task{
						{
							ID:                 taskID,
							OriginID:           taskID,
							LastTaskInstanceID: taskID,
							ManagedEndpointID:  taskID,
						}}, nil)
				return tp
			},
		},
		{
			name:                 "testCase 3 - can't get list of Tasks",
			URL:                  "/" + partnerID + "/tasks/managed-endpoints/" + taskID.String(),
			method:               http.MethodGet,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetListOfTasksByManagedEndpoint,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByPartnerAndManagedEndpointID(gomock.Any(), partnerID, taskID, common.UnlimitedCount).
					Return([]models.Task{}, errors.New("error"))
				return tp
			},
		},
		{
			name:                 "testCase 4 - StatusInternalServerError",
			URL:                  "/" + partnerID + "/tasks/managed-endpoints/" + taskID.String(),
			method:               http.MethodGet,
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetListOfTasksByManagedEndpoint,
			// First we're Get tasks ByPartnerAndManagedEndpointID with no err, and then get lists of internalTasks by IDs but get err
			taskPersistenceMock: func(ti *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				ti.EXPECT().
					GetByPartnerAndManagedEndpointID(gomock.Any(), partnerID, taskID, common.UnlimitedCount).
					Return([]models.Task{}, nil)
				ti.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, gomock.Any()).
					Return([]models.Task{
						{
							ID:                 taskID,
							OriginID:           taskID,
							LastTaskInstanceID: taskID,
							ManagedEndpointID:  taskID,
						}}, errors.New("error"))
				return ti
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, nil, nil, nil, tc.taskPersistenceMock, nil, nil,
				nil, nil, tc.userServiceMock, nil, nil, nil, nil, nil)
			router := getTaskServiceRouter(service)

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, nil)
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}

			if w.Body.String() != tc.expectedBody && tc.expectedBody != nil {
				t.Errorf("Wanted body: %v but got: %v", tc.expectedBody, w.Body)
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
		})
	}
}

func TestGetTasksSummaryData(t *testing.T) {
	expectedBody := []models.TaskSummaryData{{
		TaskID: taskID,
	}}

	expectedBodyJSON, err := json.Marshal(expectedBody)
	if err != nil {
		t.Errorf("Can't marshal custom script: %v", expectedBody)
		return
	}

	testCases := []struct {
		name                       string
		method                     string
		URL                        string
		expectedBody               interface{}
		expectedErrorMessage       string
		expectedCode               int
		taskSummaryPersistenceMock mock.TaskSummaryPersistenceConf
		userServiceMock            user.Service
	}{
		{
			name:         "testCase 1 - good",
			URL:          "/" + partnerID + "/tasks/data",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			taskSummaryPersistenceMock: func(ts *mock.MockTaskSummaryPersistence) *mock.MockTaskSummaryPersistence {
				ts.EXPECT().
					GetTasksSummaryData(gomock.Any(), false, nil, partnerID).
					Return([]models.TaskSummaryData{{TaskID: taskID}}, nil)
				return ts
			},
			expectedBody: string(expectedBodyJSON),
		},
		{
			name:         "testCase 2 - bad",
			URL:          "/" + partnerID + "/tasks/data",
			method:       http.MethodGet,
			expectedCode: http.StatusInternalServerError,
			taskSummaryPersistenceMock: func(ts *mock.MockTaskSummaryPersistence) *mock.MockTaskSummaryPersistence {
				ts.EXPECT().
					GetTasksSummaryData(gomock.Any(), false, nil, partnerID).
					Return([]models.TaskSummaryData{}, errors.New("error"))
				return ts
			},
			expectedErrorMessage: errorcode.ErrorCantGetTasksSummaryData,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, nil, tc.taskSummaryPersistenceMock, nil, nil, nil,
				nil, nil, nil, tc.userServiceMock, nil, nil, nil, nil, nil)
			router := getTaskServiceRouter(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, nil)
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}

			if w.Body.String() != tc.expectedBody && tc.expectedBody != nil {
				t.Errorf("Wanted body: %v but got: %v", tc.expectedBody, w.Body)
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
		})
	}
}

func TestEnableTaskForAllTargetsEnableSelectedTargets(t *testing.T) {
	type EnableTaskBody struct {
		Active bool `json:"active"`
	}

	enableTaskBody := EnableTaskBody{Active: true}
	testCases := []struct {
		name                 string
		URL                  string
		method               string
		mockTaskPersistence  mock.TaskPersistenceConf
		bodyToSend           EnableTaskBody
		expectedCode         int
		expectedErrorMessage string
	}{
		{
			name:                 "testCase 1 - Cannot parse taskID for all targets",
			method:               http.MethodPut,
			URL:                  "/" + partnerID + "/tasks/" + taskID.String() + "baduuid" + "/enable",
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTaskIDHasBadFormat,
		},
		{
			name:                 "testCase 2 - Cannot parse body, empty for selected targets",
			method:               http.MethodPut,
			URL:                  "/" + partnerID + "/tasks/" + taskID.String() + "/enable/targets",
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantDecodeInputData,
		},
		{
			name:       "testCase 3 - Cannot update task, Cassandra error",
			method:     http.MethodPut,
			URL:        "/" + partnerID + "/tasks/" + taskID.String() + "/enable",
			bodyToSend: enableTaskBody,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					UpdateTask(gomock.Any(), gomock.Any(), partnerID, taskID).
					Return(errors.New("Cassandra error"))
				return tp
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
		},
		{
			name:       "testCase 4 - Cannot update task, TaskIsExpiredError",
			method:     http.MethodPut,
			URL:        "/" + partnerID + "/tasks/" + taskID.String() + "/enable",
			bodyToSend: enableTaskBody,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					UpdateTask(gomock.Any(), gomock.Any(), partnerID, taskID).
					Return(models.TaskIsExpiredError{PartnerID: partnerID, TaskID: taskID})
				return tp
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
		},
		{
			name:       "testCase 5 - Cannot update task, TaskNotFound",
			method:     http.MethodPut,
			URL:        "/" + partnerID + "/tasks/" + taskID.String() + "/enable",
			bodyToSend: enableTaskBody,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					UpdateTask(gomock.Any(), gomock.Any(), partnerID, taskID).
					Return(models.TaskNotFoundError{ErrorParameters: "TaskNotFound"})
				return tp
			},
			expectedCode:         http.StatusNotFound,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
		},
		{
			name:       "testCase 6 - Successes",
			method:     http.MethodPut,
			URL:        "/" + partnerID + "/tasks/" + taskID.String() + "/enable",
			bodyToSend: enableTaskBody,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					UpdateTask(gomock.Any(), gomock.Any(), partnerID, taskID).
					Return(nil)
				return tp
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, nil, nil, nil, tc.mockTaskPersistence,
				nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
			router := getTaskServiceRouter(service)

			body, err := json.Marshal(tc.bodyToSend)
			if err != nil {
				t.Errorf("Can't encode body")
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, bytes.NewReader(body))

			router.ServeHTTP(w, r)
			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
		})
	}
}

func TestGetByPartnerAndID(t *testing.T) {
	logger.Load(config.Config.Log)
	validRunNowTask := models.Task{
		TargetsByType: targetManagedEndpoint,
		OriginID:      taskID,
		Type:          config.ScriptTaskType,
		Schedule:      apiModels.Schedule{Regularity: apiModels.RunNow},
	}

	testCases := []struct {
		name                 string
		URL                  string
		method               string
		mockTaskPersistence  mock.TaskPersistenceConf
		userServiceMock      user.Service
		expectedCode         int
		expectedErrorMessage string
	}{
		{
			name:                 "testCase 1 - Cannot parse taskID",
			method:               http.MethodGet,
			URL:                  "/" + partnerID + "/tasks/" + taskID.String() + "badUUID",
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTaskIDHasBadFormat,
		},
		{
			name:   "testCase 2 - Cannot get task, Cassandra error",
			method: http.MethodGet,
			URL:    "/" + partnerID + "/tasks/" + taskID.String(),
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{}, errors.New("Cassandra error"))
				return tp
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
		},
		{
			name:   "testCase 3 - Cannot make taskOutput, db returned empty task slice",
			method: http.MethodGet,
			URL:    "/" + partnerID + "/tasks/" + taskID.String(),
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID)
				return tp
			},
			expectedCode:         http.StatusNotFound,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
		},
		{
			name:   "testCase 4 - Successes",
			method: http.MethodGet,
			URL:    "/" + partnerID + "/tasks/" + taskID.String(),
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{validRunNowTask}, nil)
				return tp
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, nil, nil, nil, tc.mockTaskPersistence, nil,
				nil, nil, nil, tc.userServiceMock, nil, nil, nil, nil, nil)
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
		})
	}
}

// For test coverage
func TestSendTaskOnExecutionREST(t *testing.T) {

	var mockController *gomock.Controller
	validTaskInstanceStatusesSuccesses := make(map[gocql.UUID]statuses.TaskInstanceStatus)
	validTaskInstanceStatusesSuccesses[validUUID] = statuses.TaskInstanceSuccess

	testCases := []struct {
		name                        string
		URL                         string
		mockTaskPersistence         mock.TaskPersistenceConf
		mockTaskInstancePersistence mock.TaskInstancePersistenceConf
		mockExecResPersistence      mock.ExecutionResultPersistenceConf
		mockExecResultUpdateService func() ExecutionResultUpdateUC
	}{
		{
			name: "testCase 1 - no errors",
			URL:  "/URL",
			// this flow is used in ProcessExecutionResults (which in getExecutionURL func, that is in executionResultsUpdate.ProcessExecutionResults)
			mockTaskInstancePersistence: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				return ti
			},
			mockExecResultUpdateService: func() ExecutionResultUpdateUC {
				executionResultUpdateUCMock := mockusecases.NewMockExecutionResultUpdateUC(mockController)
				executionResultUpdateUCMock.EXPECT().ProcessExecutionResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).Times(1)
				return executionResultUpdateUCMock
			},
		},
		{
			name: "testCase 2 -  error",
			URL:  "/URL",
			mockTaskInstancePersistence: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				return ti
			},
			mockExecResultUpdateService: func() ExecutionResultUpdateUC {
				executionResultUpdateUCMock := mockusecases.NewMockExecutionResultUpdateUC(mockController)
				executionResultUpdateUCMock.EXPECT().ProcessExecutionResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("some_err")).Times(1)
				return executionResultUpdateUCMock
			},
		},
	}

	config.Config.RetryStrategy.MaxNumberOfRetries = 1
	config.Config.RetryStrategy.RetrySleepIntervalSec = 1

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockController = gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, nil, nil, nil, nil, nil, tc.mockExecResPersistence, nil,
				nil, userServiceMockNoErr, nil, nil, nil, nil, nil)

			service.executionResultUpdateService = tc.mockExecResultUpdateService()
			// to mock global variables
			if tc.mockTaskInstancePersistence != nil {
				taskInstancePersistenceMock := mock.NewMockTaskInstancePersistence(mockController)
				taskInstancePersistenceMock = tc.mockTaskInstancePersistence(taskInstancePersistenceMock)
				models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
			}
			if tc.mockTaskPersistence != nil {
				taskPersistenceMock := mock.NewMockTaskPersistence(mockController)
				taskPersistenceMock = tc.mockTaskPersistence(taskPersistenceMock)
				models.TaskPersistenceInstance = taskPersistenceMock
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			customHTTPMock(tc.URL, http.MethodPost, nil, http.StatusBadRequest)

			ctx := context.TODO()
			service.sendTaskOnExecutionREST(ctx, tc.URL, apiModels.ExecutionPayload{}, models.TaskInstance{}, models.Task{})
		})
	}
}

func TestDelete(t *testing.T) {

	testCases := []struct {
		name                 string
		URL                  string
		method               string
		mockTaskPersistence  mock.TaskPersistenceConf
		mockTaskInstance     mock.TaskInstancePersistenceConf
		mockExecResults      mock.ExecutionResultPersistenceConf
		mockExecExpirations  mock.ExecutionExpirationPersistenceConf
		mockTriggerUC        mock.TriggerUCConf
		expectedCode         int
		expectedErrorMessage string
	}{
		{
			name:                 "testCase 0 - bad UUID in request, cannot parse",
			URL:                  "/" + partnerID + "/tasks/badUUID",
			method:               http.MethodDelete,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTaskIDHasBadFormat,
		},
		{
			name:   "testCase 1 - GetByIDs server error",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr,
			method: http.MethodDelete,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{}, errors.New("err"))
				return tp
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
		},
		{
			name:   "testCase 2 - entitlement error, no user doesnt have noc Access, but task has",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr,
			method: http.MethodDelete,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{PartnerID: partnerID, IsRequireNOCAccess: true, Targets: models.Target{Type: models.ManagedEndpoint}}}, nil)
				return tp
			},
			expectedCode:         http.StatusForbidden,
			expectedErrorMessage: errorcode.ErrorAccessDenied,
		},
		{
			name:   "testCase 3 - Cant get TaskInstance",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr,
			method: http.MethodDelete,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{PartnerID: partnerID, Targets: models.Target{Type: models.ManagedEndpoint}}}, nil)
				return tp
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByTaskID(gomock.Any(), taskID).
					Return([]models.TaskInstance{}, errors.New("err"))
				return ti
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskInstances,
		},
		{
			name:   "testCase 4 - Cant get Execution results",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr,
			method: http.MethodDelete,
			// get task by id, get taskInstance by taskID , get Exec results but get error
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{PartnerID: partnerID, Targets: models.Target{Type: models.ManagedEndpoint}}}, nil)
				return tp
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByTaskID(gomock.Any(), taskID).
					Return([]models.TaskInstance{{PartnerID: partnerID, ID: taskID}}, nil)
				return ti
			},
			mockExecResults: func(er *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				er.EXPECT().
					GetByTaskInstanceIDs(gomock.Any()).
					Return([]models.ExecutionResult{}, errors.New("err"))
				return er
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskExecutionResults,
		},
		{
			name:   "testCase 5 - Cant get exec expirations",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr,
			method: http.MethodDelete,
			// get task by id, get taskInstance by taskID , get Exec results but get error
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{PartnerID: partnerID, Targets: models.Target{Type: models.ManagedEndpoint}}}, nil)
				return tp
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByTaskID(gomock.Any(), taskID).
					Return([]models.TaskInstance{{PartnerID: partnerID, ID: taskID}}, nil)
				return ti
			},
			mockExecResults: func(er *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				er.EXPECT().
					GetByTaskInstanceIDs(gomock.Any())
				return er
			},
			mockTriggerUC: func(us *mock.MockUsecase) *mock.MockUsecase {
				us.EXPECT().GetActiveTriggersByTaskID(gomock.Any(), gomock.Any())
				return us
			},
			mockExecExpirations: func(ee *mock.MockExecutionExpirationPersistence) *mock.MockExecutionExpirationPersistence {
				ee.EXPECT().
					GetByTaskInstanceIDs(gomock.Any(), gomock.Any()).
					Return([]models.ExecutionExpiration{}, errors.New("err"))
				return ee
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantDeleteExecutionResults,
		},
		{
			name:   "testCase 5.1 - Cant get active triggers",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr,
			method: http.MethodDelete,
			// get task by id, get taskInstance by taskID , get Exec results but get error
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{PartnerID: partnerID, Targets: models.Target{Type: models.ManagedEndpoint}}}, nil)
				return tp
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByTaskID(gomock.Any(), taskID).
					Return([]models.TaskInstance{{PartnerID: partnerID, ID: taskID}}, nil)
				return ti
			},
			mockExecResults: func(er *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				er.EXPECT().
					GetByTaskInstanceIDs(gomock.Any())
				return er
			},
			mockTriggerUC: func(us *mock.MockUsecase) *mock.MockUsecase {
				us.EXPECT().GetActiveTriggersByTaskID(gomock.Any(), gomock.Any()).Return([]entities.ActiveTrigger{}, errors.New("err"))
				return us
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetActiveTriggers,
		},
		{
			name:   "testCase 6 - Can not execute deletion, all deletions returned error",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr,
			method: http.MethodDelete,
			// get task by id, get taskInstance by taskID , get Exec results but get error
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{PartnerID: partnerID, Targets: models.Target{Type: models.ManagedEndpoint}}}, nil)
				tp.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return tp
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByTaskID(gomock.Any(), taskID).
					Return([]models.TaskInstance{{PartnerID: partnerID, ID: taskID}}, nil)
				ti.EXPECT().
					DeleteBatch(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return ti
			},
			mockExecResults: func(er *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				er.EXPECT().
					GetByTaskInstanceIDs(gomock.Any()).
					Return([]models.ExecutionResult{}, nil)
				er.EXPECT().
					DeleteBatch(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return er
			},
			mockTriggerUC: func(us *mock.MockUsecase) *mock.MockUsecase {
				us.EXPECT().GetActiveTriggersByTaskID(gomock.Any(), gomock.Any())
				us.EXPECT().DeleteActiveTriggers(gomock.Any()).Return(errors.New("err"))
				return us
			},
			mockExecExpirations: func(ee *mock.MockExecutionExpirationPersistence) *mock.MockExecutionExpirationPersistence {
				ee.EXPECT().
					GetByTaskInstanceIDs(gomock.Any(), gomock.Any()).Return([]models.ExecutionExpiration{{PartnerID: partnerID}}, nil)
				ee.EXPECT().
					Delete(gomock.Any()).Return(errors.New("err"))
				return ee
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantDeleteTask,
		},
		{
			name:   "testCase 7 - Cant update counters, external task, but StatusOk",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr,
			method: http.MethodDelete,
			// get task by id, get taskInstance by taskID , get Exec results but get error
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{PartnerID: partnerID, ExternalTask: true, Targets: models.Target{Type: models.ManagedEndpoint}}}, nil)
				tp.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByTaskID(gomock.Any(), taskID).
					Return([]models.TaskInstance{{PartnerID: partnerID, ID: taskID}}, nil)
				ti.EXPECT().
					DeleteBatch(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			mockExecResults: func(er *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				er.EXPECT().
					GetByTaskInstanceIDs(gomock.Any()).
					Return([]models.ExecutionResult{}, nil)
				er.EXPECT().
					DeleteBatch(gomock.Any(), gomock.Any()).
					Return(nil)
				return er
			},
			mockTriggerUC: func(us *mock.MockUsecase) *mock.MockUsecase {
				us.EXPECT().GetActiveTriggersByTaskID(gomock.Any(), gomock.Any())
				us.EXPECT().DeleteActiveTriggers(gomock.Any())
				return us
			},
			mockExecExpirations: func(ee *mock.MockExecutionExpirationPersistence) *mock.MockExecutionExpirationPersistence {
				ee.EXPECT().
					GetByTaskInstanceIDs(gomock.Any(), gomock.Any()).Return([]models.ExecutionExpiration{{PartnerID: partnerID}}, nil)
				ee.EXPECT().Delete(gomock.Any())
				return ee
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:   "testCase 8 - Successes",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr,
			method: http.MethodDelete,
			// get task by id, get taskInstance by taskID , get Exec results but get error
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{PartnerID: partnerID, Targets: models.Target{Type: models.ManagedEndpoint}}}, nil)
				tp.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByTaskID(gomock.Any(), taskID).
					Return([]models.TaskInstance{{PartnerID: partnerID, ID: taskID}}, nil)
				ti.EXPECT().
					DeleteBatch(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			mockExecResults: func(er *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				er.EXPECT().
					GetByTaskInstanceIDs(gomock.Any()).
					Return([]models.ExecutionResult{}, nil)
				er.EXPECT().
					DeleteBatch(gomock.Any(), gomock.Any()).
					Return(nil)
				return er
			},
			mockTriggerUC: func(us *mock.MockUsecase) *mock.MockUsecase {
				us.EXPECT().GetActiveTriggersByTaskID(gomock.Any(), gomock.Any())
				us.EXPECT().DeleteActiveTriggers(gomock.Any())
				return us
			},
			mockExecExpirations: func(ee *mock.MockExecutionExpirationPersistence) *mock.MockExecutionExpirationPersistence {
				ee.EXPECT().
					GetByTaskInstanceIDs(gomock.Any(), gomock.Any()).Return([]models.ExecutionExpiration{{PartnerID: partnerID}}, nil)
				ee.EXPECT().Delete(gomock.Any())
				return ee
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:   "testCase 9 - GetByIDs returns empty slice of tasks",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr,
			method: http.MethodDelete,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{}, nil)
				return tp
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTaskIsNotFoundByTaskID,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, nil,
				nil,
				tc.mockTaskInstance,
				tc.mockTaskPersistence,
				nil,
				tc.mockExecResults,
				tc.mockExecExpirations,
				nil,
				userServiceMockNoErr,
				nil, tc.mockTriggerUC, nil, nil, nil)
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
		})
	}

}

func TestDeleteExecutions(t *testing.T) {

	testCases := []struct {
		name                 string
		URL                  string
		method               string
		mockTaskPersistence  mock.TaskPersistenceConf
		mockTaskInstance     mock.TaskInstancePersistenceConf
		mockExecResults      mock.ExecutionResultPersistenceConf
		expectedCode         int
		expectedErrorMessage string
	}{
		{
			name:                 "testCase 0 - bad taskUUID in request",
			URL:                  "/" + partnerID + "/tasks/badUUID/task-instances/" + validTaskInstanceUUIDstr,
			method:               http.MethodDelete,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTaskIDHasBadFormat,
		},
		{
			name:         "testCase 1 - bad taskInstanceUUID in request",
			URL:          "/" + partnerID + "/tasks/" + taskIDstr + "/task-instances/badUUID",
			method:       http.MethodDelete,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "testCase 2 - GetByIDs server error",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr + "/task-instances/" + validTaskInstanceUUIDstr,
			method: http.MethodDelete,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, true, taskID).
					Return([]models.Task{}, errors.New("err"))
				return tp
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
		},
		{
			name:   "testCase 3 - GetByIDs bad request error",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr + "/task-instances/" + validTaskInstanceUUIDstr,
			method: http.MethodDelete,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, true, taskID).
					Return([]models.Task{}, nil)
				return tp
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
		},
		{
			name:   "testCase 4 - accessDenied error",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr + "/task-instances/" + validTaskInstanceUUIDstr,
			method: http.MethodDelete,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, true, taskID).
					Return([]models.Task{{PartnerID: partnerID, IsRequireNOCAccess: true, Targets: models.Target{Type: models.ManagedEndpoint}}}, nil)
				return tp
			},
			expectedCode:         http.StatusForbidden,
			expectedErrorMessage: errorcode.ErrorAccessDenied,
		},
		{
			name:   "testCase 5 - Cant get TaskInstance",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr + "/task-instances/" + validTaskInstanceUUIDstr,
			method: http.MethodDelete,
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, true, taskID).
					Return([]models.Task{{PartnerID: partnerID, Targets: models.Target{Type: models.ManagedEndpoint}}}, nil)
				return tp
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByIDs(gomock.Any(), validTaskInstanceUUID).
					Return([]models.TaskInstance{}, errors.New("err"))
				return ti
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskInstances,
		},
		{
			name:   "testCase 6 - Cant get Execution results",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr + "/task-instances/" + validTaskInstanceUUIDstr,
			method: http.MethodDelete,

			// get task by id, get taskInstance by taskID , get Exec results but get error
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, true, taskID).
					Return([]models.Task{{PartnerID: partnerID, Targets: models.Target{Type: models.ManagedEndpoint}}}, nil)
				return tp
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByIDs(gomock.Any(), validTaskInstanceUUID).
					Return([]models.TaskInstance{{PartnerID: partnerID, ID: validTaskInstanceUUID}}, nil)
				return ti
			},
			mockExecResults: func(er *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				er.EXPECT().
					GetByTaskInstanceIDs(gomock.Any()).
					Return([]models.ExecutionResult{}, errors.New("err"))
				return er
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskExecutionResults,
		},
		{
			name:   "testCase 7 - Can not execute deletion, deleting of task instance returned error",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr + "/task-instances/" + validTaskInstanceUUIDstr,
			method: http.MethodDelete,
			// get task by id, get taskInstance by taskInstanceID but get error while its deleting, get Exec results
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, true, taskID).
					Return([]models.Task{{PartnerID: partnerID, Targets: models.Target{Type: models.ManagedEndpoint}}}, nil)
				return tp
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByIDs(gomock.Any(), validTaskInstanceUUID).
					Return([]models.TaskInstance{{PartnerID: partnerID, ID: validTaskInstanceUUID}}, nil)
				ti.EXPECT().
					DeleteBatch(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return ti
			},
			mockExecResults: func(er *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				er.EXPECT().
					GetByTaskInstanceIDs(gomock.Any()).
					Return([]models.ExecutionResult{}, nil)
				er.EXPECT().
					DeleteBatch(gomock.Any(), gomock.Any()).
					Return(nil)
				return er
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantDeleteTaskInstance,
		},
		{
			name:   "testCase 8: Can not execute deletion, deleting of task results returned error",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr + "/task-instances/" + validTaskInstanceUUIDstr,
			method: http.MethodDelete,

			// get task by id, get taskInstance by taskID , get Exec results but get error
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, true, taskID).
					Return([]models.Task{{PartnerID: partnerID, Targets: models.Target{Type: models.ManagedEndpoint}}}, nil)
				return tp
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByIDs(gomock.Any(), validTaskInstanceUUID).
					Return([]models.TaskInstance{{PartnerID: partnerID, ID: validTaskInstanceUUID}}, nil)
				ti.EXPECT().
					DeleteBatch(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			mockExecResults: func(er *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				er.EXPECT().
					GetByTaskInstanceIDs(gomock.Any()).
					Return([]models.ExecutionResult{}, nil)
				er.EXPECT().
					DeleteBatch(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return er
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantDeleteTaskInstance,
		},
		{
			name:   "testCase 9 - Success",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr + "/task-instances/" + validTaskInstanceUUIDstr,
			method: http.MethodDelete,

			// get task by id, get taskInstance by taskID , get Exec results and successfully deleted everything, except task
			mockTaskPersistence: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, true, taskID).
					Return([]models.Task{{PartnerID: partnerID, Targets: models.Target{Type: models.ManagedEndpoint}}}, nil)
				return tp
			},
			mockTaskInstance: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetByIDs(gomock.Any(), validTaskInstanceUUID).
					Return([]models.TaskInstance{{PartnerID: partnerID, ID: validTaskInstanceUUID}}, nil)
				ti.EXPECT().
					DeleteBatch(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			mockExecResults: func(er *mock.MockExecutionResultPersistence) *mock.MockExecutionResultPersistence {
				er.EXPECT().
					GetByTaskInstanceIDs(gomock.Any()).
					Return([]models.ExecutionResult{}, nil)
				er.EXPECT().
					DeleteBatch(gomock.Any(), gomock.Any()).
					Return(nil)
				return er
			},
			expectedCode: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := getMockedTaskService(mockController, nil,
				nil,
				tc.mockTaskInstance,
				tc.mockTaskPersistence,
				nil,
				tc.mockExecResults,
				nil,
				nil,
				userServiceMockNoErr,
				nil, nil, nil, nil, nil)
			router := getTaskServiceRouter(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.URL, nil)
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v: test case is %s", tc.expectedCode, w.Code, tc.name)
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}
		})
	}
}

func customHTTPMock(url, method string, data interface{}, expectedResponseStatus int) {
	httpmock.RegisterResponder(method, url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(expectedResponseStatus, data)
		},
	)
}

func getTaskServiceRouter(service TaskService) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/{partnerID}/tasks", service.Create).Methods(http.MethodPost)
	router.HandleFunc("/{partnerID}/tasks/{taskID}", service.Edit).Methods(http.MethodPost)
	router.HandleFunc("/{partnerID}/tasks/managed-endpoints/{managedEndpointID}", service.GetByPartnerAndManagedEndpointID).Methods(http.MethodGet)
	router.HandleFunc("/{partnerID}/tasks/data", service.GetTasksSummaryData).Methods(http.MethodGet)
	router.HandleFunc("/{partnerID}/tasks/{taskID}", service.GetByPartnerAndID).Methods(http.MethodGet)
	router.HandleFunc("/{partnerID}/tasks/{taskID}/enable", service.EnableTaskForAllTargets).Methods(http.MethodPut)
	router.HandleFunc("/{partnerID}/tasks/{taskID}/enable/targets", service.EnableTaskForSelectedTargets).Methods(http.MethodPut)
	router.HandleFunc("/{partnerID}/tasks/task-instances/stop", service.StopTaskInstanceExecutions).Methods(http.MethodPut)
	router.HandleFunc("/{partnerID}/tasks/{taskID}", service.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/{partnerID}/tasks/{taskID}/task-instances/{taskInstanceID}", service.DeleteExecutions).Methods(http.MethodDelete)
	router.HandleFunc("/{partnerID}/tasks/{taskID}/task-instances/count", service.TaskInstancesCountByTaskID).Methods(http.MethodGet)
	router.HandleFunc("/{partnerID}/tasks/{taskID}/postpone", service.PostponeNearestExecution).Methods(http.MethodPut)
	router.HandleFunc("/{partnerID}/tasks/{taskID}/managed-endpoints/{managedEndpointID}/postpone", service.PostponeDeviceNearestExecution).Methods(http.MethodPut)
	router.HandleFunc("/{partnerID}/tasks/{taskID}/cancel", service.CancelNearestExecution).Methods(http.MethodPut)
	router.HandleFunc("/{partnerID}/tasks/{taskID}/managed-endpoint/{managedEndpointID}/cancel", service.CancelNearestExecutionForEndpoint).Methods(http.MethodPut)
	router.HandleFunc("/{partnerID}/tasks/data/{originID}", service.GetByOriginID).Methods(http.MethodGet)
	return
}
