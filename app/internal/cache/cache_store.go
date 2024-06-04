package internal_cache

import (
	"context"
	"time"
)

type CacheStore interface {
	Get(ctx context.Context, key string, result interface{}) bool
	GetAll(ctx context.Context, obj interface{}, key ...string) (map[string]any, bool)
	Put(ctx context.Context, key string, obj interface{}, expiration time.Duration) bool
	PutAll(ctx context.Context, values map[string]any, expiration time.Duration) bool
	Delete(ctx context.Context, key ...string) bool
}
