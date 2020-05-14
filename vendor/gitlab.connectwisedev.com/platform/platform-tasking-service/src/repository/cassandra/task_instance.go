package cassandra

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
)

const (
	taskInstancesTableName          = "task_instances"
	taskInstancesStartedAtTableName = "task_instances_started_at_mv"
	taskInstancesByIDTableName      = "task_instances_by_id_mv"
	UpdateKeyWord                   = "UPDATE "
)

// NewTaskInstance - returns new instance of Task Instances repository
func NewTaskInstance(conn cassandra.ISession) *TaskInstance {
	return &TaskInstance{
		conn: conn,
	}
}

// TaskInstance - Task Instances repository
type TaskInstance struct {
	conn cassandra.ISession
}

// GetInstancesForScheduled returns Instances with small part of info needed only for scheduled panel.
// retrieves info concurrently
func (t *TaskInstance) GetInstancesForScheduled(instanceIDs []string) ([]entities.TaskInstance, error) {
	if len(instanceIDs) == 0 {
		return nil, nil
	}

	query := `SELECT 
                    id, 
                    task_id, 
                    last_run_time, 
                    device_statuses, 
                    failure_count,
                    success_count,
            		endpoints
              FROM task_instances_by_id_mv 
              WHERE id = ?`
	return t.getInstancesForScheduled(instanceIDs, query)
}

// GetByStartedAtAfter returns slice TaskInstance found by PartnerID and timestamp (started_at)
func (t *TaskInstance) GetByStartedAtAfter(partnerID string, from, to time.Time) ([]entities.TaskInstance, error) {
	query := `SELECT 
				partner_id,
				id,
				task_id,
				started_at,
				last_run_time, 
				device_statuses,
				failure_count,
				success_count,
				endpoints
       		FROM task_instances_started_at_mv WHERE partner_id = ? AND started_at > ? AND started_at < ?`
	return t.selectTaskInstances(
		query,
		partnerID,
		from,
		to,
	)
}

func (t *TaskInstance) selectTaskInstances(selectQuery string, values ...interface{}) ([]entities.TaskInstance, error) {
	var (
		taskInstance   entities.TaskInstance
		queries        = make([]entities.TaskInstance, 0)
		cassandraQuery = t.conn.Query(selectQuery, values...)
		deviceStatuses = make(map[string]statuses.TaskInstanceStatus)
		endpoints      = make(map[string]statuses.TaskInstanceStatus)
	)

	iterator := cassandraQuery.Iter()
	for iterator.Scan(
		&taskInstance.PartnerID,
		&taskInstance.ID,
		&taskInstance.TaskID,
		&taskInstance.StartedAt,
		&taskInstance.LastRunTime,
		&deviceStatuses,
		&taskInstance.FailureCount,
		&taskInstance.SuccessCount,
		&endpoints,
	) {

		ti := entities.TaskInstance{
			PartnerID:    taskInstance.PartnerID,
			ID:           taskInstance.ID,
			TaskID:       taskInstance.TaskID,
			StartedAt:    taskInstance.StartedAt,
			LastRunTime:  taskInstance.LastRunTime,
			Statuses:     t.mergeDeviceStatuses(deviceStatuses, endpoints),
			FailureCount: taskInstance.FailureCount,
			SuccessCount: taskInstance.SuccessCount,
		}

		ti.PreparePendingStatuses()

		queries = append(queries, ti)
	}

	if err := iterator.Close(); err != nil {
		return nil, fmt.Errorf("error while working with found entities: %v", err)
	}
	return queries, nil
}

// GetTopInstancesForScheduledByTaskIDs ..
func (t *TaskInstance) GetTopInstancesForScheduledByTaskIDs(taskIDs []string) ([]entities.TaskInstance, error) {
	query := `SELECT 
                    id, 
                    task_id, 
                    last_run_time, 
                    device_statuses, 
                    failure_count,
                    success_count, 
					endpoints
        	  FROM task_instances 
              WHERE task_id = ?
              LIMIT 2`
	return t.getInstancesForScheduled(taskIDs, query)
}

