package module

import (
	"context"
	"github.com/beego/beego/v2/server/web"
	"log/slog"
	"os"
	internalLog "std-library/app/internal/log"
	internal "std-library/app/internal/module"
	internalWeb "std-library/app/internal/web"
	"std-library/app/internal/web/sys"
	"std-library/app/property"
	"sync"
)

type Context struct {
	StartupHook       *internal.StartupHook
	ShutdownHook      *internal.ShutdownHook
	Probe             *internal.ReadinessProbe
	PropertyManager   *property.Manager
	propertyValidator *property.Validator
	configs           sync.Map // map[string]Config
	listenPorts       sync.Map // map[int]bool
	httpServer        *internalWeb.HTTPServer
	httpConfigAdded   bool
}

func (m *Context) Initialize() {
	handler := internalLog.NewHandler(os.Stdout)
	handler.SetLevel(slog.LevelDebug)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	m.StartupHook = &internal.StartupHook{}
	m.ShutdownHook = &internal.ShutdownHook{}
	m.Probe = &internal.ReadinessProbe{}
	m.ShutdownHook.Initialize()
	m.PropertyManager = property.NewManager()
	m.propertyValidator = property.NewValidator()
	m.httpServer = m.createHTTPServer()

	web.Handler("/_sys/property", internal_sys.NewPropertyController(m.PropertyManager))
}

func (m *Context) createHTTPServer() *internalWeb.HTTPServer {
	httpServer := internalWeb.NewHTTPServer()

	m.StartupHook.StartStage2 = append(m.StartupHook.StartStage2, httpServer)
	m.ShutdownHook.Add(internal.STAGE_0, func(ctx context.Context, timeoutInMs int64) {
		httpServer.Shutdown(ctx)
	})
	m.ShutdownHook.Add(internal.STAGE_1, func(ctx context.Context, timeoutInMs int64) {
		httpServer.AwaitRequestCompletion(ctx, timeoutInMs)
	})
	m.ShutdownHook.Add(internal.STAGE_8, func(ctx context.Context, timeoutInMs int64) {
		httpServer.AwaitTermination(ctx)
	})

	return httpServer
}

func (m *Context) AddListenPort(port int) {
	if _, exists := m.listenPorts.Load(port); exists {
		slog.Warn("Port already added", "port", port)
		return
	}
	m.listenPorts.Store(port, true)
}

func (m *Context) ConfigByType(configType, name string, newConfig func() Config) Config {
	cfg, ok := m.configs.Load(configType + ":" + name)
	if !ok {
		config := newConfig()
		config.Initialize(m, name)
		m.configs.Store(configType+":"+name, config)
		cfg = config
	}
	return cfg.(Config)
}

func (m *Context) Config(name string, newConfig func() Config) Config {
	cfg, ok := m.configs.Load(name)
	if !ok {
		config := newConfig()
		config.Initialize(m, name)
		m.configs.Store(name, config)
		cfg = config
	}
	return cfg.(Config)
}

func (m *Context) Property(key string) string {
	m.propertyValidator.Add(key)
	return m.PropertyManager.Get(key)
}

func (m *Context) Validate() {
	keys := m.PropertyManager.Keys()
	m.propertyValidator.Validate(keys)

	m.configs.Range(func(key, value interface{}) bool {
		value.(Config).Validate()
		return true
	})
}

func (m *Context) Cleanup() {
	m.StartupHook = nil
	m.PropertyManager = nil
	m.propertyValidator = nil
	m.configs.Range(func(key, value interface{}) bool {
		m.configs.Delete(key)
		return true
	})
}
