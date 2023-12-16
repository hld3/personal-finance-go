package utility

import (
	"log"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("There was an error loading .env:", err)
	}
}

func TestCreateJWTToken(t *testing.T) {
	userId := uuid.NewString()
	token, err := CreateJWTToken(userId, time.Hour)
	if err != nil {
		t.Error("Failed to create token:", err)
	}

	parsedToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if claims, ok := parsedToken.Claims.(*jwt.StandardClaims); ok && parsedToken.Valid {
		if claims.Subject != userId {
			t.Errorf("Expected user id %s, got %s", userId, claims.Subject)
		}
		if claims.Issuer != "auth0" {
			t.Errorf("Expected issuer %s, got %s", "auth0", claims.Issuer)
		}
		if claims.ExpiresAt < time.Now().Unix() {
			t.Error("Token is expired")
		}
	} else {
		t.Error("There was an error parsing the token:", err)
	}
}
