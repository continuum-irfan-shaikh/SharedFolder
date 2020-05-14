package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
)

var (
	instance      *provider
	providerLock  sync.Mutex
	providerCache = make(map[string]DatabaseProvider)
)

type provider struct {
	driver     string
	datasource string
	db         *sql.DB
	config     Config
}

//GetDbProvider - Fetching and initializing Database Provider using db configurations.
//Once it is invoked it will return instance of provider struct to access db.
//  Returns - DatabaseProvider : instance of provider struct to access db.
//error : incase db gives error creating connection for given configs.
func GetDbProvider(config Config) (DatabaseProvider, error) {
	initializeCache(config)

	dialect, ok := getDialect(config.Driver)
	if !ok {
		return nil, fmt.Errorf("GetDbProvider: failed to get dialect instance for " + config.Driver)
	}
	dbConnInfo, err := dialect.GetConnectionString(config)
	if err != nil {
		return nil, fmt.Errorf("GetDbProvider: failed to get database connection config. " + err.Error())
	}

	providerLock.Lock()
	defer providerLock.Unlock()
	if providerInstance, ok := providerCache[dbConnInfo]; !ok {
		database, err := getConnection(config.Driver, dbConnInfo)
		if err != nil {
			return nil, err
		}

		instance = &provider{
			datasource: dbConnInfo,
			driver:     config.Driver,
			db:         database,
			config:     config,
		}

		providerCache[dbConnInfo] = instance
	} else {
		instance = providerInstance.(*provider)
	}

	return instance, nil
}

//convertSQLRowsToMap returns rows from db as []map[string]interface{}
func convertSQLRowsToMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	err := errors.New("ConvertSQLRowsToMap: Invalid sql rows")
	count := 0
	if rows == nil {
		return nil, err
	}
	rowsHolder := make([]map[string]interface{}, 0, 1)
	defer rows.Close() //nolint:errcheck
	for rows.Next() {
		count++
		cols, er := rows.Columns()
		if er != nil {
			return nil, er
		}
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		err = rows.Scan(columnPointers...)

		if err != nil {
			return nil, err
		}
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		rowsHolder = append(rowsHolder, m)
		err = nil
	}
	if count == 0 {
		return rowsHolder, nil
	}
	return rowsHolder, err

}

var getConnection = func(driver string, datasource string) (*sql.DB, error) {
	db, err := sql.Open(driver, datasource)
	if err == nil {
		err = db.Ping()
	}
	return db, err
}

//ExecWithPrepare is used to execute query that does not return data rows using prepared statement ex: INSERT, UPDATE or DELETE.
//  Ex: ExecWithPrepare(someQuery,val1,val2,val3)
//Returns - error: incase the database get error creating prepared statement or executing query.
func (c *provider) ExecWithPrepare(query string, value ...interface{}) error {
	stmt := getStatement(query)
	if stmt == nil {
		Logger().Info("", "Creating new prepared statement for exec query: "+query)
		st, err := c.db.Prepare(query)
		if err != nil {
			return err
		}
		addStatement(query, st)
		stmt = st
	}
	_, err := stmt.Exec(value...)
	return err
}

//SelectWithPrepare is used to execute select query using prepared statement.
//This will fetch the results on the basis of query provided and the values.
//  Ex: SelectWithPrepare(someQuery,val1,val2,val3)
//Returns - []map[string]interface{}: this will contain rows of data.
//error: incase the database gets error creating prepared statement or query execution.
//  Note: Incase query returns large result data, all the data rows will be returned at once.
func (c *provider) SelectWithPrepare(query string, value ...interface{}) ([]map[string]interface{}, error) {
	stmt := getStatement(query)
	if stmt == nil {
		Logger().Info("", "Creating new prepared statement for select query: "+query)
		st, err := c.db.Prepare(query)
		if err != nil {
			return nil, err
		}
		addStatement(query, st)
		stmt = st
	}
	rows, err := stmt.Query(value...)
	if err != nil {
		return nil, err
	}
	records, err := convertSQLRowsToMap(rows)
	if err != nil {
		return nil, err
	}
	return records, nil
}

//Exec is used to execute plaintext query that does not return data rows ex: INSERT, UPDATE or DELETE.
//  Ex: Exec(someQuery).
//Returns - error: incase database gets error executing the query.
func (c *provider) Exec(query string) error {
	_, err := c.db.Exec(query)
	return err
}

//Select is used to execute plaintext select query.
//  Ex: Select(someQuery)
//Returns - []map[string]interface{}: this contains rows of data.
//error: incase the database gets error executing the query
//  Note: Incase query returns large result data, all the data rows will be returned at once.
func (c *provider) Select(query string) ([]map[string]interface{}, error) {
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	records, err := convertSQLRowsToMap(rows)
	if err != nil {
		return nil, err
	}
	return records, nil
}

//CloseStatement is used to close prepared statement created for given query.
//Returns - error: if database get error closing prepared statement.
func (c *provider) CloseStatement(query string) error {
	st := getStatement(query)
	if st == nil {
		return nil
	}
	err := st.Close()
	if err != nil {
		return err
	}
	delete(query)
	return nil
}

//SelectAndProcess is used to execute plaintext select query.
//SelectAndProcess will process each row returned by query using a callback function.
//	Ex: SelectAndProcess(someQuery, callbackFunction)
func (c *provider) SelectAndProcess(query string, callback ProcessRow) {
	rows, err := c.db.Query(query)
	if err != nil {
		callback(Row{Error: err})
		return
	}

	processRows(rows, callback)
}

//SelectWithPrepareAndProcess is used to execute select query using prepared statement.
//SelectWithPrepareAndProcess will process each row returned by query using a callback function
//  Ex: SelectWithPrepareAndProcess(someQuery, callbackFunction, val1,val2...)
func (c *provider) SelectWithPrepareAndProcess(query string, callback ProcessRow, value ...interface{}) {
	stmt := getStatement(query)
	if stmt == nil {
		Logger().Info("", "Creating new prepared statement for select query: "+query)
		st, err := c.db.Prepare(query)
		if err != nil {
			callback(Row{Error: err})
			return
		}
		addStatement(query, st)
		stmt = st
	}
	rows, err := stmt.Query(value...)
	if err != nil {
		callback(Row{Error: err})
		return
	}

	processRows(rows, callback)
}

//processRows process row returned by query using callback function one row at time
func processRows(rows *sql.Rows, callback ProcessRow) {
	defer rows.Close() //nolint
	columnNames, err := rows.Columns()
	if err != nil {
		callback(Row{Error: fmt.Errorf("Failed to find column %+v", err)})
		return
	}

	colCount := len(columnNames)
	columnValues := make([]interface{}, colCount)
	columnPtrs := make([]interface{}, colCount)

	for rows.Next() {
		for i := 0; i < colCount; i++ {
			columnPtrs[i] = &columnValues[i]
		}
		err = rows.Scan(columnPtrs...)
		if err != nil {
			callback(Row{Error: err})
			continue
		}

		columns := make([]Column, colCount)
		for i := 0; i < colCount; i++ {
			columns[i] = Column{Name: columnNames[i], Value: columnValues[i]}
		}
		callback(Row{Columns: columns})
	}
}
