package taskCounterCassandra

import (
	"context"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository"
	"github.com/gocql/gocql"
)

// TaskCounterCassandra is the instance of the TaskCounterPersistence
type taskCounterCassandra struct {
	BatchSize int
}

// New us a function to return new TaskCounter repository
func New(cassandraBatchSize int) repository.TaskCounter {
	return taskCounterCassandra{
		BatchSize: cassandraBatchSize,
	}
}

// GetCounters fetches all task's counters
func (taskCounterCassandra) GetCounters(ctx context.Context, partnerID string, endpointID gocql.UUID) (counts []models.TaskCount, err error) {
	var (
		endpoint      gocql.UUID
		internalTasks int
		empty         = gocql.UUID{}
	)
	counts = make([]models.TaskCount, 0)

	if endpointID == empty { // looks for all managed endpoints
		iterator := cassandra.QueryCassandra(ctx, `
	SELECT
		endpoint_id,
		external_tasks,
		internal_tasks 
		FROM 
	task_counters WHERE partner_id = ?
`,
			partnerID).Iter()

		for iterator.Scan(&endpoint, nil, &internalTasks) {
			counts = append(counts, models.TaskCount{
				ManagedEndpointID: endpoint,
				Count:             internalTasks,
			})
		}

		if err = iterator.Close(); err != nil || len(counts) == 0 {
			return []models.TaskCount{}, err
		}
		return

	}

	queryString := `SELECT
		endpoint_id, 
		external_tasks, 
		internal_tasks 
		FROM 
	task_counters WHERE partner_id = ? AND endpoint_id = ?`

	// looks for specific managedEndpoint
	err = cassandra.QueryCassandra(ctx, queryString, partnerID, endpointID).Scan(nil, nil, &internalTasks)
	if err != nil {
		return []models.TaskCount{}, err
	}

	return []models.TaskCount{{
		ManagedEndpointID: endpointID,
		Count:             internalTasks,
	}}, nil
}

// IncreaseCounter increase task's counters
func (tcc taskCounterCassandra) IncreaseCounter(partnerID string, counters []models.TaskCount, isExternal bool) error {
	var queryString string
	if isExternal {
		queryString = `
	UPDATE task_counters SET
		external_tasks = external_tasks +?
	WHERE partner_id = ? AND 
		  endpoint_id = ?`
	} else {
		queryString = `
	UPDATE task_counters SET 
		internal_tasks = internal_tasks +? 
	WHERE partner_id = ? AND 
	  endpoint_id = ?`
	}

	batch := cassandra.Session.NewBatch(gocql.UnloggedBatch)
	for i, counter := range counters {
		batch.Query(queryString, counter.Count, partnerID, counter.ManagedEndpointID)

		if (i+1)%tcc.BatchSize == 0 || i+1 == len(counters) {
			err := cassandra.Session.ExecuteBatch(batch)
			if err != nil {
				return err
			}

			batch = cassandra.Session.NewBatch(gocql.UnloggedBatch)
		}
	}

	return nil
}

// DecreaseCounter is expected to use when task deletion will be implemented
func (tcc taskCounterCassandra) DecreaseCounter(partnerID string, counters []models.TaskCount, isExternal bool) error {
	var queryString string
	if isExternal {
		queryString = `
		UPDATE task_counters SET 
			external_tasks = external_tasks -? 
		WHERE partner_id = ? AND 
			endpoint_id = ?`
	} else {
		queryString = `
		UPDATE task_counters SET
			internal_tasks = internal_tasks -?
		WHERE partner_id = ? AND 
			endpoint_id = ?`
	}

	batch := cassandra.Session.NewBatch(gocql.UnloggedBatch)
	for i, counter := range counters {
		batch.Query(queryString, counter.Count, partnerID, counter.ManagedEndpointID)

		if (i+1)%tcc.BatchSize == 0 || i+1 == len(counters) {
			err := cassandra.Session.ExecuteBatch(batch)
			if err != nil {
				return err
			}

			batch = cassandra.Session.NewBatch(gocql.UnloggedBatch)
		}
	}

	return nil
}

// GetAllPartners ...
func (taskCounterCassandra) GetAllPartners(ctx context.Context) (map[string]struct{}, error) {
	var (
		iteratorTask        = cassandra.QueryCassandra(ctx, `SELECT distinct partner_id FROM tasks`).Iter()
		iteratorTaskCounter = cassandra.QueryCassandra(ctx, `SELECT distinct partner_id FROM task_counters`).Iter()
		partnersMap         = make(map[string]struct{})
		partner             string
	)

	for iteratorTask.Scan(&partner) {
		partnersMap[partner] = struct{}{}
	}
	err := iteratorTask.Close()
	if err != nil {
		return partnersMap, err
	}

	for iteratorTaskCounter.Scan(&partner) {
		partnersMap[partner] = struct{}{}
	}

	err = iteratorTaskCounter.Close()
	return partnersMap, err
}
