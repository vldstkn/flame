package account

import (
	"flame/internal/config"
	"flame/internal/interfaces"
	"flame/internal/models"
	"flame/pkg/jwt"
	"flame/pkg/logger"
	"flame/tests/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"os"
	"testing"
	"time"
)

func TestService_Register(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	log := logger.NewLogger(os.Stdout)
	service := NewService(&ServiceDeps{
		Repository: repo,
		Logger:     log,
	})
	type response struct {
		id    int64
		isErr bool
	}
	badPassword := ""
	for i := 0; i < 75; i++ {
		badPassword += "1"
	}
	validData := &interfaces.AccountSRegisterDeps{
		Name:       "test",
		Password:   "123456",
		Email:      "test@gmail.com",
		Gender:     "male",
		LookingFor: "female",
	}
	tests := []struct {
		name  string
		input *interfaces.AccountSRegisterDeps
		res   response
		repo  func()
	}{
		{
			name:  "success",
			input: validData,
			res: response{
				id:    1,
				isErr: false,
			},
			repo: func() {
				repo.On("GetByEmail", mock.Anything).Return(nil)
				repo.On("Create", mock.Anything).Return(1, nil)
			},
		},
		{
			name:  "user exists",
			input: validData,
			res: response{
				id:    -1,
				isErr: true,
			},
			repo: func() {
				repo.On("GetByEmail", mock.Anything).Return(&models.User{})
				repo.On("Create", mock.Anything).Return(1, nil)
			},
		},
		{
			name:  "bad data",
			input: validData,
			res: response{
				id:    -1,
				isErr: true,
			},
			repo: func() {
				repo.On("GetByEmail", mock.Anything).Return(nil)
				repo.On("Create", mock.Anything).Return(-1, errors.New(""))
			},
		},
		{
			name: "bad hash password",
			input: &interfaces.AccountSRegisterDeps{
				Name:       "test",
				Password:   badPassword,
				Email:      "test@gmail.com",
				Gender:     "male",
				LookingFor: "female",
			},
			res: response{
				id:    -1,
				isErr: true,
			},
			repo: func() {
				repo.On("GetByEmail", mock.Anything).Return(nil)
				repo.On("Create", mock.Anything).Return(1, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.repo()
			t.Cleanup(func() {
				repo.ExpectedCalls = nil
			})
			res, err := service.Register(tt.input)
			assert.Equal(t, err != nil, tt.res.isErr)
			assert.Equal(t, res, tt.res.id)
		})
	}
}

func TestService_IssueToken(t *testing.T) {
	log := logger.NewLogger(os.Stdout)
	conf := config.LoadConfig(configPath, mode)
	service := NewService(&ServiceDeps{
		Logger: log,
	})
	type response struct {
		tokens *interfaces.AccountSIssueToken
		isErr  bool
	}
	jwtData := jwt.Data{
		Id: 1,
	}
	j := jwt.NewJWT(conf.Auth.Jwt)
	validAccessToken, _ := j.Create(jwtData, time.Now().Add(time.Hour*2).Add(time.Minute*10))
	validRefreshToken, _ := j.Create(jwtData, time.Now().AddDate(0, 0, 2).Add(time.Hour*2))

	tests := []struct {
		name   string
		secret string
		data   jwt.Data
		res    response
	}{
		{
			name:   "success",
			secret: conf.Auth.Jwt,
			data:   jwtData,
			res: response{
				tokens: &interfaces.AccountSIssueToken{
					AccessToken:  validAccessToken,
					RefreshToken: validRefreshToken,
				},
				isErr: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := service.IssueToken(tt.secret, tt.data)
			assert.Equal(t, res, tt.res.tokens)
			assert.Equal(t, err != nil, tt.res.isErr)
		})
	}
}
func TestService_Login(t *testing.T) {
	log := logger.NewLogger(os.Stdout)
	repo := new(mocks.MockAccountRepository)
	service := NewService(&ServiceDeps{
		Logger:     log,
		Repository: repo,
	})
	type response struct {
		id    int64
		isErr bool
	}
	validInput := struct {
		password string
		email    string
	}{
		password: "123456",
		email:    "test@gmail.com",
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(validInput.password), bcrypt.DefaultCost)
	user := &models.User{
		Id:         1,
		Email:      "test@gmail.com",
		Password:   string(hashPassword),
		Gender:     "male",
		LookingFor: "female",
		Name:       "test",
	}
	tests := []struct {
		name     string
		email    string
		password string
		res      response
		repo     func()
	}{
		{
			name:     "success",
			email:    validInput.email,
			password: validInput.password,
			res: response{
				id:    1,
				isErr: false,
			},
			repo: func() {
				repo.On("GetByEmail", mock.Anything).Return(user)
			},
		},
		{
			name:     "user does not exist",
			email:    validInput.email,
			password: validInput.password,
			res: response{
				id:    -1,
				isErr: true,
			},
			repo: func() {
				repo.On("GetByEmail", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:     "bad password",
			email:    validInput.email,
			password: "bad password",
			res: response{
				id:    -1,
				isErr: true,
			},
			repo: func() {
				repo.On("GetByEmail", mock.Anything, mock.Anything).Return(user)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.repo != nil {
				tt.repo()
				t.Cleanup(func() {
					repo.ExpectedCalls = nil
				})
			}
			id, err := service.Login(tt.email, tt.password)
			assert.Equal(t, id, tt.res.id)
			assert.Equal(t, err != nil, tt.res.isErr)
		})

	}
}

func TestService_GetTokens(t *testing.T) {
	log := logger.NewLogger(os.Stdout)
	repo := new(mocks.MockAccountRepository)
	conf := config.LoadConfig(configPath, mode)
	service := NewService(&ServiceDeps{
		Logger:     log,
		Repository: repo,
	})
	type response struct {
		data  *interfaces.AccountSIssueToken
		isErr bool
	}
	jwtData := jwt.Data{
		Id: 1,
	}
	j := jwt.NewJWT(conf.Auth.Jwt)
	validAccessToken, _ := j.Create(jwtData, time.Now().Add(time.Hour*2).Add(time.Minute*10))
	validRefreshToken, _ := j.Create(jwtData, time.Now().AddDate(0, 0, 2).Add(time.Hour*2))
	user := &models.User{
		Id:         1,
		Email:      "test@gmail.com",
		Gender:     "male",
		LookingFor: "female",
		Name:       "test",
	}
	tests := []struct {
		name   string
		secret string
		data   jwt.Data
		res    response
		repo   func()
	}{
		{
			name:   "success",
			secret: conf.Auth.Jwt,
			data:   jwtData,
			res: response{
				data: &interfaces.AccountSIssueToken{
					AccessToken:  validAccessToken,
					RefreshToken: validRefreshToken,
				},
				isErr: false,
			},
			repo: func() {
				repo.On("GetById", mock.Anything).Return(user)
			},
		},
		{
			name:   "user does not exist",
			secret: conf.Auth.Jwt,
			data:   jwtData,
			res: response{
				data:  nil,
				isErr: true,
			},
			repo: func() {
				repo.On("GetById", mock.Anything, mock.Anything).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.repo != nil {
				tt.repo()
				t.Cleanup(func() {
					repo.ExpectedCalls = nil
				})
			}
			tokens, err := service.GetTokens(tt.secret, tt.data)
			assert.Equal(t, tokens, tt.res.data)
			assert.Equal(t, err != nil, tt.res.isErr)
		})

	}
}
