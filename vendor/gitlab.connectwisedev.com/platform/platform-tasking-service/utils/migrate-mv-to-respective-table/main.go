package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/utils/migrate-mv-to-respective-table/entities"
	"github.com/gocql/gocql"
	"github.com/urfave/cli"
)

const (
	mode  = "mode"
	table = "table"

	csvMode                = "csvMode"
	singleInsertMode       = "insertMode"
	scriptExecResultsTable = "script_execution_results"
	taskInstancesTable     = "task_instances"
	size                   = "size"

	logName = "_copying_materialised_views.log"
)

var payloads map[string]map[string]func(c *cli.Context)

func main() {
	app := cli.NewApp()
	app.Name = "Continuum Juno Automation Tasks DB migration tool"
	app.Usage = "This piece of software will copy consistency of defined materialised views into respective tables with the same structure"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     mode,
			Usage:    "Mode to process migration activities. Available parameters: " + csvMode + " or " + singleInsertMode,
			Required: true,
		},
		cli.StringFlag{
			Name:     table,
			Usage:    "Target table data of which should be migrated. Available parameters: " + scriptExecResultsTable + " or " + taskInstancesTable,
			Required: true,
		},
		cli.IntFlag{
			Name:  size,
			Usage: "Size of data set fetched by one query from Cassandra",
			Value: 100,
		},
	}

	createPayload()

	setUpLogging()

	app.Action = copyDataFromMVToTable

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("ERROR while running app: %v", err)
	}
	log.Println("INFO: app finished its work")
}

