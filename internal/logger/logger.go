package logger

import (
	"context"
	"log/slog"
	"os"
)

type Log struct {
	original *slog.Logger
}

func New() Log {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return Log{logger}
}

func (l *Log) Error(msg string, args ...any) {
	l.original.Error(msg, args...)
}

func (l *Log) WithError(err error, msg string, args ...any) {
	passArgs := make([]any, len(args)+2)
	passArgs[0] = "error"
	passArgs[1] = err
	for i, arg := range args {
		passArgs[i+2] = arg
	}
	l.original.Error(msg, passArgs...)
}

func (l *Log) Info(msg string, args ...any) {
	l.original.Info(msg, args...)
}

func (l *Log) Debug(msg string, args ...any) {
	l.original.Debug(msg, args...)
}

func (l *Log) DebugContext(ctx context.Context, msg string, args ...any) {
	l.original.DebugContext(ctx, msg, args...)
}
