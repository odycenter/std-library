package module

import (
	"context"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"log"
	internal "std-library/app/internal/module"
	"std-library/app/internal/web"
	server "std-library/app/web/grpc"
)

type GrpcServerConfig struct {
	moduleContext  *Context
	grpcServer     *server.Server
	maxConnections int32
	listen         string
	opt            []grpc.ServerOption
}

func (c *GrpcServerConfig) Initialize(moduleContext *Context, _ string) {
	c.moduleContext = moduleContext
}

func (c *GrpcServerConfig) Validate() {
	if c.grpcServer == nil {
		log.Fatal("grpc server not configured, please configure first!")
	}
}

func (c *GrpcServerConfig) AddOpt(opt grpc.ServerOption) *GrpcServerConfig {
	if c.grpcServer != nil {
		log.Fatal("grpc server already configured, cannot add option!")
	}
	c.opt = append(c.opt, opt)

	return c
}

func (c *GrpcServerConfig) MaxConnections(maxConnections int32) {
	if c.grpcServer != nil {
		log.Fatal("grpc server already configured, cannot set maxConnections!")
	}
	c.maxConnections = maxConnections
}

func (c *GrpcServerConfig) Listen(listen string) {
	c.listen = listen
}

func (c *GrpcServerConfig) Server() *grpc.Server {
	if c.listen == "" {
		log.Fatal("grpc listen is not configured, please configure first")
	}
	if c.grpcServer != nil {
		return c.grpcServer.Srv
	}

	c.grpcServer = server.NewServer(c.opt...)
	if c.maxConnections > 0 {
		c.grpcServer.MaxConnections(c.maxConnections)
	}

	host := internal_web.Parse(c.listen)
	c.moduleContext.AddListenPort(host.Port)
	c.grpcServer.HttpListen = host.String()
	c.moduleContext.StartupHook.StartStage2 = append(c.moduleContext.StartupHook.StartStage2, c.grpcServer)
	c.moduleContext.ShutdownHook.Add(internal.STAGE_0, func(ctx context.Context, timeoutInMs int64) {
		c.grpcServer.Shutdown(ctx)
	})
	c.moduleContext.ShutdownHook.Add(internal.STAGE_1, func(ctx context.Context, timeoutInMs int64) {
		c.grpcServer.AwaitRequestCompletion(ctx, timeoutInMs)
	})
	c.moduleContext.ShutdownHook.Add(internal.STAGE_8, func(ctx context.Context, timeoutInMs int64) {
		c.grpcServer.AwaitTermination(ctx)
	})

	grpcPrometheus.Register(c.grpcServer.Srv)

	return c.grpcServer.Srv
}
