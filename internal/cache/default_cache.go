package cache

import "github.com/redis/go-redis/v9"

type DefaultCache struct {
	cache Cache
}

func NewDefaultCache(cache Cache) *DefaultCache {
	return &DefaultCache{cache: cache}
}

func (dc *DefaultCache) NewClient() *redis.Client {
	return dc.cache.NewClient()
}
