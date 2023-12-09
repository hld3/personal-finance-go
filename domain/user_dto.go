package domain

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserDTO struct {
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	DateOfBirth int64  `json:"dateOfBirth" validate:"required"`
	Password    string `json:"password" validate:"required"`

	// for returning the user profile data.
	UserId       uuid.UUID
	CreationDate int64
}

type UserLoginDTO struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

type UserProfileDTO struct {
	UserDTO  UserDTO
	JWTToken string
}

type UserData struct {
	Validator *validator.Validate
	User      *UserDTO
	Login     *UserLoginDTO
}

func (u *UserData) ValidateUserDTO() error {
	err := u.Validator.Struct(u.User)
	if err != nil {
		log.Printf("User validation failed, %v. UserDTO: %v\n", err, u.User)
		return err
	}
	return nil
}

func (u *UserData) ValidateUserLoginDTO() error {
	err := u.Validator.Struct(u.Login)
	if err != nil {
		log.Printf("Login validation failed, %v. LoginDTO: %v\n", err, u.Login)
		return err
	}
	return nil
}
