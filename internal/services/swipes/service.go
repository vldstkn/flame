package swipes

import (
	"flame/internal/interfaces"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
)

type ServiceDeps struct {
	Repository interfaces.SwipesRepository
	Logger     *slog.Logger
}
type Service struct {
	Logger     *slog.Logger
	Repository interfaces.SwipesRepository
}

func NewService(deps *ServiceDeps) *Service {
	return &Service{
		Logger:     deps.Logger,
		Repository: deps.Repository,
	}
}

func (service *Service) CreateOrUpdate(userId1, userId2 int64, isLike bool) error {
	if userId1 == userId2 {
		return status.Errorf(codes.InvalidArgument, http.StatusText(http.StatusBadRequest))
	}
	err := service.Repository.CreateOrUpdate(userId1, userId2, isLike)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, http.StatusText(http.StatusBadRequest))
	}
	candidateListKey := fmt.Sprintf("user:%d:candidates", userId1)
	err = service.Repository.RemoveSwipeFromRedis(candidateListKey, userId2)
	if err != nil {
		service.Logger.Error(err.Error(),
			slog.String("Error location", "service.Repository.RemoveSwipeFromRedis"),
			slog.String("CandidateListKey", candidateListKey),
			slog.Int64("UserId2", userId2),
		)
	}
	return nil
}

func (service *Service) GetUnreadSwipes(userId int64) []int64 {
	userIds := service.Repository.GetUnreadSwipes(userId)
	return userIds
}
