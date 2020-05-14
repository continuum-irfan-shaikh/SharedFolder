package modelMocks

import (
	"context"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
	"reflect"
	"testing"
)

func TestGetTasksSummaryData(t *testing.T) {
	defaultTaskID2 := str2uuid("11111111-2222-2222-1111-111111111111")
	tests := []struct {
		name                string
		partnerID           string
		taskSummaries       []models.TaskSummaryData
		taskIDs             []gocql.UUID
		wantInternalErr     bool
		wantTaskNotFoundErr bool
		wantRepoFill        bool
		want                []models.TaskSummaryData
	}{
		{
			name:            "GetTasksSummaryData - good case",
			partnerID:       defaultPartner,
			taskSummaries:   DefaultTaskSummaries,
			wantInternalErr: false,
			wantRepoFill:    true,
			want:            DefaultTaskSummaries,
		},
		{
			name:            "GetTasksSummaryData - internal error",
			partnerID:       defaultPartner,
			taskSummaries:   []models.TaskSummaryData{},
			wantInternalErr: true,
			wantRepoFill:    true,
			want:            []models.TaskSummaryData{},
		},
		{
			name:                "GetTasksSummaryData by IDs - good",
			wantRepoFill:        true,
			wantInternalErr:     false,
			wantTaskNotFoundErr: false,
			partnerID:           defaultPartner,
			taskIDs:             []gocql.UUID{defaultTaskID},
			want:                DefaultTaskSummaries,
		},
		{
			name:                "GetTasksSummaryData - internal error",
			wantRepoFill:        true,
			wantInternalErr:     true,
			wantTaskNotFoundErr: false,
			partnerID:           defaultPartner,
			taskIDs:             []gocql.UUID{defaultTaskID},
			want:                []models.TaskSummaryData{},
		},
		{
			name:                "GetTasksSummaryData by IDs - internal error, empty repo",
			wantRepoFill:        false,
			wantInternalErr:     true,
			wantTaskNotFoundErr: false,
			partnerID:           defaultPartner,
			taskIDs:             []gocql.UUID{defaultTaskID},
			want:                []models.TaskSummaryData{},
		},
		{
			name:                "GetTasksSummaryData by IDs - internal error, not valid partner",
			wantRepoFill:        true,
			wantInternalErr:     true,
			wantTaskNotFoundErr: false,
			partnerID:           "bad partner",
			taskIDs:             []gocql.UUID{defaultTaskID},
			want:                []models.TaskSummaryData{},
		},
		{
			name:                "GetTasksSummaryData by IDs - TaskNotFoundErr",
			wantRepoFill:        false,
			wantInternalErr:     false,
			wantTaskNotFoundErr: true,
			partnerID:           defaultPartner,
			taskIDs:             []gocql.UUID{defaultTaskID},
			want:                []models.TaskSummaryData{},
		},
		{
			name:                "GetTasksSummaryData by IDs - TaskNotFoundErr2",
			wantRepoFill:        false,
			wantInternalErr:     false,
			wantTaskNotFoundErr: true,
			partnerID:           defaultPartner,
			taskIDs:             []gocql.UUID{defaultTaskID, defaultTaskID2},
			want:                []models.TaskSummaryData{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				summaryData []models.TaskSummaryData
				err         error
			)
			ctx := context.WithValue(context.TODO(), IsNeedError, test.wantInternalErr)
			repo := NewTaskSummaryRepoMock(test.wantRepoFill)

			if len(test.taskIDs) == 0 {
				summaryData, err = repo.GetTasksSummaryData(ctx, false, nil, defaultPartner)
			} else {
				summaryData, err = repo.GetTasksSummaryData(ctx, false, nil, defaultPartner, test.taskIDs...)
			}

			if errType, ok := err.(models.TaskNotFoundError); ok != test.wantTaskNotFoundErr {
				t.Fatalf("%s: got %v, but wantTaskNotFoundErr == %v", test.name, errType, test.wantTaskNotFoundErr)
			}
			if !test.wantTaskNotFoundErr && (err != nil) != test.wantInternalErr {
				t.Fatalf("%s: got err=%v, want err=%v", test.name, err, test.wantInternalErr)
			}
			if !reflect.DeepEqual(summaryData, test.want) {
				t.Fatalf("%s: got %v, want %v", test.name, summaryData, test.want)
			}
		})
	}
}

