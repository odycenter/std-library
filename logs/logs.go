package logs

import (
	"fmt"
	"strings"
	"sync"
)

type Adapter string
type LogLevel int

type newLoggerFunc func() Logger

// Option log配置结构
type Option struct {
	Adapter  Adapter
	LogLevel LogLevel
	//-----log to kafka options-----
	KafkaBrokersAddr []string //kafka集群节点
	Topic            string   //订阅
	GroupName        string   //组名
}

type loggerFunc func() Logger

// Logger log抽象结构
type Logger interface {
	Init(opt *Option) error
	WriteMsg(msg *Msg) error
	Destroy()
}

type nameLogger struct {
	Logger
	name Adapter
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
