package cassandra

import (
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
)

// NewTaskExecutionHistory - returns new instance of tasks history repository
func NewTaskExecutionHistory(conn cassandra.ISession) *TaskExecutionHistory {
	return &TaskExecutionHistory{
		conn: conn,
	}
}

// TaskExecutionHistory - tasks history cassandra repository
type TaskExecutionHistory struct {
	conn cassandra.ISession
}

// Insert - inserts data about task execution history into database
func (t *TaskExecutionHistory) Insert(task entities.TaskExecHistory) error {
	return t.conn.Query(
		`INSERT INTO tasks_execution_history (
				exec_year,
				exec_month,
				date,
				executed_time, 
				endpoint_id,
				script_name,
				completed_time,
				script_id,
         	    status,
				partner_id,
				site_id,
				machine_name,
				executed_by,
				output
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) `,
		task.ExecYear,
		task.ExecMonth,
		task.ExecDate,
		task.ExecTime,
		task.EndpointID,
		task.ScriptName,
		task.CompletedTime,
		task.ScriptID,
		task.ExecStatus,
		task.PartnerID,
		task.SiteID,
		task.MachineName,
		task.ExecBy,
		task.Output,
	).Exec()
}
