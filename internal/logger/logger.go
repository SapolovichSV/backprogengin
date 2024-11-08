package logger

import "log/slog"

type logger struct {
	slog *slog.Logger
}

func New(h slog.Handler) *logger {
	logr := slog.New(h)
	return &logger{logr}
}
func (l *logger) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}
func (l *logger) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}
func (l *logger) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}
