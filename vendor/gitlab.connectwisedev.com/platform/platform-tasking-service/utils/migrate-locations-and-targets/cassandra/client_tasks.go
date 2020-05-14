package cassandra

import (
	"fmt"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/utils/migrate-locations-and-targets/migrate-models"
	"github.com/gocql/gocql"
)

const (
	tasksTblName                               = "tasks"
	tempTasksTblName                           = "temp_tasks"
	dropTasksTable                             = `DROP TABLE IF EXISTS ` + tasksTblName + `;`
	dropIndexCqlQuery                          = `DROP INDEX IF EXISTS tasks_targets_idx;`
	dropTaskByRunTimeMaterializedViewCqlQuery  = `DROP MATERIALIZED VIEW IF EXISTS tasks_by_runtime;`
	dropTaskByTargetMaterializedViewCqlQuery   = `DROP MATERIALIZED VIEW IF EXISTS tasks_by_target;`
	selectTasksCqlQuery                        = `SELECT id, name, description, task_targets, created_at, created_by, partner_id, origin_id, state, regularity, schedule, run_time, start_run_time, end_run_time, trigger, type, parameters FROM ` + tasksTblName + `;`
	selectTempTasksCqlQuery                    = `SELECT id, name, description, target, created_at, created_by, partner_id, origin_id, state, regularity, schedule, run_time, run_time_unix, start_run_time, end_run_time, trigger, type, parameters, location FROM ` + tempTasksTblName + `;`
	insertTempTaskQuery                        = `INSERT INTO ` + tempTasksTblName + ` (id, name, description, target, created_at, created_by, partner_id, origin_id, state, regularity, schedule, run_time, run_time_unix, start_run_time, end_run_time, trigger, type, parameters, location) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	insertTaskQuery                            = `INSERT INTO ` + tasksTblName + ` (id, name, description, target, created_at, created_by, partner_id, origin_id, state, regularity, schedule, run_time, run_time_unix, start_run_time, end_run_time, trigger, type, parameters, location) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	createTaskByTargetMaterializedViewCqlQuery = `CREATE MATERIALIZED VIEW IF NOT EXISTS tasks_by_target AS
                                                   SELECT * FROM ` + tasksTblName + ` WHERE partner_id IS NOT NULL AND id IS NOT NULL AND target IS NOT NULL
    											   PRIMARY KEY (partner_id, target, id);`
	createTaskByRuntimeMaterializedViewCqlQuery = `CREATE MATERIALIZED VIEW IF NOT EXISTS tasks_by_runtime AS
												   SELECT run_time_unix, id, partner_id, target, created_at, created_by, description,
												   end_run_time, location, name, origin_id, parameters, regularity, run_time, schedule,
                                                   start_run_time, state, trigger, type FROM ` + tasksTblName + ` WHERE run_time_unix IS NOT NULL AND partner_id IS NOT NULL AND id IS NOT NULL AND target IS NOT NULL
                                                   PRIMARY KEY (run_time_unix, id, partner_id, target);`
)

func insertTasksBatch(tasks []migrateModels.NewTask, insertQuery string) (err error) {
	session, err := GetSession()
	if err != nil {
		return fmt.Errorf("can't create cassandra session: %v", err)
	}

	for _, v := range tasks {
		batch := gocql.NewBatch(gocql.LoggedBatch)
		batch.Query(insertQuery,
			v.ID,
			v.Name,
			v.Description,
			v.Target,
			v.CreatedAt,
			v.CreatedBy,
			v.PartnerID,
			v.OriginID,
			v.State,
			v.Regularity,
			v.Schedule,
			v.RunTime,
			v.RunTimeUnix,
			v.StartRunTime,
			v.EndRunTime,
			v.Trigger,
			v.Type,
			v.Parameters,
			v.TaskLocation,
		)

		err = session.ExecuteBatch(batch)

		if err != nil {
			return fmt.Errorf("can't execute batch: %v", err)
		}
	}

	return nil
}

func selectTempTasksTable() (tasks []migrateModels.NewTask, err error) {
	cassandraQuery, session, err := QueryCassandra(selectTempTasksCqlQuery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer session.Close()

	tasks = make([]migrateModels.NewTask, 0)
	var task migrateModels.NewTask

	iter := cassandraQuery.Iter()

	for iter.Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Target,
		&task.CreatedAt,
		&task.CreatedBy,
		&task.PartnerID,
		&task.OriginID,
		&task.State,
		&task.Regularity,
		&task.Schedule,
		&task.RunTime,
		&task.RunTimeUnix,
		&task.StartRunTime,
		&task.EndRunTime,
		&task.Trigger,
		&task.Type,
		&task.Parameters,
		&task.TaskLocation,
	) {
		tasks = append(
			tasks,
			migrateModels.NewTask{
				Task: migrateModels.Task{
					ID:           task.ID,
					Name:         task.Name,
					Description:  task.Description,
					CreatedAt:    task.CreatedAt,
					CreatedBy:    task.CreatedBy,
					PartnerID:    task.PartnerID,
					OriginID:     task.OriginID,
					State:        task.State,
					Regularity:   task.Regularity,
					Schedule:     task.Schedule,
					RunTime:      task.RunTime,
					StartRunTime: task.StartRunTime,
					EndRunTime:   task.EndRunTime,
					Trigger:      task.Trigger,
					Type:         task.Type,
					Parameters:   task.Parameters,
				},
				RunTimeUnix:  task.RunTimeUnix,
				TaskLocation: task.TaskLocation,
				Target:       task.Target,
			},
		)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("error while working with found entities: %v", err)
	}

	return tasks, err
}

