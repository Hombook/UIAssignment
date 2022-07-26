package auth

import (
	_ "embed"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

//go:embed resources/jwt.key
var jwtSecretKey []byte

type Claims struct {
	Account string `json:"acct"`
	jwt.StandardClaims
}

func IsAccessTokenValid(accessToken string) bool {
	claims := &Claims{}
	// Parse the token
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("Received a token with invalid signature: ", accessToken)
			return false
		}
		log.Println("Fail to parse token: ", accessToken)
		return false
	}
	if !token.Valid {
		log.Println("Invalid token: ", accessToken)
		return false
	}

	return true
}

func CreateAccessTokenForUser(userAccount string) (string, int64, error) {
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		Account: userAccount,
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}
