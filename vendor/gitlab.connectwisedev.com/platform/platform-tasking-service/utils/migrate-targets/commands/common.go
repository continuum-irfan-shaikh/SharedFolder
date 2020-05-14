package commands

import (
	"encoding/json"
	"fmt"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/utils/migrate-targets/cassandra"
	"github.com/urfave/cli"
	"io/ioutil"
)

//Config keeps user data
type Config struct {
	CassandraURL      string `json:"cassandraURL"`
	CassandraKeyspace string `json:"cassandraKeyspace"`
}

//TransformTasksTable transforms targets to task_targets  field.
func TransformTasksTable(c *cli.Context) error {
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

	cassandra.Load(config.CassandraURL, config.CassandraKeyspace)
	cassandra.TransformTasksTable()

	return nil
}
