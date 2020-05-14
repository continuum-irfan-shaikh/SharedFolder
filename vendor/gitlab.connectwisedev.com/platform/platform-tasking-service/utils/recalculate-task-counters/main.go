package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/urfave/cli"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
)

const (
	allPartnerIDs                 = "SELECT DISTINCT partner_id FROM tasks"
	endpointIDsByPartnerID        = "SELECT managed_endpoint_id FROM tasks WHERE partner_id = ?"
	selectTasksByEndpointID       = "SELECT partner_id, external_task, managed_endpoint_id, require_noc_access FROM tasks WHERE partner_id = ? AND external_task = ? AND managed_endpoint_id = ?"
	selectTaskCounterByEndpointID = "SELECT * FROM task_counters WHERE partner_id = ? AND endpoint_id = ?"
	updateTaskCounterByEndpointID = "UPDATE task_counters SET external_tasks = external_tasks + ?, internal_tasks = internal_tasks + ? WHERE partner_id = ? AND endpoint_id = ?"
)

type taskCounter struct {
	PartnerID     string
	EndpointID    gocql.UUID
	ExternalTasks int
	InternalTasks int
}

func main() {
	app := cli.NewApp()
	app.Name = "Continuum Juno Automation TaskCounters updating"
	app.Usage = "This piece of software will recalculate task counters for each endpoint"
	app.Version = "0.0.1"

	setUpLogging()
	log.Println("INFO: Logging was set up successfully")

	config.Load()
	log.Println("INFO: config was loaded successfully")

	config.Config.CassandraTimeoutSec = 1500

	cassandra.Load()
	log.Println("INFO: Cassandra session was created successfully")

	app.Action = Recalculate

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("ERROR: failed to run app, err: %v", err)
	}
	log.Println("INFO: app finished the work successfully")
}

func setUpLogging() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.Create(currentDir + "/" + time.Now().Truncate(time.Minute).Format(time.RFC3339) + "_recalculation.log")
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func Recalculate(c *cli.Context) error {
	start := time.Now()
	log.Println("Started recalculation in ", time.Now().Truncate(time.Minute).Format(time.RFC3339))
	defer log.Printf("Finished recalculation in: %v min \n", time.Now().Sub(start).Minutes())

	partnerIDsMap, err := getPartnerIDsMap()
	if err != nil {
		log.Printf("ERROR: failed to get partnerIDs, err: %v", err)
		return err
	}

	log.Println("INFO: number of partnerIDs: ", len(partnerIDsMap))
	log.Println("INFO: partnerIDs map: ", partnerIDsMap)

	var resultErr error

	for partnerID := range partnerIDsMap {
		err = processEndpointsByPartner(partnerID)
		if err != nil {
			resultErr = fmt.Errorf("%v : %v", resultErr, err)
			log.Printf("ERROR: failed to process endpoints for partnerID(%s), err: %v", partnerID, err)
			continue
		}
		log.Printf("INFO: taskCounters for partnerID(%s) were successfully updated", partnerID)
	}

	return resultErr
}

func getPartnerIDsMap() (partnerIDsMap map[string]struct{}, err error) {
	partnerIDsMap = make(map[string]struct{})
	query := cassandra.Session.Query(allPartnerIDs)
	iterator := query.Iter()

	var task models.Task
	for iterator.Scan(
		&task.PartnerID,
	) {
		partnerIDsMap[task.PartnerID] = struct{}{}
	}

	err = iterator.Close()
	return
}

func processEndpointsByPartner(partnerID string) error {
	log.Println("INFO: started processing endpoints for partnerID ", partnerID)
	endpointIDsMap, err := getEndpointIDsMap(partnerID)
	if err != nil {
		return err
	}

	log.Printf("INFO: number of endpointIDs for partnerID[%s]: %v \n", partnerID, len(endpointIDsMap))
	log.Printf("INFO: endpointIDs map for partnerID[%s]: %v \n", partnerID, endpointIDsMap)

	var resultErr error
	for endpointID := range endpointIDsMap {
		err = processEndpoint(partnerID, endpointID)
		if err != nil {
			resultErr = fmt.Errorf("%v : %v", resultErr, err)
			log.Printf("ERROR: failed to process endpointID(%v) for partnerID(%s), err: %v", endpointID, partnerID, err)
			continue
		}
		log.Printf("INFO: taskCounters for endpointID(%v) for partnerID(%s) were successfully updated", endpointID, partnerID)
	}

	return resultErr
}

