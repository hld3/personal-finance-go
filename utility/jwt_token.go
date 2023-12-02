package utility

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret")

func CreateJWTToken(userId string, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	claims := &jwt.StandardClaims{
		Issuer:    "auth0",
		Subject:   userId,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
