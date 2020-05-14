package main

import (
	"io"
	"log"
	"os"
	"time"

	taskCFG "gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	db "gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"github.com/urfave/cli"
)

const (
	cassandraSleepTimeSec = 1
	maxInsertCount        = 1000
)

func main() {
	app := cli.NewApp()
	app.Name = "Tartets migrator"
	app.Usage = "This piece of software will migrate targets from internal tasks to new table but only for MANAGED_ENDPOINT_TYPE"
	app.Version = "231.24"

	taskCFG.Load()
	taskCFG.Config.CassandraTimeoutSec = 60

	logger.Load(false)
	db.Load()

	app.Action = MigrateTargets
	app.Run(os.Args)
}

func MigrateTargets(c *cli.Context) error {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.Create(currentDir + "/migrate_targets.log." + time.Now().Format(time.RFC3339Nano))
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	migrate()
	return nil
}

func migrate() {
	start := time.Now()
	defer log.Println("Finised in: ", time.Now().Sub(start).Seconds(), "sec")

	partners, err := getPartnerIDs()
	if err != nil {
		log.Fatalf("error during retreiving partners: %v", err)
	}
	log.Printf("Found %v partners", len(partners))
	getTasksByPartner(partners)
	log.Print("Finished")
}

type taskData struct {
	taskID     string
	partnerID  string
	targetType models.TargetType
	targets    []string
}

func getTasksByPartner(partners []string) {
	selectQueryInternal := `SELECT managed_endpoint_id, target_type, id
              FROM task_by_runtime_unix_mv
              WHERE partner_id = ?
              AND external_task = false
              AND run_time_unix > toUnixTimestamp(now())`

	selectQueryExternal := `SELECT managed_endpoint_id, target_type, id
              FROM task_by_runtime_unix_mv
              WHERE partner_id = ?
              AND external_task = true
              AND run_time_unix > toUnixTimestamp(now())`

	// internal tasks
	for _, partnerID := range partners {
		internalTasks, err := getTaskByPartner(selectQueryInternal, partnerID)
		if err != nil {
			log.Printf("Can't get tasks for partner %v, reason: %v", partnerID, err)
			continue
		}
		log.Printf("Found %v internal tasks to insert for partner %v", len(internalTasks), partnerID)

		if err = insertTasksTargets(internalTasks); err != nil {
			log.Printf("error during inserting internalTasks tasks: %v", err)
		}
	}

	// external tasks
	for _, partnerID := range partners {
		externalTasks, err := getTaskByPartner(selectQueryExternal, partnerID)
		if err != nil {
			log.Printf("Can't get tasks for partner %v, reason %v", partnerID, err)
			continue
		}

		log.Printf("Found %v external tasks to insert for partner %v", len(externalTasks), partnerID)
		if err = insertTasksTargets(externalTasks); err != nil {
			log.Printf("error during inserting externalTasks tasks: %v", err)
		}
	}
}

/*
CREATE TABLE IF NOT EXISTS targets (
     partner_id           text,
     task_id              uuid,
     target_type          int,
     targets              set<text>,
     PRIMARY KEY          (partner_id, task_id)
);
*/
func getTaskByPartner(query, partnerID string) ([]taskData, error) {
	var (
		iter              = db.Session.Query(query, partnerID).Iter()
		endpointID        string
		taskID            string
		targetType        models.TargetType
		endpointsByTaskID = make(map[string][]string)
	)
	params := []interface{}{
		&endpointID,
		&targetType,
		&taskID,
	}

	for iter.Scan(
		params...,
	) {
		if targetType != models.ManagedEndpoint {
			continue
		}
		endpointsByTaskID[taskID] = append(endpointsByTaskID[taskID], endpointID)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	var tasks []taskData
	for tID, endpoints := range endpointsByTaskID {
		tasks = append(tasks, taskData{
			taskID:     tID,
			partnerID:  partnerID,
			targetType: models.ManagedEndpoint,
			targets:    endpoints,
		})
	}
	return tasks, nil
}

func getPartnerIDs() ([]string, error) {
	selectQuery := `SELECT DISTINCT partner_id FROM tasks`
	iter := db.Session.Query(selectQuery).Iter()

	var partner string
	var partners []string
	for iter.Scan(&partner) {
		partners = append(partners, partner)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	return partners, nil
}

func insertTasksTargets(tasks []taskData) error {
	query := `INSERT INTO targets (partner_id, task_id, target_type, targets) VALUES (?,?,?,?)`
	counter := 0
	for _, t := range tasks {
		if counter == maxInsertCount {
			counter = 0
			// make cassandra sleep for a bit before new cycle of HEAVY LOAD  \*_*/
			time.Sleep(time.Second * cassandraSleepTimeSec)
		}

		params := []interface{}{
			t.partnerID,
			t.taskID,
			t.targetType,
			t.targets,
		}

		if err := db.Session.Query(query, params...).Exec(); err != nil {
			log.Printf("error during Inserting %v", err)
			continue
		}
		counter++
	}
	return nil
}
