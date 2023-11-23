package service

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/domain"
	"golang.org/x/crypto/bcrypt"
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
	userId := uuid.New()
	hashedPass := hashPassword(from.Password, userId)
	return domain.UserModel{
		UserId:       userId,
		FirstName:    from.FirstName,
		LastName:     from.LastName,
		Phone:        from.Phone,
		Email:        from.Email,
		DateOfBirth:  from.DateOfBirth,
		PasswordHash: hashedPass,
		CreationDate: time.Now().UnixMilli(),
	}
}

func hashPassword(password string, userId uuid.UUID) string {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %s, for user: %v", password, userId)
		return password
	}
	return string(hashedPass)
}
