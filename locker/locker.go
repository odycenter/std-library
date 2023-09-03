package locker

import (
	"crypto/tls"
	"log"
	"time"
)

var globLocker Locker

// Locker 分布式锁
type Locker interface {
	Lock(string, ...time.Duration) bool
	Unlock(string)
}

type locker int

const (
	Redis = locker(iota + 1) //redis实现
	Etcd                     //etcd实现 TODO 涉及grpc版本问题，暂不支持
)

type Option struct {
	Url        []string //连接分布式锁实现server的地址
	Username   string
	Password   string
	UseCluster bool   //使用集群
	Locker     locker //锁实现类型
	TLS        *tls.Config
}

// Init 初始化
func Init(opt *Option) Locker {
	if opt == nil {
		log.Panicln("locker init need option")
	}
	switch opt.Locker {
	case Redis:
		globLocker = newRedis(opt)
	case Etcd:
		//globLocker = newEtcd(opt)
	}
	return nil
}

func Lock(k string, ex ...time.Duration) bool {
	if globLocker == nil {
		return false
	}
	return globLocker.Lock(k, ex...)
}

func Unlock(k string) {
	if globLocker == nil {
		return
	}
	globLocker.Unlock(k)
}
