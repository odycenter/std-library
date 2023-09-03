package logs

import (
	"errors"
	"fmt"
	"github.com/shiena/ansicolor"
	"os"
)

// consoleWriter 实现 LoggerInterface 并将消息写入终端。
type consoleWriter struct {
	lg        *logWriter
	formatter LogFormatter
	Formatter string
	Level     LogLevel
}

func (c *consoleWriter) Format(lm *Msg) string {

	return lm.Format()
}

func (c *consoleWriter) SetFormatter(f LogFormatter) {
	c.formatter = f
}

// NewConsole 创建 LoggerInterface 返回的 ConsoleWriter。
func NewConsole() Logger {
	return newConsole()
}

func newConsole() *consoleWriter {
	cw := &consoleWriter{
		lg:    newLogWriter(ansicolor.NewAnsiColorWriter(os.Stdout)),
		Level: LevelDebug,
	}
	cw.formatter = cw
	return cw
}

// Init 初始化 console Logger
func (c *consoleWriter) Init(opt *Option) error {
	if opt == nil {
		return nil
	}
	c.Formatter = opt.Formatter
	if len(c.Formatter) > 0 {
		formatter, ok := GetFormatter(c.Formatter)
		if !ok {
			return errors.New(fmt.Sprintf("the formatter with name: %s not found", c.Formatter))
		}
		c.formatter = formatter
	}
	c.Level = opt.LogLevel
	return nil
}

// WriteMsg 在控制台中写入消息
func (c *consoleWriter) WriteMsg(lm *Msg) error {
	if lm.Level > c.Level {
		return nil
	}
	msg := c.formatter.Format(lm)
	_, _ = c.lg.writeln(msg)
	return nil
}

// Destroy 实现，未设置
func (c *consoleWriter) Destroy() {
}

// Flush 实现，未设置
func (c *consoleWriter) Flush() {
}

func init() {
	Register(AdapterConsole, NewConsole)
}
