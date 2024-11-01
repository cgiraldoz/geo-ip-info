package services

import (
	"github.com/redis/go-redis/v9"
	"sync"
)

type DefaultRedisService struct {
	client *redis.Client
	once   sync.Once
}

func NewDefaultRedisService() *DefaultRedisService {
	return &DefaultRedisService{}
}

func (rc *DefaultRedisService) NewClient() *redis.Client {
	rc.once.Do(func() {
		rc.client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	})
	return rc.client
}
