package main

import (
	"flame/internal/config"
	"flame/internal/services/account"
	"flame/pkg/db"
	"flame/pkg/logger"
	"github.com/go-redis/redis/v8"
	"log/slog"
	"os"
)

func main() {
	mode := os.Getenv("APP_ENV")
	if mode == "" {
		mode = "dev"
	}
	conf := config.LoadConfig("configs", mode)
	log := logger.NewLogger(os.Stdout)
	database := db.NewDb(conf.Database.Account.Dsn)
	rdb := db.NewRedis(&redis.Options{
		Addr:     conf.GetRedisAddr(),
		Password: conf.Database.Redis.Password,
		DB:       conf.Database.Redis.Db,
		Username: conf.Database.Redis.Username,
	})
	app := account.NewApp(&account.AppDeps{
		Config: conf,
		Logger: log,
		Db:     database,
		Redis:  rdb,
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
