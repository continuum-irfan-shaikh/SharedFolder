package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/urfave/cli"
)

const (
	csvFileName         = "tasks_copy.csv"
	transformedFileName = "transformed_tasks.csv"

	numProcFlag      = "numProc"
	pageTimeoutFlag  = "pageTimeout"
	maxAttemptsFlag  = "maxAttempts"
	maxBatchSizeFlag = "maxBatchSize"
	pageSizeFlag     = "pageSize"
)

func main() {
	app := cli.NewApp()
	app.Name = "Continuum Juno Automation Tasks migration tool"
	app.Usage = "This piece of software will copy tasks from tasks table to tables that emulate behaviour of materialized views"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  numProcFlag,
			Value: 16,
			Usage: "Number of worker processes. Maximum value is 16. Default value: -1.",
		},
		cli.IntFlag{
			Name:  pageTimeoutFlag,
			Value: 10,
			Usage: "Page timeout for fetching results. Default value: 10.",
		},
		cli.IntFlag{
			Name:  maxAttemptsFlag,
			Value: 5,
			Usage: "Maximum number of attempts for errors. Default value: 5.",
		},
		cli.IntFlag{
			Name:  maxBatchSizeFlag,
			Value: 20,
			Usage: "Maximum size of an import batch. Default value:20",
		},
		cli.IntFlag{
			Name:  pageSizeFlag,
			Value: 1000,
			Usage: "Page size for fetching results. Default value: 1000.",
		},
	}

	setUpLogging()
	log.Println("INFO: Logging was set up successfully")

	app.Action = CopyTasks

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("ERROR: failed to run app, err: %v", err)
	}
	log.Println("INFO: app finished the work successfully")

}

func CopyTasks(c *cli.Context) error {
	start := time.Now()
	log.Println("Started copy in ", time.Now().Truncate(time.Minute).Format(time.RFC3339))
	defer log.Printf("Finished copy in: %v min \n", time.Now().Sub(start).Minutes())

	cmd := exec.Command("cqlsh", "-e", fmt.Sprintf("use platform_tasking_db; COPY tasks (run_time_unix, definition_id, iteration_ids, id, name, description,managed_endpoint_id, created_at, created_by, partner_id, origin_id,state, trigger, type, parameters, external_task, result_webhook, last_task_instance_id, require_noc_access, modified_by, modified_at, original_next_run_time, run_time, targets, target_type, schedule) TO 'tasks_copy.csv' WITH HEADER = TRUE AND NUMPROCESSES = %v AND PAGETIMEOUT = %v AND MAXATTEMPTS = %v AND PAGESIZE = %v;", c.Int(numProcFlag), c.Int(pageTimeoutFlag), c.Int(maxAttemptsFlag), c.Int(pageSizeFlag)))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("ERROR: can not execute command copy to, err: %v", err)
	}

	result := make(chan string, 100)
	done := make(chan struct{})
	transformedFile, err := os.Create(transformedFileName)
	if err != nil {
		log.Printf("ERROR: can not create file, err: %v", err)
	}
	defer transformedFile.Close()

	writer := bufio.NewWriter(transformedFile)

	go func() {
		for line := range result {
			writer.WriteString(line)
		}
		writer.Flush()
		done <- struct{}{}

	}()

	file, err := os.Open(csvFileName)
	if err != nil {
		log.Printf("ERROR: can not open file, err: %v", err)
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Printf("ERROR: error during reading the file: %v", err)
				return err
			}
		}

		if strings.HasPrefix(line, ",") {
			line = "0001-01-01T00:00:00Z" + line
		}

		result <- line
	}

	close(result)
	<-done

	cmd = exec.Command("cqlsh", "-e", fmt.Sprintf("use platform_tasking_db; COPY task_by_runtime_unix_mv (run_time_unix, definition_id, iteration_ids, id, name, description,managed_endpoint_id, created_at, created_by, partner_id, origin_id,state, trigger, type, parameters, external_task, result_webhook, last_task_instance_id, require_noc_access, modified_by, modified_at, original_next_run_time, run_time, targets, target_type, schedule) FROM 'transformed_tasks.csv' WITH HEADER = TRUE AND MAXBATCHSIZE = %v;", c.Int(maxBatchSizeFlag)))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("ERROR: can not execute command copy from for task_by_runtime_unix_mv table, err: %v", err)
	}

	cmd = exec.Command("cqlsh", "-e", fmt.Sprintf("use platform_tasking_db; COPY tasks_order_by_last_task_instance_id_mv (run_time_unix, definition_id, iteration_ids, id, name, description,managed_endpoint_id, created_at, created_by, partner_id, origin_id,state, trigger, type, parameters, external_task, result_webhook, last_task_instance_id, require_noc_access, modified_by, modified_at, original_next_run_time, run_time, targets, target_type, schedule) FROM 'transformed_tasks.csv' WITH HEADER = TRUE AND MAXBATCHSIZE = %v;", c.Int(maxBatchSizeFlag)))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("ERROR: can not execute command copy from for tasks_order_by_last_task_instance_id_mv table, err: %v", err)
	}

	cmd = exec.Command("cqlsh", "-e", fmt.Sprintf("use platform_tasking_db; COPY tasks_by_runtime_mv (run_time_unix, definition_id, iteration_ids, id, name, description,managed_endpoint_id, created_at, created_by, partner_id, origin_id,state, trigger, type, parameters, external_task, result_webhook, last_task_instance_id, require_noc_access, modified_by, modified_at, original_next_run_time, run_time, targets, target_type, schedule) FROM 'transformed_tasks.csv' WITH HEADER = TRUE AND MAXBATCHSIZE = %v;", c.Int(maxBatchSizeFlag)))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("ERROR: can not execute command copy from for tasks_by_runtime_mv table, err: %v", err)
	}

	cmd = exec.Command("cqlsh", "-e", fmt.Sprintf("use platform_tasking_db; COPY tasks_by_id_mv (run_time_unix, definition_id, iteration_ids, id, name, description,managed_endpoint_id, created_at, created_by, partner_id, origin_id,state, trigger, type, parameters, external_task, result_webhook, last_task_instance_id, require_noc_access, modified_by, modified_at, original_next_run_time, run_time, targets, target_type, schedule) FROM 'transformed_tasks.csv' WITH HEADER = TRUE AND MAXBATCHSIZE = %v;", c.Int(maxBatchSizeFlag)))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("ERROR: can not execute command copy from for tasks_by_id_mv table, err: %v", err)
	}

	return nil
}

func setUpLogging() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.Create(currentDir + "/" + time.Now().Truncate(time.Minute).Format(time.RFC3339) + "_copy_tasks.log")
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}
