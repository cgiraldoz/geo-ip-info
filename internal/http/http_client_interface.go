package http

import (
	"context"
	"net/http"
)

type Client interface {
	Get(ctx context.Context, url string) (*http.Response, error)
}
