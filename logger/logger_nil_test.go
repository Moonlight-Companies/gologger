package logger

import (
	"bytes"
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

	// Create logger with a buffer writer and zero time
	logger := NewLogger("TEST", WithLevel(LogLevelDebug), WithZeroTime(), WithWriter(&buf))

	// Test cases
	testCases := []struct {
		name     string
		args     []interface{}
		expected string // Expected exact output
	}{
		{
			name:     "Untyped nil",
			args:     []interface{}{nil},
			expected: "0001/01/01 00:00:00.000000 DEBUG: TEST <nil arg 0>\n",
		},
		{
			name:     "Typed nil pointer",
			args:     []interface{}{(*TestStruct)(nil)},
			expected: "0001/01/01 00:00:00.000000 DEBUG: TEST <nil *logger.TestStruct at arg 0>\n",
		},
		{
			name:     "Multiple nil types",
			args:     []interface{}{nil, (*TestStruct)(nil), "valid string", (*int)(nil)},
			expected: "0001/01/01 00:00:00.000000 DEBUG: TEST <nil arg 0> <nil *logger.TestStruct at arg 1> valid string <nil *int at arg 3>\n",
		},
		{
			name:     "Mixed valid and nil values",
			args:     []interface{}{123, nil, "test", (*TestStruct)(nil)},
			expected: "0001/01/01 00:00:00.000000 DEBUG: TEST 123 <nil arg 1> test <nil *logger.TestStruct at arg 3>\n",
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

			// Check for exact output match
			if output != tc.expected {
				t.Errorf("Expected output %q, got: %q", tc.expected, output)
			}
		})
	}
}

// TestLoggerNilStringMethod tests that typed nil pointers with String() methods
// are handled properly and don't cause panics
func TestLoggerNilStringMethod(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TEST", WithLevel(LogLevelDebug), WithZeroTime(), WithWriter(&buf))

	// Create a nil pointer to our test struct
	var nilStruct *TestStruct = nil

	// This should call our custom logln method and detect the nil pointer
	// rather than calling the String() method on the nil pointer
	logger.Debugln(nilStruct)

	output := buf.String()
	expected := "0001/01/01 00:00:00.000000 DEBUG: TEST <nil *logger.TestStruct at arg 0>\n"

	if output != expected {
		t.Errorf("Expected output %q, but got: %q", expected, output)
	}

	// Verify that the String() method was not called (which would output a different message)
	unexpectedContent := "TestStruct is nil but String() was called"
	if strings.Contains(output, unexpectedContent) {
		t.Errorf("String() method was called on nil pointer, output contains: %q", unexpectedContent)
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

type TestPanicStruct struct {
	Name string
}

func (t *TestPanicStruct) String() string {
	panic("hello world")
}

// TestDirectStringCall verifies that our String method works correctly when called directly
func TestDirectStringPanicCall(t *testing.T) {
	// Test with a String fn that panics
	panicStruct := &TestPanicStruct{Name: "World"}

	var buf bytes.Buffer
	logger := NewLogger("TEST", WithLevel(LogLevelDebug), WithZeroTime(), WithWriter(&buf))
	logger.Debugln(panicStruct)

	output := buf.String()
	expected := "0001/01/01 00:00:00.000000 DEBUG: TEST %!v(PANIC=String method: hello world)\n"

	if output != expected {
		t.Errorf("Expected output to contain %q, but got: %q", expected, output)
	}
}
