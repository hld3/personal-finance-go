package service

import (
	"database/sql"
	"log"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/domain"
	_ "github.com/mattn/go-sqlite3"
)

// Test setup using SQLite.
func TestRegisterNewUser_Integration(t *testing.T) {
	db := setUpDatabase()
	defer db.Close()
	udb := database.SQLManager{DB: db}
	userService := UserService{UDBI: &udb}

	tests := []struct {
		name    string
		userDTO domain.UserDTO
		wantErr bool
	}{
		{
			name:    "Valid user",
			userDTO: domain.UserDTOBuilder().Build(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userData := domain.UserData{User: &tt.userDTO, Validator: validator.New()}
			err := userService.RegisterNewUser(&userData)
			if (err != nil) && !tt.wantErr {
				t.Errorf("RegisterNewUser: %s, unexpected error: %v", tt.name, err)
			}
		})
	}
}

func TestConfirmUserLogin_Integration(t *testing.T) {
	db := setUpDatabase()
	defer db.Close()
	udb := database.SQLManager{DB: db}
	userService := UserService{UDBI: &udb}

	tests := []struct {
		name          string
		loginDTO      domain.UserLoginDTO
		wantErr       bool
		expectedError string
	}{
		{
			name:          "Successful login",
			loginDTO:      domain.UserLoginDTOBuilder().Build(),
			wantErr:       false,
			expectedError: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userData := domain.UserData{Login: &test.loginDTO, Validator: validator.New()}
			err := saveUserToDatabase(test.loginDTO, db)
			if err != nil {
				t.Error("Error saving user to the database.", err)
			}

			result, err := userService.ConfirmUserLogin(&userData)
			if err != nil {
				t.Error("Error logging user in:", err)
			}

			if test.wantErr && !strings.Contains(err.Error(), test.expectedError) {
				t.Errorf("Error did not match the expected error, got %v want %s", err, test.expectedError)
			} else if err != nil && !test.wantErr {
				t.Error("There was an unexpected error", err)
			} else if err == nil && result == "" {
				t.Error("Token expected but was missing, got", result)
			}
		})
	}
}

func saveUserToDatabase(loginDTO domain.UserLoginDTO, db *sql.DB) error {
	ping := db.Ping()
	if ping != nil {
		return ping
	}
	hashPW := HashPassword(loginDTO.Password, uuid.New())
	um := domain.UserModelBuilder().WithPasswordHash(hashPW).Build()
	stmt := `insert into user_model (user_id, first_name, last_name, email, phone, date_of_birth, password_hash, creation_date) values (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(stmt, um.UserId, um.FirstName, um.LastName, loginDTO.Email, um.Phone, um.DateOfBirth, um.PasswordHash, um.CreationDate)
	if err != nil {
		return err
	}
	return nil
}

func setUpDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	stmt := `create table user_model (
		id integer primary key autoincrement,
		user_id text not null,
		first_name text not null,
		last_name text not null,
		email text not null,
		phone text not null,
		password_hash text not null,
		date_of_birth integer not null,
		creation_date integer not null
	)`

	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal("There was an error creating user_model table:", err)
	}

	return db
}
