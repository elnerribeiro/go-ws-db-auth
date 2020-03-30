package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/elnerribeiro/go-mustache-db"
	dbx "github.com/go-ozzo/ozzo-dbx"
	ozzolog "github.com/go-ozzo/ozzo-log"
)

//Logger is the logging system
var Logger *ozzolog.Logger

//DataBase is the database connections
var DataBase *dbx.DB

//Message returns a Json message
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

//Respond encodes the response
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

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
	Logger = ozzolog.NewLogger()
	Logger.Targets = []ozzolog.Target{ozzolog.NewConsoleTarget()}
	Logger.Open()
	db.InitDb(Logger)
	DataBase = db.Database
}
