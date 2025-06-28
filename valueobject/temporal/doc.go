// Copyright (c) 2025 A Bit of Help, Inc.

// Package temporal provides immutable value objects for representing and manipulating
// temporal concepts such as dates, times, durations, and versions.
//
// This package implements various temporal value objects that follow the Value Object
// pattern from Domain-Driven Design (DDD). All temporal value objects are immutable,
// providing type safety, built-in validation, and comparison operations.
//
// Key features:
//   - Immutable values: All temporal value objects are immutable
//   - Type safety: Strong typing for temporal concepts
//   - Validation: Built-in validation for temporal values
//   - Comparison operations: Methods for comparing temporal values
//   - Formatting: Methods for formatting temporal values
//
// The package provides several types of temporal value objects:
//
//   - Duration: Represents a duration of time
//     Example:
//     duration, err := temporal.NewDuration(5 * time.Minute)
//
//   - Time: Represents a specific point in time
//     Example:
//     timeVO, err := temporal.NewTime(time.Now())
//
//   - Version: Represents a semantic version (major.minor.patch)
//     Example:
//     version, err := temporal.NewVersion(1, 2, 3)
//
// Best practices when working with temporal value objects:
//  1. Never modify temporal value objects; always create new ones
//  2. Always validate input when creating temporal value objects
//  3. Use the provided comparison methods instead of comparing fields directly
//  4. Be aware of time zone implications when working with time value objects
//  5. Use the provided formatting methods for consistent output
//
// Example usage:
//
//	// Create a duration value object
//	duration, err := temporal.NewDuration(30 * time.Second)
//	if err != nil {
//	    // Handle validation error
//	}
//
//	// Create a time value object
//	timeVO, err := temporal.NewTime(time.Now())
//	if err != nil {
//	    // Handle validation error
//	}
//
//	// Create a version value object
//	version, err := temporal.NewVersion(2, 1, 0)
//	if err != nil {
//	    // Handle validation error
//	}
//
//	// Compare versions
//	otherVersion, _ := temporal.NewVersion(2, 0, 0)
//	if version.IsGreaterThan(otherVersion) {
//	    fmt.Println("Version is newer")
//	}
//
// All temporal value objects implement the base.ValueObject interface and provide
// additional methods specific to their temporal nature.
package temporal
