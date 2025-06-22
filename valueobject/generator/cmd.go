// Copyright (c) 2025 A Bit of Help, Inc.

// Package generator provides code generation tools for value objects.
package generator

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// GenerateCommand is a command-line tool for generating value objects.
type GenerateCommand struct {
	// ConfigFile is the path to the configuration file.
	ConfigFile string
	// TemplateDir is the directory containing the templates.
	TemplateDir string
	// OutputDir is the directory where the generated code will be written.
	OutputDir string
}

// NewGenerateCommand creates a new GenerateCommand.
func NewGenerateCommand() *GenerateCommand {
	return &GenerateCommand{}
}

// ParseFlags parses the command-line flags.
func (c *GenerateCommand) ParseFlags() {
	flag.StringVar(&c.ConfigFile, "config", "", "path to the configuration file")
	flag.StringVar(&c.TemplateDir, "templates", "templates", "directory containing the templates")
	flag.StringVar(&c.OutputDir, "output", ".", "directory where the generated code will be written")
	flag.Parse()

	if c.ConfigFile == "" {
		fmt.Println("Error: config file is required")
		flag.Usage()
		os.Exit(1)
	}
}

// Run runs the command.
func (c *GenerateCommand) Run() error {
	// Read the configuration file
	configData, err := os.ReadFile(c.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse the configuration
	var configs []ValueObjectConfig
	if err := json.Unmarshal(configData, &configs); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Create the generator
	generator := NewGenerator(c.TemplateDir, c.OutputDir)

	// Generate the value objects
	for _, config := range configs {
		if err := generator.Generate(config); err != nil {
			return fmt.Errorf("failed to generate value object %s: %w", config.Name, err)
		}
		fmt.Printf("Generated value object %s in package %s\n", config.Name, config.Package)
	}

	return nil
}

// Execute executes the command.
func Execute() {
	cmd := NewGenerateCommand()
	cmd.ParseFlags()
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
