package e2e

import (
	"flame/internal/models"
	"flame/internal/services/api/dto"
	"flame/pkg/jwt"
	"flame/pkg/req"
	"flame/tests/env"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"testing"
	"time"
)

const getTokens = "/auth/get-tokens"

func TestGetTokens(t *testing.T) {
	e, err := env.InitEnv()
	e.Up()
	defer e.Down()
	if err != nil {
		t.Fatal(err)
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	user := models.User{
		Email:      "test@gmail.com",
		Password:   string(hashPassword),
		Gender:     "male",
		LookingFor: "female",
		Name:       "test",
	}
	var id int64
	err = e.DB.QueryRow(`INSERT INTO users (email, password, gender, looking_for, name) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		user.Email, user.Password, user.Gender, user.LookingFor, user.Name).Scan(&id)
	if err != nil {
		t.Fatal(err)
	}
	j := jwt.NewJWT(e.Jwt)

	validToken, err := j.Create(jwt.Data{
		Id: id,
	}, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}

	badToken, err := j.Create(jwt.Data{
		Id: 99999,
	}, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		status int
		token  string
	}{
		{
			name:   "success",
			status: http.StatusOK,
			token:  validToken,
		},
		{
			name:   "user not found",
			status: http.StatusBadRequest,
			token:  badToken,
		},
		{
			name:   "invalid token",
			status: http.StatusUnauthorized,
			token:  "invalid token",
		},
	}
	client := http.Client{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest("GET", e.ApiAddress+getTokens, nil)
			if err != nil {
				t.Error(err)
			}
			r.AddCookie(&http.Cookie{
				Name:  "refresh_token",
				Value: tt.token,
			})
			response, err := client.Do(r)

			if err != nil {
				t.Error(err)
			}
			if response.StatusCode != tt.status {
				t.Errorf("expected %d got %d", tt.status, response.StatusCode)
			}
			if response.StatusCode == 200 {
				body, err := req.Decode[dto.AccountGetTokensRes](response.Body)
				if err != nil {
					t.Error(err)
				}
				if body.AccessToken == "" {
					t.Error("token is empty")
				}
			}
		})
	}
}
