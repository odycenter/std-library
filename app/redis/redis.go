// Package redis redis操作封装
package redis

import (
	"context"
	"time"
)

type Redis interface {
	Close()
	Get(ctx context.Context, key string) (string, error)
	MultiGet(ctx context.Context, key ...string) (map[string]string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) (string, error)
	MultiSet(ctx context.Context, values map[string]interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) (int64, error)
}
