package grpc

import (
	"context"
	"google.golang.org/grpc"
	"std-library/app/web"
	"std-library/app/web/errors"
)

type ShutdownInterceptor struct {
	shutdownHandler *web.ShutdownHandler
}

func NewShutdownHandler() *ShutdownInterceptor {
	return &ShutdownInterceptor{
		shutdownHandler: web.NewShutdownHandler(),
	}
}

func (s *ShutdownInterceptor) handle(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	s.shutdownHandler.Increment()
	defer s.shutdownHandler.Decrement()

	if s.shutdownHandler.IsShutdown() {
		return nil, errors.NewServiceUnavailable(503)
	}

	return handler(ctx, req)
}
