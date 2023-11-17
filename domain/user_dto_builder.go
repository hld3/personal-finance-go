package domain

import (
	"time"

	gen "github.com/pallinder/go-randomdata"
)

type UserDTOBuild struct{}

func UserDTOBuilder() *UserDTOBuild {
	return &UserDTOBuild{}
}

func (b *UserDTOBuild) Build() UserDTO {
	return UserDTO{
		FirstName: gen.FirstName(0),
		LastName: gen.LastName(),
		Email: gen.Email(),
		Phone: gen.PhoneNumber(),
		DateOfBirth: time.Now().UnixMilli(),
		Password: gen.Letters(15),
	}
}
