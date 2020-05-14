package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
	"strings"
	"time"
)

// CassandraClient is a component providing access to Cassandra storage
var CassandraClient *gocql.ClusterConfig

const (
	dropTempTasksTable         = `DROP TABLE IF EXISTS ` + tempTasksTblName + `;`
	dropTempTaskInstancesTable = `DROP TABLE IF EXISTS ` + tempTaskInstancesTblName + `;`
)

// Load creates a Cassandra Client and populates it with initial data
func Load(cassandraURL, cassandraKeyspace string, cassandraTimeoutSec int) {
	urls := strings.SplitN(cassandraURL, ",", -1)
	clusterConfig := gocql.NewCluster(urls...)
	clusterConfig.ProtoVersion = 4
	clusterConfig.Consistency = gocql.Quorum
	clusterConfig.Timeout = time.Duration(cassandraTimeoutSec) * time.Second
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

func executeQuery(query string) error {
	cassandraQuery, session, err := QueryCassandra(query)
	if err != nil {
		fmt.Printf("Can't fulfil the query: %s\n", query)
		return err
	}
	defer session.Close()

	return cassandraQuery.Exec()
}

func dropObsoleteTasks() error {
	err := executeQuery(dropIndexCqlQuery)
	if err != nil {
		fmt.Println("can't execute query:", dropIndexCqlQuery)
		return err
	}
	err = executeQuery(dropTaskByRunTimeMaterializedViewCqlQuery)
	if err != nil {
		fmt.Println("can't execute query:", dropTaskByRunTimeMaterializedViewCqlQuery)
		return err
	}
	err = executeQuery(dropTaskByTargetMaterializedViewCqlQuery)
	if err != nil {
		fmt.Println("can't execute query:", dropTaskByTargetMaterializedViewCqlQuery)
		return err
	}
	err = executeQuery(dropTasksTable)
	if err != nil {
		fmt.Println("can't execute query:", dropTasksTable)
		return err
	}
	return nil
}

func transformTasksTable() error {
	err := createTasksTables(tempTasksTblName)
	if err != nil {
		fmt.Println("can't create table with the name:", tempTasksTblName)
		return err
	}

	err = saveTasks(true)
	if err != nil {
		fmt.Println("can't save data to table with the name:", tempTasksTblName)
		return err
	}

	err = dropObsoleteTasks()
	if err != nil {
		fmt.Println("can't drop obsolete tasks data:", err)
		return err
	}
	err = createTasksTables(tasksTblName)
	if err != nil {
		fmt.Println("can't save data to table with the name:", tasksTblName)
		return err
	}

	err = saveTasks(false)
	if err != nil {
		fmt.Println("can't save data to table with the name:", tasksTblName)
		return err
	}
	err = executeQuery(createTaskByRuntimeMaterializedViewCqlQuery)
	if err != nil {
		fmt.Println("can't execute query:", createTaskByRuntimeMaterializedViewCqlQuery)
		return err
	}
	err = executeQuery(createTaskByTargetMaterializedViewCqlQuery)
	if err != nil {
		fmt.Println("can't execute query:", createTaskByTargetMaterializedViewCqlQuery)
		return err
	}

	return nil
}

func dropObsoleteTaskInstances() error {
	err := executeQuery(dropTaskInstancesIndexCqlQuery)
	if err != nil {
		fmt.Println("can't execute query:", dropTaskInstancesIndexCqlQuery)
		return err
	}
	err = executeQuery(dropTaskInstanceByIDMaterializedViewCqlQuery)
	if err != nil {
		fmt.Println("can't execute query:", dropTaskInstanceByIDMaterializedViewCqlQuery)
		return err
	}
	err = executeQuery(dropTaskInstanceByStartedMVCqlQuery)
	if err != nil {
		fmt.Println("can't execute query:", dropTaskInstanceByStartedMVCqlQuery)
		return err
	}
	err = executeQuery(dropTaskInstancesTable)
	if err != nil {
		fmt.Println("can't execute query:", dropTaskInstancesTable)
		return err
	}

	return nil
}

func transformTaskInstancesTable() error {
	err := createTaskInstancesTables(tempTaskInstancesTblName)
	if err != nil {
		fmt.Println("can't create table with the name:", tempTaskInstancesTblName)
		return err
	}
	err = saveTaskInstances(true)
	if err != nil {
		fmt.Println("can't save data to table with the name:", tempTaskInstancesTblName)
		return err
	}

	err = dropObsoleteTaskInstances()
	if err != nil {
		fmt.Println("can't drop obsolete task instances data", err)
		return err
	}

	err = createTaskInstancesTables(taskInstancesTblName)
	if err != nil {
		fmt.Println("can't create table with the name:", taskInstancesTblName)
		return err
	}
	err = executeQuery(createTaskInstancesIndexCqlQuery)
	if err != nil {
		fmt.Println("can't execute query:", createTaskInstancesIndexCqlQuery)
		return err
	}
	err = executeQuery(createTaskInstancesMaterializedViewCqlQuery)
	if err != nil {
		fmt.Println("can't execute query:", createTaskInstancesMaterializedViewCqlQuery)
		return err
	}

	err = saveTaskInstances(false)
	if err != nil {
		fmt.Println("can't save data to table with the name:", taskInstancesTblName)
		return err
	}

	err = executeQuery(dropTempTasksTable)
	if err != nil {
		fmt.Println("can't execute query:", dropTempTasksTable)
		return err
	}
	err = executeQuery(dropTempTaskInstancesTable)
	if err != nil {
		fmt.Println("can't execute query:", dropTempTasksTable)
		return err
	}

	return nil
}

// TransformTaskingDBTables  fulfils transformations with tasks and task_instances tables
func TransformTaskingDBTables() error {
	err := transformTasksTable()
	if err != nil {
		fmt.Println("can't migrate tasks table:")
		return nil
	}
	err = transformTaskInstancesTable()
	if err != nil {
		fmt.Println("can't migrate task_instances table:")
		return nil
	}

	return nil
}
