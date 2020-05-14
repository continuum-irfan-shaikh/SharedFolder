package db

import "testing"

func TestInitializeCache1(t *testing.T) {

	t.Run("Cache limit correct", func(t *testing.T) {
		expectedCacheLimit := 200
		data = nil
		cacheLimit = 0
		initializeCache(Config{CacheLimit: 200})

		if cacheLimit != expectedCacheLimit {
			t.Errorf("execpting cache limit to be %v but got %v", expectedCacheLimit, cacheLimit)
		}
	})
}
