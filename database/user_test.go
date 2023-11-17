package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hld3/personal-finance-go/domain"
)

func TestAddNewUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Error creating SQL stub", err)
	}
	defer db.Close()
	udb := DB{db}

	user := domain.UserModelBuilder().Build() 
	mock.ExpectExec("insert into user").WithArgs(user.UserId, user.FirstName, user.LastName, user.Email, user.Phone, user.DateOfBirth, user.CreationDate, user.PasswordHash).WillReturnResult(sqlmock.NewResult(1, 1))

	err = udb.AddNewUser(&user)
	if err != nil {
		t.Error("Error saving user", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("Expectations were not met", err)
	}
}
