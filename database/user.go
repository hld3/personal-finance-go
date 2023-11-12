package database

import (
	"log"

	"github.com/hld3/personal-finance-go/domain"
)

func AddNewUser(user *domain.UserModel) error {
	stmt := `insert into user_model (user_id, first_name, last_name, email, phone, date_of_birth, creation_date, password_hash) values (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := Conn.Exec(stmt, user.UserId, user.FirstName, user.LastName, user.Email, user.Phone, user.DateOfBirth, user.CreationDate, user.PasswordHash)
	if err != nil {
		log.Println("Error saving user:", err)
		return err
	}
	log.Println("New user successfully created:", user.UserId) 
	return nil
}
