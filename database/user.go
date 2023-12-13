package database

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/hld3/personal-finance-go/domain"
)

type UserDatabaseInterface interface {
	AddNewUser(user *domain.UserModel) error
	RetrieveUserByEmail(email string) (domain.UserModel, error)
	RetrieveUserByUserId(userId uuid.UUID) (domain.UserModel, error)
}

type SQLManager struct {
	DB *sql.DB
}

func (db *SQLManager) AddNewUser(user *domain.UserModel) error {
	stmt := `insert into user_model (user_id, first_name, last_name, email, phone, date_of_birth, creation_date, password_hash) values (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.DB.Exec(stmt, user.UserId, user.FirstName, user.LastName, user.Email, user.Phone, user.DateOfBirth, user.CreationDate, user.PasswordHash)
	if err != nil {
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
		return user, err
	}
	return user, nil
}

func (db *SQLManager) RetrieveUserByUserId(userId uuid.UUID) (domain.UserModel, error) {
	stmt := `select user_id, first_name, last_name, email, phone, date_of_birth, creation_date, password_hash from user_model where user_id = ?`
	var user domain.UserModel
	err := db.DB.QueryRow(stmt, userId).Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.DateOfBirth, &user.CreationDate, &user.PasswordHash)
	if err != nil {
		return user, err
	}
	return user, nil
}
