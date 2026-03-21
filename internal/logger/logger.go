package logger

import (
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"autotests/internal/config"
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

	logFile, err := os.OpenFile(
		logDir+"/session_"+time.Now().Format("2006-01-02_15-04-05")+".log",
		os.O_CREATE|os.O_WRONLY, 0o600,
	)

	var output io.Writer
	if err != nil {
		output = os.Stdout
	} else {
		output = io.MultiWriter(os.Stdout, logFile)
	}

	return slog.New(slog.NewTextHandler(output, &slog.HandlerOptions{
		Level: lvl,
	}))
}

// ForTest returns a logger with test name context.
func ForTest(t *testing.T) *slog.Logger {
	return New().With("test", t.Name())
}
