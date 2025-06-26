// Copyright (c) 2025 A Bit of Help, Inc.

// Package main provides a utility for fixing and standardizing README.md files.
//
// This tool scans README.md files in the project and ensures they follow the
// project's standardized format and structure. It can automatically fix common
// issues such as missing sections, incorrect formatting, and inconsistent styling.
//
// The readme_fixer tool helps maintain consistency across all documentation in the
// project, making it easier for developers and users to find information. It enforces
// the documentation standards defined in the project's templates.
//
// Key features:
//   - Automatic detection of README.md files
//   - Validation against template structure
//   - Automatic fixing of common issues
//   - Reporting of issues that require manual intervention
//
// This tool is typically run as part of the CI/CD pipeline to ensure all
// documentation remains consistent as the project evolves.
//
// Usage:
//
//	go run main.go [flags]
//
// Flags:
//   -dir string    Directory to scan for README.md files (default ".")
//   -fix           Automatically fix issues (default false)
//   -verbose       Enable verbose output (default false)
package main