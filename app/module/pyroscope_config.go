package module

import (
	"context"
	"log"
	app "std-library/app/conf"
	internal "std-library/app/internal/module"
	"std-library/pyroscope"
	"strings"
)

type PyroscopeConfig struct {
	uri             string
	config          *pyroscope.Config
	forceLocalStart bool
}

func (c *PyroscopeConfig) Initialize(moduleContext *Context, name string) {
	c.config = &pyroscope.Config{
		ApplicationName: app.Name,
		LogLevel:        1,
		OpenGMB:         false,
	}
	moduleContext.StartupHook.Add(c)
	moduleContext.ShutdownHook.Add(internal.STAGE_7, func(ctx context.Context, timeoutInMs int64) {
		c.stop(ctx, timeoutInMs)
	})
}

func (c *PyroscopeConfig) Validate() {
	if len(c.uri) == 0 {
		log.Fatalf("pyroscope uri is not configured")
	}
}

func (c *PyroscopeConfig) ForceLocalStart() {
	c.forceLocalStart = true
}

func (c *PyroscopeConfig) Uri(uri string) {
	if c.uri != "" {
		log.Fatalf("pyroscope uri is already configured, uri=%s, previous=%s", uri, c.uri)
	}
	uri = strings.TrimSpace(uri)
	if !strings.Contains(uri, ":") {
		uri = "http://" + uri + ":4040"
	}
	c.uri = uri
	c.config.ServerAddress = uri
}

func (c *PyroscopeConfig) LogLevel(level int) {
	c.config.LogLevel = level
}

func (c *PyroscopeConfig) OpenGMB(openGMB bool) {
	c.config.OpenGMB = openGMB
}

func (c *PyroscopeConfig) AuthToken(token string) {
	c.config.AuthToken = token
}

func (c *PyroscopeConfig) start(_ context.Context) {
	if c.forceLocalStart || !app.Local() {
		go pyroscope.Start(c.config)
	}
}

func (c *PyroscopeConfig) stop(_ context.Context, timeoutInMs int64) {
	go pyroscope.Stop()
}

func (c *PyroscopeConfig) Execute(ctx context.Context) {
	c.start(ctx)
}
