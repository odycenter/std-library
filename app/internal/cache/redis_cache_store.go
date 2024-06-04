package internal_cache

import (
	"context"
	"encoding/json"
	"errors"
	redisV9 "github.com/redis/go-redis/v9"
	"reflect"
	"std-library/app/redis"
	"std-library/logs"
	"time"
)

type RedisCacheStore struct {
	redis redis.Redis
}

func (c *RedisCacheStore) Initialize(redis redis.Redis) {
	if c.redis != nil {
		logs.Error("redisImpl is already configured, please configure only once")
		return
	}
	c.redis = redis
}

func (c *RedisCacheStore) Get(ctx context.Context, key string, result interface{}) bool {
	body, err := c.redis.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redisV9.Nil) {
			logs.DebugWithCtx(ctx, "[Cache][Get] Key not found:<%s>", key)
		} else {
			logs.ErrorWithCtx(ctx, "[Cache][Get] key:<%s> failed:%v", key, err)
		}
		return false
	}

	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		logs.ErrorWithCtx(ctx, "[Cache][Get] Json Unmarshal Error, key:%s, value:%s, failed:%v", key, body, err)
		return false
	}

	return true
}

func (c *RedisCacheStore) GetAll(ctx context.Context, obj interface{}, key ...string) (map[string]any, bool) {
	data, err := c.redis.MultiGet(ctx, key...)
	if err != nil {
		logs.ErrorWithCtx(ctx, "[Cache][GetAll] keys:<%v> failed:%v", key, err)
		return nil, false
	}

	result := make(map[string]any, len(data))
	for k, v := range data {
		value := reflect.New(reflect.TypeOf(obj)).Interface()
		err = json.Unmarshal([]byte(v), &value)
		if err != nil {
			logs.ErrorWithCtx(ctx, "[Cache][GetAll] Json Unmarshal Error, key:%s, value:%s, failed:%v", k, v, err)
			return nil, false
		}
		result[k] = value
	}

	return result, true
}

func (c *RedisCacheStore) Put(ctx context.Context, key string, obj interface{}, expiration time.Duration) bool {
	bs, err := json.Marshal(obj)
	if err != nil {
		logs.ErrorWithCtx(ctx, "[Cache][Put] Json Marshal Error, key: %s, value:%s, failed:%v", key, obj, err)
		return false
	}

	_, err = c.redis.Set(ctx, key, string(bs), expiration)
	if err != nil {
		logs.ErrorWithCtx(ctx, "[Cache][Put] key: %s, failed:%v", key, err)
	}
	return err == nil
}

func (c *RedisCacheStore) PutAll(ctx context.Context, values map[string]any, expiration time.Duration) bool {
	data := make(map[string]any, len(values))
	for k, v := range values {
		bs, err := json.Marshal(v)
		if err != nil {
			logs.ErrorWithCtx(ctx, "[Cache][PutAll] Json Marshal Error, key: %s, value:%s, failed:%v", k, v, err)
			return false
		}
		data[k] = string(bs)

	}
	err := c.redis.MultiSet(ctx, data, expiration)
	if err != nil {
		logs.ErrorWithCtx(ctx, "[Cache][PutAll] failed:%v", err)
	}
	return err == nil
}

func (c *RedisCacheStore) Delete(ctx context.Context, key ...string) bool {
	_, err := c.redis.Del(ctx, key...)
	if err != nil {
		logs.ErrorWithCtx(ctx, "[Cache][Delete] keys: %v, failed:%v", key, err)
	}
	return err == nil
}
