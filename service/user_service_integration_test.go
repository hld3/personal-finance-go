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
			userModel, err := saveUserLoginToDatabase(test.loginDTO, db)
			if err != nil {
				t.Fatal("Error saving user to the database.", err)
			}

			result, err := userService.ConfirmUserLogin(&userData)
			if err != nil {
				t.Fatal("Error logging user in:", err)
			}

			if test.wantErr && !strings.Contains(err.Error(), test.expectedError) {
				t.Fatalf("Error did not match the expected error, got %v want %s", err, test.expectedError)
			} else if err != nil && !test.wantErr {
				t.Fatal("There was an unexpected error:", err)
			} else if err == nil && result.JWTToken == "" {
				t.Fatal("Token expected but was missing.")
			} else if !confirmUserData(userModel, result.UserDTO) {
				t.Fatalf("The user data does not match. got %v want %v", result.UserDTO, userModel)
			}
		})
	}
}

// TODO needs testing and maybe a test that fails to retrieve the data.
func TestRetrieveUserProfileData_Integration(t *testing.T) {
	db := setUpDatabase()
	defer db.Close()
	udb := database.SQLManager{DB: db}
	userService := UserService{UDBI: &udb}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "User found",
			wantErr: false,
		},
		{
			name:    "User not found",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// save the user
			userModel := domain.UserModelBuilder().Build()
			err := saveUserModelToDatabase(userModel, db)
			if err != nil {
				t.Fatal("Error saving user model to the database:", err)
			}

			if !test.wantErr {
				// pass the user id
				profileData, err := userService.RetrieveUserProfileData(userModel.UserId)
				if err != nil {
					t.Fatal("Error retrieving profile data:", err)
				}

				// confirm the correct user
				if profileData.FirstName != userModel.FirstName ||
					profileData.LastName != userModel.LastName ||
					profileData.Email != userModel.Email ||
					profileData.Phone != userModel.Phone ||
					profileData.DateOfBirth != userModel.DateOfBirth ||
					profileData.CreationDate != userModel.CreationDate {
					t.Fatalf("The profile data does not match the expected data, got %v, want %v", profileData, userModel)
				}
			} else {
				_, err := userService.RetrieveUserProfileData(uuid.New())
				if err == nil {
					t.Fatal("Expected error when no user profile data is found.")
				}
			}
		})
	}

}

func TestUpdateUserProfileData_Integration(t *testing.T) {
	db := setUpDatabase()
	defer db.Close()
	udb := database.SQLManager{DB: db}
	userService := UserService{UDBI: &udb}

	// User is saved to the database
	savedUser := domain.UserModelBuilder().Build()
	err := saveUserModelToDatabase(savedUser, db)
	if err != nil {
		t.Error("Error saving user model to the database:", err)
	}

	// The data to update the user is sent
	expected := domain.UserDTOBuilder().WithUserId(savedUser.UserId).Build()
	err = userService.UpdateUserProfileData(&expected)
	if err != nil {
		t.Error("Error updating user data:", err)
	}

	// Get the updated user.
	updatedUser, err := udb.RetrieveUserByUserId(savedUser.UserId)
	if err != nil {
		t.Error("Error retrieving user by userId:", err)
	}

	// Check the user for the updated data.
	if updatedUser.FirstName != expected.FirstName ||
		updatedUser.LastName != expected.LastName ||
		updatedUser.Email != expected.Email ||
		updatedUser.Phone != expected.Phone ||
		updatedUser.DateOfBirth != expected.DateOfBirth {
		t.Errorf("The user was not updated as expected, want %v, got %v.", expected, updatedUser)
	}
}

func confirmUserData(fromUser domain.UserModel, toUser domain.UserDTO) bool {
	return fromUser.UserId == toUser.UserId &&
		fromUser.FirstName == toUser.FirstName &&
		fromUser.LastName == toUser.LastName &&
		fromUser.Phone == toUser.Phone &&
		fromUser.Email == toUser.Email &&
		fromUser.DateOfBirth == toUser.DateOfBirth &&
		fromUser.CreationDate == toUser.CreationDate
}

func saveUserLoginToDatabase(loginDTO domain.UserLoginDTO, db *sql.DB) (domain.UserModel, error) {
	hashPW := HashPassword(loginDTO.Password, uuid.New())
	um := domain.UserModelBuilder().WithPasswordHash(hashPW).Build()
	um.Email = loginDTO.Email //TODO do I want to create a builder function?

	stmt := `insert into user_model (user_id, first_name, last_name, email, phone, date_of_birth, password_hash, creation_date) values (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(stmt, um.UserId, um.FirstName, um.LastName, loginDTO.Email, um.Phone, um.DateOfBirth, um.PasswordHash, um.CreationDate)
	if err != nil {
		return um, err
	}
	return um, nil
}

func saveUserModelToDatabase(um domain.UserModel, db *sql.DB) error {
	stmt := `insert into user_model (user_id, first_name, last_name, email, phone, date_of_birth, password_hash, creation_date) values (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(stmt, um.UserId, um.FirstName, um.LastName, um.Email, um.Phone, um.DateOfBirth, um.PasswordHash, um.CreationDate)
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
