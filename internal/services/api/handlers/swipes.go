package api

import (
	"context"
	"flame/internal/config"
	"flame/internal/services/api/dto"
	"flame/internal/services/api/middleware"
	http_errors "flame/pkg/errors"
	grpc_conn "flame/pkg/grpc-conn"
	"flame/pkg/pb"
	"flame/pkg/req"
	"flame/pkg/res"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type SwipesHandlerDeps struct {
	Logger *slog.Logger
	Config *config.Config
}
type SwipesHandler struct {
	Logger       *slog.Logger
	Config       *config.Config
	SwipesClient pb.SwipesClient
}

func NewSwipesHandler(router chi.Router, deps *SwipesHandlerDeps) error {
	swipesConn, err := grpc_conn.NewClientConn(deps.Config.Services.Swipes.Address)
	if err != nil {
		deps.Logger.Error(err.Error(),
			slog.String("Error location", "NewAccountHandler.grpc_conn.NewClientConn"),
			slog.String("Swipes address", deps.Config.Services.Swipes.Address),
		)
		return err
	}
	swipesClient := pb.NewSwipesClient(swipesConn)
	if err != nil {
		deps.Logger.Error(err.Error(),
			slog.String("Error location", "NewAccountHandler.config.NewS3Client"),
		)
		return err
	}
	handler := &SwipesHandler{
		Logger:       deps.Logger,
		Config:       deps.Config,
		SwipesClient: swipesClient,
	}
	router.Route("/swipes", func(r chi.Router) {
		r.Use(middleware.IsAuthed(handler.Config.Auth.Jwt))
		r.Post("/", handler.CreateSwipe())
		r.Get("/unread", handler.GetUnreadSwipes())
	})
	return nil

}

func (handler *SwipesHandler) CreateSwipe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId1 := r.Context().Value("authData").(middleware.AuthData).Id
		body, err := req.HandleBody[dto.CreateSwipesReq](r)
		if err != nil {
			res.Json(w, dto.ErrorRes{
				Error: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		_, err = handler.SwipesClient.CreateOrUpdateSwipe(context.Background(), &pb.CreateOrUpdateSwipeReq{
			UserId1: userId1,
			UserId2: body.UserId,
			IsLike:  *body.IsLike,
		})
		if err != nil {
			mes, code := http_errors.HandleError(err)
			res.Json(w, dto.ErrorRes{
				Error: mes,
			}, code)
			return
		}
		res.Json(w, nil, http.StatusCreated)
	}
}

func (handler *SwipesHandler) GetUnreadSwipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("authData").(middleware.AuthData).Id
		response, err := handler.SwipesClient.GetUnreadSwipes(context.Background(), &pb.GetUnreadSwipesReq{
			UserId: userId,
		})
		if err != nil {
			mes, code := http_errors.HandleError(err)
			res.Json(w, dto.ErrorRes{
				Error: mes,
			}, code)
			return
		}
		res.Json(w, dto.GetUnreadSwipes{
			Users: response.UserIds,
		}, http.StatusOK)
	}
}