func createPayload() {
	payloads = make(map[string]map[string]func(c *cli.Context))
	payloads[csvMode] = map[string]func(c *cli.Context){
		scriptExecResultsTable: func(c *cli.Context) {
			cmdTmpl := exec.Command("cqlsh", "--cqlshrc", "./cqlshrc")

			cmd := *cmdTmpl
			createDumpPayload := fmt.Sprintf("COPY platform_tasking_db.script_execution_results" +
				"(managed_endpoint_id, task_instance_id, execution_status, std_err, std_out, updated_at)" +
				"TO 'script_execution_results_dump.csv'")
			cmd.Args = append(cmd.Args, "-e", createDumpPayload)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.Fatalf("ERROR: can't execute dump creation command: %v", err)
			}

			cmd = *cmdTmpl
			putDataFromDumpPayload := fmt.Sprintf("COPY platform_tasking_db.script_execution_results_by_task_instance_id_mv" +
				"(managed_endpoint_id, task_instance_id, execution_status, std_err, std_out, updated_at)" +
				"FROM 'script_execution_results_dump.csv'")
			cmd.Args = append(cmd.Args, "-e", putDataFromDumpPayload)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				log.Printf("ERROR: can't execute copying command to task_instances_by_id_mv: %v", err)
			}
		},

		taskInstancesTable: func(c *cli.Context) {
			cmdTmpl := exec.Command("cqlsh", "--cqlshrc", "./cqlshrc")

			cmd := *cmdTmpl
			createDumpPayload := "COPY platform_tasking_db.task_instances" +
				"(task_id, started_at, id, device_statuses, failure_count, last_run_time, name, origin_id, overall_status, partner_id, status, success_count, targets, triggered_by)" +
				"TO 'task_instances_dump.csv';"
			cmd.Args = append(cmd.Args, "-e", createDumpPayload)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.Fatalf("ERROR: can't execute dump creation command: %v", err)
			}

			cmd = *cmdTmpl
			putDataFromDumpPayload := "COPY platform_tasking_db.task_instances_started_at_mv" +
				"(task_id, started_at, id, device_statuses, failure_count, last_run_time, name, origin_id, overall_status, partner_id, status, success_count, targets, triggered_by)" +
				"FROM 'task_instances_dump.csv';"
			cmd.Args = append(cmd.Args, "-e", putDataFromDumpPayload)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				log.Printf("ERROR: can't execute copying command to task_instances_started_at_mv: %v", err)
			}

			cmd = *cmdTmpl
			putDataFromDumpPayload = fmt.Sprintf("COPY platform_tasking_db.task_instances_by_id_mv" +
				"(task_id, started_at, id, device_statuses, failure_count, last_run_time, name, origin_id, overall_status, partner_id, status, success_count, targets, triggered_by)" +
				"FROM 'task_instances_dump.csv'")
			cmd.Args = append(cmd.Args, "-e", putDataFromDumpPayload)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				log.Printf("ERROR: can't execute copying command to task_instances_by_id_mv: %v", err)
			}
		},
	}
	payloads[singleInsertMode] = map[string]func(c *cli.Context){
		scriptExecResultsTable: func(c *cli.Context) {
			s := createCassandraSession()
			defer s.Close()

			selectSTMT := `SELECT
                               managed_endpoint_id, task_instance_id, execution_status, std_err, std_out, updated_at
                           FROM platform_tasking_db.script_execution_results;`

			insertSTMT := `INSERT INTO platform_tasking_db.script_execution_results_by_task_instance_id_mv
                               (managed_endpoint_id, task_instance_id, execution_status, std_err, std_out, updated_at)
                           VALUES (?, ?, ?, ?, ?, ?);`

			retryPolicy := gocql.ExponentialBackoffRetryPolicy{
				NumRetries: 3,
				Min:        30 * time.Second,
				Max:        2 * time.Minute,
			}

			q := s.Query(selectSTMT).PageSize(c.Int(size)).RetryPolicy(&retryPolicy)
			iter := q.Iter()
			er := entities.ExecutionResult{}
			params := []interface{}{
				&er.ManagedEndpointID,
				&er.TaskInstanceID,
				&er.ExecutionStatus,
				&er.StdErr,
				&er.StdOut,
				&er.UpdatedAt,
			}

			for iter.Scan(params...) {
				q := s.Query(insertSTMT, params...).RetryPolicy(&retryPolicy)
				err := q.Exec()
				if err != nil {
					log.Printf("ERROR: can't execute query: %v, reason: %v", q, err)
				}
			}
		},

		taskInstancesTable: func(c *cli.Context) {
			s := createCassandraSession()
			defer s.Close()

			partnersIDs := make([]string, 0)
			var partnerID string
			iterator := s.Query("SELECT DISTINCT partner_id FROM platform_tasking_db.task_instances_started_at").Iter()

			for iterator.Scan(&partnerID) {
				partnersIDs = append(partnersIDs, partnerID)
			}
			if err := iterator.Close(); err != nil {
				log.Printf("ERROR: can't close iterator: %v", err)
			}

			selectSTMT := `SELECT
                               task_id, started_at, id,
                               device_statuses, failure_count, last_run_time,
                               name, origin_id, overall_status, partner_id,
                               status, success_count, targets,
                               triggered_by
                           FROM platform_tasking_db.task_instances_started_at
                           WHERE partner_id = ?;`

			insertStartedAtSTMT := `INSERT INTO platform_tasking_db.task_instances_started_at_mv
                                        (task_id, started_at, id,
                                         device_statuses, failure_count, last_run_time,
                                         name, origin_id, overall_status, partner_id,
                                         status, success_count, targets,
                                         triggered_by)
                                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

			insertByIDSTMT := `INSERT INTO platform_tasking_db.task_instances_by_id_mv
                                   (task_id, started_at, id,
                                    device_statuses, failure_count, last_run_time,
                                    name, origin_id, overall_status, partner_id,
                                    status, success_count, targets,
                                    triggered_by)
                               VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

			for _, pID := range partnersIDs {
				retryPolicy := gocql.ExponentialBackoffRetryPolicy{
					NumRetries: 5,
					Min:        30 * time.Second,
					Max:        5 * time.Minute,
				}

				q := s.Query(selectSTMT, pID).PageSize(c.Int(size)).RetryPolicy(&retryPolicy)
				iter := q.Iter()
				ti := entities.TaskInstance{}
				params := []interface{}{
					&ti.TaskID,
					&ti.StartedAt,
					&ti.ID,
					&ti.Statuses,
					&ti.FailureCount,
					&ti.LastRunTime,
					&ti.Name,
					&ti.OriginID,
					&ti.OverallStatus,
					&ti.PartnerID,
					&ti.Status,
					&ti.SuccessCount,
					&ti.Targets,
					&ti.TriggeredBy,
				}

				for iter.Scan(params...) {
					q := s.Query(insertStartedAtSTMT, params...).RetryPolicy(&retryPolicy)
					err := q.Exec()
					if err != nil {
						log.Printf("ERROR: can't execute query: %v, reason: %v", q, err)
					}

					q = s.Query(insertByIDSTMT, params...).RetryPolicy(&retryPolicy)
					err = q.Exec()
					if err != nil {
						log.Printf("ERROR: can't execute query: %v, reason: %v", q, err)
					}
				}
			}
		},
	}
}

var copyDataFromMVToTable cli.ActionFunc = func(c *cli.Context) error {
	start := time.Now()
	log.Println("Copying data process started at ", start.Truncate(time.Minute).Format(time.RFC3339))
	defer func() {
		log.Printf("Copying data process finished in %v min \n", time.Now().Sub(start).Minutes())
	}()

	payload := payloads[c.String(mode)][c.String(table)]
	if payload == nil {
		return errors.New("there are no such function to execute")
	}
	payload(c)

	return nil
}

func createCassandraSession() *gocql.Session {
	client := gocql.NewCluster("127.0.0.1:9042")
	client.ProtoVersion = 4
	client.Consistency = gocql.Quorum
	client.Keyspace = "platform_tasking_db"
	client.NumConns = 20
	client.Timeout = 15 * time.Second

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
