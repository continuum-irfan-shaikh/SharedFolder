package main

import (
	"context"
	"log"
	"os"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"github.com/gocql/gocql"
	"github.com/urfave/cli"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"

	"io"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
)

func main() {
	app := cli.NewApp()
	app.Name = "Continuum Juno Automation TaskTargets converting"
	app.Usage = "This piece of software will recalculate tasks for each endpoint"
	app.Version = "0.0.1"

	config.Load()
	config.Config.CassandraTimeoutSec = 1500
	logger.Load(false)
	cassandra.Load()

	app.Action = Recalculate
	app.Run(os.Args)
}

func Recalculate(c *cli.Context) error {
	ctx := context.Background()
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	logFile, err := os.Create(currentDir + "/recalculation.log." + time.Now().String())
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	partners, err := models.TaskCounter.GetAllPartners(ctx)
	if err != nil {
		log.Printf("RecalculateAllCounters: GetAllPartners: err: %v", err)
		return err
	}
	log.Printf("RecalculateAllCounters: GetAllPartners: got partners [%v]", partners)

	for currentPartner := range partners {
		var (
			endpointsForUpdate   = make(map[gocql.UUID]bool)         // for current partner
			countersForResetting = map[gocql.UUID]models.TaskCount{} // this counters should be set to 0
		)

		taskCountsFromTasksTable, err := models.TaskPersistenceInstance.GetCountsByPartner(ctx, currentPartner)
		if err != nil {
			log.Printf("TaskPersistenceInstance.GetCountsByPartner: %v", err)
			continue
		}

		// gocql.UUID{} means that all counter for each ManagedEndpoint of current partner will be fetched form DB
		taskCountsFromTaskCounterTable, err := models.TaskCounter.GetCounters(ctx, currentPartner, gocql.UUID{})
		if err != nil {
			log.Printf("TaskCounter.GetCounters: %v", err)
			continue
		}

		// prepare set of TaskCounters by ManagedEndpointIDs
		for _, c := range taskCountsFromTaskCounterTable {
			endpointsForUpdate[c.ManagedEndpointID] = false // is not updated yet
			countersForResetting[c.ManagedEndpointID] = c
		}

		for _, counter := range taskCountsFromTasksTable {
			currentCount := counter.Count
			currentEndpoint := counter.ManagedEndpointID

			if countersForResetting[currentEndpoint].Count == currentCount {
				endpointsForUpdate[currentEndpoint] = true // marks it like updated 'cause they equal already
				log.Println("RecalculateAllCounters: counters are equal")
				continue
			}

			log.Printf("RecalculateAllCounters: processing partnerID [%s] and endpointID [%v]", currentPartner, currentEndpoint)
			err := updateCounters(ctx, currentPartner, currentEndpoint, currentCount)
			if err != nil {
				log.Printf("RecalculateAllCounters: updateCounters: %v", err)
				continue
			}
			endpointsForUpdate[currentEndpoint] = true
		}

		resetCounters := make([]models.TaskCount, 0)
		for endpoint, isUpdated := range endpointsForUpdate {
			// set counter for this particular ManagedEndpoint to 0
			if !isUpdated {
				log.Printf("Removing TaskCount for partnerID [%s] and endpointID [%s]\n", currentPartner, endpoint.String())

				resetCounters = append(resetCounters, models.TaskCount{
					ManagedEndpointID: endpoint,
					Count:             countersForResetting[endpoint].Count,
				})
			}
		}

		if len(resetCounters) > 0 {
			err = models.TaskCounter.DecreaseCounter(currentPartner, resetCounters, false)
			if err != nil {
				logger.Log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "got error while removing count: %v\n", err)
			}
		}
	}
	return nil
}

// updateCounters sets `correctCount` as actual one for current partner and endpoint
func updateCounters(ctx context.Context, partner string, endpoint gocql.UUID, correctCount int) (err error) {
	counts, err := models.TaskCounter.GetCounters(ctx, partner, endpoint)
	if err != nil {
		counts = []models.TaskCount{{
			ManagedEndpointID: endpoint,
			Count:             0,
		}}
	}

	err = models.TaskCounter.DecreaseCounter(partner, counts, false)
	if err != nil {
		return err
	}

	err = models.TaskCounter.IncreaseCounter(partner, []models.TaskCount{{ManagedEndpointID: endpoint, Count: correctCount}}, false)
	if err != nil {
		return err
	}

	return
}
