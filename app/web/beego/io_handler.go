package beego

import (
	"net/http"
	"std-library/app/web"
	"std-library/logs"
)

const HealthCheckPath = "/health-check"

type IOHandler struct {
	shutdownHandler *web.ShutdownHandler
}

func (f *IOHandler) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == HealthCheckPath {
			// end exchange will send 200 / content-length=0
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
			return
		}

		f.shutdownHandler.Increment()
		defer f.shutdownHandler.Decrement()

		if f.shutdownHandler.IsShutdown() {
			logs.WarnWithCtx(r.Context(), "reject request due to server is shutting down, requestURL="+r.URL.Path)
			// ask client not set keep alive for current connection, with send header "connection: close",
			// this does no effect with http/2.0, only for http/1.1
			w.Header().Set("Connection", "close")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte{})
			return
		}

		next.ServeHTTP(w, r)
	})
}
