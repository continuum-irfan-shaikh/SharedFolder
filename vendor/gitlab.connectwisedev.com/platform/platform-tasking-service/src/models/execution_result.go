package models

//go:generate mockgen -destination=../mocks/mocks-gomock/executionResultPersistence_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/models ExecutionResultPersistence

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
)

const (
	neededScriptExecutionResultFields         = "managed_endpoint_id, task_instance_id, updated_at, execution_status, std_out, std_err"
	deleteScriptExecutionResultCQL            = `DELETE FROM script_execution_results WHERE managed_endpoint_id = ? AND task_instance_id = ?`
	deleteScriptExecutionResultMatViewCQL     = `DELETE FROM script_execution_results_by_task_instance_id_mv WHERE managed_endpoint_id = ? AND task_instance_id = ?`
	getByManagedEndpointAndTaskInstanceIDsCQL = selectDb + neededScriptExecutionResultFields + ` FROM script_execution_results WHERE managed_endpoint_id = ? AND task_instance_id IN `
	selectByTargetsAndTaskInstanceIDsCQL      = selectDb + neededScriptExecutionResultFields + ` FROM script_execution_results_by_task_instance_id_mv WHERE task_instance_id IN `
)

// ExecutionResult is defined by UI requirements.
type ExecutionResult struct {
	ManagedEndpointID gocql.UUID                  `json:"managedEndpointId"`
	TaskInstanceID    gocql.UUID                  `json:"taskInstanceId"`
	UpdatedAt         time.Time                   `json:"updatedAt"`
	ExecutionStatus   statuses.TaskInstanceStatus `json:"executionStatus"`
	StdErr            string                      `json:"stdErr"`
	StdOut            string                      `json:"stdOut"`
}

// ExecutionResultPersistence defines interface for the ExecutionResult repository
type ExecutionResultPersistence interface {
	GetByTargetAndTaskInstanceIDs(managedEndpointID gocql.UUID, taskInstanceID ...gocql.UUID) ([]ExecutionResult, error)
	GetByTaskInstanceIDs(taskInstanceIDs []gocql.UUID) ([]ExecutionResult, error)
	Upsert(ctx context.Context, partnerID, taskName string, results ...ExecutionResult) error
	DeleteBatch(ctx context.Context, executionResults []ExecutionResult) error
}

// ExecutionResultRepoCassandra is a realisation of ExecutionResultPersistence interface for Cassandra
type ExecutionResultRepoCassandra struct{}

var (
	// ExecutionResultPersistenceInstance is an instance presented ExecutionResultRepoCassandra
	ExecutionResultPersistenceInstance ExecutionResultPersistence = ExecutionResultRepoCassandra{}
)

// GetByTaskInstanceIDs gets script execution results by Task Instance IDs
func (ExecutionResultRepoCassandra) GetByTaskInstanceIDs(taskInstanceIDs []gocql.UUID) ([]ExecutionResult, error) {
	uuidInterface := common.ConvertUUIDsToInterfaces(taskInstanceIDs)
	cql := selectByTargetsAndTaskInstanceIDsCQL + "(" + common.GetQuestionMarkString(len(taskInstanceIDs)) + ")"
	return selectExecutionResults(cql, uuidInterface...)
}

// GetByTargetAndTaskInstanceIDs returns ExecutionResult found in Cassandra repository
func (ExecutionResultRepoCassandra) GetByTargetAndTaskInstanceIDs(managedEndpointID gocql.UUID, taskInstanceIDs ...gocql.UUID) ([]ExecutionResult, error) {
	uuidInterface := make([]interface{}, 0, len(taskInstanceIDs)+1)
	uuidInterface = append(uuidInterface, managedEndpointID)
	for _, id := range taskInstanceIDs {
		uuidInterface = append(uuidInterface, id)
	}
	cql := getByManagedEndpointAndTaskInstanceIDsCQL + "(" + common.GetQuestionMarkString(len(taskInstanceIDs)) + ")"
	return selectExecutionResults(cql, uuidInterface...)
}

