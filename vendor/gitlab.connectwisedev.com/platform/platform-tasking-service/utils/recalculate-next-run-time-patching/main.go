package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"github.com/gocql/gocql"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Action = reanimatePatchingTasks
	app.Run(os.Args)
}

func reanimatePatchingTasks(c *cli.Context) error {
	if err := setUp(); err != nil {
		log.Fatal(err)
	}

	data := getTasksData(tasksPerPartner)
	for taskID, partner := range data {
		if err := reanimateTask(partner, taskID); err != nil {
			log.Printf("Failed to update %v for partner %v", taskID, partner)
			log.Println("Sleeping for 10 seconds")
			time.Sleep(time.Second * 10)
			if err := reanimateTask(partner, taskID); err != nil {
				log.Printf("Failed to update %v for partner %v", taskID, partner)
			}
		}
	}
	return nil
}

type tasksData map[string]string // taskID - partnerID

func getTasksData(dataToParse string) tasksData {
	data := make(tasksData)
	temp := strings.ReplaceAll(dataToParse, "\n", " | ")
	strangeSlice := strings.Split(temp, " | ")
	for i := range strangeSlice {
		if i%2 != 0 {
			data[strangeSlice[i]] = strangeSlice[i-1]
		}
	}
	return data
}

func reanimateTask(partnerID, taskID string) error {
	if partnerID == "" || taskID == "" {
		log.Fatal(fmt.Errorf("partner and id cannot be empty"))
	}

	partnerID = strings.TrimSpace(partnerID)
	taskID = strings.TrimSpace(taskID)

	uuid, err := gocql.ParseUUID(taskID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("you've chosen %v partnerID and %v taskID", partnerID, uuid)

	ctx := context.TODO()
	tasks, err := models.TaskPersistenceInstance.GetByIDs(ctx, nil, partnerID, false, uuid)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Found %v tasks for %v", len(tasks), taskID)

	if len(tasks) == 0 {
		return nil
	}

	now := time.Now().UTC().Truncate(time.Minute)
	nextTimeMap := make(map[time.Time]time.Time) // runTimeUTc - next run time
	toUpdate := make([]models.Task, 0, len(tasks))
	mustBeUpdated := false

	wrongTasksCount := 0
	for _, t := range tasks {
		if !t.ExternalTask || t.Type != "patching" || t.RunTimeUTC.After(now) || t.State == statuses.TaskStateInactive {
			log.Printf("Task has wrong parameters %v", t)
			wrongTasksCount++
			continue
		}
		mustBeUpdated = true
		if next, ok := nextTimeMap[t.RunTimeUTC]; ok {
			t.RunTimeUTC = next
			toUpdate = append(toUpdate, t)
			continue
		}

		nextRunTime, err := common.CalcFirstNextRunTime(now, t.Schedule)
		if err != nil {
			log.Println(err)
		}

		nextTimeMap[t.RunTimeUTC] = nextRunTime.UTC()
		t.RunTimeUTC = nextRunTime.UTC()
		toUpdate = append(toUpdate, t)
	}

	if wrongTasksCount == len(tasks) {
		log.Printf("All tasks has wrong parameters, id %v", tasks[0].ID)
	}

	if mustBeUpdated {
		taskInstance := models.NewTaskInstance(toUpdate, true)

		for i := range toUpdate {
			toUpdate[i].LastTaskInstanceID = taskInstance.ID
			toUpdate[i].CreatedAt = now
		}

		if err = models.TaskInstancePersistenceInstance.Insert(ctx, taskInstance); err != nil {
			return err
		}
	}

	if err = models.TaskPersistenceInstance.InsertOrUpdate(ctx, toUpdate...); err != nil {
		log.Println(err)
		return err
	}

	log.Println("Finished")
	return nil
}

func setUp() error {
	setConfig()

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.Create(currentDir + "/recalculation.log." + time.Now().String())
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	return err
}

func setConfig() {
	config.Config.CassandraURL = "127.0.0.1:9042"
	config.Config.CassandraKeyspace = "platform_tasking_db"
	config.Config.CassandraConnNumber = 20
	config.Config.CassandraTimeoutSec = 300
	config.Config.CassandraBatchSize = 5
	config.Config.Log.FileName = "platform-tasking-service.log"
	config.Config.Log.CallDepth = 5
	config.Config.Log.ServiceName = "platform_tasking_service"

	config.Load()
	if err := logger.Load(config.Config.Log); err != nil {
		log.Println("LoadApplicationServices: error during loading logger: ", err)
	}
	cassandra.Load()
}
