/*
	Copyright (c) 2016 Gleez Technologies, Sandeep Sangamreddi, contributors
	The use of this source code is governed by a MIT style license found in the LICENSE file

	Package log provides a generic interface for leveled logging and a default implementation
	of the interface as a facade to Go stdlib log.Logger.
*/

package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

var defaultTimeFormat = time.RFC3339 // 2006-01-02T15:04:05Z07:00

var log = New(os.Stdout, INFO)

// Logger provides a struct with fields that describe the details of log.
type Logger struct {
	mu     sync.Mutex // ensures atomic writes; protects the following fields
	prefix string     // prefix to write at beginning of each line
	out    io.Writer  // destination for output
	level  Level
	file   *os.File
}

// Constructs a new instance of Log.
func New(out io.Writer, lvl Level) *Logger {
	return &Logger{out: out, level: lvl}
}

// Prefix returns the output prefix for the logger.
func (l *Logger) Prefix() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.prefix
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// SetLevel sets the logging level of the logger. Messages with the level lower than currently set
// are ignored. Messages with the level FATAL or PANIC are never ignored.
func (l *Logger) SetLevel(lvl Level) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.level = lvl
}

// Set File creates or appends to the file path passed
func (l *Logger) SetFile(filePath string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		// use specified log file
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("Failed to create %s: %v\n", filePath, err)
		}

		// assign it to the standard logger
		l.out = file
		l.file = file

		Debug("Opened log file")
	}

	return nil
}

// Close Log File closes the current logfile
func (l *Logger) CloseFile() error {
	if l.file != nil {
		Debug("Closing log file")
		return l.file.Close()
	}

	return nil
}

// Format formats the logs as "time [level] line message"
func (l *Logger) format(r *Record) (b []byte, err error) {
	s := fmt.Sprintf(colors[r.Lvl]+"%s [%s] %s", r.Time.Format(defaultTimeFormat), levels[r.Lvl], "\033[0m")

	// Show file name and line in debug and trace
	if l.level < INFO {
		if len(r.Line) != 0 {
			s = s + "[" + r.Line + "] "
		}
	}

	if len(r.Msg) != 0 {
		s = s + r.Msg
	}

	b = []byte(s)

	if len(b) == 0 || b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}

	return b, nil
}

func (l *Logger) output(record *Record) error {
	b, err := l.format(record)
	if err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	_, err = l.out.Write(b)

	return err
}

func (l *Logger) log(msg string, lvl Level, args ...interface{}) {
	if lvl < l.level && lvl < PANIC {
		return
	}

	m := msg
	if len(args) > 0 {
		m = fmt.Sprintf(m, args...)
	}

	record := NewRecord(time.Now(), m, line(4), lvl)
	l.output(record)

	switch lvl {
	case PANIC:
		panic(m)
	case FATAL:
		os.Exit(1)
	default:
		// ignore
	}
}

// Log logs a formatted message with a given level. If level is below the one currently set
// the message will be ignored.
func (l *Logger) Log(message string, lvl Level, args ...interface{}) {
	l.log(message, lvl, args...)
}

// Trace logs a formatted message with the TRACE level.
func (l *Logger) Trace(message string, args ...interface{}) {
	l.log(message, TRACE, args...)
}

// Debug logs a formatted message with the DEBUG level.
func (l *Logger) Debug(message string, args ...interface{}) {
	l.log(message, DEBUG, args...)
}

// Info logs a formatted message with the INFO level.
func (l *Logger) Info(message string, args ...interface{}) {
	l.log(message, INFO, args...)
}

// Notice logs a formatted message with the NOTICE level.
func (l *Logger) Notice(message string, args ...interface{}) {
	l.log(message, NOTICE, args...)
}

// Warn logs a formatted message with the WARN level.
func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(message, WARN, args...)
}

// Error logs a formatted message with the ERROR level.
func (l *Logger) Error(message string, args ...interface{}) {
	l.log(message, ERROR, args...)
}

// Panic logs a formatted message with the PANIC level and calls panic.
func (l *Logger) Panic(message string, args ...interface{}) {
	l.log(message, PANIC, args...)
}

// Alert logs a formatted message with the ALERT level.
func (l *Logger) Alert(message string, args ...interface{}) {
	l.log(message, ALERT, args...)
}

// Fatal logs a formatted message with the FATAL level and exits the application with an error.
func (l *Logger) Fatal(message string, args ...interface{}) {
	l.log(message, FATAL, args...)
}

// Trace logs a formatted message with the TRACE level.
func Trace(message string, args ...interface{}) {
	log.Trace(message, args...)
}

// Debug logs a formatted message with the DEBUG level.
func Debug(message string, args ...interface{}) {
	log.Debug(message, args...)
}

// Info logs a formatted message with the INFO level.
func Info(message string, args ...interface{}) {
	log.Info(message, args...)
}

// Notice logs a formatted message with the NOTICE level.
func Notice(message string, args ...interface{}) {
	log.Notice(message, args...)
}

// Warn logs a formatted message with the WARN level.
func Warn(message string, args ...interface{}) {
	log.Warn(message, args...)
}

// Error logs a formatted message with the ERROR level.
func Error(message string, args ...interface{}) {
	log.Error(message, args...)
}

// Panic logs a formatted message with the PANIC level and calls panic.
func Panic(message string, args ...interface{}) {
	log.Panic(message, args...)
}

// Alert logs a formatted message with the ALERT level.
func Alert(message string, args ...interface{}) {
	log.Alert(message, args...)
}

// Fatal logs a formatted message with the FATAL level and exits the application with an error.
func Fatal(message string, args ...interface{}) {
	log.Fatal(message, args...)
}

//SetLevel sets the level of default Logger
func SetLevel(lvl Level) {
	log.SetLevel(lvl)
}

// Prefix returns the output prefix for the standard logger.
func Prefix() string {
	return log.Prefix()
}

// SetPrefix sets the output prefix for the standard logger.
func SetPrefix(prefix string) {
	log.SetPrefix(prefix)
}

// Set File creates or appends to the file path passed
func SetFile(filePath string) error {
	return log.SetFile(filePath)
}

// Close Log File closes the current logfile
func CloseFile() error {
	return log.CloseFile()
}

func line(calldepth int) string {
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}

	return fmt.Sprintf("%s:%d", path.Base(file), line)
}
