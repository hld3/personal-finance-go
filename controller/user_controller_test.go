package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hld3/personal-finance-go/domain"
	"github.com/stretchr/testify/mock"
)

type StubSQLManager struct {
	mock.Mock
}

func (db *StubSQLManager) AddNewUser(user *domain.UserModel) error {
	return nil
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterNewUser(userData *domain.UserData) error {
	args := m.Called(userData)
	return args.Error(0)
}

func TestUserController_Register(t *testing.T) {
	mockService := new(MockUserService)

	userDTO := domain.UserDTOBuilder().Build()
	userJSON, err := json.Marshal(userDTO)
	if err != nil {
		t.Fatal("Error marshaling the user DTO.", err)
	}

	req, err := http.NewRequest("POST", "/register", bytes.NewBufferString(string(userJSON)))
	if err != nil {
		t.Fatal("Error building the request.", err)
	}
	rr := httptest.NewRecorder()

	mockService.On("RegisterNewUser", mock.AnythingOfType("*domain.UserData")).Return(nil)

	handler := http.HandlerFunc(RegisterNewUserControl(mockService, validator.New()))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, expected %v", status, http.StatusOK)
	}
}
