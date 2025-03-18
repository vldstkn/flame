package mathcing

import (
	"context"
	"flame/internal/config"
	"flame/internal/interfaces"
	"flame/internal/mappers"
	"flame/pkg/pb"
	"log/slog"
)

type Handler struct {
	Logger  *slog.Logger
	Config  *config.Config
	Service interfaces.MatchingService
	pb.UnsafeMatchingServer
}

type HandlerDeps struct {
	Logger  *slog.Logger
	Config  *config.Config
	Service interfaces.MatchingService
}

func NewHandler(deps *HandlerDeps) *Handler {
	return &Handler{
		Logger:  deps.Logger,
		Config:  deps.Config,
		Service: deps.Service,
	}
}

func (handler *Handler) GetMatchingUsers(ctx context.Context, r *pb.GetMatchingUsersReq) (*pb.GetMatchingUsersRes, error) {
	users, lonLat, err := handler.Service.GetMatchingUsers(r.Id, r.Location)
	if err != nil {
		return nil, err
	}
	return &pb.GetMatchingUsersRes{
		Users: mappers.FromModelGetMatchingUsersToGrpc(users, lonLat),
	}, nil
}
