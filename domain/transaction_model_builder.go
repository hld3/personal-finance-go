package domain

import (
	"strings"

	"github.com/google/uuid"
	gen "github.com/pallinder/go-randomdata"
)

var types = []TransactionType{INCOME, EXPENSE}
var payments = []TransactionMethod{CASH, CREDIT_CARD, BANK_TRANSFER}
var statuses = []TransactionStatus{PENDING, CLEARED, CANCELLED}

type TransactionModelBuild struct{}

func TransactionModelBuilder() *TransactionModelBuild {
	return &TransactionModelBuild{}
}

func (b *TransactionModelBuild) Build() TransactionModel {
	return TransactionModel{
		UserId:        uuid.New(),
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

type CategoryModelBuild struct{}

func CategoryModelBuilder() *CategoryModelBuild {
	return &CategoryModelBuild{}
}

func (b *CategoryModelBuild) Build() CategoryModel {
	return CategoryModel{
		UserId:      uuid.New(),
		Name:        gen.Letters(10),
		Description: strings.Split(gen.Paragraph(), ".")[0],
	}
}
