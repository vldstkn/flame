package account

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
	Db     *db.DB
	Mode   string
}
type App struct {
	Config *config.Config
	Logger *slog.Logger
	Db     *db.DB
	Mode   string
}

func NewApp(deps *AppDeps) *App {
	return &App{
		Config: deps.Config,
		Logger: deps.Logger,
		Db:     deps.Db,
		Mode:   deps.Mode,
	}
}

func (app *App) Run() error {
	var opts []grpc.ServerOption
	lis, err := net.Listen("tcp", app.Config.Services.Account.Address)
	if err != nil {
		app.Logger.Error(err.Error(),
			slog.String("Error location", "net.Listen"),
			slog.String("Account address", app.Config.Services.Account.Address),
		)
		return err
	}
	defer lis.Close()

	repository := NewRepository(app.Db)
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
	pb.RegisterAccountServer(server, handler)
	app.Logger.Info("Service starts",
		slog.String("Name", "Account"),
		slog.String("Address", app.Config.Services.Account.Address),
		slog.String("Mode", app.Mode),
	)
	err = server.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}
