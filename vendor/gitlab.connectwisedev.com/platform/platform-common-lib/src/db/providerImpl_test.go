package db

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type mockStruct struct {
}

func (s mockStruct) GetConnectionString(config Config) (string, error) {
	if config.DbName == "" || config.Password == "" || config.Server == "" || config.UserID == "" {
		return "", fmt.Errorf("getDbConnInfo: One or more required db configuration  missing")
	}
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", config.Server, config.UserID, config.Password, config.DbName)
	return connString, nil
}

var isError bool

//Callback func to read rows
func process(row Row) {
	if row.Error != nil {
		isError = true
	}
}

func TestGetDbProvider(t *testing.T) {
	dialectsMap["mockMssql"] = mockStruct{}

	t.Run("Error failed to get dialect", func(t *testing.T) {
		_, err := GetDbProvider(Config{DbName: "NOCBO",
			Server:     "10.2.27.41",
			Password:   "its",
			UserID:     "its",
			CacheLimit: 200})

		if err == nil {
			t.Errorf("Expecting error but found nil")
		}

	})

	t.Run("Error to get connection Config", func(t *testing.T) {
		_, err := GetDbProvider(Config{
			Driver: "mockMssql"})
		if err == nil {
			t.Errorf("Expecting error but found nil")
		}
	})

	t.Run("Error to get connection Config", func(t *testing.T) {
		old := getConnection
		getConnection = func(driver string, datasource string) (*sql.DB, error) {
			return nil, errors.New("Error getting db connection")
		}
		defer func() {
			getConnection = old
		}()
		_, err := GetDbProvider(Config{DbName: "NOCBO",
			Server:     "10.2.27.41",
			Password:   "its",
			UserID:     "its",
			Driver:     "mockMssql",
			CacheLimit: 200})
		if err == nil {
			t.Errorf("Expecting error but found nil")
		}
	})

	t.Run("Success getting db provider", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		expectedProvider := &provider{
			driver:     "mockMssql",
			datasource: "server=10.2.27.41;user id=its;password=its;database=NOCBO",
			db:         db,
			config: Config{DbName: "NOCBO",
				Server:     "10.2.27.41",
				Password:   "its",
				UserID:     "its",
				Driver:     "mockMssql",
				CacheLimit: 200},
		}

		old := getConnection
		getConnection = func(driver string, datasource string) (*sql.DB, error) {

			return db, nil
		}
		defer func() {
			getConnection = old
		}()
		gotPrivder, er := GetDbProvider(Config{DbName: "NOCBO",
			Server:     "10.2.27.41",
			Password:   "its",
			UserID:     "its",
			Driver:     "mockMssql",
			CacheLimit: 200})
		if er != nil {
			t.Errorf("Expecting no error but found err %v", err)
		}

		if !reflect.DeepEqual(expectedProvider, gotPrivder) {
			t.Errorf("Expected provider %v but got provider %v", expectedProvider, gotPrivder)
		}
	})

}

func TestExec(t *testing.T) {
	t.Run("Error running exec query", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectExec("Insert into Ticket_Tracing values(10,GETDATE(),GETDATE())").WillReturnError(fmt.Errorf("some error"))

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}

		err = p.Exec("Insert into Ticket_Tracing values(10,GETDATE(),GETDATE())")
		if err == nil {
			t.Errorf("Expecting error but found nil")
		}
	})

	t.Run("success exec query", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectExec("Update Ticket set ID = 1").WillReturnResult(sqlmock.NewResult(1, 1))

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}

		err = p.Exec("Update Ticket set ID = 1")
		if err != nil {
			t.Errorf("Expecting nil but found err %v", err)
		}
	})
}

