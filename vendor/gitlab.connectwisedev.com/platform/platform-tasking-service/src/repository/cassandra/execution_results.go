package cassandra

import (
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

// NewScriptExecutionResults - returns new instance of ScriptExecutionResults repository
func NewScriptExecutionResults(conn cassandra.ISession) *ScriptExecutionResults {
	return &ScriptExecutionResults{
		conn: conn,
	}
}

// ScriptExecutionResults - ScriptExecutionResults repo representation
type ScriptExecutionResults struct {
	conn cassandra.ISession
}

// GetLastResultByEndpointID - returns last script execution result for endpoint
func (s *ScriptExecutionResults) GetLastResultByEndpointID(endpointID string) (res entities.ExecutionResult, err error) {
	query := `SELECT 
				task_instance_id,
				execution_status,
				updated_at 
			FROM script_execution_results WHERE managed_endpoint_id = ? LIMIT 1;`
	params := []interface{}{&res.TaskInstanceID, &res.ExecutionStatus, &res.UpdatedAt}

	q := s.conn.Query(query, endpointID)
	q.SetConsistency(gocql.One)
	defer q.Release()

	if err := q.Scan(params...); err != nil && err != gocql.ErrNotFound {
		return res, errors.Wrapf(err, "can't get script execution results by endpoint id %v", endpointID)
	}
	return res, nil
}

// GetLastExecutions - returns last script execution result for endpoint
func (s *ScriptExecutionResults) GetLastExecutions(partnerID string, endpointIDs map[string]struct{}) (execs []entities.LastExecution, err error) {
	query := `SELECT 
				endpoint_id,
				run_time,
				name,
				status 
			FROM last_task_executions WHERE partner_id = ? AND endpoint_id IN ?;`
	execs = make([]entities.LastExecution, 0)
	e := entities.LastExecution{}
	params := []interface{}{&e.EndpointID, &e.RunTime, &e.Name, &e.Status}

	ids := make([]string, 0)
	for k := range endpointIDs {
		ids = append(ids, k)
	}

	q := s.conn.Query(query, partnerID, ids)
	q.SetConsistency(gocql.One)
	defer q.Release()

	iter := q.Iter()
	for iter.Scan(params...) {
		exec := entities.LastExecution{
			EndpointID: e.EndpointID,
			RunTime:    e.RunTime,
			Name:       e.Name,
			Status:     e.Status,
		}
		execs = append(execs, exec)
	}

	if err = iter.Close(); err != nil {
		return execs, errors.Wrapf(err, "can't get last executions for partner %v", partnerID)
	}
	return execs, nil
}
