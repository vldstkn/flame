package e2e

import (
	"bytes"
	"encoding/json"
	dto2 "flame/internal/services/api/dto"
	http_errors "flame/pkg/errors"
	"flame/pkg/req"
	"flame/tests/env"
	"net/http"
	"testing"
)

const register = "/auth/register"

func TestRegister(t *testing.T) {
	e, err := env.InitEnv()
	e.Up()
	defer e.Down()
	if err != nil {
		t.Fatal(err)
	}
	_, err = e.DB.Exec(`INSERT INTO users (email, password, gender, looking_for, name) VALUES ($1, $2, $3, $4, $5)`,
		"exists@gmail.com", "123456", "female", "male", "test2")

	tests := []struct {
		name   string
		input  dto2.AccountRegisterReq
		status int
		err    dto2.ErrorRes
	}{
		{
			name: "success",
			input: dto2.AccountRegisterReq{
				Name:       "test",
				Email:      "test@gmail.com",
				Password:   "123456",
				Gender:     "male",
				LookingFor: "female",
			},
			status: http.StatusCreated,
		},
		{
			name: "user exists",
			input: dto2.AccountRegisterReq{
				Name:       "test",
				Email:      "exists@gmail.com",
				Password:   "123456",
				Gender:     "male",
				LookingFor: "female",
			},
			status: http.StatusBadRequest,
			err: dto2.ErrorRes{
				Error: http_errors.UserExists,
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
			r, err := http.NewRequest("POST", e.ApiAddress+register, bytes.NewReader(data))
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
			if response.StatusCode != 201 {
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
