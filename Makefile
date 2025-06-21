# Makefile for github.com/abitofhelp/servicelib
# 
# This Makefile provides targets for building, testing, and maintaining the ServiceLib library.
# Run 'make help' to see all available targets.

# ==================================================================================
# Variables
# ==================================================================================

# Project metadata
PROJECT_NAME=servicelib
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GO_VERSION=1.24

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
GOINSTALL=$(GOCMD) install

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
TEST_FLAGS=-coverprofile=coverage.out
TEST_PACKAGES=$(shell go list ./... | grep -v /examples)
INTEGRATION_TEST_FLAGS=-coverprofile=integration_coverage.out -tags=integration
PACKAGE_COVERAGE_FLAGS=-coverprofile=

# Linter parameters
GOLANGCI_LINT=golangci-lint

# Documentation parameters
GODOC=godoc
GODOC_PORT=6060

# Tools
TOOLS_DIR=$(shell pwd)/tools
TOOLS_BIN_DIR=$(TOOLS_DIR)/bin
GOIMPORTS=$(TOOLS_BIN_DIR)/goimports
GOFUMPT=$(TOOLS_BIN_DIR)/gofumpt

.PHONY: all build build-darwin-amd64 build-darwin-arm64 build-windows-amd64 build-linux-amd64 build-linux-arm64 build-all clean test test-integration test-all test-package coverage coverage-integration coverage-package lint fmt vet tidy vendor generate generate-package security bench bench-package help check-go-version tools docs serve-docs ci pre-commit update-deps release

# ==================================================================================
# Main targets
# ==================================================================================

all: check-go-version clean lint test build

# Check Go version
check-go-version:
	@echo "Checking Go version..."
	@if [ "$(shell go version | awk '{print $$3}' | sed 's/go//')" != "$(GO_VERSION)" ]; then \
		echo "ERROR: Required Go version is $(GO_VERSION), but you have $(shell go version | awk '{print $$3}' | sed 's/go//')"; \
		echo "Please update your Go version."; \
		exit 1; \
	fi
	@echo "Go version $(GO_VERSION) confirmed."

