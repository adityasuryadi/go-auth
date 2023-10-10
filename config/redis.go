package config

import (
	"runtime"

	"github.com/gofiber/storage/redis/v3"
)

func NewRedis() *redis.Storage {
	store := redis.New(redis.Config{
		Host:      "redis-service",
		Port:      6379,
		Username:  "",
		Password:  "",
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})
	return store
}
