package domain

import (
	"time"

	gen "github.com/pallinder/go-randomdata"
)

type UserDTOBuild struct{
	firstName string
}

func UserDTOBuilder() *UserDTOBuild {
	return &UserDTOBuild{
		firstName: gen.FirstName(1),
	}
}

func (b *UserDTOBuild) Build() UserDTO {
	return UserDTO{
		FirstName: b.firstName,
		LastName: gen.LastName(),
		Email: gen.Email(),
		Phone: gen.PhoneNumber(),
		DateOfBirth: time.Now().UnixMilli(),
		Password: gen.Letters(15),
	}
}

func (b *UserDTOBuild) WithFirstName(name string) *UserDTOBuild {
	b.firstName = name
	return b
}
