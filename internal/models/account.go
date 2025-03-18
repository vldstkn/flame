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
	Location  *string `db:"location"`
}

type UserPhoto struct {
	Id         int64   `db:"id"`
	UploadedAt *string `db:"uploaded_at"`
	UserId     *int64  `db:"user_id"`
	PhotoUrl   string  `db:"photo_url"`
	IsMain     *bool   `db:"is_main"`
}

type GetMatchingUser struct {
	User
	PhotoId  *int64  `db:"photo_id"`
	PhotoUrl *string `db:"photo_url"`
	Lon      float64 `db:"lon"`
	Lat      float64 `db:"lat"`
}

type UserPreferences struct {
	UserId   int64   `db:"user_id"`
	Distance *int    `db:"distance"`
	Age      *int    `db:"age"`
	Gender   *string `db:"gender"`
	City     *string `db:"city"`
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