func (t *TaskInstance) getInstancesForScheduled(ids []string, query string) ([]entities.TaskInstance, error) {
	var (
		taskInstances = make([]entities.TaskInstance, 0)
		results       = make(chan entities.TaskInstance, config.Config.CassandraConcurrentCallNumber)
		limit         = make(chan struct{}, config.Config.CassandraConcurrentCallNumber)
		done          = make(chan struct{}, 2)
		errChan       = make(chan error, len(ids))
		err           error
	)

	go func() {
		for e := range errChan {
			err = fmt.Errorf("%v %v", err, e)
		}
		done <- struct{}{}
	}()

	go func() {
		for ti := range results {
			taskInstances = append(taskInstances, ti)
		}
		done <- struct{}{}
	}()

	var wg sync.WaitGroup
	for _, id := range ids {
		limit <- struct{}{}
		wg.Add(1)
		go func(id string) {
			defer func() { <-limit }()
			defer wg.Done()
			var (
				deviceStatuses = make(map[string]statuses.TaskInstanceStatus)
				endpoints      = make(map[string]statuses.TaskInstanceStatus)
				ti             entities.TaskInstance
			)
			iter := t.conn.Query(query, id).Iter()
			for iter.Scan(
				&ti.ID,
				&ti.TaskID,
				&ti.LastRunTime,
				&deviceStatuses,
				&ti.FailureCount,
				&ti.SuccessCount,
				&endpoints,
			) {
				instance := entities.TaskInstance{
					ID:           ti.ID,
					TaskID:       ti.TaskID,
					LastRunTime:  ti.LastRunTime,
					Statuses:     t.mergeDeviceStatuses(deviceStatuses, endpoints),
					FailureCount: ti.FailureCount,
					SuccessCount: ti.SuccessCount,
				}
				instance.PreparePendingStatuses()
				results <- instance
			}

			if err := iter.Close(); err != nil {
				if err == gocql.ErrNotFound {
					return
				}
				errChan <- err
				return
			}
		}(id)
	}

	wg.Wait()
	close(results)
	<-done
	close(errChan)
	<-done
	close(limit)
	return taskInstances, err
}

// GetInstance - returns Task Instances by id
func (t *TaskInstance) GetInstance(id gocql.UUID) (models.TaskInstance, error) {
	deviceStatuses := make(map[gocql.UUID]statuses.TaskInstanceStatus)
	endpoints := make(map[gocql.UUID]statuses.TaskInstanceStatus)

	query := `SELECT 
				partner_id, 
				id,
				task_id,
				origin_id,
				started_at,
				last_run_time,
				device_statuses,
				failure_count, 
				success_count, 
				triggered_by, 
				endpoints 
			  FROM task_instances_by_id_mv WHERE id = ?`

	var ti models.TaskInstance
	params := []interface{}{
		&ti.PartnerID,
		&ti.ID,
		&ti.TaskID,
		&ti.OriginID,
		&ti.StartedAt,
		&ti.LastRunTime,
		&deviceStatuses,
		&ti.FailureCount,
		&ti.SuccessCount,
		&ti.TriggeredBy,
		&endpoints,
	}

	if err := t.conn.Query(query, id).Scan(params...); err != nil {
		return ti, fmt.Errorf("can't get task instance by id %v: err=%s", id, err)
	}
	ti.Statuses = t.mergeDeviceStatusesLegacy(deviceStatuses, endpoints)
	ti.PreparePendingStatuses()
	return ti, nil
}

//Insert new TaskInstance in repository with parameters from task instance
func (t *TaskInstance) Insert(ti models.TaskInstance, ttl int) (err error) {
	err = t.insertTaskInstance(ti, ttl)
	if err != nil {
		return
	}
	err = t.insertTaskInstanceStartedAt(ti, ttl)
	if err != nil {
		return
	}
	return t.insertInstancesByID(ti, ttl)
}

func (t *TaskInstance) insertTaskInstance(taskInstance models.TaskInstance, ttl int) error {
	fields := []interface{}{
		// TTL
		ttl,
		// INSERT
		taskInstance.PartnerID,
		taskInstance.Name,
		taskInstance.OriginID,
		taskInstance.LastRunTime,
		taskInstance.Statuses,
		taskInstance.FailureCount,
		taskInstance.SuccessCount,
		taskInstance.OverallStatus,
		taskInstance.TriggeredBy,
		// PK
		taskInstance.TaskID,
		taskInstance.StartedAt,
		taskInstance.ID,
	}

	query := UpdateKeyWord + taskInstancesTableName + ` 
			  USING TTL ?
			  SET 
 					partner_id = ?, 
			        name = ?, 
			        origin_id = ?, 
			        last_run_time = ?,
					endpoints = endpoints + ?, 
			        failure_count = ?, 
			        success_count = ?, 
			        overall_status = ?, 
			        triggered_by = ?
			  WHERE 
				task_id = ? AND 
				started_at = ? AND 
				id = ?;`

	if err := t.conn.Query(query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "error while inserting task instance to %s table", taskInstancesTableName)
	}
	return nil
}

