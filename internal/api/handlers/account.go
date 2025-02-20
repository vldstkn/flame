package api

import (
	"context"
	"flame/internal/api/dto"
	"flame/internal/config"
	"flame/internal/interfaces"
	"flame/internal/models"
	http_errors "flame/pkg/errors"
	grpc_conn "flame/pkg/grpc-conn"
	"flame/pkg/pb"
	"flame/pkg/req"
	"flame/pkg/res"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
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
}

func NewAccountHAndler(router chi.Router, deps *AccountHandlerDeps) error {
	acccountConn, err := grpc_conn.NewClientConn(deps.Config.Services.Account.Address)
	if err != nil {
		deps.Logger.Error(err.Error(),
			slog.String("Error location", "NewAccountHAndler.grpc_conn.NewClientConn"),
			slog.String("Account address", deps.Config.Services.Account.Address),
		)
		return err
	}
	accountClient := pb.NewAccountClient(acccountConn)
	handler := &AccountHandler{
		Logger:        deps.Logger,
		Config:        deps.Config,
		ApiService:    deps.ApiService,
		AccountClient: accountClient,
	}
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", handler.Register())
		r.Post("/login", handler.Login())
		r.Get("/get-tokens", handler.GetTokens())
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
		if !models.GenderIsValid(body.Gender) {
			res.Json(w, dto.ErrorRes{
				Error: fmt.Sprintf("gender can not be '%s'", body.Gender),
			}, http.StatusBadRequest)
			return
		}
		if !models.GenderIsValid(body.LookingFor) {
			res.Json(w, dto.ErrorRes{
				Error: fmt.Sprintf("looking_for can not be '%s'", body.LookingFor),
			}, http.StatusBadRequest)
			return
		}
		response, err := handler.AccountClient.Register(context.Background(), &pb.RegisterReq{
			Email:      body.Email,
			Password:   body.Password,
			Name:       body.Name,
			Gender:     body.Gender,
			LookingFor: body.LookingFor,
		})
		if err != nil {
			msg, code := http_errors.HandleError(err)
			res.Json(w, dto.ErrorRes{
				Error: msg,
			}, code)
			return
		}
		handler.ApiService.AddCookie(&w, "refresh_token", response.RefreshToken, int((time.Hour * 6).Seconds()))
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
		handler.ApiService.AddCookie(&w, "refresh_token", response.RefreshToken, int((time.Hour * 6).Seconds()))
		res.Json(w, dto.AccountLoginRes{
			AccessToken: response.AccessToken,
		}, http.StatusOK)
	}
}
func (handler *AccountHandler) GetTokens() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetTokens")
		handler.AccountClient.GetTokens(context.Background(), &pb.GetTokensReq{
			Id: -1,
		})
	}
}
