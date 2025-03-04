package api

import (
	"bytes"
	"context"
	"encoding/json"
	"flame/internal/config"
	"flame/internal/services/api/dto"
	"flame/pkg/jwt"
	"flame/pkg/logger"
	"flame/pkg/pb"
	"flame/tests/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const configPath = "../../../../configs"

func TestHandlerRegister(t *testing.T) {
	const accessToken = "access_token"
	log := logger.NewLogger(os.Stdout)
	conf := config.LoadConfig(configPath, "test")
	apiService := new(mocks.MockApiService)
	accountClient := new(mocks.MockAccountClient)
	handler := AccountHandler{
		Logger:        log,
		Config:        conf,
		ApiService:    apiService,
		AccountClient: accountClient,
	}
	tests := []struct {
		name           string
		data           dto.AccountRegisterReq
		code           int
		accountService func()
	}{
		{
			name: "bad email",
			data: dto.AccountRegisterReq{
				Name:       "test",
				Email:      "gmail.com",
				Password:   "123456",
				Gender:     "male",
				LookingFor: "female",
			},
			code: 400,
		},
		{
			name: "short password",
			data: dto.AccountRegisterReq{
				Name:       "test",
				Email:      "test@gmail.com",
				Password:   "12345",
				Gender:     "male",
				LookingFor: "female",
			},
			code: 400,
		},
		{
			name: "bad gender",
			data: dto.AccountRegisterReq{
				Name:       "test",
				Email:      "test@gmail.com",
				Password:   "123456",
				Gender:     "bad gender",
				LookingFor: "female",
			},
			code: 400,
		},
		{
			name: "bad LookingFor",
			data: dto.AccountRegisterReq{
				Name:       "test",
				Email:      "test@gmail.com",
				Password:   "123456",
				Gender:     "male",
				LookingFor: "bad LookingFor",
			},
			code: 400,
		},
		{
			name: "success",
			data: dto.AccountRegisterReq{
				Name:       "test",
				Email:      "test@gmail.com",
				Password:   "123456",
				Gender:     "male",
				LookingFor: "female",
			},
			code: 201,
			accountService: func() {
				accountClient.On("Register", context.Background(), mock.Anything, mock.Anything).
					Return(&pb.RegisterRes{
						AccessToken: accessToken,
					}, nil)
			},
		},
		{
			name: "bad account client",
			data: dto.AccountRegisterReq{
				Name:       "test",
				Email:      "test@gmail.com",
				Password:   "123456",
				Gender:     "male",
				LookingFor: "female",
			},
			code: 500,
			accountService: func() {
				accountClient.On("Register", context.Background(), mock.Anything, mock.Anything).
					Return(nil, errors.New("internal"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.accountService != nil {
				tt.accountService()
				t.Cleanup(func() {
					accountClient.ExpectedCalls = nil
				})
			}
			body, err := json.Marshal(tt.data)
			if err != nil {
				t.Error(err)
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/", bytes.NewReader(body))
			handler.Register()(w, r)
			assert.Equal(t, tt.code, w.Result().StatusCode)
			if w.Result().StatusCode == http.StatusCreated {
				var data dto.AccountRegisterRes
				err = json.Unmarshal(w.Body.Bytes(), &data)
				if err != nil {
					t.Error(err)
				}
				require.Equal(t, accessToken, data.AccessToken)
			} else {
				var data dto.ErrorRes
				err = json.Unmarshal(w.Body.Bytes(), &data)
				if err != nil {
					t.Error(err)
				}
				if data.Error == "" {
					assert.Equal(t, data.Error, "")
				}
			}
			if tt.accountService != nil {
				accountClient.AssertExpectations(t)
			}
		})
	}
}

func TestApiHandlerLogin(t *testing.T) {
	const accessToken = "access_token"
	log := logger.NewLogger(os.Stdout)
	conf := config.LoadConfig(configPath, "test")
	apiService := new(mocks.MockApiService)
	accountClient := new(mocks.MockAccountClient)
	handler := AccountHandler{
		Logger:        log,
		Config:        conf,
		ApiService:    apiService,
		AccountClient: accountClient,
	}
	tests := []struct {
		name           string
		data           dto.AccountLoginReq
		code           int
		accountService func()
	}{
		{
			name: "success",
			data: dto.AccountLoginReq{
				Email:    "test@gmail.com",
				Password: "123456",
			},
			code: http.StatusOK,
			accountService: func() {
				accountClient.On("Login", context.Background(), mock.Anything, mock.Anything).
					Return(&pb.LoginRes{
						AccessToken: accessToken,
					}, nil)
			},
		},
		{
			name: "short password",
			data: dto.AccountLoginReq{
				Email:    "test@gmail.com",
				Password: "12345",
			},
			code: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			data: dto.AccountLoginReq{
				Email:    "gmail.com",
				Password: "123456",
			},
			code: http.StatusBadRequest,
		},
		{
			name: "bad account client",
			data: dto.AccountLoginReq{
				Email:    "test@gmail.com",
				Password: "123456",
			},
			code: http.StatusInternalServerError,
			accountService: func() {
				accountClient.On("Login", context.Background(), mock.Anything, mock.Anything).
					Return(nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError)))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.accountService != nil {
				tt.accountService()
				t.Cleanup(func() {
					accountClient.ExpectedCalls = nil
				})
			}
			body, err := json.Marshal(tt.data)
			if err != nil {
				t.Error(err)
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/", bytes.NewReader(body))
			handler.Login()(w, r)
			assert.Equal(t, tt.code, w.Result().StatusCode)
			if w.Result().StatusCode == http.StatusOK {
				var data dto.AccountRegisterRes
				err = json.Unmarshal(w.Body.Bytes(), &data)
				if err != nil {
					t.Error(err)
				}
				require.Equal(t, accessToken, data.AccessToken)
			} else {
				var data dto.ErrorRes
				err = json.Unmarshal(w.Body.Bytes(), &data)
				if err != nil {
					t.Error(err)
				}
				if data.Error == "" {
					assert.Equal(t, data.Error, "")
				}
			}
			if tt.accountService != nil {
				accountClient.AssertExpectations(t)
			}
		})
	}
}
func TestApiHandlerGetTokens(t *testing.T) {
	const accessToken = "access_token"
	const refreshToken = "refresh_token"

	log := logger.NewLogger(os.Stdout)
	conf := config.LoadConfig(configPath, "test")
	apiService := new(mocks.MockApiService)
	accountClient := new(mocks.MockAccountClient)
	handler := AccountHandler{
		Logger:        log,
		Config:        conf,
		ApiService:    apiService,
		AccountClient: accountClient,
	}
	validToken, err := jwt.NewJWT(conf.Auth.Jwt).Create(jwt.Data{
		Id: 1,
	}, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name           string
		tokenValue     string
		tokenName      string
		code           int
		accountService func()
	}{
		{
			name:       "success",
			tokenValue: validToken,
			tokenName:  refreshToken,
			code:       http.StatusOK,
			accountService: func() {
				accountClient.On("GetTokens", context.Background(), mock.Anything, mock.Anything).
					Return(&pb.GetTokensRes{
						AccessToken: accessToken,
					}, nil)
			},
		},
		{
			name:       "bad token",
			tokenValue: "bad_token",
			tokenName:  refreshToken,
			code:       http.StatusUnauthorized,
			accountService: func() {
				accountClient.On("GetTokens", context.Background(), mock.Anything, mock.Anything).
					Return(&pb.GetTokensRes{
						AccessToken: accessToken,
					}, nil)
			},
		},
		{
			name:       "token is empty",
			tokenValue: "empty",
			tokenName:  "empty",
			code:       http.StatusUnauthorized,
			accountService: func() {
				accountClient.On("GetTokens", context.Background(), mock.Anything, mock.Anything).
					Return(&pb.GetTokensRes{
						AccessToken: accessToken,
					}, nil)
			},
		},
		{
			name:       "bad account client",
			tokenValue: validToken,
			tokenName:  refreshToken,
			code:       http.StatusInternalServerError,
			accountService: func() {
				accountClient.On("GetTokens", context.Background(), mock.Anything, mock.Anything).
					Return(nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError)))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.accountService != nil {
				tt.accountService()
				t.Cleanup(func() {
					accountClient.ExpectedCalls = nil
				})
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			c := &http.Cookie{
				Name:  tt.tokenName,
				Value: tt.tokenValue,
			}
			r.AddCookie(c)

			handler.GetTokens()(w, r)
			assert.Equal(t, tt.code, w.Result().StatusCode)
			if w.Result().StatusCode == http.StatusOK {
				var data dto.AccountGetTokensRes
				err = json.Unmarshal(w.Body.Bytes(), &data)
				if err != nil {
					t.Error(err)
				}
				require.Equal(t, accessToken, data.AccessToken)
			} else {
				var data dto.ErrorRes
				err = json.Unmarshal(w.Body.Bytes(), &data)
				if err != nil {
					t.Error(err)
				}
				if data.Error == "" {
					assert.Equal(t, data.Error, "")
				}
			}
			if tt.accountService != nil {
				accountClient.AssertExpectations(t)
			}
		})
	}
}
