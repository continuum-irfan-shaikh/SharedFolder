package cassandra

import (
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"github.com/gocql/gocql"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

// NewTargets - returns new instance of Targets repository
func NewTargets(conn cassandra.ISession) *Targets {
	return &Targets{
		conn: conn,
	}
}

// Targets - Targets repository
type Targets struct {
	conn cassandra.ISession
}

// GetTargetsByTaskID returns targets by task ID
func (t *Targets) GetTargetsByTaskID(partnerID string, taskID gocql.UUID) ([]gocql.UUID, error) {
	ids := make([]gocql.UUID, 0)
	query := `SELECT 
				target_type,
				targets
			FROM targets WHERE partner_id = ? AND task_id = ?`
	var trg models.Target
	params := []interface{}{
		&trg.Type,
		&trg.IDs,
	}

	if err := t.conn.Query(query, partnerID, taskID).Scan(params...); err != nil {
		if err == gocql.ErrNotFound {
			return ids, err
		}
		return ids, fmt.Errorf("can't get targets for task_id %v ; cause: err=%s", taskID, err.Error())
	}

	if trg.Type != models.ManagedEndpoint {
		return ids, fmt.Errorf("this method can't be used for other target types except ManagedEndpoint")
	}

	for _, str := range trg.IDs {
		uuid, err := gocql.ParseUUID(str)
		if err != nil {
			return ids, fmt.Errorf("wrong format of managed endpoint id: %s, cause : %s ", str, err.Error())
		}

		ids = append(ids, uuid)
	}

	return ids, nil
}

// Insert new Targets in repository
func (t *Targets) Insert(partnerID string, taskID gocql.UUID, targets models.Target) error {
	ttl := 0
	query := `INSERT INTO targets (
				partner_id,
				task_id,
				target_type,
				targets
			) VALUES (?, ?, ?, ?) USING TTL ?`
	fields := []interface{}{
		// INSERT
		partnerID,
		taskID,
		targets.Type,
		targets.IDs,
		// TTL
		ttl,
	}

	if err := t.conn.Query(query, fields...).Exec(); err != nil {
		return fmt.Errorf("can't update targets of task with id :%v; cause: %s", taskID, err)
	}
	return nil
}
