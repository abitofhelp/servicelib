// Copyright (c) 2025 A Bit of Help, Inc.

// Package signal provides utilities for handling OS signals and implementing graceful shutdown.
//
// This package helps applications respond properly to operating system signals (such as
// SIGINT, SIGTERM) by providing a structured way to handle shutdown procedures. It ensures
// that applications can clean up resources, finish in-flight operations, and terminate
// gracefully when signaled to stop.
//
// Key features:
//   - Signal handling for common termination signals (SIGINT, SIGTERM, SIGHUP, SIGQUIT)
//   - Graceful shutdown with configurable timeout
//   - Support for registering multiple shutdown callbacks
//   - Concurrent execution of shutdown callbacks
//   - Proper handling of multiple signals (force exit on second signal)
//   - Context-based notification system for shutdown events
//
// The package is designed around the GracefulShutdown type, which manages the shutdown
// process and provides methods for registering callbacks to be executed during shutdown.
//
// Example usage:
//
//	// Create a logger
//	logger := logging.NewContextLogger("my-service")
//
//	// Set up graceful shutdown with a 30-second timeout
//	ctx, gs := signal.SetupSignalHandler(30*time.Second, logger)
//
//	// Register shutdown callbacks
//	gs.RegisterCallback(func(ctx context.Context) error {
//	    logger.Info(ctx, "Closing database connection...")
//	    return db.Close()
//	})
//
//	gs.RegisterCallback(func(ctx context.Context) error {
//	    logger.Info(ctx, "Stopping HTTP server...")
//	    return server.Shutdown(ctx)
//	})
//
//	// Start your application
//	server := startServer()
//
//	// Wait for shutdown signal
//	<-ctx.Done()
//	logger.Info(ctx, "Shutdown signal received, stopping application...")
//
//	// The registered callbacks will be executed automatically
//	// when a shutdown signal is received
//
// The package also provides simpler helper functions for common use cases:
//
//	// Just wait for a shutdown signal
//	ctx := signal.WaitForShutdown(30*time.Second, logger)
//
//	// Start your application
//	server := startServer()
//
//	// Wait for shutdown signal
//	<-ctx.Done()
//	logger.Info(ctx, "Shutdown signal received, stopping application...")
//
//	// Perform shutdown manually
//	server.Shutdown(context.Background())
package signal
