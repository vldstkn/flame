package interfaces

import "flame/internal/models"

type SwipesService interface {
	CreateOrUpdate(UserId1, userId2 int64, isLike bool) error
	GetUnreadSwipes(userId int64) []int64
}

type SwipesRepository interface {
	CreateOrUpdate(UserId1, userId2 int64, isLike bool) error
	GetUnreadSwipes(userId int64) []int64
	GetSwipeById(userId1, userId2 int64) *models.Swipe
	RemoveSwipeFromRedis(candidateListKey string, userId int64) error
}
