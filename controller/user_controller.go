package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/domain"
	"github.com/hld3/personal-finance-go/service"
)

func RegisterNewUser(db database.UserDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body:", err)
			return
		}

		var user domain.UserDTO
		err = json.Unmarshal(bodyBytes, &user)
		if err != nil {
			log.Println("Error converting to user DTO:", err)
			return
		}

		err = service.RegisterNewUser(db, &user)
		if err != nil {
			log.Println("Error registering a new user:", err)
			return
		}
	}
}
