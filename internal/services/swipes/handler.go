package swipes

import (
	"context"
	"flame/internal/config"
	"flame/internal/interfaces"
	"flame/pkg/pb"
	"log/slog"
)

type Handler struct {
	Logger  *slog.Logger
	Config  *config.Config
	Service interfaces.SwipesService
	pb.UnsafeSwipesServer
}

type HandlerDeps struct {
	Logger  *slog.Logger
	Config  *config.Config
	Service interfaces.SwipesService
}

func NewHandler(deps *HandlerDeps) *Handler {
	return &Handler{
		Logger:  deps.Logger,
		Config:  deps.Config,
		Service: deps.Service,
	}
}

func (handler *Handler) CreateOrUpdateSwipe(ctx context.Context, r *pb.CreateOrUpdateSwipeReq) (*pb.CreateOrUpdateSwipeRes, error) {
	err := handler.Service.CreateOrUpdate(r.UserId1, r.UserId2, r.IsLike)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdateSwipeRes{}, nil
}

func (handler *Handler) GetUnreadSwipes(ctx context.Context, r *pb.GetUnreadSwipesReq) (*pb.GetUnreadSwipesRes, error) {
	ids := handler.Service.GetUnreadSwipes(r.UserId)
	return &pb.GetUnreadSwipesRes{
		UserIds: ids,
	}, nil
}
