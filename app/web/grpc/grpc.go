package grpc

import (
	"context"
	prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	app "std-library/app/conf"
	"std-library/logs"
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
	logs.WarnWithCtx(ctx, "grpc server Running on http://%v", s.HttpListen)
	go s.Start(listener)
}

func (s *Server) Shutdown(ctx context.Context) {
	logs.InfoWithCtx(ctx, "shutting down grpc server")
	s.shutdownInterceptor.shutdownHandler.Shutdown()
}

func (s *Server) AwaitRequestCompletion(ctx context.Context, timeoutInMs int64) {
	success := s.shutdownInterceptor.shutdownHandler.AwaitTermination(timeoutInMs)
	if !success {
		logs.WarnWithCtx(ctx, "[FAILED_TO_STOP], failed to wait active grpc requests to complete, due to timeout, canceledRequests=%d", s.shutdownInterceptor.shutdownHandler.ActiveRequests())
	} else {
		logs.InfoWithCtx(ctx, "active grpc requests completed")
	}
}

func (s *Server) AwaitTermination(ctx context.Context) {
	logs.InfoWithCtx(ctx, "shutting down grpc server")
	s.Srv.GracefulStop()
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
