# 🍅 Pomodoro CLI

A command-line Pomodoro timer.

## Installation

### Prerequisites

- Go 1.24 or later
- Make (optional, for using Makefile targets)

### Build from Source

```bash
git clone <repository-url>
cd pomodoro_cli
make build
```

### Alternative Build Methods

```bash
# Using Go directly
go build -o bin/pomodoro ./cmd/pomodoro

# Install globally to GOPATH/bin
make install
# or
go install ./cmd/pomodoro
```

## Usage

### Basic Usage

```bash
# 25 minutes work, 5 minutes break
./bin/pomodoro 25 5

# 45 minutes work, 15 minutes break
./bin/pomodoro 45 15
```

### Hour Mode

Use the `-h` flag to specify times in hours instead of minutes:

```bash
# 1 hour work, 30 minutes break (use minutes for fractional hours)
./bin/pomodoro 60 30

# 2 hours work, 1 hour break
./bin/pomodoro 2 1 -h
```

### Examples

```bash
# Standard Pomodoro technique
./bin/pomodoro 25 5

# Extended focus sessions
./bin/pomodoro 45 15

# Short bursts
./bin/pomodoro 15 5

# Long work sessions
./bin/pomodoro 2 1 -h

# Using make to run with defaults
make run
```

## Features in Detail

### Terminal Compatibility

- **Color Support**: Automatically detects terminal color capabilities
- **Fallback Mode**: Works in terminals without ANSI support
- **Progress Bar**: Visual progress indication with fallbacks
- **Cross-Platform**: Works on Linux, macOS, and Windows

### Error Handling

- **Structured Errors**: Detailed error codes and messages
- **Input Validation**: Comprehensive argument validation
- **Graceful Degradation**: Continues operation when possible
- **User-Friendly Messages**: Clear error descriptions with examples

### Logging

- **Session Tracking**: Logs all session starts, completions, and cancellations
- **Statistics**: Tracks work/break time and session counts
- **JSON Format**: Structured logs for easy parsing
- **Error Logging**: Detailed error information for debugging

### Signal Handling

- **Graceful Shutdown**: Responds to Ctrl+C and SIGTERM
- **Session Summary**: Shows statistics on exit
- **Resource Cleanup**: Properly closes files and connections

## Configuration

### Available Make Targets

```bash
make build          # Build the application
make test           # Run all tests
make test-coverage  # Run tests with coverage report
make clean          # Clean build artifacts
make install        # Install to GOPATH/bin
make run            # Run with default settings (25min work, 5min break)
make lint           # Run linter (requires golangci-lint)
make format         # Format code with go fmt and goimports
make help           # Show help information
```

### Time Limits

- **Minimum**: 1 minute/hour
- **Maximum**: 999 minutes/hours (configurable)
- **Validation**: Prevents invalid time specifications

### Display Settings

- **Progress Bar Width**: Automatically adjusts to terminal width
- **Color Scheme**: Consistent purple theme with fallbacks
- **Update Frequency**: 1-second intervals

## Development

### Project Structure

```
pomodoro_cli/
├── cmd/
│   └── pomodoro/
│       └── main.go       # Application entry point
├── internal/
│   ├── config/
│   │   ├── config.go     # Configuration and argument parsing
│   │   └── config_test.go # Configuration tests
│   ├── errors/
│   │   └── errors.go     # Error handling and types
│   ├── session/
│   │   └── session.go    # Session management and tracking
│   └── terminal/
│       └── terminal.go   # Terminal operations and display
├── bin/                  # Built binaries
├── Makefile             # Build automation
├── go.mod               # Go module dependencies
├── go.sum               # Go module checksums
└── README.md            # Documentation
```

### Running Tests

```bash
# Run all tests using make
make test

# Run tests with coverage
make test-coverage

# Run all tests using go directly
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test -v ./internal/config

# Clean build artifacts
make clean
```

### Code Quality

- **Go Standards**: Follows Go coding conventions
- **Error Handling**: Comprehensive error coverage
- **Documentation**: Inline documentation for all public functions
- **Testing**: Unit tests for core functionality

## Troubleshooting

### Common Issues

**Colors not showing**

- Check if your terminal supports ANSI colors
- The app automatically falls back to text-mode

**Progress bar looks wrong**

- Terminal width detection may have failed
- Try resizing your terminal window

**Application won't start**

- Check Go version (requires 1.24+)
- Verify command syntax: `./bin/pomodoro <work> <break> [-h]`
- Make sure you've built the application: `make build`

### Getting Help

```bash
# Show usage information
./bin/pomodoro

# Check for invalid arguments
./bin/pomodoro invalid args
```

### Logs and Debugging

- Errors are logged to stderr in JSON format
- Session information is logged for troubleshooting
- Use `-v` flag for verbose logging (if implemented)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

### Development Guidelines

- Maintain test coverage above 80%
- Follow Go naming conventions
- Add error handling for all operations
- Update documentation for new features

## License

[Add your license here]

## Changelog

### v1.0.0

- Initial production release
- Complete rewrite with modular architecture
- Added comprehensive error handling
- Implemented terminal compatibility
- Added session tracking and logging
- Full test coverage
- Production-ready signal handling
