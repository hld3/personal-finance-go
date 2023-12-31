package main

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hld3/personal-finance-go/controller"
	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load env file:", err)
	}
	db := database.ConnectDB()
	defer db.Close()

	dbManager := database.SQLManager{DB: db} // implementation of UserDatabase interface
	userService := service.UserService{UDBI: &dbManager} // implementation of UserServiceInterface
	newValidator := validator.New()
	
	http.HandleFunc("/user/register", controller.RegisterNewUserControl(&userService, newValidator))
	http.HandleFunc("/user/login", controller.ConfirmUserLoginControl(&userService, newValidator))
	http.HandleFunc("/user/profile", controller.RetrieveUserProfileDataControl(&userService))
	http.HandleFunc("/user/update", controller.UpdateUserProfileDataControl(&userService))
	log.Fatal(http.ListenAndServe(":8083", nil))
}

