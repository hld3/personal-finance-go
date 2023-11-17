package domain

import "log"

type UserDTO struct {
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	DateOfBirth int64  `json:"dateOfBirth" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

func (user *UserDTO) ValidateUserDTO() error {
	err := Validate.Struct(user)
	if err != nil {
		log.Printf("User validation failed, %v. UserDTO: %v\n", err, user)
	}
	return nil
}
