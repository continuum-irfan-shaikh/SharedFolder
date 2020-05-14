package tasks

import (
	"testing"
	"time"

	"github.com/gocql/gocql"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-repository"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

func TestNewInternalTasks(t *testing.T) {
	var (
		taskID       = gocql.TimeUUID()
		originID     = gocql.TimeUUID()
		endpointID0  = gocql.TimeUUID()
		endpointID1  = gocql.TimeUUID()
		testTemplate = models.TemplateDetails{
			OriginID: originID,
		}
		now = time.Now().UTC().Truncate(time.Minute)

		testTask = models.Task{
			ID:       taskID,
			OriginID: originID,
			Name:     "Test task",
			Targets: models.Target{
				IDs: []string{
					endpointID0.String(),
					endpointID1.String(),
				},
			},
			Schedule: apiModels.Schedule{
				Regularity:   apiModels.OneTime,
				StartRunTime: now.Add(time.Minute),
			},
		}
		meWithTimeZones = map[gocql.UUID]*time.Location{
			endpointID0: time.UTC,
			endpointID1: time.UTC,
		}
		expectedTasks = []models.Task{
			{
				ID:       taskID,
				OriginID: originID,
				Name:     "Test task",
				Targets: models.Target{
					IDs: []string{
						endpointID0.String(),
						endpointID1.String(),
					},
				},
				ManagedEndpointID: endpointID0,
				RunTimeUTC:        now.Add(time.Minute),
			},

			{
				ID:       taskID,
				OriginID: originID,
				Name:     "Test task",
				Targets: models.Target{
					IDs: []string{
						endpointID0.String(),
						endpointID1.String(),
					},
				},
				ManagedEndpointID: endpointID1,
				RunTimeUTC:        now.Add(time.Minute),
			},
		}
		expectedRunTime = now.Add(time.Minute)
	)

	tasksService := New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &mockrepositories.EncryptionServiceMock{}, &mockrepositories.AgentEncryptionServiceMock{})

	gotInternalTask, err := tasksService.buildInternalTasks(&testTask, testTemplate, meWithTimeZones, map[gocql.UUID]models.TargetType{endpointID0: 0, endpointID1: 0})
	if err != nil {
		t.Fatal("unexpected err: ", err)
	}

	if len(gotInternalTask) != len(expectedTasks) {
		t.Fatal("got wrong number of tasks: ", len(gotInternalTask))
	}

	if gotInternalTask[0].Name != testTask.Name {
		t.Errorf("Expected task name [%s] but got [%s]\n", testTask.Name, gotInternalTask[0].Name)
	}

	if gotInternalTask[0].Schedule.Regularity != testTask.Schedule.Regularity {
		t.Errorf("tasks have wrong schedule [%+v]", gotInternalTask[0].Schedule)
	}

	if gotInternalTask[0].RunTimeUTC != expectedRunTime {
		t.Fatalf("Expected time for run [%+v] but got [%+v]\n", expectedRunTime, gotInternalTask[0].RunTimeUTC)
	}

	// bad case - wrong schedule
	testTask.Schedule = apiModels.Schedule{}
	_, err = tasksService.buildInternalTasks(&testTask, testTemplate, meWithTimeZones, map[gocql.UUID]models.TargetType{endpointID0: 0, endpointID1: 0})
	if err == nil {
		t.Fatal("unexpected error but got <nil> ")
	}
}

