package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/domain"
)

func RegisterNewUser(db database.UserDatabase, userDTO *domain.UserDTO) error {
	// validate the user fields.
	err := userDTO.ValidateUserDTO()
	if err != nil {
		return err // ValidateUser will log the error
	}

	userModel := convertDTOToModel(userDTO)

	// save the user to the database
	err = db.AddNewUser(&userModel)
	if err != nil {
		return err
	}
	return nil
}

func convertDTOToModel(from *domain.UserDTO) domain.UserModel {
	return domain.UserModel{
		UserId:       uuid.New(),
		FirstName:    from.FirstName,
		LastName:     from.LastName,
		Phone:        from.Phone,
		Email:        from.Email,
		DateOfBirth:  from.DateOfBirth,
		PasswordHash: from.Password, //TODO need to hash
		CreationDate: time.Now().UnixMilli(),
	}
}
