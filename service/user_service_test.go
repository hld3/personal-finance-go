package service

import (
	"testing"

	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/domain"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) AddNewUser(user *domain.UserModel) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestRegisterNewUser(t *testing.T) {
	mockDB := new(MockDatabase)
	database.AddNewUser = mockDB.AddNewUser
}
