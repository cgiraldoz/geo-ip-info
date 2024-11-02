package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cgiraldoz/geo-ip-info/internal/interfaces"
	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/viper"
	"io"
	"sync"
	"time"
)

type DefaultPrefetchDataService struct {
	cache      interfaces.Cache
	httpClient interfaces.Client
}

func NewDefaultPrefetchDataService(cache interfaces.Cache, httpClient interfaces.Client) *DefaultPrefetchDataService {
	return &DefaultPrefetchDataService{cache: cache, httpClient: httpClient}
}

func (pd *DefaultPrefetchDataService) PreFetchData(ctx context.Context) error {
	urlConfigs := viper.GetStringMap("prefetch.urls")

	urls := make(map[string]struct {
		url string
		ttl time.Duration
	})

	for key, config := range urlConfigs {
		conf := config.(map[string]interface{})
		url := conf["url"].(string)
		ttl, err := time.ParseDuration(conf["ttl"].(string))
		if err != nil {
			return fmt.Errorf("invalid TTL format for %s: %v", key, err)
		}
		urls[key] = struct {
			url string
			ttl time.Duration
		}{url: url, ttl: ttl}
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(urls))

	for key, config := range urls {
		wg.Add(1)

		go func(key, url string, ttl time.Duration) {
			defer wg.Done()

			exists, err := pd.cache.Exists(ctx, key)
			if err != nil {
				errCh <- fmt.Errorf("error checking existence in cache: %w", err)
				return
			}

			if exists > 0 {
				return
			}

			resp, err := pd.httpClient.Get(ctx, url)
			if err != nil {
				errCh <- fmt.Errorf("error fetching data from %s: %w", url, err)
				return
			}
			defer func(Body io.ReadCloser) {
				closingErr := Body.Close()
				if closingErr != nil {
					log.Fatalf("error closing response body: %v", closingErr)
				}
			}(resp.Body)

			var jsonData []byte

			if key == "countries" {
				var data []map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					errCh <- fmt.Errorf("error decoding response body from %s: %w", url, err)
					return
				}
				jsonData, err = json.Marshal(data)
			} else if key == "currencies" {
				var data map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					errCh <- fmt.Errorf("error decoding response body from %s: %w", url, err)
					return
				}
				jsonData, err = json.Marshal(data)
			}

			if err != nil {
				errCh <- fmt.Errorf("error marshalling data from %s: %w", url, err)
				return
			}

			if err := pd.cache.Set(ctx, key, jsonData, ttl); err != nil {
				errCh <- fmt.Errorf("error setting data in cache for key %s: %w", key, err)
				return
			}

		}(key, config.url, config.ttl)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}
