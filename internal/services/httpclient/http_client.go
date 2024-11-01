package httpclient

import (
	"context"
	"net/http"
	"time"
)

type DefaultHttpClient struct {
	client *http.Client
}

func NewDefaultHttpClient(timeout time.Duration) *DefaultHttpClient {
	return &DefaultHttpClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (hc *DefaultHttpClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return hc.client.Do(req)
}
