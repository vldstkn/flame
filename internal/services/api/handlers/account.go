package api

import (
	"context"
	"flame/internal/config"
	"flame/internal/interfaces"
	"flame/internal/services/api/dto"
	"flame/internal/services/api/middleware"
	http_errors "flame/pkg/errors"
	grpc_conn "flame/pkg/grpc-conn"
	"flame/pkg/jwt"
	"flame/pkg/pb"
	"flame/pkg/req"
	"flame/pkg/res"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/go-chi/chi/v5"
	"google.golang.org/protobuf/encoding/protojson"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type AccountHandlerDeps struct {
	Logger     *slog.Logger
	Config     *config.Config
	ApiService interfaces.ApiService
}
type AccountHandler struct {
	Logger        *slog.Logger
	Config        *config.Config
	ApiService    interfaces.ApiService
	AccountClient pb.AccountClient
	S3Client      *s3.Client
}

const (
	refreshToken = "refresh_token"
)

func NewAccountHandler(router chi.Router, deps *AccountHandlerDeps) error {
	accountConn, err := grpc_conn.NewClientConn(deps.Config.Services.Account.Address)
	if err != nil {
		deps.Logger.Error(err.Error(),
			slog.String("Error location", "NewAccountHandler.grpc_conn.NewClientConn"),
			slog.String("Account address", deps.Config.Services.Account.Address),
		)
		return err
	}
	accountClient := pb.NewAccountClient(accountConn)
	s3Client, err := config.NewS3Client()
	if err != nil {
		deps.Logger.Error(err.Error(),
			slog.String("Error location", "NewAccountHandler.config.NewS3Client"),
		)
		return err
	}
	handler := &AccountHandler{
		Logger:        deps.Logger,
		Config:        deps.Config,
		ApiService:    deps.ApiService,
		AccountClient: accountClient,
		S3Client:      s3Client,
	}
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", handler.Register())
		r.Post("/login", handler.Login())
		r.Get("/get-tokens", handler.GetTokens())
	})
	router.Route("/user", func(r chi.Router) {
		r.Use(middleware.IsAuthed(handler.Config.Auth.Jwt))
		r.Put("/profile", handler.UpdateProfile())
		r.Get("/profile", handler.GetProfile())
		r.Put("/photo", handler.UploadPhoto())
		r.Delete("/photo", handler.DeletePhoto())
	})
	return nil
}

