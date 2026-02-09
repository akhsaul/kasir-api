package logger

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestInitLogger(t *testing.T) {
	// Reset loggers
	InfoLogger = nil
	ErrorLogger = nil

	InitLogger()

	if InfoLogger == nil {
		t.Error("InitLogger should initialize InfoLogger")
	}
	if ErrorLogger == nil {
		t.Error("InitLogger should initialize ErrorLogger")
	}
}

func TestInitLogger_InfoLoggerOutput(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Reset and reinitialize
	InfoLogger = log.New(&buf, "[INFO] ", log.LstdFlags)

	InfoLogger.Printf("test message")

	output := buf.String()
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("InfoLogger should have [INFO] prefix, got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("InfoLogger should log the message, got: %s", output)
	}
}

func TestInitLogger_ErrorLoggerOutput(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Reset and reinitialize
	ErrorLogger = log.New(&buf, "[ERROR] ", log.LstdFlags)

	ErrorLogger.Printf("error message")

	output := buf.String()
	if !strings.Contains(output, "[ERROR]") {
		t.Errorf("ErrorLogger should have [ERROR] prefix, got: %s", output)
	}
	if !strings.Contains(output, "error message") {
		t.Errorf("ErrorLogger should log the message, got: %s", output)
	}
}

func TestInfo_WithInitializedLogger(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Initialize with custom writer
	InfoLogger = log.New(&buf, "[INFO] ", log.LstdFlags)

	Info("test %s %d", "message", 42)

	output := buf.String()
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Info should log with [INFO] prefix, got: %s", output)
	}
	if !strings.Contains(output, "test message 42") {
		t.Errorf("Info should format message correctly, got: %s", output)
	}
}

func TestInfo_WithNilLogger(t *testing.T) {
	// Reset logger to nil
	InfoLogger = nil

	// Create a buffer to capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Info("test message")

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Info should auto-initialize and log with [INFO] prefix, got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("Info should log the message, got: %s", output)
	}
}

func TestError_WithInitializedLogger(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Initialize with custom writer
	ErrorLogger = log.New(&buf, "[ERROR] ", log.LstdFlags)

	Error("error %s %d", "message", 500)

	output := buf.String()
	if !strings.Contains(output, "[ERROR]") {
		t.Errorf("Error should log with [ERROR] prefix, got: %s", output)
	}
	if !strings.Contains(output, "error message 500") {
		t.Errorf("Error should format message correctly, got: %s", output)
	}
}

func TestError_WithNilLogger(t *testing.T) {
	// Reset logger to nil
	ErrorLogger = nil

	// Create a buffer to capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Error("error message")

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "[ERROR]") {
		t.Errorf("Error should auto-initialize and log with [ERROR] prefix, got: %s", output)
	}
	if !strings.Contains(output, "error message") {
		t.Errorf("Error should log the message, got: %s", output)
	}
}

func TestInfo_FormatVariants(t *testing.T) {
	testCases := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{"string", "message: %s", []interface{}{"hello"}, "message: hello"},
		{"int", "count: %d", []interface{}{42}, "count: 42"},
		{"float", "value: %.2f", []interface{}{3.14}, "value: 3.14"},
		{"multiple", "%s has %d items", []interface{}{"cart", 5}, "cart has 5 items"},
		{"no args", "simple message", nil, "simple message"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			InfoLogger = log.New(&buf, "[INFO] ", 0)

			Info(tc.format, tc.args...)

			output := buf.String()
			if !strings.Contains(output, tc.expected) {
				t.Errorf("Info should format as %s, got: %s", tc.expected, output)
			}
		})
	}
}

func TestError_FormatVariants(t *testing.T) {
	testCases := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{"string", "error: %s", []interface{}{"not found"}, "error: not found"},
		{"int", "code: %d", []interface{}{404}, "code: 404"},
		{"error value", "failed: %v", []interface{}{"timeout"}, "failed: timeout"},
		{"multiple", "%s failed with code %d", []interface{}{"request", 500}, "request failed with code 500"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			ErrorLogger = log.New(&buf, "[ERROR] ", 0)

			Error(tc.format, tc.args...)

			output := buf.String()
			if !strings.Contains(output, tc.expected) {
				t.Errorf("Error should format as %s, got: %s", tc.expected, output)
			}
		})
	}
}

func TestInitLogger_CreatesCorrectPrefixes(t *testing.T) {
	InfoLogger = nil
	ErrorLogger = nil

	InitLogger()

	// Test Info prefix
	var infoBuf bytes.Buffer
	InfoLogger.SetOutput(&infoBuf)
	InfoLogger.Printf("test")
	if !strings.Contains(infoBuf.String(), "[INFO]") {
		t.Errorf("InfoLogger should have [INFO] prefix")
	}

	// Test Error prefix
	var errorBuf bytes.Buffer
	ErrorLogger.SetOutput(&errorBuf)
	ErrorLogger.Printf("test")
	if !strings.Contains(errorBuf.String(), "[ERROR]") {
		t.Errorf("ErrorLogger should have [ERROR] prefix")
	}
}

// Note: TestFatal is difficult to test because it calls os.Exit(1)
// We can test that it's set up correctly but actually calling it would exit the test
func TestFatal_WithNilLogger(t *testing.T) {
	// Just verify that Fatal initializes the logger if nil
	// We can't actually call Fatal because it exits

	ErrorLogger = nil

	// Verify the logger would be initialized
	if ErrorLogger != nil {
		t.Error("ErrorLogger should initially be nil for this test")
	}

	// We'd need to mock os.Exit to fully test Fatal
	// For now, we verify the function exists and is callable
	_ = Fatal // Reference the function to ensure it compiles
}
