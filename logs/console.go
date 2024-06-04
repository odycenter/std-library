package logs

import (
	"github.com/shiena/ansicolor"
	"os"
)

// consoleWriter 实现 LoggerInterface 并将消息写入终端。
type consoleWriter struct {
	lg       *logWriter
	Level    LogLevel
	Colorful bool
}

// NewConsole 创建 LoggerInterface 返回的 ConsoleWriter。
func NewConsole() Logger {
	return newConsole()
}

func newConsole() *consoleWriter {
	cw := &consoleWriter{
		lg:       newLogWriter(ansicolor.NewAnsiColorWriter(os.Stdout)),
		Level:    LevelDebug,
		Colorful: true,
	}
	return cw
}

// Init 初始化 console Logger
func (c *consoleWriter) Init(opt *Option) error {
	if opt == nil {
		return nil
	}
	c.Level = opt.LogLevel
	return nil
}

// WriteMsg 在控制台中写入消息
func (c *consoleWriter) WriteMsg(lm *Msg) error {
	if lm.Level > c.Level {
		return nil
	}
	_, _ = c.lg.writeln(lm.Format())
	return nil
}

// Destroy 实现，未设置
func (c *consoleWriter) Destroy() {
}

func init() {
	Register(AdapterConsole, NewConsole)
}
