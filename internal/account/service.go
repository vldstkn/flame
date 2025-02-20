package account

import (
	"flame/internal/interfaces"
	"flame/internal/models"
	http_errors "flame/pkg/errors"
	"flame/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"time"
)

type ServiceDeps struct {
	Repository interfaces.AccountRepository
	Logger     *slog.Logger
}
type Service struct {
	Logger     *slog.Logger
	Repository interfaces.AccountRepository
}

func NewService(deps *ServiceDeps) *Service {
	return &Service{
		Logger:     deps.Logger,
		Repository: deps.Repository,
	}
}

func (service *Service) IssueToken(secret string, data jwt.Data) (*interfaces.AccountSIssueToken, error) {
	op := "Service.IssueToken"
	j := jwt.NewJWT(secret)
	accessToken, err := j.Create(data, time.Now().Add(time.Hour*2).Add(time.Minute*10))
	if err != nil {
		service.Logger.Error(err.Error(),
			slog.String("Error location", op),
		)
		return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	refreshToken, err := j.Create(data, time.Now().AddDate(0, 0, 2).Add(time.Hour*2))
	if err != nil {
		service.Logger.Error(err.Error(),
			slog.String("Error location", op),
		)
		return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	return &interfaces.AccountSIssueToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func (service *Service) Login(email, password string) (int64, error) {
	user := service.Repository.GetByEmail(email)
	if user == nil {
		return -1, status.Errorf(codes.InvalidArgument, http_errors.InvalidNameOrPassword)
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return -1, status.Errorf(codes.InvalidArgument, http_errors.InvalidNameOrPassword)
	}
	return user.Id, nil
}
func (service *Service) Register(data *interfaces.AccountSRegisterDeps) (int64, error) {
	op := "Service.Register"
	existsUser := service.Repository.GetByEmail(data.Email)
	if existsUser != nil {
		return -1, status.Errorf(codes.InvalidArgument, http_errors.UserExists)
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		service.Logger.Error(err.Error(),
			slog.String("Error location", op),
			slog.String("Email", data.Email),
		)
		return -1, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	user := &models.User{
		Email:      data.Email,
		Password:   string(hashPassword),
		Gender:     data.Gender,
		LookingFor: data.LookingFor,
		Name:       data.Name,
	}
	id, err := service.Repository.Create(user)
	if err != nil {
		service.Logger.Error(err.Error(),
			slog.String("Error location", op),
			slog.String("Email", data.Email),
			slog.String("Gender", string(data.Gender)),
			slog.String("Looking for", string(data.Gender)),
			slog.String("Name", data.Name),
		)
		return -1, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	return id, nil
}

func (service *Service) GetTokens(id int64, secret string) (*interfaces.AccountSIssueToken, error) {
	user := service.Repository.GetById(id)
	if user == nil {
		return nil, status.Errorf(codes.InvalidArgument, http.StatusText(http.StatusBadRequest))
	}
	tokens, err := service.IssueToken(secret, jwt.Data{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}
