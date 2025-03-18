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
	"google.golang.org/protobuf/encoding/protojson"
	"log/slog"
	"net/http"
)

type MatchingHandlerDeps struct {
	Logger *slog.Logger
	Config *config.Config
}
type MatchingHandler struct {
	Logger      *slog.Logger
	Config      *config.Config
	MatchClient pb.MatchingClient
}

func NewMatchingHandler(router chi.Router, deps *MatchingHandlerDeps) error {
	matchConn, err := grpc_conn.NewClientConn(deps.Config.Services.Matching.Address)
	if err != nil {
		deps.Logger.Error(err.Error(),
			slog.String("Error location", "NewAccountHandler.grpc_conn.NewClientConn"),
			slog.String("Account address", deps.Config.Services.Account.Address),
		)
		return err
	}
	accountClient := pb.NewMatchingClient(matchConn)
	if err != nil {
		deps.Logger.Error(err.Error(),
			slog.String("Error location", "NewAccountHandler.config.NewS3Client"),
		)
		return err
	}
	handler := &MatchingHandler{
		Logger:      deps.Logger,
		Config:      deps.Config,
		MatchClient: accountClient,
	}
	router.Route("/match", func(r chi.Router) {
		r.Use(middleware.IsAuthed(handler.Config.Auth.Jwt))
		r.Get("/", handler.getMatchingUsers())
	})
	return nil
	
}

func (handler *MatchingHandler) getMatchingUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value("authData").(middleware.AuthData).Id
		body, err := req.HandleBody[dto.GetMatchingReq](r)
		if err != nil {
			res.Json(w, dto.ErrorRes{
				Error: http.StatusText(http.StatusBadRequest),
			}, http.StatusBadRequest)
			return
		}
		response, err := handler.MatchClient.GetMatchingUsers(context.Background(), &pb.GetMatchingUsersReq{
			Id:       id,
			Location: body.Location,
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
		data, err := opts.Marshal(response)
		if err != nil {
			res.Json(w, dto.ErrorRes{
				Error: http.StatusText(http.StatusInternalServerError),
			}, http.StatusInternalServerError)
			return
		}
		res.ProtoJson(w, data, http.StatusOK)
	}
}
