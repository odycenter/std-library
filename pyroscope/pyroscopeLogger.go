package pyroscope

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	LevelError = iota + 1
	LevelInfo
	LevelDebug
)

var mLogLevel = [LevelDebug + 1]string{"", "E", "I", "D"}

type pyroscopeLoggerImpl struct {
	Level      int
	InnerLevel int
}

func newLogger(level int) *pyroscopeLoggerImpl {
	return &pyroscopeLoggerImpl{Level: level}
}

func (logger *pyroscopeLoggerImpl) Infof(a string, b ...any) {
	if LevelInfo > logger.Level {
		return
	}
	logger.InnerLevel = LevelInfo
	_, _ = fmt.Fprintf(os.Stdout, logger.formatLog(a, b...), b...)
}
func (logger *pyroscopeLoggerImpl) Debugf(a string, b ...any) {
	if LevelDebug > logger.Level {
		return
	}
	logger.InnerLevel = LevelDebug
	_, _ = fmt.Fprintf(os.Stdout, logger.formatLog(a, b...), b...)
}
func (logger *pyroscopeLoggerImpl) Errorf(a string, b ...any) {
	if LevelError > logger.Level {
		return
	}
	logger.InnerLevel = LevelError
	_, _ = fmt.Fprintf(os.Stdout, logger.formatLog(a, b...), b...)
}

var sb strings.Builder

func (logger *pyroscopeLoggerImpl) formatLog(f string, v ...any) string {
	sb.Reset()
	sb.WriteString(time.Now().Format("2006/01/02 15:04:05.000"))
	sb.WriteString(" [")
	sb.WriteString(mLogLevel[logger.InnerLevel])
	sb.WriteString("] ")
	sb.WriteString("[")
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	sb.WriteString(filename)
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(line))
	sb.WriteString("] ")
	sb.WriteString(formatLog(f, v...))
	sb.WriteString("\n")
	return sb.String()
}

func formatLog(f any, v ...any) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
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