func TestSelect(t *testing.T) {
	t.Run("select error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectQuery("Select * from Ticket_Tracing where TicketID = 1").WillReturnError(fmt.Errorf("some error"))

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		_, err = p.Select("Select * from Ticket_Tracing where TicketID = 1")
		if err == nil {
			t.Errorf("Expecting err but found nil")
		}
	})

	t.Run("select success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		rows := sqlmock.NewRows([]string{"Id", "InTime", "OutTime"}).AddRow(1, "10 Aug", "11 Aug").AddRow(2, "12 Aug", "13 Aug")

		mock.ExpectQuery("Select Id, InTime, OutTime from Tracing").WillReturnRows(rows)

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		_, err = p.Select("Select Id, InTime, OutTime from Tracing")
		if err != nil {
			t.Errorf("Expecting nil but found err %v", err)
		}

	})
}

func TestConvertSqlRowstoMap(t *testing.T) {
	t.Run("Success", func(t *testing.T) {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"Id", "InTime", "OutTime"}).AddRow(1, "10 Aug", "11 Aug").AddRow(2, "12 Aug", "13 Aug")

		mock.ExpectQuery("Select * from Tracing").WillReturnRows(rows)
		row, _ := db.Query("Select *s from Tracing")
		_, err = convertSQLRowsToMap(row)
		if err == nil {
			t.Errorf("expected error but found nil")
		}
	})

}

