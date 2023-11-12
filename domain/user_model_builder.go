package domain

import (
	"time"

	"github.com/google/uuid"
	gen "github.com/pallinder/go-randomdata"
)

func BuildUser() *UserModel {
	return &UserModel{
		UserId: uuid.New(), 
		FirstName: gen.FirstName(1), 
		LastName: gen.LastName(), 
		Email: gen.Email(),
		Phone: gen.PhoneNumber(),
		DateOfBirth: time.Now().UnixMilli(),
		CreationDate: time.Now().UnixMilli(),
		PasswordHash: gen.Letters(15),
	}
}
