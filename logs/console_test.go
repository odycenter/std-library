package logs_test

import (
	"context"
	"github.com/odycenter/std-library/app/log/consts/logKey"
	"github.com/odycenter/std-library/logs"
	"testing"
)

func TestConsole(t *testing.T) {
	logs.SetLogFuncCallDepth(2)
	logs.SetLevel(logs.LevelDebug)
	testConsoleCalls()
}

func testConsoleCalls() {
	logs.Error("error")

	logs.Notice("notice")
	logs.Debug("debug")
	ctx := context.WithValue(context.Background(), logKey.Id, "test-id")
	logs.DebugWithCtx(ctx, "debug, log, name: %v", "test123")
}
