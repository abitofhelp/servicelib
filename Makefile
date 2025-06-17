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

# Linter parameters
GOLANGCI_LINT=golangci-lint

.PHONY: all build build-darwin-amd64 build-darwin-arm64 build-windows-amd64 build-linux-amd64 build-linux-arm64 build-all clean test coverage lint fmt vet tidy vendor help

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
	@rm -f coverage.out coverage.html
	$(GOCLEAN)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) $(TEST_FLAGS) $(TEST_PACKAGES)

# Generate test coverage report
coverage: test
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

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
	@echo "  coverage           - Generate test coverage report"
	@echo "  lint               - Run linter"
	@echo "  fmt                - Format code"
	@echo "  vet                - Run go vet"
	@echo "  tidy               - Update dependencies"
	@echo "  vendor             - Vendor dependencies"
	@echo "  security           - Check for security vulnerabilities"
	@echo "  bench              - Run benchmarks"
	@echo "  help               - Show this help message"

# Default target
.DEFAULT_GOAL := help
