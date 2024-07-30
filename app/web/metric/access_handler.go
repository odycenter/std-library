package metric

import (
	"net/http"
	internal_http "std-library/app/internal/web/http"
	"std-library/logs"
	"std-library/nets"
)

const MetricsPath = "/metrics"

type AccessHandler struct {
	accessControl *internal_http.IPv4AccessControl
}

func (h *AccessHandler) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := nets.IP(r).String()
		err := h.accessControl.Validate(ip)
		if err != nil {
			logs.WarnWithCtx(r.Context(), "access metrics denied, ip="+ip)
			w.Header().Set("Connection", "close")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("access denied"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
