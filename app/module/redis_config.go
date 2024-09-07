package module

import (
	"context"
	"fmt"
	internal "github.com/odycenter/std-library/app/internal/module"
	internalredis "github.com/odycenter/std-library/app/internal/redis"
	"log"
	"log/slog"
	"time"
)

type RedisConfig struct {
	name          string
	moduleContext *Context
	redis         *internalredis.RedisImpl
	host          string
	db            int
	password      string
}

func (c *RedisConfig) Initialize(moduleContext *Context, name string) {
	c.name = name
	c.moduleContext = moduleContext
	c.redis = c.createRedis()
}

func (c *RedisConfig) createRedis() *internalredis.RedisImpl {
	slog.Info(fmt.Sprintf("create redis client, name:%s", c.name))
	redisImpl := internalredis.New(c.name)
	redisImpl.Timeout(3 * time.Second)
	c.moduleContext.StartupHook.Initialize = append(c.moduleContext.StartupHook.Initialize, redisImpl)
	c.moduleContext.ShutdownHook.Add(internal.STAGE_6, func(ctx context.Context, timeoutInMs int64) {
		redisImpl.Close()
	})

	return redisImpl
}

func (c *RedisConfig) ForceEarlyStart() {
	c.Validate()
	hostname := internal.Hostname(c.host)
	internal.ResolveHost(context.Background(), hostname)
	c.redis.Execute(context.Background())
}

func (c *RedisConfig) Validate() {
	if c.host == "" {
		log.Fatalf("redis host must be configured, name=" + c.name)
	}
}

func (c *RedisConfig) Host(host string) {
	if c.host != "" {
		log.Fatalf("redis host is already configured, host=%s, previous=%s", host, c.host)
	}
	c.host = host
	c.redis.Host(host)
	c.moduleContext.Probe.AddHostURI(host)
}

func (c *RedisConfig) Password(password string) {
	if c.redis.Initialized() {
		log.Fatalf("redis is already initialized, can not set password! name=" + c.name)
	}
	if c.password != "" {
		log.Fatalf("redis password is already configured!")
	}
	c.password = password
	c.redis.Password(password)
}

func (c *RedisConfig) DB(db int) {
	if c.redis.Initialized() {
		log.Fatalf("redis is already initialized, can not set db! name=" + c.name)
	}
	if c.db != 0 {
		log.Fatalf("redis db is already configured! db=%d, previous=%d", db, c.db)
	}
	c.db = db
	c.redis.DB(db)
}

func (c *RedisConfig) PoolSize(minSize, maxSize int) {
	if c.redis.Initialized() {
		log.Fatalf("redis is already initialized, can not set poolSize! name=" + c.name)
	}
	c.redis.PoolSize(minSize, maxSize)
}
