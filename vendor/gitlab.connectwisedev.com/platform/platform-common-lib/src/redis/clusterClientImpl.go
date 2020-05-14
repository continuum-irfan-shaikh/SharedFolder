package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/circuit"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/utils"
)

// ClientImpl : Redis client implementation
type clusterClientImpl struct {
	config        *Config
	clusterClient *redis.ClusterClient
}

//GetClusterClientService is a function to return service instance
func GetClusterClientService(transactionID string, config *Config) Client {
	if transactionID == "" {
		transactionID = utils.GetTransactionID()
	}
	if config.CommandName == "" {
		config.CommandName = fmt.Sprintf("%s_%s", defaultCommandName, transactionID)
	}
	circuit.Register(transactionID, config.CommandName, &config.CircuitBreaker, nil)
	return &clusterClientImpl{config: config}
}

func (c *clusterClientImpl) Init() error {
	if c.clusterClient == nil {
		redisClient, err := c.genrateRedisClusterClient()
		if err != nil {
			return err
		}
		c.clusterClient = redisClient
	}
	return nil
}
func (c *clusterClientImpl) genrateRedisClusterClient() (*redis.ClusterClient, error) {

	if c.config == nil {
		return nil, fmt.Errorf(ErrInvalidConfigurationError)
	}

	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:              c.config.ServerAddress,
		Password:           c.config.Password,
		MaxRedirects:       c.config.MaxRedirects,
		ReadOnly:           c.config.ReadOnly,
		RouteByLatency:     c.config.RouteByLatency,
		RouteRandomly:      c.config.RouteRandomly,
		PoolSize:           c.config.PoolSize,
		MinIdleConns:       c.config.MinIdleConns,
		MaxConnAge:         c.config.MaxConnAge,
		PoolTimeout:        c.config.PoolTimeout,
		IdleTimeout:        c.config.IdleTimeout,
		IdleCheckFrequency: c.config.IdleCheckFrequency,
	})

	return clusterClient, nil
}

func (c *clusterClientImpl) Close() error {
	if c.clusterClient != nil {
		return c.clusterClient.Close()
	}
	return nil
}

// ZAdd: Add member to a sorted set, or update its score if it already exists
func (c *clusterClientImpl) ZAdd(key string, members ...Z) (int64, error) {
	var zAddResult int64
	err := circuit.Do(c.config.CommandName, c.config.CircuitBreaker.Enabled, func() error {
		var err error
		z := make([]redis.Z, len(members))
		for i := 0; i < len(members); i++ {
			z[i].Score = members[i].Score
			z[i].Member = members[i].Member
		}
		zAddResult, err = c.clusterClient.ZAdd(key, z...).Result()
		return err
	}, nil)
	return zAddResult, err
}

//ZRange:  Return a range of members in a sorted set, by index( Start: starting index, STOP : ending index)
func (c *clusterClientImpl) ZRange(key string, start, stop int64) ([]string, error) {
	var zRangeResult []string
	err := circuit.Do(c.config.CommandName, c.config.CircuitBreaker.Enabled, func() error {
		var err error
		zRangeResult, err = c.clusterClient.ZRange(key, start, stop).Result()
		return err
	}, nil)
	return zRangeResult, err
}

//ZRem: Remove one or more members from a sorted set
func (c *clusterClientImpl) ZRem(key string, member interface{}) (int64, error) {
	var zRemResult int64
	err := circuit.Do(c.config.CommandName, c.config.CircuitBreaker.Enabled, func() error {
		var err error
		zRemResult, err = c.clusterClient.ZRem(key, member).Result()
		return err
	}, nil)
	return zRemResult, err
}

//Exists: check existance of key in set
func (c *clusterClientImpl) Exists(key string) (int64, error) {
	var existResult int64
	err := circuit.Do(c.config.CommandName, c.config.CircuitBreaker.Enabled, func() error {
		var err error
		existResult, err = c.clusterClient.Exists(key).Result()
		return err
	}, nil)
	return existResult, err
}
func (c *clusterClientImpl) Set(key string, value interface{}) error {
	return c.clusterClient.Set(key, value, -1).Err()
}

func (c *clusterClientImpl) Get(key string) (interface{}, error) {
	return c.clusterClient.Get(key).Result()
}

func (c *clusterClientImpl) Delete(key ...string) error {
	return c.clusterClient.Del(key...).Err()
}

func (c *clusterClientImpl) Expire(key string, duration time.Duration) (bool, error) {
	return c.clusterClient.Expire(key, duration).Result()
}

func (c *clusterClientImpl) Incr(key string) (int64, error) {
	return c.clusterClient.Incr(key).Result()
}

func (c *clusterClientImpl) SetWithExpiry(key string, value interface{}, duration time.Duration) error {
	return c.clusterClient.Set(key, value, duration).Err()
}

func (c *clusterClientImpl) Scan(cursor uint64, match string, count int64) (keys []string, outCursor uint64, err error) {
	return c.clusterClient.Scan(cursor, match, count).Result()
}

func (c *clusterClientImpl) SubscribeChannel(pattern string) (<-chan *redis.Message, error) {
	pubSub := c.clusterClient.Subscribe(pattern)
	return pubSub.Channel(), nil
}

func (c *clusterClientImpl) CreatePipeline() Pipeliner {
	return &pipe{
		pipeliner: c.clusterClient.Pipeline(),
	}
}
