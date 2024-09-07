package grpc

import (
	"context"
	"github.com/odycenter/std-library/app/web"
	"github.com/odycenter/std-library/app/web/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShutdownInterceptor struct {
	shutdownHandler *web.ShutdownHandler
	maxConnections  int32
}

func NewShutdownHandler() *ShutdownInterceptor {
	return &ShutdownInterceptor{
		shutdownHandler: web.NewShutdownHandler(),
	}
}

func (s *ShutdownInterceptor) MaxConnections(max int32) {
	s.maxConnections = max
	metric.GRPCMaxConnections.Set(float64(max))
}

func (s *ShutdownInterceptor) handle(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	metric.GRPCConnectionAttempts.Inc()

	if s.maxConnections > 0 {
		current := s.shutdownHandler.ActiveRequests()
		if current >= s.maxConnections {
			metric.GRPCConnectionRejections.Inc()
			return nil, status.Errorf(codes.Unavailable, "max requests reached: current %d, max %d", current, metric.GRPCMaxConnections)
		}
	}

	s.shutdownHandler.Increment()
	metric.GRPCActiveConnections.Inc()

	defer func() {
		s.shutdownHandler.Decrement()
		metric.GRPCActiveConnections.Dec()
	}()

	if s.shutdownHandler.IsShutdown() {
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}

	return handler(ctx, req)
}
