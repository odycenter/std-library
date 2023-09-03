package logs_test

import (
	"testing"

	"github.com/odycenter/std-library/logs"
)

func TestConsole(t *testing.T) {
	l := logs.NewLogger(10000)
	l.SetLogFuncCallDepth(2)
	_ = l.SetLogger(logs.AdapterConsole, &logs.Option{
		Adapter:   logs.AdapterConsole,
		LogLevel:  logs.LevelDebug,
		Formatter: "",
	})
	testConsoleCalls(l)
}

func testConsoleCalls(dl *logs.DefaultLog) {
	dl.Emergency("emergency")
	dl.Alert("alert")
	dl.Critical("critical")
	dl.Error("error")
	dl.Notice("notice")
	dl.Debug("debug")
}
