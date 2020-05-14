package mssql

import (
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/db"

	_ "github.com/denisenkom/go-mssqldb" //To load mssql driver
)

const (

	//Dialect is a database name used for registration
	Dialect = "mssql"
)

func init() {
	db.RegisterDialect(Dialect, mssql{})
}

type mssql struct {
}

func (m mssql) GetConnectionString(config db.Config) (string, error) {
	if config.DbName == "" || config.Password == "" || config.Server == "" || config.UserID == "" {
		return "", fmt.Errorf("getDbConnInfo: One or more required db configuration  missing")
	}
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", config.Server, config.UserID, config.Password, config.DbName)
	return connString, nil
}