func (t *TaskInstance) insertTaskInstanceStartedAt(taskInstance models.TaskInstance, ttl int) error {
	fields := []interface{}{
		// TTL
		ttl,
		// INSERT
		taskInstance.Name,
		taskInstance.OriginID,
		taskInstance.LastRunTime,
		taskInstance.Statuses, // sets device_statuses for backward compatibility
		taskInstance.Statuses,
		taskInstance.FailureCount,
		taskInstance.SuccessCount,
		taskInstance.OverallStatus,
		taskInstance.TriggeredBy,
		// PK
		taskInstance.PartnerID,
		taskInstance.StartedAt,
		taskInstance.ID,
		taskInstance.TaskID,
	}

	query := UpdateKeyWord + taskInstancesStartedAtTableName + ` 
			  USING TTL ?
			  SET 
			        name = ?, 
			        origin_id = ?, 
			        last_run_time = ?,
					device_statuses = ?,
					endpoints = endpoints + ?,
					failure_count = ?, 
			        success_count = ?, 
			        overall_status = ?, 
			        triggered_by = ?
			WHERE 
				partner_id = ? AND
				started_at = ? AND 
				id = ? AND
				task_id = ?`

	if err := t.conn.Query(query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "error while inserting task instance to %s table", taskInstancesStartedAtTableName)
	}
	return nil
}

func (t *TaskInstance) insertInstancesByID(taskInstance models.TaskInstance, ttl int) error {
	fields := []interface{}{
		// TTL
		ttl,
		// INSERT
		taskInstance.PartnerID,
		taskInstance.Name,
		taskInstance.OriginID,
		taskInstance.LastRunTime,
		taskInstance.Statuses, // sets device_statuses for backward compatibility
		taskInstance.Statuses,
		taskInstance.FailureCount,
		taskInstance.SuccessCount,
		taskInstance.OverallStatus,
		taskInstance.TriggeredBy,
		// PK
		taskInstance.ID,
		taskInstance.TaskID,
		taskInstance.StartedAt,
	}

	query := UpdateKeyWord + taskInstancesByIDTableName + ` 
			  USING TTL ?
			  SET  
					partner_id = ?,
					name = ?, 
			        origin_id = ?, 
			        last_run_time = ?,
					device_statuses = ?,
					endpoints = endpoints + ?, 
			        failure_count = ?, 
			        success_count = ?, 
			        overall_status = ?, 
			        triggered_by = ?
			WHERE 
				id = ? AND
				task_id = ? AND
				started_at = ?`

	if err := t.conn.Query(query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "error while inserting task instance to %s table", taskInstancesByIDTableName)
	}
	return nil
}

// GetNearestInstanceAfter returns nearest instance after
func (t *TaskInstance) GetNearestInstanceAfter(taskID gocql.UUID, sinceDate time.Time) (models.TaskInstance, error) {
	deviceStatuses := make(map[gocql.UUID]statuses.TaskInstanceStatus)
	endpoints := make(map[gocql.UUID]statuses.TaskInstanceStatus)

	query := `SELECT 
				partner_id, 
				id, 
				task_id, 
				origin_id, 
				started_at,
				last_run_time,
				device_statuses,
				failure_count, 
				success_count, 
				triggered_by,
				endpoints
			  FROM task_instances 
			  WHERE task_id = ? 
			  AND started_at > ?
			  ORDER BY started_at ASC
			  LIMIT 1`

	var ti models.TaskInstance
	params := []interface{}{
		&ti.PartnerID,
		&ti.ID,
		&ti.TaskID,
		&ti.OriginID,
		&ti.StartedAt,
		&ti.LastRunTime,
		&deviceStatuses,
		&ti.FailureCount,
		&ti.SuccessCount,
		&ti.TriggeredBy,
		&endpoints,
	}

	if err := t.conn.Query(query, taskID, sinceDate).Scan(params...); err != nil {
		if err == gocql.ErrNotFound {
			return ti, err
		}
		return ti, fmt.Errorf("can't get task instance by task_id %v since date %v: err=%s", taskID, sinceDate, err.Error())
	}
	ti.Statuses = t.mergeDeviceStatusesLegacy(deviceStatuses, endpoints)
	ti.PreparePendingStatuses()
	return ti, nil
}

