// Copyright (c) 2025 A Bit of Help, Inc.

// Package shutdown provides functionality for graceful application shutdown.
package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// GracefulShutdown waits for termination signals and calls the provided shutdown function.
// It handles OS signals (SIGINT, SIGTERM, SIGHUP) and context cancellation to trigger
// graceful shutdown. It also handles multiple signals, forcing exit if a second signal
// is received during shutdown. A default timeout of 30 seconds is applied to the shutdown
// function to prevent hanging.
//
// Parameters:
//   - ctx: Context that can be cancelled to trigger shutdown
//   - logger: Logger for recording shutdown events
//   - shutdownFunc: Function to execute during shutdown
//
// Returns:
//   - The error from the shutdown function, if any
func GracefulShutdown(ctx context.Context, logger *logging.ContextLogger, shutdownFunc func() error) error {
	// Create a channel to receive OS signals with buffer size 2 to avoid missing signals
	quit := make(chan os.Signal, 2)

	// Register for SIGINT, SIGTERM, and SIGHUP
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Make sure to stop signal handling when we're done
	defer signal.Stop(quit)

	// Create a channel to handle multiple signals
	done := make(chan struct{})
	var shutdownErr error

	// Start a goroutine to handle signals
	go func() {
		defer close(done)

		// Wait for either interrupt signal or context cancellation
		var shutdownReason string
		var shutdownDetails []zap.Field

		select {
		case sig := <-quit:
			shutdownReason = "Received termination signal"
			shutdownDetails = []zap.Field{
				zap.String("signal", sig.String()),
				zap.String("type", "first"),
			}
		case <-ctx.Done():
			shutdownReason = "Context cancelled, shutting down"
			shutdownDetails = []zap.Field{
				zap.Error(ctx.Err()),
			}
		}

		logger.Info(ctx, shutdownReason, shutdownDetails...)

		// Create a timeout context for the shutdown function
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Execute shutdown in a separate goroutine so we can monitor for timeout
		shutdownDone := make(chan struct{})
		go func() {
			defer close(shutdownDone)
			// Call the shutdown function
			logger.Info(ctx, "Executing shutdown function")
			shutdownErr = shutdownFunc()
		}()

		// Wait for shutdown to complete or timeout
		select {
		case <-shutdownDone:
			// Shutdown completed within timeout
			if shutdownErr != nil {
				logger.Error(ctx, "Error during shutdown", zap.Error(shutdownErr))
			} else {
				logger.Info(ctx, "Graceful shutdown completed successfully")
			}
		case <-shutdownCtx.Done():
			// Shutdown timed out
			logger.Error(ctx, "Shutdown function timed out after 30 seconds")
			shutdownErr = fmt.Errorf("shutdown function timed out after 30 seconds")
		case sig := <-quit:
			// Received second signal during shutdown
			logger.Warn(ctx, "Received second termination signal during shutdown, forcing exit",
				zap.String("signal", sig.String()),
				zap.String("type", "second"))
			os.Exit(1)
		}
	}()

	// Wait for shutdown to complete
	<-done
	return shutdownErr
}

// SetupGracefulShutdown sets up a goroutine that will handle graceful shutdown.
// It creates a new context with cancellation and starts a background goroutine
// that calls GracefulShutdown. This allows for both signal-based and programmatic
// shutdown initiation.
//
// Parameters:
//   - ctx: Parent context for the shutdown context
//   - logger: Logger for recording shutdown events
//   - shutdownFunc: Function to execute during shutdown
//
// Returns:
//   - A cancel function that can be called to trigger shutdown programmatically
//   - A channel that will receive any error that occurs during shutdown
func SetupGracefulShutdown(ctx context.Context, logger *logging.ContextLogger, shutdownFunc func() error) (context.CancelFunc, <-chan error) {
	// Create a context with cancellation
	shutdownCtx, cancel := context.WithCancel(ctx)

	// Create a channel to receive shutdown errors
	// Buffer size 1 ensures we don't block if no one is listening
	errCh := make(chan error, 1)

	// Start a goroutine to handle shutdown
	go func() {
		// Execute the graceful shutdown
		err := GracefulShutdown(shutdownCtx, logger, shutdownFunc)

		// Send the error to the channel (nil if no error)
		select {
		case errCh <- err:
			// Error sent successfully
		default:
			// No one is listening, just log the error
			if err != nil {
				logger.Error(ctx, "Error during background graceful shutdown (no listener)", zap.Error(err))
			}
		}

		// Close the error channel to signal completion
		close(errCh)

		// Ensure context is cancelled when shutdown is complete
		cancel()
	}()

	return cancel, errCh
}
