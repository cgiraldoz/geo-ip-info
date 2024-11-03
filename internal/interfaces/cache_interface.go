package interfaces

import (
	"context"
	"time"
)

type Cache interface {
	Exists(ctx context.Context, key string) (int64, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
}
