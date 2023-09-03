package logs

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// redisLogWriter implements LoggerInterface.
// 将日志写入Redis list
type redisLogWriter struct {
	client            redis.UniversalClient
	input             chan string
	IsCluster         bool
	TLS               *tls.Config   //是否使用TLS链接
	RedisHost         []string      //redis地址:"127.0.0.1:679"
	RedisUsername     string        //用户名redis>6.0
	RedisPassword     string        //密码
	RedisKey          string        //推入redis的key(推入key类型仅list)
	RedisMaxConn      int           //最大链接数
	RedisMinConn      int           //最小链接数
	RedisWriteTimeout time.Duration //写入超时
	logFormatter      LogFormatter
	Formatter         string
	Level             LogLevel
}

func (w *redisLogWriter) getRedisMaxConn() int {
	if w.RedisMaxConn == 0 {
		return 100
	}
	return w.RedisMaxConn
}

func (w *redisLogWriter) getRedisMinConn() int {
	if w.RedisMinConn == 0 {
		return 10
	}
	return w.RedisMinConn
}

func (w *redisLogWriter) getRedisWriteTimeout() time.Duration {
	if w.RedisWriteTimeout == 0 {
		return time.Second
	}
	return w.RedisWriteTimeout
}

// newRedisLogWriter 创建一个作为 LoggerInterface 返回的 redisLogWriter
func newRedisLogWriter() Logger {
	w := &redisLogWriter{
		RedisMaxConn:      20,
		RedisMinConn:      2,
		RedisWriteTimeout: time.Second, //对于log来说1s足矣
		Level:             LevelDebug,
	}
	w.logFormatter = w
	return w
}

func (*redisLogWriter) Format(lm *Msg) string {
	return lm.Format()
}

func (w *redisLogWriter) SetFormatter(f LogFormatter) {
	w.logFormatter = f
}

func (w *redisLogWriter) Init(opt *Option) error {
	if opt == nil {
		return nil
	}
	w.input = make(chan string, 100000) //default 最大10wQuery
	w.RedisHost = opt.RedisHost
	w.IsCluster = opt.IsCluster
	w.TLS = opt.TLS
	w.RedisUsername = opt.RedisUsername
	w.RedisPassword = opt.RedisPassword
	w.RedisKey = opt.RedisKey
	w.RedisMaxConn = opt.RedisMinConn
	w.RedisMinConn = opt.RedisMaxConn
	w.RedisWriteTimeout = opt.RedisWriteTimeout
	w.Formatter = opt.Formatter
	w.Level = opt.LogLevel
	if len(w.Formatter) > 0 {
		fmtr, ok := GetFormatter(w.Formatter)
		if !ok {
			return fmt.Errorf("the formatter with name: %s not found", w.Formatter)
		}
		w.logFormatter = fmtr
	}
	err := w.startLogger()
	return err
}

func (w *redisLogWriter) startLogger() error {
	if w.IsCluster {
		w.client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        w.RedisHost,
			Username:     w.RedisUsername,
			Password:     w.RedisPassword,
			MaxRetries:   1,
			WriteTimeout: w.getRedisWriteTimeout(),
			PoolSize:     w.getRedisMaxConn(),
			MinIdleConns: w.getRedisMinConn(),
			TLSConfig:    w.TLS,
		})
		go w.push()
		return nil
	}
	w.client = redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         w.RedisHost[0],
		Username:     w.RedisUsername,
		Password:     w.RedisPassword,
		DB:           0,
		MaxRetries:   1,
		WriteTimeout: w.getRedisWriteTimeout(),
		PoolSize:     w.getRedisMaxConn(),
		MinIdleConns: w.getRedisMinConn(),
		TLSConfig:    w.TLS,
	})
	go w.push()
	return nil
}

// WriteMsg 此处使用 LPUSH 写入数据，如需顺序获取可使用 RPOP
func (w *redisLogWriter) WriteMsg(lm *Msg) error {
	if lm.Level > w.Level {
		return nil
	}

	msg := w.logFormatter.Format(lm)
	w.input <- msg
	return nil
}

func (w *redisLogWriter) push() {
	for {
		select {
		case msg := <-w.input:
			_, _ = w.client.LPush(context.Background(), w.RedisKey, msg).Result()
		}
	}
}

func (w *redisLogWriter) Destroy() {
	_ = w.client.Close()
	//等待日志消耗完毕
	for range time.Tick(time.Second) {
		if len(w.input) == 0 {
			close(w.input)
			break
		}
	}
}

func (w *redisLogWriter) Flush() {

}

func init() {
	Register(AdapterRedis, newRedisLogWriter)
}
