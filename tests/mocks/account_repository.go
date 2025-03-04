package mocks

import (
	"flame/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (mock *MockAccountRepository) GetById(id int64) *models.User {
	args := mock.Called(id)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*models.User)
}
func (mock *MockAccountRepository) Create(user *models.User) (int64, error) {
	args := mock.Called(user)
	return int64(args.Int(0)), args.Error(1)
}
func (mock *MockAccountRepository) GetByEmail(email string) *models.User {
	args := mock.Called(email)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*models.User)
}
