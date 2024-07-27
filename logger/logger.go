package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type Logger struct {
	logger *log.Logger
	level  int
	prefix string
	mu     sync.RWMutex
}

func NewLogger(level int, prefix string) *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds),
		level:  level,
		prefix: prefix,
	}
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
	prefix := l.GetPrefix()
	l.logger.Printf("%s: %s %s", level, prefix, fmt.Sprintf(format, v...))
}

func (l *Logger) logln(level string, v ...interface{}) {
	args := make([]string, len(v))
	for i, arg := range v {
		args[i] = fmt.Sprint(arg)
	}
	prefix := l.GetPrefix()
	l.logger.Printf("%s: %s %s", level, prefix, strings.Join(args, " "))
}
