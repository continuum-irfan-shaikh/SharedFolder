package cassandra

import (
	"fmt"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/utils/migrate-locations-and-targets/migrate-models"
	"github.com/gocql/gocql"
)

const (
	taskInstancesTblName                         = `task_instances`
	tempTaskInstancesTblName                     = `temp_task_instances`
	selectTaskInstancesCqlQuery                  = `SELECT id, task_id, origin_id, targets, started_at, status FROM ` + taskInstancesTblName + `;`
	selectTempTaskInstancesCqlQuery              = `SELECT id, task_id, origin_id, targets, started_at, status FROM ` + tempTaskInstancesTblName + `;`
	insertTempTaskInstancesQuery                 = `INSERT INTO ` + tempTaskInstancesTblName + ` (id, task_id, origin_id, targets, started_at, status) VALUES (?, ?, ?, ?, ?, ?);`
	insertTaskInstancesQuery                     = `INSERT INTO ` + taskInstancesTblName + ` (id, task_id, origin_id, targets, started_at, status) VALUES (?, ?, ?, ?, ?, ?);`
	dropTaskInstancesTable                       = `DROP TABLE IF EXISTS ` + taskInstancesTblName + `;`
	dropTaskInstancesIndexCqlQuery               = `DROP INDEX IF EXISTS task_instances_targets_idx;`
	dropTaskInstanceByIDMaterializedViewCqlQuery = `DROP MATERIALIZED VIEW IF EXISTS task_instances_by_id;`
	dropTaskInstanceByStartedMVCqlQuery          = `DROP MATERIALIZED VIEW IF EXISTS task_instances_by_started_at;`
	createTaskInstancesIndexCqlQuery             = `CREATE INDEX IF NOT EXISTS task_instances_targets_idx ON ` + taskInstancesTblName + ` (values(targets));`
	createTaskInstancesMaterializedViewCqlQuery  = `CREATE MATERIALIZED VIEW IF NOT EXISTS task_instances_by_id AS
														SELECT id, task_id, started_at, origin_id, status, targets
														FROM ` + taskInstancesTblName +
		` WHERE id IS NOT NULL AND task_id IS NOT NULL AND started_at IS NOT NULL
														PRIMARY KEY (id, task_id, started_at);`
)

func createTaskInstancesTables(taskInstancesTableName string) error {
	createTempTasksTable := `CREATE TABLE IF NOT EXISTS ` + taskInstancesTableName + ` (
										id                  uuid,
										task_id             uuid,
										targets             set<text>,
										origin_id           uuid,
										started_at          timestamp,
										status              int,
										PRIMARY KEY         (task_id, started_at, id))
										WITH CLUSTERING ORDER BY (started_at DESC, id DESC);`

	cassandraQuery, session, err := QueryCassandra(createTempTasksTable)
	if err != nil {
		fmt.Printf("Can't create a table with the name: %s\n", taskInstancesTableName)
		return err
	}
	defer session.Close()

	return cassandraQuery.Exec()
}

func selectTaskInstancesTable(isTempTable bool) (taskInstances []migrateModels.TaskInstance, err error) {
	query := selectTaskInstancesCqlQuery

	if isTempTable {
		query = selectTempTaskInstancesCqlQuery
	}

	cassandraQuery, session, err := QueryCassandra(query)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	taskInstances = make([]migrateModels.TaskInstance, 0)
	var taskInstance migrateModels.TaskInstance

	iter := cassandraQuery.Iter()

	for iter.Scan(
		&taskInstance.ID,
		&taskInstance.TaskID,
		&taskInstance.OriginID,
		&taskInstance.Targets,
		&taskInstance.StartedAt,
		&taskInstance.Status,
	) {
		taskInstances = append(
			taskInstances,
			migrateModels.TaskInstance{
				ID:        taskInstance.ID,
				TaskID:    taskInstance.TaskID,
				OriginID:  taskInstance.OriginID,
				Targets:   taskInstance.Targets,
				StartedAt: taskInstance.StartedAt,
				Status:    taskInstance.Status,
			},
		)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("error while working with found entities: %v", err)
	}

	return taskInstances, err
}

func insertTaskInstancesBatch(tasks []migrateModels.TaskInstance, insertQuery string) error {
	session, err := GetSession()
	if err != nil {
		return fmt.Errorf("can't create cassandra session: %v", err)
	}

	for _, v := range tasks {
		batch := gocql.NewBatch(gocql.LoggedBatch)
		batch.Query(insertQuery,
			v.ID,
			v.OriginID,
			v.TaskID,
			v.Targets,
			v.StartedAt,
			v.Status,
		)

		err = session.ExecuteBatch(batch)
		if err != nil {
			return fmt.Errorf("can't execute batch: %v", err)
		}
	}

	return nil
}

func saveTaskInstances(isTemp bool) (err error) {
	var (
		taskInstances []migrateModels.TaskInstance
		insertQuery   string
	)

	if isTemp {
		taskInstances, err = selectTaskInstancesTable(false)
		if err != nil {
			return fmt.Errorf("can't convert obsolete tasks to new  format: %v", err)
		}
		insertQuery = insertTempTaskInstancesQuery
	} else {
		taskInstances, err = selectTaskInstancesTable(true)
		if err != nil {
			return fmt.Errorf("can't convert obsolete tasks to new  format: %v", err)
		}
		insertQuery = insertTaskInstancesQuery
	}

	err = insertTaskInstancesBatch(taskInstances, insertQuery)

	if err != nil {
		fmt.Printf("Can't fulfil batch for task_instances table: %v", err)
		return err
	}

	return nil
}
