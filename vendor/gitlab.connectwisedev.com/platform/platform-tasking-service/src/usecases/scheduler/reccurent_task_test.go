package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

func TestIsFirstRunning(t *testing.T) {
	logger.Load(config.Config.Log)
	pr := New(config.Config, logger.Log, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	nowUTC := time.Now().UTC().Truncate(time.Minute)
	createdAt1 := nowUTC.Add(time.Minute - 1)
	runTime1 := nowUTC
	locPlus2, err := time.LoadLocation("Europe/Kiev")
	if err != nil {
		t.Fatal(err)
	}

	schedule1 := tasking.Schedule{
		Regularity:   tasking.Recurrent,
		StartRunTime: nowUTC.In(locPlus2),
		EndRunTime:   nowUTC.Add(time.Hour * 24).In(locPlus2),
		Repeat: tasking.Repeat{
			Frequency: tasking.Hourly,
			Every:     3,
		},
	}
	schedule2 := tasking.Schedule{
		Regularity:   tasking.Recurrent,
		StartRunTime: nowUTC.In(locPlus2),
		EndRunTime:   nowUTC.Add(time.Hour * 24).In(locPlus2),
		Repeat: tasking.Repeat{
			Every: 3,
		},
	}

	schedule3 := tasking.Schedule{
		Regularity:   tasking.Recurrent,
		StartRunTime: nowUTC.In(locPlus2).Add(time.Hour * -1),
		EndRunTime:   nowUTC.Add(time.Hour * 24).In(locPlus2),
		Repeat: tasking.Repeat{
			Every: 3,
		},
	}
	testCases := []struct {
		name     string
		task     models.Task
		instance models.TaskInstance
		expected bool
	}{
		{
			name:     "true, createdAt in the past. first running for +3 timezone machine, startRun TIme in loc == runTime UTC",
			expected: true,
			task: models.Task{
				LastTaskInstanceID: gocql.TimeUUID(),
				RunTimeUTC:         runTime1,
				CreatedAt:          createdAt1,
				Schedule:           schedule1,
			},
		},
		{
			name:     "true, modifiedBy in the past. first running for +3 timezone machine, startRun TIme in loc == runTime UTC",
			expected: true,
			task: models.Task{
				LastTaskInstanceID: gocql.TimeUUID(),
				RunTimeUTC:         runTime1,
				ModifiedAt:         createdAt1,
				Schedule:           schedule1,
			},
		},
		{
			name:     "true, createdAt in the past but invalid schedule. first running for +3 timezone machine, startRun TIme in loc == runTime UTC",
			expected: true,
			task: models.Task{
				LastTaskInstanceID: gocql.TimeUUID(),
				RunTimeUTC:         runTime1,
				CreatedAt:          createdAt1,
				Schedule:           schedule2,
			},
		},
		{
			name:     "false, startRunTime befoore run time utc",
			expected: false,
			task: models.Task{
				LastTaskInstanceID: gocql.TimeUUID(),
				RunTimeUTC:         runTime1,
				CreatedAt:          createdAt1,
				Schedule:           schedule3,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := pr.IsFirstRunning(context.TODO(), tc.task, tc.instance, tc.task.Schedule.StartRunTime.Location())
			if got != tc.expected {
				t.Errorf("expected %v but got %v", tc.expected, got)
			}
		})
	}

}

func TestRecalculateTasks(t *testing.T) {
	logger.Load(config.Config.Log)
	pr := New(config.Config, logger.Log, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	nowUTC := time.Now().UTC().Truncate(time.Minute)

	schedule := tasking.Schedule{
		Regularity:   tasking.Recurrent,
		StartRunTime: nowUTC,
		Location:     "Incorrect",
		Repeat: tasking.Repeat{
			Frequency: tasking.Frequency(20), // incorrect one
		},
	}
	schedule1 := tasking.Schedule{
		Regularity:   tasking.Recurrent,
		StartRunTime: nowUTC,
		Location:     "Incorrect",
		Repeat: tasking.Repeat{
			Frequency: tasking.Hourly,
			Every:     1,
		},
	}
	schedule2 := tasking.Schedule{
		Regularity:   tasking.Recurrent,
		StartRunTime: nowUTC,
		Location:     "Incorrect",
		EndRunTime:   nowUTC.Add(time.Minute * 10),
		Repeat: tasking.Repeat{
			Frequency: tasking.Hourly,
			Every:     1,
		},
	}
	testCases := []struct {
		name            string
		task            []models.Task
		expectedRunTime time.Time
		instance        models.TaskInstance
		expectedErr     bool
	}{
		{
			name:            "error happened",
			expectedErr:     true,
			expectedRunTime: nowUTC.Add(time.Hour),
			task: []models.Task{
				{
					OriginalNextRunTime: nowUTC.Add(time.Hour),
				},
				{
					RunTimeUTC: nowUTC,
					Schedule:   schedule,
				},
			},
		},
		{
			name:            "good postponed",
			expectedRunTime: nowUTC.Add(time.Minute * 30),
			task: []models.Task{
				{
					PostponedRunTime: nowUTC.Add(time.Minute * 30),
					RunTimeUTC:       nowUTC,
					Schedule:         schedule1,
				},
			},
		},
		{
			name:            "good end run time end",
			expectedRunTime: nowUTC.Add(time.Minute * 10),
			task: []models.Task{
				{
					RunTimeUTC: nowUTC,
					Schedule:   schedule2,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotTasks, err := pr.recalculateTasks(context.TODO(), tc.task)
			if err == nil && tc.expectedErr {
				t.Errorf("expected %v but got %v", tc.expectedErr, err)
			}

			if !gotTasks[0].RunTimeUTC.Equal(tc.expectedRunTime) {
				t.Errorf("expected %v but got %v", tc.expectedRunTime, gotTasks[0].RunTimeUTC)
			}
		})
	}
}

func TestGetUnexpectedRunTime(t *testing.T) {
	nowUTC := time.Now().UTC()
	t.Run("hourly", func(t *testing.T) {
		task := models.Task{
			RunTimeUTC: nowUTC,
			Schedule: tasking.Schedule{
				Repeat: tasking.Repeat{
					Frequency: tasking.Hourly,
					Every:     1,
				},
			},
		}

		expectedRunTime := nowUTC.Add(time.Hour)
		got := getUnexpectedRunTime(task.RunTimeUTC, task.Schedule)
		if !got.Equal(expectedRunTime) {
			t.Errorf("expected %v but got %v", expectedRunTime, got)
		}
	})

	t.Run("day", func(t *testing.T) {
		task := models.Task{
			RunTimeUTC: nowUTC,
			Schedule: tasking.Schedule{
				Repeat: tasking.Repeat{
					Frequency: tasking.Daily,
					Every:     1,
				},
			},
		}

		expectedRunTime := nowUTC.Add(time.Hour*24)
		got := getUnexpectedRunTime(task.RunTimeUTC, task.Schedule)
		if !got.Equal(expectedRunTime) {
			t.Errorf("expected %v but got %v", expectedRunTime, got)
		}
	})

	t.Run("week", func(t *testing.T) {
		task := models.Task{
			RunTimeUTC: nowUTC,
			Schedule: tasking.Schedule{
				Repeat: tasking.Repeat{
					Frequency: tasking.Weekly,
					Every:     1,
				},
			},
		}

		expectedRunTime := nowUTC.Add(time.Hour*24*7)
		got := getUnexpectedRunTime(task.RunTimeUTC, task.Schedule)
		if !got.Equal(expectedRunTime) {
			t.Errorf("expected %v but got %v", expectedRunTime, got)
		}
	})
	t.Run("month", func(t *testing.T) {
		task := models.Task{
			RunTimeUTC: nowUTC,
			Schedule: tasking.Schedule{
				Repeat: tasking.Repeat{
					Frequency: tasking.Monthly,
					Every:     1,
				},
			},
		}

		expectedRunTime := nowUTC.Add(time.Hour * 24 * 7 * 30)
		got := getUnexpectedRunTime(task.RunTimeUTC, task.Schedule)
		if !got.Equal(expectedRunTime) {
			t.Errorf("expected %v but got %v", expectedRunTime, got)
		}
	})
}
