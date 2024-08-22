package logs

import (
	"log"
	"sync"
)

// 适配器类型，输出方式
const (
	AdapterConsole = Adapter("console")
)

// 日志等级 RFC5424 标准
const (
	LevelError = LogLevel(iota)
	LevelWarning
	LevelNotice
	LevelInformation
	LevelDebug
)

const defaultAsyncMsgLen = 1e3

var defaultLogger = NewLogger()

var defaultLoggerMap = struct {
	sync.RWMutex
	logs map[string]*log.Logger
}{
	logs: map[string]*log.Logger{},
}

var pattern = `{"@timestamp":"%s", "id":%q, "app":"%s", "level":"%s", "file":"%s:%d", "message":%q}`
var withoutIdPattern = `{"@timestamp":"%s", "app":"%s", "level":"%s", "file":"%s:%d", "message":%q}`
var levelPrefix = [LevelDebug + 1]string{"ERROR", "WARN", "WARN", "INFO", "DEBUG"}
