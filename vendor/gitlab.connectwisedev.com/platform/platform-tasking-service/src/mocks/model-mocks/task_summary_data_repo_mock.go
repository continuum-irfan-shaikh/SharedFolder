package modelMocks

import (
	"context"
	"errors"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"github.com/gocql/gocql"
	"time"
)

const (
	defaultPartner = `1d4400c0`
	initiator      = `Andy`
)

var (
	defaultTaskID = str2uuid("11111111-2222-1111-1111-111111111111")
)

// DefaultTaskSummaryDetails is an array of predefined TaskSummaryDetails objects
var DefaultTaskSummaryDetails = []models.TaskSummaryDetails{
	{TaskSummary: DefaultTaskSummaries[0], InstanceSummaries: []models.TaskInstanceSummary{{ID: DefaultTaskInstances[0].ID, RunTime: DefaultTaskInstances[0].StartedAt, TargetSummaries: []models.TargetSummary{DefaultTargetSummaries[0]}}}},
	{TaskSummary: DefaultTaskSummaries[1], InstanceSummaries: []models.TaskInstanceSummary{{ID: DefaultTaskInstances[1].ID, RunTime: DefaultTaskInstances[1].StartedAt, TargetSummaries: []models.TargetSummary{DefaultTargetSummaries[1]}}}},
	{TaskSummary: DefaultTaskSummaries[2], InstanceSummaries: []models.TaskInstanceSummary{{ID: DefaultTaskInstances[2].ID, RunTime: DefaultTaskInstances[2].StartedAt, TargetSummaries: []models.TargetSummary{DefaultTargetSummaries[2]}}}},
	{TaskSummary: DefaultTaskSummaries[0], InstanceSummaries: []models.TaskInstanceSummary{{ID: DefaultTaskInstances[3].ID, RunTime: DefaultTaskInstances[3].StartedAt, TargetSummaries: []models.TargetSummary{DefaultTargetSummaries[3]}}}},
}

// DefaultTargetSummaries is an array of predefined TargetSummary objects
var DefaultTargetSummaries = []models.TargetSummary{
	{EndpointID: str2uuid("58a1af2f-6579-4aec-b45d-5dfde879ef01"), RunStatus: statuses.TaskInstanceSuccess, StatusDetails: "Done"},
	{EndpointID: str2uuid("58a1af2f-6579-4aec-b45d-5dfde879ef99"), RunStatus: statuses.TaskInstanceSuccess, StatusDetails: "Adobe Acrobat Reader has been installed successfully"},
	{EndpointID: str2uuid("58a1af2f-6579-4aec-b45d-5dfde879ef99"), RunStatus: statuses.TaskInstanceSuccess, StatusDetails: "507 files deleted successfully"},
	{EndpointID: str2uuid("58a1af2f-6579-4aec-b45d-5dfde879ef99"), RunStatus: statuses.TaskInstanceFailed, StatusDetails: ""},
	{EndpointID: str2uuid("74a86e19-2d9b-413a-a820-cf017e22026b"), RunStatus: statuses.TaskInstanceSuccess, StatusDetails: "Done"},
	{EndpointID: str2uuid("74a86e19-2d9b-413a-a820-cf017e22026b"), RunStatus: statuses.TaskInstanceSuccess, StatusDetails: "Adobe Acrobat Reader has been installed successfully"},
}

// DefaultTaskSummaries is an array of predefined TaskSummaryData objects
var DefaultTaskSummaries = []models.TaskSummaryData{
	{Name: "Name1", RunOn: models.TargetData{Count: 1, Type: models.ManagedEndpoint}, Regularity: apiModels.RunNow, InitiatedBy: initiator, Status: statuses.TaskStateActive, LastRunTime: time.Now()},
	{Name: "Name2", RunOn: models.TargetData{Count: 2, Type: models.ManagedEndpoint}, Regularity: apiModels.RunNow, InitiatedBy: initiator, Status: statuses.TaskStateActive, LastRunTime: time.Now()},
	{Name: "Name3", RunOn: models.TargetData{Count: 3, Type: models.ManagedEndpoint}, Regularity: apiModels.RunNow, InitiatedBy: initiator, Status: statuses.TaskStateActive, LastRunTime: time.Now()},
	{Name: "Name4", RunOn: models.TargetData{Count: 4, Type: models.ManagedEndpoint}, Regularity: apiModels.RunNow, InitiatedBy: initiator, Status: statuses.TaskStateActive, LastRunTime: time.Now()},
}

