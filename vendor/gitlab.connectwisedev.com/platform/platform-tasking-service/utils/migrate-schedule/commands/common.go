package commands

import (
	"encoding/json"
	"fmt"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/utils/migrate-schedule/cassandra"
	"github.com/urfave/cli"
	"io/ioutil"
)

// Config keeps user data
type Config struct {
	CassandraURL        string `json:"CassandraURL"`
	CassandraKeyspace   string `json:"CassandraKeyspace"`
	CassandraTimeoutSec int    `json:"CassandraTimeoutSec"`
	CassandraConnNumber int    `json:"CassandraConnNumber"`
}

//UpdateTasksTable transforms tasks and task_instances tables
func UpdateTasksTable(c *cli.Context) error {
	var configPath = c.GlobalString("config")

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("Unable to read config file %v", configPath))
	}

	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(fmt.Sprintf("Unable to unmarshal config file: %v", err))
	}

	cassandra.Load(config.CassandraURL, config.CassandraKeyspace, config.CassandraTimeoutSec, config.CassandraConnNumber)
	if err = cassandra.UpdateTasksTable(); err != nil {
		panic(fmt.Sprintf("Error while updating tasks table: %v", err))
	}
	return nil
}
