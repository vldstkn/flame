package mathcing

import (
	"context"
	"flame/internal/interfaces"
	"flame/internal/mappers"
	"flame/internal/models"
	"flame/pkg/db"
	http_errors "flame/pkg/errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"strconv"
)

type ServiceDeps struct {
	Repository interfaces.MatchingRepository
	Logger     *slog.Logger
	Redis      *db.Redis
}
type Service struct {
	Logger     *slog.Logger
	Repository interfaces.MatchingRepository
	Redis      *db.Redis
}

func NewService(deps *ServiceDeps) *Service {
	return &Service{
		Logger:     deps.Logger,
		Repository: deps.Repository,
		Redis:      deps.Redis,
	}
}

func (service *Service) GetMatchingUsers(userId int64, location string) ([]models.GetMatchingUser, *models.LonLat, error) {
	ctx := context.Background()
	lonLat := service.Repository.GetLonLat(userId)
	if lonLat == nil {
		return nil, nil, status.Errorf(codes.InvalidArgument, http_errors.LocationIsInvalid)
	}
	candidatesKey := fmt.Sprintf("user:%d:candidates", userId)
	length, err := service.Redis.SCard(ctx, candidatesKey).Result()
	if err != nil {
		service.Logger.Error(err.Error(), slog.String("Error location", "service.Redis.SCard"))
	}
	if length == 0 || err != nil {
		users, err := service.Repository.GetMatchingUsers(userId)
		validUsers := service.Repository.DeleteDuplicateMatch(userId, users)
		if err != nil {
			service.Logger.Error(err.Error(),
				slog.String("Error location", "service.Repository.GetMatchingUsers"),
			)
			return nil, nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
		}
		err = service.addUsersToRedis(candidatesKey, validUsers)
		if err != nil {
			service.Logger.Error(err.Error(), slog.String("Error location", "service.AddUsersToRedis"))
		}
		return validUsers, lonLat, nil
	} else {
		users := service.GetUsersFromRedis(ctx, candidatesKey)
		return users, lonLat, nil
	}
}

func (service *Service) GetUsersFromRedis(ctx context.Context, candidatesKey string) []models.GetMatchingUser {

	candidateIds, err := service.Redis.SRandMemberN(ctx, candidatesKey, 100).Result()
	if err != nil {
		service.Logger.Error(err.Error(), slog.String("Error location", "service.Redis.LRange"))
		return nil
	}
	var users []models.GetMatchingUser

	for _, candidateIdStr := range candidateIds {
		candidateId, _ := strconv.Atoi(candidateIdStr)
		candidateKey := fmt.Sprintf("user:%d", candidateId)
		candidateData, err := service.Redis.HGetAll(ctx, candidateKey).Result()
		if err != nil {
			service.Logger.Error(err.Error(), slog.String("Error location", "service.Redis.HGetAll"))
			return nil
		}
		user := mappers.FromMapToModelMatchingUser(candidateData)
		users = append(users, user)
	}
	return users
}

func (service *Service) addUsersToRedis(candidatesKey string, users []models.GetMatchingUser) error {
	for _, user := range users {
		userKey := fmt.Sprintf("user:%d", user.Id)

		err := service.Redis.HSet(context.Background(), userKey, mappers.FromModelToMapMatchingUser(user)).Err()
		if err != nil {
			service.Logger.Error(err.Error(), slog.String("Error location", "service.Redis.HSet"))
		}
		err = service.Redis.SAdd(context.Background(), candidatesKey, user.Id).Err()
		if err != nil {
			service.Logger.Error(err.Error(), slog.String("Error location", "service.Redis.RPush"))
		}
	}
	return nil
}

func (service *Service) UpdateRedis(userId int64) error {
	ctx := context.Background()
	candidatesKey := fmt.Sprintf("user:%d:candidates", userId)
	lonLat := service.Repository.GetLonLat(userId)
	if lonLat == nil {
		service.Logger.Error("lonLat is nil",
			slog.String("Error location", "service.Repository.GetLonLat"),
			slog.Int64("UserId", userId))
		return status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	err := service.Redis.Del(ctx, candidatesKey).Err()
	if err != nil {
		service.Logger.Error(err.Error(),
			slog.String("Error location", "service.Redis.Del"),
			slog.String("Key", candidatesKey),
		)
		return status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	users, err := service.Repository.GetMatchingUsers(userId)
	validUsers := service.Repository.DeleteDuplicateMatch(userId, users)
	if err != nil {
		service.Logger.Error(err.Error(),
			slog.String("Error location", "service.Repository.GetMatchingUsers"),
		)
		return status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	err = service.addUsersToRedis(candidatesKey, validUsers)
	if err != nil {
		service.Logger.Error(err.Error(), slog.String("Error location", "service.AddUsersToRedis"))
	}
	return nil
}
