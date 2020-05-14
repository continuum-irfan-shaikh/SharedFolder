package main

import (
	"io"
	"log"
	"os"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
	"github.com/urfave/cli"
)

const (
	size    = "size"
	host    = "host"
	logName = "_copying_tasks_to_tasks2.log"
)

func main() {
	app := cli.App{}
	app.Name = "Continuum Juno Automation Tasks DB migration tool"
	app.Usage = "This piece of software will copy consistency of tasks table to tasks2 table"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:  size,
			Usage: "Size of data set fetched by one query from Cassandra",
			Value: 100,
		},
		&cli.StringSliceFlag{
			Name:  host,
			Usage: "Cassandra hosts",
			Value: cli.NewStringSlice("127.0.0.1:9042"),
		},
	}

	setUpLogging()

	app.Action = copyData

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("ERROR while running app: %v", err)
	}
	log.Println("INFO: app finished its work")
}

var copyData cli.ActionFunc = func(c *cli.Context) error {
	start := time.Now()
	log.Println("Copying data process started at ", start.Truncate(time.Minute).Format(time.RFC3339))
	defer func() {
		log.Printf("Copying data process finished in %v min \n", time.Now().Sub(start).Minutes())
	}()

	cassandraHosts := c.StringSlice(host)
	session := createCassandraSession(cassandraHosts)
	defer session.Close()

	var partnerID string
	var taskID gocql.UUID
	query := session.Query("SELECT DISTINCT partner_id, id FROM platform_tasking_db.tasks_by_id_mv")
	iterator := query.Iter()
	defer func() {
		if w := iterator.Warnings(); len(w) != 0 {
			log.Printf("WARN: %v", w)
		}
		err := iterator.Close()
		if err != nil {
			log.Printf("ERROR: can't close iterator for query <%s>: %v", query.Statement(), err)
		}
	}()

	selectStmt := `SELECT
                       definition_id, id, name, description, created_at, created_by, partner_id, origin_id, 
                       type, parameters, external_task, result_webhook, require_noc_access,
                       modified_by, modified_at, targets, target_type, schedule, credentials
                   FROM platform_tasking_db.tasks_by_id_mv
                   WHERE partner_id = ? AND id = ? LIMIT 1;`

	insertStmt := `INSERT INTO platform_tasking_db.tasks2
                      (definition_id, id, name, description, created_at, created_by, partner_id, origin_id, 
                       type, parameters, external_task, result_webhook, require_noc_access,
                       modified_by, modified_at, targets, schedule, credentials)
                   VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	selectTargetsStmt := `SELECT target_type, targets
                          FROM platform_tasking_db.targets
                          WHERE partner_id = ? AND task_id = ?;`

	retryPolicy := gocql.ExponentialBackoffRetryPolicy{
		NumRetries: 5,
		Min:        30 * time.Second,
		Max:        5 * time.Minute,
	}

	for iterator.Scan(&partnerID, &taskID) {
		t := models.Task{}
		var scheduleString string
		targets := make(map[int][]string)

		selectParams := []interface{}{
			&t.DefinitionID,
			&t.ID,
			&t.Name,
			&t.Description,
			&t.CreatedAt,
			&t.CreatedBy,
			&t.PartnerID,
			&t.OriginID,
			&t.Type,
			&t.Parameters,
			&t.ExternalTask,
			&t.ResultWebhook,
			&t.IsRequireNOCAccess,
			&t.ModifiedBy,
			&t.ModifiedAt,
			&t.Targets.IDs,
			&t.Targets.Type,
			&scheduleString,
			&t.Credentials,
		}

		insertParams := []interface{}{
			&t.DefinitionID,
			&t.ID,
			&t.Name,
			&t.Description,
			&t.CreatedAt,
			&t.CreatedBy,
			&t.PartnerID,
			&t.OriginID,
			&t.Type,
			&t.Parameters,
			&t.ExternalTask,
			&t.ResultWebhook,
			&t.IsRequireNOCAccess,
			&t.ModifiedBy,
			&t.ModifiedAt,
			&targets,
			&scheduleString,
			&t.Credentials,
		}

		selectQuery := session.Query(selectStmt, partnerID, taskID).PageSize(c.Int(size)).RetryPolicy(&retryPolicy)
		iter := selectQuery.Iter()

		for iter.Scan(selectParams...) {
			if t.Targets.Type == models.ManagedEndpoint {
				q := session.Query(selectTargetsStmt, t.PartnerID, t.ID).RetryPolicy(&retryPolicy)
				err := q.Scan(&t.Targets.Type, &t.Targets.IDs)
				if err != nil {
					log.Printf("ERROR: can't execute query: %v, reason: %v", q, err)
					continue
				}
				if t.Targets.Type != models.ManagedEndpoint {
					log.Printf("ERROR: inconsistent target type: need: 1 but got: %v", t.Targets.Type)
					continue
				}
			}

			if len(t.Targets.IDs) < 1 {
				t.Targets.IDs = make([]string, 0)
			}
			targets[int(t.Targets.Type)] = t.Targets.IDs

			q := session.Query(insertStmt, insertParams...).RetryPolicy(&retryPolicy)
			err := q.Exec()
			if err != nil {
				log.Printf("ERROR: can't execute query: %v, reason: %v", q, err)
			}
		}

		if w := iter.Warnings(); len(w) != 0 {
			log.Printf("WARN: %v", w)
		}
		err := iter.Close()
		if err != nil {
			log.Printf("ERROR: can't close iterator for query <%s>: %v", selectQuery.Statement(), err)
		}
	}

	return nil
}

func createCassandraSession(hosts []string) *gocql.Session {
	client := gocql.NewCluster(hosts...)
	client.ProtoVersion = 4
	client.Consistency = gocql.One
	client.Keyspace = "platform_tasking_db"
	client.NumConns = 20
	client.Timeout = 90 * time.Second

	session, err := client.CreateSession()
	if err != nil {
		log.Fatalf("Cannot create cassandra session: %v", err)
	}
	return session
}

func setUpLogging() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.Create(currentDir + "/" + time.Now().Truncate(time.Minute).Format(time.RFC3339) + logName)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}
