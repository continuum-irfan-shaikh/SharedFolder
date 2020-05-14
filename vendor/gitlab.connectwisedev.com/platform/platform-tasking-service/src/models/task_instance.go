package models

//go:generate mockgen -destination=../mocks/mocks-gomock/taskInstancePersistance_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/models TaskInstancePersistence

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
)

const (
	selectDb                        = `SELECT `
	taskInstancesTableName          = "task_instances"
	taskInstancesStartedAtTableName = "task_instances_started_at_mv"
	taskInstancesByIDTableName      = "task_instances_by_id_mv"
	taskInstanceFields              = `partner_id, id, task_id, name, origin_id, started_at, last_run_time, device_statuses, failure_count, success_count, overall_status, triggered_by, endpoints`

	selectByTaskIDCqlQuery                       = selectDb + taskInstanceFields + ` FROM task_instances WHERE task_id = ?`
	selectByPartnerIDAndStartedAtAfterIDCqlQuery = selectDb + taskInstanceFields + ` FROM task_instances_started_at_mv WHERE partner_id = ? AND started_at > ? AND started_at < ?`
	selectCountByTaskID                          = `SELECT COUNT(id) FROM task_instances WHERE task_id = ?`
)

// TaskInstance describes task exemplar for particular managed endpoint ID and particular start time.
type TaskInstance struct {
	PartnerID     string                                     `json:"partnerId"`
	ID            gocql.UUID                                 `json:"id"`
	TaskID        gocql.UUID                                 `json:"taskId"`
	Name          string                                     `json:"-"`
	OriginID      gocql.UUID                                 `json:"originId"` //ScriptID or PatchID
	StartedAt     time.Time                                  `json:"startedAt"`
	LastRunTime   time.Time                                  `json:"lastRunTime"`
	Statuses      map[gocql.UUID]statuses.TaskInstanceStatus `json:"statuses"`
	OverallStatus statuses.TaskInstanceStatus                `json:"overallStatus"`
	FailureCount  int                                        `json:"failureCount"`
	SuccessCount  int                                        `json:"successCount"`
	TriggeredBy   string                                     `json:"triggeredBy"`
}

// PreparePendingStatuses - makes scheduled endpoints pending if running endpoints are existent
func (ti *TaskInstance) PreparePendingStatuses() {
	var hasRunning bool
	for _, s := range ti.Statuses {
		if s == statuses.TaskInstanceRunning ||
			s == statuses.TaskInstanceSuccess ||
			s == statuses.TaskInstanceFailed {
			hasRunning = true
			break
		}
	}
	if !hasRunning {
		return
	}
	for id, s := range ti.Statuses {
		if s == statuses.TaskInstanceScheduled {
			ti.Statuses[id] = statuses.TaskInstancePending
		}
	}
}

// NewTaskInstance - creates new task instance based on the Task
func NewTaskInstance(tasks []Task, createdByTask bool) TaskInstance {
	if len(tasks) == 0 {
		return TaskInstance{}
	}

	statusesMap := make(map[gocql.UUID]statuses.TaskInstanceStatus)
	defaultEID, err := gocql.ParseUUID(DefaultEndpointUID)
	if err != nil {
		return TaskInstance{}
	}

	for _, task := range tasks {
		var stat statuses.TaskInstanceStatus
		switch task.State {
		case statuses.TaskStateInactive:
			continue
		case statuses.TaskStateDisabled:
			stat = statuses.TaskInstanceDisabled
		case statuses.TaskStateActive:
			stat = statuses.TaskInstanceScheduled
			if task.HasPostponedTime() {
				stat = statuses.TaskInstancePostponed
			}
		}
		//if task is running on default managed endpoint ID - it means that DG/Site is empty,
		// so targets and statuses in instance should be empty
		if task.ManagedEndpointID != defaultEID {
			statusesMap[task.ManagedEndpointID] = stat
		}
	}

	// these values are the same for all tasks
	originTaskID := tasks[0].ID
	originTaskOriginID := tasks[0].OriginID
	partnerID := tasks[0].PartnerID
	name := tasks[0].Name
	triggeredBy := ""

	if !createdByTask {
		triggeredBy = strings.Join(tasks[0].Schedule.TriggerTypes, ",")
	}
	return TaskInstance{
		PartnerID:     partnerID,
		ID:            gocql.TimeUUID(),
		TaskID:        originTaskID,
		Name:          name,
		OriginID:      originTaskOriginID,
		StartedAt:     time.Now().UTC(),
		Statuses:      statusesMap,
		OverallStatus: statuses.TaskInstanceScheduled,
		TriggeredBy:   triggeredBy,
	}
}

