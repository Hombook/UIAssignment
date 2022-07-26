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

// Validates access token and returns owner account.
func IsAccessTokenValid(accessToken string) (isTokenValid bool, tokenOwnerAccount string) {
	claims := &Claims{}
	// Parse the token
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("Received a token with invalid signature: ", accessToken)
			return false, ""
		}
		log.Println("Fail to parse token: ", accessToken)
		return false, ""
	}
	if !token.Valid {
		log.Println("Invalid token: ", accessToken)
		return false, ""
	}

	return true, claims.Account
}

// Creates access token for the given user account.
func CreateAccessTokenForUser(userAccount string) (accessToken string, expiresAt int64, err error) {
	expiresAt = time.Now().Add(24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		Account: userAccount,
	})

	accessToken, err = token.SignedString(jwtSecretKey)
	if err != nil {
		return "", 0, err
	}

	return
}
