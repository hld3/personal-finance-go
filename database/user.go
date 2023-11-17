package database

import (
	"database/sql"
	"log"

	"github.com/hld3/personal-finance-go/domain"
)

type UserDatabase interface {
	AddNewUser(user *domain.UserModel) error
}

type DB struct {
	*sql.DB
}

func (db *DB) AddNewUser(user *domain.UserModel) error {
	stmt := `insert into user_model (user_id, first_name, last_name, email, phone, date_of_birth, creation_date, password_hash) values (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(stmt, user.UserId, user.FirstName, user.LastName, user.Email, user.Phone, user.DateOfBirth, user.CreationDate, user.PasswordHash)
	if err != nil {
		log.Println("Error saving user:", err)
		return err
	}
	log.Println("New user successfully created:", user.UserId) 
	return nil
}