// CalculateStatuses returns a common status for the entire instance
func (ti TaskInstance) CalculateStatuses() (map[string]int, error) {
	var statusCounts = make(map[string]int)
	if len(ti.Statuses) < 1 {
		return statusCounts, nil
	}

	for _, stat := range ti.Statuses {
		statusStr, err := statuses.TaskInstanceStatusText(stat)
		if err != nil {
			return nil, err
		}

		statusCounts[statusStr]++
	}

	return statusCounts, nil
}

// IsScheduled check if this TaskInstance has machine with scheduled status
func (ti TaskInstance) IsScheduled() bool {
	var disabledCount int

	for _, status := range ti.Statuses {
		if status == statuses.TaskInstanceDisabled {
			disabledCount++
			continue
		}

		if status == statuses.TaskInstanceScheduled {
			// FIX: RMM-36676
			return true
		}
		break
	}
	return len(ti.Statuses) == disabledCount
}

// TaskInstancePersistence interface to perform actions with Task database
type TaskInstancePersistence interface {
	Insert(ctx context.Context, taskInstance TaskInstance) error
	DeleteBatch(ctx context.Context, taskInstances []TaskInstance) error

	GetByIDs(ctx context.Context, taskInstanceIDs ...gocql.UUID) ([]TaskInstance, error)
	GetByTaskID(ctx context.Context, taskID gocql.UUID) ([]TaskInstance, error)
	GetTopInstancesByTaskID(ctx context.Context, taskID gocql.UUID) ([]TaskInstance, error)
	GetByStartedAtAfter(ctx context.Context, partnerID string, from, to time.Time) ([]TaskInstance, error)
	GetInstancesCountByTaskID(ctx context.Context, taskID gocql.UUID) (instancesCount int, err error)
	GetNearestInstanceAfter(taskID gocql.UUID, sinceDate time.Time) (TaskInstance, error)
	GetMinimalInstanceByID(ctx context.Context, id gocql.UUID) (TaskInstance, error)
	UpdateStatuses(ctx context.Context, taskInstance TaskInstance) (err error)
}

// TaskInstanceRepoCassandra is a realisation of TaskPersistence interface for Cassandra
type TaskInstanceRepoCassandra struct{}

var (
	// TaskInstancePersistenceInstance is an instance presented TaskInstanceRepoCassandra
	TaskInstancePersistenceInstance TaskInstancePersistence = TaskInstanceRepoCassandra{}
)

// GetByTaskID gets TaskInstance by TaskID
func (taskInst TaskInstanceRepoCassandra) GetByTaskID(ctx context.Context, taskID gocql.UUID) ([]TaskInstance, error) {
	return selectTaskInstances(ctx, selectByTaskIDCqlQuery, taskID)
}

// GetNearestInstanceAfter returns nearest instance after
func (taskInst TaskInstanceRepoCassandra) GetNearestInstanceAfter(taskID gocql.UUID, sinceDate time.Time) (TaskInstance, error) {
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

	var ti TaskInstance
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

	if err := cassandra.QueryCassandra(context.TODO(), query, taskID, sinceDate).Scan(params...); err != nil {
		if err == gocql.ErrNotFound {
			return ti, err
		}
		msg := "can't get task instance by task_id %v since date %v: err=%s"
		return ti, fmt.Errorf(msg, taskID, sinceDate, err.Error())
	}
	ti.Statuses = mergeDeviceStatuses(deviceStatuses, endpoints)
	ti.PreparePendingStatuses()
	return ti, nil
}

// GetTopInstancesByTaskID gets top 2 TaskInstances by TaskID
func (taskInst TaskInstanceRepoCassandra) GetTopInstancesByTaskID(ctx context.Context, taskID gocql.UUID) ([]TaskInstance, error) {
	query := `SELECT 
				partner_id,
				id, 
				task_id,
				name,
				origin_id, 
				started_at, 
				last_run_time,
				device_statuses,
				failure_count, 
				success_count,
				overall_status,
				triggered_by,
				endpoints
			FROM task_instances WHERE task_id = ? LIMIT 2`

	instances, err := selectTaskInstances(ctx, query, taskID)
	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return instances, nil
}

