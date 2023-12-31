package database

import (
	"log"

	"github.com/google/uuid"
	"github.com/hld3/personal-finance-go/domain"
)

type TransactionDatabaseInterface interface{
	AddTransaction(tm *domain.TransactionModel) error
	GetTransaction(transactionId uuid.UUID) (domain.TransactionModel, error)
}

func (db *SQLManager) AddTransaction(tm *domain.TransactionModel) error {
	stmt := `insert into transaction_model (user_id, transaction_id, category_id, amount, date, description, created_at, updated_at, type, payment_method, status) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.DB.Exec(stmt, tm.UserId, tm.TransactionId, tm.CategoryId, tm.Amount, tm.Date, tm.Description, tm.CreatedAt, tm.UpdatedAt, tm.Type, tm.PaymentMethod, tm.Status)
	if err != nil {
		log.Println("Error saving the transaction to the database:", err)
		return err
	}
	return nil
}

func (db *SQLManager) GetTransaction(transactionId uuid.UUID) (domain.TransactionModel, error) {
	stmt := `select user_id, transaction_id, category_id, amount, date, description, created_at, updated_at, type, payment_method, status from transaction_model where transaction_id = ?`
	var transaction domain.TransactionModel
	err := db.DB.QueryRow(stmt, transactionId).Scan(&transaction.UserId, &transaction.TransactionId, &transaction.CategoryId, &transaction.Amount, &transaction.Date, &transaction.Description, &transaction.CreatedAt, &transaction.UpdatedAt, &transaction.Type, &transaction.PaymentMethod, &transaction.Status)
	if err != nil {
		log.Println("Error retrieving transaction:", err)
		return transaction, err
	}
	return transaction, nil
}
