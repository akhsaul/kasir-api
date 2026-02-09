package logger

import (
	"log"
	"os"
)

var (
	// InfoLogger logs informational messages to stdout.
	InfoLogger *log.Logger
	// ErrorLogger logs error messages to stderr.
	ErrorLogger *log.Logger
)

// InitLogger initializes the global loggers.
func InitLogger() {
	InfoLogger = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	ErrorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
}

// Info logs an informational message to stdout.
func Info(format string, v ...any) {
	if InfoLogger == nil {
		InitLogger()
	}
	InfoLogger.Printf(format, v...)
}

// Error logs an error message to stderr.
func Error(format string, v ...any) {
	if ErrorLogger == nil {
		InitLogger()
	}
	ErrorLogger.Printf(format, v...)
}

// Fatal logs an error message to stderr and exits with status code 1.
func Fatal(v ...any) {
	if ErrorLogger == nil {
		InitLogger()
	}
	ErrorLogger.Fatal(v...)
}
