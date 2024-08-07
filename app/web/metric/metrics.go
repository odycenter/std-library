package metric

import (
	"context"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"net/http"
	internalweb "std-library/app/internal/web"
	internal_http "std-library/app/internal/web/http"
	"std-library/logs"
)

type Server struct {
	HttpHost *internalweb.HTTPHost
	server   *http.Server
}

func (p *Server) Start(ctx context.Context) {
	p.server = &http.Server{
		Addr: p.HttpHost.String(),
	}
	err := p.server.ListenAndServe()
	if err != nil {
		logs.ErrorWithCtx(ctx, "Failed to start monitor server %v", err)
	}
}

func (p *Server) Shutdown(ctx context.Context) {
	if p.server == nil {
		return
	}
	p.server.Shutdown(ctx)
}

func (p *Server) RegisterGRPC(grpcServer *grpc.Server) {
	grpcPrometheus.Register(grpcServer)
}

func (p *Server) Execute(ctx context.Context) {
	accessHandler := AccessHandler{
		accessControl: &internal_http.IPv4AccessControl{},
	}
	http.Handle(MetricsPath, accessHandler.Handler(promhttp.Handler()))
	logs.WarnWithCtx(ctx, "monitor server Running on http://%s", p.HttpHost.String())
	go p.Start(ctx)
}
