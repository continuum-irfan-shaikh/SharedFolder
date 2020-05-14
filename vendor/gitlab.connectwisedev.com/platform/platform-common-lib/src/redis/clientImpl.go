package redis

import (
	"fmt"
	"time"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/utils"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/circuit"
	"github.com/go-redis/redis"
)

const (
	// ErrInvalidConfigurationError : Error Invalid Configuration Error
	ErrInvalidConfigurationError = "ErrInvalidConfigurationError"
	//defaultCommandName : Default command name for redis circuit breaker
	defaultCommandName = "RedisCommand"
)

// ClientImpl : Redis client implementation
type clientImpl struct {
	config *Config
	client *redis.Client
}

// Z represents sorted set member. :Redis version
type Z struct {
	Score  float64
	Member interface{}
}

//GetService is a function to return service instance
func GetService(transactionID string, config *Config) Client {
	if transactionID == "" {
		transactionID = utils.GetTransactionID()
	}
	if config.CommandName == "" {
		config.CommandName = fmt.Sprintf("%s_%s", defaultCommandName, transactionID)
	}
	circuit.Register(transactionID, config.CommandName, &config.CircuitBreaker, nil)
	return &clientImpl{config: config}
}

func (c *clientImpl) Init() error {
	if c.client == nil {
		redisClient, err := c.genrateRedisClient()
		if err != nil {
			return err
		}
		c.client = redisClient
	}
	return nil
}

func (c *clientImpl) genrateRedisClient() (*redis.Client, error) {

	if c.config == nil {
		return nil, fmt.Errorf(ErrInvalidConfigurationError)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:               c.config.ServerAddress[0],
		Password:           c.config.Password,
		DB:                 c.config.DB,
		PoolSize:           c.config.PoolSize,
		MinIdleConns:       c.config.MinIdleConns,
		MaxConnAge:         c.config.MaxConnAge,
		PoolTimeout:        c.config.PoolTimeout,
		IdleTimeout:        c.config.IdleTimeout,
		IdleCheckFrequency: c.config.IdleCheckFrequency,
	})

	return redisClient, nil
}

func (c *clientImpl) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

func (c *clientImpl) Set(key string, value interface{}) error {
	return c.client.Set(key, value, -1).Err()
}

func (c *clientImpl) Get(key string) (interface{}, error) {
	return c.client.Get(key).Result()
}

func (c *clientImpl) Delete(key ...string) error {
	return c.client.Del(key...).Err()
}

func (c *clientImpl) Expire(key string, duration time.Duration) (bool, error) {
	return c.client.Expire(key, duration).Result()
}

func (c *clientImpl) Incr(key string) (int64, error) {
	return c.client.Incr(key).Result()
}

func (c *clientImpl) SetWithExpiry(key string, value interface{}, duration time.Duration) error {
	return c.client.Set(key, value, duration).Err()
}

func (c *clientImpl) Scan(cursor uint64, match string, count int64) (keys []string, outCursor uint64, err error) {
	return c.client.Scan(cursor, match, count).Result()
}

func (c *clientImpl) SubscribeChannel(pattern string) (<-chan *redis.Message, error) {
	pubSub := c.client.Subscribe(pattern)
	return pubSub.Channel(), nil
}

func (c *clientImpl) CreatePipeline() Pipeliner {
	return &pipe{
		pipeliner: c.client.Pipeline(),
	}
}

// ClosePipeliner : Close Pipeliner
func (p *pipe) ClosePipeliner() error {
	if p.pipeliner != nil {
		return p.pipeliner.Close()
	}
	return nil
}

// PSet : Pipeliner Set
func (p *pipe) PSet(key string, value interface{}) error {
	return p.pipeliner.Set(key, value, -1).Err()
}

// PSetWithExpiry : Pipeliner Set With Expiry
func (p *pipe) PSetWithExpiry(key string, value interface{}, duration time.Duration) error {
	return p.pipeliner.Set(key, value, duration).Err()
}

// PGet : Pipeliner Get
func (p *pipe) PGet(key string) error {
	return p.pipeliner.Get(key).Err()
}

// Exec : Pipeliner Exec
func (p *pipe) Exec() ([]CmdOut, error) {
	out, err := p.pipeliner.Exec()
	outArray := []CmdOut{}
	if len(out) > 0 {
		for _, outData := range out {
			outArray = append(outArray, CmdOut{
				Name: outData.Name(),
				Args: outData.Args(),
				Err:  outData.Err(),
			})
		}
	}
	return outArray, err
}

// ZAdd: Add member to a sorted set, or update its score if it already exists
func (c *clientImpl) ZAdd(key string, members ...Z) (int64, error) {
	var zAddResult int64
	err := circuit.Do(c.config.CommandName, c.config.CircuitBreaker.Enabled, func() error {
		var err error
		z := make([]redis.Z, len(members))
		for i := 0; i < len(members); i++ {
			z[i].Score = members[i].Score
			z[i].Member = members[i].Member
		}
		zAddResult, err = c.client.ZAdd(key, z...).Result()
		return err
	}, nil)
	return zAddResult, err
}

//ZRange:  Return a range of members in a sorted set, by index( Start: starting index, STOP : ending index)
func (c *clientImpl) ZRange(key string, start, stop int64) ([]string, error) {
	var zRangeResult []string
	err := circuit.Do(c.config.CommandName, c.config.CircuitBreaker.Enabled, func() error {
		var err error
		zRangeResult, err = c.client.ZRange(key, start, stop).Result()
		return err
	}, nil)
	return zRangeResult, err
}

//ZRem: Remove one or more members from a sorted set
func (c *clientImpl) ZRem(key string, member interface{}) (int64, error) {
	var zRemResult int64
	err := circuit.Do(c.config.CommandName, c.config.CircuitBreaker.Enabled, func() error {
		var err error
		zRemResult, err = c.client.ZRem(key, member).Result()
		return err
	}, nil)
	return zRemResult, err
}


//Exists: check existance of key in set
func (c *clientImpl) Exists(key string) (int64, error) {
	var existResult int64
	err := circuit.Do(c.config.CommandName, c.config.CircuitBreaker.Enabled, func() error {
		var err error
		existResult, err = c.client.Exists(key).Result()
		return err
	}, nil)
	return existResult, err
}
