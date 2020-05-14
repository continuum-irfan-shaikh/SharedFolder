package db

//Row is struct to hold single table row
type Row struct {
	Columns []Column
	Error   error
}

//Column is a struct to hold all the column values of row
type Column struct {
	Name  string
	Value interface{}
}

//ProcessRow is a callback function used to process table row
type ProcessRow func(row Row)

//DatabaseProvider is interface that holds all the functions related to Db
type DatabaseProvider interface {
	//SelectWithPrepare is used to execute select query using prepared statement.
	//This will fetch the results on the basis of query provided and the values.
	//  Ex: SelectWithPrepare(someQuery,val1,val2,val3)
	//Returns - []map[string]interface{}: this will contain rows of data.
	//error: incase the database gets error creating prepared statement or query execution.
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
