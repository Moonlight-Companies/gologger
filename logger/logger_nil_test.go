package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

// Custom test struct to test typed nil pointers
type TestStruct struct {
	Name string
}

func (t *TestStruct) String() string {
	if t == nil {
		return "TestStruct is nil but String() was called"
	}
	return "TestStruct: " + t.Name
}

func TestLoggerNilHandling(t *testing.T) {
	// Redirect logger output to a buffer for testing
	var buf bytes.Buffer

	// Create logger with a custom logger that writes to our buffer
	logger := NewLogger(LogLevelDebug, "TEST")
	logger.logger = log.New(&buf, "", 0) // No timestamp prefixes for simpler testing

	// Test cases
	testCases := []struct {
		name     string
		args     []interface{}
		expected string // Exact expected output
	}{
		{
			name:     "Untyped nil",
			args:     []interface{}{nil},
			expected: "DEBUG: TEST <nil arg 0>\n",
		},
		{
			name:     "Typed nil pointer",
			args:     []interface{}{(*TestStruct)(nil)},
			expected: "DEBUG: TEST <nil logger.TestStruct pointer at arg 0>\n",
		},
		{
			name:     "Multiple nil types",
			args:     []interface{}{nil, (*TestStruct)(nil), "valid string", (*int)(nil)},
			expected: "DEBUG: TEST <nil arg 0> <nil logger.TestStruct pointer at arg 1> valid string <nil int pointer at arg 3>\n",
		},
		{
			name:     "Mixed valid and nil values",
			args:     []interface{}{123, nil, "test", (*TestStruct)(nil)},
			expected: "DEBUG: TEST 123 <nil arg 1> test <nil logger.TestStruct pointer at arg 3>\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear the buffer before each test
			buf.Reset()

			// Call the logger method
			logger.Debugln(tc.args...)

			// Get the output
			output := buf.String()

			// Check exact output match
			if output != tc.expected {
				t.Errorf("Expected output %q, but got: %q", tc.expected, output)
			}
		})
	}
}

// TestLoggerNilStringMethod tests that typed nil pointers with String() methods
// are handled properly and don't cause panics
func TestLoggerNilStringMethod(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelDebug, "TEST")
	logger.logger = log.New(&buf, "", 0)

	// Create a nil pointer to our test struct
	var nilStruct *TestStruct = nil

	// This should call our custom logln method and detect the nil pointer
	// rather than calling the String() method on the nil pointer
	logger.Debugln(nilStruct)

	output := buf.String()
	expected := "DEBUG: TEST <nil logger.TestStruct pointer at arg 0>\n"

	if output != expected {
		t.Errorf("Expected output %q, but got: %q", expected, output)
	}

	// Verify that the String() method was not called (which would output a different message)
	unexpected := "TestStruct is nil but String() was called"
	if strings.Contains(output, unexpected) {
		t.Errorf("String() method was called on nil pointer, output contains: %q", unexpected)
	}
}

// TestDirectStringCall verifies that our String method works correctly when called directly
func TestDirectStringCall(t *testing.T) {
	// Test with nil pointer to ensure our String() method handles nil correctly
	var nilStruct *TestStruct = nil
	result := nilStruct.String()
	expected := "TestStruct is nil but String() was called"

	if result != expected {
		t.Errorf("Expected output %q, but got: %q", expected, result)
	}

	// Test with valid struct
	validStruct := &TestStruct{Name: "Test"}
	result = validStruct.String()
	expected = "TestStruct: Test"

	if result != expected {
		t.Errorf("Expected output %q, but got: %q", expected, result)
	}
}
