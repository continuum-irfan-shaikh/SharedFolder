package tasks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const (
	timeToPostponeDurationStr30min   = "30"
	timeToPostponeDurationStr60min   = "60"
	timeToPostponeDurationStr180min  = "180"
	timeToPostponeDurationStr1440min = "1440"
)

var tenDaysToPostpone = time.Duration(time.Hour * 24 * 10)

func TestPostponeNearestExecution(t *testing.T) {
	var (
		timeNow                           = time.Now().UTC().Truncate(time.Minute)
		runTimeWeekly, _                  = time.Parse(time.RFC3339, "2090-10-09T11:35:00Z")
		postponedRunTimeWeeklyExpected, _ = time.Parse(time.RFC3339, "2090-10-10T11:35:00Z")
		targetsStrSlice                   = []string{validUUIDstr}
		statusesMap                       = make(map[gocql.UUID]statuses.TaskInstanceStatus)
		pendingStatusesMap                = make(map[gocql.UUID]statuses.TaskInstanceStatus)

		targetsManagedEndpoint = models.Target{
			IDs:  targetsStrSlice,
			Type: models.ManagedEndpoint,
		}
		taskInstances = []models.TaskInstance{{
			ID:            validTaskInstanceUUID,
			Statuses:      statusesMap,
			OverallStatus: statuses.TaskInstanceScheduled},
		}

		taskInstancesWithPending = []models.TaskInstance{
			{
				ID:            validTaskInstanceUUID,
				Statuses:      statusesMap,
				OverallStatus: statuses.TaskInstanceScheduled},
			{
				Statuses: pendingStatusesMap,
			},
		}

		scheduleOneTime = apiModels.Schedule{
			Location:   "UTC",
			Regularity: apiModels.OneTime,
		}
		scheduleRecurrentHourly = apiModels.Schedule{
			Location:   "UTC",
			Regularity: apiModels.Recurrent,
			Repeat: apiModels.Repeat{
				Frequency: apiModels.Hourly,
				Every:     1,
			},
			EndRunTime: timeNow.Add(dayDuration * 2),
		}
		scheduleRecurrentDaily = apiModels.Schedule{
			Location:   "UTC",
			Regularity: apiModels.Recurrent,
			Repeat: apiModels.Repeat{
				Frequency: apiModels.Daily,
				Every:     1,
				RunTime:   timeNow.Truncate(time.Minute),
			},
			EndRunTime:   timeNow.Add(dayDuration * 4),
			StartRunTime: timeNow,
		}
		scheduleRecurrentWeekly = apiModels.Schedule{
			Location:   "UTC",
			Regularity: apiModels.Recurrent,
			Repeat: apiModels.Repeat{
				Period:     41,
				Frequency:  apiModels.Weekly,
				Every:      7,
				DaysOfWeek: []int{1, 2},
				RunTime:    runTimeWeekly.Truncate(time.Minute),
			},
			StartRunTime: runTimeWeekly,
		}

		scheduleRecurrentNoFreq = apiModels.Schedule{
			Location:   "UTC",
			Regularity: apiModels.Recurrent,
			Repeat:     apiModels.Repeat{},
		}

		scheduleRunNow = apiModels.Schedule{
			Location:   "UTC",
			Regularity: apiModels.RunNow,
		}

		//---------------------------------------------------------------------------------------------------------------------------------
		//	----- postponing hourly for 30 mins
		recurrentTaskThirtyMinPostponeHourly = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentHourly,
				ManagedEndpointID:  validUUID,
			},
		}
		expectedRecurrentTaskThirtyMinPostponeHourly = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PostponedRunTime:   timeNow.Add(thirtyMinutesDuration),
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentHourly,
				ManagedEndpointID:  validUUID,
			},
		}

		// ----- postponing hourly for hour
		recurrentTaskHourPostponeHourly = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentHourly,
				ManagedEndpointID:  validUUID,
			},
		}
		expectedRecurrentTaskHourPostponeHourly = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PostponedRunTime:   timeNow.Add(hourDuration),
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentHourly,
				ManagedEndpointID:  validUUID,
			},
		}
		// ----- postponing hourly for 30 mins, but it was already postponed before
		recurrentTaskThirtyMinPostponeHourlyPostponed = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PostponedRunTime:   timeNow.Add(thirtyMinutesDuration),
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentHourly,
				ManagedEndpointID:  validUUID,
			},
		}
		expectedRecurrentTaskThirtyMinPostponeHourlyPostponed = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PostponedRunTime:   timeNow.Add(hourDuration),
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentHourly,
				ManagedEndpointID:  validUUID,
			},
		}

		// ----- postponing hourly for 3 hours
		recurrentTaskThreeHoursPostponeHourly = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentHourly,
				ManagedEndpointID:  validUUID,
			},
		}
		expectedRecurrentThreeHoursPostponeHourly = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PostponedRunTime:   timeNow.Add(threeHoursDuration),
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentHourly,
				ManagedEndpointID:  validUUID,
			},
		}
		// ----- postponing Daily for 180 min
		recurrentTask3HoursPostponeDaily = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentDaily,
				ManagedEndpointID:  validUUID,
			},
		}
		expectedRecurrent3HoursPostponeDaily = []models.Task{
			{
				Name:             "task1",
				RunTimeUTC:       timeNow,
				PostponedRunTime: timeNow.Add(threeHoursDuration),
				PartnerID:        partnerID, ID: validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentDaily,
				ManagedEndpointID:  validUUID,
			},
		}
		// ----- postponing already postponed Daily for 180 min
		recurrentTask3HoursPostponeDailyPostponed = []models.Task{
			{
				Name:                "task1",
				RunTimeUTC:          timeNow,
				OriginalNextRunTime: timeNow.Add(thirtyMinutesDuration),
				PartnerID:           partnerID,
				ID:                  validUUID,
				ExternalTask:        true,
				IsRequireNOCAccess:  false,
				Targets:             targetsManagedEndpoint,
				Schedule:            scheduleRecurrentDaily,
				ManagedEndpointID:   validUUID,
			},
		}

		expectedRecurrent3HoursPostponeDailyPostponed = []models.Task{
			{
				Name:             "task1",
				RunTimeUTC:       timeNow.Add(thirtyMinutesDuration),
				PostponedRunTime: timeNow.Add(threeHoursDuration),
				PartnerID:        partnerID, ID: validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentDaily,
				ManagedEndpointID:  validUUID,
			},
		}
		// ----- postponing Daily for 24hours
		recurrentTask24HoursPostponeDaily = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentDaily,
				ManagedEndpointID:  validUUID,
			},
		}
		expectedRecurrent24HoursPostponeDaily = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PostponedRunTime:   timeNow.Add(dayDuration),
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentDaily,
				ManagedEndpointID:  validUUID,
			},
		}
		// ----- postponing Weekly for 24hours
		recurrentTask24HoursPostponeWeekly = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         runTimeWeekly,
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentWeekly,
				ManagedEndpointID:  validUUID,
			},
		}
		expectedRecurrent24HoursPostponeWeekly = []models.Task{
			{
				Name:             "task1",
				RunTimeUTC:       runTimeWeekly,
				PostponedRunTime: postponedRunTimeWeeklyExpected,
				PartnerID:        partnerID,
				ID:               validUUID,
				ExternalTask:     true, IsRequireNOCAccess: false,
				Targets:           targetsManagedEndpoint,
				Schedule:          scheduleRecurrentWeekly,
				ManagedEndpointID: validUUID,
			},
		}

		//-----days postpone
		recurrentTaskDatePostponeWeekly = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PartnerID:          partnerID,
				ID:                 validUUID,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentWeekly,
				ManagedEndpointID:  validUUID,
			},
		}
		expectedRecurrentDatePostponeWeekly = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PostponedRunTime:   timeNow.Add(tenDaysToPostpone),
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentWeekly,
				ManagedEndpointID:  validUUID,
			},
		}

		//-----date postpone has original
		recurrentTaskDatePostponeWeeklyOriginal = []models.Task{
			{
				Name:                "task1",
				RunTimeUTC:          timeNow,
				OriginalNextRunTime: timeNow.Add(thirtyMinutesDuration),
				PartnerID:           partnerID,
				ID:                  validUUID,
				IsRequireNOCAccess:  false,
				Targets:             targetsManagedEndpoint,
				Schedule:            scheduleRecurrentWeekly,
				ManagedEndpointID:   validUUID,
			},
		}
		expectedRecurrentDatePostponeWeeklyOriginal = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow.Add(thirtyMinutesDuration),
				PostponedRunTime:   timeNow.Add(tenDaysToPostpone),
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentWeekly,
				ManagedEndpointID:  validUUID,
			},
		}
		//---------------------------------------------------------------------------------------------------------------------------------

		//-----postponed + 180 days
		recurrentTaskDatePostponeWeekly180 = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow.Add(time.Hour * 24 * 1),
				PostponedRunTime:   timeNow.Add(time.Hour * 24 * 180),
				PartnerID:          partnerID,
				ID:                 validUUID,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentWeekly,
				ManagedEndpointID:  validUUID,
			},
		}
		expectedRecurrentDatePostponeWeekly180 = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow.Add(time.Hour * 24 * 1),
				PostponedRunTime:   timeNow.Add(time.Hour * 24 * 181),
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentWeekly,
				ManagedEndpointID:  validUUID,
			},
		}
		////---------------------------------------------------------------------------------------------------------------------------------
		//-----postponed + 180 days
		oneTimeTaskDatePostponeWeekly180 = []models.Task{
			{
				Name:                "task1",
				RunTimeUTC:          timeNow.Add(time.Hour * 24 * 180),
				OriginalNextRunTime: timeNow.Add(time.Hour * 24 * 1),
				PartnerID:           partnerID,
				ID:                  validUUID,
				IsRequireNOCAccess:  false,
				Targets:             targetsManagedEndpoint,
				Schedule:            scheduleOneTime,
				ManagedEndpointID:   validUUID,
			},
		}
		expectedOneTimeDatePostponeWeekly180 = []models.Task{
			{
				Name:                "task1",
				RunTimeUTC:          timeNow.Add(time.Hour * 24 * 181),
				OriginalNextRunTime: timeNow.Add(time.Hour * 24 * 1),
				PartnerID:           partnerID,
				ID:                  validUUID,
				ExternalTask:        true,
				IsRequireNOCAccess:  false,
				Targets:             targetsManagedEndpoint,
				Schedule:            scheduleOneTime,
				ManagedEndpointID:   validUUID,
			},
		}
		////-----postponed + 180 days
		oneTimeTask = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow.Add(time.Hour * 24 * 10),
				PartnerID:          partnerID,
				ID:                 validUUID,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleOneTime,
				ManagedEndpointID:  validUUID,
			},
		}
		expectedOneTime = []models.Task{
			{
				Name:                "task1",
				RunTimeUTC:          timeNow.Add(time.Hour * 24 * 20),
				OriginalNextRunTime: timeNow.Add(time.Hour * 24 * 10),
				PartnerID:           partnerID,
				ID:                  validUUID,
				ExternalTask:        true,
				IsRequireNOCAccess:  false,
				Targets:             targetsManagedEndpoint,
				Schedule:            scheduleOneTime,
				ManagedEndpointID:   validUUID,
			},
		}
		//----
		//---------------------------------------------------------------------------------------------------------------------------------
		recurrentTasksNoFrequency = []models.Task{
			{
				Name:               "task1",
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentNoFreq,
				ManagedEndpointID:  validUUID,
			},
		}

		tasksRunNowRegularity = []models.Task{
			{
				Name:               "task1",
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRunNow,
				ManagedEndpointID:  validUUID,
			},
			{
				Name:               "task2",
				PartnerID:          partnerID,
				ID:                 validUUID2,
				IsRequireNOCAccess: true,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRunNow,
				ManagedEndpointID:  validUUID,
			},
		}

		tasksOneTimeRegularity = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleOneTime,
				ManagedEndpointID:  validUUID,
			},
			{
				Name:               "task2",
				RunTimeUTC:         timeNow,
				PartnerID:          partnerID,
				ID:                 validUUID2,
				IsRequireNOCAccess: true,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleOneTime,
				ManagedEndpointID:  validUUID,
			},
		}
		inactiveTasks = []models.Task{
			{
				Name:               "task1",
				State:              statuses.TaskStateInactive,
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleOneTime,
			},
			{
				Name:               "task2",
				State:              statuses.TaskStateInactive,
				PartnerID:          partnerID,
				ID:                 validUUID2,
				IsRequireNOCAccess: true,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleOneTime,
			},
		}
	)

	testCases := []struct {
		name                   string
		URL                    string
		isNeedBadBody          bool
		postponeBody           postponeBody
		taskInstanceMock       mock.TaskInstancePersistenceConf
		taskPersistenceMock    mock.TaskPersistenceConf
		userMock               user.Service
		method                 string
		expectedErrorMessage   string
		expectedCode           int
		expectedPostponedTasks []models.Task
	}{
		{
			name: "testCase 0 - bad taskID",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			URL:                  "/" + partnerID + "/tasks/badTaskID/postpone",
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorTaskIDHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name: "testCase 1 - empty taskID",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			URL:                  "/" + partnerID + "/tasks/00000000-0000-0000-0000-000000000000/postpone",
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorTaskIDHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name:                 "testCase 2 - bad input body",
			isNeedBadBody:        true,
			URL:                  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorTimeFrameHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name: "testCase 3 - can't get Tasks",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			URL:    "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{}, errors.New("db err"))
				return tp
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
		},
		{
			name: "testCase 3.5 - negative duration",
			postponeBody: postponeBody{
				DurationString: "-30",
			},
			URL:                  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			method:               http.MethodPut,
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTimeFrameHasBadFormat,
		},
		{
			name: "testCase 4 - data base returned zero internal tasks",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{}, nil)
				return tp
			},
			expectedCode:         http.StatusNotFound,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
		},
		{
			name: "testCase 4.5 - wrong NOC access",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{IsRequireNOCAccess: true}}, nil)
				return tp
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
		},
		{
			name: "testCase 4.6 - 2 iterations, can't postpone",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return([]models.Task{{IsRequireNOCAccess: true}}, nil)
				return tp
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
		},
		{
			name: "testCase 5 - inactive tasks, cant postpone",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(inactiveTasks, nil)
				return tp
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
		},
		{
			name: "testCase 6 - RunNow regularity at schedule, cant postpone runNow",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(tasksRunNowRegularity, nil)
				return tp
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
		},
		{
			name: "testCase 7 - RunNow regularity at schedule, cant postpone runNow",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(tasksRunNowRegularity, nil)
				return tp
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
		},
		{
			name: "testCase 8 - Postponing OneTime tasks with expected postpone time, but got Insert DB error",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(tasksOneTimeRegularity, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return([]models.TaskInstance{{Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{validUUID: statuses.TaskInstanceScheduled}}}, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantSaveTaskToDB,
		},
		{
			name: "testCase 8.5 - Postponing OneTime tasks with expected postpone time, but got pending instance",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(tasksOneTimeRegularity, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return([]models.TaskInstance{{Statuses: pendingStatusesMap}}, nil)
				return ti
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
		},
		{
			name: "testCase 9 - Postponing OneTime tasks with expected postpone time, but got Insert TI DB error",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(tasksOneTimeRegularity, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return([]models.TaskInstance{{Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{validUUID: statuses.TaskInstanceScheduled}}}, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return ti
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantUpdateTaskInstances,
		},
		{
			name: "testCase 10 - Postponing OneTime tasks with expected postpone time,get task instance err",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(tasksOneTimeRegularity, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return([]models.TaskInstance{{Statuses: map[gocql.UUID]statuses.TaskInstanceStatus{validUUID: statuses.TaskInstanceScheduled}}}, errors.New("err"))
				return ti
			},
			expectedCode:         http.StatusInternalServerError,
			expectedErrorMessage: errorcode.ErrorCantGetTaskInstances,
		},
		{
			name: "testCase 10.1 - Postponing Reccurent, no frequency in schedule",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(recurrentTasksNoFrequency, nil)
				return tp
			},

			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
		},
		{
			name: "testCase 10.2 - Postponing Reccurent hourly task, 30 min, status ok",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(recurrentTaskThirtyMinPostponeHourly, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()). //
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode:           http.StatusOK,
			expectedPostponedTasks: expectedRecurrentTaskThirtyMinPostponeHourly,
		},
		{
			name: "testCase 10.3 - Postponing Reccurent hourly task, but got pending instance",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(recurrentTaskThirtyMinPostponeHourly, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstancesWithPending, nil)
				return ti
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
		},
		{
			name: "testCase 10.5 - Postponing Reccurent hourly task, 60 min, status ok",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr60min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(recurrentTaskHourPostponeHourly, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode:           http.StatusOK,
			expectedPostponedTasks: expectedRecurrentTaskHourPostponeHourly,
		},
		{
			name: "testCase 11 - Postponing Reccurent hourly task, 180 min, status ok",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr180min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(recurrentTaskThreeHoursPostponeHourly, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode:           http.StatusOK,
			expectedPostponedTasks: expectedRecurrentThreeHoursPostponeHourly,
		},
		{
			name: "testCase 12 - Postponing Reccurent daily task, 180 min, status ok",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr180min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(recurrentTask3HoursPostponeDaily, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode:           http.StatusOK,
			expectedPostponedTasks: expectedRecurrent3HoursPostponeDaily,
		},
		{
			name: "testCase 13 - Postponing Reccurent daily task, 1440 min, status ok",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(recurrentTask24HoursPostponeDaily, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode:           http.StatusOK,
			expectedPostponedTasks: expectedRecurrent24HoursPostponeDaily,
		},
		{
			name: "testCase 14 - Postponing Reccurent weekly task, 1440 min, status ok",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(recurrentTask24HoursPostponeWeekly, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode:           http.StatusOK,
			expectedPostponedTasks: expectedRecurrent24HoursPostponeWeekly,
		},
		{
			name: "testCase 15 - Postponing Reccurent Hourly task, 30 mins, task.OriginalNextRunTime == postponedTime ,status ok",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(recurrentTaskThirtyMinPostponeHourlyPostponed, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode:           http.StatusOK,
			expectedPostponedTasks: expectedRecurrentTaskThirtyMinPostponeHourlyPostponed,
		},
		{
			name: "testCase 16 - Postponing Reccurent Hourly task, 180 mins, task.OriginalNextRunTime == postponedTime ,status ok",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: timeToPostponeDurationStr180min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(recurrentTask3HoursPostponeDailyPostponed, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()). //
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode:           http.StatusOK,
			expectedPostponedTasks: expectedRecurrent3HoursPostponeDailyPostponed,
		},
		{
			name: "testCase 17 - Postponing Reccurent Hourly task, date",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: "14400", // 10 days in mins
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(recurrentTaskDatePostponeWeekly, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()). //
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode:           http.StatusOK,
			expectedPostponedTasks: expectedRecurrentDatePostponeWeekly,
		},
		{
			name: "testCase 18 - Postponing Reccurent Hourly task, date, has original Time",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: "14400", // 10 days in mins
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(recurrentTaskDatePostponeWeeklyOriginal, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode:           http.StatusOK,
			expectedPostponedTasks: expectedRecurrentDatePostponeWeeklyOriginal,
		},
		{
			name: "testCase 19 - Postponing Reccurent Hourly task, postponed, 180 days override",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: "14400", // 10 days in mins
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(recurrentTaskDatePostponeWeekly180, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode:           http.StatusOK,
			expectedPostponedTasks: expectedRecurrentDatePostponeWeekly180,
		},
		{
			name: "testCase 20 - Postponing OneTime Hourly task, postponed, 180 days override",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: "259200",
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(oneTimeTaskDatePostponeWeekly180, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode:           http.StatusOK,
			expectedPostponedTasks: expectedOneTimeDatePostponeWeekly180,
		},
		{
			name: "testCase 20.5 - Postponing OneTime Hourly task",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			postponeBody: postponeBody{
				DurationString: "14400", // 10 days in mins
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDs(gomock.Any(), nil, partnerID, false, taskID).
					Return(oneTimeTask, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedCode:           http.StatusOK,
			expectedPostponedTasks: expectedOneTime,
		},
		{
			name:   "testCase 21 - more than 180 days",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr + "/postpone",
			method: http.MethodPut,
			postponeBody: postponeBody{
				DurationString: "1438560", //999 days
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTimeFrameHasBadFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			statusesMap[validUUID] = statuses.TaskInstanceScheduled
			pendingStatusesMap[validUUID] = statuses.TaskInstancePending

			s := getMockedTaskService(mockController, nil, nil, tc.taskInstanceMock, tc.taskPersistenceMock, nil, nil,
				nil, nil, tc.userMock, nil, nil, nil, nil, nil)
			router := getTaskServiceRouter(s)

			body, err := json.Marshal(&tc.postponeBody)
			if err != nil {
				t.Fatalf("cannot parse body %v", err)
			}

			w := httptest.NewRecorder()
			var r *http.Request
			if !tc.isNeedBadBody {
				r = httptest.NewRequest(tc.method, tc.URL, bytes.NewReader(body))
			} else {
				r = httptest.NewRequest(tc.method, tc.URL, bytes.NewReader([]byte("bad")))
			}

			router.ServeHTTP(w, r)
			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}

			if len(tc.expectedPostponedTasks) > 0 {
				var gotTasks models.Task
				err = json.Unmarshal(w.Body.Bytes(), &gotTasks)
				if err != nil {
					t.Errorf("Error while unmarshalling body, %v", err)
				}

				if gotTasks.PostponedRunTime != tc.expectedPostponedTasks[0].PostponedRunTime {
					t.Errorf("Want PostponedRunTime %v but bot %v", tc.expectedPostponedTasks[0].PostponedRunTime, gotTasks.PostponedRunTime)
				}

				if gotTasks.RunTimeUTC != tc.expectedPostponedTasks[0].RunTimeUTC {
					t.Errorf("Want RunTimeUTC %v but bot %v", tc.expectedPostponedTasks[0].RunTimeUTC, gotTasks.RunTimeUTC)
				}

				if gotTasks.OriginalNextRunTime != tc.expectedPostponedTasks[0].OriginalNextRunTime {
					t.Errorf("Want OriginalNextRunTime %v but bot %v", tc.expectedPostponedTasks[0].OriginalNextRunTime, gotTasks.OriginalNextRunTime)
				}
			}
		})
	}
}

func TestPostponeDeviceNearestExecution(t *testing.T) {
	var (
		timeNow                 = time.Now().UTC().Truncate(time.Minute)
		targetsStrSlice         = []string{validUUIDstr}
		targetsManagedEndpoint  = models.Target{IDs: targetsStrSlice, Type: models.ManagedEndpoint}
		scheduledStatusesMap    = make(map[gocql.UUID]statuses.TaskInstanceStatus)
		pendingStatusesMap      = make(map[gocql.UUID]statuses.TaskInstanceStatus)
		scheduleRecurrentHourly = apiModels.Schedule{Location: "UTC", Regularity: apiModels.Recurrent, Repeat: apiModels.Repeat{Frequency: apiModels.Hourly, Every: 1}, EndRunTime: timeNow.Add(dayDuration * 2)}

		recurrentTaskThirtyMinPostponeHourly = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentHourly,
				ManagedEndpointID:  validUUID,
			},
		}
		expectedRecurrentTaskThirtyMinPostponeHourly = models.Task{
			Name:               "task1",
			RunTimeUTC:         timeNow,
			PostponedRunTime:   timeNow.Add(thirtyMinutesDuration),
			PartnerID:          partnerID,
			ID:                 validUUID,
			ExternalTask:       true,
			IsRequireNOCAccess: false,
			Targets:            targetsManagedEndpoint,
			Schedule:           scheduleRecurrentHourly,
			ManagedEndpointID:  validUUID,
		}

		//-----date postpone
		recurrentTaskDatePostponeWeekly = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow,
				PartnerID:          partnerID,
				ID:                 validUUID,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentHourly,
				ManagedEndpointID:  validUUID,
			},
		}
		expectedRecurrentTaskDatePostpone = models.Task{
			Name:               "task1",
			RunTimeUTC:         timeNow,
			PostponedRunTime:   timeNow.Add(tenDaysToPostpone),
			PartnerID:          partnerID,
			ID:                 validUUID,
			ExternalTask:       true,
			IsRequireNOCAccess: false,
			Targets:            targetsManagedEndpoint,
			Schedule:           scheduleRecurrentHourly,
			ManagedEndpointID:  validUUID,
		}

		runNowTask = []models.Task{
			{
				Name:               "task1",
				RunTimeUTC:         timeNow.Truncate(time.Minute),
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           apiModels.Schedule{Regularity: apiModels.RunNow},
				ManagedEndpointID:  validUUID,
			},
		}
		inactiveTasks = []models.Task{
			{
				Name:               "task1",
				State:              statuses.TaskStateInactive,
				RunTimeUTC:         timeNow.Truncate(time.Minute),
				PartnerID:          partnerID,
				ID:                 validUUID,
				ExternalTask:       true,
				IsRequireNOCAccess: false,
				Targets:            targetsManagedEndpoint,
				Schedule:           scheduleRecurrentHourly,
				ManagedEndpointID:  validUUID,
			},
		}
		taskInstances = []models.TaskInstance{{
			ID:            validTaskInstanceUUID,
			Statuses:      scheduledStatusesMap,
			OverallStatus: statuses.TaskInstanceScheduled},
		}
		taskInstancesWithPending = []models.TaskInstance{
			{
				ID:            validTaskInstanceUUID,
				Statuses:      scheduledStatusesMap,
				OverallStatus: statuses.TaskInstanceScheduled},
			{
				Statuses: pendingStatusesMap,
			},
		}
	)

	testCases := []struct {
		name                 string
		URL                  string
		isNeedBadBody        bool
		body                 postponeBody
		taskInstanceMock     mock.TaskInstancePersistenceConf
		taskPersistenceMock  mock.TaskPersistenceConf
		userMock             user.Service
		method               string
		expectedErrorMessage string
		expectedCode         int
		expectedTaskDetails  interface{}
	}{
		{
			name: "testCase 1 - bad taskID",
			URL:  "/" + partnerID + "/tasks/badTaskID/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorTaskIDHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name: "testCase 2 - bad managedEndpointID",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/badUUID/postpone",
			body: postponeBody{
				DurationString: "15",
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorEndpointIDHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name:                 "testCase 3 - bad input body",
			URL:                  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			isNeedBadBody:        true,
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorTimeFrameHasBadFormat,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name: "testCase 4 - GetByIDAndManagedEndpoints db error",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
			expectedCode:         http.StatusInternalServerError,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return([]models.Task{}, errors.New("db error"))
				return tp
			},
		},
		{
			name: "testCase 4.6 - 2 iterations, can't postpone",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method: http.MethodPut,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return([]models.Task{{IsRequireNOCAccess: true}}, nil)
				return tp
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
		},
		{
			name: "testCase 5 - GetByIDAndManagedEndpoints returned 0 tasks",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantGetTaskByTaskID,
			expectedCode:         http.StatusNotFound,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return([]models.Task{}, nil)
				return tp
			},
		},
		{
			name: "testCase 6 - wrong NOC access in task",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
			expectedCode:         http.StatusBadRequest,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return([]models.Task{{IsRequireNOCAccess: true}}, nil)
				return tp
			},
		},
		{
			name: "testCase 7 - can't get task instance",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantGetTaskInstances,
			expectedCode:         http.StatusInternalServerError,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return(recurrentTaskThirtyMinPostponeHourly, nil)
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
			name: "testCase 8 - got zero task instances",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantGetTaskInstances,
			expectedCode:         http.StatusNotFound,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return(recurrentTaskThirtyMinPostponeHourly, nil)
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
			name: "testCase 8.5 - got taskInstance with pending",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
			expectedCode:         http.StatusBadRequest,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return(recurrentTaskThirtyMinPostponeHourly, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return([]models.TaskInstance{{Statuses: pendingStatusesMap}}, nil)
				return ti
			},
		},
		{
			name: "testCase 8.6 - got previous taskInstance with pending",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
			expectedCode:         http.StatusBadRequest,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return(recurrentTaskThirtyMinPostponeHourly, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstancesWithPending, nil)
				return ti
			},
		},
		{
			name: "testCase 9 - task is inactive, cant update",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
			expectedCode:         http.StatusBadRequest,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return(inactiveTasks, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				return ti
			},
		},
		{
			name: "testCase 9 - task has no regularity, cant update",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
			expectedCode:         http.StatusBadRequest,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return(inactiveTasks, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				return ti
			},
		},
		{
			name: "testCase 9.5 - task has is runNow, cant update",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantUpdateTask,
			expectedCode:         http.StatusBadRequest,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return(runNowTask, nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				return ti
			},
		},
		{
			name: "testCase 10 - Insert task db error",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantSaveTaskToDB,
			expectedCode:         http.StatusInternalServerError,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return(recurrentTaskThirtyMinPostponeHourly, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				return ti
			},
		},
		{
			name: "testCase 11 - Insert task db error",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantSaveTaskToDB,
			expectedCode:         http.StatusInternalServerError,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return(recurrentTaskThirtyMinPostponeHourly, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				return ti
			},
		},
		{
			name: "testCase 12 - Insert task instances db error",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr1440min,
			},
			method:               http.MethodPut,
			expectedErrorMessage: errorcode.ErrorCantUpdateTaskInstances,
			expectedCode:         http.StatusInternalServerError,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return(recurrentTaskThirtyMinPostponeHourly, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				return ti
			},
		},
		{
			name: "testCase 13 - Ok",
			URL:  "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			body: postponeBody{
				DurationString: timeToPostponeDurationStr30min,
			},
			method:       http.MethodPut,
			expectedCode: http.StatusOK,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return(recurrentTaskThirtyMinPostponeHourly, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedTaskDetails: expectedRecurrentTaskThirtyMinPostponeHourly,
		},
		{
			name:   "testCase 14 - Ok",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			method: http.MethodPut,
			body: postponeBody{
				DurationString: "14400", // 10 days in mins
			},
			expectedCode: http.StatusOK,
			taskPersistenceMock: func(tp *mock.MockTaskPersistence) *mock.MockTaskPersistence {
				tp.EXPECT().
					GetByIDAndManagedEndpoints(gomock.Any(), partnerID, taskID, validUUID).
					Return(recurrentTaskDatePostponeWeekly, nil)
				tp.EXPECT().
					InsertOrUpdate(gomock.Any(), gomock.Any()).
					Return(nil)
				return tp
			},
			taskInstanceMock: func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
				ti.EXPECT().
					GetTopInstancesByTaskID(gomock.Any(), taskID).
					Return(taskInstances, nil)
				ti.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(nil)
				return ti
			},
			expectedTaskDetails: expectedRecurrentTaskDatePostpone,
		},
		{
			name:   "testCase 15 - more than 180 days",
			URL:    "/" + partnerID + "/tasks/" + taskIDstr + "/managed-endpoints/" + validUUIDstr + "/postpone",
			method: http.MethodPut,
			body: postponeBody{
				DurationString: time.Duration(time.Hour * 24 * 999).String(),
			},
			expectedCode:         http.StatusBadRequest,
			expectedErrorMessage: errorcode.ErrorTimeFrameHasBadFormat,
		},
	}

	for _, tc := range testCases {
		scheduledStatusesMap[validUUID] = statuses.TaskInstanceScheduled
		pendingStatusesMap[validUUID] = statuses.TaskInstancePending

		t.Run(tc.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			s := getMockedTaskService(mockController, nil, nil, tc.taskInstanceMock, tc.taskPersistenceMock, nil, nil,
				nil, nil, tc.userMock, nil, nil, nil, nil, nil)
			router := getTaskServiceRouter(s)

			body, err := json.Marshal(&tc.body)
			if err != nil {
				t.Fatalf("cannot parse body %v", err)
			}

			w := httptest.NewRecorder()
			var r *http.Request
			if !tc.isNeedBadBody {
				r = httptest.NewRequest(tc.method, tc.URL, bytes.NewReader(body))
			} else {
				r = httptest.NewRequest(tc.method, tc.URL, bytes.NewReader([]byte("")))
			}

			router.ServeHTTP(w, r)
			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}

			if !mock.IsErrorMessageValid(tc.expectedErrorMessage, w.Body.Bytes()) && tc.expectedErrorMessage != "" {
				t.Errorf("Wanted error message constant \"%v\" but got message %v",
					mock.TranslateErrorMessage(tc.expectedErrorMessage), w.Body.String())
			}

			if tc.expectedTaskDetails != nil {
				var gotTasksDetails models.Task
				expectedTaskDetails := tc.expectedTaskDetails.(models.Task)

				err = json.Unmarshal(w.Body.Bytes(), &gotTasksDetails)
				if err != nil {
					t.Errorf("Error while unmarshalling body, %v", err)
				}

				if gotTasksDetails.PostponedRunTime != expectedTaskDetails.PostponedRunTime {
					t.Errorf("Want PostponedRunTime %v but bot %v", expectedTaskDetails.PostponedRunTime, gotTasksDetails.PostponedRunTime)
				}

				if gotTasksDetails.RunTimeUTC != expectedTaskDetails.RunTimeUTC {
					t.Errorf("Want RunTimeUTC %v but bot %v", expectedTaskDetails.RunTimeUTC, gotTasksDetails.RunTimeUTC)
				}
			}
		})
	}
}
