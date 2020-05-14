package db

import (
	"database/sql"
	"sync"

	cache "github.com/patrickmn/go-cache"
)

var (
	data       *cache.Cache
	once       sync.Once
	cacheLimit = 100
)

const (
	defaultExpiration = -1
)

//initialiseCache create a new cache instance and set cache limit
func initializeCache(cfg Config) {
	once.Do(func() {
		data = cache.New(defaultExpiration, defaultExpiration)
		if cfg.CacheLimit != 0 {
			cacheLimit = cfg.CacheLimit
		}
	})
}

//getStatement is used to fetch prepared for given key from cache
var getStatement = func(key string) (stmt *sql.Stmt) {
	v, _ := data.Get(key)
	value, ok := v.(*sql.Stmt)

	if !ok {
		value = nil
	}

	return value
}

//addStatement is used to cache prepared statement
var addStatement = func(key string, stmt *sql.Stmt) {
	if data.ItemCount() == cacheLimit {
		Logger().Info("", "Flusing prepared statement cache")
		flush()
	}
	data.Set(key, stmt, defaultExpiration)
}

//Delete key from cache
var delete = func(key string) {
	data.Delete(key)
}

//Flush is used to clear the cache
func flush() {
	data.Flush()
}
