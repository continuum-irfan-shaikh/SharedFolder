package cassandra

import (
	"errors"
	"time"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/circuit"
)

const commandName = "Database-Command"

//DbConfig - configuration which is required to connect to Cassandra db.
type DbConfig struct {
	// Hosts - Addresses for the initial connections. It is recommended to use the value set in
	// the Cassandra config for broadcast_address or listen_address, an IP address not
	// a domain name. This is because events from Cassandra will use the configured IP
	// address, which is used to index connected hosts. If the domain name specified
	// resolves to more than 1 IP address then the driver may connect multiple times to
	// the same host, and will not mark the node being down or up from events.
	// This is a mandatory field
	Hosts []string

	// Keyspace - initial keyspace
	// This is a mandatory field
	Keyspace string

	// TimeoutMillisecond - connection timeout
	// (default is 1 second)
	TimeoutMillisecond time.Duration

	// NumConn - number of connections per host
	// default - 20
	NumConn int

	// CircuitBreaker - Configuration for the Circuit breaker
	// default - circuit.New()
	CircuitBreaker *circuit.Config

	// CommandName - Name for Database command
	// defaults to - Database-Command
	CommandName string

	// ValidErrors - List of error to participates in the Circuit state calculation
	// Default values are -
	ValidErrors []string
}

// NewConfig - returns a configration object having default values
func NewConfig() *DbConfig {
	return &DbConfig{
		NumConn:            20,
		TimeoutMillisecond: time.Second,
		CircuitBreaker: &circuit.Config{
			Enabled: false, TimeoutInSecond: 5, MaxConcurrentRequests: 2500,
			ErrorPercentThreshold: 25, RequestVolumeThreshold: 300, SleepWindowInSecond: 10,
		},
		CommandName: commandName,
		ValidErrors: []string{},
	}
}

func validate(conf *DbConfig) error {
	if conf == nil || len(conf.Hosts) == 0 || conf.Keyspace == "" {
		return errors.New(ErrDbHostsAndKeyspaceRequired)
	}

	if conf.NumConn == 0 {
		conf.NumConn = 20
	}

	if conf.TimeoutMillisecond == 0 {
		conf.TimeoutMillisecond = 1 * time.Second
	}

	// We are adding this additional check to avoid any failure due to
	// worng timeout configuration by microservice team
	if conf.TimeoutMillisecond < time.Millisecond {
		conf.TimeoutMillisecond = conf.TimeoutMillisecond * time.Millisecond
	}

	if conf.CircuitBreaker == nil {
		conf.CircuitBreaker = &circuit.Config{
			Enabled: false, TimeoutInSecond: 5, MaxConcurrentRequests: 2500,
			ErrorPercentThreshold: 25, RequestVolumeThreshold: 300, SleepWindowInSecond: 10,
		}
	}

	if conf.CommandName == "" {
		conf.CommandName = commandName
	}

	if len(conf.ValidErrors) == 0 {
		conf.ValidErrors = []string{}
	}
	return nil
}
