package cassandra

import (
	"fmt"
	"strings"

	"github.com/gocql/gocql"
)

//Target is a component providing the info about enables/disabled targets for a task
type Target struct {
	ID     gocql.UUID `cql:"id"`
	Active bool       `cql:"active"`
}

//TaskTargetsToAdd is component which is used to create new structure of task targets
type TaskTargetsToAdd struct {
	ID        gocql.UUID      `cql:"id"`
	PartnerID string          `cql:"partner_id"`
	Targets   map[string]bool `cql:"task_targets"`
}

//TaskTargetsToDrop is component which is used to read info about targets and move to new structure
type TaskTargetsToDrop struct {
	ID        gocql.UUID `cql:"id"`
	PartnerID string     `cql:"partner_id"`
	Targets   []Target   `cql:"targets"`
}

// CassandraClient is a component providing access to Cassandra storage
var CassandraClient *gocql.ClusterConfig

const (
	// CassandraKeyspace defines a Keyspace in Cassandra where Tasks are stored
	dropTargetsCqlQuery            = `ALTER TABLE tasks DROP targets;`
	addTaskTargetsCqlQuery         = `ALTER TABLE tasks ADD task_targets map<text, boolean>;`
	selectTasksCqlQuery            = `SELECT id, partner_id, targets FROM tasks;`
	updateTaskTargetsCqlQuery      = `UPDATE tasks SET task_targets = ? WHERE partner_id = ? AND id = ?;`
	dropIndexCqlQuery              = `DROP INDEX IF EXISTS tasks_targets_idx;`
	createIndexCqlQuery            = `CREATE INDEX IF NOT EXISTS tasks_targets_idx ON tasks (keys(task_targets));`
	dropMaterializedViewCqlQuery   = `DROP MATERIALIZED VIEW IF EXISTS tasks_by_runtime;`
	createMaterializedViewCqlQuery = `CREATE MATERIALIZED VIEW IF NOT EXISTS tasks_by_runtime AS
	                                  SELECT id, name, description, task_targets, schedule, created_at, created_by, partner_id, origin_id, started_at,
                                      state, regularity, run_time, start_run_time, end_run_time, trigger, type, last_run_at, parameters
                                      FROM tasks
                                      WHERE run_time IS NOT NULL AND partner_id IS NOT NULL AND id IS NOT NULL
                                      PRIMARY KEY (run_time, id, partner_id);`
)

// Load creates a Cassandra Client and populates it with initial data
func Load(cassandraURL, cassandraKeyspace string) {
	urls := strings.SplitN(cassandraURL, ",", -1)
	clusterConfig := gocql.NewCluster(urls...)
	clusterConfig.ProtoVersion = 4
	clusterConfig.Consistency = gocql.All
	CassandraClient = clusterConfig
	CassandraClient.Keyspace = cassandraKeyspace
}

// GetSession - returns cassandra session for Tasking Key Space.
func GetSession() (*gocql.Session, error) {
	return CassandraClient.CreateSession()
}

// QueryCassandra creates a session with cassandra and return a point on a query
func QueryCassandra(cql string, values ...interface{}) (*gocql.Query, *gocql.Session, error) {
	session, err := GetSession()
	if err != nil {
		return nil, nil, fmt.Errorf("can't create cassandra session: %v", err)
	}

	cassandraQuery := session.Query(cql, values...)

	return cassandraQuery, session, nil
}

func selectTaskTargets() (queries []TaskTargetsToDrop, err error) {
	cassandraQuery, session, err := QueryCassandra(selectTasksCqlQuery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer session.Close()

	queries = make([]TaskTargetsToDrop, 0)
	var taskTargetsToDrop TaskTargetsToDrop

	iter := cassandraQuery.Iter()

	for iter.Scan(
		&taskTargetsToDrop.ID,
		&taskTargetsToDrop.PartnerID,
		&taskTargetsToDrop.Targets,
	) {
		queries = append(
			queries,
			TaskTargetsToDrop{
				ID:        taskTargetsToDrop.ID,
				PartnerID: taskTargetsToDrop.PartnerID,
				Targets:   taskTargetsToDrop.Targets,
			},
		)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("error while working with found entities: %v", err)
	}

	return queries, err
}

func updateTaskTableCqlQuery(query string) error {
	cassandraQuery, session, err := QueryCassandra(query)
	if err != nil {
		fmt.Printf("Can't fulfil the query: %s\n", query)
		return err
	}
	defer session.Close()

	return cassandraQuery.Exec()
}

func updateTaskTargetsColumn(taskTarget TaskTargetsToAdd) error {
	cassandraQuery, session, err := QueryCassandra(updateTaskTargetsCqlQuery, taskTarget.Targets, taskTarget.PartnerID, taskTarget.ID)
	if err != nil {
		fmt.Printf("Can't update the task with ID:%s. Err: %s\n", taskTarget.ID, err)
		return err
	}
	defer session.Close()

	return cassandraQuery.Exec()
}

func convertTargetsDataToMap() error {
	queries, err := selectTaskTargets()
	if err != nil {
		fmt.Printf("Error while tasks reading. Err: %s\n", err)
		return err
	}

	err = updateTaskTableCqlQuery(addTaskTargetsCqlQuery)
	if err != nil {
		fmt.Printf("Can't add task targets column to tasks table. Err: %s\n", err)
		return err
	}

	for _, value := range queries {
		targetsMap := make(map[string]bool)

		for _, target := range value.Targets {
			targetsMap[target.ID.String()] = target.Active
		}

		taskTargetsToAdd := TaskTargetsToAdd{value.ID, value.PartnerID, targetsMap}

		err := updateTaskTargetsColumn(taskTargetsToAdd)
		if err != nil {
			fmt.Printf("Can't update task with ID:%s. Err: %s", taskTargetsToAdd.ID, err)
			return err
		}
	}

	return nil
}

// TransformTasksTable  fulfils transformations with targets field. Targets should be map.
func TransformTasksTable() error {
	if err := convertTargetsDataToMap(); err != nil {
		fmt.Printf("Error while targets converting :%s\n", err)
		return err
	}

	updateTaskTableCqlQuery(dropIndexCqlQuery)
	updateTaskTableCqlQuery(dropMaterializedViewCqlQuery)
	updateTaskTableCqlQuery(dropTargetsCqlQuery)
	updateTaskTableCqlQuery(createIndexCqlQuery)
	updateTaskTableCqlQuery(createMaterializedViewCqlQuery)

	return nil
}
