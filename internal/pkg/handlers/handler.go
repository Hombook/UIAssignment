package handlers

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type handler struct {
	DB        *gorm.DB
	Validator *validator.Validate
}

type CommonResponse struct {
	Message string `json:"message"`
}

func New(db *gorm.DB, validator *validator.Validate) handler {
	return handler{db, validator}
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
