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

// ShutdownCallback is a function that is called during the shutdown process.
// It receives a context that may include a deadline based on the shutdown timeout,
// and should return an error if the shutdown operation fails.
//
// The context passed to the callback will be canceled if the shutdown timeout is reached,
// allowing the callback to implement its own timeout handling.
type ShutdownCallback func(ctx context.Context) error

// GracefulShutdown represents a graceful shutdown handler that manages
// the process of shutting down an application in response to OS signals.
//
// It provides mechanisms for registering shutdown callbacks, handling
// termination signals, and ensuring that shutdown procedures complete
// within a specified timeout.
type GracefulShutdown struct {
	timeout    time.Duration     // Maximum time to wait for shutdown callbacks
	logger     *logging.ContextLogger // Logger for shutdown events
	callbacks  []ShutdownCallback     // Registered shutdown callbacks
	callbackMu sync.Mutex             // Mutex for thread-safe callback registration
	signals    []os.Signal            // OS signals to handle
}

// NewGracefulShutdown creates a new graceful shutdown handler.
//
// The handler is configured to listen for standard termination signals:
// SIGINT (Ctrl+C), SIGTERM (kill), SIGHUP (terminal closed), and SIGQUIT.
// When any of these signals is received, the handler will execute all registered
// shutdown callbacks concurrently and wait for them to complete, up to the specified timeout.
//
// Parameters:
//   - timeout: The maximum time to wait for shutdown callbacks to complete.
//     If callbacks take longer than this duration, they will be interrupted.
//   - logger: The logger to use for logging shutdown events and progress.
//     This should be a properly initialized ContextLogger.
//
// Returns:
//   - *GracefulShutdown: A new graceful shutdown handler ready to be used.
//     Use RegisterCallback to add shutdown procedures and HandleShutdown to
//     start listening for signals.
func NewGracefulShutdown(timeout time.Duration, logger *logging.ContextLogger) *GracefulShutdown {
	return &GracefulShutdown{
		timeout:   timeout,
		logger:    logger,
		callbacks: make([]ShutdownCallback, 0),
		signals:   []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT},
	}
}

// RegisterCallback registers a callback to be called during the shutdown process.
//
// Callbacks are executed concurrently when a shutdown signal is received.
// There is no guaranteed order of execution between callbacks, so they should be
// designed to operate independently of each other.
//
// This method is thread-safe and can be called from multiple goroutines.
//
// Parameters:
//   - callback: The function to call during shutdown. This function should
//     perform any necessary cleanup operations and return an error if the
//     cleanup fails. The function should respect the context's deadline
//     and cancel itself if the context is canceled.
func (gs *GracefulShutdown) RegisterCallback(callback ShutdownCallback) {
	gs.callbackMu.Lock()
	defer gs.callbackMu.Unlock()
	gs.callbacks = append(gs.callbacks, callback)
}

// HandleShutdown sets up signal handling for graceful shutdown.
//
// This method starts a goroutine that listens for OS signals and initiates
// the shutdown process when a signal is received. The shutdown process includes:
//   1. Canceling the returned context to notify the application
//   2. Creating a timeout context for shutdown callbacks
//   3. Executing all registered callbacks concurrently
//   4. Waiting for callbacks to complete or timeout
//   5. Logging the shutdown progress
//   6. Exiting the application when complete
//
// The method handles multiple signals appropriately - if a second signal is
// received while shutdown is in progress, it forces an immediate exit.
//
// Returns:
//   - context.Context: A context that will be canceled when a shutdown signal is received.
//     The application should monitor this context and begin its shutdown procedure when canceled.
//   - context.CancelFunc: A function that can be called to manually trigger the shutdown process
//     without waiting for an OS signal. This is useful for testing or for implementing
//     application-specific shutdown triggers.
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

// WaitForShutdown sets up signal handling and returns a context that will be canceled
// when a shutdown signal is received.
//
// This is a convenience function that creates a GracefulShutdown instance and calls
// HandleShutdown, but discards the cancel function. It's useful when you only need
// to know when a shutdown signal is received, but don't need to register callbacks
// or manually trigger shutdown.
//
// Note that this function does not block - it returns immediately. To wait for a
// shutdown signal, you should monitor the returned context:
//
//	ctx := signal.WaitForShutdown(30*time.Second, logger)
//	// Start your application...
//	<-ctx.Done() // This will block until a shutdown signal is received
//	// Perform shutdown procedures...
//
// Parameters:
//   - timeout: The maximum time to wait for any internal shutdown callbacks to complete.
//     This is primarily used for logging and doesn't affect application shutdown.
//   - logger: The logger to use for logging shutdown events.
//
// Returns:
//   - context.Context: A context that will be canceled when a shutdown signal is received.
//     The application should monitor this context and begin its shutdown procedure when canceled.
func WaitForShutdown(timeout time.Duration, logger *logging.ContextLogger) context.Context {
	gs := NewGracefulShutdown(timeout, logger)
	ctx, _ := gs.HandleShutdown()
	return ctx
}

// SetupSignalHandler sets up a signal handler for graceful shutdown and returns
// both a context and the GracefulShutdown instance.
//
// This is the recommended function to use when you need both notification of
// shutdown signals and the ability to register custom shutdown callbacks.
// It creates a GracefulShutdown instance, calls HandleShutdown, and returns
// both the resulting context and the GracefulShutdown instance.
//
// Typical usage:
//
//	ctx, gs := signal.SetupSignalHandler(30*time.Second, logger)
//	
//	// Register shutdown callbacks
//	gs.RegisterCallback(func(ctx context.Context) error {
//	    return db.Close()
//	})
//	
//	// Start your application...
//	
//	// Wait for shutdown signal
//	<-ctx.Done()
//	
//	// Perform any additional shutdown procedures not registered as callbacks...
//
// Parameters:
//   - timeout: The maximum time to wait for shutdown callbacks to complete.
//     If callbacks take longer than this duration, they will be interrupted.
//   - logger: The logger to use for logging shutdown events and progress.
//
// Returns:
//   - context.Context: A context that will be canceled when a shutdown signal is received.
//     The application should monitor this context and begin its shutdown procedure when canceled.
//   - *GracefulShutdown: The graceful shutdown handler that can be used to register
//     additional callbacks to be executed during shutdown.
func SetupSignalHandler(timeout time.Duration, logger *logging.ContextLogger) (context.Context, *GracefulShutdown) {
	gs := NewGracefulShutdown(timeout, logger)
	ctx, _ := gs.HandleShutdown()
	return ctx, gs
}
