package domain

import "github.com/google/uuid"

type UserModel struct {
	UserId       uuid.UUID
	FirstName    string
	LastName     string
	Email        string
	Phone        string
	DateOfBirth  int64
	PasswordHash string // hashed on creation
	CreationDate int64  // Added on creation
}