func (handler *AccountHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[dto.AccountRegisterReq](r)
		if err != nil {
			res.Json(w, dto.ErrorRes{
				Error: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		response, err := handler.AccountClient.Register(context.Background(), &pb.RegisterReq{
			Email:    body.Email,
			Password: body.Password,
			Name:     body.Name,
		})
		if err != nil {
			msg, code := http_errors.HandleError(err)
			res.Json(w, dto.ErrorRes{
				Error: msg,
			}, code)
			return
		}
		handler.ApiService.AddCookie(&w, refreshToken, response.RefreshToken, int((time.Hour * 6).Seconds()))
		res.Json(w, dto.AccountRegisterRes{
			AccessToken: response.AccessToken,
		}, http.StatusCreated)
	}
}
func (handler *AccountHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[dto.AccountLoginReq](r)
		if err != nil {
			res.Json(w, dto.ErrorRes{
				Error: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		response, err := handler.AccountClient.Login(context.Background(), &pb.LoginReq{
			Email:    body.Email,
			Password: body.Password,
		})
		if err != nil {
			msg, code := http_errors.HandleError(err)
			res.Json(w, dto.ErrorRes{
				Error: msg,
			}, code)
			return
		}
		handler.ApiService.AddCookie(&w, refreshToken, response.RefreshToken, int((time.Hour * 6).Seconds()))
		res.Json(w, dto.AccountLoginRes{
			AccessToken: response.AccessToken,
		}, http.StatusOK)
	}
}
func (handler *AccountHandler) GetTokens() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(refreshToken)
		if err != nil {
			res.Json(w, dto.ErrorRes{
				Error: http.StatusText(http.StatusUnauthorized),
			}, http.StatusUnauthorized)
			return
		}
		isValid, data := jwt.NewJWT(handler.Config.Auth.Jwt).Parse(c.Value)
		if !isValid {
			res.Json(w, dto.ErrorRes{
				Error: http.StatusText(http.StatusUnauthorized),
			}, http.StatusUnauthorized)
			return
		}
		response, err := handler.AccountClient.GetTokens(context.Background(), &pb.GetTokensReq{
			Id: data.Id,
		})
		if err != nil {
			mes, code := http_errors.HandleError(err)
			res.Json(w, dto.ErrorRes{
				Error: mes,
			}, code)
			return
		}
		handler.ApiService.AddCookie(&w, refreshToken, response.RefreshToken, int((time.Hour * 6).Seconds()))
		res.Json(w, dto.AccountGetTokensRes{
			AccessToken: response.AccessToken,
		}, http.StatusOK)
	}
}

func (handler *AccountHandler) UpdateProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[dto.AccountUpdateProfileReq](r)
		if err != nil {
			res.Json(w, dto.ErrorRes{
				Error: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		authData := r.Context().Value("authData").(middleware.AuthData)

		_, err = handler.AccountClient.UpdateProfile(context.Background(), &pb.UpdateProfileReq{
			Id:        authData.Id,
			Name:      body.Name,
			BirthDate: body.BirthDate,
			City:      body.City,
			Bio:       body.Bio,
			Gender:    body.Gender,
		})
		if err != nil {
			mes, code := http_errors.HandleError(err)
			res.Json(w, dto.ErrorRes{
				Error: mes,
			}, code)
			return
		}
		res.Json(w, nil, http.StatusOK)
	}
}

func (handler *AccountHandler) GetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authData := r.Context().Value("authData").(middleware.AuthData)
		resGrpc, err := handler.AccountClient.GetProfile(context.Background(), &pb.GetProfileReq{
			Id: authData.Id,
		})
		if err != nil {
			mes, code := http_errors.HandleError(err)
			res.Json(w, dto.ErrorRes{
				Error: mes,
			}, code)
			return
		}
		opts := protojson.MarshalOptions{
			EmitUnpopulated: true,
		}
		jsonData, _ := opts.Marshal(resGrpc)
		res.ProtoJson(w, jsonData, http.StatusOK)
	}
}

func (handler *AccountHandler) UploadPhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authData := r.Context().Value("authData").(middleware.AuthData)
		err := r.ParseMultipartForm(5 << 20)
		if err != nil {
			res.Json(w, dto.ErrorRes{
				Error: http.StatusText(http.StatusBadRequest),
			}, http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("photo")
		if err != nil {
			res.Json(w, dto.ErrorRes{
				Error: http.StatusText(http.StatusBadRequest),
			}, http.StatusBadRequest)
			return
		}
		exp := strings.Split(fileHeader.Filename, ".")
		if exp[len(exp)-1] != "jpg" && exp[len(exp)-1] != "png" {
			res.Json(w, dto.ErrorRes{
				Error: http.StatusText(http.StatusBadRequest),
			}, http.StatusBadRequest)
			return
		}
		defer file.Close()
		uniqueFileName := fmt.Sprintf("%d-%d-%s", authData.Id, time.Now().UnixNano(), strings.Trim(fileHeader.Filename, "/"))
		_, err = handler.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: &handler.Config.S3.Bucket,
			Key:    &uniqueFileName,
			ACL:    types.ObjectCannedACLPublicRead,
			Body:   file,
		})
		if err != nil {
			fmt.Println(err)
			res.Json(w, dto.ErrorRes{
				Error: http.StatusText(http.StatusBadRequest),
			}, http.StatusBadRequest)
			return
		}

		_, err = handler.AccountClient.UploadPhoto(context.Background(), &pb.UploadPhotoReq{
			UserId: authData.Id,
			LinkPhoto: fmt.Sprintf("%s/%s/%s",
				handler.Config.S3.Endpoint,
				handler.Config.S3.Bucket,
				uniqueFileName),
		})
		if err != nil {
			msg, code := http_errors.HandleError(err)
			res.Json(w, dto.ErrorRes{
				Error: msg,
			}, code)
			return
		}

		res.Json(w, nil, http.StatusOK)
	}
}

func (handler *AccountHandler) DeletePhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[dto.DeletePhotoReq](r)
		authData := r.Context().Value("authData").(middleware.AuthData)
		if err != nil {
			res.Json(w, dto.ErrorRes{
				Error: err.Error(),
			}, http.StatusBadRequest)
		}
		response, err := handler.AccountClient.DeletePhoto(context.Background(), &pb.DeletePhotoReq{
			PhotoId: body.PhotoId,
			UserId:  authData.Id,
		})
		if err != nil {
			mes, code := http_errors.HandleError(err)
			res.Json(w, dto.ErrorRes{
				Error: mes,
			}, code)
			return
		}
		filename := strings.Split(response.PhotoUrl, "/")
		_, err = handler.S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: &handler.Config.S3.Bucket,
			Key:    &filename[len(filename)-1],
		})
		if err != nil {
			res.Json(w, dto.ErrorRes{
				Error: http.StatusText(http.StatusInternalServerError),
			}, http.StatusInternalServerError)
			return
		}
		res.Json(w, nil, http.StatusOK)
	}
}
