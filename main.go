package main

import (
	"context"
	"github.com/cgiraldoz/geo-ip-info/cmd/cli"
	"github.com/cgiraldoz/geo-ip-info/internal/cache"
	"github.com/cgiraldoz/geo-ip-info/internal/http"
	"github.com/cgiraldoz/geo-ip-info/internal/services"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	redisCache := cache.NewRedisCache("localhost:6379", "", 0)

	httpClient := http.NewDefaultHttpClient(10 * time.Second)

	preFetchService := services.NewDefaultPrefetchDataService(redisCache, httpClient)

	if err := preFetchService.PreFetchData(ctx); err != nil {
		log.Fatalf("Error prefetching data: %v", err)
	}

	if err := cli.Execute(redisCache); err != nil {
		log.Fatalf("Error executing CLI: %v", err)
	}
}
