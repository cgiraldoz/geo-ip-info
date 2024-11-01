package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cgiraldoz/geo-ip-info/internal/cache"
	"github.com/cgiraldoz/geo-ip-info/internal/services/httpclient"
	"github.com/gofiber/fiber/v2/log"
	"io"
	"time"
)

type DefaultPrefetchDataService struct {
	cache cache.Cache
}

func NewDefaultPrefetchDataService(cache cache.Cache) *DefaultPrefetchDataService {
	return &DefaultPrefetchDataService{cache: cache}
}

func (pd *DefaultPrefetchDataService) PreFetchData(ctx context.Context) error {

	httpClient := httpclient.NewDefaultHttpClient(10 * time.Second)

	resp, err := httpClient.Get(ctx, "https://jsonplaceholder.typicode.com/todos/1")

	if err != nil {
		return fmt.Errorf("error fetching data: %w", err)
	}

	defer func(Body io.ReadCloser) {
		closingErr := Body.Close()
		if closingErr != nil {
			log.Fatalf("error closing response body: %v", closingErr)
		}
	}(resp.Body)

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("error decoding response body: %w", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}

	client := pd.cache.NewClient()

	if err := client.Set(ctx, "todos", jsonData, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("error setting data in cache: %w", err)
	}

	return nil
}
