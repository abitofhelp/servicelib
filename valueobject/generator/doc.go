// Copyright (c) 2025 A Bit of Help, Inc.

// Package generator provides code generation tools for value objects.
//
// This package contains utilities for generating value object code based on
// templates and configuration files. It is designed to automate the creation
// of value objects that follow the Value Object pattern from Domain-Driven Design.
//
// The package provides:
//   - A Generator type for generating value object code and tests
//   - Configuration types for specifying value object properties
//   - Command-line interface for using the generator from the terminal
//
// The generator supports two types of value objects:
//   - String-based value objects (wrapping a string type)
//   - Struct-based value objects (containing multiple fields)
//
// Configuration is provided through JSON files that specify:
//   - The name and package of the value object
//   - The type of value object (string or struct)
//   - Fields and their types (for struct-based value objects)
//   - Validation rules
//   - Required imports
//   - Description of the value object
//
// Example usage:
//
//	// Create a generator
//	generator := generator.NewGenerator("templates", "output")
//
//	// Configure a value object
//	config := generator.ValueObjectConfig{
//	    Name:        "Email",
//	    Package:     "contact",
//	    Type:        generator.StringBased,
//	    BaseType:    "string",
//	    Description: "Represents an email address",
//	    Validations: map[string]string{
//	        "value": "ValidateEmail",
//	    },
//	}
//
//	// Generate the value object code
//	err := generator.Generate(config)
//	if err != nil {
//	    // Handle error
//	}
//
// The package also provides a command-line interface through the Execute function,
// which is used by the valueobject/cmd/generate package to provide a standalone
// command-line tool.
package generator
