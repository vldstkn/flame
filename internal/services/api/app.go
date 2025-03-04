package api

import (
	"flame/internal/config"
	"flame/internal/services/api/handlers"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type AppDeps struct {
	Config *config.Config
	Logger *slog.Logger
	Mode   string
}
type App struct {
	Config *config.Config
	Logger *slog.Logger
	Mode   string
}

func NewApp(deps *AppDeps) *App {
	return &App{
		Config: deps.Config,
		Logger: deps.Logger,
		Mode:   deps.Mode,
	}
}

func (app *App) Run() error {
	router := chi.NewRouter()

	service := NewService(&ServiceDeps{})
	router.Route("/api", func(r chi.Router) {
		api.InitHandlers(r, &api.HandlersDeps{
			ApiService: service,
			Logger:     app.Logger,
			Config:     app.Config,
		})
	})

	server := http.Server{
		Addr:    app.Config.Services.Api.Address,
		Handler: router,
	}
	app.Logger.Info("Service starts",
		slog.String("Name", "Api"),
		slog.String("Address", app.Config.Services.Api.Address),
		slog.String("Mode", app.Mode),
	)
	defer server.Close()
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
