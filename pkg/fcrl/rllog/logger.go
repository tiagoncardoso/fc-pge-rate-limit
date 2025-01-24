package rllog

import "log/slog"

func Info(msg string, args ...any) {
	slog.Info("RateLimiter middleware> "+msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn("RateLimiter middleware> "+msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error("RateLimiter middleware> "+msg, args...)
}
