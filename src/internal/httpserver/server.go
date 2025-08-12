package httpserver

import (
	"context"
	"errors"
	"fmt"
	"lo-test-task/internal/core"
	"lo-test-task/internal/httpserver/handlers"
	"log/slog"
	"net/http"
)

type Server struct {
	server *http.Server
}

func New(service *core.Service, listenAddr string) *Server {
	mux := http.NewServeMux()
	handlers.Map(mux, service)

	server := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	srv := &Server{
		server: server,
	}

	return srv
}

func (s *Server) Launch(ctx context.Context) (err error) {
	var (
		addr          = s.server.Addr
		shutdownError error
	)

	defer func() {
		err = errors.Join(err, shutdownError)
	}()

	shutdownDone := make(chan struct{})
	go func(ctx context.Context) {
		<-ctx.Done()

		slog.Info("Shutting down HTTP server")
		shutdownError = s.server.Shutdown(ctx)
		slog.Info("HTTP server shut down")

		close(shutdownDone)
	}(ctx)

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	slog.Info("Starting to listen on specified address", slog.String("address", addr))
	err = s.server.ListenAndServe()
	if err != nil {
		if errors.Is(http.ErrServerClosed, err) {
			return nil
		}
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	<-shutdownDone

	return nil
}
