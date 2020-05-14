package main

import (
	"fmt"
	"time"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/db"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/db/mssql"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
)

const (
	query       = "Select * from Ticket_Tracing where TicketID = ?"
	insertQuery = "Insert into Ticket_Tracing values(?,?,?)"
	updateQuery = "Update Ticket_Tracing set OutTime = ? where TicketID = ?"
	deleteQuery = "Delete from Ticket_Tracing where TicketID = ?"

	selectQuery  = "Select * from Ticket_Tracing where TicketID = 1"
	insertQuery2 = "Insert into Ticket_Tracing values(10,GETDATE(),GETDATE())"
)

func main() {

	logger.Create(logger.Config{}) //no lint
	db.Logger = logger.Get

	//Get DbProvider Instance
	db, err := db.GetDbProvider(db.Config{DbName: "NOCBO",
		Server:     "10.2.27.41",
		Password:   "its",
		UserID:     "its",
		Driver:     mssql.Dialect,
		CacheLimit: 200})

	if err != nil {
		fmt.Println(err)
		return
	}

	//Select Query by Creating Prepared Statement
	rows, err := db.SelectWithPrepare(query, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rows)

	//Insert Query by Creating Prepared Statement
	err = db.ExecWithPrepare(insertQuery, 5, time.Now(), time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}

	//Update Query by Creating Prepared Statement
	err = db.ExecWithPrepare(updateQuery, time.Now(), 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Delete Query by Creating Prepared Statement
	err = db.ExecWithPrepare(deleteQuery, 5)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Close prepared statement
	err = db.CloseStatement(query)

	//Plain text select query
	rows, err = db.Select(selectQuery)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rows)

	//Plain text exec query
	err = db.Exec(insertQuery2)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Plaintext select query with callback function to read rows
	db.SelectAndProcess(selectQuery, processRowCallback)

	//Select query using prepared statement and callback function to read rows
	db.SelectWithPrepareAndProcess(query, processRowCallback, 1)

}

//Callback func to process table row
func processRowCallback(row db.Row) {
	fmt.Printf("Row: %v | Error: %v", row.Columns, row.Error)
}
