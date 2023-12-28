package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hld3/personal-finance-go/domain"
)

func TestAddTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Error creating SQL stub", err)
	}
	defer db.Close()
	udb := SQLManager{DB: db}

	tm := domain.TransactionModelBuilder().Build()
	mock.ExpectExec("insert into transaction").
		WithArgs(tm.UserId, tm.TransactionId, tm.CategoryId, tm.Amount, tm.Date, tm.Description, tm.CreatedAt, tm.UpdatedAt, tm.Type, tm.PaymentMethod, tm.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = udb.AddTransaction(&tm)
	if err != nil {
		t.Fatal("Error saving transaction:", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("Expectations were not met:", err)
	}
}

func TestGetTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Error creating SQL stub", err)
	}
	defer db.Close()
	udb := SQLManager{DB: db}

	transaction := domain.TransactionModelBuilder().Build()
	row := sqlmock.NewRows([]string{"user_id", "transaction_id", "category_id", "amount", "date", "description", "created_at", "updated_at", "type", "payment_method", "status"}).
		AddRow(transaction.UserId, transaction.TransactionId, transaction.CategoryId, transaction.Amount, transaction.Date, transaction.Description, transaction.CreatedAt, transaction.UpdatedAt, transaction.Type, transaction.PaymentMethod, transaction.Status)

	mock.ExpectQuery("select (.+) from transaction_model where transaction_id = ?").
		WithArgs(transaction.TransactionId).
		WillReturnRows(row)

	gotTransaction, err := udb.GetTransaction(transaction.TransactionId)
	if err != nil {
		t.Fatal("Error retrieving transaction:", err)
	}
	if gotTransaction != transaction {
		t.Fatalf("Retrieved transaction does not match expected, got %v, want %v", gotTransaction, transaction)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfulfilled expectations: %v", err)
	}
}
