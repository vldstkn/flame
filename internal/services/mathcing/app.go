package mathcing

import (
	"flame/internal/config"
	"flame/pkg/db"
	"flame/pkg/pb"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type AppDeps struct {
	Config    *config.Config
	Logger    *slog.Logger
	AccountDB *db.DB
	SwipeDB   *db.DB
	Redis     *db.Redis
	Mode      string
}
type App struct {
	Config    *config.Config
	Logger    *slog.Logger
	AccountDB *db.DB
	SwipeDB   *db.DB
	Redis     *db.Redis
	Mode      string
}

func NewApp(deps *AppDeps) *App {
	return &App{
		Config:    deps.Config,
		Logger:    deps.Logger,
		AccountDB: deps.AccountDB,
		SwipeDB:   deps.SwipeDB,
		Mode:      deps.Mode,
		Redis:     deps.Redis,
	}
}

func (app *App) Run() error {
	var opts []grpc.ServerOption
	lis, err := net.Listen("tcp", app.Config.Services.Matching.Address)
	if err != nil {
		app.Logger.Error(err.Error(),
			slog.String("Error location", "net.Listen"),
			slog.String("Matching address", app.Config.Services.Matching.Address),
		)
		return err
	}
	defer lis.Close()

	repository := NewRepository(&RepositoryDeps{
		AccountDB: app.AccountDB,
		SwipesDB:  app.SwipeDB,
	})
	service := NewService(&ServiceDeps{
		Repository: repository,
		Logger:     app.Logger,
		Redis:      app.Redis,
	})
	handler := NewHandler(&HandlerDeps{
		Logger:  app.Logger,
		Config:  app.Config,
		Service: service,
	})
	server := grpc.NewServer(opts...)
	defer server.Stop()
	pb.RegisterMatchingServer(server, handler)
	app.Logger.Info("Service starts",
		slog.String("Name", "Matching"),
		slog.String("Address", app.Config.Services.Matching.Address),
		slog.String("Mode", app.Mode),
	)
	err = server.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}
