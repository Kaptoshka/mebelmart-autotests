package logger

import (
	"io"
	"log/slog"
	"os"
	"testing"
	"time"
)

func New() *slog.Logger {
	level := os.Getenv("LOG_LEVEL")

	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = slog.LevelInfo
	}

	_ = os.MkdirAll("logs", 0o750)

	logFile, err := os.OpenFile(
		"logs/session_"+time.Now().Format("2006-01-02_15-04-05")+".log",
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

func ForTest(t *testing.T) *slog.Logger {
	return New().With("test", t.Name())
}