func TestCloseStatement(t *testing.T) {
	t.Run("success closing statement", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		old1 := delete
		getStatement = func(key string) (stmt *sql.Stmt) {
			st, _ := db.Prepare("Insert into Ticket values(?)")
			return st
		}
		delete = func(query string) {
			return
		}
		defer func() {
			getStatement = old
			delete = old1
		}()

		mock.ExpectPrepare("Insert into").WillBeClosed()

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}

		err = p.CloseStatement("Insert into Ticket values(?)")
		if err != nil {
			t.Errorf("Expecting nil but found err %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

}

func TestSelectWithPrepare(t *testing.T) {

	t.Run("error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		old1 := addStatement
		getStatement = func(key string) (stmt *sql.Stmt) {
			return nil
		}
		addStatement = func(key string, stmt *sql.Stmt) {
			return
		}
		defer func() {
			getStatement = old
			addStatement = old1
		}()

		mock.ExpectPrepare("Select from Ticket where ID = ?").ExpectQuery().WithArgs("2").WillReturnError(errors.New("Some Errors"))

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		_, err = p.SelectWithPrepare("Select from Ticket where ID = ?", "2")
		if err == nil {
			t.Errorf("Expecting error but found nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		old1 := addStatement
		getStatement = func(key string) (stmt *sql.Stmt) {
			return nil
		}
		addStatement = func(key string, stmt *sql.Stmt) {
			return
		}
		defer func() {
			getStatement = old
			addStatement = old1
		}()
		rows := sqlmock.NewRows([]string{"Id", "InTime", "OutTime"}).AddRow(1, "10 Aug", "11 Aug").AddRow(2, "12 Aug", "13 Aug")

		mock.ExpectPrepare("Select from Ticket where ID = ?").ExpectQuery().WithArgs("2").WillReturnRows(rows)

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		_, err = p.SelectWithPrepare("Select from Ticket where ID = ?", "2")
		if err != nil {
			t.Errorf("Expecting nil but found err %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Error cached prepared statement", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		old1 := addStatement
		getStatement = func(key string) (stmt *sql.Stmt) {
			stmt, _ = db.Prepare("Select from Ticket where ID = ?")
			return
		}
		addStatement = func(key string, stmt *sql.Stmt) {
			return
		}
		defer func() {
			getStatement = old
			addStatement = old1
		}()

		mock.ExpectPrepare("Select from Ticket where ID = ?").ExpectQuery().WithArgs("2").WillReturnError(errors.New("Some Errors"))

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		_, err = p.SelectWithPrepare("Select from Ticket where ID = ?", "2")
		if err == nil {
			t.Errorf("Expecting error but found nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Success cached prepared statement", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		old1 := addStatement
		getStatement = func(key string) (stmt *sql.Stmt) {
			stmt, _ = db.Prepare("Select from Ticket where ID = ?")
			return
		}
		addStatement = func(key string, stmt *sql.Stmt) {
			return
		}
		defer func() {
			getStatement = old
			addStatement = old1
		}()
		rows := sqlmock.NewRows([]string{"Id", "InTime", "OutTime"}).AddRow(1, "10 Aug", "11 Aug").AddRow(2, "12 Aug", "13 Aug")

		mock.ExpectPrepare("Select from Ticket where ID = ?").ExpectQuery().WithArgs("2").WillReturnRows(rows)

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		_, err = p.SelectWithPrepare("Select from Ticket where ID = ?", "2")
		if err != nil {
			t.Errorf("Expecting nil but found err %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("error creating prepared statement", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		getStatement = func(key string) (stmt *sql.Stmt) {
			return nil
		}

		defer func() {
			getStatement = old

		}()

		mock.ExpectPrepare("Select from Ticket where ID = ?").WillReturnError(errors.New("Some Error"))
		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		_, err = p.SelectWithPrepare("Select from Ticket where ID = ?", "2")
		if err == nil {
			t.Errorf("Expecting error but found nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestExecWithPrepare(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		old1 := addStatement
		getStatement = func(key string) (stmt *sql.Stmt) {
			return nil
		}

		addStatement = func(key string, stmt *sql.Stmt) {
			return
		}

		defer func() {
			getStatement = old
			addStatement = old1
		}()

		mock.ExpectPrepare("INSERT INTO Ticket VALUES(?)").ExpectExec().WithArgs("2").WillReturnResult(sqlmock.NewResult(1, 1))

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		if err = p.ExecWithPrepare("INSERT INTO Ticket VALUES(?)", "2"); err != nil {
			t.Errorf("Expecting nil but found err %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		old1 := addStatement
		getStatement = func(key string) (stmt *sql.Stmt) {
			return nil
		}

		addStatement = func(key string, stmt *sql.Stmt) {
			return
		}

		defer func() {
			getStatement = old
			addStatement = old1
		}()

		mock.ExpectPrepare("INSERT INTO Ticket VALUES(?)").ExpectExec().WithArgs("2").WillReturnError(errors.New("Some Error"))

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		if err = p.ExecWithPrepare("INSERT INTO Ticket VALUES(?)", "2"); err == nil {
			t.Errorf("Expecting error but found nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Error cached statement", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		getStatement = func(key string) (stmt *sql.Stmt) {
			stmt, _ = db.Prepare("INSERT INTO Ticket VALUES(?)")
			return
		}

		defer func() {
			getStatement = old
		}()

		mock.ExpectPrepare("INSERT INTO Ticket VALUES(?)").ExpectExec().WithArgs("2").WillReturnError(errors.New("Some Error"))

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		if err = p.ExecWithPrepare("INSERT INTO Ticket VALUES(?)", "2"); err == nil {
			t.Errorf("Expecting error but found nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Success cached statement", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		getStatement = func(key string) (stmt *sql.Stmt) {
			stmt, _ = db.Prepare("INSERT INTO Ticket VALUES(?)")
			return
		}

		defer func() {
			getStatement = old
		}()

		mock.ExpectPrepare("INSERT INTO Ticket VALUES(?)").ExpectExec().WithArgs("2").WillReturnResult(sqlmock.NewResult(1, 1))

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		if err = p.ExecWithPrepare("INSERT INTO Ticket VALUES(?)", "2"); err != nil {
			t.Errorf("Expecting nil but found err %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Error creating prepared statement", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		getStatement = func(key string) (stmt *sql.Stmt) {
			return nil
		}

		defer func() {
			getStatement = old
		}()

		mock.ExpectPrepare("INSERT INTO Ticket VALUES(?)").WillReturnError(errors.New("Some Error"))

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		if err = p.ExecWithPrepare("INSERT INTO Ticket VALUES(?)", "2"); err == nil {
			t.Errorf("Expecting error but found nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

	})

}

func TestSelectAndProcess(t *testing.T) {

	t.Run("select error", func(t *testing.T) {
		isError = false

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectQuery("Select * from Ticket_Tracing where TicketID = 1").WillReturnError(fmt.Errorf("some error"))

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		p.SelectAndProcess("Select * from Ticket_Tracing where TicketID = 1", process)
		if !isError {
			t.Errorf("Expecting err but found nil")
		}
	})

	t.Run("select success", func(t *testing.T) {
		isError = false
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		rows := sqlmock.NewRows([]string{"Id", "InTime", "OutTime"}).AddRow(1, "10 Aug", "11 Aug").AddRow(2, "12 Aug", "13 Aug")

		mock.ExpectQuery("Select Id, InTime, OutTime from Tracing").WillReturnRows(rows)

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}

		p.SelectAndProcess("Select Id, InTime, OutTime from Tracing", process)
		if isError {
			t.Errorf("Expecting nil but found err")
		}
	})
}

func TestSelectWithPrepareAndProcess(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		isError = false
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		old1 := addStatement
		getStatement = func(key string) (stmt *sql.Stmt) {
			return nil
		}
		addStatement = func(key string, stmt *sql.Stmt) {
			return
		}
		defer func() {
			getStatement = old
			addStatement = old1
		}()
		rows := sqlmock.NewRows([]string{"Id", "InTime", "OutTime"}).AddRow(1, "10 Aug", "11 Aug").AddRow(2, "12 Aug", "13 Aug")

		mock.ExpectPrepare("Select from Ticket where ID = ?").ExpectQuery().WithArgs("2").WillReturnRows(rows)

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		p.SelectWithPrepareAndProcess("Select from Ticket where ID = ?", process, "2")
		if isError {
			t.Errorf("Expecting nil but found err")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Error cached prepared statement", func(t *testing.T) {
		isError = true
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		old1 := addStatement
		getStatement = func(key string) (stmt *sql.Stmt) {
			stmt, _ = db.Prepare("Select from Ticket where ID = ?")
			return
		}
		addStatement = func(key string, stmt *sql.Stmt) {
			return
		}
		defer func() {
			getStatement = old
			addStatement = old1
		}()

		mock.ExpectPrepare("Select from Ticket where ID = ?").ExpectQuery().WithArgs("2").WillReturnError(errors.New("Some Errors"))

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		p.SelectWithPrepareAndProcess("Select from Ticket where ID = ?", process, "2")
		if !isError {
			t.Errorf("Expecting error but found nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Success cached prepared statement", func(t *testing.T) {
		isError = false
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		old1 := addStatement
		getStatement = func(key string) (stmt *sql.Stmt) {
			stmt, _ = db.Prepare("Select from Ticket where ID = ?")
			return
		}
		addStatement = func(key string, stmt *sql.Stmt) {
			return
		}
		defer func() {
			getStatement = old
			addStatement = old1
		}()
		rows := sqlmock.NewRows([]string{"Id", "InTime", "OutTime"}).AddRow(1, "10 Aug", "11 Aug").AddRow(2, "12 Aug", "13 Aug")

		mock.ExpectPrepare("Select from Ticket where ID = ?").ExpectQuery().WithArgs("2").WillReturnRows(rows)

		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		p.SelectWithPrepareAndProcess("Select from Ticket where ID = ?", process, "2")
		if isError {
			t.Errorf("Expecting nil but found err ")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("error creating prepared statement", func(t *testing.T) {
		isError = false
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		old := getStatement
		getStatement = func(key string) (stmt *sql.Stmt) {
			return nil
		}

		defer func() {
			getStatement = old

		}()

		mock.ExpectPrepare("Select from Ticket where ID = ?").WillReturnError(errors.New("Some Error"))
		p := provider{
			driver:     "mssql",
			datasource: "connection string",
			db:         db,
			config:     Config{},
		}
		p.SelectWithPrepareAndProcess("Select from Ticket where ID = ?", process, "2")
		if !isError {
			t.Errorf("Expecting error but found nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
}
