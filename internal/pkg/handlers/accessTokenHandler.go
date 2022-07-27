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

// swagger:handlers createAccessTokenRequest
type createAccessTokenRequest struct {
	// User account
	// example: myAccount100
	// required: true
	Acct string `json:"account"`
	// Password of the given account
	// example: my@pass100Word
	// required: true
	Password string `json:"password"`
}

// swagger:handlers createAccessTokenResponse
type createAccessTokenResponse struct {
	// Access token
	AccessToken string `json:"AccessToken"`
	// Unix timestamp of when the token expires
	ExpiresAt int64 `json:"ExpiresAt"`
}

// CreateAccessTokenHandler godoc
// @Description Create user access token
// @Tags accessToken
// @Produce application/json
// @Param Body body createAccessTokenRequest true "User login credentials"
// @Success 200 {object} createAccessTokenResponse
// @Failure 400 "Invalid user account credentials"
// @Failure 500 "Internal error caused by DB connection issue or JSON parsing failure"
// @Router /v1/accessToken [post]
func (h handler) CreateAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
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
