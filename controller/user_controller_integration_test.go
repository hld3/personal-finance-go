package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/domain"
	"github.com/hld3/personal-finance-go/service"
	_ "github.com/mattn/go-sqlite3"
)

func TestRegisterNewUserControl_Integration(t *testing.T) {
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

	tests := []struct {
		name           string
		userJSON       string
		wantErr        bool // TODO I could probably use this to make better tests at some point.
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

func TestConfirmUserLoginControl_Integration(t *testing.T) {
	// Setup SQLite database.
	db := setUpDatabase()
	defer db.Close()

	// Create the service tested.
	udb := database.SQLManager{DB: db}
	userService := service.UserService{UDBI: &udb}
	handler := ConfirmUserLoginControl(&userService, validator.New())

	tests := []struct {
		name         string
		password     string
		email        string
		userLogin    domain.UserLoginDTO
		hasError     bool
		expectedErr  string
		expectedCode int
	}{
		{
			name:         "Valid login",
			password:     "secret_password",
			email:        "email@email.com",
			userLogin:    domain.UserLoginDTOBuilder().WithEmailAndPassword("email@email.com", "secret_password").Build(),
			hasError:     false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Wrong email",
			password:     "secret_password",
			email:        "wrong@email.com",
			userLogin:    domain.UserLoginDTOBuilder().WithEmailAndPassword("another@email.com", "secret_password").Build(),
			hasError:     true,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Wrong password",
			password:     "wrong_password",
			email:        "your@email.com",
			userLogin:    domain.UserLoginDTOBuilder().WithEmailAndPassword("your@email.com", "secret_password").Build(),
			hasError:     true,
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// password hash needs to match the login password
			passHash := service.HashPassword(test.password, uuid.New())
			// user model with a matching email and password hash
			expectedUserModel := domain.UserModelBuilder().WithPasswordHash(passHash).Build()
			expectedUserModel.Email = test.email // TODO builder function for email? dejavu.
			// save the user model.
			saveUserModel(db, expectedUserModel)

			userLoginJSON, err := json.Marshal(test.userLogin)
			if err != nil {
				t.Error("Error marshaling user login to JSON:", err)
			}
			req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(string(userLoginJSON)))
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedCode {
				t.Errorf("Wrong status code returned, got %v, want %v", status, test.expectedCode)
			}
			if !test.hasError {
				resBody, err := io.ReadAll(rr.Body)
				if err != nil {
					t.Fatal("Error reading response body:", err)
				}

				var actualResBody domain.UserProfileDTO
				if err := json.Unmarshal(resBody, &actualResBody); err != nil {
					t.Fatal("Error converting response json to a user profile:", err)
				}

				if actualResBody.JWTToken == "" {
					t.Fatal("The reponse is missing the JWTToken.")
				}
				resDTO := actualResBody.UserDTO
				// confirm the response UserDTO to the UserModel saved to the database.
				if resDTO.UserId != expectedUserModel.UserId ||
					resDTO.FirstName != expectedUserModel.FirstName ||
					resDTO.LastName != expectedUserModel.LastName ||
					resDTO.Email != expectedUserModel.Email ||
					resDTO.Phone != expectedUserModel.Phone ||
					resDTO.CreationDate != expectedUserModel.CreationDate ||
					resDTO.DateOfBirth != expectedUserModel.DateOfBirth {
					t.Fatal("The response UserDTO does not match the expected data.")
				}
			}
		})
	}
}

func TestRetrieveUserProfileDataControl_Integration(t *testing.T) {
	// Setup SQLite database.
	db := setUpDatabase()
	defer db.Close()

	// Create the service tested.
	udb := database.SQLManager{DB: db}
	userService := service.UserService{UDBI: &udb}
	handler := RetrieveUserProfileDataControl(&userService)

	tests := []struct {
		name         string
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "User found",
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "User not found",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create the user.
			userIdString := uuid.NewString()
			userId, err := uuid.Parse(userIdString)
			if err != nil {
				t.Errorf("Error parsing the userId: %s", userIdString)
			}

			var userModel domain.UserModel
			if test.wantErr {
				userModel = domain.UserModelBuilder().Build()
			} else {
				userModel = domain.UserModelBuilder().WithUserId(userId).Build()
			}
			saveUserModel(db, userModel)

			req, _ := http.NewRequest("GET", "/profile/?user-id="+userIdString, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedCode {
				t.Errorf("Wrong status code returned, got %v, want %v", status, test.expectedCode)
			}
			if !test.wantErr {
				resBody, err := io.ReadAll(rr.Body)
				if err != nil {
					t.Fatal("Error reading response body:", err)
				}

				var resDTO domain.UserDTO
				if err := json.Unmarshal(resBody, &resDTO); err != nil {
					t.Fatal("Error converting response json to a user DTO:", err)
				}

				// confirm the response UserDTO to the UserModel saved to the database.
				if resDTO.UserId != userModel.UserId ||
					resDTO.FirstName != userModel.FirstName ||
					resDTO.LastName != userModel.LastName ||
					resDTO.Email != userModel.Email ||
					resDTO.Phone != userModel.Phone ||
					resDTO.CreationDate != userModel.CreationDate ||
					resDTO.DateOfBirth != userModel.DateOfBirth {
					t.Fatal("The response UserDTO does not match the expected data.")
				}
			}
		})
	}
}

func TestUpdateUserProfileDataControl_Integration(t *testing.T) {
	// Setup SQLite database.
	db := setUpDatabase()
	defer db.Close()

	// Create the service tested.
	udb := database.SQLManager{DB: db}
	userService := service.UserService{UDBI: &udb}
	handler := UpdateUserProfileDataControl(&userService)

	// User to udpate
	userModel := domain.UserModelBuilder().Build()
	saveUserModel(db, userModel)

	// User update data
	userDTO := domain.UserDTOBuilder().WithUserId(userModel.UserId).Build()

	// Build the request
	userDTOBytes, err := json.Marshal(userDTO)
	if err != nil {
		t.Fatal("Error marshaling the user DTO:", err)
	}
	req, err := http.NewRequest("GET", "/update", bytes.NewBufferString(string(userDTOBytes)))
	if err != nil {
		t.Fatal("Error building the request:", err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	updatedModel, err := userService.RetrieveUserProfileData(userModel.UserId)
	if err != nil {
		t.Fatal("Error retrieving updated profile data:", err)
	}

	if status := rr.Code; status == http.StatusOK {
		t.Fatalf("Wrong status code returned, got %d, want %d", status, http.StatusOK)
	}
	if updatedModel.FirstName != userDTO.FirstName ||
		updatedModel.LastName != userDTO.LastName ||
		updatedModel.Email != userDTO.Email ||
		updatedModel.Phone != userDTO.Phone ||
		updatedModel.DateOfBirth != userDTO.DateOfBirth {
		t.Fatalf("User profile data was not saved correctly. got %v, want %v", updatedModel, userDTO)
	}
}

func saveUserModel(db *sql.DB, um domain.UserModel) {
	stmt := `insert into user_model (user_id, first_name, last_name, email, phone, password_hash, date_of_birth, creation_date) values (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(stmt, um.UserId, um.FirstName, um.LastName, um.Email, um.Phone, um.PasswordHash, um.DateOfBirth, um.CreationDate)
	if err != nil {
		log.Fatal("Failed to save user model.")
	}
	log.Println("-- TEST -- User successfully saved to the database.")
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
