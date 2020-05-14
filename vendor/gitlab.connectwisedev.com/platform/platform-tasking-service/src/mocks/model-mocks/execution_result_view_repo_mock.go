package modelMocks

import (
	"context"
	"fmt"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
)

// ExecutionResultViewRepoMock is a Mocked data source with ExecutionResultsView
// its structure is PartnerID -> ManagedEndpointID -> list of ExecutionResultsView
type ExecutionResultViewRepoMock struct {
	data map[gocql.UUID][]*models.ExecutionResultView
}

// NewExecutionResultViewRepoMock creates a mocked repository for ExecutionResultsView
func NewExecutionResultViewRepoMock(isFilled bool) (repo models.ExecutionResultViewPersistence) {
	data := make(map[gocql.UUID][]*models.ExecutionResultView)

	if isFilled {
		for _, executionResultView := range TaskExecutionResultsView {
			if _, ok := data[executionResultView.ManagedEndpointID]; !ok {
				data[executionResultView.ManagedEndpointID] = []*models.ExecutionResultView{}
			}
			data[executionResultView.ManagedEndpointID] = append(data[executionResultView.ManagedEndpointID], &executionResultView)
		}
	}
	return ExecutionResultViewRepoMock{
		data: data,
	}
}

// Get returns an execution results view of all tasks running on the specified Endpoint for the specified Partner
func (mock ExecutionResultViewRepoMock) Get(ctx context.Context, partnerID string, endpointID gocql.UUID, count int, isNOCUser bool) (executionResultsView []*models.ExecutionResultView, err error) {
	if isNeedError := ctx.Value(IsNeedError).(bool); isNeedError {
		return nil, fmt.Errorf("TaskExecutionResultsRepository interrior error")
	}
	executionResultsView, ok := mock.data[endpointID]
	if !ok {
		executionResultsView = []*models.ExecutionResultView{}
	}
	return
}

// History returns a history of execution results on the specified Task for the specified Partner
func (mock ExecutionResultViewRepoMock) History(ctx context.Context, partnerID string, taskID, endpointID gocql.UUID, count int, isNOCUser bool) (executionResultsView []*models.ExecutionResultView, err error) {
	if isNeedError := ctx.Value(IsNeedError).(bool); isNeedError {
		return nil, fmt.Errorf("TaskExecutionResultsRepository interrior error")
	}
	executionResultsView, ok := mock.data[taskID]
	if !ok {
		return nil, models.TaskNotFoundError{ErrorParameters: fmt.Sprintf("no such task ID (UUID=%s) for the partner (ID=%s)", taskID, partnerID)}
	}
	return
}
