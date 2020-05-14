package redis

import (
	"time"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/circuit"
	"github.com/go-redis/redis"
)

// Config : Redis client configuration
type Config struct {
	// host:port address.
	ServerAddress []string
		// ================= Cluster Specific Options =================

	// The maximum number of retries before giving up. Command is retried
	// on network errors and MOVED/ASK redirects.
	// Default is 8 retries.
	MaxRedirects int

	// Enables read-only commands on slave nodes.
	ReadOnly bool
	// Allows routing read-only commands to the closest master or slave node.
	// It automatically enables ReadOnly.
	RouteByLatency bool
	// Allows routing read-only commands to the random master or slave node.
	// It automatically enables ReadOnly.
	RouteRandomly bool

	// ================= Common Options ======================
	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string
	// Database to be selected after connecting to the server.
	DB int
	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge time.Duration
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout time.Duration
	// Frequency of idle checks made by idle connections reaper.
	// Default is 1 minute. -1 disables idle connections reaper,
	// but idle connections are still discarded by the client
	// if IdleTimeout is set.
	IdleCheckFrequency time.Duration
	//CircuitBreaker :All the Circuit breaker related configurations
	CircuitBreaker circuit.Config
	//CommandName :  Command name for redis circuit breaker
	CommandName string
}

//go:generate mockgen -destination=redismock/mocks.go -package=redismock gitlab.connectwisedev.com/platform/platform-common-lib/src/redis Client,Pipeliner

// Client : Redis client services
type Client interface {
	Init() error
	Close() error
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key ...string) error
	Incr(key string) (int64, error)
	Expire(key string, duration time.Duration) (bool, error)
	SetWithExpiry(key string, value interface{}, duration time.Duration) error
	Scan(cursor uint64, match string, count int64) (keys []string, outCursor uint64, err error)
	SubscribeChannel(pattern string) (<-chan *redis.Message, error)
	CreatePipeline() Pipeliner
	// ZAdd: Add member to a sorted set, or update its score if it already exists
	ZAdd(key string, members ...Z) (int64, error)
	// ZRange: Return a range of members in a sorted set, by index( Start: starting index, STOP : ending index)
	ZRange(key string, start, stop int64) ([]string, error)
	// ZRem:  Remove one or more members from a sorted set
	ZRem(key string, member interface{}) (int64, error)
	//Exists: check existance of key in set
	Exists(key string) (int64, error)
}


// Pipeliner : Redis Client's Pipeliner interface
type Pipeliner interface {
	PSet(key string, value interface{}) error
	PSetWithExpiry(key string, value interface{}, duration time.Duration) error
	PGet(key string) error
	Exec() ([]CmdOut, error)
	ClosePipeliner() error
}

// CmdOut : Output of Redis Command
type CmdOut struct {
	Name string
	Args []interface{}
	Err  error
}

// Pipe : Redis Pipe object
type pipe struct {
	pipeliner redis.Pipeliner
}
