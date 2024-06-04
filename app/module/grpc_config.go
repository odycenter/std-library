package module

import (
	"context"
	"google.golang.org/grpc"
	"log"
	internal "std-library/app/internal/module"
	"std-library/app/internal/web"
	server "std-library/app/web/grpc"
)

type GrpcServerConfig struct {
	moduleContext *Context
	grpcServer    *server.Server
	listen        string
	opt           []grpc.ServerOption
}

func (c *GrpcServerConfig) Initialize(moduleContext *Context, name string) {
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

	host := internal_web.Parse(c.listen)
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

	return c.grpcServer.Srv
}
