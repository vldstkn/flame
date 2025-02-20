package main

import (
	"flame/internal/account"
	"flame/internal/config"
	"flame/pkg/db"
	"flame/pkg/logger"
	"log/slog"
	"os"
)

func main() {
	mode := os.Getenv("APP_ENV")
	if mode == "" {
		mode = "dev"
	}
	conf := config.LoadConfig("./configs", mode)
	log := logger.NewLogger(os.Stdout)
	database := db.NewDb(conf.Database.Dsn)
	app := account.NewApp(&account.AppDeps{
		Config: conf,
		Logger: log,
		Db:     database,
		Mode:   mode,
	})
	err := app.Run()
	if err != nil {
		log.Info(err.Error(),
			slog.String("Address", conf.Services.Api.Address),
			slog.String("Mode", mode),
		)
	}
}
