package memcached

//go:generate mockgen -destination=./memcachedPersistence_mock.go -package=memcached gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/memcached Cache

import (
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
)

// Cache is used to mock memcached in tests
type Cache interface {
	Get(key string) (item *memcache.Item, err error)
	Set(item *memcache.Item) error
	Delete(key string) error
}

// CacheImpl implementation of Memcached Client
type CacheImpl struct{}

var (
	// MemCacheInstance memcache instance implementation
	MemCacheInstance Cache = CacheImpl{}
	mClient          *memcache.Client
)

//Load creates a memcached Client
func Load() {
	if mClient == nil {
		urls := strings.SplitN(config.Config.Memcached.MemcachedURL, ",", -1)
		mClient = memcache.New(urls...)
		mClient.MaxIdleConns = config.Config.Memcached.MaxIdleConns
		mClient.Timeout = time.Duration(config.Config.Memcached.TimeoutSec) * time.Second
	}
}

//Get is used to get info from memcached db by the specified key
func (CacheImpl) Get(key string) (*memcache.Item, error) {
	item, err := mClient.Get(key)

	if err != nil {
		return nil, err
	}

	return item, nil
}

//Set is used to set info to memcached db with the specified key and default TTL from config file
func (CacheImpl) Set(item *memcache.Item) error {
	if err := mClient.Set(item); err != nil {
		return err
	}
	return nil
}

//Delete is used to set info to memcached db with the specified key and default TTL from config file
func (CacheImpl) Delete(key string) error {
	if err := mClient.Delete(key); err != nil {
		return err
	}
	return nil
}
