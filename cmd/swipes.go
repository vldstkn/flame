package main

import (
	"flame/internal/config"
	"flame/internal/services/swipes"
	"flame/pkg/db"
	"flame/pkg/logger"
	"log/slog"
	"os"

	"github.com/go-redis/redis/v8"
)

func main() {
	mode := os.Getenv("APP_ENV")
	if mode == "" {
		mode = "dev"
	}

	conf := config.LoadConfig("configs", mode)
	log := logger.NewLogger(os.Stdout)
	database := db.NewDb(conf.Database.Swipes.Dsn)
	rdb := db.NewRedis(&redis.Options{
		Addr:     conf.GetRedisAddr(),
		Username: conf.Database.Redis.Username,
		Password: conf.Database.Redis.Password,
		DB:       conf.Database.Redis.Db,
	})
	app := swipes.NewApp(&swipes.AppDeps{
		Config: conf,
		Logger: log,
		Mode:   mode,
		DB:     database,
		Redis:  rdb,
	})
	err := app.Run()
	if err != nil {
		log.Error(err.Error(),
			slog.String("Address", conf.Services.Api.Address),
			slog.String("Mode", mode),
		)
	}
}
