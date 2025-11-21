package http

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Server interface allows mocking for tests
type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// StartServer starts the HTTP server with graceful shutdown
func StartServer(ctx context.Context, l *zap.Logger, srv Server) error {
	go func() {
		<-ctx.Done()
		l.Info("Shutting down server...")
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctxShutdown); err != nil {
			l.Error("Shutdown error", zap.Error(err))
		}
	}()

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
