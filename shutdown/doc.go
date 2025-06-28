// Copyright (c) 2025 A Bit of Help, Inc.

// Package shutdown provides functionality for graceful application shutdown.
//
// This package implements a robust system for handling application termination,
// ensuring that resources are properly released and in-flight operations are
// completed before the application exits. It handles OS signals (SIGINT, SIGTERM, SIGHUP)
// and context cancellation to trigger graceful shutdown.
//
// Proper shutdown handling is critical for production applications to prevent:
//   - Data loss from incomplete operations
//   - Resource leaks from unclosed connections
//   - Inconsistent state from abrupt termination
//   - Service disruption for users during deployments
//
// Key features:
//   - Signal-based shutdown handling (SIGINT, SIGTERM, SIGHUP)
//   - Context-based shutdown initiation for programmatic control
//   - Timeout management to prevent hanging during shutdown
//   - Multiple signal handling with forced exit on second signal
//   - Comprehensive logging of shutdown events
//   - Error propagation from shutdown operations
//
// The package provides two main functions:
//   - GracefulShutdown: Blocks until shutdown is triggered, then executes cleanup
//   - SetupGracefulShutdown: Sets up background shutdown handling without blocking
//
// Example usage with GracefulShutdown (blocking approach):
//
//	func main() {
//	    // Initialize application components
//	    logger := logging.NewContextLogger(zapLogger)
//	    server := startServer()
//	    db := connectToDatabase()
//
//	    // Define shutdown function
//	    shutdownFunc := func() error {
//	        // Close resources in reverse order of creation
//	        serverErr := server.Shutdown(context.Background())
//	        dbErr := db.Close()
//
//	        // Return combined error if any
//	        if serverErr != nil || dbErr != nil {
//	            return fmt.Errorf("shutdown errors: server=%v, db=%v", serverErr, dbErr)
//	        }
//	        return nil
//	    }
//
//	    // Wait for shutdown signal and execute cleanup
//	    if err := shutdown.GracefulShutdown(context.Background(), logger, shutdownFunc); err != nil {
//	        logger.Error(context.Background(), "Shutdown completed with errors", zap.Error(err))
//	        os.Exit(1)
//	    }
//	}
//
// Example usage with SetupGracefulShutdown (non-blocking approach):
//
//	func main() {
//	    // Initialize application components
//	    logger := logging.NewContextLogger(zapLogger)
//	    server := startServer()
//	    db := connectToDatabase()
//
//	    // Define shutdown function
//	    shutdownFunc := func() error {
//	        // Close resources in reverse order of creation
//	        serverErr := server.Shutdown(context.Background())
//	        dbErr := db.Close()
//
//	        // Return combined error if any
//	        if serverErr != nil || dbErr != nil {
//	            return fmt.Errorf("shutdown errors: server=%v, db=%v", serverErr, dbErr)
//	        }
//	        return nil
//	    }
//
//	    // Set up graceful shutdown in the background
//	    cancel, errCh := shutdown.SetupGracefulShutdown(context.Background(), logger, shutdownFunc)
//	    defer cancel() // Ensure shutdown is triggered if main exits
//
//	    // Continue with application logic
//	    // ...
//
//	    // Optionally wait for shutdown to complete and check for errors
//	    if err := <-errCh; err != nil {
//	        logger.Error(context.Background(), "Shutdown completed with errors", zap.Error(err))
//	        os.Exit(1)
//	    }
//	}
//
// The package is designed to be used at the application's entry point (main function)
// to ensure proper resource cleanup during termination.
package shutdown
