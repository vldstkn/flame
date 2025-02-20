package main

import (
	"flame/internal/api"
	"flame/internal/config"
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
	app := api.NewApp(&api.AppDeps{
		Config: conf,
		Logger: log,
		Mode:   mode,
	})
	err := app.Run()
	if err != nil {
		log.Error(err.Error(),
			slog.String("Address", conf.Services.Api.Address),
			slog.String("Mode", mode),
		)
	}
}
