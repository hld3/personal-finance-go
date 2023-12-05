package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	var err error
	db, err := sql.Open("mysql", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection not found", err)
	}

	log.Println("Successfully connected to the database.")
	return db
}

