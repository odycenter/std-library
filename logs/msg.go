package logs

import (
	"fmt"
	"path"
	"time"
)

// Msg 日志消息结构体
type Msg struct {
	Level               LogLevel
	Msg                 string
	ID                  string
	When                time.Time
	FilePath            string
	LineNumber          int
	Title               string
	ExecDur             int64
	Args                []any
	enableFullFilePath  bool
	enableFuncCallDepth bool
}

// Format 格式化传入的msg
func (m *Msg) Format() string {
	msg := m.Msg
	if len(m.Args) > 0 {
		msg = fmt.Sprintf(msg, m.Args...)
	}
	_, f := path.Split(m.FilePath)
	h, _, _ := formatTimeHeader(m.When)
	msg = fmt.Sprintf(pattern, m.ID, h, levelPrefix[m.Level], m.Title, m.ExecDur, f, m.LineNumber, msg)
	return msg
}
