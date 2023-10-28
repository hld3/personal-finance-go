package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var conn *sql.DB

func ConnectDB() {
	var err error
	conn, err = sql.Open("mysql", "finance:finance@tcp(127.0.0.1:3306)/finance")
	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatal("Database connection not found", err)
	}

	log.Println("Successfully connected to the database.")
}

func CloseDB() {
	conn.Close()
}
