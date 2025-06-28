// Copyright (c) 2025 A Bit of Help, Inc.

// Package main provides a command-line tool for generating value objects.
//
// This tool automates the creation of value objects based on templates and
// configuration files. It leverages the valueobject/generator package to
// perform the actual generation of code.
//
// Usage:
//
//	go run main.go [flags]
//
// The tool reads configuration from JSON files that specify the properties
// and behavior of the value objects to be generated. It then creates Go files
// with the appropriate code for these value objects, following the Value Object
// pattern from Domain-Driven Design.
//
// This command-line tool is primarily used during development to:
//   - Generate new value object types quickly and consistently
//   - Ensure all value objects follow the same patterns and conventions
//   - Reduce boilerplate code when creating new value objects
//
// For more details on the configuration options and templates, see the
// documentation in the valueobject/generator package.
package main