//Insert new TaskInstance in repository with parameters from task instance
func (t TaskInstanceRepoCassandra) Insert(ctx context.Context, taskInstance TaskInstance) (err error) {
	err = t.insertTaskInstance(ctx, taskInstance)
	if err != nil {
		return
	}
	err = t.insertTaskInstanceStartedAt(ctx, taskInstance)
	if err != nil {
		return
	}
	return t.insertInstancesByID(ctx, taskInstance)
}

func (TaskInstanceRepoCassandra) insertTaskInstance(ctx context.Context, taskInstance TaskInstance) error {
	fields := []interface{}{
		// TTL
		config.Config.DataRetentionIntervalDay * secondsInDay,
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
		taskInstance.TaskID,
		taskInstance.StartedAt,
		taskInstance.ID,
	}

	query := `UPDATE ` + taskInstancesTableName + ` 
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
				task_id = ? AND 
				started_at = ? AND 
				id = ?`

	if err := cassandra.QueryCassandra(ctx, query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "error while inserting task instance to %s table", taskInstancesTableName)
	}
	return nil
}

func (TaskInstanceRepoCassandra) insertTaskInstanceStartedAt(ctx context.Context, taskInstance TaskInstance) error {
	fields := []interface{}{
		// TTL
		config.Config.DataRetentionIntervalDay * secondsInDay,
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

	query := `UPDATE ` + taskInstancesStartedAtTableName + ` 
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

	if err := cassandra.QueryCassandra(ctx, query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "error while inserting task instance to %s table", taskInstancesStartedAtTableName)
	}
	return nil
}

func (TaskInstanceRepoCassandra) insertInstancesByID(ctx context.Context, taskInstance TaskInstance) error {
	fields := []interface{}{
		// TTL
		config.Config.DataRetentionIntervalDay * secondsInDay,
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

	query := `UPDATE ` + taskInstancesByIDTableName + ` 
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

	if err := cassandra.QueryCassandra(ctx, query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "error while inserting task instance to %s table", taskInstancesByIDTableName)
	}
	return nil
}

//UpdateStatuses updates statuses in current taskinstance
func (t TaskInstanceRepoCassandra) UpdateStatuses(ctx context.Context, taskInstance TaskInstance) (err error) {
	err = t.updateStatusesInInstances(ctx, taskInstance)
	if err != nil {
		return
	}
	err = t.updateStatusesInInstanceStartedAt(ctx, taskInstance)
	if err != nil {
		return
	}
	return t.updateStatusesInInstancesByID(ctx, taskInstance)
}

func (TaskInstanceRepoCassandra) updateStatusesInInstances(ctx context.Context, taskInstance TaskInstance) error {
	fields := []interface{}{
		// TTL
		config.Config.DataRetentionIntervalDay * secondsInDay,
		// UPDATE
		taskInstance.Statuses,
		taskInstance.FailureCount,
		taskInstance.SuccessCount,
		// PK
		taskInstance.TaskID,
		taskInstance.StartedAt,
		taskInstance.ID,
	}

	query := `UPDATE ` + taskInstancesTableName + ` 
			  USING TTL ?
			  SET 
					endpoints = endpoints + ?, 
			        failure_count = ?, 
			        success_count = ? 
			WHERE 
				task_id = ? AND 
				started_at = ? AND 
				id = ?`

	if err := cassandra.QueryCassandra(ctx, query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "error while inserting task instance without targets to %s table", taskInstancesTableName)
	}
	return nil
}

func (TaskInstanceRepoCassandra) updateStatusesInInstanceStartedAt(ctx context.Context, taskInstance TaskInstance) error {
	fields := []interface{}{
		// TTL
		config.Config.DataRetentionIntervalDay * secondsInDay,
		// INSERT
		taskInstance.Statuses,
		taskInstance.FailureCount,
		taskInstance.SuccessCount,
		// PK
		taskInstance.PartnerID,
		taskInstance.StartedAt,
		taskInstance.ID,
		taskInstance.TaskID,
	}

	query := `UPDATE ` + taskInstancesStartedAtTableName + ` 
			  USING TTL ?
			  SET 
					endpoints = endpoints + ?, 
			        failure_count = ?, 
			        success_count = ? 
			WHERE 
				partner_id = ? AND
				started_at = ? AND 
				id = ? AND
				task_id = ?`
	if err := cassandra.QueryCassandra(ctx, query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "error while inserting task instance without targets to %s table", taskInstancesStartedAtTableName)
	}
	return nil
}

