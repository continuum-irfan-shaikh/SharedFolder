package modelMocks

import (
	"context"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

// ExecutionExpirationRepoMock stores executionExpiration's mocks
type ExecutionExpirationRepoMock struct {
	repo map[int64][]models.ExecutionExpiration
}

// GetByTaskInstanceIDs ..
func (m ExecutionExpirationRepoMock) GetByTaskInstanceIDs(partnerID string, tiID []gocql.UUID) ([]models.ExecutionExpiration, error) {
	return nil, nil
}

// Delete ..
func (m ExecutionExpirationRepoMock) Delete(expirationTime models.ExecutionExpiration) error {
	return nil
}

// NewExecutionExpirationRepoMock returns a new ExecutionExpirationRepoMock filled with data if needed
func NewExecutionExpirationRepoMock(isFilled bool) ExecutionExpirationRepoMock {
	mock := ExecutionExpirationRepoMock{}

	mock.repo = make(map[int64][]models.ExecutionExpiration)
	if isFilled {
		for _, exp := range DefaultExecutionExpirations {
			mock.repo[exp.ExpirationTimeUTC.Unix()] = append(mock.repo[exp.ExpirationTimeUTC.Unix()], exp)
		}
	}
	return mock
}

// InsertExecutionExpiration inserts ExecutionExpiration in Cassandra
func (m ExecutionExpirationRepoMock) InsertExecutionExpiration(ctx context.Context, exp models.ExecutionExpiration) error {
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return errors.New("Cassandra is down")
	}
	m.repo[exp.ExpirationTimeUTC.Unix()] = append(m.repo[exp.ExpirationTimeUTC.Unix()], exp)
	return nil
}

// GetByExpirationTime returns ExecutionExpirations found by ExpirationTimeUTC truncated to minutes
func (m ExecutionExpirationRepoMock) GetByExpirationTime(ctx context.Context, expirationTime time.Time) (exps []models.ExecutionExpiration, err error) {
	if isNeedError, _ := ctx.Value(IsNeedError).(bool); isNeedError {
		return exps, errors.New("Cassandra is down")
	}
	return m.repo[expirationTime.Unix()], nil
}

// GetMockRepo returns mocked repository
func (m ExecutionExpirationRepoMock) GetMockRepo() map[int64][]models.ExecutionExpiration {
	return m.repo
}

// DefaultExecutionExpirations is a list of ExecutionExpiration mocks
var (
	// CtrlExpirationTime is a control ExpirationTime
	CtrlExpirationTime = time.Date(2018, 7, 7, 7, 7, 7, 777777777, time.UTC)
	// CtrlTaskInstanceID is a control TaskInstanceID
	CtrlTaskInstanceID          = str2uuid("00000000-7777-7777-7777-000000000000")
	DefaultExecutionExpirations = []models.ExecutionExpiration{
		{
			ExpirationTimeUTC: time.Now().Add(time.Minute).Truncate(time.Minute),
			PartnerID:         partnerID,
			TaskInstanceID:    str2uuid("00000000-0000-0000-0000-000000000002"),
		},
		{
			ExpirationTimeUTC: CtrlExpirationTime,
			PartnerID:         partnerID,
			TaskInstanceID:    CtrlTaskInstanceID,
		},
	}
)
