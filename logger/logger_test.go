package logger

import (
	"bytes"
	"strings"
	"testing"
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
			logger := NewLogger("TEST", WithLevel(tc.logLevel), WithZeroTime(), WithWriter(&buf))

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
	logger := NewLogger("INITIAL", WithLevel(LogLevelDebug), WithZeroTime(), WithWriter(&buf))

	// Test initial prefix
	logger.Infoln("First message")
	output := buf.String()
	expected := "0001/01/01 00:00:00.000000 INFO: INITIAL First message\n"
	if output != expected {
		t.Errorf("Expected output %q, got: %q", expected, output)
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
	output = buf.String()
	expected = "0001/01/01 00:00:00.000000 INFO: UPDATED Second message\n"
	if output != expected {
		t.Errorf("Expected output %q, got: %q", expected, output)
	}
}

// TestLoggerDeltaTime tests the delta time option
func TestLoggerDeltaTime(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TEST", WithLevel(LogLevelDebug), WithZeroTime(), WithWriter(&buf))

	// Test without delta time
	logger.SetIncludeDeltaTime(false)
	buf.Reset()
	logger.Infoln("Without delta time")
	outputWithoutDT := buf.String()
	expectedWithoutDT := "0001/01/01 00:00:00.000000 INFO: TEST Without delta time\n"
	if outputWithoutDT != expectedWithoutDT {
		t.Errorf("Expected output %q, got: %q", expectedWithoutDT, outputWithoutDT)
	}

	// Test with delta time
	logger.SetIncludeDeltaTime(true)
	buf.Reset()
	logger.Infoln("With delta time")
	outputWithDT := buf.String()
	
	// Even with WithZeroTime(), the delta time is still calculated using time.Since
	// Just check if it has the correct format pattern
	if !strings.Contains(outputWithDT, "0001/01/01 00:00:00.000000 INFO ") ||
	   !strings.Contains(outputWithDT, "s: TEST With delta time\n") {
		t.Errorf("Expected output format not found in: %q", outputWithDT)
	}
}

// TestLoggerFormatting tests both Printf-style and Println-style logging
func TestLoggerFormatting(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TEST", WithLevel(LogLevelInfo), WithZeroTime(), WithWriter(&buf))

	// Test formatted logging (Printf-style)
	buf.Reset()
	logger.Info("Formatted %s with %d values", "message", 2)
	formatted := buf.String()
	
	expected := "0001/01/01 00:00:00.000000 INFO: TEST Formatted message with 2 values\n"
	if formatted != expected {
		t.Errorf("Expected output %q, got: %q", expected, formatted)
	}

	// Test vararg logging (Println-style)
	buf.Reset()
	logger.Infoln("Unformatted", "message", "with", 3, "values")
	unformatted := buf.String()
	
	expected = "0001/01/01 00:00:00.000000 INFO: TEST Unformatted message with 3 values\n"
	if unformatted != expected {
		t.Errorf("Expected output %q, got: %q", expected, unformatted)
	}
}

// Helper function to check if messages with a certain level are present or absent as expected
func assertMessagePresence(t *testing.T, output, level string, shouldBePresent bool) {
	t.Helper()
	
	// With zero time, we can check for a specific format with the log level
	expectedPattern := "0001/01/01 00:00:00.000000 " + level
	hasMessage := strings.Contains(output, expectedPattern)

	if shouldBePresent && !hasMessage {
		t.Errorf("Expected %s message to be present, but it was not found in: %q", level, output)
	} else if !shouldBePresent && hasMessage {
		t.Errorf("Expected %s message to be absent, but it was found in: %q", level, output)
	}
}