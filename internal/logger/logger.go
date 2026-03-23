package logger

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"autotests/internal/config"

	"github.com/lmittmann/tint"
)

// New creates a new instance of Logger with level from environment variables.
func New() *slog.Logger {
	level := os.Getenv("LOG_LEVEL")

	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = slog.LevelInfo
	}

	logDir := config.DefaultLogDir
	_ = os.MkdirAll(logDir, 0o750)

	consoleHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      lvl,
		TimeFormat: "2006-01-02T15:04:05",
	})

	logFile, err := os.OpenFile(
		logDir+"/session_"+time.Now().Format("2006-01-02_15-04-05")+".log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600,
	)

	var handler slog.Handler
	if err != nil {
		handler = consoleHandler
	} else {
		fileHandler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
			Level: lvl,
		})
		handler = &multiHandler{
			handlers: []slog.Handler{consoleHandler, fileHandler},
		}
	}

	return slog.New(handler)
}

type multiHandler struct {
	handlers []slog.Handler
}

func (h *multiHandler) Enabled(
	ctx context.Context,
	l slog.Level,
) bool {
	return h.handlers[0].Enabled(ctx, l)
}

func (h *multiHandler) Handle(
	ctx context.Context,
	r slog.Record,
) error {
	for _, handler := range h.handlers {
		if err := handler.Handle(ctx, r); err != nil {
			return err
		}
	}

	return nil
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &multiHandler{handlers: newHandlers}
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &multiHandler{handlers: newHandlers}
}

// ForTest returns a logger with test name context.
func ForTest(t *testing.T) *slog.Logger {
	return New().With("test", t.Name())
}
