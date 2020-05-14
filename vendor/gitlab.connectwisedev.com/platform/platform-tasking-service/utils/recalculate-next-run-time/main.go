package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"github.com/gocql/gocql"
	"github.com/urfave/cli"
)

const partner = "partner"
const id = "id"

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    id,
			Aliases: nil,
			Usage:   "taskID to reanimate",
		},
		&cli.StringFlag{
			Name:  partner,
			Usage: "partnerID to reanimate",
		},
	}
	app.Action = reanimateTask
	app.Run(os.Args)
}

func reanimateTask(c *cli.Context) error {
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

	partnerID := c.String(partner)
	taskID := c.String(id)

	if partnerID == "" || taskID == "" {
		log.Fatal(fmt.Errorf("partner and id cannot be empty"))
	}

	uuid, err := gocql.ParseUUID(taskID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("you've chosen %v partnerID and %v taskID", partnerID, uuid)

	ctx := context.TODO()
	tasks, err := models.TaskPersistenceInstance.GetByIDs(ctx, nil, partnerID, false, uuid)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %v tasks", len(tasks))

	now := time.Now().UTC().Truncate(time.Minute)
	nextTimeMap := make(map[time.Time]time.Time) // runTimeUTc - next run time
	for i, t := range tasks {
		if next, ok := nextTimeMap[t.RunTimeUTC]; ok {
			tasks[i].RunTimeUTC = next
			continue
		}

		nextRunTime, err := common.CalcFirstNextRunTime(now, t.Schedule)
		if err != nil {
			log.Fatal(err)
		}

		nextTimeMap[t.RunTimeUTC] = nextRunTime.UTC()
		tasks[i].RunTimeUTC = nextRunTime.UTC()
	}

	if err = models.TaskPersistenceInstance.InsertOrUpdate(ctx, tasks...); err != nil {
		log.Fatal(err)
	}
	log.Println("Finished")
	return nil
}

func setConfig() {
	config.Config.CassandraURL = "127.0.0.1:9042"
	config.Config.CassandraKeyspace = "platform_tasking_db"
	config.Config.CassandraConnNumber = 20
	config.Config.CassandraTimeoutSec = 15
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
