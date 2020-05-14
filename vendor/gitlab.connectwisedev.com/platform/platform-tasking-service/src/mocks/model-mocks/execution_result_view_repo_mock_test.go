package modelMocks

import (
	"context"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"net/http"
	"reflect"
	"testing"
)

var (
	repo                 models.ExecutionResultViewPersistence
	executionResultsView []models.ExecutionResultView
	somePartnerID        = "01234567"
	defaultHistoryCount  = 10
)

func taskExecutionResultsRepoMockStartup(isFilled bool) {
	repo = NewExecutionResultViewRepoMock(isFilled)
	executionResultsView = mockExecutionResultsView
}

func setContext(t *testing.T, isNeedError bool) (ctx context.Context) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx = req.Context()
	ctx = context.WithValue(ctx, IsNeedError, isNeedError)
	return ctx
}

func TestTaskExecutionResultsRepoMockGetWithBadRepository(t *testing.T) {
	taskExecutionResultsRepoMockStartup(true)

	if _, err := repo.Get(setContext(t, true), somePartnerID, ExistedManagedEndpointID, common.UnlimitedCount, false); err == nil {
		t.Errorf("No error getting from the Bad Repository")
	}
}

func TestTaskExecutionResultsRepoMockGetPositive(t *testing.T) {
	taskExecutionResultsRepoMockStartup(true)

	actualExecutionResultsView, err := repo.Get(setContext(t, false), somePartnerID, ExistedManagedEndpointID, common.UnlimitedCount, false)
	if err != nil {
		t.Errorf("Error for getting existed endpoint: %v", err)
	}

	var expectedExecutionResultsView []*models.ExecutionResultView
	for _, executionResultView := range executionResultsView {
		if executionResultView.ManagedEndpointID == ExistedManagedEndpointID {
			expectedExecutionResultsView = append(expectedExecutionResultsView, &executionResultView)
		}
	}

	if !reflect.DeepEqual(actualExecutionResultsView, expectedExecutionResultsView) {
		t.Errorf("actual Task list are not equal to expected Task list")
	}

}

func TestTaskExecutionResultsRepoMockGetNegative(t *testing.T) {
	taskExecutionResultsRepoMockStartup(true)

	partnerID := somePartnerID

	_, err := repo.Get(setContext(t, false), partnerID, NotExistedManagedEndpointID, common.UnlimitedCount, false)
	if err != nil {
		t.Errorf("No error was expected, but got error: %v", err)
	}
}

func TestTaskExecutionResultsRepoMockHistoryPositive(t *testing.T) {
	taskExecutionResultsRepoMockStartup(true)

	actualExecutionResultsView, err := repo.History(setContext(t, false), somePartnerID, ExistedTaskID, ExistedManagedEndpointID, defaultHistoryCount, false)
	if err != nil {
		t.Errorf("Error for getting existed execution results history: %v", err)
	}

	var expectedExecutionResultsView []*models.ExecutionResultView
	for _, executionResultView := range executionResultsView {
		if executionResultView.ManagedEndpointID == ExistedManagedEndpointID {
			expectedExecutionResultsView = append(expectedExecutionResultsView, &executionResultView)
		}
	}

	if !reflect.DeepEqual(actualExecutionResultsView, expectedExecutionResultsView) {
		t.Errorf("actual results are not equal to expected results")
	}

}

func TestTaskExecutionResultsRepoMockHistoryWithBadRepository(t *testing.T) {
	taskExecutionResultsRepoMockStartup(true)

	if _, err := repo.History(setContext(t, true), somePartnerID, ExistedTaskID, ExistedManagedEndpointID, defaultHistoryCount, false); err == nil {
		t.Errorf("No error getting from the Bad Repository")
	}
}

func TestTaskExecutionResultsRepoMockHistoryNegative(t *testing.T) {
	taskExecutionResultsRepoMockStartup(true)

	_, err := repo.History(setContext(t, false), somePartnerID, NotExistedTaskID, NotExistedManagedEndpointID, defaultHistoryCount, false)
	if err == nil {
		t.Errorf("No error returned for non existed task")
	}
}
