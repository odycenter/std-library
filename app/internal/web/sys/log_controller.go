package internal_sys

import (
	"fmt"
	"log/slog"
	"net/http"
	"std-library/app/internal/web/http"
	"std-library/app/log"
	"std-library/app/web/errors"
	"std-library/json"
	"std-library/nets"
	"strings"
	"sync"
	"time"
)

var (
	defaultManualLevelDuration = 10 * time.Minute
	maxManualLevelDuration     = 30 * time.Minute
)

type LogController struct {
	accessControl  *internal_http.IPv4AccessControl
	handler        *log.Handler
	changeTime     time.Time
	resetTimer     *time.Timer
	manualDuration time.Duration
	mutex          sync.Mutex
}

func NewLogController(logHandler *log.Handler) *LogController {
	return &LogController{
		accessControl: &internal_http.IPv4AccessControl{},
		handler:       logHandler,
	}
}

func (c *LogController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.accessControl.Validate(nets.IP(r).String())
	if err != nil {
		errors.Forbidden("access denied", "IP_ACCESS_DENIED")
	}

	if r.Method == http.MethodGet && r.URL.Path == "/_sys/log" {
		c.handleGet(w)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/_sys/log/")
	parts := strings.Split(path, "/")
	if len(parts) < 1 || len(parts) > 2 || r.Method != http.MethodPut {
		errors.NotFound("not found")
	}

	c.handlePut(w, r, parts)
}

func (c *LogController) handleGet(w http.ResponseWriter) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	result := map[string]interface{}{
		"current_level": c.handler.Level(),
		"default_level": c.handler.DefaultLevel().(slog.Level),
	}

	if c.resetTimer != nil {
		remainingTime := c.manualDuration - time.Since(c.changeTime)
		if remainingTime < 0 {
			remainingTime = 0
		}
		result["remaining_time"] = remainingTime.String()
		result["manual_duration"] = c.manualDuration.String()
		result["change_time"] = c.changeTime.Format(time.RFC3339)
	}

	w.Write(json.Stringify(result))
}

func (c *LogController) handlePut(w http.ResponseWriter, r *http.Request, parts []string) {
	levelStr := parts[0]
	if levelStr == "" {
		errors.BadRequest("invalid log level, level is empty")
	}

	level := log.ToLevel(levelStr)
	if level == c.handler.Level() {
		errors.BadRequest("invalid log level, level not changed")
	}

	duration := defaultManualLevelDuration
	if len(parts) == 2 {
		var err error
		duration, err = time.ParseDuration(parts[1])
		if err != nil {
			errors.BadRequest("invalid duration format")
		}
		if duration > maxManualLevelDuration {
			duration = maxManualLevelDuration
		}
	}

	c.setLevelWithTimer(level, duration)

	ctx := r.Context()
	slog.WarnContext(ctx, fmt.Sprintf("[MANUAL_OPERATION] change log level manually for %v, level=%s", duration, level))
	log.Context(&ctx, "manual_operation", true)
	id := log.GetId(&ctx)
	w.WriteHeader(202)

	w.Write([]byte(fmt.Sprintf("log level changed, level=%s, id=%s, duration=%v", level.String(), id, duration)))
}

func (c *LogController) setLevelWithTimer(level slog.Level, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.handler.SetLevel(level)
	c.changeTime = time.Now()
	c.manualDuration = duration

	if c.resetTimer != nil {
		c.resetTimer.Stop()
	}

	c.resetTimer = time.AfterFunc(duration, func() {
		c.resetToDefaultLevel()
	})
}

func (c *LogController) resetToDefaultLevel() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	defaultLevel, ok := c.handler.DefaultLevel().(slog.Level)
	if !ok {
		defaultLevel = slog.LevelInfo
	}

	c.handler.SetLevel(defaultLevel)
	c.resetTimer = nil
	c.manualDuration = 0
	slog.Warn("[AUTO_OPERATION] Reset log level to default", "level", defaultLevel)
}
