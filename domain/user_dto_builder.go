package domain

import (
	"time"

	gen "github.com/pallinder/go-randomdata"
)

type UserDTOStruct struct {
	firstName string
}

func UserDTOBuilder() *UserDTOStruct {
	return &UserDTOStruct{
		firstName: gen.FirstName(1),
	}
}

func (b *UserDTOStruct) Build() UserDTO {
	return UserDTO{
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

type UserLoginDTOStruct struct {
	password string
}

func UserLoginDTOBuilder() *UserLoginDTOStruct {
	return &UserLoginDTOStruct{
		password: gen.Letters(15),
	}
}

func (b *UserLoginDTOStruct) Build() UserLoginDTO {
	return UserLoginDTO{
		Email:    gen.Email(),
		Password: b.password,
	}
}

func (b *UserLoginDTOStruct) WithPassword(pw string) *UserLoginDTOStruct {
	b.password = pw
	return b
}
