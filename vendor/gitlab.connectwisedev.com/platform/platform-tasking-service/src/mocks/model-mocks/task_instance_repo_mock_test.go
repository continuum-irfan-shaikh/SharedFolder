package modelMocks

import (
	"context"
	"reflect"
	"testing"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

var (
	taskInstancesRepo TaskInstanceRepoMock
)

func taskInstanceRepoMockStartup() {
	taskInstancesRepo = NewTaskInstanceRepoMock(true)
}

func TestCreate(t *testing.T) {
	taskInstanceRepoMockStartup()
	taskInstance := models.TaskInstance{
		ID:     str2uuid("01000000-0000-0000-0000-000000000001"),
		TaskID: str2uuid("01000000-0000-0000-0000-000000000002"),
	}
	err := taskInstancesRepo.Insert(setContext(t, false), taskInstance)
	if err != nil {
		t.Errorf("Error for creating task instance: %v", err)
	}
}

func TestCreate_Incorrect(t *testing.T) {
	taskInstanceRepoMockStartup()
	taskInstance := models.TaskInstance{
		ID:     str2uuid("01000000-0000-0000-0000-000000000001"),
		TaskID: emptyUUID,
	}
	err := taskInstancesRepo.Insert(setContext(t, false), taskInstance)
	if err == nil {
		t.Errorf("Not created task instance, error: %v", err)
	}
}

func TestCreate_WithError(t *testing.T) {
	repo := NewTaskInstanceRepoMock(true)
	taskInstance := models.TaskInstance{
		ID:     str2uuid("01000000-0000-0000-0000-000000000001"),
		TaskID: str2uuid("01000000-0000-0000-0000-000000000002"),
	}
	err := repo.Insert(setContext(t, true), taskInstance)
	if err == nil {
		t.Fatalf("Expected error, but got <nil>")
	}
}

func TestCreateEmptyNotEqualID(t *testing.T) {
	taskInstanceRepoMockStartup()
	err := taskInstancesRepo.Insert(setContext(t, false), models.TaskInstance{
		TaskID: str2uuid("11111111-2222-1111-1111-111111111111"),
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGetByTaskID(t *testing.T) {
	tests := []struct {
		name     string
		taskID   gocql.UUID
		needErr  bool
		needFill bool
		want     []models.TaskInstance
	}{
		{
			name:     "Good",
			taskID:   str2uuid("00000000-0000-0000-0000-000000000000"),
			needErr:  false,
			needFill: true,
			want:     []models.TaskInstance{DefaultTaskInstances[0]},
		},
		{
			name:     "Bad - cassandra",
			taskID:   str2uuid("00000000-0000-0000-0000-000000000000"),
			needErr:  true,
			needFill: true,
			want:     []models.TaskInstance{DefaultTaskInstances[0]},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := NewTaskInstanceRepoMock(test.needFill)
			ctx := context.WithValue(context.TODO(), IsNeedError, test.needErr)
			got, err := mock.GetByTaskID(ctx, test.taskID)
			if (err != nil) != test.needErr {
				t.Fatalf("Got err=%v, want err=%v", err, test.needErr)
			}
			if !test.needErr {
				if test.needFill && !reflect.DeepEqual(got, test.want) {
					t.Fatalf("Got %v, want %v", got, test.want)
				}
				if !test.needFill && !reflect.DeepEqual(got, []models.TaskInstance{}) {
					t.Fatalf("Wanted nil slice of Task Instances, but got %v", got)
				}
			}
		})
	}
}