// DefaultTaskInstanceStatusCounts is an array of predefined TaskInstanceStatusCount objects
var DefaultTaskInstanceStatusCounts = []models.TaskInstanceStatusCount{
	{TaskInstanceID: DefaultTaskInstances[0].ID, SuccessCount: 1, FailureCount: 0},
	{TaskInstanceID: DefaultTaskInstances[1].ID, SuccessCount: 0, FailureCount: 1},
	{TaskInstanceID: DefaultTaskInstances[2].ID, SuccessCount: 1, FailureCount: 1},
	{TaskInstanceID: DefaultTaskInstances[3].ID, SuccessCount: 0, FailureCount: 0},
}

// NewTaskSummaryRepoMock creates a mock for TaskInstancesStatusesCount repository
func NewTaskSummaryRepoMock(fill bool) TaskSummaryRepoMock {
	mock := TaskSummaryRepoMock{}
	mock.RepoTaskSummaries = make(map[string]map[gocql.UUID][]models.TaskSummaryData)
	mock.RepoStatusCounts = make(map[gocql.UUID]models.TaskInstanceStatusCount)
	mock.RepoTaskSummaryDetails = make(map[string]map[gocql.UUID]models.TaskSummaryDetails)
	if fill {
		mock.RepoTaskSummaries[defaultPartner] = make(map[gocql.UUID][]models.TaskSummaryData)
		mock.RepoTaskSummaries[defaultPartner][defaultTaskID] = make([]models.TaskSummaryData, 0)
		mock.RepoTaskSummaries[defaultPartner][defaultTaskID] = append(mock.RepoTaskSummaries[defaultPartner][defaultTaskID], DefaultTaskSummaries...)
		for _, statusCount := range DefaultTaskInstanceStatusCounts {
			mock.RepoStatusCounts[statusCount.TaskInstanceID] = statusCount
		}
		for i, task := range DefaultTasks {
			if i == len(DefaultTaskSummaryDetails) {
				break
			}
			mock.RepoTaskSummaryDetails[defaultPartner] = make(map[gocql.UUID]models.TaskSummaryDetails)
			mock.RepoTaskSummaryDetails[defaultPartner][task.ID] = DefaultTaskSummaryDetails[i]
		}
	}
	return mock
}

// TaskSummaryRepoMock represents a mock for TaskInstanceStatusCount repository
type TaskSummaryRepoMock struct {
	RepoTaskSummaries      map[string]map[gocql.UUID][]models.TaskSummaryData
	RepoStatusCounts       map[gocql.UUID]models.TaskInstanceStatusCount
	RepoTaskSummaryDetails map[string]map[gocql.UUID]models.TaskSummaryDetails
}

// GetTasksSummaryData returns TaskSummaryData by PartnerID
func (mock TaskSummaryRepoMock) GetTasksSummaryData(ctx context.Context, isNOCUser bool, cache persistency.Cache, partnerID string, taskIDs ...gocql.UUID) (data []models.TaskSummaryData, err error) {
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return []models.TaskSummaryData{}, errors.New("cassandra is down")
	}

	if len(taskIDs) == 0 {
		return mock.RepoTaskSummaries[partnerID][defaultTaskID], nil
	}

	for _, taskID := range taskIDs {
		tasksSummaryData, ok := mock.RepoTaskSummaries[partnerID][taskID]
		if !ok {
			return []models.TaskSummaryData{}, models.TaskNotFoundError{}
		}
		data = append(data, tasksSummaryData...)
	}
	return data, nil

}

// UpdateTaskInstanceStatusCount updates TaskInstanceStatusCount in RepoMock
func (mock TaskSummaryRepoMock) UpdateTaskInstanceStatusCount(ctx context.Context, taskInstanceID gocql.UUID, successStatusCount, failureStatusCount int) error {
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return errors.New("cassandra is down")
	}
	if _, ok := mock.RepoStatusCounts[taskInstanceID]; !ok {
		mock.RepoStatusCounts[taskInstanceID] = models.TaskInstanceStatusCount{
			TaskInstanceID: taskInstanceID,
			SuccessCount:   0,
			FailureCount:   0,
		}
	}

	mock.RepoStatusCounts[taskInstanceID] = models.TaskInstanceStatusCount{
		SuccessCount: successStatusCount,
		FailureCount: failureStatusCount,
	}
	return nil
}

// GetStatusCountsByIDs gets Task Instances Statuses Counts from mocked repo
func (mock TaskSummaryRepoMock) GetStatusCountsByIDs(ctx context.Context, mCache persistency.Cache, taskInstancesMapByID map[gocql.UUID]models.TaskInstance, IDs []gocql.UUID) (map[gocql.UUID]models.TaskInstanceStatusCount, error) {
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return map[gocql.UUID]models.TaskInstanceStatusCount{}, errors.New("cassandra is down")
	}
	var results = make(map[gocql.UUID]models.TaskInstanceStatusCount)
	for _, id := range IDs {
		if _, ok := mock.RepoStatusCounts[id]; ok {
			results[id] = mock.RepoStatusCounts[id]
		}
	}
	return results, nil
}
