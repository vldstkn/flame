package main

import (
	"flame/internal/config"
	"flame/internal/services/mathcing"
	"flame/pkg/db"
	"flame/pkg/logger"
	"github.com/go-redis/redis/v8"
	"os"
)

func main() {
	mode := os.Getenv("APP_ENV")
	if mode == "" {
		mode = "dev"
	}
	conf := config.LoadConfig("configs", mode)
	log := logger.NewLogger(os.Stdout)

	accDb := db.NewDb(conf.Database.Account.Dsn)
	swipesDb := db.NewDb(conf.Database.Swipes.Dsn)
	rdb := db.NewRedis(&redis.Options{
		Addr:     conf.GetRedisAddr(),
		Password: conf.Database.Redis.Password,
		DB:       conf.Database.Redis.Db,
		Username: conf.Database.Redis.Username,
	})

	app := mathcing.NewApp(&mathcing.AppDeps{
		Config:    conf,
		Logger:    log,
		AccountDB: accDb,
		SwipeDB:   swipesDb,
		Mode:      mode,
		Redis:     rdb,
	})
	app.Run()
}
