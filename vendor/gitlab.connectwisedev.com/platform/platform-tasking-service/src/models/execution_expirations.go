package models

//go:generate mockgen -destination=../mocks/mocks-gomock/executionExpirationPersistence_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/models ExecutionExpirationPersistence

import (
	"context"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

const (
	executionExpirationFields                     = `expiration_time_utc, partner_id, task_instance_id, managed_endpoint_ids`
	insertExecutionExpirationCQL                  = `INSERT INTO execution_expirations (` + executionExpirationFields + `) VALUES (?, ?, ?, ?) USING TTL ?`
	selectExecutionExpirationsByExpirationTimeCQL = selectDb + executionExpirationFields + ` FROM execution_expirations WHERE expiration_time_utc = ?`
	selectExecutionExpirationsByTaskInstanceIDCQL = selectDb + executionExpirationFields + ` FROM execution_expirations WHERE partner_id = ? and task_instance_id = ? ALLOW FILTERING`
	errWorkingWithEntities                        = "error while working with found entities: %v"
)

// ExecutionExpiration describes execution expiration data for particular task instance ID
type ExecutionExpiration struct {
	ExpirationTimeUTC  time.Time
	PartnerID          string
	TaskInstanceID     gocql.UUID
	ManagedEndpointIDs []gocql.UUID
}

// ExecutionExpirationPersistence is an interface for executionExpiration persistence
type ExecutionExpirationPersistence interface {
	InsertExecutionExpiration(ctx context.Context, exp ExecutionExpiration) (err error)
	GetByExpirationTime(ctx context.Context, expirationTime time.Time) ([]ExecutionExpiration, error)
	GetByTaskInstanceIDs(partnerID string, tiID []gocql.UUID) ([]ExecutionExpiration, error)
	Delete(expirationTime ExecutionExpiration) error
}

// ExecutionExpirationRepoCassandra is executionExpiration persistence for Cassandra
type ExecutionExpirationRepoCassandra struct{}

// ExecutionExpirationPersistenceInstance is executionExpiration persistence instance
var ExecutionExpirationPersistenceInstance ExecutionExpirationPersistence = ExecutionExpirationRepoCassandra{}

// InsertExecutionExpiration inserts ExecutionExpiration in Cassandra
func (ExecutionExpirationRepoCassandra) InsertExecutionExpiration(ctx context.Context, exp ExecutionExpiration) error {
	return cassandra.QueryCassandra(ctx, insertExecutionExpirationCQL,
		exp.ExpirationTimeUTC,
		exp.PartnerID,
		exp.TaskInstanceID,
		exp.ManagedEndpointIDs,
		config.Config.DataRetentionIntervalDay*secondsInDay,
	).Exec()
}

// GetByExpirationTime returns ExecutionExpirations found by ExpirationTimeUTC truncated to minutes
func (ExecutionExpirationRepoCassandra) GetByExpirationTime(ctx context.Context, expirationTime time.Time) ([]ExecutionExpiration, error) {
	executionExpirations := make([]ExecutionExpiration, 0)
	iterator := cassandra.QueryCassandra(ctx, selectExecutionExpirationsByExpirationTimeCQL, expirationTime).Iter()

	var exp ExecutionExpiration
	for iterator.Scan(
		&exp.ExpirationTimeUTC,
		&exp.PartnerID,
		&exp.TaskInstanceID,
		&exp.ManagedEndpointIDs,
	) {
		executionExpirations = append(executionExpirations, ExecutionExpiration{
			ExpirationTimeUTC:  exp.ExpirationTimeUTC,
			PartnerID:          exp.PartnerID,
			TaskInstanceID:     exp.TaskInstanceID,
			ManagedEndpointIDs: exp.ManagedEndpointIDs,
		})
	}

	if err := iterator.Close(); err != nil {
		return nil, errors.Errorf(errWorkingWithEntities, err)
	}
	return executionExpirations, nil
}

// GetByTaskInstanceIDs returns exec expirations by tiIDs and partner
func (ExecutionExpirationRepoCassandra) GetByTaskInstanceIDs(partnerID string, tiIDs []gocql.UUID) ([]ExecutionExpiration, error) {
	var results []ExecutionExpiration
	for _, tiID := range tiIDs {
		iterator := cassandra.QueryCassandra(context.TODO(), selectExecutionExpirationsByTaskInstanceIDCQL, partnerID, tiID).Iter()
		var exp ExecutionExpiration
		for iterator.Scan(
			&exp.ExpirationTimeUTC,
			&exp.PartnerID,
			&exp.TaskInstanceID,
			&exp.ManagedEndpointIDs,
		) {
			results = append(results, ExecutionExpiration{
				ExpirationTimeUTC:  exp.ExpirationTimeUTC,
				PartnerID:          exp.PartnerID,
				TaskInstanceID:     exp.TaskInstanceID,
				ManagedEndpointIDs: exp.ManagedEndpointIDs,
			})
		}

		if err := iterator.Close(); err != nil {
			return nil, errors.Errorf(errWorkingWithEntities, err)
		}
	}
	return results, nil
}

// Delete deletes exec expirations
func (ExecutionExpirationRepoCassandra) Delete(exp ExecutionExpiration) error {
	return cassandra.QueryCassandra(context.TODO(), `DELETE FROM execution_expirations 
        WHERE partner_id = ? AND task_instance_id = ? AND expiration_time_utc = ?`,
		exp.PartnerID,
		exp.TaskInstanceID,
		exp.ExpirationTimeUTC,
	).Exec()
}
