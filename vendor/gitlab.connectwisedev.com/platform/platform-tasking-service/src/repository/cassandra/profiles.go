package cassandra

import (
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"github.com/gocql/gocql"
)

// ProfilesRepo represents struct to get profiles data
type ProfilesRepo struct {
	conn cassandra.ISession
}

// NewProfilesRepo returns new *ProfilesRepo
func NewProfilesRepo(conn cassandra.ISession) *ProfilesRepo {
	return &ProfilesRepo{conn: conn}
}

// GetByTaskID retrieves profileID by taskID
func (p *ProfilesRepo) GetByTaskID(taskID gocql.UUID) (profileID gocql.UUID, err error) {
	query := `SELECT 
				profile_id 
			FROM profile_ids where task_id = ?`
	if err = p.conn.Query(query, taskID).Scan(&profileID); err != nil {
		return profileID, err
	}
	return
}

// Insert inserts profile and task IDs to cassandra
func (p *ProfilesRepo) Insert(taskID, profileID gocql.UUID) error {
	query := `INSERT INTO profile_ids (
				task_id,
				profile_id
			) VALUES (?, ?)`
	if err := p.conn.Query(query, taskID, profileID).Exec(); err != nil {
		return fmt.Errorf("can't insert profile: %v", err)
	}
	return nil
}

// Delete removes profile by taskID
func (p *ProfilesRepo) Delete(taskID gocql.UUID) error {
	query := `DELETE FROM profile_ids WHERE task_id = ? IF EXISTS`
	if err := p.conn.Query(query, taskID).Exec(); err != nil {
		return fmt.Errorf("can't delete profile: %v", err)
	}
	return nil
}
