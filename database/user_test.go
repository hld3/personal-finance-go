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
	udb := SQLManager{DB: db}

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

func TestRetrieveUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Error creating SQL stub", err)
	}
	defer db.Close()
	udb := SQLManager{DB: db}

	// expected row
	user := domain.UserModelBuilder().Build()
	rows := sqlmock.NewRows([]string{"user_id", "first_name", "last_name", "email", "phone", "date_of_birth", "creation_date", "password_hash"}).
		AddRow(user.UserId, user.FirstName, user.LastName, user.Email, user.Phone, user.DateOfBirth, user.CreationDate, user.PasswordHash)

	mock.ExpectQuery("select (.+) from user_model where email = ?").
		WithArgs(user.Email).
		WillReturnRows(rows)

	gotUser, err := udb.RetrieveUserByEmail(user.Email)
	if err != nil {
		t.Error("Error retrieving user.", err)
	}
	if user != gotUser {
		t.Errorf("Expected %v, got %v", user, gotUser)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
