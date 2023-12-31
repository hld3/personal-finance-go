package domain

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TransactionData struct {
	Validator   *validator.Validate
	Transaction TransactionDTO
}

type TransactionDTO struct {
	UserId        uuid.UUID         `json:"userId" validate:"required"`
	TransactionId uuid.UUID         `json:"transactionId" validate:"required"`
	CategoryId    int64             `json:"categoryId" validate:"required"`
	Amount        float64           `json:"amount" validate:"required"`
	Date          int64             `json:"date" validate:"required"`
	Description   string            `json:"description" validate:"required"`
	CreatedAt     int64             `json:"createdAt" validate:"required"`
	UpdatedAt     int64             `json:"updatedAt" validate:"required"`
	Type          TransactionType   `json:"type"`
	PaymentMethod TransactionMethod `json:"paymentMethod"`
	Status        TransactionStatus `json:"status"`
}

func (t *TransactionData) ValidateTransaction() error {
	err := t.Validator.Struct(t.Transaction)
	if err != nil {
		log.Printf("Transaction validation failed, %v. ValidationDTO: %v\n", err, t.Transaction)
		return err
	}
	return nil
}
