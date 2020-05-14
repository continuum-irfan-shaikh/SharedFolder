# Description
This is a common implementation for SQL data access layer.

### Third-Party Libraties
  - **Name** : Go Mssql Driver
  - **Link** : https://github.com/denisenkom/go-mssqldb
  - **License** : [BSD 3-Clause "New" or "Revised" License] (https://github.com/denisenkom/go-mssqldb/blob/master/LICENSE.txt)
  - **Description** : Golang Microsoft SQL Server library.

  - **Name** : Go Cache
  - **Link** : https://github.com/patrickmn/go-cache
  - **License** : [MIT License] (https://github.com/patrickmn/go-cache/blob/master/LICENSE)
  - **Description** : go-cache is an in-memory key:value store/cache.

  - **Name** : go-sqlmock
  - **Link** : https://github.com/DATA-DOG/go-sqlmock
  - **License** : [BSD 3-Clause] (https://github.com/DATA-DOG/go-sqlmock/blob/master/LICENSE)
  - **Description** : go-sqlmock is a mock library implementing sql/driver.

### Use 

**Import Statement**

```go
import (
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/db"
	//Import for loading mssql driver
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/db/mssql"
)
```

**Configuration**

```go
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
```

**Supported Drivers**
* mssql

**DatabaseProvider Instance**
```go
	//GetDbProvider - Fetching and initializing Database Provider using db configurations.
	//Once it is invoked it will return instance of provider struct to access db.
	//  Returns - DatabaseProvider : instance of provider struct to access db
	//error : incase db gives error creating connection for given configs
	db, err := db.GetDbProvider(db.Config{DbName: "NOCBO",
		Server:     "10.2.27.41",
		Password:   "its",
		UserID:     "its",
		Driver:     mssql.Dialect,
		CacheLimit: 200})
```
****

**Interface Functions**
```go
//DatabaseProvider is interface that holds all the functions related to Db
type DatabaseProvider interface {
	//SelectWithPrepare is used to execute select query using prepared statment.
	//This will fetch the results on the basis of query provided and the values.
	//  Ex: SelectWithPrepare(someQuery,val1,val2,val3)
	//Returns - []map[string]interface{}: this will contain rows of data.
	//error: incase the database gets error creating prepared statement or query.
	//  Note: Incase query returns large result data, all the data rows will be returned at once.
	SelectWithPrepare(query string, value ...interface{}) ([]map[string]interface{}, error)

	//ExecWithPrepare is used to execute query that does not return data rows using prepared statement ex: INSERT, UPDATE or DELETE.
	//  Ex: ExecWithPrepare(someQuery,val1,val2,val3)
	//Returns - error: incase the database get error creating prepared statement or executing query.
	ExecWithPrepare(query string, value ...interface{}) error

	//Select is used to execute plaintext select query.
	//  Ex: Select(someQuery)
	//Returns - []map[string]interface{}: this contains rows of data.
	//error: incase the database gets error executing the query
	//  Note: Incase query returns large result data, all the data rows will be returned at once.
	Select(query string) ([]map[string]interface{}, error)

	//Exec is used to execute plaintext query that does not return data rows ex: INSERT, UPDATE or DELETE.
	//  Ex: Exec(someQuery).
	//Returns - error: incase database gets error executing the query.
	Exec(query string) error

	//SelectAndProcess is used to execute plaintext select query.
	//SelectAndProcess will process each row returned by query using a callback function.
	//	Ex: SelectAndProcess(someQuery, callbackFunction)
	SelectAndProcess(query string, callback ProcessRow)

	//SelectWithPrepareAndProcess is used to execute select query using prepared statement.
	//SelectWithPrepareAndProcess will process each row returned by query using a callback function
	//  Ex: SelectWithPrepareAndProcess(someQuery, callbackFunction, val1,val2...)
	SelectWithPrepareAndProcess(query string, callback ProcessRow, value ...interface{})

	//CloseStatement is used to close prepared statement created for given query.
	//Returns - error: if database get error closing prepared statement.
	CloseStatement(query string) error
}

```

**Registering Dialect**
- To register a dialect set config.Driver as Dialect name

**Note** : 
- Cache has limit on caching number of items (config.CacheLimit or defaultCacheLimit = 100). On exceeding this limit all the cache data will be flushed.

**Example**
	[Please refer:](gitlab.connectwisedev.com/platform/platform-common-lib/src/db/example/example.go)