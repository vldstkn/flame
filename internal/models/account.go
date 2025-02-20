package models

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type User struct {
	Id         int64  `db:"id"`
	Email      string `db:"email"`
	Password   string `db:"password"`
	Gender     Gender `db:"gender"`
	LookingFor Gender `db:"looking_for"`
	Name       string `db:"name"`
	Bio        string `db:"bio"`
}

type UserPhotos struct {
	Id         int64  `db:"id"`
	UploadedAt string `db:"uploaded_at"`
	UserId     int64  `db:"user_id"`
	PhotoUrl   string `db:"photo_url"`
	IsMain     bool   `db:"is_main"`
}

func GenderIsValid(str string) bool {
	switch str {
	case string(Male):
		return true
	case string(Female):
		return true
	default:
		return false
	}
}
