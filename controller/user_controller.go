package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
			return
		}

		userDataJson, err := json.Marshal(userProfile)
		if err != nil {
			log.Println("Error converting profile data to json:", err)
			http.Error(w, "Error converting profile data to json:", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(userDataJson)
	}
}

func RetrieveUserProfileDataControl(us service.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO add this to the other controllers
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}

		userIdStr := r.URL.Query().Get("user-id")
		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			log.Println("Error converting the given userId:", err)
			http.Error(w, fmt.Sprintf("Error converting the given userId: %s", userIdStr), http.StatusBadRequest) 
			return
		}

		profileData, err := us.RetrieveUserProfileData(userId)
		if err != nil {
			log.Println("Error retrieving profile data:", err)
			http.Error(w, "Error retrieving profile data.", http.StatusNotFound)
			return
		}

		profileDataJSON, err := json.Marshal(profileData)
		if err != nil {
			log.Println("Error marshaling profile data:", err)
			http.Error(w, "Error marshaling profile data.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(profileDataJSON)
	}
}