// GetMinimalInstanceByID - returns Task Instances minimal info (task_id and name) by id
func (t *TaskInstance) GetMinimalInstanceByID(id string) (entities.TaskInstance, error) {
	query := `SELECT 
				partner_id, 
				task_id, 
				name 
			FROM task_instances_by_id_mv WHERE id = ?`

	var ti entities.TaskInstance
	params := []interface{}{
		&ti.PartnerID,
		&ti.TaskID,
		&ti.TaskName,
	}

	q := t.conn.Query(query, id)
	q.SetConsistency(gocql.One)
	defer q.Release()
	if err := q.Scan(params...); err != nil {
		return ti, errors.Wrapf(err, "can't get minimal task instance by id %v", id)
	}
	return ti, nil
}

//Insert new TaskInstance in repository with parameters from task instance
func (t *TaskInstance) RemoveInactiveEndpoints(ti models.TaskInstance, endpoints ...gocql.UUID) (err error) {
	err = t.removeInactiveEndpoints(ti, endpoints...)
	if err != nil {
		return
	}
	err = t.removeInactiveEndpointsStartedAt(ti, endpoints...)
	if err != nil {
		return
	}
	return t.removeInactiveEndpointsByID(ti, endpoints...)
}

func (t *TaskInstance) removeInactiveEndpoints(taskInstance models.TaskInstance, endpoints ...gocql.UUID) error {
	fields := []interface{}{
		// DELETE
		taskInstance.Statuses, // sets device_statuses for backward compatibility
		endpoints,
		// PK
		taskInstance.TaskID,
		taskInstance.StartedAt,
		taskInstance.ID,
	}

	query := UpdateKeyWord + taskInstancesTableName + ` 
			  SET
					device_statuses = ?,
					endpoints = endpoints - ?
			  WHERE 
				task_id = ? AND 
				started_at = ? AND 
				id = ?;`

	if err := t.conn.Query(query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "removeInactiveEndpoints: error while removing inactive endpoints in %s table", taskInstancesTableName)
	}
	return nil
}

func (t *TaskInstance) removeInactiveEndpointsStartedAt(taskInstance models.TaskInstance, endpoints ...gocql.UUID) error {
	fields := []interface{}{
		// DELETE
		taskInstance.Statuses, // sets device_statuses for backward compatibility
		endpoints,
		// PK
		taskInstance.PartnerID,
		taskInstance.StartedAt,
		taskInstance.ID,
		taskInstance.TaskID,
	}

	query := UpdateKeyWord + taskInstancesStartedAtTableName + ` 
			  SET
					device_statuses = ?,
					endpoints = endpoints - ?
			  WHERE 
				partner_id = ? AND
				started_at = ? AND 
				id = ? AND
				task_id = ?`

	if err := t.conn.Query(query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "removeInactiveEndpointsStartedAt: error while removing inactive endpoints in %s table", taskInstancesStartedAtTableName)
	}
	return nil
}

func (t *TaskInstance) removeInactiveEndpointsByID(taskInstance models.TaskInstance, endpoints ...gocql.UUID) error {
	fields := []interface{}{
		// DELETE
		taskInstance.Statuses, // sets device_statuses for backward compatibility
		endpoints,
		// PK
		taskInstance.ID,
		taskInstance.TaskID,
		taskInstance.StartedAt,
	}

	query := UpdateKeyWord + taskInstancesByIDTableName + ` 
			  SET  
					device_statuses = ?,
					endpoints = endpoints - ?
			  WHERE 
				id = ? AND
				task_id = ? AND
				started_at = ?`

	if err := t.conn.Query(query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "error while removing inactive endpoints in %s table", taskInstancesByIDTableName)
	}
	return nil
}

