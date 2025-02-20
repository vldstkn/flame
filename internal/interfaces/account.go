package interfaces

import (
	"flame/internal/models"
	"flame/pkg/jwt"
)

type AccountService interface {
	IssueToken(secret string, data jwt.Data) (*AccountSIssueToken, error)
	Login(email, password string) (int64, error)
	Register(data *AccountSRegisterDeps) (int64, error)
	GetTokens(id int64, secret string) (*AccountSIssueToken, error)
}
type AccountRepository interface {
	GetById(id int64) *models.User
	Create(user *models.User) (int64, error)
	GetByEmail(email string) *models.User
}

type AccountSRegisterDeps struct {
	Name       string
	Password   string
	Email      string
	Gender     models.Gender
	LookingFor models.Gender
}

type AccountSIssueToken struct {
	AccessToken  string
	RefreshToken string
}
