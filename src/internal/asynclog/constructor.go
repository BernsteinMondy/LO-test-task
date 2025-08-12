package asynclog

import (
	"context"
	"log/slog"
)

type AsyncLogger interface {
	Start(ctx context.Context)
	Stop()

	Error(ctx context.Context, msg string, err error, attrs ...slog.Attr)
	Warn(ctx context.Context, msg string, err error, attrs ...slog.Attr)
	Info(ctx context.Context, msg string, attrs ...slog.Attr)
	Debug(ctx context.Context, msg string, attrs ...slog.Attr)
}

func NewAsyncLogger(bufferSize int) AsyncLogger {
	return newAsyncLogger(bufferSize)
}
