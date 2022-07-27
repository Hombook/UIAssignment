package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"uiassignment/internal/pkg/auth"
	"uiassignment/internal/pkg/db"
	"uiassignment/internal/pkg/models"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

var queryDecoder = schema.NewDecoder()

func (h handler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	type listUserQuery struct {
		FullName string `schema:"fullName" validate:"omitempty,min=1,max=50"`
		Limit    int    `schema:"limit" validate:"omitempty,gte=5,lte=100"`
		Page     int    `schema:"page" validate:"omitempty,gt=0"`
		OrderBy  string `schema:"orderBy" validate:"omitempty,oneof=acct fullname"`
		Order    string `schema:"order" validate:"omitempty,oneof=asc desc"`
	}
	var luQuery listUserQuery

	err := queryDecoder.Decode(&luQuery, r.URL.Query())
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Validator.Struct(luQuery)
	if err != nil {
		var errResponse CommonResponse
		errResponse.Message = ValidatorErrorMessageBuilder(err)

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errResponse)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	pagination := db.Pagination{Limit: luQuery.Limit, Page: luQuery.Page}

	var order string
	if len(luQuery.OrderBy) > 0 {
		if len(luQuery.Order) > 0 {
			order = luQuery.OrderBy + " " + luQuery.Order
		} else {
			order = luQuery.OrderBy + " asc"
		}
	}

	var usersList = []models.UsersList{}
	users := models.Users{FullName: luQuery.FullName}
	if result := h.DB.Model(&users).
		Scopes(db.Paginate(users, &pagination, h.DB)).
		Order(order).Find(&usersList); result.Error != nil {
		log.Println(result.Error)
	}
	pagination.Rows = usersList

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(pagination)
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
		Acct     string `json:"account" validate:"required,alphanum"`
		Password string `json:"password" validate:"required,alphanum,min=6,max=40"`
		FullName string `json:"fullName" validate:"required,min=1,max=50"`
	}
	var cuRequest createUserRequest

	err := json.NewDecoder(r.Body).Decode(&cuRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Validator.Struct(cuRequest)
	if err != nil {
		var errResponse CommonResponse
		errResponse.Message = ValidatorErrorMessageBuilder(err)

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errResponse)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	encryptedPassword, err := auth.EncryptPassword(cuRequest.Password)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if result := h.DB.Create(&models.Users{
		Acct:     cuRequest.Acct,
		Password: encryptedPassword,
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

func (h handler) DeleteUserByAccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	account := vars["account"]

	if result := h.DB.Delete(&models.Users{Acct: account}); result.Error != nil {
		log.Println(result.Error)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (h handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	type updateUserRequest struct {
		Password string `json:"password" validate:"omitempty,alphanum,min=6,max=40"`
		FullName string `json:"fullName" validate:"omitempty,min=1,max=50"`
	}
	var uuRequest updateUserRequest

	err := json.NewDecoder(r.Body).Decode(&uuRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Validator.Struct(uuRequest)
	if err != nil {
		var errResponse CommonResponse
		errResponse.Message = ValidatorErrorMessageBuilder(err)

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errResponse)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	var encryptedPassword string
	encryptedPassword, err = auth.EncryptPassword(uuRequest.Password)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	account := vars["account"]

	user := models.Users{Acct: account}
	if result := h.DB.Model(&user).Updates(models.Users{
		Password: encryptedPassword,
		FullName: uuRequest.FullName}); result.Error != nil {
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
}
