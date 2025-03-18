package db

import (
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	*redis.Client
}

func NewRedis(opts *redis.Options) *Redis {
	rdb := redis.NewClient(opts)

	return &Redis{
		rdb,
	}
}
