package service

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/domain"
	"github.com/hld3/personal-finance-go/utility"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	RegisterNewUser(userData *domain.UserData) error
	ConfirmUserLogin(userData *domain.UserData) (*domain.UserProfileDTO, error)
}

type UserService struct {
	UDBI database.UserDatabaseInterface
}

func (us *UserService) RegisterNewUser(userData *domain.UserData) error {
	// validate the user fields.
	err := userData.ValidateUserDTO()
	if err != nil {
		return err // ValidateUser will log the error
	}

	userModel := convertDTOToModel(userData.User)

	// save the user to the database
	err = us.UDBI.AddNewUser(&userModel)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) ConfirmUserLogin(userData *domain.UserData) (*domain.UserProfileDTO, error) {
	err := userData.ValidateUserLoginDTO()
	if err != nil {
		return nil, err
	}
	// retrieve the user by email
	userModel, err := us.UDBI.RetrieveUserByEmail(userData.Login.Email)
	if err != nil {
		return nil, err
	}
	// compare password hashes. CompareHashAndPassword returns nil if they match.
	err = bcrypt.CompareHashAndPassword([]byte(userModel.PasswordHash), []byte(userData.Login.Password))
	if err != nil {
		return nil, err
	}
	// create JWT token.
	token, err := utility.CreateJWTToken(userModel.UserId.String(), time.Hour) // TODO env variable for time?
	if err != nil {
		return nil, err
	}
	// update the user profile DTO
	userDTO := convertModelToDTO(&userModel)
	userProfile := domain.UserProfileDTO{UserDTO: userDTO, JWTToken: token}
	// return result
	return &userProfile, nil // when to use a pointer or not? this time was in order to return nil for the UserProfileDTO.
}

func convertDTOToModel(from *domain.UserDTO) domain.UserModel {
	userId := uuid.New()
	hashedPass := HashPassword(from.Password, userId)
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

func convertModelToDTO(from *domain.UserModel) domain.UserDTO {
	return domain.UserDTO{
		UserId: from.UserId,
		FirstName: from.FirstName,
		LastName: from.LastName,
		Email: from.Email,
		Phone: from.Phone,
		DateOfBirth: from.DateOfBirth,
		CreationDate: from.CreationDate,
	}
}

// TODO move this to utility some day.
func HashPassword(password string, userId uuid.UUID) string {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %s, for user: %v", password, userId)
		return password
	}
	return string(hashedPass)
}
