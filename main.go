package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elnerribeiro/go-ws-db-auth/app"
	"github.com/elnerribeiro/go-ws-db-auth/controllers"
	"github.com/elnerribeiro/go-ws-db-auth/repositories"

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

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	fmt.Println("Server started on port 8000")

	<-done
	fmt.Println("Server Stopped")
	repositories.FinalizeDB()
	repositories.FinalizeLog()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown Failed:%+v", err)
	}
	fmt.Println("Server Exited Properly")
}
