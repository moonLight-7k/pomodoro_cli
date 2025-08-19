package errors

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type ErrorCode int

const (
	ErrInvalidArgs ErrorCode = iota + 1000
	ErrInvalidFlag
	ErrInvalidNumber
	ErrInvalidDuration
	ErrTerminalNotSupported
	ErrSessionInterrupted
	ErrConfigLoad
	ErrLogWrite
)

type AppError struct {
	Code      ErrorCode              `json:"code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewAppError(code ErrorCode, message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}
}

type Logger struct {
	file *os.File
}

func NewLogger(logPath string) (*Logger, error) {
	if logPath == "" {
		return &Logger{file: nil}, nil
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &Logger{file: file}, nil
}

func (l *Logger) LogError(err error) {
	logEntry := map[string]interface{}{
		"level":     "error",
		"timestamp": time.Now().Format(time.RFC3339),
		"error":     err.Error(),
	}

	if appErr, ok := err.(*AppError); ok {
		logEntry["code"] = appErr.Code
		logEntry["details"] = appErr.Details
	}

	l.writeLog(logEntry)
}

func (l *Logger) LogInfo(message string, details map[string]interface{}) {
	logEntry := map[string]interface{}{
		"level":     "info",
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   message,
		"details":   details,
	}

	l.writeLog(logEntry)
}

func (l *Logger) writeLog(entry map[string]interface{}) {
	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}

	if l.file != nil {
		l.file.WriteString(string(data) + "\n")
		l.file.Sync()
	} else {
		if entry["level"] == "error" {
			fmt.Fprintf(os.Stderr, "%s\n", data)
		}
	}
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func HandleError(logger *Logger, err error) {
	if logger != nil {
		logger.LogError(err)
	}

	if appErr, ok := err.(*AppError); ok {
		fmt.Fprintf(os.Stderr, "Error: %s\n", appErr.Message)

		switch appErr.Code {
		case ErrInvalidArgs:
			if examples, ok := appErr.Details["examples"].([]string); ok {
				fmt.Fprintf(os.Stderr, "\nExamples:\n")
				for _, example := range examples {
					fmt.Fprintf(os.Stderr, "  %s\n", example)
				}
			}
		case ErrInvalidFlag:
			if flags, ok := appErr.Details["valid_flags"].([]string); ok {
				fmt.Fprintf(os.Stderr, "Valid flags: %v\n", flags)
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
