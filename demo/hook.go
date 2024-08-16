package demo

import (
	"context"
	"log/slog"
	"std-library/app/module"
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
