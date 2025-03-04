package account

import (
	"errors"
	"flame/internal/models"
	"flame/pkg/db"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepository_GetById(t *testing.T) {
	database, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer database.Close()
	sqlxDB := sqlx.NewDb(database, "postgres")
	repo := NewRepository(&db.DB{
		DB: sqlxDB,
	})
	user := &models.User{
		Id:         1,
		Email:      "test@gmail.com",
		Password:   "123456",
		Gender:     "male",
		LookingFor: "female",
		Name:       "test",
	}
	tests := []struct {
		name string
		id   int64
		res  *models.User
		db   func()
	}{
		{
			name: "success",
			id:   1,
			res:  user,
			db: func() {
				row := sqlmock.NewRows([]string{"id", "email", "password", "gender", "looking_for", "name"}).
					AddRow(user.Id, user.Email, user.Password, user.Gender, user.LookingFor, user.Name)
				mock.ExpectQuery("SELECT * FROM users WHERE").WillReturnRows(row)
			},
		},
		{
			name: "user does not exist",
			id:   1,
			res:  nil,
			db: func() {
				mock.ExpectQuery("SELECT * FROM users WHERE").WillReturnError(errors.New(""))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.db()
			res := repo.GetById(tt.id)
			assert.Equal(t, res, tt.res)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
