package freecache

import (
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"github.com/coocood/freecache"
	"log"
	"sync"
)

var (
	cacheInstance persistency.Cache
	mu            = &sync.Mutex{}
)

type cache struct {
	c *freecache.Cache
}

// Set is set, life is life
func (c cache) Set(key, value []byte, expireSeconds int) (err error) {
	return c.c.Set(key, value, expireSeconds)
}

// Get is set, life is life
func (c cache) Get(key []byte) (value []byte, err error) {
	return c.c.Get(key)
}

// New returns concrete instance for the persistency.Cache interface
// it's a blocking operation
func New() persistency.Cache {
	mu.Lock()
	defer mu.Unlock()

	if cacheInstance == nil {
		log.Printf("Setting up cahce with size [%d] bytes", config.Config.InMemoryCacheSize)
		cacheInstance = cache{
			c: freecache.NewCache(config.Config.InMemoryCacheSize),
		}
	}

	return cacheInstance
}
