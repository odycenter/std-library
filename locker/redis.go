package locker

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// RedisLocker 分布式锁redis实现
type RedisLocker struct {
	cli redis.UniversalClient
}

func newRedis(opt *Option) *RedisLocker {
	if opt.UseCluster {
		return &RedisLocker{redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        opt.Url,
			Username:     opt.Username,
			Password:     opt.Password,
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
			PoolSize:     10,
			TLSConfig:    opt.TLS,
		})}
	}
	return &RedisLocker{redis.NewClient(&redis.Options{
		Addr:         opt.Url[0],
		Username:     opt.Username,
		Password:     opt.Password,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		PoolSize:     10,
		TLSConfig:    opt.TLS,
	})}
}

// Lock 加锁
func (r *RedisLocker) Lock(k string, ex ...time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ex = append(ex, time.Minute*10)
	return r.cli.SetNX(ctx, k, time.Now().Format(time.RFC3339), ex[0]).Val()
}

// Unlock 解锁
func (r *RedisLocker) Unlock(k string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r.cli.Del(ctx, k)
}
