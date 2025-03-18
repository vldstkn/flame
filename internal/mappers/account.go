package mappers

import (
	"flame/internal/models"
	"flame/pkg/pb"
	"github.com/umahmood/haversine"
	"time"
)

func FromModelPhotoToGrpc(photo models.UserPhoto) *pb.UserPhoto {
	return &pb.UserPhoto{
		Id:         photo.Id,
		UploadedAt: photo.UploadedAt,
		UserId:     photo.UserId,
		PhotoUrl:   photo.PhotoUrl,
		IsMain:     photo.IsMain,
	}
}
func FromModelPhotosToGrpc(photos []models.UserPhoto) []*pb.UserPhoto {
	res := make([]*pb.UserPhoto, len(photos))
	if res == nil {
		return nil
	}
	for i, p := range photos {
		res[i] = FromModelPhotoToGrpc(p)
	}
	return res
}

func FromModelGetMatchingUserToGrpc(user models.GetMatchingUser, lonLat *models.LonLat) *pb.UserMatch {
	var age *int32
	if user.BirthDate != nil {
		birthTime, err := time.Parse(time.RFC3339, *user.BirthDate)
		if err == nil {
			now := time.Now()
			a := int32(now.Year() - birthTime.Year())
			if now.YearDay() < birthTime.YearDay() {
				a--
			}
			age = &a
		}

	}
	var distance float64
	p1 := haversine.Coord{Lat: lonLat.Lat, Lon: lonLat.Lon}
	p2 := haversine.Coord{Lat: user.Lat, Lon: user.Lon}
	_, distance = haversine.Distance(p1, p2)

	if user.PhotoUrl == nil || *user.PhotoUrl == "" {
		return &pb.UserMatch{
			Id:       user.Id,
			Name:     user.Name,
			Age:      age,
			City:     user.City,
			Gender:   user.Gender,
			Photo:    nil,
			Distance: int32(distance * 1000),
		}
	}
	return &pb.UserMatch{
		Id:     user.Id,
		Name:   user.Name,
		Age:    age,
		City:   user.City,
		Gender: user.Gender,
		Photo: &pb.UserPhoto{
			Id:       *user.PhotoId,
			PhotoUrl: *user.PhotoUrl,
		},
		Distance: int32(distance * 1000),
	}
}
func FromModelGetMatchingUsersToGrpc(users []models.GetMatchingUser, lonLat *models.LonLat) []*pb.UserMatch {
	res := make([]*pb.UserMatch, len(users))
	if res == nil {
		return nil
	}
	for i, u := range users {
		res[i] = FromModelGetMatchingUserToGrpc(u, lonLat)
	}
	return res
}
