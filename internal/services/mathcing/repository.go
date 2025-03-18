package mathcing

import (
	"flame/internal/models"
	"flame/pkg/db"
)

type RepositoryDeps struct {
	AccountDB *db.DB
	SwipesDB  *db.DB
}
type Repository struct {
	AccountDB *db.DB
	SwipesDB  *db.DB
}

func NewRepository(deps *RepositoryDeps) *Repository {
	return &Repository{
		AccountDB: deps.AccountDB,
		SwipesDB:  deps.SwipesDB,
	}
}

func (repo *Repository) GetMatchingUsers(userId int64) ([]models.GetMatchingUser, error) {
	var users []models.GetMatchingUser
	err := repo.AccountDB.Select(&users,
		`SELECT u1.id, u1.name, u1.birth_date, u1.city, u1.gender, st_x(ST_AsText(u1.location)::geometry) as lon, st_y(ST_AsText(u1.location)::geometry) as lat, up.photo_url, up.id as photo_id
       			FROM users u
       			JOIN preferences p ON	u.id = p.user_id
       			JOIN users u1 ON u1.location IS NOT NULL AND st_dwithin(u1.location, u.location, p.distance * 1000) AND
       			(p.age IS NULL OR (EXTRACT(YEAR FROM AGE(u1.birth_date)) BETWEEN  GREATEST(p.age + ROUND(p.age * 0.8), 16) AND p.age + ROUND(p.age * 1.2)))AND
						(u1.city IS NULL OR p.city IS NULL OR u1.city = p.city) AND (u1.gender IS NULL OR p.gender IS NULL OR u1.gender = p.gender) AND u.id != u1.id
       			LEFT JOIN user_photos up ON u1.id = up.user_id
       			WHERE u.id=$1`, userId)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *Repository) DeleteDuplicateMatch(userId int64, users []models.GetMatchingUser) []models.GetMatchingUser {
	var removeIds []int64
	repo.SwipesDB.Select(&removeIds, `SELECT 
    CASE 
    	WHEN user_id1=$1 THEN user_id2
    	WHEN user_id2=$1 THEN user_id1
    END FROM swipes WHERE (user_id1=$1 AND user_is_liked1 IS NOT NULL) OR (user_id2=$1 AND user_is_liked2 IS NOT NULL)`, userId)
	removeMap := make(map[int64]struct{}, len(removeIds))
	for _, id := range removeIds {
		removeMap[id] = struct{}{}
	}
	var res []models.GetMatchingUser

	for _, el := range users {
		if _, found := removeMap[el.Id]; !found {
			res = append(res, el)
		}
	}

	return res
}

func (repo *Repository) GetLonLat(userId int64) *models.LonLat {
	var lonLat models.LonLat

	err := repo.AccountDB.Get(&lonLat, `SELECT st_x(ST_AsText(location)::geometry) as lon,st_y(ST_AsText(location)::geometry) as lat FROM users WHERE id=$1`, userId)
	if err != nil {
		return nil
	}
	return &lonLat
}
