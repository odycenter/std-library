package internal_cache

import (
	"context"
	"time"
)

type LocalCacheStore struct {
	MaxSize int
}

func (c *LocalCacheStore) Get(ctx context.Context, key string, result interface{}) bool {
	return false // TODO
}

func (c *LocalCacheStore) GetAll(ctx context.Context, obj interface{}, key ...string) (map[string]any, bool) {
	return nil, false // TODO
}

func (c *LocalCacheStore) Put(ctx context.Context, key string, obj interface{}, expiration time.Duration) bool {
	return false // TODO
}

func (c *LocalCacheStore) PutAll(ctx context.Context, values map[string]any, expiration time.Duration) bool {
	return false // TODO
}

func (c *LocalCacheStore) Delete(ctx context.Context, key ...string) bool {
	return false // TODO
}

func (c *LocalCacheStore) Cleanup() {
	// TODO
}
