package service

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/gofiber/storage/redis/v3"
)

func NewRedisConfig(redis *redis.Storage, time float64) RedisService {
	var timeCfg float64
	timeCfg = 7200
	if time > 0 {
		timeCfg = time
	}
	return &RedisConfigImpl{
		Redis: redis,
		Time:  timeCfg,
	}
}

type RedisConfigImpl struct {
	Time  float64
	Redis *redis.Storage
}

// Delete implements RedisService.
func (service *RedisConfigImpl) Delete(key string) error {
	err := service.Redis.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

// Get implements RedisConfig.
func (service *RedisConfigImpl) Get(key string) *[]byte {
	result, err := service.Redis.Get(key)
	if err != nil {
		return nil
	}
	return &result
}

// Set implements RedisConfig.
func (service *RedisConfigImpl) Set(key string, value interface{}) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value)
	if err != nil {
		// return nil, err
		log.Fatal(err)
	}
	service.Redis.Set(key, buf.Bytes(), time.Hour*time.Duration(service.Time))
}
