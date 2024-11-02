package main

import (
	"context"
	"fmt"
	"github.com/cgiraldoz/geo-ip-info/cmd/cli"
	"github.com/cgiraldoz/geo-ip-info/config"
	"github.com/cgiraldoz/geo-ip-info/internal/cache"
	"github.com/cgiraldoz/geo-ip-info/internal/http"
	"github.com/cgiraldoz/geo-ip-info/internal/services"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := config.LoadConfigurations(); err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	ctx, cancel := createContext()
	defer cancel()

	redisCache := createRedisCache()
	httpClient := createHTTPClient()

	preFetchService := services.NewDefaultPrefetchDataService(redisCache, httpClient)

	if err := preFetchService.PreFetchData(ctx); err != nil {
		log.Fatalf("Error prefetching data: %v", err)
	}

	if err := cli.Execute(redisCache, httpClient); err != nil {
		log.Fatalf("Error executing CLI: %v", err)
	}
}

func createContext() (context.Context, context.CancelFunc) {
	contextTimeout := viper.GetDuration("context.timeout")
	return context.WithTimeout(context.Background(), contextTimeout)
}

func createRedisCache() *cache.RedisCache {
	redisHost := viper.GetString("redis.host")
	redisPort := viper.GetInt("redis.port")
	redisAddress := fmt.Sprintf("%s:%d", redisHost, redisPort)
	redisPassword := viper.GetString("redis.password")
	redisDB := viper.GetInt("redis.db")
	return cache.NewRedisCache(redisAddress, redisPassword, redisDB)
}

func createHTTPClient() *http.DefaultHttpClient {
	httpTimeout := viper.GetDuration("http.timeout")
	return http.NewDefaultHttpClient(httpTimeout)
}