func selectExecutionResults(cql string, values ...interface{}) (executionResults []ExecutionResult, err error) {
	cassandraQuery := cassandra.QueryCassandra(context.Background(), cql, values...)
	var res ExecutionResult
	iter := cassandraQuery.Iter()
	for iter.Scan(
		&res.ManagedEndpointID,
		&res.TaskInstanceID,
		&res.UpdatedAt,
		&res.ExecutionStatus,
		&res.StdOut,
		&res.StdErr,
	) {
		executionResults = append(
			executionResults,
			ExecutionResult{
				ManagedEndpointID: res.ManagedEndpointID,
				TaskInstanceID:    res.TaskInstanceID,
				UpdatedAt:         res.UpdatedAt,
				ExecutionStatus:   res.ExecutionStatus,
				StdOut:            res.StdOut,
				StdErr:            res.StdErr,
			},
		)
	}

	if err = iter.Close(); err != nil {
		return nil, fmt.Errorf("error scanning script_execution_results: %v", err)
	}
	return executionResults, nil
}

const errTotalFmt = "%v; %v"

// Upsert updates/inserts ScriptExecutionResults in Cassandra, updates related
// tasks and taskInstances
func (r ExecutionResultRepoCassandra) Upsert(ctx context.Context, partnerID, taskName string, results ...ExecutionResult) (errTotal error) {
	for _, result := range results {
		batch := cassandra.Session.NewBatch(gocql.UnloggedBatch)

		resultFields := []interface{}{
			result.ManagedEndpointID,
			result.TaskInstanceID,
			result.UpdatedAt,
			result.ExecutionStatus,
			result.StdOut,
			result.StdErr,
			config.Config.DataRetentionIntervalDay * secondsInDay,
		}

		query := `INSERT INTO script_execution_results (
				managed_endpoint_id, 
				task_instance_id, 
				updated_at, 
				execution_status, 
				std_out, 
				std_err
			) VALUES (?, ?, ?, ?, ?, ?) USING TTL ?`

		batch.Query(query, resultFields...)

		queryMatView := `INSERT INTO script_execution_results_by_task_instance_id_mv (
				managed_endpoint_id, 
				task_instance_id, 
				updated_at, 
				execution_status, 
				std_out, 
				std_err
			) VALUES (?, ?, ?, ?, ?, ?) USING TTL ?`

		batch.Query(queryMatView, resultFields...)

		queryLastExecutionTable := `INSERT INTO last_task_executions (
				partner_id, 
				endpoint_id, 
				run_time, 
				name, 
				status)
			  VALUES (?, ?, ?, ?, ?)`

		resultFieldsLastExecution := []interface{}{
			partnerID,
			result.ManagedEndpointID,
			result.UpdatedAt,
			taskName,
			result.ExecutionStatus,
		}

		batch.Query(queryLastExecutionTable, resultFieldsLastExecution...)

		if err := cassandra.Session.ExecuteBatch(batch); err != nil {
			//IF batch too large we have to insert one by one
			err = cassandra.QueryCassandra(ctx, query, resultFields...).Exec()
			if err != nil {
				errTotal = fmt.Errorf(errTotalFmt, errTotal, err)
			}

			err = cassandra.QueryCassandra(ctx, queryMatView, resultFields...).Exec()
			if err != nil {
				errTotal = fmt.Errorf(errTotalFmt, errTotal, err)
			}

			err = cassandra.QueryCassandra(ctx, queryLastExecutionTable, resultFieldsLastExecution...).Exec()
			if err != nil {
				errTotal = fmt.Errorf(errTotalFmt, errTotal, err)
			}
		}

	}
	return
}

// DeleteBatch deletes a batch of ExecutionResults from repository
func (ExecutionResultRepoCassandra) DeleteBatch(ctx context.Context, executionResults []ExecutionResult) error {
	batch := cassandra.Session.NewBatch(gocql.UnloggedBatch)

	for i, e := range executionResults {
		eFields := []interface{}{
			e.ManagedEndpointID,
			e.TaskInstanceID,
		}

		batch.Query(deleteScriptExecutionResultCQL, eFields...)
		batch.Query(deleteScriptExecutionResultMatViewCQL, eFields...)

		// no more than config.Config.CassandraBatchSize in 1 batch or all if it's the last iteration
		if (i+1)%config.Config.CassandraBatchSize == 0 || i+1 == len(executionResults) {
			err := cassandra.Session.ExecuteBatch(batch)
			if err != nil {
				return err
			}

			batch = cassandra.Session.NewBatch(gocql.UnloggedBatch)
		}
	}

	return nil
}
