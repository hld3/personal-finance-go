package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/hld3/personal-finance-go/domain"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("There was an error loading .env:", err)
	}
}

// This also essentially tests AddTransaction.
func TestGetTransactionIntegration(t *testing.T) {
	// Open a connection to the database.
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		t.Fatal("Missing DB_URL env variable")
	}
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		t.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()
	udb := SQLManager{DB: db}

	// Save the transaction.
	expected := domain.TransactionModelBuilder().Build()
	err = udb.AddTransaction(&expected)
	if err != nil {
		t.Fatal("Error saving transaction:", err)
	}

	// Retrieve the transaction.
	result, err := udb.GetTransaction(expected.TransactionId)
	if err != nil {
		t.Fatal("Error retrieving the transaction:", err)
	}

	if result.UserId != expected.UserId ||
		result.TransactionId != expected.TransactionId ||
		result.CategoryId != expected.CategoryId ||
		result.Amount != expected.Amount ||
		result.Date != expected.Date ||
		result.Description != expected.Description ||
		result.CreatedAt != expected.CreatedAt ||
		result.UpdatedAt != expected.UpdatedAt ||
		result.Type != expected.Type ||
		result.PaymentMethod != expected.PaymentMethod ||
		result.Status != expected.Status {
		t.Errorf("Transaction data does not match. want %v, got %v", expected, result)
	}
}
