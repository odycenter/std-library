package module

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	app "std-library/app/conf"
	actionlog "std-library/app/log"
	"std-library/app/log/consts/logKey"
	"std-library/app/log/dto"
	"strconv"
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
	a.startDefaultHttpServer()
	a.ModuleContext.Validate()

	a.ModuleContext.Probe.Check(a.ctx)
	slog.InfoContext(a.ctx, "execute startup tasks")
	a.ModuleContext.StartupHook.DoInitialize(a.ctx)
	a.ModuleContext.StartupHook.DoStart(a.ctx)
	slog.InfoContext(a.ctx, "startup completed", "elapsed", a.actionLog.Elapsed())
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

func (a *App) startDefaultHttpServer() {
	if a.ModuleContext.httpConfigAdded {
		return
	}

	listenPort := 18080

	if app.Local() {
		listenPort = a.getAvailablePort(listenPort)
	} else if _, ok := a.ModuleContext.listenPorts.Load(listenPort); ok || !isPortAvailable(listenPort) {
		slog.ErrorContext(a.ctx, "no available ports found for HTTP server")
		return
	}

	slog.WarnContext(a.ctx, "No http listen port configured, using candidate port to start HTTP server", "port", listenPort, "local", app.Local())

	a.ModuleContext.PropertyManager.DefaultHTTPPort = listenPort
	a.Common.Http().Listen(strconv.Itoa(listenPort))
}

func (a *App) getAvailablePort(port int) int {
	for {
		if isPortAvailable(port) {
			break
		}
		port++
	}
	return port
}

func isPortAvailable(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}
