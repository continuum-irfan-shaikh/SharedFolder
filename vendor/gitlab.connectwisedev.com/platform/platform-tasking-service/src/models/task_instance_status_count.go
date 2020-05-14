package models

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const (
	taskInstanceStatusFields            = `task_instance_id, success_count, failure_count`
	updateTaskInstanceStatusCountQuery  = `UPDATE task_instances_status_count SET success_count = success_count + ?, failure_count = failure_count + ? WHERE task_instance_id = ?`
	selectTaskInstanceStatusesByIDQuery = `SELECT ` + taskInstanceStatusFields + ` FROM task_instances_status_count WHERE task_instance_id = ?`
	selectTaskInstanceIDsQuery          = `SELECT task_instance_id FROM task_instances_status_count`
	deleteTaskInstanceStatusesByIDQuery = `DELETE FROM task_instances_status_count WHERE task_instance_id IN (%s)`
)

// TaskInstanceStatusCount stores succeeded and failed execution results' count for every particular task_instance_id
type TaskInstanceStatusCount struct {
	TaskInstanceID gocql.UUID
	SuccessCount   int
	FailureCount   int
}

// GetStatusCountsByIDs gets Task Instances Statuses Counts from Cassandra
func (t TaskSummaryRepoCassandra) GetStatusCountsByIDs(ctx context.Context, cache persistency.Cache, taskInstancesMapByID map[gocql.UUID]TaskInstance, taskInstanceIDs []gocql.UUID) (map[gocql.UUID]TaskInstanceStatusCount, error) {
	selectQuery := selectTaskInstanceStatusesByIDQuery
	statusCounts := make(map[gocql.UUID]TaskInstanceStatusCount)
	results := make(chan TaskInstanceStatusCount, config.Config.CassandraConcurrentCallNumber)
	limit := make(chan struct{}, config.Config.CassandraConcurrentCallNumber)
	done := make(chan struct{})
	var resultErr error

	go func() {
		for statusCount := range results {
			statusCounts[statusCount.TaskInstanceID] = statusCount
		}
		done <- struct{}{}
	}()

	var wg sync.WaitGroup
	for _, taskInstanceID := range taskInstanceIDs {
		limit <- struct{}{}
		wg.Add(1)

		go func(taskInstanceID gocql.UUID) {
			defer func() {
				<-limit
				wg.Done()
			}()

			if config.Config.AssetCacheEnabled && cache != nil {
				err := t.getCountsFromCache(taskInstanceID, cache, ctx, selectQuery, results, taskInstancesMapByID)
				if err != nil {
					resultErr = fmt.Errorf("%v : %v", resultErr, err)
					return
				}
				return
			}

			if err := t.getCountsFromDB(ctx, selectQuery, taskInstanceID, results); err != nil {
				resultErr = fmt.Errorf("%v : %v", resultErr, err)
				return
			}
			logger.Log.InfofCtx(ctx, "%v : instance with id [%s] not found", resultErr, taskInstanceID.String())
		}(taskInstanceID)
	}

	wg.Wait()
	close(results)
	<-done

	return statusCounts, resultErr
}

func (t TaskSummaryRepoCassandra) getCountsFromDB(ctx context.Context, selectQuery string, taskInstanceID gocql.UUID, results chan TaskInstanceStatusCount) error {
	statusCounts, err := selectTaskInstanceStatusCount(ctx, selectQuery, taskInstanceID)
	if err != nil {
		return err
	}

	if len(statusCounts) > 0 {
		results <- statusCounts[0]
		return nil
	}
	return err
}

func (t TaskSummaryRepoCassandra) getCountsFromCache(taskInstanceID gocql.UUID, cache persistency.Cache, ctx context.Context, selectQuery string, results chan TaskInstanceStatusCount, taskInstancesMapByID map[gocql.UUID]TaskInstance) (resultErr error) {
	keyForCache := []byte("TKS_STATUS_COUNT_BY_ID_" + taskInstanceID.String())
	statusCountBin, err := cache.Get(keyForCache)
	taskInstanceStatusCount := TaskInstanceStatusCount{}
	if err == nil {
		if err = json.Unmarshal(statusCountBin, &taskInstanceStatusCount); err == nil {
			results <- taskInstanceStatusCount
			return
		}
	}

	statusCountsFromCassandra, err := selectTaskInstanceStatusCount(ctx, selectQuery, taskInstanceID)
	if err != nil {
		return err
	}

	if len(statusCountsFromCassandra) <= 0 {
		return
	}

	results <- statusCountsFromCassandra[0]
	if instance, ok := taskInstancesMapByID[taskInstanceID]; !ok {
		return
	} else if len(instance.Statuses) != statusCountsFromCassandra[0].SuccessCount+statusCountsFromCassandra[0].FailureCount {
		return
	}

	statusCounterBytes, err := json.Marshal(statusCountsFromCassandra[0])
	if err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantMarshall, "couldn't marshal  statusCount  for instance id=%s, err:%s", taskInstanceID.String(), err.Error())
		return
	}

	err = cache.Set(keyForCache, statusCounterBytes, 0)
	if err != nil {
		logger.Log.WarnfCtx(ctx, "couldn't set statusCount  with instanceId=%s,  err: %s", taskInstanceID.String(), err.Error())
	}
	return
}

