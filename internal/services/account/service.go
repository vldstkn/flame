package account

import (
	"flame/internal/interfaces"
	"flame/internal/mappers"
	"flame/internal/models"
	http_errors "flame/pkg/errors"
	"flame/pkg/jwt"
	"flame/pkg/pb"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"strconv"
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
func (service *Service) Login(email, password, location string) (int64, error) {
	user := service.Repository.GetByEmail(email)
	if user == nil {
		return -1, status.Errorf(codes.InvalidArgument, http_errors.InvalidNameOrPassword)
	}
	if location != "" {
		u := &models.User{
			Location: getLocation(location),
		}
		service.Repository.UpdateProfile(u)
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
	loc := getLocation(data.Location)
	user := &models.User{
		Email:    data.Email,
		Password: string(hashPassword),
		Name:     data.Name,
		Location: loc,
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
			data.Name = nil
		} else {
			user.Name = name
		}
	}
	if data.City != nil {
		city := strings.TrimSpace(*data.City)
		if len(city) == 0 {
			data.City = nil
		} else {
			user.City = &city
		}
	}
	if data.Bio != nil {
		bio := strings.TrimSpace(*data.Bio)
		if len(bio) == 0 {
			data.Bio = nil
		} else {
			user.Bio = &bio
		}
	}
	if data.Gender != nil {
		gender := strings.TrimSpace(*data.Gender)
		if !models.GenderIsValid(*data.Gender) {
			data.Bio = nil
		} else {
			user.Gender = &gender
		}
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
			Location:  user.Location,
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

func (service *Service) UpdateLocation(userId int64, location string) error {
	if location == "" {
		return status.Errorf(codes.InvalidArgument, http_errors.LocationIsInvalid)
	}
	user := &models.User{
		Id:       userId,
		Location: getLocation(location),
	}
	distance, err := service.Repository.GetDistance(user)
	if err != nil {
		service.Logger.Error(err.Error(),
			slog.String("Error location", "service.Repository.GetDistance"),
			slog.Int64("User id", userId),
			slog.String("Location", location),
		)
		return status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	pref := service.Repository.GetPreferences(userId)
	if pref == nil {
		service.Logger.Error("preference is nil",
			slog.String("Error location", "service.Repository.GetPreferences"),
			slog.Int64("User id", userId),
			slog.String("Location", location),
		)
		return status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	if distance == nil || int32(*distance)/1000 >= *pref.Distance {
		err = service.Repository.UpdateProfile(user)
		if err != nil {
			service.Logger.Error(err.Error(),
				slog.String("Error location", "service.UpdateLocation"),
				slog.String("Location", location),
			)
			return status.Errorf(codes.Internal, http_errors.LocationIsInvalid)
		}

		//TODO: поиск пользователей
	}
	key := fmt.Sprintf("user:%d", userId)
	lonLat := getLonLat(location)
	err = service.Repository.UpdateLocationRedis(key, lonLat)
	if err != nil {
		service.Logger.Error(err.Error(),
			slog.String("Error location", "service.Repository.UpdateLocationRedis"),
			slog.String("Location", location),
			slog.Float64("Lon", lonLat.Lon),
			slog.Float64("Lat", lonLat.Lat),
		)
	}
	return nil
}
func getLocation(loc string) *string {
	if loc == "" {
		return nil
	} else {
		l := fmt.Sprintf("SRID=4326;POINT%s", loc)
		return &l
	}
}

func getLonLat(loc string) models.LonLat {
	loc = loc[1 : len(loc)-1]
	locArr := strings.Split(loc, " ")
	lon, _ := strconv.ParseFloat(locArr[0], 64)
	lat, _ := strconv.ParseFloat(locArr[1], 64)
	return models.LonLat{
		Lon: lon,
		Lat: lat,
	}
}

func (service *Service) UpdatePreferences(r *pb.UpdatePreferencesReq) error {
	user := service.Repository.GetById(r.UserId)
	if user == nil {
		return status.Errorf(codes.InvalidArgument, http.StatusText(http.StatusBadRequest))
	}
	var pref models.UserPreferences
	pref.UserId = r.UserId
	if r.Age != nil && (*r.Age > 110 || *r.Age < 16) {
		return status.Errorf(codes.InvalidArgument, http_errors.InvalidAge)
	}
	pref.Age = r.Age

	if r.Distance != nil && (*r.Distance > 50 || *r.Distance < 3) {
		return status.Errorf(codes.InvalidArgument, http_errors.InvalidDistance)
	}
	pref.Distance = r.Distance

	if r.Gender != nil && !models.GenderIsValid(*r.Gender) {
		return status.Errorf(codes.InvalidArgument, http_errors.InvalidGender)
	}
	pref.Gender = r.Gender

	if r.City != nil && *r.City != "" && len(*r.City) < 2 {
		return status.Errorf(codes.InvalidArgument, http_errors.InvalidCity)
	}
	pref.City = r.City

	if r.Age != nil {
		if *r.Age > 110 || *r.Age < 16 {
			return status.Errorf(codes.InvalidArgument, http_errors.InvalidAge)
		}
		pref.Age = r.Age
	}

	err := service.Repository.UpdatePreferences(&pref)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, http.StatusText(http.StatusBadRequest))
	}
	return nil
}
