package modelMocks

import (
	"context"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

const (
	badPartnerID     = "bad"
	expectedErrorMsg = "expected_error"
)

// TaskCounterDefaultMock mock to perform actions with task_counters table
type TaskCounterDefaultMock struct{}

// GetCounters fetches all task's counters
func (TaskCounterDefaultMock) GetCounters(ctx context.Context, partnerID string, endpointID gocql.UUID) (counts []models.TaskCount, err error) {
	if partnerID == badPartnerID {
		return []models.TaskCount{{
			ManagedEndpointID: ExistedManagedEndpointID,
			Count:             0,
		}}, errors.New(expectedErrorMsg)
	}
	return []models.TaskCount{{
		ManagedEndpointID: ExistedManagedEndpointID,
		Count:             5,
	}}, nil
}

// IncreaseCounter increase task's counters
func (TaskCounterDefaultMock) IncreaseCounter(partnerID string, counters []models.TaskCount, isExternal bool) error {
	if partnerID == badPartnerID {
		return errors.New(expectedErrorMsg)
	}
	return nil
}

// DecreaseCounter is expected to use when task deletion will be implemented
func (TaskCounterDefaultMock) DecreaseCounter(partnerID string, counters []models.TaskCount, isExternal bool) error {
	if partnerID == badPartnerID {
		return errors.New(expectedErrorMsg)
	}
	return nil
}

// GetAllPartners ...
func (TaskCounterDefaultMock) GetAllPartners(ctx context.Context) (partnerIDs map[string]struct{}, err error) {
	return map[string]struct{}{"partner1": {}}, nil
}

// TaskCounterCustomizableMock mock to perform actions with task_counters table
type TaskCounterCustomizableMock struct {
	GetCountersF     func(ctx context.Context, partnerID string, endpointID gocql.UUID) ([]models.TaskCount, error)
	IncreaseCounterF func(partnerID string, counters []models.TaskCount, isExternal bool) error
	DecreaseCounterF func(partnerID string, counters []models.TaskCount, isExternal bool) error
	GetAllPartnersF  func(ctx context.Context) (partnerIDs map[string]struct{}, err error)
}

// GetCounters fetches all task's counters
func (mock TaskCounterCustomizableMock) GetCounters(ctx context.Context, partnerID string, endpointID gocql.UUID) (counts []models.TaskCount, err error) {
	return mock.GetCountersF(ctx, partnerID, endpointID)
}

// IncreaseCounter increase task's counters
func (mock TaskCounterCustomizableMock) IncreaseCounter(partnerID string, counters []models.TaskCount, isExternal bool) error {
	return mock.IncreaseCounterF(partnerID, counters, isExternal)
}

// DecreaseCounter is expected to use when task deletion will be implemented
func (mock TaskCounterCustomizableMock) DecreaseCounter(partnerID string, counters []models.TaskCount, isExternal bool) error {
	return mock.DecreaseCounterF(partnerID, counters, isExternal)
}

// GetAllPartners ...
func (mock TaskCounterCustomizableMock) GetAllPartners(ctx context.Context) (partnerIDs map[string]struct{}, err error) {
	return mock.GetAllPartnersF(ctx)
}
