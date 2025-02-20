package account

import (
	"context"
	"flame/internal/config"
	"flame/internal/interfaces"
	"flame/internal/models"
	"flame/pkg/jwt"
	"flame/pkg/pb"
	"log/slog"
)

type Handler struct {
	Logger  *slog.Logger
	Config  *config.Config
	Service interfaces.AccountService
	pb.UnsafeAccountServer
}

type HandlerDeps struct {
	Logger  *slog.Logger
	Config  *config.Config
	Service interfaces.AccountService
}

func NewHandler(deps *HandlerDeps) *Handler {
	return &Handler{
		Logger:  deps.Logger,
		Config:  deps.Config,
		Service: deps.Service,
	}
}

func (handler *Handler) Register(ctx context.Context, r *pb.RegisterReq) (*pb.RegisterRes, error) {
	id, err := handler.Service.Register(&interfaces.AccountSRegisterDeps{
		Name:       r.Name,
		Password:   r.Password,
		Email:      r.Email,
		Gender:     models.Gender(r.Gender),
		LookingFor: models.Gender(r.LookingFor),
	})
	if err != nil {
		return nil, err
	}
	tokens, err := handler.Service.IssueToken(handler.Config.Auth.Jwt, jwt.Data{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &pb.RegisterRes{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
func (handler *Handler) Login(ctx context.Context, r *pb.LoginReq) (*pb.LoginRes, error) {
	id, err := handler.Service.Login(r.Email, r.Password)
	if err != nil {
		return nil, err
	}
	tokens, err := handler.Service.IssueToken(handler.Config.Auth.Jwt, jwt.Data{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &pb.LoginRes{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
func (handler *Handler) GetTokens(ctx context.Context, r *pb.GetTokensReq) (*pb.GetTokensRes, error) {
	tokens, err := handler.Service.GetTokens(r.Id, handler.Config.Auth.Jwt)
	if err != nil {
		return nil, err
	}

	return &pb.GetTokensRes{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