func TestUpdateTaskInstanceStatusCount(t *testing.T) {
	tests := []struct {
		name               string
		wantFill           bool
		wantErr            bool
		taskInstanceID     gocql.UUID
		successStatusCount int
		failureStatusCount int
	}{
		{
			name:               "Case 1: Update status of existing record",
			wantFill:           true,
			wantErr:            false,
			taskInstanceID:     DefaultTaskInstanceStatusCounts[0].TaskInstanceID,
			successStatusCount: 4,
			failureStatusCount: 1,
		},
		{
			name:               "Case 2: Update status of non-existing record",
			wantFill:           true,
			wantErr:            false,
			taskInstanceID:     gocql.TimeUUID(),
			successStatusCount: 4,
			failureStatusCount: 1,
		},
		{
			name:               "Case 3: Cassandra is down",
			wantFill:           false,
			wantErr:            true,
			taskInstanceID:     DefaultTaskInstanceStatusCounts[0].TaskInstanceID,
			successStatusCount: 4,
			failureStatusCount: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewTaskSummaryRepoMock(test.wantFill)
			ctx := context.TODO()
			ctx = context.WithValue(ctx, IsNeedError, test.wantErr)
			testSuccessCount := repo.RepoStatusCounts[test.taskInstanceID].SuccessCount
			testFailureCount := repo.RepoStatusCounts[test.taskInstanceID].FailureCount

			expectedSuccessCount := testSuccessCount + test.successStatusCount
			expectedFailureCount := testFailureCount + test.failureStatusCount

			err := repo.UpdateTaskInstanceStatusCount(ctx, test.taskInstanceID, expectedSuccessCount, expectedFailureCount)
			if (err != nil) != test.wantErr {
				t.Fatalf("%s: wantErr == %v, got err == %v", test.name, test.wantErr, err)
			}

			if repo.RepoStatusCounts[test.taskInstanceID].SuccessCount != expectedSuccessCount && !test.wantErr {
				t.Fatalf("%s: UpdateTaskInstanceStatusCount failed: expected final successCount == %v, got %v", test.name, expectedSuccessCount, repo.RepoStatusCounts[test.taskInstanceID].SuccessCount)
			}

			if repo.RepoStatusCounts[test.taskInstanceID].FailureCount != expectedFailureCount && !test.wantErr {
				t.Fatalf("%s: UpdateTaskInstanceStatusCount failed: expected final successCount == %v, got %v", test.name, expectedFailureCount, repo.RepoStatusCounts[test.taskInstanceID].FailureCount)
			}
		})
	}
}

func TestGetStatusCountsByTaskInstanceIDs(t *testing.T) {
	tests := []struct {
		name                 string
		taskInstanceIDs      []gocql.UUID
		wantErr              bool
		wantInstanceStatuses map[gocql.UUID]models.TaskInstanceStatusCount
	}{
		{
			name: `Good`,
			taskInstanceIDs: []gocql.UUID{
				DefaultTaskInstanceStatusCounts[0].TaskInstanceID,
				DefaultTaskInstanceStatusCounts[1].TaskInstanceID,
			},
			wantErr: false,
			wantInstanceStatuses: map[gocql.UUID]models.TaskInstanceStatusCount{
				DefaultTaskInstanceStatusCounts[0].TaskInstanceID: DefaultTaskInstanceStatusCounts[0],
				DefaultTaskInstanceStatusCounts[1].TaskInstanceID: DefaultTaskInstanceStatusCounts[1],
			},
		},
		{
			name: `Bad`,
			taskInstanceIDs: []gocql.UUID{
				DefaultTaskInstanceStatusCounts[0].TaskInstanceID,
				DefaultTaskInstanceStatusCounts[1].TaskInstanceID,
			},
			wantErr:              true,
			wantInstanceStatuses: map[gocql.UUID]models.TaskInstanceStatusCount{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.WithValue(context.TODO(), IsNeedError, test.wantErr)
			mock := NewTaskSummaryRepoMock(true)
			got, err := mock.GetStatusCountsByIDs(ctx, nil, nil, test.taskInstanceIDs)
			if (err != nil) != test.wantErr {
				t.Fatalf("Got err = %v, but want err = %v", err, test.wantErr)
			}
			if !reflect.DeepEqual(got, test.wantInstanceStatuses) {
				t.Fatalf("Got: \n%v\nWant: \n%v\n", got, test.wantInstanceStatuses)
			}
		})
	}
}