//Insert new TaskInstance in repository with parameters from task instance
func (t *TaskInstance) AppendNewEndpoints(ti models.TaskInstance, endpoints map[gocql.UUID]statuses.TaskInstanceStatus) (err error) {
	err = t.appendNewEndpoints(ti, endpoints)
	if err != nil {
		return
	}
	err = t.appendNewEndpointsStartedAt(ti, endpoints)
	if err != nil {
		return
	}
	return t.appendNewEndpointsByID(ti, endpoints)
}

func (t *TaskInstance) appendNewEndpoints(taskInstance models.TaskInstance, endpoints map[gocql.UUID]statuses.TaskInstanceStatus) error {
	fields := []interface{}{
		// SET
		endpoints,
		taskInstance.OverallStatus,
		taskInstance.LastRunTime,
		// PK
		taskInstance.TaskID,
		taskInstance.StartedAt,
		taskInstance.ID,
	}

	query := UpdateKeyWord + taskInstancesTableName + ` 
			  SET
					endpoints = endpoints + ?,
			        overall_status = ?,
					last_run_time = ?
			  WHERE 
				task_id = ? AND 
				started_at = ? AND 
				id = ?;`

	if err := t.conn.Query(query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "appendNewEndpoints: error while appending new endpoints in %s table", taskInstancesTableName)
	}
	return nil
}

func (t *TaskInstance) appendNewEndpointsStartedAt(taskInstance models.TaskInstance, endpoints map[gocql.UUID]statuses.TaskInstanceStatus) error {
	fields := []interface{}{
		// SET
		endpoints,
		taskInstance.OverallStatus,
		taskInstance.LastRunTime,
		// PK
		taskInstance.PartnerID,
		taskInstance.StartedAt,
		taskInstance.ID,
		taskInstance.TaskID,
	}

	query := UpdateKeyWord + taskInstancesStartedAtTableName + ` 
			  SET 
					endpoints = endpoints + ?,
			        overall_status = ?,
					last_run_time = ?
			  WHERE 
				partner_id = ? AND
				started_at = ? AND 
				id = ? AND
				task_id = ?`

	if err := t.conn.Query(query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "appendNewEndpointsStartedAt: error while appending new endpoints in %s table", taskInstancesStartedAtTableName)
	}
	return nil
}

func (t *TaskInstance) appendNewEndpointsByID(taskInstance models.TaskInstance, endpoints map[gocql.UUID]statuses.TaskInstanceStatus) error {
	fields := []interface{}{
		// SET
		endpoints,
		taskInstance.OverallStatus,
		taskInstance.LastRunTime,
		// PK
		taskInstance.ID,
		taskInstance.TaskID,
		taskInstance.StartedAt,
	}

	query := UpdateKeyWord + taskInstancesByIDTableName + ` 
			  SET  
					endpoints = endpoints + ?,
			        overall_status = ?,
					last_run_time = ?
			  WHERE 
				id = ? AND
				task_id = ? AND
				started_at = ?`

	if err := t.conn.Query(query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "appendNewEndpointsByID: error while appending new endpoints in %s table", taskInstancesByIDTableName)
	}
	return nil
}

func (TaskInstance) mergeDeviceStatuses(deviceStatuses map[string]statuses.TaskInstanceStatus, endpoints map[string]statuses.TaskInstanceStatus) map[string]statuses.TaskInstanceStatus {
	if len(deviceStatuses) == 0 {
		if endpoints == nil {
			return make(map[string]statuses.TaskInstanceStatus)
		}
		return endpoints
	}
	if len(endpoints) == 0 {
		return deviceStatuses
	}
	result := deviceStatuses
	for key, val := range endpoints {
		result[key] = val
	}
	return result
}

func (TaskInstance) mergeDeviceStatusesLegacy(deviceStatuses map[gocql.UUID]statuses.TaskInstanceStatus, endpoints map[gocql.UUID]statuses.TaskInstanceStatus) map[gocql.UUID]statuses.TaskInstanceStatus {
	if len(deviceStatuses) == 0 {
		if endpoints == nil {
			return make(map[gocql.UUID]statuses.TaskInstanceStatus)
		}
		return endpoints
	}
	if len(endpoints) == 0 {
		return deviceStatuses
	}
	result := deviceStatuses
	for key, val := range endpoints {
		result[key] = val
	}
	return result
}
