package modelMocks

import (
	"context"
	"reflect"
	"testing"

	"encoding/binary"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"github.com/gocql/gocql"
	"sort"
	"time"
)

type byTargetID []models.TaskCount

func (s byTargetID) Len() int {
	return len(s)
}

func (s byTargetID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byTargetID) Less(i, j int) bool {
	return binary.BigEndian.Uint64(s[i].ManagedEndpointID.Bytes()) < binary.BigEndian.Uint64(s[j].ManagedEndpointID.Bytes())
}

func TestGetByIDs_Success(t *testing.T) {
	repo := NewTaskRepoMock(true)
	taskID := ExistedTaskID
	partnerID := PartnerID

	if newTasks, err := repo.GetByIDs(getContextWithTransactionID(t), nil, partnerID, false, taskID); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(newTasks[0], DefaultTasks[0]) {
		t.Fatal("Invalid task returned")
	}
}

func TestGetByIDs_Fail(t *testing.T) {
	repo := NewTaskRepoMock(true)
	taskID := ExistedTaskID
	partnerID := PartnerID

	if _, err := repo.GetByIDs(getContextWithNeedErr(t, true), nil, partnerID, false, taskID); err == nil {
		t.Fatalf("Expected error, but no error returned")
	}
}

func TestGetByIDs_Negative(t *testing.T) {
	repo := NewTaskRepoMock(true)
	if _, err := repo.GetByIDs(setContext(t, true), nil, somePartnerID, false, ExistedTaskID); err == nil {
		t.Fatal(err)
	}
}

func TestInsertOrUpdate_Fail(t *testing.T) {
	var repo TaskRepoMock
	var task models.Task

	if err := repo.InsertOrUpdate(getContextWithTransactionID(t), task); err == nil {
		t.Fatalf("1. Expected 'bad input data' error, but error is nil")
	}

	repo = NewTaskRepoMock(false)
	if err := repo.InsertOrUpdate(getContextWithTransactionID(t), task); err == nil {
		t.Fatalf("2. Expected 'bad input data' error, but no error returned")
	}

	task.ID = gocql.TimeUUID()
	if err := repo.InsertOrUpdate(getContextWithTransactionID(t), task); err == nil {
		t.Fatalf("3. Expected 'bad input data' error, but error is nil")
	}

	ctx := context.WithValue(getContextWithTransactionID(t), IsNeedError, true)
	if err := repo.InsertOrUpdate(ctx, task); err == nil {
		t.Fatalf("1. Expected 'Cassandra is down' error, but error is nil")
	}
}

func TestInsertOrUpdate_Success(t *testing.T) {
	var task models.Task
	repo := NewTaskRepoMock(false)

	task.ID = str2uuid("11111111-1111-1111-1111-111111111111")
	task.ManagedEndpointID = ExistedManagedEndpointID
	task.OriginID = gocql.TimeUUID()

	if err := repo.InsertOrUpdate(getContextWithTransactionID(t), task); err != nil {
		t.Fatal(err)
	}
}

func TestGetByRunTimeRange(t *testing.T) {
	repo := NewTaskRepoMock(true)
	tasks, err := repo.GetByRunTimeRange(getContextWithTransactionID(t), []time.Time{someTime.Add(time.Minute)})
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks, found %v", len(tasks))
	}

	tasks, err = repo.GetByRunTimeRange(context.Context(context.WithValue(context.TODO(), IsNeedError, true)), []time.Time{{}})
	if err == nil {
		t.Fatal("Error expected, but got nil")
	}
	if len(tasks) > 0 {
		t.Fatalf("No tasks expected")
	}
}

func TestGetByTargetID_Negative(t *testing.T) {
	repo := NewTaskRepoMock(true)
	if _, err := repo.GetByPartnerAndManagedEndpointID(setContext(t, true), somePartnerID, ExistedManagedEndpointID, common.UnlimitedCount); err == nil {
		t.Fatalf("Error expected, but got <nil>")
	}
}

func TestGetByTargetID_Positive(t *testing.T) {
	repo := NewTaskRepoMock(true)
	if tasks, err := repo.GetByPartnerAndManagedEndpointID(setContext(t, false), PartnerID, DefaultTasks[1].ManagedEndpointID, common.UnlimitedCount); err != nil {
		t.Fatal(err)
	} else if len(tasks) == 0 {
		t.Fatalf("Non-nil slice expected")
	}
}

func TestGetByTargetID_Equal(t *testing.T) {
	repo := NewTaskRepoMock(true)
	if _, err := repo.GetByPartnerAndManagedEndpointID(setContext(t, false), somePartnerID, emptyUUID, common.UnlimitedCount); err != nil {
		t.Fatal(err)
	}
}

func TestGetCountByTargetID_WithNeedErr(t *testing.T) {
	repo := NewTaskRepoMock(true)
	_, err := repo.GetCountByManagedEndpointID(setContext(t, true), somePartnerID, ExistedManagedEndpointID)
	if err == nil {
		t.Fatalf("Expected an error")
	}
}

