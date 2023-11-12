package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/domain"
)

func RegisterNewUser(user *domain.UserModel) error {
	// validate the user fields.
	err := user.ValidateUser()
	if err != nil {
		return err // ValidateUser will log the error
	}
	// add new user id
	user.UserId = uuid.New()
	// add creation date
	user.CreationDate = time.Now().UnixMilli()

	// TODO hash password

	// save the user to the database
	err = database.AddNewUser(user)
	if err != nil {
		return err
	}
	return nil
}