func selectTasksTable() (tasks []migrateModels.ObsoleteTask, err error) {
	cassandraQuery, session, err := QueryCassandra(selectTasksCqlQuery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer session.Close()

	tasks = make([]migrateModels.ObsoleteTask, 0)
	var task migrateModels.ObsoleteTask

	iter := cassandraQuery.Iter()

	for iter.Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.TaskTargets,
		&task.CreatedAt,
		&task.CreatedBy,
		&task.PartnerID,
		&task.OriginID,
		&task.State,
		&task.Regularity,
		&task.Schedule,
		&task.RunTime,
		&task.StartRunTime,
		&task.EndRunTime,
		&task.Trigger,
		&task.Type,
		&task.Parameters,
	) {
		tasks = append(
			tasks,
			migrateModels.ObsoleteTask{
				Task: migrateModels.Task{
					ID:           task.ID,
					Name:         task.Name,
					Description:  task.Description,
					CreatedAt:    task.CreatedAt,
					CreatedBy:    task.CreatedBy,
					PartnerID:    task.PartnerID,
					OriginID:     task.OriginID,
					State:        task.State,
					Regularity:   task.Regularity,
					Schedule:     task.Schedule,
					StartRunTime: task.StartRunTime,
					EndRunTime:   task.EndRunTime,
					RunTime:      task.RunTime,
					Trigger:      task.Trigger,
					Type:         task.Type,
					Parameters:   task.Parameters,
				},
				TaskTargets: task.TaskTargets,
			},
		)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("error while working with found entities: %v", err)
	}

	return tasks, nil
}

func convertObsoleteTasksToNew() (convertedTasks []migrateModels.NewTask, err error) {
	tasks, err := selectTasksTable()

	if err != nil {
		fmt.Printf("Can't select obsolete task objects: %v", err)
		return nil, fmt.Errorf("Can't select obsolete task objects: %v", err)
	}

	for _, task := range tasks {
		for target := range task.TaskTargets {
			convertedTask := migrateModels.NewTask{
				Task: migrateModels.Task{ID: task.ID,
					Name:         task.Name,
					Description:  task.Description,
					Schedule:     task.Schedule,
					CreatedAt:    task.CreatedAt,
					CreatedBy:    task.CreatedBy,
					PartnerID:    task.PartnerID,
					OriginID:     task.OriginID,
					Regularity:   task.Regularity,
					RunTime:      task.RunTime,
					StartRunTime: task.StartRunTime,
					EndRunTime:   task.EndRunTime,
					Trigger:      task.Trigger,
					Type:         task.Type,
					Parameters:   task.Parameters,
				},
				Target:      target,
				RunTimeUnix: task.RunTime,
			}

			//checkTaskState
			if !task.TaskTargets[target] {
				convertedTask.State = migrateModels.TaskState(3)
			} else {
				if task.State == migrateModels.TaskStateCompleted {
					convertedTask.State = migrateModels.TaskState(2)
				} else {
					convertedTask.State = migrateModels.TaskState(1)
				}
			}
			convertedTask.TaskLocation = "Europe/London"

			convertedTasks = append(convertedTasks, convertedTask)
		}
	}

	return convertedTasks, nil
}

func saveTasks(isTempTable bool) (err error) {
	var (
		tasks       []migrateModels.NewTask
		insertQuery string
	)

	if isTempTable {
		tasks, err = convertObsoleteTasksToNew()
		if err != nil {
			return fmt.Errorf("can't convert obsolete tasks to new  format: %v", err)
		}
		insertQuery = insertTempTaskQuery
	} else {
		tasks, err = selectTempTasksTable()
		if err != nil {
			return fmt.Errorf("can't select tasks from table with name: %v", tempTasksTblName)
		}
		insertQuery = insertTaskQuery
	}

	err = insertTasksBatch(tasks, insertQuery)

	if err != nil {
		fmt.Printf("Can't fulfil batch for tasks table: %v", err)
		return err
	}

	return nil
}

func createTasksTables(tasksTableName string) error {
	createTempTasksTable := `CREATE TABLE IF NOT EXISTS ` + tasksTableName + ` (
										id uuid,
										partner_id text,
										target text,
										created_at timestamp,
										created_by text,
										description text,
										end_run_time timestamp,
										location text,
										name text,
										origin_id uuid,
										parameters text,
										regularity int,
										run_time timestamp,
										run_time_unix timestamp,
										schedule text,
										start_run_time timestamp,
										state int,
										trigger text,
										type text,
										PRIMARY KEY (partner_id, id, target)
									);`

	cassandraQuery, session, err := QueryCassandra(createTempTasksTable)
	if err != nil {
		fmt.Printf("Can't create a table with the name: %s\n", tasksTableName)
		return err
	}
	defer session.Close()

	return cassandraQuery.Exec()
}
