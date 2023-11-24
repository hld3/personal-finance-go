package main

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hld3/personal-finance-go/controller"
	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/domain"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()
	dbManager := database.SQLManager{DB: db}

	domain.Validate = validator.New()
	
	http.HandleFunc("/register", controller.RegisterNewUser(&dbManager))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

