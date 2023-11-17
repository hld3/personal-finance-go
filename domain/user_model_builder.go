package domain

import (
	"time"

	"github.com/google/uuid"
	gen "github.com/pallinder/go-randomdata"
)

type UserModelBuild struct {
	userId uuid.UUID
}

func UserModelBuilder() *UserModelBuild {
	return &UserModelBuild{
		userId: uuid.New(),
	}
}

func (b *UserModelBuild) Build() UserModel {
	return UserModel{
		UserId: b.userId,
		FirstName: gen.FirstName(1),
		LastName: gen.LastName(),
		Email: gen.Email(),
		Phone: gen.PhoneNumber(),
		DateOfBirth: time.Now().UnixMilli(),
		CreationDate: time.Now().UnixMilli(),
		PasswordHash: gen.Letters(15),
	}
}

func (b *UserModelBuild) WithUserId(userId uuid.UUID) *UserModelBuild {
	b.userId = userId
	return b
}
