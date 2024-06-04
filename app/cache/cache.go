package cache

import (
	"context"
)

type Cache interface {
	Get(ctx context.Context, key string, obj interface{}, f func(key string) (interface{}, error)) error
	GetAll(ctx context.Context, key []string, obj interface{}, f func(key string) (interface{}, error)) (map[string]interface{}, error)
	Put(ctx context.Context, key string, obj any) bool
	PutAll(ctx context.Context, values map[string]any) bool
	Evict(ctx context.Context, key string) bool
	EvictAll(ctx context.Context, key ...string) bool
}
