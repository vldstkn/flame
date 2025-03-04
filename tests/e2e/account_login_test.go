package e2e

import (
	"bytes"
	"encoding/json"
	"flame/internal/models"
	dto2 "flame/internal/services/api/dto"
	http_errors "flame/pkg/errors"
	"flame/pkg/req"
	"flame/tests/env"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"testing"
)

const login = "/auth/login"

func TestLogin(t *testing.T) {
	e, err := env.InitEnv()
	e.Up()
	defer e.Down()
	if err != nil {
		t.Fatal(err)
	}
	password := "123456"
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{
		Email:      "test@gmail.com",
		Password:   string(hashPassword),
		Gender:     "male",
		LookingFor: "female",
		Name:       "test",
	}
	_, err = e.DB.Exec(`INSERT INTO users (email, password, gender, looking_for, name) VALUES ($1,$2,$3,$4,$5)`,
		user.Email, user.Password, user.Gender, user.LookingFor, user.Name)

	tests := []struct {
		name   string
		input  dto2.AccountLoginReq
		status int
		err    dto2.ErrorRes
	}{
		{
			name: "success",
			input: dto2.AccountLoginReq{
				Email:    user.Email,
				Password: password,
			},
			status: http.StatusOK,
		},
		{
			name: "invalid email",
			input: dto2.AccountLoginReq{
				Email:    "bad@gmail.com",
				Password: password,
			},
			status: http.StatusBadRequest,
			err: dto2.ErrorRes{
				Error: http_errors.InvalidNameOrPassword,
			},
		},
		{
			name: "invalid password",
			input: dto2.AccountLoginReq{
				Email:    user.Email,
				Password: "badpassword",
			},
			status: http.StatusBadRequest,
			err: dto2.ErrorRes{
				Error: http_errors.InvalidNameOrPassword,
			},
		},
	}
	client := http.Client{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			if err != nil {
				t.Error(err)
			}
			r, err := http.NewRequest("POST", e.ApiAddress+login, bytes.NewReader(data))
			if err != nil {
				t.Error(err)
			}

			response, err := client.Do(r)
			if err != nil {
				t.Error(err)
			}
			if response.StatusCode != tt.status {
				t.Errorf("%s: expected %d got %d", tt.name, tt.status, response.StatusCode)
			}
			if response.StatusCode != 200 {
				body, err := req.Decode[dto2.ErrorRes](response.Body)
				if err != nil {
					t.Error(err)
				}
				if body != tt.err {
					t.Errorf("error: expected %+v got %+v", tt.err, body)
				}
			}
		})
	}
}
