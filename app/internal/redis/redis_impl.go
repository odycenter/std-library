package internal_redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	app "std-library/app/conf"
	"std-library/logs"
	redismigration "std-library/redis"
	"time"
)

type RedisImpl struct {
	name        string
	host        *RedisHost
	password    string
	db          int
	timeout     time.Duration
	minSize     int
	maxSize     int
	client      redis.UniversalClient
	initialized bool
}

func New(name string) *RedisImpl {
	return &RedisImpl{
		name: name,
	}
}

func (r *RedisImpl) Execute(_ context.Context) {
	if r.initialized {
		return
	}

	logs.Debug("redis Initialize, name=%s", r.name)
	r.Initialize()

	if r.name == "redis" {
		redismigration.InitMigration("default", r.client)
	} else {
		redismigration.InitMigration(r.name, r.client)
	}

	r.initialized = true
}

func (r *RedisImpl) Initialized() bool {
	return r.initialized
}

func (r *RedisImpl) Initialize() {
	if r.minSize <= 0 {
		r.minSize = 5
	}
	if r.maxSize <= 0 {
		r.maxSize = 50
	}
	client := redis.NewClient(&redis.Options{
		ClientName:      app.Name + ":" + r.name,
		Addr:            r.host.String(),
		Password:        r.password,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     r.timeout,
		WriteTimeout:    r.timeout,
		ConnMaxIdleTime: 30 * time.Minute,
		ConnMaxLifetime: 120 * time.Minute,
		MinIdleConns:    r.minSize,
		PoolSize:        r.maxSize,
		DB:              r.db,
	})
	r.client = client
}

func (r *RedisImpl) Host(host string) {
	r.host = Host(host)
}

func (r *RedisImpl) Password(password string) {
	r.password = password
}

func (r *RedisImpl) Timeout(timeout time.Duration) {
	r.timeout = timeout
}

func (r *RedisImpl) DB(db int) {
	r.db = db
}

func (r *RedisImpl) PoolSize(minSize, maxSize int) {
	if r.Initialized() {
		log.Fatalf("redis is already initialized, can not set pool size! name=" + r.name)
	}
	r.minSize = minSize
	r.maxSize = maxSize
}

func (r *RedisImpl) Client() redis.UniversalClient {
	if !r.initialized {
		log.Fatalf("redis must be initialized, name=" + r.name)
		return nil
	}
	return r.client
}

func (r *RedisImpl) Close() {
	logs.Info("close redis client, name=%s, host=%s", r.name, r.host)
	err := r.client.Close()
	if err != nil {
		logs.Error("close redis client error, name=%s, host=%s, err=%v", r.name, r.host, err)
	}
}

func (r *RedisImpl) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisImpl) MultiGet(ctx context.Context, key ...string) (map[string]string, error) {
	if key == nil || len(key) == 0 {
		return nil, errors.New("key must not be empty")
	}
	values, err := r.client.MGet(ctx, key...).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for i, v := range values {
		if v == nil {
			continue
		}
		str, ok := v.(string)
		if ok {
			result[key[i]] = str
		}
	}
	return result, nil
}

func (r *RedisImpl) Set(ctx context.Context, key string, value string, expiration time.Duration) (string, error) {
	arg := redis.SetArgs{
		TTL: expiration,
	}
	return r.client.SetArgs(ctx, key, value, arg).Result()
}

func (r *RedisImpl) MultiSet(ctx context.Context, values map[string]interface{}, expiration time.Duration) error {
	arg := redis.SetArgs{
		TTL: expiration,
	}
	for key, value := range values {
		_, err := r.client.SetArgs(ctx, key, value, arg).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RedisImpl) Del(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Del(ctx, keys...).Result()
}