# Install required tools
tools:
	@echo "Installing required tools..."
	@mkdir -p $(TOOLS_BIN_DIR)
	@GOBIN=$(TOOLS_BIN_DIR) $(GOINSTALL) golang.org/x/tools/cmd/goimports@latest
	@GOBIN=$(TOOLS_BIN_DIR) $(GOINSTALL) mvdan.cc/gofumpt@latest
	@GOBIN=$(TOOLS_BIN_DIR) $(GOINSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@GOBIN=$(TOOLS_BIN_DIR) $(GOINSTALL) golang.org/x/vuln/cmd/govulncheck@latest
	@GOBIN=$(TOOLS_BIN_DIR) $(GOINSTALL) github.com/wadey/gocovmerge@latest
	@echo "Tools installed successfully."

# ==================================================================================
# Build targets
# ==================================================================================

# Build the project (compile without producing a binary since this is a library)
build: check-go-version
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

# Clean build artifacts and temporary files
clean:
	@echo "Cleaning..."
	@rm -f coverage.out coverage.html integration_coverage.out integration_coverage.html unit_coverage.out all_coverage.out all_coverage.html *_coverage.out *_coverage.html
	@rm -rf $(TOOLS_DIR)
	@find . -type f -name "*.tmp" -delete
	@find . -type d -name "vendor" -prune -o -type f -name ".DS_Store" -delete
	$(GOCLEAN)

# ==================================================================================
# Test targets
# ==================================================================================

# Run tests
test:
	@echo "Running tests..."
	CGO_ENABLED=1 $(GOTEST) $(TEST_FLAGS) $(TEST_PACKAGES)

# Run integration tests (with integration build tag)
test-integration:
	@echo "Running integration tests..."
	CGO_ENABLED=1 $(GOTEST) $(INTEGRATION_TEST_FLAGS) $(TEST_PACKAGES)

# Run tests for a specific package
# Usage: make test-package PACKAGE=./path/to/package
test-package:
	@if [ -z "$(PACKAGE)" ]; then \
		echo "Error: PACKAGE parameter is required. Usage: make test-package PACKAGE=./path/to/package"; \
		exit 1; \
	fi
	@echo "Running tests for package $(PACKAGE)..."
	CGO_ENABLED=1 $(GOTEST) $(TEST_FLAGS) $(PACKAGE)

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
	CGO_ENABLED=1 $(GOTEST) $(TEST_FLAGS) $(TEST_PACKAGES)
	@mv coverage.out unit_coverage.out
	CGO_ENABLED=1 $(GOTEST) $(INTEGRATION_TEST_FLAGS) $(TEST_PACKAGES)
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
	CGO_ENABLED=1 $(GOTEST) $(PACKAGE_COVERAGE_FLAGS)$(OUTPUT).out $(PACKAGE)
	@echo "Generating coverage report for $(PACKAGE)..."
	$(GOCMD) tool cover -html=$(OUTPUT).out -o $(OUTPUT).html

# ==================================================================================
# Code quality targets
# ==================================================================================

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
	@if [ -f $(GOIMPORTS) ]; then \
		echo "Running goimports..."; \
		$(GOIMPORTS) -w -local github.com/abitofhelp/servicelib ./; \
	fi
	@if [ -f $(GOFUMPT) ]; then \
		echo "Running gofumpt..."; \
		$(GOFUMPT) -w ./; \
	fi

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Run all code quality checks
quality: fmt vet lint security
	@echo "All code quality checks passed!"

# ==================================================================================
# Dependency management targets
# ==================================================================================

# Update dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

# Update all dependencies to their latest versions
update-deps:
	@echo "Updating dependencies to latest versions..."
	@go get -u ./...
	@$(GOMOD) tidy
	@echo "Dependencies updated. Please run tests to verify compatibility."

# Vendor dependencies
vendor: tidy
	@echo "Vendoring dependencies..."
	$(GOMOD) vendor

# ==================================================================================
# Code generation targets
# ==================================================================================

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

# ==================================================================================
# Security targets
# ==================================================================================

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

# ==================================================================================
# Benchmarking targets
# ==================================================================================

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

# ==================================================================================
# Documentation targets
# ==================================================================================

# Generate documentation
docs:
	@echo "Generating documentation..."
	@if ! command -v $(GODOC) > /dev/null; then \
		echo "godoc not found, installing..."; \
		go install golang.org/x/tools/cmd/godoc@latest; \
	fi
	@echo "Documentation generated. Run 'make serve-docs' to view."

# Serve documentation locally
serve-docs:
	@echo "Serving documentation at http://localhost:$(GODOC_PORT)/pkg/github.com/abitofhelp/servicelib/"
	@$(GODOC) -http=:$(GODOC_PORT)

# ==================================================================================
# CI/CD targets
# ==================================================================================

# CI target for continuous integration
ci: check-go-version tidy lint test-all

# Pre-commit hook
pre-commit: fmt vet lint security test
	@echo "Pre-commit checks passed!"

# Release the library
release:
	@echo "Preparing release $(VERSION)..."
	@echo "Ensuring all tests pass..."
	@make test-all
	@echo "Ensuring code quality..."
	@make quality
	@echo "Checking for security vulnerabilities..."
	@make security
	@echo "Updating dependencies..."
	@make tidy
	@echo "Release $(VERSION) is ready!"
	@echo "To complete the release, tag the repository:"
	@echo "git tag -a v$(VERSION) -m 'Release $(VERSION)'"
	@echo "git push origin v$(VERSION)"

# ==================================================================================
# Help
# ==================================================================================

# Help target
help:
	@echo "Available targets:"
	@echo "  all                - Run check-go-version, clean, lint, test, and build"
	@echo "  check-go-version   - Verify the correct Go version is installed"
	@echo "  tools              - Install required development tools"
	@echo "  build              - Compile the code for current platform (no binary produced)"
	@echo "  build-darwin-amd64 - Compile for macOS Intel (amd64)"
	@echo "  build-darwin-arm64 - Compile for macOS Apple Silicon (arm64)"
	@echo "  build-windows-amd64 - Compile for Windows (amd64)"
	@echo "  build-linux-amd64  - Compile for Linux (amd64)"
	@echo "  build-linux-arm64  - Compile for Linux (arm64)"
	@echo "  build-all          - Compile for all platforms"
	@echo "  clean              - Remove test artifacts, tools, and clean Go cache"
	@echo "  test               - Run tests"
	@echo "  test-integration   - Run integration tests (with integration build tag)"
	@echo "  test-all           - Run all tests (unit and integration) with combined coverage report"
	@echo "  test-package       - Run tests for a specific package (PACKAGE=./path/to/package)"
	@echo "  coverage           - Generate test coverage report"
	@echo "  coverage-integration - Generate test coverage report for integration tests"
	@echo "  coverage-package   - Generate test coverage report for a specific package (PACKAGE=./path/to/package OUTPUT=name)"
	@echo "  lint               - Run linter"
	@echo "  fmt                - Format code with gofmt, goimports, and gofumpt"
	@echo "  vet                - Run go vet"
	@echo "  quality            - Run all code quality checks (fmt, vet, lint, security)"
	@echo "  tidy               - Update dependencies"
	@echo "  update-deps        - Update all dependencies to their latest versions"
	@echo "  vendor             - Vendor dependencies"
	@echo "  generate           - Generate mocks using go:generate"
	@echo "  generate-package   - Generate mocks for a specific package (PACKAGE=./path/to/package)"
	@echo "  security           - Check for security vulnerabilities"
	@echo "  bench              - Run benchmarks"
	@echo "  bench-package      - Run benchmarks for a specific package (PACKAGE=./path/to/package)"
	@echo "  docs               - Generate documentation"
	@echo "  serve-docs         - Serve documentation locally"
	@echo "  ci                 - Run continuous integration checks"
	@echo "  pre-commit         - Run pre-commit checks"
	@echo "  release            - Prepare a release"
	@echo "  help               - Show this help message"

# Default target
.DEFAULT_GOAL := help
