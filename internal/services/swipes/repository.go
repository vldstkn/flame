package swipes

import (
	"context"
	"flame/internal/models"
	"flame/pkg/db"
)

type RepositoryDeps struct {
	DB    *db.DB
	Redis *db.Redis
}
type Repository struct {
	DB    *db.DB
	Redis *db.Redis
}

func NewRepository(deps *RepositoryDeps) *Repository {
	return &Repository{
		DB:    deps.DB,
		Redis: deps.Redis,
	}
}

func (repo *Repository) CreateOrUpdate(userId1, userId2 int64, isLike bool) error {
	var query string
	if userId1 > userId2 {
		id1 := userId1
		userId1 = userId2
		userId2 = id1
		query = "INSERT INTO swipes (user_id1, user_id2, user_is_liked2) VALUES($1,$2,$3) ON CONFLICT (user_id1, user_id2) DO UPDATE SET user_is_liked2=$3"
	} else {
		query = "INSERT INTO swipes (user_id1, user_id2, user_is_liked1) VALUES($1,$2,$3) ON CONFLICT (user_id1, user_id2) DO UPDATE SET user_is_liked1=$3"
	}
	_, err := repo.DB.Exec(query, userId1, userId2, isLike)
	return err
}

func (repo *Repository) GetUnreadSwipes(userId int64) []int64 {
	var ids []int64
	err := repo.DB.Select(&ids, `SELECT
    CASE 
        WHEN user_id1=$1 THEN user_id2
        WHEN user_id2=$1 THEN user_id1
    END FROM swipes 
        WHERE (user_id1=$1 AND user_is_liked1 IS NULL) 
           OR (user_id2=$1 AND user_is_liked2 IS NULL)`, userId)
	if err != nil {
		return nil
	}
	return ids
}

func (repo *Repository) GetSwipeById(userId1, userId2 int64) *models.Swipe {
	var swipe models.Swipe
	if userId1 > userId2 {
		id1 := userId1
		userId1 = userId2
		userId2 = id1
	}
	err := repo.DB.Get(&swipe, `SELECT * FROM swipes WHERE user_id1=$1 AND user_id2=$2`, userId1, userId2)
	if err != nil {
		return nil
	}
	return &swipe
}

func (repo *Repository) RemoveSwipeFromRedis(candidateListKey string, userId int64) error {
	err := repo.Redis.SRem(context.Background(), candidateListKey, userId).Err()
	return err
}
