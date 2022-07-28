package handlers

import (
	"strings"
	"uiassignment/internal/pkg/websocket"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type handler struct {
	DB        *gorm.DB
	Validator *validator.Validate
	Hub       *websocket.Hub
}

// swagger:handlers CommonResponse
// @Description A single element JSON for returning a human readable message.
type CommonResponse struct {
	// Human readable message
	Message string `json:"message"`
}

func New(db *gorm.DB, validator *validator.Validate, hub *websocket.Hub) handler {
	return handler{db, validator, hub}
}

// Helper function for generating message from ValidationErrors.
func ValidatorErrorMessageBuilder(err error) string {
	var errorMessage strings.Builder
	var lastIndex = len(err.(validator.ValidationErrors)) - 1
	for i, err := range err.(validator.ValidationErrors) {
		errorMessage.WriteString(err.StructField())
		errorMessage.WriteString(" ")
		errorMessage.WriteString(err.Tag())
		// Not meeting the size range, print the suggestion.
		if len(err.Param()) > 0 {
			errorMessage.WriteString(":")
			errorMessage.WriteString(err.Param())
		}
		if i < lastIndex {
			errorMessage.WriteString(", ")
		}
	}
	return errorMessage.String()
}
