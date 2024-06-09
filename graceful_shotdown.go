package xhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// ServeWithShutdown starts the s server and waits for syscall.SIGINT or syscall.SIGTERM,
// then preforms a graceful shutdown
func ServeWithShutdown(s *http.Server) error {
	ctx := context.Background()
	ctx, cFunc := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cFunc()

	errChan := make(chan error, 1)
	defer close(errChan)

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, shutdownCFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCFunc()
		if err := s.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
	case err := <-errChan:
		return fmt.Errorf("server listen and serve failed: %w", err)
	}
	return nil
}
