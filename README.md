# 🍅 Pomodoro CLI (Vibe Coded)

A command-line Pomodoro timer built in Go.

## 📦 Installation

### Prerequisites

- Go 1.24 or later
- Make (optional, for using Makefile targets)

### Quick Install

```bash
git clone https://github.com/moonLight-7k/pomodoro_cli.git
cd pomodoro_cli
make build
```

### Alternative Build Methods

```bash
# Build using Go directly
go build -o bin/pomodoro ./cmd/pomodoro

# Install globally to GOPATH/bin
make install
# or
go install ./cmd/pomodoro

# Build for multiple platforms
make build-all
```

## 🚀 Usage

### Basic Usage

```bash
# Standard Pomodoro technique (25 min work, 5 min break)
./bin/pomodoro 25 5

# Extended focus sessions (45 min work, 15 min break)
./bin/pomodoro 45 15

# Short bursts (15 min work, 5 min break)
./bin/pomodoro 15 5
```

### Hour Mode

Use the `-h` flag to specify times in hours:

```bash
# 2 hours work, 1 hour break
./bin/pomodoro 2 1 -h

# 1.5 hours work, 30 min break (mix units by converting to minutes)
./bin/pomodoro 90 30
```

### What You'll See

During a work session:

```
🍅 Pomodoro CLI v1.0.0
Work: 25m0s, Break: 5m0s

work
2:30 PM - 0m45s
████████████████                    75%

Press Ctrl+C to exit
```

When transitioning between sessions:

```
Work complete!
✓ Session finished successfully

break
2:55 PM - 0m15s
█████                               30%

Press Ctrl+C to exit
```

## 🛠️ Available Commands

### Make Targets

```bash
make build          # Build the application
make test           # Run all tests
make test-coverage  # Run tests with coverage report
make clean          # Clean build artifacts
make install        # Install to GOPATH/bin
make run            # Run with default settings (25min work, 5min break)
make lint           # Run linter (requires golangci-lint)
make format         # Format code with go fmt and goimports
make vet            # Run go vet for code analysis
make check          # Run all quality checks (format, vet, lint, test)
make build-all      # Build for multiple platforms
make package        # Create release packages
make dev-setup      # Set up development environment
make help           # Show all available targets
```

### Direct Go Commands

```bash
# Run with custom times
go run ./cmd/pomodoro 30 10

# Build and run
go build ./cmd/pomodoro && ./pomodoro 25 5

# Run tests
go test ./...
```

## ⚙️ Configuration

### Time Limits & Validation

- **Range**: 1 to 999 minutes or hours
- **Input Validation**: Prevents invalid time specifications
- **Maximum Session**: 12 hours (configurable in code)
- **Positive Values Only**: Rejects zero or negative times

### Display Settings

- **Progress Bar**: Automatically adjusts to terminal width (minimum 10, maximum terminal width - 10)
- **Color Scheme**: Purple theme with white accents and fallbacks for non-color terminals
- **Update Frequency**: Real-time updates every second
- **Terminal Compatibility**: Works with most terminal types

### Error Codes

The application uses structured error codes for better debugging:

- `1000`: Invalid arguments
- `1001`: Invalid flag
- `1002`: Invalid number format
- `1003`: Invalid duration
- `1004`: Terminal not supported
- `1005`: Session interrupted
- `1006`: Configuration load error
- `1007`: Log write error

## 🏗️ Development

### Project Structure

```
pomodoro_cli/
├── cmd/
│   └── pomodoro/
│       └── main.go           # Modern application entry point (recommended)
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration management and CLI parsing
│   ├── errors/
│   │   └── errors.go         # Structured error handling and logging
│   ├── session/
│   │   └── session.go        # Session management and cycle control
│   └── terminal/
│       └── terminal.go       # Terminal operations and display logic
├── test/
│   └── config_test.go        # Unit tests
├── bin/                      # Built binaries (created during build)
├── main.go                   # Legacy simple implementation
├── Makefile                  # Build automation and development tools
├── go.mod                    # Go module dependencies
├── go.sum                    # Go module checksums
├── Dockerfile                # Container deployment
└── README.md                 # This documentation
```

### Architecture Overview

- **`cmd/pomodoro/main.go`**: Modern entry point with full feature set
- **`main.go`**: Simple legacy implementation (basic functionality only)
- **`internal/config`**: Command-line argument parsing and validation
- **`internal/session`**: Core Pomodoro session logic and cycle management
- **`internal/terminal`**: Terminal capabilities detection and UI rendering
- **`internal/errors`**: Centralized error handling with structured logging

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests using go directly
go test ./...

# Run tests with verbose output and race detection
go test -v -race ./...

# Run specific package tests
go test -v ./internal/config
go test -v ./test

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Development Workflow

```bash
# Set up development environment
make dev-setup

# Format code
make format

# Run linter
make lint

# Run all quality checks
make check

# Build and test
make build && make test
```

### Debugging

- **Logs**: Errors are logged to stderr in JSON format for easy parsing
- **Session Data**: All session information is logged with timestamps
- **Verbose Mode**: Use structured logging for detailed debugging information
- **Exit Codes**: Application returns appropriate exit codes for scripting

## 🤝 Contributing

I welcome contributions! Please follow these guidelines:

### Getting Started

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `make check`
6. Commit with clear messages
7. Push to your fork
8. Submit a pull request

### Development Guidelines

- **Test Coverage**: Maintain or improve test coverage
- **Code Style**: Run `make format` before committing
- **Documentation**: Update README for new features
- **Error Handling**: Add proper error handling for all operations
- **Logging**: Include appropriate logging for debugging