// UpdateTaskInstanceStatusCount updates execution result counter for particular Task Instance
func (TaskSummaryRepoCassandra) UpdateTaskInstanceStatusCount(ctx context.Context, taskInstanceID gocql.UUID, successStatusCount, failureStatusCount int) error {
	if successStatusCount == 0 && failureStatusCount == 0 {
		return nil
	}
	err := cassandra.QueryCassandra(ctx, updateTaskInstanceStatusCountQuery, successStatusCount, failureStatusCount, taskInstanceID).Exec()
	if err != nil {
		return fmt.Errorf("error while updating taskInstanceStatusCount: %v", err)
	}

	return nil
}

// DeleteOldTaskInstanceCounts removes retained data from Cassandra
func DeleteOldTaskInstanceCounts(ctx context.Context) {
	var (
		query = cassandra.QueryCassandra(ctx, selectTaskInstanceIDsQuery).
			PageState(nil).PageSize(config.Config.DataRetentionRemoveBatchSize)
		iter             = query.Iter()
		taskInstancesIDs = make([]gocql.UUID, 0)
	)

	for {
		var id gocql.UUID
		for iter.Scan(&id) {
			if id.Time().AddDate(0, 0, config.Config.DataRetentionIntervalDay).Before(time.Now().UTC()) {
				taskInstancesIDs = append(taskInstancesIDs, id)
			}
		}
		if len(iter.PageState()) > 0 {
			iter = query.PageState(iter.PageState()).Iter()
		} else {
			break
		}
	}

	// Get Task Instances which are still exist in the DB
	existedTaskInstances, err := TaskInstancePersistenceInstance.GetByIDs(ctx, taskInstancesIDs...)
	if err != nil {
		logger.Log.WarnfCtx(ctx, "DeleteOldTaskInstanceCounts: can't find task instances by IDs (%v), error: %v", taskInstancesIDs, err)
		return
	}

	var existedTaskInstanceIDs = make(map[gocql.UUID]struct{})
	for _, taskInstance := range existedTaskInstances {
		existedTaskInstanceIDs[taskInstance.ID] = struct{}{}
	}

	var (
		deprecatedTaskInstanceIDs = make([]interface{}, len(taskInstancesIDs)-len(existedTaskInstanceIDs))
		i                         = 0
	)

	for _, taskInstanceID := range taskInstancesIDs {
		if _, ok := existedTaskInstanceIDs[taskInstanceID]; !ok {
			deprecatedTaskInstanceIDs[i] = taskInstanceID
			i++
		}
	}

	queryString := fmt.Sprintf(deleteTaskInstanceStatusesByIDQuery, common.GetQuestionMarkString(len(deprecatedTaskInstanceIDs)))
	if err := cassandra.QueryCassandra(ctx, queryString, deprecatedTaskInstanceIDs...).Exec(); err != nil {
		logger.Log.WarnfCtx(ctx, "DeleteOldTaskInstanceCounts: can't delete rows with task_instance_ids (%v) from task_instance_status_count, error: %v", deprecatedTaskInstanceIDs, err)
		return
	}
	logger.Log.InfofCtx(ctx, "DeleteOldTaskInstanceCounts: deleted %d deprecated rows from task_instance_status_count", len(deprecatedTaskInstanceIDs))
}

// Getting a slice of Task Instance Status Counts from Cassandra
func selectTaskInstanceStatusCount(ctx context.Context, selectQuery string, values ...interface{}) ([]TaskInstanceStatusCount, error) {
	cassandraQuery := cassandra.QueryCassandra(ctx, selectQuery, values...)

	var (
		statusCount              TaskInstanceStatusCount
		taskInstanceStatusCounts = make([]TaskInstanceStatusCount, 0)
		iterator                 = cassandraQuery.Iter()
	)

	for iterator.Scan(
		&statusCount.TaskInstanceID,
		&statusCount.SuccessCount,
		&statusCount.FailureCount,
	) {

		taskInstanceStatusCounts = append(taskInstanceStatusCounts,
			TaskInstanceStatusCount{
				TaskInstanceID: statusCount.TaskInstanceID,
				SuccessCount:   statusCount.SuccessCount,
				FailureCount:   statusCount.FailureCount,
			},
		)

	}
	if err := iterator.Close(); err != nil {
		return taskInstanceStatusCounts, fmt.Errorf("error while working with found entities: %v", err)
	}
	return taskInstanceStatusCounts, nil
}
