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

// ListUsersHandler godoc
// @Description Get a list of user accounts and names with paging
// @Tags user
// @Produce application/json
// @Param X-Accesstoken header string true "Access token"
// @Param fullName query string false "Filter by user's full name"
// @Param limit query int false "Max items per page(min=5, max=100, default=5)"
// @Param page query int false "Requested page"
// @Param orderBy query string false "Select attribute to sort the list(acct: account, fullname: full name)"
// @Param order query string false "Sort order(asc: ascending, desc: descending )"
// @Success 200 {object} db.Pagination
// @Failure 400 {object} CommonResponse "Invalid query parameter"
// @Failure 401 "Missing valid acces token for accessing this resource"
// @Failure 500 "Internal error caused by DB connection issue or JSON parsing failure"
// @Router /v1/users [get]
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
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
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

// GetUserByAccountHandler godoc
// @Description Get user details by the selected account
// @Tags user
// @Produce application/json
// @Param X-Accesstoken header string true "Access token"
// @Param account path string true "User account"
// @Success 200 {object} models.Users
// @Failure 401 "Missing valid acces token for accessing this resource"
// @Failure 404 "Account doesn't exist"
// @Failure 500 "Internal error caused by DB connection issue or JSON parsing failure"
// @Router /v1/users/{account} [get]
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

// swagger:handlers createUserRequest
// @Description JSON request body for creating user
type createUserRequest struct {
	// User account, alphanumeric only
	Acct string `json:"account" validate:"required,alphanum"`
	// Password, alphanumeric only(Length: min=6, max=40)
	Password string `json:"password" validate:"required,alphanum,min=6,max=40"`
	// User's full name(Length: min=1, max=50)
	FullName string `json:"fullName" validate:"required,min=1,max=50"`
}

// CreateUserHandler godoc
// @Description Create user
// @Tags user
// @Produce application/json
// @Param Body body createUserRequest true "Data for creating the user"
// @Success 201 "User created"
// @Failure 400 {object} CommonResponse "Invalid request body or duplicated account"
// @Failure 500 "Internal error caused by DB connection issue or JSON parsing failure"
// @Router /v1/users [post]
func (h handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
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

// DeleteUserByAccountHandler godoc
// @Description Delete user by the given account
// @Tags user
// @Produce application/json
// @Param X-Accesstoken header string true "Access token"
// @Param account path string true "User account"
// @Success 200 "Successfully deleted the user"
// @Failure 401 "Missing valid acces token for accessing this resource"
// @Failure 403 "Current token owner has no right to access this resource"
// @Failure 500 "Internal error caused by DB connection issue"
// @Router /v1/users/{account} [delete]
func (h handler) DeleteUserByAccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	account := vars["account"]

	if result := h.DB.Delete(&models.Users{Acct: account}); result.Error != nil {
		log.Println(result.Error)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

// swagger:handlers updateUserRequest
// @Description JSON request body for updating user
type updateUserRequest struct {
	// Password, alphanumeric only(Length: min=6, max=40)
	Password string `json:"password" validate:"omitempty,alphanum,min=6,max=40"`
	// User's full name(Length: min=1, max=50)
	FullName string `json:"fullName" validate:"omitempty,min=1,max=50"`
}

// UpdateUserHandler godoc
// @Description Update selected account's user data
// @Tags user
// @Produce application/json
// @Param X-Accesstoken header string true "Access token"
// @Param account path string true "User account"
// @Param Body body updateUserRequest true "Data for updating the user"
// @Success 200 "Successfully updated the user"
// @Failure 400 {object} CommonResponse "Invalid request body"
// @Failure 401 "Missing valid acces token for accessing this resource"
// @Failure 403 "Current token owner has no right to access this resource"
// @Failure 500 "Internal error caused by DB connection issue"
// @Router /v1/users/{account} [patch]
func (h handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(http.StatusOK)
}
