package config

import (
	"runtime"

	"github.com/gofiber/storage/redis/v3"
)

func NewRedis(configuration Config) *redis.Storage {
	store := redis.New(redis.Config{
		Host:      configuration.Get("REDIS_HOST"),
		Port:      6379,
		Username:  configuration.Get("REDIS_USER"),
		Password:  configuration.Get("REDIS_PASSWORD"),
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})
	return store
}
