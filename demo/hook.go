package demo

import (
	"context"
	"std-library/app/module"
	"std-library/logs"
)

type OnStartupTask struct {
}

func (t OnStartupTask) Execute(ctx context.Context) {
	logs.InfoWithCtx(ctx, "call on startup")
}

type OnShutdownTask struct {
}

func (t OnShutdownTask) Execute(ctx context.Context) {
	logs.InfoWithCtx(ctx, "call on Shutdown")
}

type HookModule struct {
	module.Common
}

func (m *HookModule) Initialize() {
	m.OnStartup(&OnStartupTask{})
	m.OnShutdown(&OnShutdownTask{})
}
