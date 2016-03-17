// Copyright (c) 2016 Gleez Technologies, Sandeep Sangamreddi, contributors
// The use of this source code is governed by a MIT style license found in the LICENSE file

// Package log provides a generic interface for leveled logging and a default implementation
// of the interface as a facade to Go stdlib log.Logger.
package log

import (
	stdlog "log"
	"os"
)

// Level defines the log level. The lowest is TRACE and the highest are FATAL and PANIC (treated as
// being of the same level). Messages with the level lower than currently set are ignored.
// Messages with the level FATAL or PANIC are never ignored.
type Level int

const (
	// TRACE log level (TRACE < DEBUG).
	TRACE Level = iota
	// DEBUG log level (TRACE < DEBUG < INFO).
	DEBUG
	// INFO log level (DEBUG < INFO < WARN).
	INFO
	// NOTICE log level (DEBUG < INFO < WARN).
	NOTICE
	// WARN log level (NOTICE < WARN < ERROR).
	WARN
	// ERROR log level (WARN < ERROR < PANIC/FATAL).
	ERROR
	// PANIC log level (ERROR < PANIC), application panics.
	PANIC
	// ALERT log level (ERROR < ALERT), application panics.
	ALERT
	// FATAL log level (ERROR < FATAL), application exits with 1.
	FATAL
)

var log Logger

func init() {
	log = NewGoLog(stdlog.New(os.Stderr, "", stdlog.Ldate|stdlog.Ltime|stdlog.LUTC), INFO, "", false)
}

// Logger represents a generic leveled logger interface.
type Logger interface {
	// SetLevel sets the logging level of the logger. Messages with the level lower than currently set
	// are ignored. Messages with the level FATAL or PANIC are never ignored.
	SetLevel(level Level)

	// GetLevel retrieves the currently set logger level.
	GetLevel() Level

	// Log logs a formatted message with a given level. If level is below the one currently set
	// the message will be ignored.
	Log(message string, level Level, args ...interface{})

	// Trace logs a formatted message with the TRACE level.
	Trace(message string, args ...interface{})

	// Debug logs a formatted message with the DEBUG level.
	Debug(message string, args ...interface{})

	// Info logs a formatted message with the INFO level.
	Info(message string, args ...interface{})

	// Notice logs a formatted message with the NOTICE level.
	Notice(message string, args ...interface{})

	// Warn logs a formatted message with the WARN level.
	Warn(message string, args ...interface{})

	// Error logs a formatted message with the ERROR level.
	Error(message string, args ...interface{})

	// Panic logs a formatted message with the PANIC level and calls panic.
	Panic(message string, args ...interface{})

	// Alert logs a formatted message with the ALERT level.
	Alert(message string, args ...interface{})

	// Fatal logs a formatted message with the FATAL level and exits the application with an error.
	Fatal(message string, args ...interface{})
}

// SetLogger replaces the default logger with a third-party implementation.
func SetLogger(logger Logger) {
	log = logger
}

// GetLogger retrieves a logger for the given topic.
func GetLogger(topic string) Logger {
	return log
}
