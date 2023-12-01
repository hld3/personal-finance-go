package database

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hld3/personal-finance-go/domain"
)

func TestAddNewUserIntegration(t *testing.T) {
	// Open a connection to the database.
	db, err := sql.Open("mysql", "finance:finance@tcp(127.0.0.1:3307)/finance")
	if err != nil {
		t.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()
	udb := SQLManager{DB: db}

	// Attempt to save the expectedUser.
	dob := time.Date(1980, time.December, 19, 0, 0, 0, 0, time.UTC).UnixMilli()
	createStart := time.Now().UnixMilli()
	expectedUser := domain.UserModelBuilder().Build()
	expectedUser.DateOfBirth = dob
	createEnd := time.Now().UnixMilli()

	err = udb.AddNewUser(&expectedUser)
	if err != nil {
		t.Fatal("Error adding a new user:", err)
	}
	// Remove the user that was just added.
	defer cleanUpDatabase(&udb, expectedUser.UserId)

	// Check the database for the user.
	var savedUser domain.UserModel

	stmt := `select user_id, first_name, last_name, email, phone, password_hash, creation_date, date_of_birth from user_model where user_id = ?`
	err = udb.DB.QueryRow(stmt, expectedUser.UserId).Scan( &savedUser.UserId, &savedUser.FirstName, &savedUser.LastName, &savedUser.Email, 
		&savedUser.Phone, &savedUser.PasswordHash, &savedUser.CreationDate, &savedUser.DateOfBirth)
	if err != nil {
		t.Fatal("Failed to retrieve user:", err)
	}

	if !compareUserData(expectedUser, savedUser, dob, createStart, createEnd) {
		t.Fatalf("Saved user data: %v does not match expected user data: %v", savedUser, expectedUser)
	}
}

func TestRetrieveUserByEmailIntegration(t *testing.T) {
	// Open a connection to the database.
	db, err := sql.Open("mysql", "finance:finance@tcp(127.0.0.1:3307)/finance")
	if err != nil {
		t.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()
	udb := SQLManager{DB: db}
	
	// Save a user to find
	user := domain.UserModelBuilder().Build()
	_, err = db.Exec("insert into user_model (user_id, first_name, last_name, email, phone, date_of_birth, creation_date, password_hash) values (?, ?, ?, ?, ?, ?, ?, ?)", user.UserId, user.FirstName, user.LastName, user.Email, user.Phone, user.DateOfBirth, user.CreationDate, user.PasswordHash)
	if err != nil {
		t.Error("Error saving test user:", err)
	}

	um, err := udb.RetrieveUserByEmail(user.Email)
	if err != nil {
		t.Error("Error retrieving user:", err)
	}

	if (um != user) {
		t.Errorf("User data does not match, got %v, want %v", um, user)
	}
}

func compareUserData(a, b domain.UserModel, dob, createStart, createEnd int64) bool {
	return a.FirstName == b.FirstName &&
		a.LastName == b.LastName &&
		a.Email == b.Email &&
		a.PasswordHash == b.PasswordHash &&
		b.DateOfBirth == dob &&
		(b.CreationDate >= createStart && b.CreationDate <= createEnd)
}

func cleanUpDatabase(db *SQLManager, userId uuid.UUID) {
	stmt := `delete from user_model where user_id = ?`
	_, err := db.DB.Exec(stmt, userId)
	if err != nil {
		log.Println("Error removing user:", userId)
	}
}
