package module

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"log"
	"log/slog"
	"reflect"
	"std-library/app/cache"
	internalcache "std-library/app/internal/cache"
	internal "std-library/app/internal/module"
	internalredis "std-library/app/internal/redis"
	"std-library/app/internal/web/sys"
	reflects "std-library/reflect"
	"strings"
	"sync"
	"time"
)

type CacheConfig struct {
	name            string
	moduleContext   *Context
	redisCacheStore *internalcache.RedisCacheStore
	localCacheStore *internalcache.LocalCacheStore
	options         *cache.Options
	maxLocalSize    int
	caches          map[string]*internalcache.CacheImpl
	mu              sync.Mutex
}

func (c *CacheConfig) Initialize(moduleContext *Context, name string) {
	c.name = name
	c.moduleContext = moduleContext
	c.caches = make(map[string]*internalcache.CacheImpl)
	controller := internal_sys.NewCacheController(c.caches)
	web.Handler("/_sys/cache", controller)
	web.Handler("/_sys/cache/*/*", controller)
}

func (c *CacheConfig) Validate() {
	if len(c.caches) == 0 {
		log.Fatal("cache is configured but no cache added, please remove unnecessary config")
	}

	// maxLocalSize() can be configured before localCacheStore is created, so set max size at end
	if c.maxLocalSize > 0 && c.localCacheStore != nil {
		c.localCacheStore.MaxSize = c.maxLocalSize
	}
}

func (c *CacheConfig) Local() {
	if c.localCacheStore != nil || c.redisCacheStore != nil {
		log.Fatal("cache store is already configured, please configure only once")
	}
	c.configureLocalCacheStore()
}

func (c *CacheConfig) Redis(host string, password ...string) {
	if c.localCacheStore != nil || c.redisCacheStore != nil {
		log.Fatal("cache store is already configured, please configure only once")
	}

	c.configureRedis(host, password...)
}

func (c *CacheConfig) Options(options *cache.Options) {
	if c.localCacheStore != nil || c.redisCacheStore != nil {
		log.Fatalf("cache is already initialized, can not set options! name=" + c.name)
	}
	c.options = options
}

func (c *CacheConfig) Add(obj interface{}, expiration time.Duration) cache.Cache {
	if c.localCacheStore == nil && c.redisCacheStore == nil {
		log.Fatal("cache store is not configured, please configure first")
	}
	// validate obj class, only struct is allowed
	if reflect.TypeOf(obj).Kind() != reflect.Struct {
		log.Fatalf("%v is not supported, illegal Add, only struct is allowed", obj)
	}

	typeName := reflects.StructFullName(obj)
	name := strings.ToLower(typeName)
	slog.Info(fmt.Sprintf("add cache, struct=%s, expiration=%v", typeName, expiration))
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.caches[name]; ok {
		log.Fatal("found duplicate cache name, name=" + name)
	}

	var cacheImpl = internalcache.CacheImpl{
		Expiration: expiration}
	cacheImpl.TypeName(typeName)
	if c.redisCacheStore != nil {
		cacheImpl.CacheStore = c.redisCacheStore
	} else {
		cacheImpl.CacheStore = c.localCacheStore
	}
	if c.options != nil {
		cacheImpl.Options = *c.options
	}

	c.caches[name] = &cacheImpl
	return &cacheImpl

}

func (c *CacheConfig) MaxLocalSize(size int) {
	c.maxLocalSize = size
}

func (c *CacheConfig) configureRedis(host string, password ...string) {
	slog.Info(fmt.Sprintf("create redis cache store, host=%v", host))
	redisImpl := internalredis.New("redis-cache")
	hostname := internal.Hostname(host)
	internal.ResolveHost(context.Background(), hostname)
	redisImpl.Host(host)
	if password != nil && len(password) > 0 && password[0] != "" {
		redisImpl.Password(password[0])
	}
	redisImpl.Timeout(1 * time.Second) // for cache, use shorter timeout than default redis config
	c.moduleContext.ShutdownHook.Add(internal.STAGE_6, func(ctx context.Context, timeoutInMs int64) {
		redisImpl.Close()
	})
	if c.options != nil {
		redisImpl.PoolSize(c.options.MinPoolSize, c.options.MaxPoolSize)
	}

	redisImpl.Initialize()
	c.redisCacheStore = &internalcache.RedisCacheStore{}
	c.redisCacheStore.Initialize(redisImpl)
}

func (c *CacheConfig) configureLocalCacheStore() {
	if c.localCacheStore == nil {
		slog.Info("create local cache store")
		c.localCacheStore = &internalcache.LocalCacheStore{}
		// TODO cleanup local cache store
		// c.localCacheStore.Cleanup()
	}
}
