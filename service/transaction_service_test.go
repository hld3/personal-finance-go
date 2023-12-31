package service

import (
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hld3/personal-finance-go/domain"
)

func (m *StubDatabase) AddTransaction(tm *domain.TransactionModel) error {
	return nil
}

func (m *StubDatabase) GetTransaction(transactionId uuid.UUID) (domain.TransactionModel, error) {
	transaction := domain.TransactionModelBuilder().Build()
	return transaction, nil
}

// useful, mostly for validation.
func TestAddTransaction(t *testing.T) {
	stubDB := new(StubDatabase)
	transactionService := TransactionService{UDBI: stubDB}

	tests := []struct {
		name        string
		transaction domain.TransactionDTO
		wantErr     bool
		expectedErr string
	}{
		{
			name: "Valid transaction",
			transaction: domain.TransactionDTOBuilder().Build(),
			wantErr: false,
			expectedErr: "",
		},
		{
			name: "Validation error",
			transaction: domain.TransactionDTOBuilder().WithDescription("").Build(),
			wantErr: true,
			expectedErr: "Key: 'TransactionDTO.Description' Error:Field validation for 'Description' failed on the 'required' tag",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transactionData := domain.TransactionData{Transaction: test.transaction, Validator: validator.New()}
			err := transactionService.AddTransaction(&transactionData)

			if test.wantErr && err == nil {
				t.Fatal("Expected error, but it was nil")
			}

			if test.wantErr && !strings.Contains(err.Error(), test.expectedErr) {
				t.Fatalf("Incorrect error returned. want %s, got %v", test.expectedErr, err)
			}

			if !test.wantErr && err != nil {
				t.Fatal("Unexpected error:", err)
			}
		})
	}
}
