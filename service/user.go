package service

import (
	"time"

	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/domain"
)

func RegisterNewUser(user *domain.UserModel) error {
	// validate the user fields.
	err := user.ValidateUser()
	if err != nil {
		return err // ValidateUser will log the error
	}
	// add creation date
	user.CreateDate = time.Now()

	// TODO hash password

	// save the user to the database
	err = database.AddNewUser(user)
	if err != nil {
		return err
	}
	return nil
}
