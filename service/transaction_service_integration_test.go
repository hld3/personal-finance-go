package service

import (
	"database/sql"
	"log"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/domain"
)

func TestServiceAddGetTransaction(t *testing.T) {
	db := setUpTransactionModel()
	defer db.Close()
	udb := database.SQLManager{DB: db}
	transactionService := TransactionService{UDBI: &udb}

	transactionDTO := domain.TransactionDTOBuilder().Build()
	transactionData := domain.TransactionData{Transaction: transactionDTO, Validator: validator.New()}

	err := transactionService.AddTransaction(&transactionData)
	if err != nil {
		t.Fatal("Error adding the transaction:", err)
	}

	transactionModel, err := transactionService.GetTransaction(transactionDTO.TransactionId)
	if err != nil {
		t.Fatal("Error retrieving the transaction:", err)
	}

	if transactionDTO.UserId != transactionModel.UserId ||
		transactionDTO.TransactionId != transactionModel.TransactionId ||
		transactionDTO.CategoryId != transactionModel.CategoryId ||
		transactionDTO.Amount != transactionModel.Amount ||
		transactionDTO.Date != transactionModel.Date ||
		transactionDTO.Description != transactionModel.Description ||
		transactionDTO.CreatedAt != transactionModel.CreatedAt ||
		transactionDTO.UpdatedAt != transactionModel.UpdatedAt ||
		transactionDTO.Type != transactionModel.Type ||
		transactionDTO.PaymentMethod != transactionModel.PaymentMethod ||
		transactionDTO.Status != transactionModel.Status {
		t.Fatalf("The saved data is not the expected data: got %v, want %v", transactionModel, transactionDTO)
	}
}

func setUpTransactionModel() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	stmt := `create table transaction_model (
		id integer primary key autoincrement,
		user_id text not null,
		transaction_id text not null,
		category_id integer not null,
		amount float not null,
		date integer not null,
		description text not null,
		created_at integer not null,
		updated_at integer not null,
		type integer not null,
		payment_method integer not null,
		status integer not null
	)`

	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal("There was an error creating transaction_model table:", err)
	}

	return db
}
