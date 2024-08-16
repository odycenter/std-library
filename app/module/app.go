package module

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	actionlog "std-library/app/log"
	"std-library/app/log/consts/logKey"
	"std-library/app/log/dto"
	"syscall"
)

type Module interface {
	SetContext(moduleContext *Context)
	Initialize()
}

type AppInterface interface {
	Module
	Configure()
	Start()
}

type App struct {
	Common
	actionLog dto.ActionLog
	ctx       context.Context
}

func Start(app AppInterface) {
	app.Configure()
	app.Initialize()
	app.Start()
}

func (a *App) Configure() {
	a.actionLog = actionlog.Begin("app:start", "app")
	a.ctx = context.WithValue(context.Background(), logKey.Id, a.actionLog.Id)
	moduleContext := &Context{}
	moduleContext.Initialize()
	a.SetContext(moduleContext)
}

func (a *App) Start() {
	a.ModuleContext.Validate()

	a.ModuleContext.Probe.Check(a.ctx)
	slog.InfoContext(a.ctx, "execute startup tasks")
	a.ModuleContext.StartupHook.DoInitialize(a.ctx)
	a.ModuleContext.StartupHook.DoStart(a.ctx)
	slog.InfoContext(a.ctx, fmt.Sprintf("startup completed, elapsed=%v", a.actionLog.Elapsed()))
	a.cleanup()
	actionlog.End(a.actionLog, "ok")
	a.shutdownHook()
}

func (a *App) cleanup() {
	a.ModuleContext.Cleanup()
}

func (a *App) shutdownHook() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Wait for a shutdown signal
	<-stop

	a.ModuleContext.ShutdownHook.Run()
}
