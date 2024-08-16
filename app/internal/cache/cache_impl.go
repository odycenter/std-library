package internal_cache

import (
	"context"
	"fmt"
	"github.com/mohae/deepcopy"
	"log/slog"
	"reflect"
	"std-library/app/cache"
	actionlog "std-library/app/log"
	reflects "std-library/reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type CacheImpl struct {
	name       string
	typeName   string
	Expiration time.Duration
	Options    cache.Options
	CacheStore CacheStore
}

func (c *CacheImpl) GetTypeName() string {
	return c.typeName
}

func (c *CacheImpl) Name() string {
	return c.name
}

func (c *CacheImpl) TypeName(typeName string) {
	if c.typeName != "" {
		return
	}
	c.typeName = typeName
	c.name = strings.ToLower(typeName)
}

func (c *CacheImpl) Get(ctx context.Context, key string, obj interface{}, f func(key string) (interface{}, error)) (err error) {
	resultVal := reflect.ValueOf(obj)
	if resultVal.Kind() != reflect.Ptr {
		return fmt.Errorf("key: %s, obj argument must be a pointer, but was a %s", key, resultVal.String())
	}

	err = c.checkType(obj)
	if err != nil {
		return err
	}

	cacheKey := c.cacheKey(key)
	success := c.CacheStore.Get(ctx, cacheKey, obj)
	if success {
		c.stat(ctx, "cache_hits", 1)
		return
	}

	var o interface{}
	defer func() {
		if err == nil {
			err = c.copy(o, obj)
		}
	}()

	slog.DebugContext(ctx, fmt.Sprintf("load value, key=%s", key))
	o, err = c.load(key, f)
	if err != nil {
		return
	}
	c.CacheStore.Put(ctx, cacheKey, o, c.Expiration)
	c.stat(ctx, "cache_misses", 1)
	return
}

func (c *CacheImpl) GetByKey(ctx context.Context, key string, obj interface{}) bool {
	cacheKey := c.cacheKey(key)
	return c.CacheStore.Get(ctx, cacheKey, obj)
}

func (c *CacheImpl) GetAll(ctx context.Context, keys []string, obj interface{}, f func(key string) (interface{}, error)) (map[string]interface{}, error) {
	values := make(map[string]interface{}, len(keys))
	newValues := make(map[string]interface{})
	cacheKeys := c.cacheKeys(keys...)
	cacheValues, _ := c.CacheStore.GetAll(ctx, obj, cacheKeys...)
	c.stat(ctx, "cache_hits", float64(len(cacheValues)))
	for i, key := range keys {
		cacheKey := cacheKeys[i]
		cacheValue := cacheValues[cacheKey]
		if cacheValue == nil {
			slog.DebugContext(ctx, fmt.Sprintf("load value, key=%s", key))
			val, err := c.load(key, f)
			if err != nil {
				return nil, err
			}
			cacheValue = val
			newValues[cacheKey] = cacheValue
		}
		if reflect.TypeOf(cacheValue).Kind() == reflect.Ptr {
			cacheValue = reflect.ValueOf(cacheValue).Elem().Interface()
		}
		values[key] = cacheValue
	}
	if len(newValues) > 0 {
		c.CacheStore.PutAll(ctx, newValues, c.Expiration)
		c.stat(ctx, "cache_misses", float64(len(newValues)))
	}
	return values, nil
}

func (c *CacheImpl) load(key string, f func(key string) (interface{}, error)) (interface{}, error) {
	val, err := f(key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, fmt.Errorf("value must not be null, key=%s", key)
	}

	err = c.checkLoaderReturnType(val)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (c *CacheImpl) Put(ctx context.Context, key string, obj interface{}) bool {
	err := c.checkType(obj)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("key: %s, error: %s", key, err.Error()))
		return false
	}

	cacheKey := c.cacheKey(key)

	return c.CacheStore.Put(ctx, cacheKey, obj, c.Expiration)
}

func (c *CacheImpl) PutAll(ctx context.Context, values map[string]any) bool {
	cacheValues := make(map[string]interface{}, len(values))
	for key, value := range values {
		err := c.checkType(value)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("key: %s, error: %s", key, err.Error()))
			return false
		}
		cacheValues[c.cacheKey(key)] = value
	}
	return c.CacheStore.PutAll(ctx, cacheValues, c.Expiration)
}

func (c *CacheImpl) Evict(ctx context.Context, key string) bool {
	cacheKey := c.cacheKey(key)
	return c.CacheStore.Delete(ctx, cacheKey)
}
func (c *CacheImpl) EvictAll(ctx context.Context, key ...string) bool {
	cacheKeys := c.cacheKeys(key...)
	return c.CacheStore.Delete(ctx, cacheKeys...)
}

func (c *CacheImpl) checkType(obj interface{}) error {
	if c.typeName != reflects.StructFullName(obj) {
		return fmt.Errorf("illegal usage, cache type is not correct, must be %s, please check", c.typeName)
	}
	return nil
}

func (c *CacheImpl) checkLoaderReturnType(obj interface{}) error {
	if c.typeName != reflects.StructFullName(obj) {
		return fmt.Errorf("illegal usage, cache loader return type is not correct, must be %s, please check", c.typeName)
	}
	return nil
}

// copy object to return, to avoid dirty data
func (c *CacheImpl) copy(src, dst interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = errors.WithStack(v)
			default:
				err = errors.New(fmt.Sprint(r))
			}
			c.Options.OnError(err)
		}
	}()

	v := deepcopy.Copy(src)
	if reflect.ValueOf(v).IsValid() {
		reflect.ValueOf(dst).Elem().Set(reflect.Indirect(reflect.ValueOf(v)))
	}
	return
}

func (c *CacheImpl) cacheKey(key string) string {
	return c.name + ":" + key
}

func (c *CacheImpl) cacheKeys(key ...string) []string {
	var keys = make([]string, len(key))
	for i := range key {
		keys[i] = c.cacheKey(key[i])
	}
	return keys
}

func (c *CacheImpl) stat(ctx context.Context, key string, value float64) {
	val := actionlog.GetStat(&ctx, key)
	if val != 0 {
		value += val
	}
	actionlog.Stat(&ctx, key, value)
}
