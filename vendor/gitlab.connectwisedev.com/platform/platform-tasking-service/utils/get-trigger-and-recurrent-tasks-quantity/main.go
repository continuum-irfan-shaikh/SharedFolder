package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/urfave/cli"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
)

const (
	size    = "size"
	host    = "host"
	env     = "env"
	logName = "_retrieving_tasks_data.log"

	multiTrigger = "trigger+trigger"
	comboTrigger = "recurrent+trigger"
)

var sitesService = map[string]string{
	"dt":   "https://rmmitswebapi.dtitsupport247.net/v1",
	"prod": "http://rmmitswebapi.itsupport247.local/v1",
	"qa":   "http://rmmitswebapi.qaitsupport247.local/v1",
}

func main() {
	app := cli.App{}
	app.Name = "Continuum Juno Automation Tasks DB migration tool"
	app.Usage = "This piece of software will retrieve information about tasks that uses any kind of combination (recurrent+trigger, trigger+trigger)"
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
			Value: &cli.StringSlice{"127.0.0.1:9042"},
		},
		&cli.StringFlag{
			Name:  env,
			Usage: "Environment for getting data",
			Value: "prod",
		},
	}

	setUpLogging()

	app.Action = getData

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("ERROR while running app: %v", err)
	}
	log.Println("INFO: app finished its work")

}

var getData cli.ActionFunc = func(c *cli.Context) error {
	start := time.Now()
	log.Println("Copying data process started at ", start.Truncate(time.Minute).Format(time.RFC3339))
	defer func() {
		log.Printf("Copying data process finished in %v min \n", time.Now().Sub(start).Minutes())
	}()

	sitesServiceURL, ok := sitesService[c.String(env)]
	if !ok {
		log.Printf("ERROR: there is no sitesServiceURL speicified for env %v", c.String(env))
		return fmt.Errorf("can't get sites service url for env %v", c.String(env))
	}
	log.Printf("INFO: url for getting tasks: %v", c.String(env))

	httpClient := &http.Client{
		Timeout: time.Duration(60) * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout:     time.Duration(60) * time.Second,
			MaxIdleConns:        2 * 100,
			MaxIdleConnsPerHost: 100,
			DisableKeepAlives:   false,
		},
	}

	response, err := httpClient.Get(sitesServiceURL + "/partner/0/partners")
	if err != nil {
		log.Printf("ERROR: error while getting list of partners %v", err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		msg := fmt.Sprintf("ERROR: response status %v during the attempt to get list of partners ", response.StatusCode)
		log.Printf(msg)
		return fmt.Errorf(msg)
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		msg := fmt.Sprintf("ERROR: can't read the body %v. err: %v ", response.Body, err)
		log.Printf(msg)
		return fmt.Errorf(msg)
	}

	partnerIDs := parseGetAllPartnersResponseBody(respBody)
	if len(partnerIDs) < 1 {
		log.Fatalf("There are no partners IDs")
	}

	log.Printf("INFO: partners count: %v", len(partnerIDs))

	cassandraHosts := c.StringSlice(host)
	session := createCassandraSession(cassandraHosts)
	defer session.Close()

	selectTasksStmt := "SELECT id, schedule FROM platform_tasking_db.tasks WHERE partner_id = ?"

	retryPolicy := gocql.ExponentialBackoffRetryPolicy{
		NumRetries: 5,
		Min:        30 * time.Second,
		Max:        5 * time.Minute,
	}

	out := setUpOutPut()
	defer func() {
		err := out.Close()
		if err != nil {
			log.Printf("ERROR: can't close file, reason: %v", err)
		}
	}()

	for _, partnerID := range partnerIDs {
		partnerTasks := make(map[gocql.UUID]struct{})
		tasksByCombTriggerTypeCount := make(map[string]int64)
		query := session.Query(selectTasksStmt, partnerID).RetryPolicy(&retryPolicy)

		iterator := query.Iter()

		var taskID gocql.UUID
		var scheduleString string
		for iterator.Scan(&taskID, &scheduleString) {
			if _, ok := partnerTasks[taskID]; ok {
				continue
			}

			schedule := apiModels.Schedule{}

			err := json.Unmarshal([]byte(scheduleString), &schedule)
			if err != nil {
				log.Println("ERROR: can't unmarshal schedule")
				continue
			}

			if schedule.Regularity == apiModels.Recurrent && len(schedule.TriggerFrames) > 0 {
				tasksByCombTriggerTypeCount[comboTrigger] += 1
				partnerTasks[taskID] = struct{}{}
				continue
			}

			if len(schedule.TriggerFrames) > 1 {
				tasksByCombTriggerTypeCount[multiTrigger] += 1
				partnerTasks[taskID] = struct{}{}
				continue
			}
		}

		if len(tasksByCombTriggerTypeCount) != 0 {
			combined := tasksByCombTriggerTypeCount[comboTrigger]
			multi := tasksByCombTriggerTypeCount[multiTrigger]

			result := fmt.Sprintf("Partner %s has %v specific tasks, recurrent+trigger = %v, trigger+trigger = %v\n", partnerID, combined+multi, combined, multi)
			_, err := out.WriteString(result)
			if err != nil {
				log.Printf("ERROR: can't wtite result to file, reason: %v", err)
			}
		}

		err = iterator.Close()
		if err != nil {
			log.Printf("ERROR: can't close iterator for query <%s>: %v", query.Statement(), err)
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

func setUpOutPut() *os.File {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(currentDir + "/specific_tasks_count.txt")
	if err != nil {
		log.Fatalf("ERROR: can't create a file, reason: %v", err)
	}
	return out
}

func parseGetAllPartnersResponseBody(body []byte) []string {
	comaSeparatedPartnerIDs := strings.Trim(string(body), "[]")
	if len(comaSeparatedPartnerIDs) > 0 {
		return strings.Split(comaSeparatedPartnerIDs, ",")
	}
	return []string{}
}
