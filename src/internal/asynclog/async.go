package asynclog

import (
	"context"
	"log/slog"
	"sync"
)

type asyncLogger struct {
	logCh      chan logEntry
	stopOnce   sync.Once
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
// sync.Once actually guarantees, that Stop() method can be called only once, this
// is a safety measure to prevent user from calling Stop() more, than once.
func (a *asyncLogger) Stop() {
	a.stopOnce.Do(func() {
		close(a.stopCh)
		close(a.logCh)
	})
}

func (a *asyncLogger) processLogEntry(entry logEntry) {
	attrs := append([]slog.Attr{
		slog.String("level", entry.Level.String()),
	}, entry.Attrs...)

	if entry.Err != nil {
		attrs = append(attrs, slog.Any("error", entry.Err))
	}

	slog.LogAttrs(context.Background(), entry.Level, entry.Message, attrs...)
}

func (a *asyncLogger) Log(ctx context.Context, level slog.Level, msg string, err error, attrs ...slog.Attr) {
	select {
	case a.logCh <- logEntry{Level: level, Message: msg, Err: err, Attrs: attrs}:
	case <-a.stopCh:
		slog.Warn("async logger stopped, message dropped", "msg", msg)
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
