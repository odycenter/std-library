package logs_test

import (
	"context"
	"std-library/app/log/consts/logKey"
	"std-library/logs"
	"testing"
)

func TestConsole(t *testing.T) {
	l := logs.NewLogger(10000)
	l.SetLogFuncCallDepth(2)
	_ = l.SetLogger(logs.AdapterConsole, &logs.Option{
		Adapter:  logs.AdapterConsole,
		LogLevel: logs.LevelDebug,
	})
	testConsoleCalls(l)
}

func testConsoleCalls(dl *logs.DefaultLog) {
	dl.Alert("alert")
	dl.Critical("critical")
	dl.Error("error")
	dl.Notice("notice")
	dl.Debug("debug")
	ctx := context.WithValue(context.Background(), logKey.Id, "test-id")
	logs.DebugWithCtx(ctx, "debug, log, name: %v", "test123")
}
