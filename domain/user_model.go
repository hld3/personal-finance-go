package domain

import (
	"log"

	"github.com/google/uuid"
)

type UserModel struct {
	UserId       uuid.UUID
	FirstName    string `json:"firstName" validate:"required"`
	LastName     string `json:"lastName" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone" validate:"required"`
	DateOfBirth  int64  `json:"dateOfBirth" validate:"required"`
	PasswordHash string `json:"password" validate:"required"` // hashed on creation
	CreationDate int64  // Added on creation
}

func (user *UserModel) ValidateUser() error {
	err := Validate.Struct(user)
	if err != nil {
		log.Printf("User validation failed, %v. UserModel: %v\n", err, user)
		return err
	}
	return nil
}
