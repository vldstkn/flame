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
