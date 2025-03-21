package account

import (
	"context"
	"database/sql"
	"flame/internal/models"
	"flame/pkg/db"
	"fmt"
	"reflect"
)

type Repository struct {
	DB    *db.DB
	Redis *db.Redis
}
type RepositoryDeps struct {
	DB    *db.DB
	Redis *db.Redis
}

func NewRepository(deps *RepositoryDeps) *Repository {
	return &Repository{
		DB:    deps.DB,
		Redis: deps.Redis,
	}
}

func (repo *Repository) GetById(id int64) *models.User {
	var user models.User
	err := repo.DB.Get(&user, `SELECT *,
      '(' || to_char(ST_X(location::geometry), 'FM999990.000000') || ' ' ||
            to_char(ST_Y(location::geometry), 'FM999990.000000') || ')'
    AS location FROM users WHERE id=$1`, id)

	if err != nil {
		return nil
	}
	return &user
}

func (repo *Repository) GetByEmail(email string) *models.User {
	var user models.User
	err := repo.DB.Get(&user, `SELECT * FROM users WHERE email=$1`, email)
	if err != nil {
		return nil
	}
	return &user
}
func (repo *Repository) Create(user *models.User) (int64, error) {
	var id int64
	tr, err := repo.DB.Beginx()
	if err != nil {
		return -1, err
	}
	err = repo.DB.QueryRow(`INSERT INTO users (email, password, name, location) 
																   VALUES ($1,$2,$3, ST_GeographyFromText($4)) RETURNING id`,
		user.Email, user.Password, user.Name, user.Location).Scan(&id)
	if err != nil {
		tr.Rollback()
		return -1, err
	}
	_, err = repo.DB.Exec(`INSERT INTO preferences (user_id) VALUES ($1)`, id)
	if err != nil {
		tr.Rollback()
		return -1, err
	}
	tr.Commit()
	return id, nil
}

func (repo *Repository) UpdateProfile(user *models.User) error {
	flag := false
	query := `UPDATE users SET`
	count := 1

	var args []interface{}

	val := reflect.ValueOf(*user)
	typ := reflect.TypeOf(*user)
	for i := 0; i < val.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)
		if fieldType.Name == "Id" {
			continue
		}
		if fieldValue.Kind() == reflect.Ptr {
			if !fieldValue.IsNil() {
				if flag {
					query += ","
				} else {
					flag = true
				}
				query += fmt.Sprintf(" %s=$%d", fieldType.Tag.Get("db"), count)
				args = append(args, fieldValue.Interface().(*string))

				count++
			}
		} else if fieldValue.Kind() == reflect.String {
			if len(fieldValue.Interface().(string)) != 0 {
				if flag {
					query += ","
				} else {
					flag = true
				}
				query += fmt.Sprintf(" %s=$%d", fieldType.Tag.Get("db"), count)
				args = append(args, fieldValue.Interface().(string))
				count++
			}
		}
	}
	if !flag {
		return nil
	}
	query += fmt.Sprintf(" WHERE id=$%d", count)
	args = append(args, user.Id)
	_, err := repo.DB.Exec(query, args...)
	return err
}

func (repo *Repository) UploadPhoto(userId int64, link string) (*int64, error) {
	var id int64
	err := repo.DB.QueryRow(`INSERT INTO user_photos (user_id, photo_url) VALUES ($1, $2) RETURNING id`, userId, link).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (repo *Repository) SetMainPhoto(userId int64, mainPhotoId int64) error {

	tr, err := repo.DB.Beginx()
	if err != nil {
		return err
	}
	_, err = repo.DB.Exec(`UPDATE user_photos SET is_main=false WHERE user_id=$1 AND is_main=true`, userId)
	if err != nil {
		tr.Rollback()
		return err
	}

	_, err = repo.DB.Exec(`UPDATE user_photos SET is_main=true WHERE id=$1`, mainPhotoId)
	if err != nil {
		tr.Rollback()
		return err
	}
	tr.Commit()
	return nil
}

func (repo *Repository) GetUserProfilePhotos(userId int64) []models.UserPhoto {
	var photos []models.UserPhoto
	err := repo.DB.Select(&photos, `SELECT id, uploaded_at, photo_url, is_main FROM user_photos WHERE user_id=$1`, userId)
	if err != nil {
		return nil
	}
	return photos
}
func (repo *Repository) GetPhoto(photoId int64) *models.UserPhoto {
	var photo models.UserPhoto
	err := repo.DB.Get(&photo, `SELECT * FROM user_photos WHERE id=$1`, photoId)
	if err != nil {
		return nil
	}
	return &photo
}

func (repo *Repository) DeletePhoto(photoId int64) error {
	_, err := repo.DB.Exec(`DELETE FROM user_photos WHERE id=$1`, photoId)
	return err
}

func (repo *Repository) GetLastUserPhoto(userId int64) *models.UserPhoto {
	var photo models.UserPhoto
	err := repo.DB.Get(&photo, `SELECT * FROM user_photos 
       													 WHERE user_id=$1
       													 ORDER BY id DESC 
       													 LIMIT 1`, userId)
	if err != nil {
		return nil
	}

	return &photo
}

func (repo *Repository) GetDistance(user *models.User) (*float64, error) {
	var distance sql.NullFloat64
	err := repo.DB.QueryRow(`SELECT CASE WHEN location IS NOT NULL THEN st_distance($1, location) ELSE NULL END FROM users WHERE id=$2`,
		user.Location, user.Id).Scan(&distance)
	if err != nil {
		return nil, err
	}
	v, _ := distance.Value()
	if v == nil {
		return nil, nil
	}
	return &distance.Float64, nil
}
func (repo *Repository) GetPreferences(userId int64) *models.UserPreferences {
	var pref models.UserPreferences
	err := repo.DB.Get(&pref, `SELECT * FROM preferences WHERE user_id=$1`, userId)
	if err != nil {
		return nil
	}
	return &pref
}
func (repo *Repository) UpdateLocationRedis(key string, lonLat models.LonLat) error {
	err := repo.Redis.HSet(context.Background(), key, map[string]interface{}{
		"lon": lonLat.Lon,
		"lat": lonLat.Lat,
	}).Err()
	return err
}

func (repo *Repository) UpdatePreferences(pref *models.UserPreferences) error {
	flag := false
	query := "UPDATE preferences SET"
	count := 1

	var args []interface{}

	val := reflect.ValueOf(*pref)
	typ := reflect.TypeOf(*pref)
	for i := 0; i < val.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)
		if fieldType.Name == "UserId" {
			continue
		}
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			if flag {
				query += ","
			} else {
				flag = true
			}
			query += fmt.Sprintf(" %s=$%d", fieldType.Tag.Get("db"), count)
			args = append(args, fieldValue.Elem().Interface())
			count++

		}
	}
	if !flag {
		return nil
	}
	query += fmt.Sprintf(" WHERE user_id=$%d", count)
	args = append(args, pref.UserId)
	_, err := repo.DB.Exec(query, args...)
	return err
}