func getEndpointIDsMap(partnerID string) (endpointIDsMap map[gocql.UUID]struct{}, err error) {
	endpointIDsMap = make(map[gocql.UUID]struct{})
	query := cassandra.Session.Query(endpointIDsByPartnerID, partnerID)
	iterator := query.Iter()

	var task models.Task
	for iterator.Scan(
		&task.ManagedEndpointID,
	) {
		endpointIDsMap[task.ManagedEndpointID] = struct{}{}
	}

	err = iterator.Close()
	return
}

func processEndpoint(partnerID string, endpointID gocql.UUID) (err error) {
	var externalTasksCounter, internalTasksCounter int

	externalTasks, err := getTasksByEndpointID(partnerID, endpointID, true)
	if err != nil {
		return
	}

	for i := range externalTasks {
		if !externalTasks[i].IsRequireNOCAccess {
			externalTasksCounter++
		}
	}

	internalTasks, err := getTasksByEndpointID(partnerID, endpointID, false)
	if err != nil {
		return
	}

	for i := range internalTasks {
		if !internalTasks[i].IsRequireNOCAccess {
			internalTasksCounter++
		}
	}

	tc, err := getTaskCounterByEndpointID(partnerID, endpointID)
	if err != nil {
		return
	}

	// updating broken taskCounter
	// if for some reason the value of taskCounter for external tasks or for internal tasks is < 0
	// taskCounter value will be updated to 0
	var (
		brokenTCShouldUpdate  bool
		brokenExternalTCDelta int
		brokenInternalTCDelta int
	)

	if tc.ExternalTasks < 0 {
		brokenExternalTCDelta = mathAbsInt(tc.ExternalTasks)
		brokenTCShouldUpdate = true
	}

	if tc.InternalTasks < 0 {
		brokenInternalTCDelta = mathAbsInt(tc.InternalTasks)
		brokenTCShouldUpdate = true
	}

	if brokenTCShouldUpdate {
		err = updateTaskCounter(partnerID, endpointID, brokenExternalTCDelta, brokenInternalTCDelta)
		if err != nil {
			return
		}
	}

	// updating taskCounter according to number of internal and external tasks
	var (
		shouldUpdate       bool
		externalTasksDelta int
		internalTasksDelta int
	)

	if externalTasksCounter != tc.ExternalTasks {
		externalTasksDelta = externalTasksCounter - tc.ExternalTasks
		shouldUpdate = true
	}

	if internalTasksCounter != tc.InternalTasks {
		internalTasksDelta = internalTasksCounter - tc.InternalTasks
		shouldUpdate = true
	}

	if shouldUpdate {
		err = updateTaskCounter(partnerID, endpointID, externalTasksDelta, internalTasksDelta)
	}
	return
}

func getTasksByEndpointID(partnerID string, endpointID gocql.UUID, isExternal bool) (tasks []models.Task, err error) {
	iterator := cassandra.Session.Query(selectTasksByEndpointID, partnerID, isExternal, endpointID).Iter()

	var task models.Task
	for iterator.Scan(
		&task.PartnerID,
		&task.ExternalTask,
		&task.ManagedEndpointID,
		&task.IsRequireNOCAccess,
	) {
		tasks = append(tasks, task)
	}

	err = iterator.Close()
	return
}

func getTaskCounterByEndpointID(partnerID string, endpointID gocql.UUID) (tc taskCounter, err error) {
	iterator := cassandra.Session.Query(selectTaskCounterByEndpointID, partnerID, endpointID).Iter()
	iterator.Scan(
		&tc.PartnerID,
		&tc.EndpointID,
		&tc.ExternalTasks,
		&tc.InternalTasks,
	)

	err = iterator.Close()
	return
}

func updateTaskCounter(partnerID string, endpointID gocql.UUID, externalTasksDelta, internalTasksDelta int) (err error) {
	return cassandra.Session.Query(updateTaskCounterByEndpointID,
		externalTasksDelta, internalTasksDelta, partnerID, endpointID).Exec()
}

func mathAbsInt(i int) int {
	if i < 0 {
		return -i
	}

	return i
}
