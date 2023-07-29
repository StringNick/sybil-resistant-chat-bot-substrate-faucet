package redis

import (
	"substrate-faucet/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis(conf config.Redis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Endpoint,
		Password: "",
		DB:       0,
	})

	return rdb, nil
}
