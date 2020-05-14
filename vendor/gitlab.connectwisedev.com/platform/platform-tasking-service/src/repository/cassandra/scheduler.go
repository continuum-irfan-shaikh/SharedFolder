package cassandra

import (
	"fmt"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
)

// NewScheduler  returns new Scheduler repo client
func NewScheduler(conn cassandra.ISession) *Scheduler {
	return &Scheduler{
		conn: conn,
	}
}

// Scheduler  repo client to work with scheduler table
type Scheduler struct {
	conn cassandra.ISession
}

// GetLastUpdate  returns Last update time for id = 1
func (s *Scheduler) GetLastUpdate() (time.Time, error) {
	query := `SELECT 
				last_update 
			FROM scheduler where id = 1`
	lastUpdate := time.Time{}
	if err := s.conn.Query(query).Scan(&lastUpdate); err != nil {
		return lastUpdate, err
	}
	return lastUpdate, nil
}

// UpdateScheduler updates last update
func (s *Scheduler) UpdateScheduler(time time.Time) error {
	query := `UPDATE scheduler SET last_update = ? where id = 1`
	if err := s.conn.Query(query, time).Exec(); err != nil {
		return fmt.Errorf("can't update scheduler: %v", err)
	}
	return nil
}

//GetLastExpiredExecutionCheck returns last expired execution check time
func (s *Scheduler) GetLastExpiredExecutionCheck() (lastUpdate time.Time, err error) {
	err = s.conn.Query(`SELECT
						last_update
					   FROM scheduler WHERE id = 2`).Scan(&lastUpdate)
	return
}

//UpdateLastExpiredExecutionCheck  updates last expired execution check time
func (s *Scheduler) UpdateLastExpiredExecutionCheck(time time.Time) (err error) {
	return s.conn.Query(`UPDATE scheduler SET last_update = ? WHERE id = 2`, time).Exec()
}
