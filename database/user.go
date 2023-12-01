package database

import (
	"database/sql"
	"log"

	"github.com/hld3/personal-finance-go/domain"
)

type UserDatabaseInterface interface {
	AddNewUser(user *domain.UserModel) error
}

type SQLManager struct {
	DB *sql.DB
}

func (db *SQLManager) AddNewUser(user *domain.UserModel) error {
	stmt := `insert into user_model (user_id, first_name, last_name, email, phone, date_of_birth, creation_date, password_hash) values (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.DB.Exec(stmt, user.UserId, user.FirstName, user.LastName, user.Email, user.Phone, user.DateOfBirth, user.CreationDate, user.PasswordHash)
	if err != nil {
		log.Println("Error saving user:", err)
		return err
	}
	log.Println("New user successfully created:", user.UserId)
	return nil
}

func (db *SQLManager) RetrieveUserByEmail(email string) (domain.UserModel, error) {
	stmt := `select user_id, first_name, last_name, email, phone, date_of_birth, creation_date, password_hash from user_model where email = ?`
	var user domain.UserModel
	err := db.DB.QueryRow(stmt, email).Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.DateOfBirth, &user.CreationDate, &user.PasswordHash)
	if err != nil {
		log.Println("Error retrieving user:", err)
		return user, err
	}
	return user, nil
}