func (TaskInstanceRepoCassandra) updateStatusesInInstancesByID(ctx context.Context, taskInstance TaskInstance) error {
	fields := []interface{}{
		// TTL
		config.Config.DataRetentionIntervalDay * secondsInDay,
		// INSERT
		taskInstance.Statuses,
		taskInstance.FailureCount,
		taskInstance.SuccessCount,
		// PK
		taskInstance.ID,
		taskInstance.TaskID,
		taskInstance.StartedAt,
	}

	query := `UPDATE ` + taskInstancesByIDTableName + ` 
			  USING TTL ?
			  SET 
					endpoints = endpoints + ?, 
			        failure_count = ?, 
			        success_count = ? 
			WHERE 
				id = ? AND
				task_id = ? AND
				started_at = ?`

	if err := cassandra.QueryCassandra(ctx, query, fields...).Exec(); err != nil {
		return errors.Wrapf(err, "error while inserting task instance without targets to %s table", taskInstancesByIDTableName)
	}
	return nil
}

const deleteFrom = `DELETE FROM `

// DeleteBatch deletes a batch of TaskInstances from repository
func (TaskInstanceRepoCassandra) DeleteBatch(ctx context.Context, taskInstances []TaskInstance) error {
	tables := []string{taskInstancesTableName, taskInstancesByIDTableName, taskInstancesStartedAtTableName}
	for _, ti := range taskInstances {
		for _, table := range tables {
			tiFields := []interface{}{
				ti.TaskID,
				ti.StartedAt,
				ti.ID,
			}

			var query = `DELETE FROM `
			switch table {
			case taskInstancesTableName, taskInstancesByIDTableName:
				query = deleteFrom + table +
					` WHERE task_id = ? AND started_at = ? AND id = ? IF EXISTS`
			case taskInstancesStartedAtTableName:
				query = deleteFrom + table +
					` WHERE task_id = ? AND started_at = ? AND id = ? AND partner_id = ? IF EXISTS`
				tiFields = append(tiFields, ti.PartnerID)
			}

			if err := cassandra.QueryCassandra(ctx, query, tiFields...).Exec(); err != nil {
				return fmt.Errorf("error while deleting set of task instances: %s", err.Error())
			}
		}
	}
	return nil
}

// GetByIDs returns slice TaskInstance found by IDs
func (TaskInstanceRepoCassandra) GetByIDs(ctx context.Context, IDs ...gocql.UUID) ([]TaskInstance, error) {
	var (
		selectQuery = `SELECT 
							partner_id,
							id, 
							task_id, 
							name, 
							origin_id, 
							started_at,
							last_run_time,
							device_statuses, 
							failure_count,
							success_count,
							overall_status,
							triggered_by,
							endpoints
						FROM task_instances_by_id_mv WHERE id = ?`

		taskInstances = make([]TaskInstance, 0, len(IDs))
		results       = make(chan []TaskInstance, config.Config.CassandraConcurrentCallNumber)
		limit         = make(chan struct{}, config.Config.CassandraConcurrentCallNumber)
		done          = make(chan struct{})
	)

	//avoid data race
	var resultErr error
	go func() {
		for taskInstance := range results {
			taskInstances = append(taskInstances, taskInstance...)
		}
		done <- struct{}{}
	}()

	var wg sync.WaitGroup
	for _, id := range IDs {
		taskInstanceID := id
		limit <- struct{}{}
		wg.Add(1)
		go func(id gocql.UUID, results chan []TaskInstance) {
			defer func() {
				<-limit
			}()
			defer wg.Done()
			instances, err := selectTaskInstances(ctx, selectQuery, id)
			if err != nil {
				resultErr = fmt.Errorf("%v : %v", resultErr, err)
				return
			}

			if len(instances) > 0 {
				results <- instances
			} else {
				resultErr = fmt.Errorf("%v : instance with id [%s] not found", resultErr, id.String())
			}
		}(taskInstanceID, results)
	}

	wg.Wait()
	close(results)
	<-done
	return taskInstances, resultErr
}

