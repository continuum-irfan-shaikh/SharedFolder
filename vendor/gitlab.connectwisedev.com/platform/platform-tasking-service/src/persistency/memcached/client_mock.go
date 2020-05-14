package memcached

import (
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// MockCacheConf is a function type for configuring custom mock of Cache interface
type MockCacheConf func(mc *MockCache) *MockCache

// MClientMock mock for memcached
type MClientMock struct {
	Cache map[string]*memcache.Item
}

// Get obtain memcahced item by key
func (mcm MClientMock) Get(key string) (*memcache.Item, error) {
	value, ok := mcm.Cache[key]
	if !ok {
		return nil, fmt.Errorf("Item not found by key %v", key)
	}
	if value.Expiration < int32(time.Now().Unix()) {
		return nil, fmt.Errorf("Item not found by key %v", key)
	}
	return value, nil
}

// Set memcached item to mock
func (mcm MClientMock) Set(item *memcache.Item) error {
	mcm.Cache[item.Key] = item
	return nil
}

// Delete obtain memcahced item by key
func (mcm MClientMock) Delete(key string) error {
	return nil
}
