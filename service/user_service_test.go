package service

import (
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hld3/personal-finance-go/domain"
	"github.com/stretchr/testify/mock"
)

type StubDatabase struct {
	mock.Mock
}

func (m *StubDatabase) AddNewUser(user *domain.UserModel) error {
	return nil
}

func TestRegisterNewUser(t *testing.T) {
	stubDB := new(StubDatabase)
	domain.Validate = validator.New()

	tests := []struct {
		name        string
		user        domain.UserDTO
		wantErr     bool
		expectedErr string
	}{
		{
			name: "Valid user data",
			user: domain.UserDTOBuilder().Build(),
			wantErr: false,
			expectedErr: "",
		},
		{
			name: "Invalid user data",
			user: domain.UserDTOBuilder().WithFirstName("").Build(),
			wantErr: true,
			expectedErr: "Key: 'UserDTO.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RegisterNewUser(stubDB, &tt.user)

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
