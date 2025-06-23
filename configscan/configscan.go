// Copyright (c) 2025 A Bit of Help, Inc.

// Package configscan provides functionality to scan packages for configuration requirements
// and check if they have default values set.
package configscan

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ConfigRequirement represents a configuration requirement for a package
type ConfigRequirement struct {
	PackageName string   // Name of the package
	ConfigType  string   // Name of the configuration struct
	Fields      []string // Fields of the configuration struct
	HasDefault  bool     // Whether the package provides default values for the configuration
}

// Scanner scans packages for configuration requirements
type Scanner struct {
	rootDir string // Root directory to scan
}

// NewScanner creates a new Scanner
func NewScanner(rootDir string) *Scanner {
	return &Scanner{
		rootDir: rootDir,
	}
}

// Scan scans packages for configuration requirements
func (s *Scanner) Scan() ([]ConfigRequirement, error) {
	var requirements []ConfigRequirement

	// Walk through all directories in the root directory
	err := filepath.Walk(s.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files and directories
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		// Parse the Go file
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse file %s: %w", path, err)
		}

		// Extract package name
		packageName := file.Name.Name

		// Find configuration structs and check if they have default values
		for _, decl := range file.Decls {
			// Check if it's a type declaration
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				// Check if it's a struct type
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				// Check if the struct name contains "Config"
				if !strings.Contains(typeSpec.Name.Name, "Config") {
					continue
				}

				// Extract fields from the struct
				var fields []string
				if structType.Fields != nil {
					for _, field := range structType.Fields.List {
						for _, name := range field.Names {
							fields = append(fields, name.Name)
						}
					}
				}

				// Check if there's a default function for this config type
				hasDefault := s.hasDefaultFunction(file, typeSpec.Name.Name)

				// Add the configuration requirement
				requirements = append(requirements, ConfigRequirement{
					PackageName: packageName,
					ConfigType:  typeSpec.Name.Name,
					Fields:      fields,
					HasDefault:  hasDefault,
				})
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan packages: %w", err)
	}

	return requirements, nil
}

// hasDefaultFunction checks if there's a default function for a config type
func (s *Scanner) hasDefaultFunction(file *ast.File, configType string) bool {
	for _, decl := range file.Decls {
		// Check if it's a function declaration
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Check if the function name is "Default" + configType or "Default" + strings.TrimSuffix(configType, "Config") + "Config"
		funcName := funcDecl.Name.Name
		if funcName == "Default"+configType || funcName == "Default"+strings.TrimSuffix(configType, "Config")+"Config" {
			return true
		}
	}

	return false
}

// Report generates a report of packages with configuration requirements that don't have default values set
func Report(requirements []ConfigRequirement) []ConfigRequirement {
	var missingDefaults []ConfigRequirement

	for _, req := range requirements {
		if !req.HasDefault {
			missingDefaults = append(missingDefaults, req)
		}
	}

	return missingDefaults
}