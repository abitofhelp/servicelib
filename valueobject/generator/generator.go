// Copyright (c) 2025 A Bit of Help, Inc.

// Package generator provides code generation tools for value objects.
package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// ValueObjectType represents the type of value object to generate.
type ValueObjectType string

const (
	// StringBased represents a string-based value object.
	StringBased ValueObjectType = "string"
	// StructBased represents a struct-based value object.
	StructBased ValueObjectType = "struct"
)

// ValueObjectConfig contains the configuration for generating a value object.
type ValueObjectConfig struct {
	// Name is the name of the value object (e.g., "Email").
	Name string
	// Package is the package name (e.g., "contact").
	Package string
	// Type is the type of value object (string or struct).
	Type ValueObjectType
	// Fields is a map of field names to field types (for struct-based value objects).
	Fields map[string]string
	// Validations is a map of field names to validation functions.
	Validations map[string]string
	// Imports is a list of additional imports.
	Imports []string
	// BaseType is the base type for string-based value objects (e.g., "string").
	BaseType string
	// Description is a description of the value object.
	Description string
}

// Generator generates code for value objects.
type Generator struct {
	// TemplateDir is the directory containing the templates.
	TemplateDir string
	// OutputDir is the directory where the generated code will be written.
	OutputDir string
}

// NewGenerator creates a new Generator.
func NewGenerator(templateDir, outputDir string) *Generator {
	return &Generator{
		TemplateDir: templateDir,
		OutputDir:   outputDir,
	}
}

// Generate generates code for a value object.
func (g *Generator) Generate(config ValueObjectConfig) error {
	// Create the output directory if it doesn't exist
	outputDir := filepath.Join(g.OutputDir, config.Package)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate the value object code
	code, err := g.generateValueObject(config)
	if err != nil {
		return fmt.Errorf("failed to generate value object code: %w", err)
	}

	// Write the code to a file
	outputFile := filepath.Join(outputDir, strings.ToLower(config.Name)+".go")
	if err := os.WriteFile(outputFile, code, 0644); err != nil {
		return fmt.Errorf("failed to write value object code: %w", err)
	}

	// Generate the test code
	testCode, err := g.generateValueObjectTest(config)
	if err != nil {
		return fmt.Errorf("failed to generate value object test code: %w", err)
	}

	// Write the test code to a file
	testOutputFile := filepath.Join(outputDir, strings.ToLower(config.Name)+"_test.go")
	if err := os.WriteFile(testOutputFile, testCode, 0644); err != nil {
		return fmt.Errorf("failed to write value object test code: %w", err)
	}

	return nil
}

// createTemplateFuncs creates the template functions.
func createTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// ToLower converts a string to lowercase.
		"ToLower": strings.ToLower,
		// Keys returns the keys of a map.
		"Keys": func(m map[string]string) []string {
			keys := make([]string, 0, len(m))
			for k := range m {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			return keys
		},
		// First returns the first element of a slice.
		"First": func(s []string) string {
			if len(s) == 0 {
				return ""
			}
			return s[0]
		},
	}
}

// generateValueObject generates code for a value object.
func (g *Generator) generateValueObject(config ValueObjectConfig) ([]byte, error) {
	var templateFile string
	switch config.Type {
	case StringBased:
		templateFile = filepath.Join(g.TemplateDir, "string_value_object.tmpl")
	case StructBased:
		templateFile = filepath.Join(g.TemplateDir, "struct_value_object.tmpl")
	default:
		return nil, fmt.Errorf("unsupported value object type: %s", config.Type)
	}

	// Create a new template with functions
	tmpl := template.New(filepath.Base(templateFile)).Funcs(createTemplateFuncs())

	// Parse the template
	tmpl, err := tmpl.ParseFiles(templateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Format the code
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format code: %w", err)
	}

	return formattedCode, nil
}

// generateValueObjectTest generates test code for a value object.
func (g *Generator) generateValueObjectTest(config ValueObjectConfig) ([]byte, error) {
	var templateFile string
	switch config.Type {
	case StringBased:
		templateFile = filepath.Join(g.TemplateDir, "string_value_object_test.tmpl")
	case StructBased:
		templateFile = filepath.Join(g.TemplateDir, "struct_value_object_test.tmpl")
	default:
		return nil, fmt.Errorf("unsupported value object type: %s", config.Type)
	}

	// Create a new template with functions
	tmpl := template.New(filepath.Base(templateFile)).Funcs(createTemplateFuncs())

	// Parse the template
	tmpl, err := tmpl.ParseFiles(templateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Format the code
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format code: %w", err)
	}

	return formattedCode, nil
}
