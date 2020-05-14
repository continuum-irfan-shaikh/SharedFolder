package commands

import (
	"encoding/json"
	"fmt"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/utils/migrate-locations-and-targets/cassandra"
	"github.com/urfave/cli"
	"io/ioutil"
)

//Config keeps user data
type Config struct {
	CassandraURL        string `json:"cassandraURL"`
	CassandraKeyspace   string `json:"cassandraKeyspace"`
	CassandraTimeoutSec int    `json:"cassandraTimeoutSec"`
}

//TransformTaskingDBTables transforms tasks and task_instances tables
func TransformTaskingDBTables(c *cli.Context) error {
	var (
		configPath = c.GlobalString("config")
	)

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("Unable to read config file %v", configPath))
	}

	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(fmt.Sprintf("Unable to unmarshal config file: %v", err))
	}

	cassandra.Load(config.CassandraURL, config.CassandraKeyspace, config.CassandraTimeoutSec)
	cassandra.TransformTaskingDBTables()

	return nil
}
