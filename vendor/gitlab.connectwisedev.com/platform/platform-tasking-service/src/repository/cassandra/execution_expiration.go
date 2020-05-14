package cassandra

import (
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
)

// NewExecutionExpiration returns new ExecutionExpiration
func NewExecutionExpiration(conn cassandra.ISession) *ExecutionExpiration {
	return &ExecutionExpiration{
		conn: conn,
	}
}

// ExecutionExpiration repo struct for exec expirations
type ExecutionExpiration struct {
	conn cassandra.ISession
}

// Insert inserts new exec expirations
func (e *ExecutionExpiration) Insert(ex entities.ExecutionExpiration, ttl int) error {
	query := `INSERT INTO execution_expirations (
				expiration_time_utc,
				partner_id,
				task_instance_id,
				managed_endpoint_ids
				)
			  VALUES (?, ?, ?, ?) USING TTL ?`

	params := []interface{}{
		ex.ExpirationTimeUTC,
		ex.PartnerID,
		ex.TaskInstanceID,
		ex.ManagedEndpointIDs,
		ttl,
	}

	if err := e.conn.Query(query, params...).Exec(); err != nil {
		return fmt.Errorf("can't insert execution expirations: %v", err)
	}
	return nil
}
