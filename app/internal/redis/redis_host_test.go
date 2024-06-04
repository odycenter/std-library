package internal_redis_test

import (
	"github.com/stretchr/testify/assert"
	internal_redis "std-library/app/internal/redis"
	"std-library/app/internal/web"
	"testing"
)

func TestParse(t *testing.T) {
	assert.Equal(t, "redis", internal_redis.Host("redis:6379").String())
	assert.Panics(t, func() {
		internal_redis.Host(":8080")
	})
	assert.Equal(t, "redis:1234", web.Parse("redis:1234").String())
}
