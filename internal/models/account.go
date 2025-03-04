package models

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type User struct {
	Id        int64   `db:"id"`
	CreatedAt string  `db:"created_at"`
	UpdatedAt string  `db:"updated_at"`
	Email     string  `db:"email"`
	Password  string  `db:"password"`
	BirthDate *string `db:"birth_date"`
	City      *string `db:"city"`
	Gender    *string `db:"gender"`
	Name      string  `db:"name"`
	Bio       *string `db:"bio"`
}

type UserPhoto struct {
	Id         int64   `db:"id"`
	UploadedAt *string `db:"uploaded_at"`
	UserId     *int64  `db:"user_id"`
	PhotoUrl   string  `db:"photo_url"`
	IsMain     *bool   `db:"is_main"`
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
