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

func TestUserController_Register(t *testing.T) {
	mockDB := new(StubSQLManager)
	domain.Validate = validator.New()

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
	handler := http.HandlerFunc(RegisterNewUser(mockDB))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, expected %v", status, http.StatusOK)
	}
}