// GetInstancesCountByTaskID gets count of instances by task id
func (TaskInstanceRepoCassandra) GetInstancesCountByTaskID(ctx context.Context, taskID gocql.UUID) (instancesCount int, err error) {
	cassandraQuery := cassandra.QueryCassandra(ctx, selectCountByTaskID, taskID)
	iterator := cassandraQuery.Iter()
	iterator.Scan(&instancesCount)

	if err = iterator.Close(); err != nil {
		return instancesCount, fmt.Errorf(errWorkingWithEntities, err)
	}
	return
}

// GetByStartedAtAfter returns slice TaskInstance found by PartnerID and timestamp (started_at)
func (TaskInstanceRepoCassandra) GetByStartedAtAfter(ctx context.Context, partnerID string, from, to time.Time) ([]TaskInstance, error) {
	return selectTaskInstances(ctx, selectByPartnerIDAndStartedAtAfterIDCqlQuery, partnerID, from, to)
}

// Returns set of found rows. All select queries should use this function.
// selectTaskInstances is the select template in which the values are inserted
func selectTaskInstances(ctx context.Context, selectQuery string, values ...interface{}) ([]TaskInstance, error) {
	var (
		taskInstance   TaskInstance
		queries        = make([]TaskInstance, 0)
		cassandraQuery = cassandra.QueryCassandra(ctx, selectQuery, values...)
		emptyTime      = time.Time{}
		deviceStatuses = make(map[gocql.UUID]statuses.TaskInstanceStatus)
		endpoints      = make(map[gocql.UUID]statuses.TaskInstanceStatus)
	)

	iterator := cassandraQuery.Iter()
	for iterator.Scan(
		&taskInstance.PartnerID,
		&taskInstance.ID,
		&taskInstance.TaskID,
		&taskInstance.Name,
		&taskInstance.OriginID,
		&taskInstance.StartedAt,
		&taskInstance.LastRunTime,
		&deviceStatuses,
		&taskInstance.FailureCount,
		&taskInstance.SuccessCount,
		&taskInstance.OverallStatus,
		&taskInstance.TriggeredBy,
		&endpoints,
	) {

		//this check can be removed 03/17/2019
		if emptyTime == taskInstance.LastRunTime {
			taskInstance.LastRunTime = taskInstance.StartedAt
		}

		ti := TaskInstance{
			PartnerID:     taskInstance.PartnerID,
			ID:            taskInstance.ID,
			TaskID:        taskInstance.TaskID,
			Name:          taskInstance.Name,
			OriginID:      taskInstance.OriginID,
			StartedAt:     taskInstance.StartedAt,
			LastRunTime:   taskInstance.LastRunTime,
			Statuses:      mergeDeviceStatuses(deviceStatuses, endpoints),
			FailureCount:  taskInstance.FailureCount,
			SuccessCount:  taskInstance.SuccessCount,
			OverallStatus: taskInstance.OverallStatus,
			TriggeredBy:   taskInstance.TriggeredBy,
		}

		ti.PreparePendingStatuses()

		queries = append(queries, ti)
	}

	if err := iterator.Close(); err != nil {
		return nil, fmt.Errorf(errWorkingWithEntities, err)
	}
	return queries, nil
}

func mergeDeviceStatuses(deviceStatuses map[gocql.UUID]statuses.TaskInstanceStatus, endpoints map[gocql.UUID]statuses.TaskInstanceStatus) map[gocql.UUID]statuses.TaskInstanceStatus {
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

// GetMinimalInstanceByID - returns Task Instances minimal info (task_id and name) by id
func (TaskInstanceRepoCassandra) GetMinimalInstanceByID(ctx context.Context, id gocql.UUID) (TaskInstance, error) {
	query := `SELECT
				partner_id,
				task_id,
				origin_id,
				name,
				started_at,
				success_count,
				failure_count
			  FROM task_instances_by_id_mv WHERE id = ?`

	var ti TaskInstance
	params := []interface{}{
		&ti.PartnerID,
		&ti.TaskID,
		&ti.OriginID,
		&ti.Name,
		&ti.StartedAt,
		&ti.SuccessCount,
		&ti.FailureCount,
	}

	q := cassandra.QueryCassandra(ctx, query, id)
	defer q.Release()
	if err := q.Scan(params...); err != nil {
		return ti, errors.Wrapf(err, "can't get minimal task instance by id %v", id)
	}
	ti.ID = id
	return ti, nil
}
