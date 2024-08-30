package module

import (
	"context"
	"log"
	internal "std-library/app/internal/module"
	"std-library/app/web"
	"std-library/app/web/metric"
)

type MetricConfig struct {
	moduleContext *Context
	server        *metric.Server
	listen        string
}

func (c *MetricConfig) Initialize(moduleContext *Context, _ string) {
	c.moduleContext = moduleContext
	c.server = &metric.Server{
		HttpHost: &web.HTTPHost{
			Host: "0.0.0.0",
			Port: 8000,
		},
	}

	c.moduleContext.StartupHook.StartStage2 = append(c.moduleContext.StartupHook.StartStage2, c.server)
	c.moduleContext.ShutdownHook.Add(internal.STAGE_8, func(ctx context.Context, timeoutInMs int64) {
		c.server.Shutdown(ctx)
	})
}

func (c *MetricConfig) Validate() {
	if c.server.HttpHost == nil {
		log.Fatal("metrics http listen is not configured, please configure first")
	}
}

func (c *MetricConfig) Listen(listen string) {
	c.server.HttpHost = web.Parse(listen)
}
