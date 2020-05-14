package modelMocks

import (
	"context"
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
)

// ExecutionResultRepoMock is a Mocked data source with ScriptExecutionResults
// its structure is ManagedEndpointID -> TaskInstanceID -> list of ScriptExecutionResults
type ExecutionResultRepoMock struct {
	isNeedError bool
	data        map[gocql.UUID]map[gocql.UUID]models.ExecutionResult
}

// NewExecutionResultRepoMock returns a new ExecutionResultsRepoMock filled with data if needed
func NewExecutionResultRepoMock(isFilled, isNeedError bool) models.ExecutionResultPersistence {
	data := make(map[gocql.UUID]map[gocql.UUID]models.ExecutionResult)

	if isFilled {
		for _, executionResultView := range TaskExecutionResultsView {
			if _, ok := data[executionResultView.ManagedEndpointID]; !ok {
				data[executionResultView.ManagedEndpointID] = make(map[gocql.UUID]models.ExecutionResult)
			}
			data[executionResultView.ManagedEndpointID][executionResultView.ExecutionID] = models.ExecutionResult{
				ManagedEndpointID: executionResultView.ManagedEndpointID,
				TaskInstanceID:    executionResultView.ExecutionID,
				UpdatedAt:         executionResultView.LastRunTime,
				StdOut:            executionResultView.LastRunStdOut,
				StdErr:            executionResultView.LastRunStdErr,
				ExecutionStatus:   executionResultView.LastRunStatus,
			}
		}
	}

	return ExecutionResultRepoMock{
		isNeedError: isNeedError,
		data:        data,
	}
}

// GetByTaskInstanceIDs gets a slice of Script Execution Results found by task instance IDs in mocked repo
func (mock ExecutionResultRepoMock) GetByTaskInstanceIDs(taskInstanceIDs []gocql.UUID) ([]models.ExecutionResult, error) {
	if mock.isNeedError {
		return nil, fmt.Errorf("repository error")
	}

	var results = make([]models.ExecutionResult, 0)

	for _, mapExecResultsByTaskInstanceID := range mock.data {
		for _, taskInstanceID := range taskInstanceIDs {
			scriptExecutionResult, ok := mapExecResultsByTaskInstanceID[taskInstanceID]
			if ok {
				results = append(results, scriptExecutionResult)
			}
		}
	}
	return results, nil
}

// GetByTargetAndTaskInstanceIDs returns ExecutionResult found in Mocked repository
func (mock ExecutionResultRepoMock) GetByTargetAndTaskInstanceIDs(managedEndpointID gocql.UUID, taskInstanceIDs ...gocql.UUID) (results []models.ExecutionResult, err error) {
	if mock.isNeedError {
		return nil, fmt.Errorf("repository error")
	}
	taskInstances, ok := mock.data[managedEndpointID]
	if !ok {
		return nil, fmt.Errorf("no TaskInstance for ManagedEndpoint ID=%s found", managedEndpointID)
	}
	for _, taskInstanceID := range taskInstanceIDs {
		res, ok := taskInstances[taskInstanceID]
		if !ok {
			return nil, fmt.Errorf("no Results for TaskInstance ID=%s of ManagedEndpoint ID=%s found", taskInstanceID, managedEndpointID)
		}
		results = append(results, res)
	}
	return
}

// Upsert updates/inserts ExecutionResult in Mocked repository
func (mock ExecutionResultRepoMock) Upsert(ctx context.Context, partnerID, taskName string, results ...models.ExecutionResult) (err error) {
	if mock.isNeedError {
		return fmt.Errorf("repository error")
	}
	for _, res := range results {
		_, ok := mock.data[res.ManagedEndpointID]
		if !ok {
			mock.data[res.ManagedEndpointID] = make(map[gocql.UUID]models.ExecutionResult)
		}
		mock.data[res.ManagedEndpointID][res.TaskInstanceID] = res
	}
	return
}

// DeleteBatch ...
func (mock ExecutionResultRepoMock) DeleteBatch(ctx context.Context, executionResults []models.ExecutionResult) (err error) {
	return
}

// ExecutionResultCustomizableMock provides ability to mock any method of ExecutionResultPersistence interface
type ExecutionResultCustomizableMock struct {
	GetByTargetAndTaskInstanceIDsF func(managedEndpointID gocql.UUID, taskInstanceID ...gocql.UUID) ([]models.ExecutionResult, error)
	GetByTaskInstanceIDsF          func(taskInstanceIDs []gocql.UUID) ([]models.ExecutionResult, error)
	UpsertF                        func(ctx context.Context, partnerID, taskName string, results ...models.ExecutionResult) error
	DeleteBatchF                   func(context.Context, []models.ExecutionResult) error
}

// GetByTargetAndTaskInstanceIDs ...
func (mock ExecutionResultCustomizableMock) GetByTargetAndTaskInstanceIDs(managedEndpointID gocql.UUID, taskInstanceID ...gocql.UUID) ([]models.ExecutionResult, error) {
	return mock.GetByTargetAndTaskInstanceIDsF(managedEndpointID, taskInstanceID...)
}

// GetByTaskInstanceIDs ...
func (mock ExecutionResultCustomizableMock) GetByTaskInstanceIDs(taskInstanceIDs []gocql.UUID) ([]models.ExecutionResult, error) {
	return mock.GetByTaskInstanceIDsF(taskInstanceIDs)
}

// Upsert ...
func (mock ExecutionResultCustomizableMock) Upsert(ctx context.Context, partnerID, taskName string, results ...models.ExecutionResult) error {
	return mock.UpsertF(ctx, partnerID, taskName, results...)
}

// DeleteBatch ...
func (mock ExecutionResultCustomizableMock) DeleteBatch(ctx context.Context, executionResults []models.ExecutionResult) error {
	return mock.DeleteBatchF(ctx, executionResults)
}
