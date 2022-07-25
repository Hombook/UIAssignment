package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"uiassignment/internal/pkg/models"

	"github.com/gorilla/mux"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func (h handler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fullName := vars["fullName"]

	var users []models.Users
	if result := h.DB.Where(&models.Users{FullName: fullName}).Find(&users); result.Error != nil {
		log.Println(result.Error)
	}

	var userAcctList []string
	for _, user := range users {
		userAcctList = append(userAcctList, user.Acct)
	}

	type listUsersResp struct {
		Users []string `json:"users"`
	}
	response := &listUsersResp{Users: userAcctList}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h handler) GetUserByAccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	account := vars["account"]

	var user models.Users
	if result := h.DB.Where(&models.Users{Acct: account}).First(&user); result.Error != nil {
		log.Println(result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		Acct     string `json:"account"`
		Password string `json:"password"`
		FullName string `json:"fullName"`
	}
	var cuRequest createUserRequest

	err := json.NewDecoder(r.Body).Decode(&cuRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if result := h.DB.Create(&models.Users{
		Acct:     cuRequest.Acct,
		Password: cuRequest.Password,
		FullName: cuRequest.FullName}); result.Error != nil {
		log.Println(result.Error)

		var duplicateEntryError = &pgconn.PgError{Code: "23505"}
		if errors.As(result.Error, &duplicateEntryError) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}
