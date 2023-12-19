package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hld3/personal-finance-go/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterNewUser(userData *domain.UserData) error {
	args := m.Called(userData)
	return args.Error(0)
}

func (m *MockUserService) ConfirmUserLogin(userData *domain.UserData) (*domain.UserProfileDTO, error) {
	args := m.Called(userData)
	return args.Get(0).(*domain.UserProfileDTO), args.Error(1)
}

func (m *MockUserService) RetrieveUserProfileData(userId uuid.UUID) (*domain.UserDTO, error) {
	args := m.Called(userId)
	return args.Get(0).(*domain.UserDTO), args.Error(1)
}

func (m *MockUserService) UpdateUserProfileData(user *domain.UserDTO) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestRegisterNewUserControl(t *testing.T) {
	validDTOJson, err := json.Marshal(domain.UserDTOBuilder().Build())
	if err != nil {
		t.Error("Error marshaling DTO", err)
	}
	tests := []struct {
		name           string
		userJSON       string
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name:           "Valid DTO",
			userJSON:       string(validDTOJson),
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Malformed JSON",
			userJSON:       `{"someField" "Some Value",}`,
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "UserService error",
			userJSON:       string(validDTOJson),
			mockReturnErr:  errors.New("Some user service error."),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockService := new(MockUserService)
			req, err := http.NewRequest("POST", "/register", bytes.NewBufferString(string(test.userJSON)))
			if err != nil {
				t.Fatal("Error building the request.", err)
			}
			rr := httptest.NewRecorder()

			mockService.On("RegisterNewUser", mock.AnythingOfType("*domain.UserData")).Return(test.mockReturnErr)

			handler := http.HandlerFunc(RegisterNewUserControl(mockService, validator.New()))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("Wrong status code for %v: got %v, expected %v", test.name, status, test.expectedStatus)
			}
		})
	}
}

func TestConfirmUserLoginControl(t *testing.T) {
	validLoginJson, err := json.Marshal(domain.UserLoginDTOBuilder().Build())
	validUserDTO := domain.UserDTOBuilder().Build()
	validUserProfileDTO := domain.UserProfileDTO{UserDTO: validUserDTO, JWTToken: "some_token"}
	if err != nil {
		t.Error("Error marshaling DTO:", err)
	}

	mockService := new(MockUserService)
	req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(string(validLoginJson)))
	if err != nil {
		t.Fatal("Error building request:", err)
	}
	rr := httptest.NewRecorder()

	mockService.On("ConfirmUserLogin", mock.AnythingOfType("*domain.UserData")).Return(&validUserProfileDTO, nil)

	handler := http.HandlerFunc(ConfirmUserLoginControl(mockService, validator.New()))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}
}

func TestRetrieveUserProfileDataControl(t *testing.T) {
	userIdString := uuid.NewString()
	validUserDTO := domain.UserDTOBuilder().Build()

	// Do I need the same UUID? Checked and the answer is yes.
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		t.Fatalf("Error parsing userId: %v, err: %v", userIdString, err)
	}
	validUserDTO.UserId = userId

	mockService := new(MockUserService)
	req, err := http.NewRequest("GET", "/profile/?user-id="+userIdString, nil)
	if err != nil {
		t.Fatal("Error building request:", err)
	}

	rr := httptest.NewRecorder()
	mockService.On("RetrieveUserProfileData", userId).Return(&validUserDTO, nil)

	handler := http.HandlerFunc(RetrieveUserProfileDataControl(mockService))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}
}

func TestUpdateUserProfileDataControl(t *testing.T) {
	userDTO := domain.UserDTOBuilder().Build()
	userDTOByte, err := json.Marshal(userDTO)
	if err != nil {
		t.Fatal("Error marshaling the user DTO:", err)
	}

	mockService := new(MockUserService)
	req, err := http.NewRequest("PUT", "/update", bytes.NewBufferString(string(userDTOByte)))
	if err != nil {
		t.Fatal("Error creating the request:", err)
	}

	rr := httptest.NewRecorder()
	mockService.On("UpdateUserProfileData", &userDTO).Return(nil)

	handler := http.HandlerFunc(UpdateUserProfileDataControl(mockService))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}
}
