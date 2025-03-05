package mappers

import (
	"flame/internal/models"
	"flame/pkg/pb"
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

func FromModelGetMatchingUserToGrpc(user models.GetMatchingUser) *pb.UserProfile {
	return &pb.UserProfile{
		Id:        user.Id,
		Name:      user.Name,
		BirthDate: user.BirthDate,
		City:      user.City,
		Bio:       user.Bio,
		Gender:    user.Gender,
		Photos:    []*pb.UserPhoto{FromModelPhotoToGrpc(user.Photo)},
	}
}
func FromModelGetMatchingUsersToGrpc(users []models.GetMatchingUser) []*pb.UserProfile {
	res := make([]*pb.UserProfile, len(users))
	if res == nil {
		return nil
	}
	for i, u := range users {
		res[i] = FromModelGetMatchingUserToGrpc(u)
	}
	return res
}
