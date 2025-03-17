package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"
)

// TestLoggerLevels tests that log levels work correctly (messages are shown or suppressed based on level)
func TestLoggerLevels(t *testing.T) {
	tests := []struct {
		name     string
		logLevel LogLevel
		debug    bool // whether debug messages should appear
		info     bool // whether info messages should appear
		warn     bool // whether warn messages should appear
		error    bool // whether error messages should appear
	}{
		{
			name:     "Debug level shows all messages",
			logLevel: LogLevelDebug,
			debug:    true,
			info:     true,
			warn:     true,
			error:    true,
		},
		{
			name:     "Info level hides debug messages",
			logLevel: LogLevelInfo,
			debug:    false,
			info:     true,
			warn:     true,
			error:    true,
		},
		{
			name:     "Warn level hides debug and info messages",
			logLevel: LogLevelWarn,
			debug:    false,
			info:     false,
			warn:     true,
			error:    true,
		},
		{
			name:     "Error level only shows error messages",
			logLevel: LogLevelError,
			debug:    false,
			info:     false,
			warn:     false,
			error:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger(tc.logLevel, "TEST")
			logger.logger = log.New(&buf, "", 0) // No timestamps for easier testing

			// Log messages at each level
			logger.Debugln("debug message")
			logger.Infoln("info message")
			logger.Warnln("warn message")
			logger.Errorln("error message")

			output := buf.String()

			// Check if each message type appears as expected
			assertMessagePresence(t, output, "DEBUG", tc.debug)
			assertMessagePresence(t, output, "INFO", tc.info)
			assertMessagePresence(t, output, "WARN", tc.warn)
			assertMessagePresence(t, output, "ERROR", tc.error)
		})
	}
}

// TestLoggerPrefix tests that the prefix functionality works correctly
func TestLoggerPrefix(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelDebug, "INITIAL")
	logger.logger = log.New(&buf, "", 0)

	// Test initial prefix
	logger.Infoln("First message")
	if !strings.Contains(buf.String(), "INITIAL") {
		t.Errorf("Expected output to contain initial prefix 'INITIAL', got: %q", buf.String())
	}

	// Test setting new prefix
	buf.Reset()
	logger.SetPrefix("UPDATED")

	// Verify GetPrefix returns the updated value
	if prefix := logger.GetPrefix(); prefix != "UPDATED" {
		t.Errorf("GetPrefix() = %q, want %q", prefix, "UPDATED")
	}

	// Verify the updated prefix appears in log output
	logger.Infoln("Second message")
	if !strings.Contains(buf.String(), "UPDATED") {
		t.Errorf("Expected output to contain updated prefix 'UPDATED', got: %q", buf.String())
	}
}

// TestLoggerDeltaTime tests the delta time option
func TestLoggerDeltaTime(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelDebug, "TEST")
	logger.logger = log.New(&buf, "", 0)

	// Test without delta time
	logger.SetIncludeDeltaTime(false)
	buf.Reset()
	logger.Infoln("Without delta time")
	outputWithoutDT := buf.String()

	// Ensure no timestamp format is present
	if strings.Contains(outputWithoutDT, "INFO ") && strings.Contains(outputWithoutDT, "s") {
		t.Errorf("Expected no delta time, but found timing information: %q", outputWithoutDT)
	}

	// Test with delta time
	logger.SetIncludeDeltaTime(true)
	buf.Reset()

	// Sleep a small amount to ensure measurable time passes
	time.Sleep(10 * time.Millisecond)

	logger.Infoln("With delta time")
	outputWithDT := buf.String()

	// The format should be like "INFO 10.123ms: TEST message"
	// We're checking for a space after INFO and a time unit (s, ms, µs, ns)
	if !strings.Contains(outputWithDT, "INFO ") ||
		!(strings.Contains(outputWithDT, "s:") ||
			strings.Contains(outputWithDT, "ms:") ||
			strings.Contains(outputWithDT, "µs:") ||
			strings.Contains(outputWithDT, "ns:")) {
		t.Errorf("Expected delta time format not found in: %q", outputWithDT)
	}
}

// TestLoggerFormatting tests both Printf-style and Println-style logging
func TestLoggerFormatting(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelDebug, "TEST")
	logger.logger = log.New(&buf, "", 0)

	// Test formatted logging (Printf-style)
	buf.Reset()
	logger.Info("Formatted %s with %d values", "message", 2)
	formatted := buf.String()
	expected := "INFO: TEST Formatted message with 2 values\n"
	if formatted != expected {
		t.Errorf("Expected formatted output %q, got: %q", expected, formatted)
	}

	// Test vararg logging (Println-style)
	buf.Reset()
	logger.Infoln("Unformatted", "message", "with", 3, "values")
	unformatted := buf.String()
	expected = "INFO: TEST Unformatted message with 3 values\n"
	if unformatted != expected {
		t.Errorf("Expected unformatted output %q, got: %q", expected, unformatted)
	}
}

// TestNilHandling tests that nil values are handled correctly in Println-style logging
func TestNilHandling(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelDebug, "TEST")
	logger.logger = log.New(&buf, "", 0)

	// Test with nil values
	buf.Reset()
	var nilSlice []int = nil
	logger.Debugln("A nil value:", nil, "and a typed nil:", nilSlice)

	output := buf.String()
	expected := "DEBUG: TEST A nil value: <nil arg 1> and a typed nil: <nil int slice at arg 3>\n"

	if output != expected {
		t.Errorf("Expected nil handling output %q, got: %q", expected, output)
	}
}

// Helper function to check if messages with a certain level are present or absent as expected
func assertMessagePresence(t *testing.T, output, level string, shouldBePresent bool) {
	t.Helper()
	hasMessage := strings.Contains(output, level)

	if shouldBePresent && !hasMessage {
		t.Errorf("Expected %s message to be present, but it was not found in: %q", level, output)
	} else if !shouldBePresent && hasMessage {
		t.Errorf("Expected %s message to be absent, but it was found in: %q", level, output)
	}
}
