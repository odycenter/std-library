package redis

import "github.com/redis/go-redis/v9"

// Pipeliner 管道
type Pipeliner struct {
	redis.Pipeliner
}