func TestGetCountByTargetID_Positive(t *testing.T) {
	repo := NewTaskRepoMock(true)
	taskCountObj, err := repo.GetCountByManagedEndpointID(setContext(t, false), PartnerID, ExistedManagedEndpointID)
	if err != nil {
		t.Error(err)
	}

	var repoTasksByTargetID []models.Task
	for _, task := range DefaultTasks {
		if task.ManagedEndpointID == ExistedManagedEndpointID && task.PartnerID == PartnerID {
			repoTasksByTargetID = append(repoTasksByTargetID, task)
		}
	}

	if taskCountObj.Count != len(repoTasksByTargetID) {
		t.Fatalf("GetCountByTargetID(): got count = %v, want %v", taskCountObj.Count, len(repoTasksByTargetID))
	}
}

func Test_str2uuid(t *testing.T) {
	perfectUUID, _ := gocql.ParseUUID("55555555-5555-5555-5555-555555555555")

	tests := []struct {
		name       string
		stringUUID string
		want       gocql.UUID
	}{
		{
			name:       "good",
			stringUUID: "55555555-5555-5555-5555-555555555555",
			want:       perfectUUID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := str2uuid(tt.stringUUID); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("str2uuid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTaskRepoMock(t *testing.T) {
	filledTaskRepoMock, unfilledTaskRepoMock := TaskRepoMock{}, TaskRepoMock{}

	unfilledTaskRepoMock.Repo = make(map[gocql.UUID]models.Task)
	filledTaskRepoMock.Repo = make(map[gocql.UUID]models.Task)

	for _, task := range DefaultTasks {
		filledTaskRepoMock.Repo[task.ID] = task
	}

	tests := []struct {
		name string
		fill bool
		want TaskRepoMock
	}{
		{
			name: "filled",
			fill: true,
			want: filledTaskRepoMock,
		},
		{
			name: "unfilled",
			fill: false,
			want: unfilledTaskRepoMock,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTaskRepoMock(tt.fill); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("NewTaskRepoMock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name  string
		slice []interface{}
		value interface{}
		want  bool
	}{
		{
			name:  "does contain",
			slice: []interface{}{"1", "2", "3", "4", "value", "5"},
			value: "value",
			want:  true,
		},
		{
			name:  "does not contain",
			slice: []interface{}{"1", "2", "3", "4", "5"},
			value: "value",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.slice, tt.value); got != tt.want {
				t.Fatalf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name        string
		mock        TaskRepoMock
		ctx         context.Context
		inputStruct interface{}
		partnerID   string
		taskID      gocql.UUID
		wantErr     bool
	}{
		{
			name:        "good",
			mock:        NewTaskRepoMock(true),
			ctx:         context.Context(context.WithValue(context.TODO(), IsNeedError, false)),
			inputStruct: nil,
			partnerID:   "0",
			taskID:      str2uuid("22222222-2222-2222-2222-222222222222"),
			wantErr:     false,
		},
		{
			name:        "bad",
			mock:        NewTaskRepoMock(false),
			ctx:         context.Context(context.WithValue(context.TODO(), IsNeedError, true)),
			inputStruct: nil,
			partnerID:   "0",
			taskID:      str2uuid("22222222-2222-2222-2222-222222222222"),
			wantErr:     true,
		},
		{
			name:        "bad",
			mock:        NewTaskRepoMock(false),
			ctx:         context.Context(context.WithValue(context.TODO(), IsNeedError, false)),
			inputStruct: nil,
			partnerID:   "0",
			taskID:      str2uuid("22222222-2222-2222-2222-222222222222"),
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.mock.UpdateTask(tt.ctx, tt.inputStruct, tt.partnerID, tt.taskID); (err != nil) != tt.wantErr {
				t.Fatalf("TaskRepoMock.UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCountForAllTargets(t *testing.T) {
	tests := []struct {
		name      string
		mock      TaskRepoMock
		ctx       context.Context
		partnerID string
		taskCount []models.TaskCount
		wantErr   bool
	}{
		{
			name:      "good",
			mock:      NewTaskRepoMock(true),
			ctx:       context.Context(context.WithValue(context.TODO(), IsNeedError, false)),
			partnerID: PartnerID,
			taskCount: []models.TaskCount{
				{
					ManagedEndpointID: ExistedManagedEndpointID,
					Count:             6,
				},
				{
					ManagedEndpointID: TargetID,
					Count:             1,
				},
			},
			wantErr: false,
		},
		{
			name:      "bad",
			mock:      NewTaskRepoMock(false),
			ctx:       context.Context(context.WithValue(context.TODO(), IsNeedError, true)),
			partnerID: PartnerID,
			taskCount: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskCount, err := tt.mock.GetCountForAllTargets(tt.ctx, tt.partnerID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("TaskRepoMock.GetCountForAllTargets() error = %v, wantErr %v", err, tt.wantErr)
			}
			sort.Sort(byTargetID(taskCount))
			sort.Sort(byTargetID(tt.taskCount))
			if !reflect.DeepEqual(taskCount, tt.taskCount) {
				t.Fatalf("Expected %v, but got %v", tt.taskCount, taskCount)
			}
		})
	}
}

func TestGetByIDAndTargets(t *testing.T) {
	tests := []struct {
		name               string
		ctx                context.Context
		partnerID          string
		taskID             gocql.UUID
		managedEndpointIDs []gocql.UUID
		want               []models.Task
		wantErr            bool
	}{
		{
			name:      "good",
			ctx:       context.Context(context.WithValue(context.TODO(), IsNeedError, false)),
			partnerID: PartnerID,
			taskID:    DefaultTasks[0].ID,
			managedEndpointIDs: []gocql.UUID{
				DefaultTasks[0].ManagedEndpointID,
				DefaultTasks[1].ManagedEndpointID,
				DefaultTasks[2].ManagedEndpointID,
				DefaultTasks[3].ManagedEndpointID,
				DefaultTasks[4].ManagedEndpointID,
				DefaultTasks[5].ManagedEndpointID,
			},
			want:    []models.Task{DefaultTasks[0]},
			wantErr: false,
		},
		{
			name:      "cassandra error",
			ctx:       context.Context(context.WithValue(context.TODO(), IsNeedError, true)),
			partnerID: PartnerID,
			taskID:    DefaultTasks[0].ID,
			managedEndpointIDs: []gocql.UUID{
				DefaultTasks[0].ManagedEndpointID,
				DefaultTasks[1].ManagedEndpointID,
				DefaultTasks[2].ManagedEndpointID,
				DefaultTasks[3].ManagedEndpointID,
				DefaultTasks[4].ManagedEndpointID,
				DefaultTasks[5].ManagedEndpointID,
			},
			want:    []models.Task{DefaultTasks[0]},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewTaskRepoMock(true)
			got, err := mock.GetByIDAndManagedEndpoints(tt.ctx, tt.partnerID, tt.taskID, tt.managedEndpointIDs...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("TaskRepoMock.GetByIDAndManagedEndpoints() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Fatalf("TaskRepoMock.GetByIDAndManagedEndpoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCountsByPartner(t *testing.T) {
	tests := []struct {
		name    string
		partner string
		wantErr bool
		want    []models.TaskCount
	}{
		{
			name:    "Good case",
			partner: PartnerID,
			wantErr: false,
			want: []models.TaskCount{
				{
					ManagedEndpointID: ExistedManagedEndpointID,
					Count:             6,
				},
				{
					ManagedEndpointID: TargetID,
					Count:             1,
				},
			},
		},
		{
			name:    "Bad case",
			partner: PartnerID,
			wantErr: true,
			want:    nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := NewTaskRepoMock(true)
			ctx := context.WithValue(context.TODO(), IsNeedError, test.wantErr)
			got, err := mock.GetCountsByPartner(ctx, test.partner)
			if (err != nil) != test.wantErr {
				t.Fatalf("TaskRepoMock.GetCountsByPartner() error = %v, wantErr = %v", err, test.wantErr)
			}
			sort.Sort(byTargetID(got))
			sort.Sort(byTargetID(test.want))
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("Got != want. Got %v, want %v", got, test.want)
			}
		})
	}
}

func TestGetByPartner(t *testing.T) {
	tests := []struct {
		name      string
		partner   string
		wantFill  bool
		wantErr   bool
		wantTasks []models.Task
	}{
		{
			name:      "Good",
			partner:   PartnerID,
			wantFill:  true,
			wantErr:   false,
			wantTasks: DefaultTasks,
		},
		{
			name:      "Bad",
			partner:   "bad partner",
			wantFill:  false,
			wantErr:   false,
			wantTasks: []models.Task{},
		},
		{
			name:      "Bad/Cassandra",
			partner:   PartnerID,
			wantFill:  true,
			wantErr:   true,
			wantTasks: []models.Task{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := NewTaskRepoMock(test.wantFill)
			ctx := context.WithValue(context.TODO(), IsNeedError, test.wantErr)
			got, err := mock.GetByPartner(ctx, test.partner)
			if (err != nil) != test.wantErr {
				t.Fatalf("TaskRepoMock.GetByPartner() error = %v, wantErr = %v", err, test.wantErr)
			}

			var mapTasksByPartner = make(map[string][]models.Task)
			for _, task := range mock.Repo {
				mapTasksByPartner[task.PartnerID] = append(mapTasksByPartner[task.PartnerID], task)
			}

			if (!test.wantErr || !test.wantFill) && len(got) != len(mapTasksByPartner[test.partner]) {
				t.Fatalf("len(got) = %v, want len = %v", len(got), len(mapTasksByPartner[test.partner]))
			}

			for _, task := range got {
				if _, ok := mock.Repo[task.ID]; !ok {
					t.Fatalf("Got task %v, but it doesn't exist in mock", task)
				}
			}
		})
	}
}
