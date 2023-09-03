package logs

import (
	"log"
	"sync"
)

// 适配器类型，输出方式
const (
	AdapterConsole = Adapter("console")
	AdapterFile    = Adapter("file")
	AdapterRedis   = Adapter("redis")
	AdapterEs      = Adapter("elastic")
	AdapterKafka   = Adapter("kafka")
)

// 日志等级 RFC5424 标准
const (
	LevelEmergency = LogLevel(iota)
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformation
	LevelDebug
)

// levelLogLogger 被定义来实现 log.Logger 真正的日志级别将是 LevelEmergency
const levelLoggerImpl = -1

const defaultAsyncMsgLen = 1e3

var levelNames = [...]string{"emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"}

var defaultLogger = NewLogger()

var defaultLoggerMap = struct {
	sync.RWMutex
	logs map[string]*log.Logger
}{
	logs: map[string]*log.Logger{},
}

var pattern = `{"ID":%q,Date":"%s","Level":"%s","Title":"%s","ExecDur":%d,"File":"%s:%d","Msg":%q}`
var levelPrefix = [LevelDebug + 1]string{"[M]", "[A]", "[C]", "[E]", "[W]", "[N]", "[I]", "[D]"}
var levelColorPrefix = [LevelDebug + 1]string{"[\x1b[37m[M]\x1b[0m", "[\x1b[36m[A]\x1b[0m", "[\x1b[31m[E]\x1b[0m", "[\x1b[33m[W]\x1b[0m", "[\x1b[32m[N]\x1b[0m", "[\x1b[34m[I]\x1b[0m", "[\x1b[44m[D]\x1b[0m"}
