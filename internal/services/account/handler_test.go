package account

import (
	"context"
	"flame/internal/config"
	"flame/internal/interfaces"
	"flame/pkg/logger"
	"flame/pkg/pb"
	"flame/tests/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"os"
	"testing"
)

const configPath = "../../../configs"
const mode = "test"

func TestHandler_Register(t *testing.T) {
	const accessToken = "access_token"
	const refreshToken = "refresh_token"
	log := logger.NewLogger(os.Stdout)
	conf := config.LoadConfig(configPath, mode)
	service := new(mocks.MockAccountService)
	handler := NewHandler(&HandlerDeps{
		Logger:  log,
		Config:  conf,
		Service: service,
	})
	type grpcRes struct {
		pb    *pb.RegisterRes
		isErr bool
	}
	validData := &pb.RegisterReq{
		Email:      "test@gmail.com",
		Password:   "123456",
		Name:       "test",
		Gender:     "male",
		LookingFor: "female",
	}
	tests := []struct {
		name    string
		request *pb.RegisterReq
		service func()
		res     grpcRes
	}{
		{
			name:    "success",
			request: validData,
			res: grpcRes{
				pb: &pb.RegisterRes{
					AccessToken:  accessToken,
					RefreshToken: refreshToken,
				},
				isErr: false,
			},
			service: func() {
				service.On("Register", mock.Anything).Return(1, nil)
				service.On("IssueToken", mock.Anything, mock.Anything).Return(&interfaces.AccountSIssueToken{
					AccessToken:  accessToken,
					RefreshToken: refreshToken,
				}, nil)
			},
		},
		{
			name:    "bad service register",
			request: validData,
			res: grpcRes{
				pb:    nil,
				isErr: true,
			},
			service: func() {
				service.On("Register", mock.Anything).
					Return(-1, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError)))
				service.On("IssueToken", mock.Anything, mock.Anything).
					Return(&interfaces.AccountSIssueToken{
						AccessToken:  accessToken,
						RefreshToken: refreshToken,
					}, nil)
			},
		},
		{
			name:    "bad service issue tokens",
			request: validData,
			res: grpcRes{
				pb:    nil,
				isErr: true,
			},
			service: func() {
				service.On("Register", mock.Anything).Return(1, nil)
				service.On("IssueToken", mock.Anything, mock.Anything).
					Return(nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError)))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.service()
			t.Cleanup(func() {
				service.ExpectedCalls = nil
			})
			response, err := handler.Register(context.Background(), tt.request)
			assert.Equal(t, err != nil, tt.res.isErr)
			assert.Equal(t, response, tt.res.pb)
		})
	}
}

func TestHandler_Login(t *testing.T) {
	const accessToken = "access_token"
	const refreshToken = "refresh_token"
	log := logger.NewLogger(os.Stdout)
	conf := config.LoadConfig(configPath, mode)
	service := new(mocks.MockAccountService)
	handler := NewHandler(&HandlerDeps{
		Logger:  log,
		Config:  conf,
		Service: service,
	})
	type grpcRes struct {
		pb    *pb.LoginRes
		isErr bool
	}
	validData := &pb.LoginReq{
		Email:    "test@gmail.com",
		Password: "123456",
	}
	tests := []struct {
		name    string
		request *pb.LoginReq
		service func()
		res     grpcRes
	}{
		{
			name:    "success",
			request: validData,
			res: grpcRes{
				pb: &pb.LoginRes{
					AccessToken:  accessToken,
					RefreshToken: refreshToken,
				},
				isErr: false,
			},
			service: func() {
				service.On("Login", mock.Anything, mock.Anything).Return(1, nil)
				service.On("IssueToken", mock.Anything, mock.Anything).Return(&interfaces.AccountSIssueToken{
					AccessToken:  accessToken,
					RefreshToken: refreshToken,
				}, nil)
			},
		},
		{
			name:    "bad service login",
			request: validData,
			res: grpcRes{
				pb:    nil,
				isErr: true,
			},
			service: func() {
				service.On("Login", mock.Anything, mock.Anything).
					Return(-1, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError)))
				service.On("IssueToken", mock.Anything, mock.Anything).
					Return(&interfaces.AccountSIssueToken{
						AccessToken:  accessToken,
						RefreshToken: refreshToken,
					}, nil)
			},
		},
		{
			name:    "bad service issue tokens",
			request: validData,
			res: grpcRes{
				pb:    nil,
				isErr: true,
			},
			service: func() {
				service.On("Login", mock.Anything, mock.Anything).Return(1, nil)
				service.On("IssueToken", mock.Anything, mock.Anything).
					Return(nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError)))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.service()
			t.Cleanup(func() {
				service.ExpectedCalls = nil
			})
			response, err := handler.Login(context.Background(), tt.request)
			assert.Equal(t, err != nil, tt.res.isErr)
			assert.Equal(t, response, tt.res.pb)
		})
	}
}

func TestHandler_GetTokens(t *testing.T) {
	const accessToken = "access_token"
	const refreshToken = "refresh_token"
	log := logger.NewLogger(os.Stdout)
	conf := config.LoadConfig(configPath, mode)
	service := new(mocks.MockAccountService)
	handler := NewHandler(&HandlerDeps{
		Logger:  log,
		Config:  conf,
		Service: service,
	})
	type grpcRes struct {
		pb    *pb.GetTokensRes
		isErr bool
	}
	tests := []struct {
		name    string
		request *pb.GetTokensReq
		service func()
		res     grpcRes
	}{
		{
			name: "success",
			request: &pb.GetTokensReq{
				Id: 1,
			},
			service: func() {
				service.On("GetTokens", mock.Anything, mock.Anything).Return(&interfaces.AccountSIssueToken{
					AccessToken:  accessToken,
					RefreshToken: refreshToken,
				}, nil)
			},
			res: grpcRes{
				pb: &pb.GetTokensRes{
					AccessToken:  accessToken,
					RefreshToken: refreshToken,
				},
				isErr: false,
			},
		},
		{
			name: "bad service get tokens",
			request: &pb.GetTokensReq{
				Id: 1,
			},
			service: func() {
				service.On("GetTokens", mock.Anything, mock.Anything).Return(nil, errors.New(""))
			},
			res: grpcRes{
				pb:    nil,
				isErr: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.service()
			t.Cleanup(func() {
				service.ExpectedCalls = nil
			})
			response, err := handler.GetTokens(context.Background(), tt.request)
			assert.Equal(t, err != nil, tt.res.isErr)
			assert.Equal(t, response, tt.res.pb)
		})
	}
}
