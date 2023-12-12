package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hld3/personal-finance-go/domain"
	"github.com/hld3/personal-finance-go/service"
)

func RegisterNewUserControl(us service.UserServiceInterface, validator *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body:", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var user domain.UserDTO
		err = json.Unmarshal(bodyBytes, &user)
		if err != nil {
			log.Println("Error converting to user DTO:", err)
			http.Error(w, "Error converting to user DTO", http.StatusBadRequest)
			return
		}

		userData := domain.UserData{User: &user, Validator: validator}
		err = us.RegisterNewUser(&userData)
		if err != nil {
			log.Println("Error registering a new user:", err)
			http.Error(w, "Error registering a new user", http.StatusInternalServerError)
			return
		}
	}
}

func ConfirmUserLoginControl(us service.UserServiceInterface, validator *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body:", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var userLogin domain.UserLoginDTO
		err = json.Unmarshal(bodyBytes, &userLogin)
		if err != nil {
			log.Println("Error converting to user DTO:", err)
			http.Error(w, "Error converting to user DTO", http.StatusBadRequest)
			return
		}

		userData := domain.UserData{Login: &userLogin, Validator: validator}
		userProfile, err := us.ConfirmUserLogin(&userData)
		if err != nil {
			log.Println("Error logging in:", err)
			http.Error(w, "Error loggin in:", http.StatusUnauthorized)
		}

		userDataJson, err := json.Marshal(userProfile)
		if err != nil {
			log.Println("Error converting profile data to json:", err)
			http.Error(w, "Error converting profile data to json:", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(userDataJson)
	}
}
