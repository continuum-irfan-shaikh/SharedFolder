package managedEndpoints

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	modelMocks "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/model-mocks"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var (
	partnerID             = "1"
	userID                = "2"
	timeUUID              = gocql.TimeUUID()
	managedEndpointTarget = models.TargetsByType{models.ManagedEndpoint: []string{timeUUID.String()}}
	dynamicGroupTarget    = models.TargetsByType{models.DynamicGroup: []string{timeUUID.String()}}
	siteTarget            = models.TargetsByType{models.Site: []string{timeUUID.String()}}
	uuid, _               = gocql.ParseUUID("ab5fa564-5a19-43dd-bacf-e3a26d7f5c8a")
)

func TestGetManagedEndpointsFromTargets(t *testing.T) {
	testCases := []struct {
		name                     string
		target                   models.TargetsByType
		httpMockURL              string
		httpMockResponseStatus   int
		httpMockResponseBody     []byte
		httpMockResponseError    error
		userSitesMock            mock.UserSitesConf
		expectedError            bool
		expectedManagedEndpoints map[gocql.UUID]models.TargetType
	}{
		{
			name:                     "testCase 1 - ManagedEndpoints target",
			target:                   managedEndpointTarget,
			expectedManagedEndpoints: nil,
		},
		{
			name:                     "testCase 2 - Invalid target",
			target:                   models.TargetsByType{gocql.BatchSizeMaximum: []string{}},
			expectedManagedEndpoints: map[gocql.UUID]models.TargetType{},
			expectedError:            true,
		},
		{
			name:   "testCase 3 - Site target, httpMock response err",
			target: siteTarget,
			httpMockURL: fmt.Sprintf("%s/partner/%s/sites/%s/summary",
				config.Config.AssetMsURL, partnerID, strings.Join(siteTarget[models.Site], siteDelimiter)),
			expectedError:            true,
			httpMockResponseError:    errors.New("err"),
			expectedManagedEndpoints: map[gocql.UUID]models.TargetType{},
		},
		{
			name:   "testCase 4 - DynamicGroup target, err getting userSites from db",
			target: dynamicGroupTarget,
			httpMockURL: fmt.Sprintf("%s/partners/%s/dynamic-groups/managed-endpoints/set?expression=%s",
				config.Config.DynamicGroupsMsURL, partnerID, strings.Join(dynamicGroupTarget[models.DynamicGroup], dgDelimiter)),
			httpMockResponseStatus: http.StatusOK,
			httpMockResponseBody:   []byte("[{\"id\":\"ab5fa564-5a19-43dd-bacf-e3a26d7f5c8a\"}]"),
			expectedError:          true,
			userSitesMock: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					Sites(gomock.Any(), partnerID, userID).
					Return(entities.UserSites{}, errors.New("err"))
				return us
			},
		},
		{
			name:   "testCase 5 - DynamicGroup target, bad sideID, error",
			target: dynamicGroupTarget,
			httpMockURL: fmt.Sprintf("%s/partners/%s/dynamic-groups/managed-endpoints/set?expression=%s",
				config.Config.DynamicGroupsMsURL, partnerID, strings.Join(dynamicGroupTarget[models.DynamicGroup], dgDelimiter)),
			httpMockResponseStatus: http.StatusOK,
			httpMockResponseBody:   []byte("[{\"id\":\"ab5fa564-5a19-43dd-bacf-e3a26d7f5c8a\"},{\"id\":\"ab5fa562-5a19-43dd-bacf-e3a26d7f5c8a\"}]"),
			expectedError:          true,
			userSitesMock: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					Sites(gomock.Any(), partnerID, userID)
				return us
			},
		},
		{
			name:   "testCase 6 - DynamicGroup target, good",
			target: dynamicGroupTarget,
			httpMockURL: fmt.Sprintf("%s/partners/%s/dynamic-groups/managed-endpoints/set?expression=%s",
				config.Config.DynamicGroupsMsURL, partnerID, strings.Join(dynamicGroupTarget[models.DynamicGroup], dgDelimiter)),
			httpMockResponseStatus: http.StatusOK,
			httpMockResponseBody:   []byte("[{\"id\":\"ab5fa564-5a19-43dd-bacf-e3a26d7f5c8a\", \"site\":\"1\"}]"),
			userSitesMock: func(us *mock.MockUserSitesPersistence) *mock.MockUserSitesPersistence {
				us.EXPECT().
					Sites(gomock.Any(), partnerID, userID).
					Return(entities.UserSites{SiteIDs: []int64{1}}, nil)
				return us
			},
			expectedManagedEndpoints: map[gocql.UUID]models.TargetType{uuid: models.DynamicGroup},
		},
		{
			name:   "testCase 7 - Site target, successes",
			target: siteTarget,
			httpMockURL: fmt.Sprintf("%s/partner/%s/sites/%s/summary",
				config.Config.AssetMsURL, partnerID, strings.Join(siteTarget[models.Site], siteDelimiter)),
			httpMockResponseStatus:   http.StatusOK,
			httpMockResponseBody:     []byte("[{\"id\":\"ab5fa564-5a19-43dd-bacf-e3a26d7f5c8a\", \"endpointID\":\"ab5fa564-5a19-43dd-bacf-e3a26d7f5c8a\"}]"),
			expectedManagedEndpoints: map[gocql.UUID]models.TargetType{uuid: models.Site},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.httpMockURL) > 0 {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder(http.MethodGet, tc.httpMockURL,
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewBytesResponse(tc.httpMockResponseStatus, tc.httpMockResponseBody), tc.httpMockResponseError
					},
				)
			}

			if tc.userSitesMock != nil {
				mockController := gomock.NewController(t)
				defer mockController.Finish()

				userSitesMock := mock.NewMockUserSitesPersistence(mockController)
				userSitesMock = tc.userSitesMock(userSitesMock)
				models.UserSitesPersistenceInstance = userSitesMock
			}

			task := models.Task{
				PartnerID:     partnerID,
				CreatedBy:     userID,
				TargetsByType: tc.target,
			}

			gotManagedEndpointsMap, err := GetManagedEndpointsFromTargets(context.Background(), task, http.DefaultClient)
			if (err != nil) != tc.expectedError {
				t.Errorf("Got err:%v but expected err was %v", err, tc.expectedError)
			}
			if !reflect.DeepEqual(gotManagedEndpointsMap, tc.expectedManagedEndpoints) && tc.expectedManagedEndpoints != nil {
				t.Errorf("Want %v but gotManagedEndpointsMap %v", tc.expectedManagedEndpoints, gotManagedEndpointsMap)
			}
		})
	}
	fmt.Println(dgDelimiter, siteDelimiter)
}

