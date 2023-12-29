package domain

import (
	"strings"

	"github.com/google/uuid"
	gen "github.com/pallinder/go-randomdata"
)

type TransactionDTOStruct struct{}

func TransactionDTOBuilder() *TransactionDTOStruct {
	return &TransactionDTOStruct{}
}

func (b *TransactionDTOStruct) Build() TransactionDTO {
	return TransactionDTO{
		UserId: uuid.New(),
		TransactionId: uuid.New(),
		CategoryId:    int64(gen.Number(15)),
		Amount:        float64(gen.Number(3)),
		Date:          int64(gen.Number(13)),
		Description:   strings.Split(gen.Paragraph(), ".")[0],
		CreatedAt:     int64(gen.Number(13)),
		UpdatedAt:     int64(gen.Number(13)),
		Type:          types[gen.Number(0, len(types))],
		PaymentMethod: payments[gen.Number(0, len(payments))],
		Status:        statuses[gen.Number(0, len(statuses))],
	}
}
