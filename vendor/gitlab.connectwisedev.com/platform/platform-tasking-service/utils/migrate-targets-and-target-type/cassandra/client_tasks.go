package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

// TargetType is used for targets definition
type TargetType int

// These constants describe the type of entities the task was run on
const (
	_ TargetType = iota
	managedEndpoint
	dynamicGroup
)

const (
	selectAllTasksCqlQuery = `SELECT partner_id, id, managed_endpoint_id, external_task, targets FROM tasks`
	updateTaskCqlQuery     = `UPDATE tasks USING TTL 7775984 SET targets = ?, target_type = ? WHERE partner_id = ? AND external_task = ? AND managed_endpoint_id = ? AND id = ?`
)

// TaskFields is a set of Primary Keys of tasks table
type TaskFields struct {
	PartnerID         string
	ID                gocql.UUID
	ManagedEndpointID gocql.UUID
	ExternalTask      bool
	Targets           []string
}

func selectAllTasks() ([]TaskFields, map[gocql.UUID][]string, error) {
	var (
		cassandraQuery             = Session.Query(selectAllTasksCqlQuery)
		taskFields                 TaskFields
		taskFieldsList             []TaskFields
		managedEndpointIDsByTaskID = make(map[gocql.UUID][]string)
		iter                       = cassandraQuery.Iter()
	)
	for iter.Scan(
		&taskFields.PartnerID,
		&taskFields.ID,
		&taskFields.ManagedEndpointID,
		&taskFields.ExternalTask,
		&taskFields.Targets,
	) {
		taskFieldsList = append(
			taskFieldsList,
			TaskFields{
				PartnerID:         taskFields.PartnerID,
				ID:                taskFields.ID,
				ManagedEndpointID: taskFields.ManagedEndpointID,
				ExternalTask:      taskFields.ExternalTask,
				Targets:           taskFields.Targets,
			},
		)
		managedEndpointIDsByTaskID[taskFields.ID] = append(
			managedEndpointIDsByTaskID[taskFields.ID],
			taskFields.ManagedEndpointID.String(),
		)
	}

	if err := iter.Close(); err != nil {
		return nil, nil, errors.Wrapf(err, "can't perform query '%s'", selectAllTasksCqlQuery)
	}

	return taskFieldsList, managedEndpointIDsByTaskID, nil
}

func updateTasks(taskFieldsList []TaskFields, managedEndpointIDsByTaskID map[gocql.UUID][]string) (err error) {
	var targetType TargetType
	for _, task := range taskFieldsList {
		if len(task.Targets) > 0 {
			// Then the task was created for the list of DynamicGroups
			targetType = dynamicGroup
		} else {
			// Then the task was created for the list of ManagedEndpoints
			// So they should be duplicated in the targets field
			task.Targets = managedEndpointIDsByTaskID[task.ID]
			targetType = managedEndpoint
		}
		errExec := Session.Query(updateTaskCqlQuery, task.Targets, targetType, task.PartnerID, task.ExternalTask, task.ManagedEndpointID, task.ID).Exec()
		if errExec != nil {
			err = fmt.Errorf("\ntask:%v, err:%s. %s", task, errExec, err)
		}
	}

	return err
}

// UpdateTasksTable sets external_task to False for all existing tasks
func UpdateTasksTable() error {
	taskFieldsList, managedEndpointIDsByTaskID, err := selectAllTasks()
	if err != nil {
		return err
	}
	if err = updateTasks(taskFieldsList, managedEndpointIDsByTaskID); err != nil {
		return errors.Wrap(err, "error while updating tasks")
	}
	return nil
}
