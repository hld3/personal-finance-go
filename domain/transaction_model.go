package domain

import "github.com/google/uuid"

type TransactionType int

const (
	INCOME TransactionType = iota
	EXPENSE
)

type TransactionMethod int

const (
	CASH TransactionMethod = iota
	CREDIT_CARD
	BANK_TRANSFER
)

type TransactionStatus int

const (
	PENDING TransactionStatus = iota
	CLEARED
	CANCELLED
)

type TransactionModel struct {
	UserId        uuid.UUID
	CategoryId    int64
	Amount        float64
	Date          int64
	Description   string
	CreatedAt     int64
	UpdatedAt     int64
	Type          TransactionType
	PaymentMethod TransactionMethod
	Status        TransactionStatus
}

type CategoryModel struct {
	UserId      uuid.UUID
	Name        string
	Description string
}
