package account

import (
	"context"
	"flame/internal/config"
	"flame/internal/interfaces"
	"flame/internal/mappers"
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
		Name:     r.Name,
		Password: r.Password,
		Email:    r.Email,
		Location: r.Location,
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
	id, err := handler.Service.Login(r.Email, r.Password, r.Location)
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
	tokens, err := handler.Service.GetTokens(handler.Config.Auth.Jwt, jwt.Data{
		Id: r.Id,
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetTokensRes{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (handler *Handler) UpdateProfile(ctx context.Context, r *pb.UpdateProfileReq) (*pb.UpdateProfileRes, error) {
	err := handler.Service.UpdateProfile(r)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateProfileRes{}, nil
}

func (handler *Handler) GetProfile(ctx context.Context, r *pb.GetProfileReq) (*pb.GetProfileRes, error) {
	response, err := handler.Service.GetProfile(r.Id)
	return response, err
}

func (handler *Handler) UploadPhoto(ctx context.Context, r *pb.UploadPhotoReq) (*pb.UploadPhotoRes, error) {
	err := handler.Service.UploadPhoto(r.UserId, r.LinkPhoto)
	if err != nil {
		return nil, err
	}
	return &pb.UploadPhotoRes{}, nil
}

func (handler *Handler) DeletePhoto(ctx context.Context, r *pb.DeletePhotoReq) (*pb.DeletePhotoRes, error) {
	url, err := handler.Service.DeletePhoto(r.UserId, r.PhotoId)
	if err != nil {
		return nil, err
	}
	return &pb.DeletePhotoRes{
		PhotoUrl: url,
	}, nil
}

func (handler *Handler) GetMatchingUsers(ctx context.Context, r *pb.GetMatchingUsersReq) (*pb.GetMatchingUsersRes, error) {
	users, err := handler.Service.GetMatchingUsers(r.Id, r.Location)
	if err != nil {
		return nil, err
	}
	return &pb.GetMatchingUsersRes{
		Users: mappers.FromModelGetMatchingUsersToGrpc(users),
	}, nil
}
