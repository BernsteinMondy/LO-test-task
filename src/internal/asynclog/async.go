package asynclog

import (
	"context"
	"log/slog"
)

type asyncLogger struct {
	logCh      chan logEntry
	stopCh     chan struct{}
	bufferSize int
}

type logEntry struct {
	Level   slog.Level
	Message string
	Err     error
	Attrs   []slog.Attr
}

// NewAsyncLogger creates new asynchronous logger with buffer size.
func newAsyncLogger(bufferSize int) *asyncLogger {
	return &asyncLogger{
		logCh:      make(chan logEntry, bufferSize),
		stopCh:     make(chan struct{}),
		bufferSize: bufferSize,
	}
}

// Start starts logs event handler. Start launches an infinite for loop for
// reading messages.
func (a *asyncLogger) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-a.stopCh:
			return
		case entry := <-a.logCh:
			a.processLogEntry(entry)
		}
	}
}

// Stop gracefully stops logger.
func (a *asyncLogger) Stop() {
	close(a.logCh)
	close(a.stopCh)
}

func (a *asyncLogger) processLogEntry(entry logEntry) {
	attrs := append([]slog.Attr{}, entry.Attrs...)

	if entry.Err != nil {
		attrs = append(attrs, slog.Any("error", entry.Err))
	}

	slog.LogAttrs(context.Background(), entry.Level, entry.Message, attrs...)
}

func (a *asyncLogger) Log(ctx context.Context, level slog.Level, msg string, err error, attrs ...slog.Attr) {
	select {
	case <-a.stopCh:
		return
	case a.logCh <- logEntry{Level: level, Message: msg, Err: err, Attrs: attrs}:
	default:
		slog.Warn("async logger buffer full, message dropped", "msg", msg)
	}
}

func (a *asyncLogger) Error(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	a.Log(ctx, slog.LevelError, msg, err, attrs...)
}

func (a *asyncLogger) Warn(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	a.Log(ctx, slog.LevelWarn, msg, err, attrs...)
}

func (a *asyncLogger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	a.Log(ctx, slog.LevelInfo, msg, nil, attrs...)
}

func (a *asyncLogger) Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	a.Log(ctx, slog.LevelDebug, msg, nil, attrs...)
}
