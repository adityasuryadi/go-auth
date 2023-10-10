package service

type RedisService interface {
	Set(key string, value interface{})
	Get(key string) *[]byte
}
