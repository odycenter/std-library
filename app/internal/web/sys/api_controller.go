package internal_sys

import (
	"embed"
	"log/slog"
	"net/http"
	"std-library/app/internal/web/http"
	"std-library/app/web/errors"
	"std-library/nets"
	"strings"
	"sync"
)

//go:embed api.html
var fs embed.FS
var once sync.Once
var htmlContent []byte
var ApiJsonContent []byte
var ApiJsonGzipped []byte

type APIController struct {
	AccessControl *internal_http.IPv4AccessControl
}

func NewAPIController() *APIController {
	once.Do(func() {
		content, err := fs.ReadFile("api.html")
		if err == nil {
			htmlContent = content
		} else {
			slog.Error("failed to read api.html")
		}
	})
	return &APIController{
		AccessControl: &internal_http.IPv4AccessControl{},
	}
}

func (c *APIController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.AccessControl.Validate(nets.IP(r).String())
	if err != nil {
		errors.Forbidden("access denied", "IP_ACCESS_DENIED")
	}

	if r.URL.Path != "/_sys/api" && r.URL.Path != "/_sys/api.json" {
		errors.NotFound("not found")
		return
	}

	if r.Method == http.MethodGet && r.URL.Path == "/_sys/api.json" {
		if len(ApiJsonContent) == 0 {
			errors.NotFound("api.json not configured")
			return
		}

		w.Header().Set("Cache-Control", "public, max-age=60")
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && len(ApiJsonGzipped) > 0 {
			w.Header().Set("Content-Encoding", "gzip")
			w.Write(ApiJsonGzipped)
		} else {
			w.Write(ApiJsonContent)
		}
		return
	}

	if htmlContent == nil {
		errors.Internal("api.html content not available")
		return
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "text/html")
	w.Write(htmlContent)
}
