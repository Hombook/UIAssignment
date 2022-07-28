package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"uiassignment/internal/pkg/db"
	"uiassignment/internal/pkg/handlers"
	"uiassignment/internal/pkg/middlewares"
	"uiassignment/internal/pkg/websocket"
	"uiassignment/web/pkg/webhandlers"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// @title        uiassignment REST API
// @version      v1
// @description  uiassignment REST service
// @BasePath     /
// @schemes      http
// @tag.name     uiassignment.
func main() {
	DB := db.Init()
	Validator := validator.New()
	hub := websocket.NewHub()
	go hub.Run()
	handler := handlers.New(DB, Validator, hub)

	router := mux.NewRouter()
	router.HandleFunc("/health", handlers.HealthCheckHandler)
	// Websocket demo
	router.HandleFunc("/web/chat", webhandlers.ChatWebHandler)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	// Paths without access control
	subRouter := router.PathPrefix("/api/v1/").Subrouter()
	subRouter.HandleFunc("/accessToken", handler.CreateAccessTokenHandler).Methods(http.MethodPost)
	subRouter.HandleFunc("/users", handler.CreateUserHandler).Methods(http.MethodPost)

	// Paths that requires access token
	accessControledSR := router.PathPrefix("/api/v1/").Subrouter()
	accessControledSR.Use(middlewares.AccessTokenCheckMW())
	accessControledSR.HandleFunc("/users/{account}", handler.GetUserByAccountHandler).Methods(http.MethodGet)
	accessControledSR.HandleFunc("/users", handler.ListUsersHandler).Methods(http.MethodGet)

	// Paths that requires resource owner access
	ownerAccessSR := router.PathPrefix("/api/v1/").Subrouter()
	ownerAccessSR.Use(middlewares.OwnerAccessCheckMW())
	ownerAccessSR.HandleFunc("/users/{account}", handler.DeleteUserByAccountHandler).Methods(http.MethodDelete)
	ownerAccessSR.HandleFunc("/users/{account}", handler.UpdateUserHandler).Methods(http.MethodPatch)

	// TLS
	enableTls := true
	sslCrtPath := "/app/uiassignment/tls/tls.crt"
	sslKeyPath := "/app/uiassignment/tls/tls.key"
	if _, err := os.Stat(sslCrtPath); errors.Is(err, os.ErrNotExist) {
		enableTls = false
	}
	if _, err := os.Stat(sslKeyPath); errors.Is(err, os.ErrNotExist) {
		enableTls = false
	}

	var err error
	if enableTls {
		log.Println("TLS certificates found. Starting TLS server at port 443")
		err = http.ListenAndServeTLS(":443", sslCrtPath, sslKeyPath, router)
	} else {
		log.Println("TLS certificates not found. Starting server at port 80")
		err = http.ListenAndServe(":80", router)
	}
	if err != nil {
		log.Panicln("FATAL - HTTP server startup failure: ", err)
	}
}
