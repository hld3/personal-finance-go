package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/domain"
	"github.com/hld3/personal-finance-go/service"
	_ "github.com/mattn/go-sqlite3"
)

func TestUserControllerIntegration_Register(t *testing.T) {
	// Setup SQLite database.
	db := setUpDatabase()
	defer db.Close()

	// Create the service tested, the controller.
	udb := database.SQLManager{DB: db}
	userService := service.UserService{UDBI: &udb}
	handler := RegisterNewUserControl(&userService, validator.New())

	// A valid user DTO as a json.
	validDTO := domain.UserDTOBuilder().Build()
	validJSON, err := json.Marshal(validDTO)
	if err != nil {
		t.Error("Error marshaling valid user DTO to JSON.", err)
	}

	// Invalid user DTO as a json.
	invalidDTO := domain.UserDTOBuilder().WithFirstName("").Build()
	invalidJSON, err := json.Marshal(invalidDTO)
	if err != nil {
		t.Error("Error marshaling invalid user DTO to JSON.", err)
	}
	fmt.Println("boom", string(invalidJSON))

	tests := []struct {
		name           string
		userJSON       string
		wantErr        bool
		expectedStatus int
	}{
		{
			name:           "Valid user",
			userJSON:       string(validJSON),
			wantErr:        false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Malformed JSON",
			userJSON:       `{test: "test"}`,
			wantErr:        false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid UserDTO (Missing FirstName)",
			userJSON:       string(invalidJSON),
			wantErr:        false,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(test.userJSON))
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("handler returned the wrong status code: got %v, want %v", rr.Code, test.expectedStatus)
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
