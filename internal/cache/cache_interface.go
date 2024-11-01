package cache

import "github.com/redis/go-redis/v9"

type Cache interface {
	NewClient() *redis.Client
}
