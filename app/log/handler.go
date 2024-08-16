package log

import (
	"context"
	"io"
	"log/slog"
	"path/filepath"
	"std-library/app/log/consts/logKey"
	"std-library/logs"
	"strings"
	"sync/atomic"
)

var loggerLevel atomic.Value
var defaultLevel atomic.Value

type Handler struct {
	*slog.JSONHandler
}

func NewHandler(w io.Writer) *Handler {
	handler := &Handler{}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     handler,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case logKey.Id, "app", slog.LevelKey, "function", "file", "line":
				return a
			case slog.TimeKey:
				return slog.Attr{Key: "@timestamp", Value: a.Value}
			case slog.MessageKey:
				return slog.Attr{Key: "message", Value: a.Value}
			case slog.SourceKey:
				if source, ok := a.Value.Any().(*slog.Source); ok {
					source.File = filepath.Base(source.File)
					return slog.Attr{Key: a.Key, Value: slog.AnyValue(source)}
				}
			default:
				if !strings.HasPrefix(a.Key, "context.") {
					return slog.Attr{Key: "context." + a.Key, Value: a.Value}
				}
			}
			return a
		},
	}

	return &Handler{
		JSONHandler: slog.NewJSONHandler(w, opts),
	}
}

func (h *Handler) SetDefaultLevel(levelStr string) {
	level := ToLevel(levelStr)
	defaultLevel.Store(level)
	h.SetLevel(level)
}

func (h *Handler) DefaultLevel() interface{} {
	return defaultLevel.Load()
}

func (h *Handler) SetLevel(level slog.Level) {
	loggerLevel.Store(level)

	switch level { // TODO remove when logs replace to slog
	case slog.LevelDebug:
		logs.SetLevel(logs.LevelDebug)
	case slog.LevelInfo:
		logs.SetLevel(logs.LevelInformation)
	case slog.LevelWarn:
		logs.SetLevel(logs.LevelWarning)
	case slog.LevelError:
		logs.SetLevel(logs.LevelError)
	}
}

func (h *Handler) Level() slog.Level {
	level, ok := loggerLevel.Load().(slog.Level)
	if !ok {
		return slog.LevelInfo
	}
	return level
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level < h.Level() {
		return nil
	}

	if ctx != nil {
		if id, ok := ctx.Value(logKey.Id).(string); ok {
			r.AddAttrs(slog.String("id", id))
		}
	}
	return h.JSONHandler.Handle(ctx, r)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		JSONHandler: h.JSONHandler.WithAttrs(attrs).(*slog.JSONHandler),
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		JSONHandler: h.JSONHandler.WithGroup(name).(*slog.JSONHandler),
	}
}

func ToLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}

	slog.Warn("Invalid log level, return default INFO level", "level", level)
	return slog.LevelInfo
}
