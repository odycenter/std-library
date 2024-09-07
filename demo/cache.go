package demo

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"github.com/odycenter/std-library/app/cache"
	"github.com/odycenter/std-library/app/module"
	"log/slog"
	"net/http"
	"time"
)

type CacheTest struct {
	PlayerName string
}

type CacheModule struct {
	module.Common
}

func (m *CacheModule) Initialize() {
	c := m.Cache().Add(CacheTest{}, time.Second*300)
	service := &cacheService{
		cache: c,
	}
	handler := &getAllHandler{
		service: service,
	}
	web.Handler("/cache-getall", handler)
	web.Handler("/cache-get", handler)
	web.Handler("/cache-put", handler)
	web.Handler("/cache-put-error", handler)
	web.Handler("/cache-evict", handler)
	web.Handler("/cache-evict-all", handler)
}

type getAllHandler struct {
	service *cacheService
}

func (h *getAllHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("r: " + r.Method + ":" + r.URL.Path)
	if r.URL.Path == "/cache-getall" {
		_ = h.service.GetAll(r.Context(), []string{"key1", "test2"})
	} else if r.URL.Path == "/cache-get" {
		_ = h.service.Get(r.Context(), "key1")
	} else if r.URL.Path == "/cache-put" {
		_ = h.service.Put(r.Context(), "key1", CacheTest{PlayerName: "put test..."})
	} else if r.URL.Path == "/cache-put-error" {
		_ = h.service.Put(r.Context(), "key1", "")
	} else if r.URL.Path == "/cache-evict" {
		h.service.Evict(r.Context(), "key1")
	} else if r.URL.Path == "/cache-evict-all" {
		h.service.EvictAll(r.Context(), "key1", "test2")
	}
}

type cacheService struct {
	cache cache.Cache
}

func (c *cacheService) Get(ctx context.Context, key string) error {
	var tt CacheTest
	err := c.cache.Get(ctx, key, &tt, func(key string) (interface{}, error) {
		return CacheTest{PlayerName: "from get loader:" + key}, nil
	})
	return err
}

func (c *cacheService) GetAll(ctx context.Context, keys []string) error {
	result, err := c.cache.GetAll(ctx, keys, CacheTest{}, func(key string) (interface{}, error) {
		return CacheTest{PlayerName: "from get all loader:" + key}, nil
	})
	if err != nil {
		panic(err)
	}
	for k, v := range result {
		if t, ok := (v).(CacheTest); ok {
			fmt.Printf("key: %s, value:%v \r\n", k, t.PlayerName)
		}
	}
	return nil
}

func (c *cacheService) Put(ctx context.Context, key string, value interface{}) bool {
	return c.cache.Put(ctx, key, value)
}

func (c *cacheService) PutAll(ctx context.Context, values map[string]interface{}) bool {
	return c.cache.PutAll(ctx, values)
}

func (c *cacheService) Evict(ctx context.Context, key string) bool {
	return c.cache.Evict(ctx, key)
}

func (c *cacheService) EvictAll(ctx context.Context, key ...string) bool {
	return c.cache.EvictAll(ctx, key...)
}
