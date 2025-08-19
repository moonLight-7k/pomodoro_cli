# Pomodoro CLI Makefile

.PHONY: build test clean install run lint format help

# Variables
BINARY_NAME=pomodoro
MAIN_PACKAGE=./cmd/pomodoro
INSTALL_PACKAGE=pomodoro_cli/cmd/pomodoro
BUILD_DIR=./bin
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -cover ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) $(INSTALL_PACKAGE)

# Run the application with default arguments
run:
	@echo "Running $(BINARY_NAME) with default settings (25min work, 5min break)..."
	go run $(MAIN_PACKAGE) 25 5

# Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run

# Format code
format:
	@echo "Formatting code..."
	go fmt ./...
	@which goimports > /dev/null && goimports -w . || echo "goimports not found, skipping import formatting"

# Vet code
vet:
	@echo "Running go vet..."
	go vet ./...

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	
	# macOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	
	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	
	@echo "Built binaries for multiple platforms in $(BUILD_DIR)/"

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	go mod download
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@which goimports > /dev/null || go install golang.org/x/tools/cmd/goimports@latest

# Quick quality check
check: format vet lint test
	@echo "All checks passed!"

# Package for release
package: clean build-all
	@echo "Creating release packages..."
	@mkdir -p $(BUILD_DIR)/releases
	
	# Create tar.gz for Unix systems
	cd $(BUILD_DIR) && tar -czf releases/$(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	cd $(BUILD_DIR) && tar -czf releases/$(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	cd $(BUILD_DIR) && tar -czf releases/$(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	
	# Create zip for Windows
	cd $(BUILD_DIR) && zip releases/$(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	
	@echo "Release packages created in $(BUILD_DIR)/releases/"

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  clean        - Clean build artifacts"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  run          - Run with default settings"
	@echo "  lint         - Run linter"
	@echo "  format       - Format code"
	@echo "  vet          - Run go vet"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  dev-setup    - Set up development environment"
	@echo "  check        - Run all quality checks"
	@echo "  package      - Create release packages"
	@echo "  help         - Show this help"
