// Copyright (c) 2025 A Bit of Help, Inc.

// Package signal provides utilities for handling OS signals and graceful shutdown.
package signal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// ShutdownCallback is a function that is called during shutdown
type ShutdownCallback func(ctx context.Context) error

// GracefulShutdown represents a graceful shutdown handler
type GracefulShutdown struct {
	timeout    time.Duration
	logger     *logging.ContextLogger
	callbacks  []ShutdownCallback
	callbackMu sync.Mutex
	signals    []os.Signal
}

// NewGracefulShutdown creates a new graceful shutdown handler
//
// Parameters:
//   - timeout: The maximum time to wait for shutdown callbacks to complete
//   - logger: The logger to use for logging shutdown events
//
// Returns:
//   - *GracefulShutdown: A new graceful shutdown handler
func NewGracefulShutdown(timeout time.Duration, logger *logging.ContextLogger) *GracefulShutdown {
	return &GracefulShutdown{
		timeout:   timeout,
		logger:    logger,
		callbacks: make([]ShutdownCallback, 0),
		signals:   []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT},
	}
}

// RegisterCallback registers a callback to be called during shutdown
//
// Parameters:
//   - callback: The function to call during shutdown
func (gs *GracefulShutdown) RegisterCallback(callback ShutdownCallback) {
	gs.callbackMu.Lock()
	defer gs.callbackMu.Unlock()
	gs.callbacks = append(gs.callbacks, callback)
}

// HandleShutdown handles graceful shutdown
// It returns a context that will be canceled when a shutdown signal is received
// and a cancel function that can be called to cancel the context manually
//
// Returns:
//   - context.Context: A context that will be canceled when a shutdown signal is received
//   - context.CancelFunc: A function that can be called to cancel the context manually
func (gs *GracefulShutdown) HandleShutdown() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// Create a channel to receive OS signals
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, gs.signals...)

		// Wait for a signal
		sig := <-sigCh
		gs.logger.Info(ctx, "Received shutdown signal", zap.String("signal", sig.String()))

		// Cancel the context to notify all services to shut down
		cancel()

		// Create a timeout context for shutdown
		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), gs.timeout)
		defer timeoutCancel()

		// Execute registered callbacks
		gs.callbackMu.Lock()
		callbacks := make([]ShutdownCallback, len(gs.callbacks))
		copy(callbacks, gs.callbacks)
		gs.callbackMu.Unlock()

		var wg sync.WaitGroup
		for i, callback := range callbacks {
			wg.Add(1)
			go func(i int, cb ShutdownCallback) {
				defer wg.Done()
				gs.logger.Info(ctx, "Executing shutdown callback", zap.Int("callback_index", i))
				if err := cb(timeoutCtx); err != nil {
					gs.logger.Error(ctx, "Shutdown callback failed",
						zap.Int("callback_index", i),
						zap.Error(err))
				} else {
					gs.logger.Info(ctx, "Shutdown callback completed successfully",
						zap.Int("callback_index", i))
				}
			}(i, callback)
		}

		// Wait for callbacks to complete or timeout
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		// Wait for callbacks to complete, timeout, or second signal
		select {
		case <-done:
			gs.logger.Info(ctx, "All shutdown callbacks completed successfully")
		case <-timeoutCtx.Done():
			gs.logger.Warn(ctx, "Graceful shutdown timed out, forcing exit")
		case sig := <-sigCh:
			gs.logger.Warn(ctx, "Received second signal, forcing exit",
				zap.String("signal", sig.String()))
		}

		// Force exit
		os.Exit(0)
	}()

	return ctx, cancel
}

// WaitForShutdown blocks until a shutdown signal is received
// It returns a context that will be canceled when a shutdown signal is received
//
// Parameters:
//   - timeout: The maximum time to wait for shutdown callbacks to complete
//   - logger: The logger to use for logging shutdown events
//
// Returns:
//   - context.Context: A context that will be canceled when a shutdown signal is received
func WaitForShutdown(timeout time.Duration, logger *logging.ContextLogger) context.Context {
	gs := NewGracefulShutdown(timeout, logger)
	ctx, _ := gs.HandleShutdown()
	return ctx
}

// SetupSignalHandler sets up a signal handler for graceful shutdown
// It returns a context that will be canceled when a shutdown signal is received
// and the GracefulShutdown instance for registering callbacks
//
// Parameters:
//   - timeout: The maximum time to wait for shutdown callbacks to complete
//   - logger: The logger to use for logging shutdown events
//
// Returns:
//   - context.Context: A context that will be canceled when a shutdown signal is received
//   - *GracefulShutdown: The graceful shutdown handler for registering callbacks
func SetupSignalHandler(timeout time.Duration, logger *logging.ContextLogger) (context.Context, *GracefulShutdown) {
	gs := NewGracefulShutdown(timeout, logger)
	ctx, _ := gs.HandleShutdown()
	return ctx, gs
}
