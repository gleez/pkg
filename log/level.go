/*
	Copyright (c) 2016 Gleez Technologies, Sandeep Sangamreddi, contributors
	The use of this source code is governed by a MIT style license found in the LICENSE file
*/

package log

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

var levels = []string{
	TRACE:  "TRACE",
	DEBUG:  "DEBUG",
	INFO:   "INFO",
	NOTICE: "NOTICE",
	WARN:   "WARN",
	ERROR:  "ERROR",
	PANIC:  "PANIC",
	ALERT:  "ALERT",
	FATAL:  "FATAL",
}

var colors = []string{
	TRACE:  "\033[36m;1m", // Cyan Bold
	DEBUG:  "\033[36m",    // Cyan
	INFO:   "\033[32m",    // Green
	NOTICE: "\033[34m",    // Bluew
	WARN:   "\033[33;1m",  // Yellow
	ERROR:  "\033[31m",    // Red
	PANIC:  "\033[31;1m",  // Red Bold
	ALERT:  "\033[35m",    // Magenta
	FATAL:  "\033[35;1m",  // Magenta Bold
}
