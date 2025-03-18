package interfaces

import "flame/internal/models"

type MatchingService interface {
	GetMatchingUsers(userId int64, location string) ([]models.GetMatchingUser, *models.LonLat, error)
}

type MatchingRepository interface {
	GetMatchingUsers(userId int64) ([]models.GetMatchingUser, error)
	GetLonLat(userId int64) *models.LonLat
	DeleteDuplicateMatch(userId int64, users []models.GetMatchingUser) []models.GetMatchingUser
}
