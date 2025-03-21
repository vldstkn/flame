package interfaces

import (
	"flame/internal/models"
	"flame/pkg/jwt"
	"flame/pkg/pb"
)

type AccountService interface {
	IssueToken(secret string, data jwt.Data) (*AccountSIssueToken, error)
	Login(email, password, location string) (int64, error)
	Register(data *AccountSRegisterDeps) (int64, error)
	GetTokens(secret string, data jwt.Data) (*AccountSIssueToken, error)
	UpdateProfile(data *pb.UpdateProfileReq) error
	GetProfile(id int64) (*pb.GetProfileRes, error)
	UploadPhoto(userId int64, link string) error
	DeletePhoto(userId, photoId int64) (string, error)
	UpdateLocation(userId int64, location string) error
	UpdatePreferences(prefer *pb.UpdatePreferencesReq) error
}
type AccountRepository interface {
	GetById(id int64) *models.User
	Create(user *models.User) (int64, error)
	GetByEmail(email string) *models.User
	UpdateProfile(user *models.User) error
	UploadPhoto(userId int64, link string) (*int64, error)
	SetMainPhoto(userId int64, mainPhotoId int64) error
	GetUserProfilePhotos(userId int64) []models.UserPhoto
	DeletePhoto(photoId int64) error
	GetPhoto(photoId int64) *models.UserPhoto
	GetLastUserPhoto(userId int64) *models.UserPhoto
	GetDistance(user *models.User) (*float64, error)
	GetPreferences(userId int64) *models.UserPreferences
	UpdateLocationRedis(key string, lonLat models.LonLat) error
	UpdatePreferences(prefer *models.UserPreferences) error
}

type AccountSRegisterDeps struct {
	Name     string
	Password string
	Email    string
	Location string
}

type AccountSIssueToken struct {
	AccessToken  string
	RefreshToken string
}
