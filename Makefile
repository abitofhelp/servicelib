# Makefile for github.com/abitofhelp/servicelib

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOGENERATE=$(GOCMD) generate

# Platform-specific parameters
# macOS (Darwin)
DARWIN_AMD64=GOOS=darwin GOARCH=amd64  # Intel
DARWIN_ARM64=GOOS=darwin GOARCH=arm64  # Apple Silicon
# Windows
WINDOWS_AMD64=GOOS=windows GOARCH=amd64
# Linux
LINUX_AMD64=GOOS=linux GOARCH=amd64
LINUX_ARM64=GOOS=linux GOARCH=arm64

# This is a library project, so no binary is produced

# Test parameters
TEST_FLAGS=-race -coverprofile=coverage.out
TEST_PACKAGES=./...
INTEGRATION_TEST_FLAGS=-race -coverprofile=integration_coverage.out -tags=integration
PACKAGE_COVERAGE_FLAGS=-race -coverprofile=

# Linter parameters
GOLANGCI_LINT=golangci-lint

.PHONY: all build build-darwin-amd64 build-darwin-arm64 build-windows-amd64 build-linux-amd64 build-linux-arm64 build-all clean test test-integration test-all test-package coverage coverage-integration coverage-package lint fmt vet tidy vendor generate generate-package security bench bench-package help

all: clean lint test build

# Build the project (compile without producing a binary since this is a library)
build:
	@echo "Building for current platform..."
	$(GOBUILD) -v ./...

# Build for macOS Intel
build-darwin-amd64:
	@echo "Building for macOS Intel (amd64)..."
	$(DARWIN_AMD64) $(GOBUILD) -v ./...

# Build for macOS Apple Silicon
build-darwin-arm64:
	@echo "Building for macOS Apple Silicon (arm64)..."
	$(DARWIN_ARM64) $(GOBUILD) -v ./...

# Build for Windows
build-windows-amd64:
	@echo "Building for Windows (amd64)..."
	$(WINDOWS_AMD64) $(GOBUILD) -v ./...

# Build for Linux amd64
build-linux-amd64:
	@echo "Building for Linux (amd64)..."
	$(LINUX_AMD64) $(GOBUILD) -v ./...

# Build for Linux arm64
build-linux-arm64:
	@echo "Building for Linux (arm64)..."
	$(LINUX_ARM64) $(GOBUILD) -v ./...

# Build for all platforms
build-all: build-darwin-amd64 build-darwin-arm64 build-windows-amd64 build-linux-amd64 build-linux-arm64
	@echo "Building for all platforms completed."

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f coverage.out coverage.html integration_coverage.out integration_coverage.html unit_coverage.out all_coverage.out all_coverage.html *_coverage.out *_coverage.html
	$(GOCLEAN)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) $(TEST_FLAGS) $(TEST_PACKAGES)

# Run integration tests (with integration build tag)
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) $(INTEGRATION_TEST_FLAGS) $(TEST_PACKAGES)

# Run tests for a specific package
# Usage: make test-package PACKAGE=./path/to/package
test-package:
	@if [ -z "$(PACKAGE)" ]; then \
		echo "Error: PACKAGE parameter is required. Usage: make test-package PACKAGE=./path/to/package"; \
		exit 1; \
	fi
	@echo "Running tests for package $(PACKAGE)..."
	$(GOTEST) $(TEST_FLAGS) $(PACKAGE)

# Generate test coverage report
coverage: test
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Generate test coverage report for integration tests
coverage-integration: test-integration
	@echo "Generating integration test coverage report..."
	$(GOCMD) tool cover -html=integration_coverage.out -o integration_coverage.html

# Run all tests (unit and integration) and generate a combined coverage report
test-all:
	@echo "Running all tests (unit and integration)..."
	$(GOTEST) $(TEST_FLAGS) $(TEST_PACKAGES)
	@mv coverage.out unit_coverage.out
	$(GOTEST) $(INTEGRATION_TEST_FLAGS) $(TEST_PACKAGES)
	@echo "Merging coverage reports..."
	@if command -v gocovmerge > /dev/null; then \
		gocovmerge unit_coverage.out integration_coverage.out > all_coverage.out; \
	else \
		echo "gocovmerge not found, installing..."; \
		go install github.com/wadey/gocovmerge@latest; \
		gocovmerge unit_coverage.out integration_coverage.out > all_coverage.out; \
	fi
	@echo "Generating combined coverage report..."
	$(GOCMD) tool cover -html=all_coverage.out -o all_coverage.html
	@echo "Combined coverage report generated: all_coverage.html"

