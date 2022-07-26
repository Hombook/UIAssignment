package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"uiassignment/internal/pkg/auth"
	"uiassignment/internal/pkg/models"

	"gorm.io/gorm"
)

func (h handler) CreateAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	type createAccessTokenRequest struct {
		Acct     string `json:"account"`
		Password string `json:"password"`
	}
	var catRequest createAccessTokenRequest

	err := json.NewDecoder(r.Body).Decode(&catRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user models.Users
	if result := h.DB.Where("acct = ?", catRequest.Acct).First(&user); result.Error != nil {
		log.Println(result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if !auth.IsPasswordMatched(user.Password, catRequest.Password) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accessToken, expiresAt, err := auth.CreateAccessTokenForUser(user.Acct)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type createAccessTokenResponse struct {
		AccessToken string `json:"AccessToken"`
		ExpiresAt   int64  `json:"ExpiresAt"`
	}
	var catResponse createAccessTokenResponse
	catResponse.AccessToken = accessToken
	catResponse.ExpiresAt = expiresAt

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(catResponse)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
