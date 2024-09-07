package demo

import (
	"context"
	"github.com/odycenter/std-library/app/module"
	"log/slog"
)

type OnStartupTask struct {
}

func (t OnStartupTask) Execute(ctx context.Context) {
	slog.InfoContext(ctx, "call on startup")
}

type OnShutdownTask struct {
}

func (t OnShutdownTask) Execute(ctx context.Context) {
	slog.InfoContext(ctx, "call on Shutdown")
}

type HookModule struct {
	module.Common
}

func (m *HookModule) Initialize() {
	m.OnStartup(&OnStartupTask{})
	m.OnShutdown(&OnShutdownTask{})
}
