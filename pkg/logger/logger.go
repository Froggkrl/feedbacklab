// Package logger provides logging functionality for the application.
package logger

import (
	"log/slog"
	"os"
)

var Global *slog.Logger

func Init() {
	Global = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	Global.Info("logger initialized")
}

func Info(msg string, args ...interface{}) {
	Global.Info(msg, args...)
}

func Error(msg string, args ...interface{}) {
	Global.Error(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	Global.Warn(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	Global.Debug(msg, args...)
}
