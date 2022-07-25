package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"uiassignment/internal/pkg/models"

	"github.com/gorilla/mux"
)

type ListUsersResp struct {
	Users []string `json:"users"`
}

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

	response := &ListUsersResp{Users: userAcctList}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
