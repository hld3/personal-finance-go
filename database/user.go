package database

import (
	"log"

	"github.com/hld3/personal-finance-go/domain"
)

func AddNewUser(user *domain.UserModel) error {
	// user id made automatically on creation
	stmt := `insert into user (first_name, last_name, email, phone, date_of_birth, creation_date, password_hash) values (?, ?, ?, ?, ?, ?, ?)`
	_, err := Conn.Exec(stmt, user.FirstName, user.LastName, user.Email, user.Phone, user.DateOfBirth, user.CreateDate, user.PasswordHash)
	if err != nil {
		log.Println("Error saving user:", err)
		return err
	}
	log.Println("New user successfully created:", user.UserId) // TODO change to set the user id in code and not the database.
	return nil
}
