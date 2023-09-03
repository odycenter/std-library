package redis

import (
	"github.com/redis/go-redis/v9"
)

// Tx redis事务结构
// 注意：redis任何多命令执行都不支持回滚，redis事务并不是原子操作
type Tx struct {
	*redis.Tx
}
