package main

import (
	"log"
	"net/http"
	"uiassignment/internal/pkg/db"
	"uiassignment/internal/pkg/handlers"
	"uiassignment/internal/pkg/middlewares"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func main() {
	DB := db.Init()
	Validator := validator.New()
	handler := handlers.New(DB, Validator)

	router := mux.NewRouter()
	router.HandleFunc("/health", handlers.HealthCheckHandler)

	// Paths without access control
	subRouter := router.PathPrefix("/v1/").Subrouter()
	subRouter.HandleFunc("/accessToken", handler.CreateAccessTokenHandler).Methods(http.MethodPost)
	subRouter.HandleFunc("/users", handler.CreateUserHandler).Methods(http.MethodPost)

	// Paths that requires access token
	accessControledSR := router.PathPrefix("/v1/").Subrouter()
	accessControledSR.Use(middlewares.AccessTokenCheckMW())
	accessControledSR.HandleFunc("/users/{account}", handler.GetUserByAccountHandler).Methods(http.MethodGet)
	accessControledSR.HandleFunc("/users", handler.ListUsersHandler).Methods(http.MethodGet).Queries("fullName", "{fullName}")
	accessControledSR.HandleFunc("/users", handler.ListUsersHandler).Methods(http.MethodGet)

	// Paths that requires resource owner access
	ownerAccessSR := router.PathPrefix("/v1/").Subrouter()
	ownerAccessSR.Use(middlewares.OwnerAccessCheckMW())
	ownerAccessSR.HandleFunc("/users/{account}", handler.DeleteUserByAccountHandler).Methods(http.MethodDelete)
	ownerAccessSR.HandleFunc("/users/{account}", handler.UpdateUserHandler).Methods(http.MethodPatch)

	err := http.ListenAndServe(":80", router)
	if err != nil {
		log.Panicln("FATAL - HTTP server startup failure: ", err)
	}
}
