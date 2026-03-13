package logger

import (
	"io"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	globalLogger *slog.Logger
	once         sync.Once
)

func initRoot(level string) {
	once.Do(func() {
		var lvl slog.Level
		if err := lvl.UnmarshalText([]byte(level)); err != nil {
			lvl = slog.LevelInfo
		}

		_ = os.MkdirAll("logs", os.ModePerm)
		logFile, err := os.OpenFile(
			"logs/session_"+time.Now().Format("2006-01-02_15-04-05")+".log",
			os.O_CREATE|os.O_WRONLY, 0o666,
		)
		if err != nil {
			panic(err)
		}

		handler := slog.NewTextHandler(
			io.MultiWriter(
				os.Stdout,
				logFile,
			),

			&slog.HandlerOptions{
				Level: lvl,
			},
		)

		globalLogger = slog.New(handler)
		slog.SetDefault(globalLogger)
	})
}

func ForTest(t *testing.T) *slog.Logger {
	initRoot("info")
	return globalLogger.With("test", t.Name())
}
