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
	"github.com/elnerribeiro/go-ws-db-auth/utils"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(router)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: handler,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[main] Error while executing server: %s", err)
		}
	}()

	fmt.Println("[main] Server started on port 8000")

	<-done
	fmt.Println("[main] Server Stopped")
	utils.FinalizeDB()
	utils.FinalizeLog()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("[main] Server Shutdown Failed:%+v", err)
	}
	fmt.Println("[main] Server Exited Properly")
}
