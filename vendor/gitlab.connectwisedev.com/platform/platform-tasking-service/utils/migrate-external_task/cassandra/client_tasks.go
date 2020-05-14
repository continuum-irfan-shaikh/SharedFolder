package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
)

const (
	selectAllTasksCqlQuery             = `SELECT partner_id, id, target FROM tasks`
	updateTasksSetExternalTaskCqlQuery = `UPDATE tasks SET external_task = false WHERE partner_id = ? AND id = ? AND target = ?`
)

// TaskPrimaryKey is a set of Primary Keys of tasks table
type TaskPrimaryKey struct {
	PartnerID string
	ID        gocql.UUID
	Target    string
}

func selectAllTasks() (taskPrimaryKeys []TaskPrimaryKey, err error) {
	cassandraQuery := Session.Query(selectAllTasksCqlQuery)

	var taskPrimaryKey TaskPrimaryKey
	iter := cassandraQuery.Iter()
	for iter.Scan(
		&taskPrimaryKey.PartnerID,
		&taskPrimaryKey.ID,
		&taskPrimaryKey.Target,
	) {
		taskPrimaryKeys = append(
			taskPrimaryKeys,
			TaskPrimaryKey{
				PartnerID: taskPrimaryKey.PartnerID,
				ID:        taskPrimaryKey.ID,
				Target:    taskPrimaryKey.Target,
			},
		)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("can't perform query '%s'. Error: %v", selectAllTasksCqlQuery, err)
	}

	return taskPrimaryKeys, nil
}

func updateTasks(tasksPrimaryKeys []TaskPrimaryKey) (err error) {
	for _, task := range tasksPrimaryKeys {
		errExec := Session.Query(updateTasksSetExternalTaskCqlQuery, task.PartnerID, task.ID, task.Target).Exec()
		if errExec != nil {
			err = fmt.Errorf("\ntask:%v, err:%s. %s", task, errExec, err)
		}
	}

	return err
}

// UpdateTasksTable sets external_task to False for all existing tasks
func UpdateTasksTable() error {
	tasksPrimaryKeys, err := selectAllTasks()
	if err != nil {
		return err
	}
	if err = updateTasks(tasksPrimaryKeys); err != nil {
		return fmt.Errorf("error while updating tasks: %s", err)
	}
	return nil
}
