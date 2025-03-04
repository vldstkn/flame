package account

import (
	"flame/internal/interfaces"
	"flame/internal/mappers"
	"flame/internal/models"
	http_errors "flame/pkg/errors"
	"flame/pkg/jwt"
	"flame/pkg/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"strings"
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
	refreshToken, _ := j.Create(data, time.Now().AddDate(0, 0, 2).Add(time.Hour*2))
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
		Email:    data.Email,
		Password: string(hashPassword),
		Name:     data.Name,
	}
	id, err := service.Repository.Create(user)
	if err != nil {
		service.Logger.Error(err.Error(),
			slog.String("Error location", op),
			slog.String("Email", data.Email),
			slog.String("Name", data.Name),
		)
		return -1, status.Errorf(codes.InvalidArgument, http.StatusText(http.StatusBadRequest))
	}
	return id, nil
}

func (service *Service) GetTokens(secret string, data jwt.Data) (*interfaces.AccountSIssueToken, error) {
	user := service.Repository.GetById(data.Id)
	if user == nil {
		return nil, status.Errorf(codes.InvalidArgument, http.StatusText(http.StatusBadRequest))
	}
	tokens, err := service.IssueToken(secret, jwt.Data{
		Id: data.Id,
	})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (service *Service) UpdateProfile(data *pb.UpdateProfileReq) error {
	existUser := service.Repository.GetById(data.Id)
	if existUser == nil {
		return status.Errorf(codes.InvalidArgument, http.StatusText(http.StatusBadRequest))
	}
	layout := "2006-01-02"
	user := &models.User{}
	user.Id = data.Id
	if data.BirthDate != nil {
		date, err := time.Parse(layout, *data.BirthDate)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "birth_date is bad")
		}
		if time.Now().Year()-date.Year() < 16 {
			return status.Errorf(codes.InvalidArgument, "user must be over 16 years old")
		}
		user.BirthDate = data.BirthDate
	}
	if data.Name != nil {
		name := strings.TrimSpace(*data.Name)
		if len(name) == 0 {
			return status.Errorf(codes.InvalidArgument, "name can not be empty")
		}
		user.Name = name
	}
	if data.City != nil {
		city := strings.TrimSpace(*data.City)
		if len(city) == 0 {
			return status.Errorf(codes.InvalidArgument, "city can not be empty")
		}
		user.City = &city
	}
	if data.Bio != nil {
		bio := strings.TrimSpace(*data.Bio)
		if len(bio) == 0 {
			return status.Errorf(codes.InvalidArgument, "bio can not be empty")
		}
		user.Bio = &bio
	}
	if data.Gender != nil {
		gender := strings.TrimSpace(*data.Gender)
		if !models.GenderIsValid(*data.Gender) {
			return status.Errorf(codes.InvalidArgument, "gender is bad")
		}
		user.Gender = &gender
	}
	err := service.Repository.UpdateProfile(user)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (service *Service) GetProfile(id int64) (*pb.GetProfileRes, error) {
	user := service.Repository.GetById(id)
	if user == nil {
		return nil, status.Errorf(codes.InvalidArgument, http_errors.UserDoesNotExist)
	}
	userPhotos := service.Repository.GetUserProfilePhotos(id)
	return &pb.GetProfileRes{
		Profile: &pb.UserProfile{
			Id:        user.Id,
			Name:      user.Name,
			BirthDate: user.BirthDate,
			City:      user.City,
			Bio:       user.Bio,
			Gender:    user.Gender,
			Photos:    mappers.FromModelPhotosToGrpc(userPhotos),
		},
	}, nil
}

func (service *Service) UploadPhoto(userId int64, link string) error {
	user := service.Repository.GetById(userId)
	if user == nil {
		return status.Errorf(codes.InvalidArgument, http_errors.UserDoesNotExist)
	}
	photoId, err := service.Repository.UploadPhoto(userId, link)
	if err != nil {
		return status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	err = service.Repository.SetMainPhoto(userId, *photoId)
	if err != nil {
		return status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}
func (service *Service) DeletePhoto(userId, photoId int64) (string, error) {
	photo := service.Repository.GetPhoto(photoId)
	if photo == nil {
		return "", status.Errorf(codes.InvalidArgument, http.StatusText(http.StatusBadRequest))
	}
	if *photo.UserId != userId {
		return "", status.Errorf(codes.PermissionDenied, http.StatusText(http.StatusForbidden))
	}
	err := service.Repository.DeletePhoto(photoId)
	if err != nil {
		return "", status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	if *photo.IsMain {
		lastPhoto := service.Repository.GetLastUserPhoto(userId)
		if lastPhoto == nil {
			return photo.PhotoUrl, nil
		}
		service.Repository.SetMainPhoto(userId, lastPhoto.Id)
	}
	return photo.PhotoUrl, nil
}
