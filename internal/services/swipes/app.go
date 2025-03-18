package swipes

import (
	"flame/internal/config"
	"flame/pkg/db"
	"flame/pkg/pb"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type AppDeps struct {
	Config *config.Config
	Logger *slog.Logger
	DB     *db.DB
	Mode   string
	Redis  *db.Redis
}
type App struct {
	Config *config.Config
	Logger *slog.Logger
	DB     *db.DB
	Mode   string
	Redis  *db.Redis
}

func NewApp(deps *AppDeps) *App {
	return &App{
		Config: deps.Config,
		Logger: deps.Logger,
		DB:     deps.DB,
		Mode:   deps.Mode,
		Redis:  deps.Redis,
	}
}

func (app *App) Run() error {
	var opts []grpc.ServerOption
	lis, err := net.Listen("tcp", app.Config.Services.Swipes.Address)
	if err != nil {
		app.Logger.Error(err.Error(),
			slog.String("Error location", "net.Listen"),
			slog.String("Swipes address", app.Config.Services.Swipes.Address),
		)
		return err
	}
	defer lis.Close()

	repository := NewRepository(&RepositoryDeps{
		DB:    app.DB,
		Redis: app.Redis,
	})
	service := NewService(&ServiceDeps{
		Repository: repository,
		Logger:     app.Logger,
	})
	handler := NewHandler(&HandlerDeps{
		Logger:  app.Logger,
		Config:  app.Config,
		Service: service,
	})
	server := grpc.NewServer(opts...)
	defer server.Stop()
	pb.RegisterSwipesServer(server, handler)
	app.Logger.Info("Service starts",
		slog.String("Name", "Swipes"),
		slog.String("Address", app.Config.Services.Swipes.Address),
		slog.String("Mode", app.Mode),
	)
	err = server.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}
