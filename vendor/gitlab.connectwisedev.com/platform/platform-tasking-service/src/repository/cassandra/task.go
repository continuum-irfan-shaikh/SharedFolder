package cassandra

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
)

// NewTask - returns new instance of Task repository
func NewTask(conn cassandra.ISession) *Task {
	return &Task{
		Conn: conn,
	}
}

// Task - task repo representation
type Task struct {
	Conn cassandra.ISession
}

// GetScheduledTasks returns future tasks for ScheduledTasks entity
func (t *Task) GetScheduledTasks(partnerID string) (tasks []entities.ScheduledTasks, err error) {
	now := time.Now().UTC().Truncate(time.Minute)
	queryStmt := `SELECT
                  id,
                  run_time_unix,
				  state,
				  last_task_instance_id,
				  run_time
              FROM task_by_runtime_unix_mv
              WHERE partner_id = ?
              AND external_task = false
              AND run_time_unix >= ?`

	commonFieldsQuery := `SELECT
                              name,
                              description,
                              created_by,
                              created_at,
                              modified_by,
                              modified_at,
                              require_noc_access,
                              type,
	 			              origin_id,
	 			              schedule
                          FROM tasks2
                          WHERE partner_id = ?
                          AND id = ?`

	iter := t.Conn.Query(queryStmt, partnerID, now).Iter()
	tasks = make([]entities.ScheduledTasks, 0)
	tasksByID := make(map[string][]entities.ScheduledTasks)

	var st entities.ScheduledTasks
	var scheduleString string

	params := []interface{}{
		&st.ID,
		&st.RunTimeUTC,
		&st.State,
		&st.LastTaskInstanceID,
		&st.PostponedTime,
	}

	commonParams := []interface{}{
		&st.Name,
		&st.Description,
		&st.CreatedBy,
		&st.CreatedAt,
		&st.ModifiedBy,
		&st.ModifiedAt,
		&st.IsNOC,
		&st.TaskType,
		&st.OriginID,
		&scheduleString,
	}

	for iter.Scan(params...) {
		schedule := apiModels.Schedule{}
		mappedTasks, ok := tasksByID[st.ID]
		if !ok {
			q := t.Conn.Query(commonFieldsQuery, partnerID, st.ID)
			err := q.Scan(commonParams...)
			if err != nil {
				return nil, err
			}

			if err = json.Unmarshal([]byte(scheduleString), &schedule); err != nil {
				return nil, fmt.Errorf("wrong format of schedule string %s: %s", scheduleString, err.Error())
			}
		} else {
			mappedTask := mappedTasks[0]
			st.Name = mappedTask.Name
			st.Description = mappedTask.Description
			st.CreatedBy = mappedTask.CreatedBy
			st.CreatedAt = mappedTask.CreatedAt
			st.ModifiedBy = mappedTask.ModifiedBy
			st.ModifiedAt = mappedTask.ModifiedAt
			st.IsNOC = mappedTask.IsNOC
			st.TaskType = mappedTask.TaskType
			st.OriginID = mappedTask.OriginID
			schedule.Regularity = mappedTask.Regularity
			schedule.TriggerFrames = mappedTask.TriggerFrames
			schedule.TriggerTypes = mappedTask.TriggerTypes
		}

		task := entities.ScheduledTasks{
			ID:                 st.ID,
			Name:               st.Name,
			RunTimeUTC:         st.RunTimeUTC,
			Description:        st.Description,
			CreatedBy:          st.CreatedBy,
			CreatedAt:          st.CreatedAt,
			ModifiedBy:         st.ModifiedBy,
			ModifiedAt:         st.ModifiedAt,
			State:              st.State,
			IsNOC:              st.IsNOC,
			LastTaskInstanceID: st.LastTaskInstanceID,
			TaskType:           t.getTaskType(st.TaskType, st.OriginID),
			Regularity:         schedule.Regularity,
			TriggerFrames:      schedule.TriggerFrames,
			TriggerTypes:       schedule.TriggerTypes,
			PostponedTime:      st.PostponedTime,
		}
		mappedTasks = append(mappedTasks, task)
		tasksByID[task.ID] = mappedTasks
	}

	if err = iter.Close(); err != nil {
		return nil, errors.Wrap(err, "can't get scheduled tasks")
	}

	for _, tasksGroup := range tasksByID {
		tasks = append(tasks, tasksGroup...)
	}

	return tasks, nil
}

// returns task type based on origin ID if needed
func (t *Task) getTaskType(taskType, originID string) string {
	switch originID {
	case entities.CustomBash, entities.CustomCMD, entities.CustomPowershell:
		return models.TaskTypeScript
	}
	// if it's sequence or so
	if taskType != models.TaskTypeScript {
		return taskType
	}
	return models.TaskTypeAction
}

// GetNext - returns next tasks to be done
func (t *Task) GetNext(partnerID string) ([]entities.Task, error) {
	query := `SELECT
                  id,
                  run_time_unix,
                  run_time,
                  managed_endpoint_id,
				  state
              FROM task_by_runtime_unix_mv
              WHERE partner_id = ?
              AND external_task = false
              AND run_time_unix > toUnixTimestamp(now())`

	commonFieldsQuery := `SELECT
                              name
                          FROM tasks2
                          WHERE partner_id = ?
                          AND id = ?`

	q := t.Conn.Query(query, partnerID)
	q.SetConsistency(gocql.One)
	defer q.Release()

	iter := q.Iter()
	tasks := make([]entities.Task, 0)
	tasksByID := make(map[string][]entities.Task)

	var tc entities.Task

	params := []interface{}{
		&tc.ID,
		&tc.RunTimeUTC,
		&tc.PostponedRunTime,
		&tc.ManagedEndpointID,
		&tc.State,
	}

	commonParams := []interface{}{
		&tc.Name,
	}

	for iter.Scan(params...) {
		mappedTasks, ok := tasksByID[tc.ID]
		if !ok {
			q := t.Conn.Query(commonFieldsQuery, partnerID, tc.ID)
			err := q.Scan(commonParams...)
			if err != nil {
				return nil, err
			}
		} else {
			mappedTask := mappedTasks[0]
			tc.Name = mappedTask.Name
		}

		task := entities.Task{
			Name:              tc.Name,
			RunTimeUTC:        tc.RunTimeUTC,
			PostponedRunTime:  tc.PostponedRunTime,
			ManagedEndpointID: tc.ManagedEndpointID,
			State:             tc.State,
		}

		mappedTasks = append(mappedTasks, task)
		tasksByID[tc.ID] = mappedTasks
	}

	if err := iter.Close(); err != nil {
		return nil, errors.Wrap(err, "can't get next task runs")
	}

	for _, tasksGroup := range tasksByID {
		tasks = append(tasks, tasksGroup...)
	}

	return tasks, nil
}

// GetName - returns task name
func (t *Task) GetName(partnerID string, id string) (string, error) {
	query := `SELECT
                  name
              FROM tasks2
              WHERE partner_id = ?
              AND id = ?`

	var task entities.Task
	params := []interface{}{
		&task.Name,
	}

	q := t.Conn.Query(query, partnerID, id)
	q.SetConsistency(gocql.One)
	defer q.Release()

	if err := q.Scan(params...); err != nil && err != gocql.ErrNotFound {
		return "", errors.Wrap(err, "can't get task name")
	}

	return task.Name, nil
}
