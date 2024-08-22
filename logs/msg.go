package logs

import (
	"fmt"
	"path"
	"time"
)

var AppName = ""

// Msg 日志消息结构体
type Msg struct {
	Level               LogLevel
	Msg                 string
	ID                  string
	When                time.Time
	FilePath            string
	LineNumber          int
	Args                []any
	enableFuncCallDepth bool
}

// Format 格式化传入的msg
func (m *Msg) Format() string {
	msg := m.Msg
	if len(m.Args) > 0 {
		msg = fmt.Sprintf(msg, m.Args...)
	}
	_, f := path.Split(m.FilePath)
	var level = levelPrefix[m.Level]

	if m.ID == "" {
		return fmt.Sprintf(withoutIdPattern, m.When.Format(time.RFC3339Nano), AppName, level, f, m.LineNumber, msg)
	}

	return fmt.Sprintf(pattern, m.When.Format(time.RFC3339Nano), m.ID, AppName, level, f, m.LineNumber, msg)
}
