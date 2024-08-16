package logs_test

import (
	"std-library/logs"
	"testing"
)

func TestDefault(t *testing.T) {
	_ = logs.SetLogger(logs.AdapterConsole, &logs.Option{
		Adapter:  logs.AdapterConsole,
		LogLevel: logs.LevelDebug,
	})
	logs.SetLogFuncCallDepth(3)
	logs.Error("error", "Ex", "Ex")
	logs.Debug("debug", "Ex", "Ex")
}
