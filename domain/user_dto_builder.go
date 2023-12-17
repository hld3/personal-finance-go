package domain

import (
	"time"

	"github.com/google/uuid"
	gen "github.com/pallinder/go-randomdata"
)

type UserDTOStruct struct {
	userId    uuid.UUID
	firstName string
}

func UserDTOBuilder() *UserDTOStruct {
	return &UserDTOStruct{
		userId:    uuid.New(),
		firstName: gen.FirstName(1),
	}
}

func (b *UserDTOStruct) Build() UserDTO {
	return UserDTO{
		UserId:      b.userId,
		FirstName:   b.firstName,
		LastName:    gen.LastName(),
		Email:       gen.Email(),
		Phone:       gen.PhoneNumber(),
		DateOfBirth: time.Now().UnixMilli(),
		Password:    gen.Letters(15),
	}
}

func (b *UserDTOStruct) WithFirstName(name string) *UserDTOStruct {
	b.firstName = name
	return b
}

func (b *UserDTOStruct) WithUserId(userId uuid.UUID) *UserDTOStruct {
	b.userId = userId
	return b
}

type UserLoginDTOStruct struct {
	email    string
	password string
}

func UserLoginDTOBuilder() *UserLoginDTOStruct {
	return &UserLoginDTOStruct{
		email:    gen.Email(),
		password: gen.Letters(15),
	}
}

func (b *UserLoginDTOStruct) Build() UserLoginDTO {
	return UserLoginDTO{
		Email:    b.email,
		Password: b.password,
	}
}

func (b *UserLoginDTOStruct) WithEmailAndPassword(email string, pw string) *UserLoginDTOStruct {
	b.email = email
	b.password = pw
	return b
}

func (b *UserLoginDTOStruct) WithPassword(pw string) *UserLoginDTOStruct {
	b.password = pw
	return b
}
