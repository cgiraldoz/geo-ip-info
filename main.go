package main

import (
	"context"
	"github.com/cgiraldoz/geo-ip-info/cmd/cli"
	"github.com/cgiraldoz/geo-ip-info/internal/cache"
	"github.com/cgiraldoz/geo-ip-info/internal/services"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	preFetchService := newPreFetchService()

	if err := preFetchService.PreFetchData(ctx); err != nil {
		log.Fatalf("Error prefetching data: %v", err)
	}

	if err := cli.Execute(); err != nil {
		log.Fatalf("Error executing CLI: %v", err)
	}
}

func newPreFetchService() *services.DefaultPrefetchDataService {
	redisService := services.NewDefaultRedisService()
	cacheService := cache.NewDefaultCache(redisService)
	return services.NewDefaultPrefetchDataService(cacheService)
}
