package account

import (
	"flame/internal/models"
	"flame/pkg/db"
	"fmt"
)

type Repository struct {
	DB *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{
		DB: database,
	}
}

func (repo *Repository) GetById(id int64) *models.User {
	var user models.User
	err := repo.DB.Get(&user, `SELECT FROM users WHERE id=$1`, id)
	if err != nil {
		return nil
	}
	return &user
}
func (repo *Repository) GetByEmail(email string) *models.User {
	var user models.User
	err := repo.DB.Get(&user, `SELECT FROM users WHERE email=$1`, email)
	if err != nil {
		return nil
	}
	fmt.Println(user, email)
	return &user
}
func (repo *Repository) Create(user *models.User) (int64, error) {
	var id int64
	row, err := repo.DB.NamedQuery(`INSERT INTO users (email, password, name, gender, looking_for) 
																   VALUES (:email,:password,:name,:gender,:looking_for) RETURNING id`, user)
	if err != nil {
		return -1, err
	}
	row.Scan(&id)
	return id, nil
}