func TestGetDifference(t *testing.T) {
	var (
		testTimeNow                 = time.Now().UTC()
		deprecatedManagedEndpointID = gocql.TimeUUID()
		oldValidManagedEndpointID   = gocql.TimeUUID()
		newManagedEndpointID        = gocql.TimeUUID()
		managedEndpointID1          = gocql.TimeUUID()
		managedEndpointID2          = gocql.TimeUUID()
		targetsManagedEndpoint      = models.Target{
			IDs:  []string{managedEndpointID1.String(), managedEndpointID2.String()},
			Type: models.ManagedEndpoint,
		}
		deprecatedTask = models.Task{
			ID:                gocql.TimeUUID(),
			Name:              "Name",
			Description:       "Description",
			CreatedAt:         testTimeNow.AddDate(0, 0, -2),
			CreatedBy:         "Admin",
			PartnerID:         "PartnerID",
			OriginID:          gocql.TimeUUID(),
			State:             statuses.TaskStateActive,
			Schedule:          apiModels.Schedule{Regularity: apiModels.Recurrent, StartRunTime: testTimeNow.AddDate(0, 0, -1), EndRunTime: testTimeNow.AddDate(0, 0, 4), Location: "Europe/Kiev", Repeat: apiModels.Repeat{Every: 1, Frequency: apiModels.Hourly}},
			RunTimeUTC:        testTimeNow,
			Type:              "script",
			ManagedEndpointID: deprecatedManagedEndpointID,
		}
		oldValidTask                    = *deprecatedTask.CopyWithRunTime(oldValidManagedEndpointID)
		newTask                         = *deprecatedTask.CopyWithRunTime(newManagedEndpointID)
		taskWithManagedEndpointTargets1 = *copyTaskAndSetTargets(deprecatedTask, managedEndpointID1, targetsManagedEndpoint)
		taskWithManagedEndpointTargets2 = *taskWithManagedEndpointTargets1.CopyWithRunTime(managedEndpointID2)
		tests                           = []struct {
			name                    string
			tasksToRunGroup         []models.Task
			disabledTasksGroup      []models.Task
			newManagedEndpoints     map[gocql.UUID]struct{}
			expectedNewTasks        []models.Task
			expectedDeprecatedTasks []models.Task
			expectedOldActiveTasks  []models.Task
			allManagedEndpoints     map[gocql.UUID]struct{}
		}{
			{
				name:            "one deleted, one added endpoints",
				tasksToRunGroup: []models.Task{deprecatedTask, oldValidTask},
				newManagedEndpoints: map[gocql.UUID]struct{}{
					oldValidManagedEndpointID: {},
					newManagedEndpointID:      {},
				},
				expectedNewTasks:        []models.Task{newTask},
				expectedDeprecatedTasks: []models.Task{*deactivateTask(deprecatedTask)},
				expectedOldActiveTasks:  []models.Task{oldValidTask},
				allManagedEndpoints: map[gocql.UUID]struct{}{
					oldValidManagedEndpointID: {},
				},
			},
			{
				name:            "targets are managed endpoints",
				tasksToRunGroup: []models.Task{taskWithManagedEndpointTargets1, taskWithManagedEndpointTargets2},
				newManagedEndpoints: map[gocql.UUID]struct{}{
					managedEndpointID1: {},
					managedEndpointID2: {}},
				expectedOldActiveTasks: []models.Task{taskWithManagedEndpointTargets1, taskWithManagedEndpointTargets2},
				allManagedEndpoints: map[gocql.UUID]struct{}{
					managedEndpointID1: {},
					managedEndpointID2: {}},
			},
			{
				name:               "one deleted, one added endpoints, one disabled",
				tasksToRunGroup:    []models.Task{deprecatedTask},
				disabledTasksGroup: []models.Task{oldValidTask},
				newManagedEndpoints: map[gocql.UUID]struct{}{
					oldValidManagedEndpointID: {},
					newManagedEndpointID:      {},
				},
				expectedNewTasks:        []models.Task{newTask},
				expectedDeprecatedTasks: []models.Task{*deactivateTask(deprecatedTask)},
				expectedOldActiveTasks:  nil,
				allManagedEndpoints: map[gocql.UUID]struct{}{
					oldValidManagedEndpointID: {},
				},
			},
		}
	)
	models.TaskInstancePersistenceInstance = modelMocks.NewTaskInstanceRepoMock(false)
	logger.Load(config.Config.Log)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newTasks, deprecatedTasks, oldActiveTasks := GetDifference(context.Background(), tt.tasksToRunGroup, tt.disabledTasksGroup, tt.newManagedEndpoints, tt.allManagedEndpoints)
			if !reflect.DeepEqual(newTasks, tt.expectedNewTasks) {
				t.Fatalf("GetDifference: test case '%v'. Expected newTasks = %v, but got \n%v", tt.name, tt.expectedNewTasks, newTasks)
				return
			}
			if !reflect.DeepEqual(deprecatedTasks, tt.expectedDeprecatedTasks) {
				t.Fatalf("GetDifference: test case '%v'. Expected deprecatedTasks = %v, but got \n%v", tt.name, tt.expectedDeprecatedTasks, deprecatedTasks)
				return
			}
			if !reflect.DeepEqual(oldActiveTasks, tt.expectedOldActiveTasks) {
				t.Fatalf("GetDifference: test case '%v'. Expected oldActiveTasks = %v, but got %v", tt.name, tt.expectedOldActiveTasks, oldActiveTasks)
				return
			}
		})
	}

	getCounterForTask(newTask)
	newTask.ExternalTask = true
	getCounterForTask(newTask)
	insertNewDeviceIntoInstance(gocql.TimeUUID(), gocql.TimeUUID(), gocql.TimeUUID())

}

func deactivateTask(task models.Task) *models.Task {
	task.State = statuses.TaskStateInactive
	return &task
}

func copyTaskAndSetTargets(task models.Task, managedEndpoint gocql.UUID, targets models.Target) *models.Task {
	task.ManagedEndpointID = managedEndpoint
	task.Targets = targets
	return &task
}