# Generate test coverage report for a specific package
# Usage: make coverage-package PACKAGE=./path/to/package OUTPUT=package_coverage
coverage-package:
	@if [ -z "$(PACKAGE)" ]; then \
		echo "Error: PACKAGE parameter is required. Usage: make coverage-package PACKAGE=./path/to/package OUTPUT=package_coverage"; \
		exit 1; \
	fi
	@if [ -z "$(OUTPUT)" ]; then \
		echo "Error: OUTPUT parameter is required. Usage: make coverage-package PACKAGE=./path/to/package OUTPUT=package_coverage"; \
		exit 1; \
	fi
	@echo "Running tests for package $(PACKAGE) with coverage..."
	$(GOTEST) $(PACKAGE_COVERAGE_FLAGS)$(OUTPUT).out $(PACKAGE)
	@echo "Generating coverage report for $(PACKAGE)..."
	$(GOCMD) tool cover -html=$(OUTPUT).out -o $(OUTPUT).html

# Run linter
lint:
	@echo "Running linter..."
	@if command -v $(GOLANGCI_LINT) > /dev/null; then \
		$(GOLANGCI_LINT) run; \
	else \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		$(GOLANGCI_LINT) run; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Update dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

# Vendor dependencies
vendor:
	@echo "Vendoring dependencies..."
	$(GOMOD) vendor

# Generate mocks
generate:
	@echo "Generating mocks..."
	$(GOGENERATE) ./...

# Generate mocks for a specific package
# Usage: make generate-package PACKAGE=./path/to/package
generate-package:
	@if [ -z "$(PACKAGE)" ]; then \
		echo "Error: PACKAGE parameter is required. Usage: make generate-package PACKAGE=./path/to/package"; \
		exit 1; \
	fi
	@echo "Generating mocks for package $(PACKAGE)..."
	$(GOGENERATE) $(PACKAGE)

# Check for security vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	@if command -v govulncheck > /dev/null; then \
		govulncheck ./...; \
	else \
		echo "govulncheck not found, installing..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
		govulncheck ./...; \
	fi

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Run benchmarks for a specific package
# Usage: make bench-package PACKAGE=./path/to/package
bench-package:
	@if [ -z "$(PACKAGE)" ]; then \
		echo "Error: PACKAGE parameter is required. Usage: make bench-package PACKAGE=./path/to/package"; \
		exit 1; \
	fi
	@echo "Running benchmarks for package $(PACKAGE)..."
	$(GOTEST) -bench=. -benchmem $(PACKAGE)

# Help target
help:
	@echo "Available targets:"
	@echo "  all                - Run clean, lint, test, and build"
	@echo "  build              - Compile the code for current platform (no binary produced)"
	@echo "  build-darwin-amd64 - Compile for macOS Intel (amd64)"
	@echo "  build-darwin-arm64 - Compile for macOS Apple Silicon (arm64)"
	@echo "  build-windows-amd64 - Compile for Windows (amd64)"
	@echo "  build-linux-amd64  - Compile for Linux (amd64)"
	@echo "  build-linux-arm64  - Compile for Linux (arm64)"
	@echo "  build-all          - Compile for all platforms"
	@echo "  clean              - Remove test artifacts and clean Go cache"
	@echo "  test               - Run tests"
	@echo "  test-integration   - Run integration tests (with integration build tag)"
	@echo "  test-all           - Run all tests (unit and integration) with combined coverage report"
	@echo "  test-package       - Run tests for a specific package (PACKAGE=./path/to/package)"
	@echo "  coverage           - Generate test coverage report"
	@echo "  coverage-integration - Generate test coverage report for integration tests"
	@echo "  coverage-package   - Generate test coverage report for a specific package (PACKAGE=./path/to/package OUTPUT=name)"
	@echo "  lint               - Run linter"
	@echo "  fmt                - Format code"
	@echo "  vet                - Run go vet"
	@echo "  tidy               - Update dependencies"
	@echo "  vendor             - Vendor dependencies"
	@echo "  generate           - Generate mocks using go:generate"
	@echo "  generate-package   - Generate mocks for a specific package (PACKAGE=./path/to/package)"
	@echo "  security           - Check for security vulnerabilities"
	@echo "  bench              - Run benchmarks"
	@echo "  bench-package      - Run benchmarks for a specific package (PACKAGE=./path/to/package)"
	@echo "  help               - Show this help message"

# Default target
.DEFAULT_GOAL := help
