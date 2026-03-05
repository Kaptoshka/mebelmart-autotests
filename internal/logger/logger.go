package logger

import (
	"log/slog"
	"os"
)

var globalLogger *slog.Logger

func Init(level string) {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = slog.LevelInfo
	}

	globalLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))

	slog.SetDefault(globalLogger)
}

func Get() *slog.Logger {
	if globalLogger == nil {
		Init("info")
	}

	return globalLogger
}

func WithTest(testName string) *slog.Logger {
	return Get().With("test", testName)
}

func WithPage(pageName string) *slog.Logger {
	return Get().With("page", pageName)
}

func WithAction(actionName string) *slog.Logger {
	return Get().With("action", actionName)
}

func Step(msg string, args ...any) {
	Get().Info("-- STEP: "+msg, args...)
}

func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	Get().Error(msg, args...)
	os.Exit(1)
}
