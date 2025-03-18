package logger

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

// LogLevel defines severity levels for logging
type LogLevel int

// String returns the string representation of a LogLevel
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Log levels in ascending order of severity
const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// Logger provides a simple space-delimited logging capability with prefixes and levels
type Logger struct {
	writer        io.Writer
	level         LogLevel
	createTime    time.Time
	includeDeltaT bool
	zeroT         bool
	prefix        string
	mu            sync.RWMutex
}

// LoggerOption defines a functional option for configuring a Logger
type LoggerOption func(*Logger)

// WithDeltaTime enables or disables including time since logger creation
func WithDeltaTime(include bool) LoggerOption {
	return func(l *Logger) {
		l.includeDeltaT = include
	}
}

// WithLevel sets the initial log level
func WithLevel(level LogLevel) LoggerOption {
	return func(l *Logger) {
		l.level = level
	}
}

func WithWriter(w io.Writer) LoggerOption {
	return func(l *Logger) {
		l.writer = w
	}
}

func WithZeroTime() LoggerOption {
	return func(l *Logger) {
		l.zeroT = true
	}
}

// NewLogger creates a new Logger with the specified prefix and options
func NewLogger(prefix string, options ...LoggerOption) *Logger {
	l := &Logger{
		writer:     os.Stdout,
		level:      LogLevelDebug, // Default level
		prefix:     prefix,
		createTime: time.Now(),
	}

	// Apply options
	for _, option := range options {
		option(l)
	}

	return l
}

func (l *Logger) now() time.Time {
	if l.zeroT {
		return time.Time{}
	}
	return time.Now()
}

// SetIncludeDeltaTime configures whether to include time since logger creation
func (l *Logger) SetIncludeDeltaTime(include bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.includeDeltaT = include
}

// SetPrefix updates the logger prefix
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// GetPrefix returns the current logger prefix
func (l *Logger) GetPrefix() string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.prefix
}

// SetLevel updates the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() LogLevel {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

// Debug logs a formatted message at DEBUG level
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.GetLevel() <= LogLevelDebug {
		l.log(LogLevelDebug, format, v...)
	}
}

// Info logs a formatted message at INFO level
func (l *Logger) Info(format string, v ...interface{}) {
	if l.GetLevel() <= LogLevelInfo {
		l.log(LogLevelInfo, format, v...)
	}
}

// Warn logs a formatted message at WARN level
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.GetLevel() <= LogLevelWarn {
		l.log(LogLevelWarn, format, v...)
	}
}

// Error logs a formatted message at ERROR level
func (l *Logger) Error(format string, v ...interface{}) {
	if l.GetLevel() <= LogLevelError {
		l.log(LogLevelError, format, v...)
	}
}

// Debugln logs a space-separated list of values at DEBUG level
func (l *Logger) Debugln(v ...interface{}) {
	if l.GetLevel() <= LogLevelDebug {
		l.logln(LogLevelDebug, v...)
	}
}

// Infoln logs a space-separated list of values at INFO level
func (l *Logger) Infoln(v ...interface{}) {
	if l.GetLevel() <= LogLevelInfo {
		l.logln(LogLevelInfo, v...)
	}
}

// Warnln logs a space-separated list of values at WARN level
func (l *Logger) Warnln(v ...interface{}) {
	if l.GetLevel() <= LogLevelWarn {
		l.logln(LogLevelWarn, v...)
	}
}

// Errorln logs a space-separated list of values at ERROR level
func (l *Logger) Errorln(v ...interface{}) {
	if l.GetLevel() <= LogLevelError {
		l.logln(LogLevelError, v...)
	}
}

// log handles formatted logging
func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	levelStr := level.String()
	prefix := l.getLogPrefix(levelStr)
	message := fmt.Sprintf(format, v...)

	// Get current time for timestamping
	now := l.now()
	timeStr := now.Format("2006/01/02 15:04:05.000000")

	// Format the full log line
	logLine := fmt.Sprintf("%s %s %s\n", timeStr, prefix, message)

	// Write to the writer
	l.mu.Lock() // Lock to ensure atomic writes
	defer l.mu.Unlock()
	fmt.Fprint(l.writer, logLine)
}

// logln handles unformatted logging with space-separated values
func (l *Logger) logln(level LogLevel, v ...interface{}) {
	levelStr := level.String()
	prefix := l.getLogPrefix(levelStr)
	message := l.formatArgs(v...)

	// Get current time for timestamping
	now := l.now()
	timeStr := now.Format("2006/01/02 15:04:05.000000")

	// Format the full log line
	logLine := fmt.Sprintf("%s %s %s\n", timeStr, prefix, message)

	// Write to the writer
	l.mu.Lock() // Lock to ensure atomic writes
	defer l.mu.Unlock()
	fmt.Fprint(l.writer, logLine)
}

// getLogPrefix builds the log prefix with optional delta time
func (l *Logger) getLogPrefix(levelStr string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.includeDeltaT {
		dt := time.Since(l.createTime)
		return fmt.Sprintf("%s %s: %s", levelStr, dt, l.prefix)
	}
	return fmt.Sprintf("%s: %s", levelStr, l.prefix)
}

func formatArgIntoString(arg interface{}) (s string) {
	defer func() {
		if r := recover(); r != nil {
			s = fmt.Sprintf("<error printing arg: %v>", r)
		}
	}()
	return fmt.Sprint(arg)
}

// formatArgs handles special formatting for nil values
func (l *Logger) formatArgs(v ...interface{}) string {
	args := make([]string, len(v))

	for i, arg := range v {
		if arg == nil {
			args[i] = fmt.Sprintf("<nil arg %d>", i)
			continue
		}

		val := reflect.ValueOf(arg)
		if (val.Kind() == reflect.Ptr || val.Kind() == reflect.Slice || val.Kind() == reflect.Map) && val.IsNil() {
			args[i] = fmt.Sprintf("<nil %s at arg %d>", val.Type(), i)
			continue
		}

		args[i] = formatArgIntoString(arg)
	}

	return strings.Join(args, " ")
}
