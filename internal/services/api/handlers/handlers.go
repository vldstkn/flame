package api

import (
	"flame/internal/config"
	"flame/internal/interfaces"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

type HandlersDeps struct {
	ApiService interfaces.ApiService
	Logger     *slog.Logger
	Config     *config.Config
}

func InitHandlers(router chi.Router, deps *HandlersDeps) {
	_ = NewAccountHandler(router, &AccountHandlerDeps{
		Logger:     deps.Logger,
		Config:     deps.Config,
		ApiService: deps.ApiService,
	})
	_ = NewMatchingHandler(router, &MatchingHandlerDeps{
		Logger: deps.Logger,
		Config: deps.Config,
	})
	_ = NewSwipesHandler(router, &SwipesHandlerDeps{
		Logger: deps.Logger,
		Config: deps.Config,
	})
}
