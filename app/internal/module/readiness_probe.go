package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

const maxWaitTime = 27 * time.Second

type ReadinessProbe struct {
	hostURIs []string
	urls     []string
}

func (r *ReadinessProbe) AddHostURI(hostURI string) {
	r.hostURIs = append(r.hostURIs, hostURI)
}

func (r *ReadinessProbe) AddURL(url string) {
	r.urls = append(r.urls, url)

}

func (r *ReadinessProbe) Check(ctx context.Context) {
	slog.InfoContext(ctx, "check readiness")
	start := time.Now()

	r.checkDNS(ctx)
	r.checkHTTP(ctx)

	r.hostURIs = nil
	r.urls = nil

	elapsed := time.Since(start)
	slog.InfoContext(ctx, fmt.Sprintf("Readiness check completed in %s", elapsed))
}

func (r *ReadinessProbe) checkDNS(ctx context.Context) {
	for _, hostURI := range r.hostURIs {
		hostname := Hostname(hostURI)
		ResolveHost(ctx, hostname)
	}
}

func (r *ReadinessProbe) checkHTTP(ctx context.Context) {
	for _, url := range r.urls {
		err := r.sendHTTPRequest(ctx, url)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (r *ReadinessProbe) sendHTTPRequest(ctx context.Context, url string) error {
	start := time.Now()
	for {
		resp, err := http.Get(url)
		if err != nil {
			if time.Since(start) >= maxWaitTime {
				return errors.New("readiness check failed, url=" + url)
			}
			slog.WarnContext(ctx, fmt.Sprintf("[NOT_READY] http probe failed, retry soon, url=%s", url))
			time.Sleep(5 * time.Second)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
	}
}

func ResolveHost(ctx context.Context, hostname string) {
	start := time.Now()
	for {
		_, err := net.LookupHost(hostname)
		if err != nil {
			if time.Since(start) >= maxWaitTime {
				log.Fatal("readiness check failed, host=" + hostname)
			}
			slog.WarnContext(ctx, fmt.Sprintf("[NOT_READY] dns probe failed, retry soon, host=%s", hostname))
			time.Sleep(5 * time.Second)
			continue
		}
		return
	}
}

func Hostname(hostURI string) string {
	index := strings.Index(hostURI, ":")
	if index == -1 {
		return hostURI
	}
	return hostURI[:index]
}
