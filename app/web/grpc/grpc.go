package grpc

import (
	"context"
	"fmt"
	prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	app "github.com/odycenter/std-library/app/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"log/slog"
	"net"
	"time"
)

// Server 对grpc package server结构的封装
type Server struct {
	Srv                 *grpc.Server
	HttpListen          string
	Health              *health.Server
	shutdownInterceptor *ShutdownInterceptor
}

func (s *Server) Start(listener net.Listener) {
	err := s.Srv.Serve(listener)
	if err != nil {
		log.Panicf("GRPC service start failed %s\n", err.Error())
	}
}

func (s *Server) Execute(ctx context.Context) {
	defer func() {
		time.Sleep(300 * time.Millisecond)
		s.Health.Resume()
	}()
	listener, err := net.Listen("tcp", s.HttpListen)
	if err != nil {
		log.Panicf("GRPC service listen failed %s\n", err.Error())
	}
	slog.Warn(fmt.Sprintf("grpc server Running on http://%v", s.HttpListen))
	go s.Start(listener)
}

func (s *Server) Shutdown(ctx context.Context) {
	slog.InfoContext(ctx, "shutting down grpc server")
	s.shutdownInterceptor.shutdownHandler.Shutdown()
}

func (s *Server) AwaitRequestCompletion(ctx context.Context, timeoutInMs int64) {
	success := s.shutdownInterceptor.shutdownHandler.AwaitTermination(timeoutInMs)
	if !success {
		slog.WarnContext(ctx, fmt.Sprintf("[FAILED_TO_STOP], failed to wait active grpc requests to complete, due to timeout, canceledRequests=%d", s.shutdownInterceptor.shutdownHandler.ActiveRequests()))
	} else {
		slog.InfoContext(ctx, "active grpc requests completed")
	}
}

func (s *Server) AwaitTermination(ctx context.Context) {
	slog.InfoContext(ctx, "shutting down grpc server")
	s.Srv.GracefulStop()
}

func (s *Server) MaxConnections(maxConnections int32) {
	s.shutdownInterceptor.MaxConnections(maxConnections)
}

// NewServer 创建新的GRPC服务
func NewServer(opt ...grpc.ServerOption) *Server {
	shutdownHandler := NewShutdownHandler()
	options := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(shutdownHandler.handle, serverInterceptor, prometheus.UnaryServerInterceptor),
		grpc.ChainStreamInterceptor(prometheus.StreamServerInterceptor),
	}
	options = append(options, opt...)
	grpcServer := grpc.NewServer(options...)
	healthServer := health.NewServer()
	healthServer.SetServingStatus(app.Name, healthpb.HealthCheckResponse_NOT_SERVING)
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	return &Server{
		grpcServer,
		"",
		healthServer,
		shutdownHandler,
	}
}
