package db

import "gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"

// Logger : Logger instance used for logging
// Defaults to Discard
var Logger = logger.DiscardLogger

//Config is struct to define db configurations
type Config struct {
	//DbName - Db to be selected after connecting to server
	//Required
	DbName string

	//Server - ip address of db host server
	//Required
	Server string

	//UserID - UserId for db server
	//Required
	UserID string

	//Password - Password for db server
	//Required
	Password string

	//Driver - Name of db driver
	//Required
	Driver string

	//Map to hold additional db config
	AdditionalConfig map[string]string

	//CacheLimit - CacheLimit sets limit on number of prepared statements to be cached
	//Default CacheLimit: 100
	CacheLimit int
}
