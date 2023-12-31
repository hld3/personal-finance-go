package domain

import (
	"strings"

	"github.com/google/uuid"
	gen "github.com/pallinder/go-randomdata"
)

type TransactionDTOStruct struct {
	description string
}

func TransactionDTOBuilder() *TransactionDTOStruct {
	return &TransactionDTOStruct{
		description: strings.Split(gen.Paragraph(), ".")[0],
	}
}

func (b *TransactionDTOStruct) Build() TransactionDTO {
	return TransactionDTO{
		UserId:        uuid.New(),
		TransactionId: uuid.New(),
		CategoryId:    int64(gen.Number(1, 15)),
		Amount:        float64(gen.Number(1, 5)),
		Date:          int64(gen.Number(1, 13)),
		Description:   b.description,
		CreatedAt:     int64(gen.Number(1, 13)),
		UpdatedAt:     int64(gen.Number(1, 13)),
		Type:          types[gen.Number(0, len(types))],
		PaymentMethod: payments[gen.Number(0, len(payments))],
		Status:        statuses[gen.Number(0, len(statuses))],
	}
}

func (b *TransactionDTOStruct) WithDescription(description string) *TransactionDTOStruct {
	b.description = description
	return b
}