func TestParseTaskTimeInLoc(t *testing.T) {
	now := time.Now().UTC()
	kievLoc, _ := time.LoadLocation("EET")

	testCases := []struct {
		name         string
		schedule     apiModels.Schedule
		loc          *time.Location
		expectedTime time.Time
		expectedErr  bool
	}{
		{
			name: "bad - oneTime UTC+1",
			schedule: apiModels.Schedule{
				StartRunTime: now.Add(-30 * time.Hour),
				Regularity:   apiModels.Recurrent,
				Location:     "+0100",
				EndRunTime:   now.Add(time.Hour * 24),
			},
			loc:         time.FixedZone("UTC+1", 60*60),
			expectedErr: true,
		},
		{
			name: "good - oneTime UTC+1",
			schedule: apiModels.Schedule{
				StartRunTime: now.Add(time.Hour),
				Regularity:   apiModels.OneTime,
				Location:     "+0100",
				EndRunTime:   now.Add(time.Hour * 24),
			},
			loc:          time.FixedZone("UTC+1", 60*60),
			expectedTime: now,
		},
		{
			name: "good - recurrent UTC+2",
			schedule: apiModels.Schedule{
				StartRunTime: now.Add(time.Hour * 2).Add(time.Minute),
				Regularity:   apiModels.Recurrent,
				Location:     "Europe/Kiev",
				Repeat: apiModels.Repeat{
					Every:     1,
					Frequency: apiModels.Hourly,
				},
			},
			loc:          kievLoc,
			expectedTime: now.Add(time.Minute),
		},
		{
			name: "good - runNow UTC+1",
			schedule: apiModels.Schedule{
				Regularity: apiModels.RunNow,
				Location:   "+0100",
				EndRunTime: now.Add(time.Hour * 24),
			},
			loc:          time.FixedZone("UTC+1", 60*60),
			expectedTime: now,
		},
	}

	for _, tc := range testCases {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			gotRunTime, gotSchedule, err := parseTaskTimeInLoc(test.schedule, test.loc)
			if tc.expectedErr {
				if err == nil {
					t.Fatalf("expected err but got nil")
				}
				return
			} else if err != nil {
				t.Errorf("%s: unexpected err: %v", test.name, err)
			}

			if gotSchedule.Regularity != apiModels.RunNow {
				if gotRunTime != test.expectedTime {
					t.Errorf("%s: unexpected time [%s] but got [%s]", test.name, test.expectedTime.String(), gotRunTime.String())
				}

				if gotSchedule.StartRunTime.Location().String() != test.loc.String() {
					t.Errorf("%s: expected location [%v] but got [%v]", test.name, test.loc, gotSchedule.StartRunTime.Location())
				}
			} else {
				if gotSchedule.EndRunTime.Location().String() != test.loc.String() {
					t.Errorf("%s: expected location [%v] but got [%v]", test.name, test.loc, gotSchedule.EndRunTime.Location())
				}
			}
		})
	}
}

func TestTaskName(t *testing.T) {
	testCases := []struct {
		name         string
		taskName     string
		scriptName   string
		expectedName string
	}{
		{
			name:         "Task name is NOT empty",
			taskName:     "test task",
			scriptName:   "name from template",
			expectedName: "test task",
		},
		{
			name:         "Task name is empty",
			taskName:     "",
			scriptName:   "name from template",
			expectedName: "name from template",
		},
	}

	for _, tc := range testCases {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			gotName := getTaskName(test.taskName, test.scriptName)
			if gotName != test.expectedName {
				t.Fatalf("%s: expected name [%s] but got [%s]", test.name, test.expectedName, gotName)
			}
		})
	}

}

func TestGetPeriod(t *testing.T) {
	testCases := []struct {
		name           string
		frequency      apiModels.Frequency
		startTime      time.Time
		expectedPeriod int
	}{
		{
			name:           "hourly",
			frequency:      apiModels.Hourly,
			expectedPeriod: 0,
		},
		{
			name:      "daily",
			frequency: apiModels.Daily,
			// year: 2100, month: 2, day: 5
			startTime:      time.Date(2100, 2, 5, 0, 0, 0, 0, time.UTC),
			expectedPeriod: 36,
		},
		{
			name:      "weekly",
			frequency: apiModels.Weekly,
			// year: 2100, month: 1, day: 7
			startTime:      time.Date(2100, 1, 7, 0, 0, 0, 0, time.UTC),
			expectedPeriod: 1,
		},
		{
			name:      "monthly",
			frequency: apiModels.Monthly,
			// year: 2100, month: 10, day: 1
			startTime:      time.Date(2100, 10, 1, 0, 0, 0, 0, time.UTC),
			expectedPeriod: 10,
		},
		{
			name:           "wrong frequency",
			expectedPeriod: 0,
		},
	}

	for _, tc := range testCases {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			gotPeriod := getPeriod(test.frequency, test.startTime, time.UTC)
			if gotPeriod != test.expectedPeriod {
				t.Log(test.startTime)
				t.Fatalf("%s: extected period [%d] but got [%d]", test.name, test.expectedPeriod, gotPeriod)
			}
		})
	}
}
