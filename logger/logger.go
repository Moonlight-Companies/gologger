package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type LogLevel int

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

const (
	LogLevelDebug = LogLevel(iota)
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type Logger struct {
	logger            *log.Logger
	level             LogLevel
	create            time.Time
	option_include_dt bool
	prefix            string
	mu                sync.RWMutex
}

func NewLogger(level LogLevel, prefix string) *Logger {
	return &Logger{
		logger:            log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds),
		level:             level,
		prefix:            prefix,
		create:            time.Now(),
		option_include_dt: false,
	}
}

func (l *Logger) SetIncludeDeltaTime(include bool) {
	l.option_include_dt = include
}

func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

func (l *Logger) GetPrefix() string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.prefix
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= LogLevelDebug {
		l.log("DEBUG", format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= LogLevelInfo {
		l.log("INFO", format, v...)
	}
}

func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= LogLevelWarn {
		l.log("WARN", format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= LogLevelError {
		l.log("ERROR", format, v...)
	}
}

// Log methods without formatting (Println-like behavior)
func (l *Logger) Debugln(v ...interface{}) {
	if l.level <= LogLevelDebug {
		l.logln("DEBUG", v...)
	}
}

func (l *Logger) Infoln(v ...interface{}) {
	if l.level <= LogLevelInfo {
		l.logln("INFO", v...)
	}
}

func (l *Logger) Warnln(v ...interface{}) {
	if l.level <= LogLevelWarn {
		l.logln("WARN", v...)
	}
}

func (l *Logger) Errorln(v ...interface{}) {
	if l.level <= LogLevelError {
		l.logln("ERROR", v...)
	}
}

func (l *Logger) log(level, format string, v ...interface{}) {
	text := level
	if l.option_include_dt {
		dt := time.Since(l.create)
		text = fmt.Sprintf("%s %s", text, dt)
	}
	prefix := l.GetPrefix()
	l.logger.Printf("%s: %s %s", text, prefix, fmt.Sprintf(format, v...))
}

func (l *Logger) logln(level string, v ...interface{}) {
	text := level
	if l.option_include_dt {
		dt := time.Since(l.create)
		text = fmt.Sprintf("%s %s", text, dt)
	}
	args := make([]string, len(v))
	for i, arg := range v {
		args[i] = fmt.Sprint(arg)
	}
	prefix := l.GetPrefix()
	l.logger.Printf("%s: %s %s", text, prefix, strings.Join(args, " "))
}
