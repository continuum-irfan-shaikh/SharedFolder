package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"time"

	"encoding/json"

	"io/ioutil"
	"net/http"
	"strconv"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/utils/monthly-report/models"
	"github.com/mohae/struct2csv"
	"github.com/urfave/cli"
)

const (
	csvFileName       = "tasking_monthly_report.csv"
	monthDuration     = -31
	selectPartnersIDs = "select distinct partner_id from tasks"
	querySelectData   = "select partner_id, managed_endpoint_id, name, schedule, created_by, created_at, modified_at, modified_by, state from tasks where created_at > ? and created_at < ?  allow filtering "
	startDateFlag     = "startDate"
	endDateFlag       = "endDate"
)

var httpClient *http.Client

func main() {
	app := cli.NewApp()
	app.Name = "Continuum Juno Automation Tasks tool for retrieving monthly report"
	app.Usage = "This piece of software will get tasks's information for 1 month from tasks table and asset ms to csv"
	app.Version = "0.0.1"
	setUpLogging()
	log.Println("INFO: Logging was set up successfully")

	config.Load()
	log.Println("INFO: config was loaded successfully")

	cassandra.Load()
	log.Println("INFO: Cassandra session was created successfully")

	httpClient = &http.Client{
		Timeout: time.Duration(config.Config.HTTPClientTimeoutSec) * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout:     time.Duration(config.Config.HTTPClientTimeoutSec) * time.Second,
			MaxIdleConns:        2 * config.Config.HTTPClientMaxIdleConnPerHost,
			MaxIdleConnsPerHost: config.Config.HTTPClientMaxIdleConnPerHost,
			DisableKeepAlives:   false,
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  startDateFlag,
			Usage: "Date to begin export from.",
		},
		cli.StringFlag{
			Name:  endDateFlag,
			Usage: "End date to finish exporting.",
		},
	}

	app.Action = getExportedData
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("ERROR: failed to run app, err: %v", err)
	}
	log.Println("Data was successfully exported")
}

func getExportedData(c *cli.Context) error {
	var (
		queries []models.Data
		err     error
		names   map[string]string
	)
	checkedID := map[string]bool{}
	groupedSites := make(map[string][]string)
	queries, err = getDataCassandra(c)
	if err != nil {
		log.Fatalf("can't get data from DB, err - %s", err.Error())
		return err
	}

	for i, val := range queries {
		if !checkedID[val.PartnerID] {
			sites, err := GetUserSites(val.PartnerID)
			if err != nil {
				log.Printf("error while getting user sites, err - %s", err.Error())
				continue
			}
			names, err = GetFriendlyNameByID(val.PartnerID)
			if err != nil {
				log.Printf("error while getting machines name. Err=%s", err.Error())
				continue
			}
			groupedSites[val.PartnerID] = sites
			checkedID[val.PartnerID] = true
		}
		queries[i].UserSites = groupedSites[val.PartnerID]
		queries[i].MachineName = names[queries[i].MEID]
	}
	WriteToCSV(queries)
	return nil
}

func getDataCassandra(c *cli.Context) (queries []models.Data, err error) {
	var (
		emptyTime        time.Time
		data             models.Data
		schedule         models.Schedule
		timeData         models.TimeData
		sites            []string
		scheduleString   string
		timeToSearchFrom time.Time
		timeToSearchTo   time.Time
	)

	timeToSearchFrom = time.Now().AddDate(0, 0, monthDuration)
	timeToSearchTo = time.Now()

	if c.String(startDateFlag) != "" {
		startString := c.String(startDateFlag)
		timeToSearchFrom, err = time.Parse(time.RFC3339Nano, startString)
		if err != nil {
			return queries, err
		}
	}

	if c.String(endDateFlag) != "" {
		endString := c.String(endDateFlag)
		timeToSearchTo, err = time.Parse(time.RFC3339Nano, endString)
		if err != nil {
			return queries, err
		}
	}

	groupedSites := make(map[string][]string)
	queries = make([]models.Data, 0)
	query := cassandra.Session.Query(querySelectData, timeToSearchFrom, timeToSearchTo).
		PageState(nil).
		PageSize(300)
	iter := query.Iter()
	for {
		for iter.Scan(
			&data.PartnerID,
			&data.MEID,
			&data.TaskName,
			&scheduleString,
			&data.CreatedBy,
			&timeData.CreatedAt,
			&timeData.ModifiedAt,
			&data.ModifiedBy,
			&data.State,
		) {
			err := json.Unmarshal([]byte(scheduleString), &schedule)
			if err != nil {
				log.Printf("wrong format of schedule string %s. Err=%s", scheduleString, err.Error())
				continue
			}
			gotData := models.Data{
				PartnerID:    data.PartnerID,
				MEID:         data.MEID,
				TaskName:     data.TaskName,
				CreatedAt:    timeData.CreatedAt.String(),
				CreatedBy:    data.CreatedBy,
				StartRunTime: schedule.StartRunTime.String(),
				EndRunTime:   schedule.EndRunTime.String(),
				ModifiedAt:   timeData.ModifiedAt.String(),
				ModifiedBy:   data.ModifiedBy,
				State:        models.GetTaskStateText(data.State),
			}
			if gotData.ModifiedAt == emptyTime.String() {
				gotData.ModifiedAt = "-"
			}
			if gotData.ModifiedBy == "" {
				gotData.ModifiedBy = "-"
			}
			queries = append(queries, gotData)
			groupedSites[data.PartnerID] = sites
			data = models.Data{}
		}
		if len(iter.PageState()) > 0 {
			iter = query.PageState(iter.PageState()).Iter()
		} else {
			break
		}
	}
	if err := iter.Close(); err != nil {
		return queries, err
	}
	return queries, nil
}

func setUpLogging() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	logFile, err := os.Create(currentDir + "/" + time.Now().Truncate(time.Minute).Format(time.RFC3339) + "_get_monthly_report.log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

//UserSitesFromAsset gets user sites from asset
func GetUserSites(partnerID string) ([]string, error) {
	var sitesSlice []string
	url := fmt.Sprintf("%s/partner/%s/sites", config.Config.SitesMsURL, partnerID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	sitesList := models.Sites{}
	if err := json.Unmarshal(body, &sitesList); err != nil {
		return nil, err
	}

	sitesSlice = make([]string, 0)
	for _, v := range sitesList.SiteList {
		sitesSlice = append(sitesSlice, strconv.FormatInt(v.ID, 10))
	}
	return sitesSlice, nil
}

//GetFriendlyNameByID gets user sites from asset by partnerID
func GetFriendlyNameByID(partnerID string) (groupedNames map[string]string, err error) {

	var asset []models.Asset
	groupedNames = make(map[string]string)
	url := fmt.Sprintf("http://internal-continuum-asset-service-elb-int-1972580147.ap-south-1.elb.amazonaws.com/asset/v1/partner/%s/endpoints", partnerID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return groupedNames, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return groupedNames, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return groupedNames, err
	}
	if err := json.Unmarshal(body, &asset); err != nil {
		return groupedNames, err
	}
	groupedNames = models.StructToMap(asset)
	return groupedNames, nil
}

func WriteToCSV(data []models.Data) {
	file, err := os.Create(csvFileName)
	if err != nil {
		log.Fatal("Cannot create file", err)
		return
	}
	defer file.Close()
	w := struct2csv.NewWriter(file)
	err = w.WriteStructs(data)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}
}
