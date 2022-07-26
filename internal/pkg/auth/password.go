package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
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
