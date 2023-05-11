package ding

import (
	"fmt"
	"os"
)

// LeveledLogger is an interface that can be implemented by any logger or a
// logger wrapper to provide leveled logging. The methods accept a message
// string and a variadic number of key-value pairs.
type LeveledLogger interface {
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelNull
)

// DefaultLeveledLogger is the default logger that the library will use to log
// errors, warnings, and informational messages.
//
// LeveledLoggerInterface is implemented by LeveledLogger, and one can be
// initialized at the desired level of logging.  LeveledLoggerInterface also
// provides out-of-the-box compatibility with a Logrus Logger, but may require
// a thin shim for use with other logging libraries that use less standard
// conventions like Zap.
//
// This Logger will be inherited by any backends created by default, but will
// be overridden if a backend is created with GetBackendWithConfig with a
// custom LeveledLogger set.
var DefaultLeveledLogger Logger = Logger{
	Level: LevelError,
}

// Logger is a leveled logger implementation.
//
// It prints warnings and errors to `os.Stderr` and other messages to
// `os.Stdout`.
type Logger struct {
	// Level is the minimum logging level that will be emitted by this logger.
	//
	// For example, a Level set to LevelWarn will emit warnings and errors, but
	// not informational or debug messages.
	//
	// Always set this with a constant like LevelWarn because the individual
	// values are not guaranteed to be stable.
	Level Level
}

// Debugf logs a debug message using Printf conventions.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.Level >= LevelDebug {
		fmt.Fprintf(os.Stdout, "[DEBUG] "+format+"\n", v...)
	}
}

// Errorf logs a warning message using Printf conventions.
func (l *Logger) Errorf(format string, v ...interface{}) {
	// Infof logs a debug message using Printf conventions.
	if l.Level >= LevelError {
		fmt.Fprintf(os.Stderr, "[ERROR] "+format+"\n", v...)
	}
}

// Infof logs an informational message using Printf conventions.
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.Level >= LevelInfo {
		fmt.Fprintf(os.Stdout, "[INFO] "+format+"\n", v...)
	}
}

// Warnf logs a warning message using Printf conventions.
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.Level >= LevelWarn {
		fmt.Fprintf(os.Stderr, "[WARN] "+format+"\n", v...)
	}
}
