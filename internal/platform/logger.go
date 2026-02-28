package platform

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// Logger provides structured logging functionality
type Logger struct {
	level  LogLevel
	logger *log.Logger
}

var defaultLogger *Logger

// InitLogger initializes the default logger with the specified level and output file
func InitLogger(level LogLevel, logFilePath string) error {
	// Create log directory if it doesn't exist
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Write to both file and stdout
	multiWriter := io.MultiWriter(file, os.Stdout)

	defaultLogger = &Logger{
		level:  level,
		logger: log.New(multiWriter, "", log.LstdFlags),
	}

	return nil
}

// GetLogger returns the default logger instance
func GetLogger() *Logger {
	if defaultLogger == nil {
		// Fallback to stdout if not initialized
		defaultLogger = &Logger{
			level:  LogLevelInfo,
			logger: log.New(os.Stdout, "", log.LstdFlags),
		}
	}
	return defaultLogger
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= LogLevelDebug {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= LogLevelInfo {
		l.logger.Printf("[INFO] "+format, v...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= LogLevelWarn {
		l.logger.Printf("[WARN] "+format, v...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= LogLevelError {
		l.logger.Printf("[ERROR] "+format, v...)
	}
}

// Fatal logs a fatal error and exits
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.logger.Fatalf("[FATAL] "+format, v...)
}
