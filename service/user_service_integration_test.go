package service

import (
	"database/sql"
	"log"
	"testing"

	"github.com/go-playground/validator/v10"
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

