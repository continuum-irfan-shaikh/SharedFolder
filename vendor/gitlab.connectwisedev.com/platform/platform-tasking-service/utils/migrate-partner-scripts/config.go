package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gocql/gocql"

	"log"
	"os"
)

var (
	// Config is a package variable, which is populated during init() execution and shared to whole application
	Config Configuration
	// ConfigFilePath defines a path to JSON-config file
	ConfigFilePath = "config.json"

	// MsSQLConn is a component to access MsSQL server
	MsSQLConn *sql.DB

	// CassandraSession is a component providing access to Cassandra storage
	CassandraSession *gocql.Session
)

// Configuration options
type Configuration struct {
	CassandraURL        string `json:"CassandraURL"                              default:"localhost:9042"`
	CassandraTimeoutSec int    `json:"CassandraTimeoutSec"                       default:"1"`
	CassandraKeyspace   string `json:"CassandraKeyspace"                         default:"platform_tasking_db"`
	CassandraConnNumber int    `json:"CassandraConnNumber"                       default:"20"`
	MsSQLDataBaseName   string `json:"MsSQLDataBaseName"`
	MsSQLUser           string `json:"MsSQLUser"`
	MsSQLPassword       string `json:"MsSQLPassword"`
	MsSQLAddress        string `json:"MsSQLAddress"`
	MsSQLPort           int    `json:"MsSQLPort"`
}

func LoadConfig() {
	var err error

	confLen := len(ConfigFilePath)
	if confLen != 0 {
		err = readConfigFromJSON(ConfigFilePath)
	}

	if err != nil {
		panic(`Configuration not found. Please specify configuration`)
	}
}

// nolint: gosec
func readConfigFromJSON(configFilePath string) error {
	log.Printf("Looking for JSON config file (%s)", configFilePath)

	cfgFile, err := os.Open(configFilePath)
	if err != nil {
		log.Printf("Reading configuration from JSON (%s) failed: %v\n", configFilePath, err)
		return err
	}
	defer func() {
		err := cfgFile.Close()
		if err != nil {
			log.Printf("Cannot close the configuration file [%s]: %v\n", cfgFile.Name(), err)
		}
	}()

	err = json.NewDecoder(cfgFile).Decode(&Config)
	if err != nil {
		log.Printf("Reading configuration from JSON (%s) failed: %s\n", configFilePath, err)
		return err
	}

	log.Printf("Configuration has been read from JSON (%s) successfully\n", configFilePath)
	return nil
}

func LoadMsSQLDBFromConfig() (err error) {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		Config.MsSQLAddress, Config.MsSQLUser, Config.MsSQLPassword, Config.MsSQLPort, Config.MsSQLDataBaseName)

	MsSQLConn, err = sql.Open("sqlserver", connString)
	if err != nil {
		return fmt.Errorf("error creating connection pool: %v", err.Error())
	}
	return
}

func LoadMsSQLDBFromString(connString string) (err error) {
	MsSQLConn, err = sql.Open("sqlserver", connString)
	if err != nil {
		return fmt.Errorf("error creating connection pool: %v", err.Error())
	}
	return
}

func LoadCassandra() (err error) {
	urls := strings.SplitN(Config.CassandraURL, ",", -1)
	client := gocql.NewCluster(urls...)
	client.ProtoVersion = 4
	client.Consistency = gocql.Quorum
	client.Keyspace = Config.CassandraKeyspace
	client.NumConns = Config.CassandraConnNumber
	client.Timeout = time.Duration(Config.CassandraTimeoutSec) * time.Second
	gocql.Logger = logger.CassandraLogger

	CassandraSession, err = client.CreateSession()
	if err != nil {
		return err
	}
	return
}

// QueryCassandra creates a session with cassandra and return a point on a query
func QueryCassandra(cql string, values ...interface{}) *gocql.Query {
	cassandraQuery := CassandraSession.Query(cql, values...)
	return cassandraQuery
}
