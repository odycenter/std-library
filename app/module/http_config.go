package module

import (
	"context"
	"log"
	internal "std-library/app/internal/module"
	"std-library/app/internal/web"
	"std-library/app/web/beego"
)

type HTTPConfig struct {
	moduleContext *Context
	server        *beego.HTTPServer
}

func (c *HTTPConfig) Initialize(moduleContext *Context, name string) {
	c.moduleContext = moduleContext
	c.server = beego.NewHTTPServer()

	c.moduleContext.StartupHook.StartStage2 = append(c.moduleContext.StartupHook.StartStage2, c.server)

	c.moduleContext.ShutdownHook.Add(internal.STAGE_0, func(ctx context.Context, timeoutInMs int64) {
		c.server.Shutdown(ctx)
	})
	c.moduleContext.ShutdownHook.Add(internal.STAGE_1, func(ctx context.Context, timeoutInMs int64) {
		c.server.AwaitRequestCompletion(ctx, timeoutInMs)
	})
	c.moduleContext.ShutdownHook.Add(internal.STAGE_8, func(ctx context.Context, timeoutInMs int64) {
		c.server.AwaitTermination(ctx)
	})
}

func (c *HTTPConfig) Validate() {
	if c.server.HttpHost == nil {
		log.Fatal("http listen is not configured, please configure first")
	}
}

func (c *HTTPConfig) GZip() {
	c.server.GZip()
}

func (c *HTTPConfig) Listen(listen string) {
	c.server.HttpHost = internal_web.Parse(listen)
}

func (c *HTTPConfig) ErrorWithOkStatus(val bool) {
	c.server.ErrorWithOkStatus(val)
}

func (c *HTTPConfig) CustomErrorResponseMessage(f func(code int, message string) map[string]interface{}) {
	c.server.CustomErrorResponseMessage(f)
}
