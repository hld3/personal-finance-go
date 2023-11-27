package domain

import (
	"log"

	"github.com/go-playground/validator/v10"
)

type UserDTO struct {
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	DateOfBirth int64  `json:"dateOfBirth" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type UserData struct {
	Validator *validator.Validate
	User      *UserDTO
}

func (u *UserData) ValidateUserDTO() error {
	err := u.Validator.Struct(u.User)
	if err != nil {
		log.Printf("User validation failed, %v. UserDTO: %v\n", err, u.User)
		return err
	}
	return nil
}
