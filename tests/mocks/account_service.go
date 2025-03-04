package mocks

import (
	"flame/internal/interfaces"
	"flame/pkg/jwt"
	"github.com/stretchr/testify/mock"
)

type MockAccountService struct {
	mock.Mock
}

func (mock *MockAccountService) IssueToken(secret string, data jwt.Data) (*interfaces.AccountSIssueToken, error) {
	args := mock.Called(secret, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.AccountSIssueToken), args.Error(1)
}
func (mock *MockAccountService) Login(email, password string) (int64, error) {
	args := mock.Called(email, password)
	return int64(args.Int(0)), args.Error(1)
}
func (mock *MockAccountService) Register(data *interfaces.AccountSRegisterDeps) (int64, error) {
	args := mock.Called(data)
	return int64(args.Int(0)), args.Error(1)
}
func (mock *MockAccountService) GetTokens(secret string, data jwt.Data) (*interfaces.AccountSIssueToken, error) {
	args := mock.Called(secret, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.AccountSIssueToken), args.Error(1)
}
