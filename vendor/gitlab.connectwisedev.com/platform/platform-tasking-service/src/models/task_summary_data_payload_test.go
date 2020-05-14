package models

import (
	"reflect"
	"testing"

	"github.com/gocql/gocql"
	. "github.com/onsi/gomega"
)

var taskInstID, _ = gocql.ParseUUID("18f66dfc-9be2-48eb-9451-1e80cfc967b5")

func TestGroupTaskByID(t *testing.T) {
	testTasksInternal := []Task{
		{
			ID: gocql.TimeUUID(),
		},
		{
			ID: gocql.TimeUUID(),
		},
		{
			ID: gocql.TimeUUID(),
		},
	}

	testGroup := groupTasksByID(testTasksInternal)

	for _, want := range testTasksInternal {
		got := testGroup[want.ID]

		if !reflect.DeepEqual(got, []Task{want}) {
			t.Fatalf("Got %v, want %v", got, want)
		}
	}
}

func TestGetLastTaskInstanceIDs(t *testing.T) {
	testTasksInternal := []Task{
		{
			LastTaskInstanceID: gocql.TimeUUID(),
		},
		{
			LastTaskInstanceID: gocql.TimeUUID(),
		},
		{
			LastTaskInstanceID: gocql.TimeUUID(),
		},
	}

	testGroup := getLastTaskInstanceIDs(testTasksInternal)

	for _, task := range testTasksInternal {
		found := false

		for _, got := range testGroup {
			if task.LastTaskInstanceID == got {
				found = true
			}
		}

		if !found {
			t.Fatalf("LastTaskInstanceID %s not found", task.LastTaskInstanceID)
		}
	}
}

func TestGroupTaskInstancesByID(t *testing.T) {
	testTaskInstances := []TaskInstance{
		{
			ID: gocql.TimeUUID(),
		},
		{
			ID: gocql.TimeUUID(),
		},
		{
			ID: gocql.TimeUUID(),
		},
	}

	testGroup := groupTaskInstancesByID(testTaskInstances)
	for i, v := range testTaskInstances {
		got := testGroup[v.ID].ID
		want := testTaskInstances[i].ID
		if got != want {
			t.Fatalf("Got %s, want %s", got.String(), want.String())
		}
	}
}

func TestGroupLastTaskInstancesByTaskID(t *testing.T) {
	var (
		taskID              = gocql.TimeUUID()
		lastTaskInstanceID1 = gocql.TimeUUID()
		lastTaskInstanceID2 = gocql.TimeUUID()
		taskInstance1       = TaskInstance{
			ID: lastTaskInstanceID1,
		}
		taskInstance2 = TaskInstance{
			ID: lastTaskInstanceID2,
		}
		taskInstancesByID = map[gocql.UUID]TaskInstance{
			lastTaskInstanceID1: taskInstance1,
			lastTaskInstanceID2: taskInstance2,
		}
		testCases = []struct {
			name                          string
			tasksInternal                 []Task
			taskInstancesByID             map[gocql.UUID]TaskInstance
			expectedTaskInstancesByTaskID map[gocql.UUID][]TaskInstance
		}{
			{
				name: "different LastTaskInstanceID for each internal task",
				tasksInternal: []Task{
					{
						ID:                 taskID,
						LastTaskInstanceID: lastTaskInstanceID1,
					},
					{
						ID:                 taskID,
						LastTaskInstanceID: lastTaskInstanceID2,
					},
				},
				taskInstancesByID: taskInstancesByID,
				expectedTaskInstancesByTaskID: map[gocql.UUID][]TaskInstance{
					taskID: {taskInstance1, taskInstance2},
				},
			},
			{
				name: "single LastTaskInstanceID for each internal task",
				tasksInternal: []Task{
					{
						ID:                 taskID,
						LastTaskInstanceID: lastTaskInstanceID1,
					},
					{
						ID:                 taskID,
						LastTaskInstanceID: lastTaskInstanceID1,
					},
				},
				taskInstancesByID: taskInstancesByID,
				expectedTaskInstancesByTaskID: map[gocql.UUID][]TaskInstance{
					taskID: {taskInstance1},
				},
			},
			{
				name: "internal task with no LastTaskInstanceID",
				tasksInternal: []Task{
					{
						ID: taskID,
					},
				},
				taskInstancesByID:             taskInstancesByID,
				expectedTaskInstancesByTaskID: map[gocql.UUID][]TaskInstance{},
			},
		}
	)
	for _, testCase := range testCases {
		test := testCase
		t.Run(test.name, func(t *testing.T) {
			actualTaskInstancesByTaskID := groupLastTaskInstancesByTaskID(test.tasksInternal, test.taskInstancesByID)
			if !reflect.DeepEqual(actualTaskInstancesByTaskID, test.expectedTaskInstancesByTaskID) {
				t.Fatalf(`Expected %v, but got %v`, test.expectedTaskInstancesByTaskID, actualTaskInstancesByTaskID)
			}
		})
	}

}

func TestGroupTaskInstancesByTaskID(t *testing.T) {
	RegisterTestingT(t)

	type payload struct {
		taskInstances []TaskInstance
	}

	type expected struct {
		taskInstancesMap map[gocql.UUID][]TaskInstance
	}

	tc := []struct {
		payload
		expected
	}{
		{
			payload: payload{
				taskInstances: []TaskInstance{{TaskID: taskInstID}},
			},
			expected: expected{
				taskInstancesMap: map[gocql.UUID][]TaskInstance{
					taskInstID: {{TaskID: taskInstID}},
				},
			},
		},
	}

	for _, test := range tc {
		m := GroupTaskInstancesByTaskID(test.payload.taskInstances)
		Ω(m[taskInstID]).To(ConsistOf(test.expected.taskInstancesMap[taskInstID]))
	}
}
