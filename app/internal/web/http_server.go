package web

import (
	"context"
	"fmt"
	beegoWeb "github.com/beego/beego/v2/server/web"
	"github.com/odycenter/std-library/app/web"
	"log/slog"
	"time"
)

type HTTPServer struct {
	server          *beegoWeb.HttpServer
	shutdownHandler *web.ShutdownHandler
	ioHandler       *IOHandler
	handler         *HTTPHandler
	HttpHost        *web.HTTPHost
}

func NewHTTPServer() *HTTPServer {
	shutdownHandler := web.NewShutdownHandler()
	ioHandler := &IOHandler{
		shutdownHandler: shutdownHandler,
	}
	server := &HTTPServer{
		server:          beegoWeb.BeeApp,
		shutdownHandler: shutdownHandler,
		ioHandler:       ioHandler,
		handler: &HTTPHandler{
			ErrorWithOkStatus: false,
		},
	}
	beegoWeb.BConfig.RecoverPanic = false
	beegoWeb.BConfig.RunMode = beegoWeb.PROD
	beegoWeb.ErrorController(&errorController{})

	return server
}

func (s *HTTPServer) GZip() {
	beegoWeb.BConfig.EnableGzip = true
}

func (s *HTTPServer) ErrorWithOkStatus(val bool) {
	s.handler.ErrorWithOkStatus = val
}

func (s *HTTPServer) CustomErrorResponseMessage(f func(code int, message string) map[string]interface{}) {
	s.handler.CustomErrorResponseMessage = f
}

func (s *HTTPServer) Execute(ctx context.Context) {
	slog.WarnContext(ctx, fmt.Sprintf("web server Running on http://%v", s.HttpHost.String()))
	go s.Start()
}

func (s *HTTPServer) Start() {
	s.server.Run(s.HttpHost.String(), s.ioHandler.Handler, s.handler.Handler)
}

func (s *HTTPServer) Shutdown(ctx context.Context) {
	slog.InfoContext(ctx, "shutting down web server")
	s.shutdownHandler.Shutdown()
}

func (s *HTTPServer) AwaitRequestCompletion(ctx context.Context, timeoutInMs int64) {
	success := s.shutdownHandler.AwaitTermination(timeoutInMs)
	if !success {
		slog.Warn(fmt.Sprintf("[FAILED_TO_STOP], failed to wait active http requests to complete, due to timeout, canceledRequests=%d", s.shutdownHandler.ActiveRequests()))
	} else {
		slog.InfoContext(ctx, "active web requests completed")
	}
}

func (s *HTTPServer) AwaitTermination(ctx context.Context) {
	slog.InfoContext(ctx, "shutting down http server")
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	s.server.Server.Shutdown(ctx)
}
