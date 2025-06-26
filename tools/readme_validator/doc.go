// Copyright (c) 2025 A Bit of Help, Inc.

// Package main provides a utility for validating README.md files against templates.
//
// This tool checks README.md files in the project to ensure they conform to the
// project's documentation standards and templates. Unlike the readme_fixer tool,
// which attempts to fix issues automatically, this tool focuses on validation and
// reporting, making it suitable for use in CI/CD pipelines where failures should
// prevent merges of non-compliant documentation.
//
// The readme_validator tool helps maintain high-quality documentation by enforcing
// consistent structure, required sections, and formatting standards across all
// README files in the project. This ensures that users and developers can find
// information easily and that all components are documented to the same standard.
//
// Key features:
//   - Strict validation against template requirements
//   - Detailed error reporting with line numbers and context
//   - Support for different templates based on component type
//   - Exit codes for integration with CI/CD pipelines
//
// This tool is designed to be run both locally during development and as part of
// automated checks in the CI/CD pipeline to prevent merging documentation that
// doesn't meet the project's standards.
//
// Usage:
//
//	go run main.go [flags]
//
// Flags:
//   -dir string       Directory to scan for README.md files (default ".")
//   -template string  Path to the template file (default "README_TEMPLATE.md")
//   -verbose          Enable verbose output (default false)
package main