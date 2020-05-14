package freecache

import (
	"log"
	"sync"

	"github.com/coocood/freecache"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
)

var (
	cacheInstance apiModels.Cacher
	mu            = &sync.Mutex{}
)

type cache struct {
	c *freecache.Cache
}

// Set writes the value by the key to freecache
// If the key is larger than 65535 or value is larger than 1/1024 of the cache size,
// the entry will not be written to the cache. expireSeconds <= 0 means no expire,
// but it can be evicted when cache is full
func (c cache) Set(key, value []byte, expireSeconds int) (err error) {
	return c.c.Set(key, value, expireSeconds)
}

// Get return the value from freecache by the key or not found error
func (c cache) Get(key []byte) (value []byte, err error) {
	return c.c.Get(key)
}


// New returns concrete instance for the apiModels.Cacher interface
// it's a blocking operation
func New(cacheSize int) apiModels.Cacher {
	mu.Lock()
	defer mu.Unlock()

	if cacheInstance == nil {
		log.Printf("Setting up cahce with size [%d] bytes", cacheSize)
		cacheInstance = cache{
			c: freecache.NewCache(cacheSize),
		}
	}

	return cacheInstance
}
