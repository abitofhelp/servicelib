// Copyright (c) 2025 A Bit of Help, Inc.

// Package interfaces defines the interfaces for the logging package.
//
// This package provides abstractions for logging functionality that can be implemented
// by different logging backends. It defines the core Logger interface that all logger
// implementations must satisfy, enabling dependency inversion and making it easier to
// swap logging implementations.
//
// The primary interface defined in this package is:
//   - Logger: A context-aware logging interface that provides methods for logging
//     at different levels (Debug, Info, Warn, Error, Fatal) with context information.
//
// This package is designed to be used in conjunction with the parent logging package,
// which provides concrete implementations of these interfaces. By depending on these
// interfaces rather than concrete implementations, application code can remain decoupled
// from specific logging backends.
//
// Example usage:
//
//	// Function that depends on the Logger interface
//	func ProcessData(ctx context.Context, logger interfaces.Logger, data []byte) error {
//	    logger.Info(ctx, "Processing data", zap.Int("bytes", len(data)))
//	    
//	    // Process data...
//	    
//	    if err := validateData(data); err != nil {
//	        logger.Error(ctx, "Data validation failed", zap.Error(err))
//	        return err
//	    }
//	    
//	    logger.Debug(ctx, "Data processed successfully")
//	    return nil
//	}
//
// By using the interfaces defined in this package, application code can be tested
// more easily with mock implementations and can adapt to different logging backends
// without requiring changes to the code that uses the logger.
package interfaces