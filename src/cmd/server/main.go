package main

import (
	"context"
	"fmt"
	"lo-test-task/internal/asynclog"
	"lo-test-task/internal/core"
	"lo-test-task/internal/httpserver"
	"lo-test-task/internal/impl/storage"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	err := run()
	if err != nil {
		slog.Error("run() returned error:", slog.String("error", err.Error()))
	}
}

func run() (err error) {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGSTOP, syscall.SIGTERM)
	defer cancel()

	slog.Info("Loading config")
	cfg, err := loadConfigFromEnv()
	if err != nil {
		return fmt.Errorf("load config from env: %w", err)
	}
	slog.Info("Config loaded")
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	slog.Info("Creating storage")
	store := storage.New()
	slog.Info("Storage created")

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	asyncLogger := asynclog.NewAsyncLogger(cfg.Logger.BufferSize)

	stopWg := &sync.WaitGroup{}
	stopWg.Add(1)
	go func(ctx context.Context) {
		defer func() {
			slog.Info("Stopping async logger")
			asyncLogger.Stop()
			slog.Info("Async logger stopped")
			stopWg.Done()
		}()

		slog.Info("Starting async logger")
		asyncLogger.Start(ctx)
	}(ctx)

	service := core.NewService(
		store,
		asyncLogger,
	)

	srv := httpserver.New(service, cfg.HTTPServer.ListenAddr)

	stopWg.Add(1)
	go func(ctx context.Context) {
		defer stopWg.Done()
		srvErr := srv.Launch(ctx)
		if srvErr != nil {
			slog.Error("Server launch error", slog.String("err", srvErr.Error()))
		}
	}(ctx)

	<-ctx.Done()
	stopWg.Wait()
	return nil
}
