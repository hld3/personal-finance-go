package utility

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("There was an error loading .env:", err)
	}
}

func CreateJWTToken(userId string, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	claims := &jwt.StandardClaims{
		Issuer:    os.Getenv("JWT_ISSUER"),
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
