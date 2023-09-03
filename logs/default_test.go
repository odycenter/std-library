package logs_test

import (
	"testing"

	"github.com/odycenter/std-library/logs"
)

func TestDefault(t *testing.T) {
	_ = logs.SetLogger(logs.AdapterConsole, &logs.Option{
		Adapter:  logs.AdapterConsole,
		LogLevel: logs.LevelDebug,
	})
	logs.SetLogFuncCallDepth(3)
	logs.Emergency("emergency", "Ex", "Ex")
	logs.Alert("alert", "Ex", "Ex")
	logs.Critical("critical", "Ex", "Ex")
	logs.Error("error", "Ex", "Ex")
	logs.Notice("notice", "Ex", "Ex")
	logs.Debug("debug", "Ex", "Ex")
	logs.Ex(logs.LevelError, map[string]any{"Title": "a/b/c", "ExecDur": 10}, "Ex Log", "asdsda")
	logs.ID(logs.LevelError, "ID128024809123", "WithID", "Ex")
}
