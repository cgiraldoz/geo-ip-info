package httpclient

import (
	"context"
	"net/http"
)

type HttpClient interface {
	Get(ctx context.Context, url string) (*http.Response, error)
}
