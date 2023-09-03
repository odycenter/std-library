package logs

import (
	"crypto/tls"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Adapter string
type LogLevel int

type newLoggerFunc func() Logger

// Option log配置结构
type Option struct {
	Adapter   Adapter
	LogLevel  LogLevel
	Formatter string
	//-----log to file options-----
	Filename string //文件名:"logs/beego.log",
	//轮转方式
	Rotate bool //开启轮转:true,
	Daily  bool //按日轮转:true,
	Hourly bool //按小时轮转:true,
	//按行轮转
	MaxLines int //:10000,
	MaxFiles int //:10,
	//按大小轮转
	MaxSize int //:1 << 28,
	//按日轮转
	MaxDays int64 //:15,
	//按小时轮转
	MaxHours int64 //:24
	//日志文件权限
	Perm string //:"0600"
	//如果在 FileName 中指定了目录的权限
	DirPerm    string //:"0770"
	RotatePerm string //:"0440"
	//-----log to kafka options-----
	KafkaBrokersAddr []string //kafka集群节点
	Topic            string   //订阅
	GroupName        string   //组名
	//-----log to redis options-----
	RedisHost         []string      //redis地址:"127.0.0.1:679"
	RedisUsername     string        //用户名redis>6.0
	RedisPassword     string        //密码
	IsCluster         bool          //是否为集群
	TLS               *tls.Config   //是否使用TLS链接
	RedisKey          string        //推入redis的key(推入key类型仅list)
	RedisMaxConn      int           //最大链接数
	RedisMinConn      int           //最小链接数
	RedisWriteTimeout time.Duration //写入超时
}

func (o *Option) getPerm() string {
	if o.Perm == "" {
		return "0660"
	}
	return o.Perm
}

type loggerFunc func() Logger

// Logger log抽象结构
type Logger interface {
	Init(opt *Option) error
	WriteMsg(msg *Msg) error
	Destroy()
	Flush()
	SetFormatter(f LogFormatter)
}

type nameLogger struct {
	Logger
	name Adapter
}

type msg struct {
	level int
	msg   string
	when  time.Time
}

var (
	adapters sync.Map //map[string]Adapter
	msgPool  *sync.Pool
)

// Register 通过提供的名称提供日志。
// 如果 Register 以相同的 name 被调用两次，或 driver 为 nil，它会panic。
func Register(name Adapter, log newLoggerFunc) {
	if log == nil {
		panic("logs: Register provide is nil")
	}
	if _, dup := adapters.Load(name); dup {
		panic("logs: Register called twice for provider <" + name + ">")
	}
	adapters.Store(name, log)
}

func formatPattern(f any, v ...any) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if !strings.Contains(msg, "%") {
			// do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return msg
}
