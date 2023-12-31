package service

import (
	"github.com/google/uuid"
	"github.com/hld3/personal-finance-go/database"
	"github.com/hld3/personal-finance-go/domain"
)

type TransactionServiceInterface interface {
	GetTransaction(transactionId uuid.UUID) (*domain.TransactionModel, error)
}

type TransactionService struct {
	UDBI database.TransactionDatabaseInterface
}

func (t *TransactionService) AddTransaction(transactionData *domain.TransactionData) error {
	err := transactionData.ValidateTransaction()
	if err != nil {
		return err
	}

	tm := convertTransactionDTOToModel(&transactionData.Transaction)
	err = t.UDBI.AddTransaction(&tm)
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionService) GetTransaction(transactionId uuid.UUID) (*domain.TransactionModel, error) {
	transaction, err := t.UDBI.GetTransaction(transactionId)
	if err != nil {
		return &domain.TransactionModel{}, err
	}
	return &transaction, nil
}

func convertTransactionDTOToModel(from *domain.TransactionDTO) domain.TransactionModel {
	return domain.TransactionModel{
		UserId:        from.UserId,
		TransactionId: from.TransactionId,
		CategoryId:    from.CategoryId,
		Amount:        from.Amount,
		Date:          from.Date,
		Description:   from.Description,
		CreatedAt:     from.CreatedAt,
		UpdatedAt:     from.UpdatedAt,
		Type:          from.Type,
		PaymentMethod: from.PaymentMethod,
		Status:        from.Status,
	}
}
