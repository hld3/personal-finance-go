package main

import "github.com/hld3/personal-finance-go/database"

func main() {
	database.ConnectDB()
	defer database.CloseDB()
}
