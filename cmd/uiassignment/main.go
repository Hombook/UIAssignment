package main

import (
	"log"
	"net/http"
	"uiassignment/internal/pkg/db"
	"uiassignment/internal/pkg/handlers"
	"uiassignment/internal/pkg/middlewares"

	"github.com/gorilla/mux"
)

func main() {
	DB := db.Init()
	handler := handlers.New(DB)

	router := mux.NewRouter()
	router.HandleFunc("/health", handlers.HealthCheckHandler)

	// Paths without access control
	//subRouter := router.PathPrefix("/v1/").Subrouter()

	// Paths that requires access token
	accessControledSR := router.PathPrefix("/v1/").Subrouter()
	accessControledSR.Use(middlewares.AccessTokenCheckMW())
	accessControledSR.HandleFunc("/users", handler.ListUsersHandler).Methods(http.MethodGet)

	err := http.ListenAndServe(":80", router)
	if err != nil {
		log.Panicln("FATAL - HTTP server startup failure: ", err)
	}
}
