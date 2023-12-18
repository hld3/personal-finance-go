package service

import (
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hld3/personal-finance-go/domain"
	"golang.org/x/crypto/bcrypt"
)

var pw = "super_private_password"

type StubDatabase struct{}

func (m *StubDatabase) AddNewUser(user *domain.UserModel) error {
	return nil
}

// This stub doesn't allow me to compare the UserProfileDTO result at the end of the tests.
// However the test are better with a stub, primarily when it comes to validation.
// The test to compare the UserProfileDTO result is in the intergration test.
func (m *StubDatabase) RetrieveUserByEmail(email string) (domain.UserModel, error) {
	pw := HashPassword(pw, uuid.New())
	return domain.UserModelBuilder().WithPasswordHash(pw).Build(), nil
}

func (m *StubDatabase) RetrieveUserByUserId(userId uuid.UUID) (domain.UserModel, error) {
	return domain.UserModelBuilder().Build(), nil
}

func (m *StubDatabase) UpdateUserByUserId(user *domain.UserDTO) error {
	return nil
}

func TestRegisterNewUser_Validation(t *testing.T) {
	stubDB := new(StubDatabase)
	userService := UserService{UDBI: stubDB}

	tests := []struct {
		name        string
		user        domain.UserDTO
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "Valid user data",
			user:        domain.UserDTOBuilder().Build(),
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "Invalid user data",
			user:        domain.UserDTOBuilder().WithFirstName("").Build(),
			wantErr:     true,
			expectedErr: "Key: 'UserDTO.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userData := domain.UserData{User: &tt.user, Validator: validator.New()}
			err := userService.RegisterNewUser(&userData)

			if tt.wantErr {
				if err == nil {
					t.Errorf("RegisterNewUser: %s, expected %v, got nil", tt.name, tt.expectedErr)
				} else {
					if !strings.Contains(err.Error(), tt.expectedErr) {
						t.Errorf("RegisterNewUser: %s, expected %s, got %s", tt.name, tt.expectedErr, err.Error())
					}
				}
			} else if err != nil {
				t.Errorf("RegisterNewUser: %s, unexpected error %s", tt.name, err)
			}
		})
	}
}

func TestConfirmUserLogin(t *testing.T) {
	stubDB := new(StubDatabase)
	userService := UserService{UDBI: stubDB}

	tests := []struct {
		name          string
		loginDTO      domain.UserLoginDTO
		wantError     bool
		expectedError string
	}{
		{
			name:          "Successful login",
			loginDTO:      domain.UserLoginDTOBuilder().WithPassword(pw).Build(),
			wantError:     false,
			expectedError: "",
		},
		{
			name:          "Validation failure",
			loginDTO:      domain.UserLoginDTOBuilder().WithPassword("").Build(),
			wantError:     true,
			expectedError: "Key: 'UserLoginDTO.Password' Error:Field validation for 'Password' failed on the 'required' tag",
		},
		{
			name:          "Wrong password",
			loginDTO:      domain.UserLoginDTOBuilder().WithPassword("wrong_password").Build(),
			wantError:     true,
			expectedError: "hashedPassword is not the hash of the given password",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userData := domain.UserData{Login: &test.loginDTO, Validator: validator.New()}

			result, err := userService.ConfirmUserLogin(&userData)

			if test.wantError && !strings.Contains(err.Error(), test.expectedError) {
				t.Errorf("Error did not match expected error, got %v want %v", err, test.expectedError)
			} else if err != nil && !test.wantError {
				t.Error("There was an unexpected error", err)
			} else if err == nil && result.JWTToken == "" {
				t.Error("Token expected but was missing, got", result)
			}
		})
	}
}

// This function would only pass or fail because of what I tell it to return.
// It doesn't seem useful as a unit test.
// func TestRetrieveUserProfileData(t *testing.T) { }

// This test would be useless with a stub as the result is dependant on the returned error.
// func TestUpdateUserProfileData(t *testing.T) { }

func TestConvertDTOToModel(t *testing.T) {
	fromDTO := domain.UserDTOBuilder().Build()
	toModel := convertDTOToModel(&fromDTO)
	err := bcrypt.CompareHashAndPassword([]byte(toModel.PasswordHash), []byte(fromDTO.Password))
	if err != nil {
		t.Error("Password hash does not match")
	}

	if toModel.FirstName != fromDTO.FirstName ||
		toModel.LastName != fromDTO.LastName ||
		toModel.Phone != fromDTO.Phone ||
		toModel.Email != fromDTO.Email ||
		toModel.DateOfBirth != fromDTO.DateOfBirth ||
		toModel.UserId == uuid.Nil ||
		toModel.CreationDate == 0 {
		t.Error("Conversion of UserDTO to UserModel failed")
	}
}

func TestConvertModelToDTO(t *testing.T) {
	fromModel := domain.UserModelBuilder().Build()
	toDTO := convertModelToDTO(&fromModel)

	if toDTO.UserId != fromModel.UserId ||
		toDTO.FirstName != fromModel.FirstName ||
		toDTO.LastName != fromModel.LastName ||
		toDTO.Phone != fromModel.Phone ||
		toDTO.Email != fromModel.Email ||
		toDTO.DateOfBirth != fromModel.DateOfBirth ||
		toDTO.CreationDate != fromModel.CreationDate {
		t.Error("Conversion of UserModel to UserDTO failed")
	}
}
