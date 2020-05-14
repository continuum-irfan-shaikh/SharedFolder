package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/urfave/cli"

	taskCFG "gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	taskDB "gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
)

const (
	scriptExecutionResultFields = "task_instance_id, managed_endpoint_id, updated_at, execution_status, std_out, std_err"
	allExecResults              = `SELECT ` + scriptExecutionResultFields + ` FROM script_execution_results`
)

func main() {
	app := cli.NewApp()
	app.Name = "Continuum Juno Automation TaskInstances updating"
	app.Usage = "This piece of software will recalculate statuses for Task instances"
	app.Version = "0.0.1"

	taskCFG.Load()
	taskCFG.Config.CassandraTimeoutSec = 1500

	logger.Load(false)
	taskDB.Load()

	app.Action = UpdateTaskInstancesStatuses
	app.Run(os.Args)
}

func UpdateTaskInstancesStatuses(c *cli.Context) error {
	start := time.Now()
	defer log.Println("Finised in: ", time.Now().Sub(start).Seconds(), "sec")

	ctx := context.Background()
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	logFile, err := os.Create(currentDir + "/status_updating.log." + time.Now().Format(time.RFC3339Nano))
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	updateStatusesForPartner(ctx)

	return nil
}

func updateStatusesForPartner(ctx context.Context) {
	query := taskDB.
		Session.
		Query(allExecResults).
		PageState(nil).
		PageSize(1000)
	iterator := query.Iter()

	for {
		var (
			executionResult             models.ExecutionResult
			taskInstanceID              string
			resultsByTaskInstanceIDsMap = make(map[gocql.UUID][]models.ExecutionResult, 0)
		)

		for iterator.Scan(
			&taskInstanceID,
			&executionResult.ManagedEndpointID,
			&executionResult.UpdatedAt,
			&executionResult.ExecutionStatus,
			&executionResult.StdOut,
			&executionResult.StdErr,
		) {
			tiID, err := gocql.ParseUUID(taskInstanceID)
			if err != nil {
				log.Println("updateStatusesForPartner: cannot parse UUID: ", taskInstanceID)
				continue
			}
			executionResult.TaskInstanceID = tiID

			resultsByTaskInstanceIDsMap[tiID] = append(resultsByTaskInstanceIDsMap[tiID], executionResult)
			executionResult = models.ExecutionResult{}
		}

		for taskInstID, execResults := range resultsByTaskInstanceIDsMap {
			log.Printf("updateStatusesForPartner: Processing TaskInstance [%s] and ExecResult %+v\n\n", taskInstID, execResults)

			err := ProcessExecutionResults(
				ctx,
				taskInstID,
				execResults...,
			)
			if err != nil {
				log.Printf("updateStatusesForPartner: executionResultsUpdate.ProcessExecutionResults err: %v\n\n", err)
			}
		}

		if len(iterator.PageState()) > 0 {
			iterator = query.PageState(iterator.PageState()).Iter()
		} else {
			break
		}
	}

	if err := iterator.Close(); err != nil {
		log.Fatal(err)
	}
}

func ProcessExecutionResults(
	ctx context.Context,
	taskInstanceID gocql.UUID,
	results ...models.ExecutionResult,
) error {
	var (
		executionResults   = make([]models.ExecutionResult, len(results))
		managedEndpointIDs = make([]gocql.UUID, len(results))
	)

	taskInstance, err := getTaskInstance(ctx, taskInstanceID)
	if err != nil {
		return err
	}

	if taskInstance.Statuses == nil {
		taskInstance.Statuses = make(map[gocql.UUID]statuses.TaskInstanceStatus)
	}

	for i, result := range results {
		executionResult := models.ExecutionResult{
			TaskInstanceID:    taskInstanceID,
			ExecutionStatus:   result.ExecutionStatus,
			StdErr:            result.StdErr,
			StdOut:            result.StdOut,
			UpdatedAt:         result.UpdatedAt,
			ManagedEndpointID: result.ManagedEndpointID,
		}
		executionResults[i] = executionResult
		managedEndpointIDs[i] = result.ManagedEndpointID

		if taskInstance.Statuses[result.ManagedEndpointID] == statuses.TaskInstanceSuccess ||
			taskInstance.Statuses[result.ManagedEndpointID] == statuses.TaskInstanceFailed {
			continue // because we've already got results
		}

		taskInstance.Statuses[result.ManagedEndpointID] = result.ExecutionStatus
		if result.ExecutionStatus == statuses.TaskInstanceSuccess {
			taskInstance.SuccessCount++
		} else {
			taskInstance.FailureCount++
		}
	}

	taskInstance.OverallStatus, err = taskInstance.CalculateOverallStatus()
	if err != nil {
		return err
	}

	err = models.TaskInstancePersistenceInstance.Insert(ctx, taskInstance)
	if err != nil {
		return err
	}
	return nil
}

func getTaskInstance(ctx context.Context, taskInstanceID gocql.UUID) (models.TaskInstance, error) {
	taskInstances, err := models.TaskInstancePersistenceInstance.GetByIDs(ctx, taskInstanceID)
	if err != nil {
		return models.TaskInstance{}, err
	}
	if len(taskInstances) == 0 {
		return models.TaskInstance{}, fmt.Errorf("no TaskInstance found by TaskInstanceID %v", taskInstanceID)
	}
	return taskInstances[0], nil
}
