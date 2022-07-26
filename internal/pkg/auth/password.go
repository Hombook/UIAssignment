package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Hash the password. Returns empty string if an empty password string is given.
func EncryptPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", nil
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func IsPasswordMatched(storedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
