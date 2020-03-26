package repositories

import (
	"fmt"

	db "github.com/elnerribeiro/go-mustache-db"
	dbx "github.com/go-ozzo/ozzo-dbx"
	ozzolog "github.com/go-ozzo/ozzo-log"
)

//Logger is the logging system
var Logger *ozzolog.Logger

//DataBase is the database connections
var DataBase *dbx.DB

//FinalizeLog closes the logging system
func FinalizeLog() {
	fmt.Println("Logger finalizer executed.")
	if Logger != nil {
		Logger.Close()
	}
}

//FinalizeDB closes DB connection
func FinalizeDB() {
	fmt.Println("DB finalizer executed.")
	if DataBase != nil {
		DataBase.Close()
	}
}

func init() {
	Logger := ozzolog.NewLogger()
	Logger.Targets = []ozzolog.Target{ozzolog.NewConsoleTarget()}
	Logger.Open()
	db.InitDb(Logger)
	DataBase = db.Database
}
