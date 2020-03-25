package main

import (
	"fmt"
	"net/http"

	"github.com/elnerribeiro/go-ws-db-auth/app"
	"github.com/elnerribeiro/go-ws-db-auth/controllers"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api/users", controllers.ListUsers).Methods("POST")
	router.HandleFunc("/api/user/{id:[0-9]+}", controllers.GetUserByID).Methods("GET")
	router.HandleFunc("/api/user", controllers.Upsert).Methods("PUT")
	router.HandleFunc("/api/user/{id:[0-9]+}", controllers.Delete).Methods("DELETE")
	router.HandleFunc("/api/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/validate", controllers.Validate).Methods("GET")
	router.HandleFunc("/api/insert/{id:[0-9]+}", controllers.ListInsert).Methods("GET")
	router.HandleFunc("/api/insert/sync/{qty:[0-9]+}", controllers.InsertSync).Methods("PUT")
	router.HandleFunc("/api/insert/async/{qty:[0-9]+}", controllers.InsertASync).Methods("PUT")
	router.HandleFunc("/api/insert", controllers.ClearInserts).Methods("DELETE")
	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	port := "8000"

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
