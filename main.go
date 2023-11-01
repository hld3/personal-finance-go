package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/domain"
	"github.com/hld3/personal-finance-go/service"
)

func main() {
	database.ConnectDB()
	defer database.CloseDB()

	domain.Validate = validator.New()
	
	http.HandleFunc("/register", registerNewUser)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerNewUser(w http.ResponseWriter, r *http.Request) {
	// convert request body to UserModel
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		return
	}

	var user domain.UserModel
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		log.Println("Error converting to user:", err)
		return
	}
	log.Println(user)

	err = service.RegisterNewUser(&user)
	if err != nil {
		//TODO idk the control flow I need here yet.
	}
}
