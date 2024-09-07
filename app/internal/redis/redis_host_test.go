package internal_redis_test

import (
	internal_redis "github.com/odycenter/std-library/app/internal/redis"
	"github.com/odycenter/std-library/app/web"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	assert.Equal(t, "redis", internal_redis.Host("redis:6379").String())
	assert.Panics(t, func() {
		internal_redis.Host(":8080")
	})
	assert.Equal(t, "redis:1234", web.Parse("redis:1234").String())
}
