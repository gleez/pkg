/*
	Copyright (c) 2016 Gleez Technologies, Sandeep Sangamreddi, contributors
	The use of this source code is governed by a MIT style license found in the LICENSE file

	Package log provides a generic interface for leveled logging and a default implementation
	of the interface as a facade to Go stdlib log.Logger.
*/

package log

import (
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var defaultTimeFormat = time.RFC3339 // 2006-01-02T15:04:05Z07:00
var format = "[%s] %s [%s:%d] %s\n"  // level application [file:line]: message newline

var log = New(stdlog.New(os.Stderr, "", stdlog.Ldate|stdlog.Ltime), INFO, "", false)

// Log provides a struct with fields that describe the details of log.
type Log struct {
	logger *stdlog.Logger
	level  Level
	name   string
	path   string
	file   *os.File
	signal bool
	mu     sync.Mutex
}

// NewLog constructs a new instance of Log.
func New(logger *stdlog.Logger, level Level, name string, osSignal bool) *Log {
	glog := &Log{
		logger: logger,
		level:  level,
		name:   name,
		signal: osSignal,
	}

	if glog.signal {
		// Setup signal handler
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan,
			os.Interrupt,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		go glog.SignalProcessor(sigChan)
	}

	return glog
}

// SetLevel sets the logging level of the logger. Messages with the level lower than currently set
// are ignored. Messages with the level FATAL or PANIC are never ignored.
func (l *Log) SetLevel(lvl Level) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.level = lvl
}

func (gl *Log) log(msg string, level Level, args ...interface{}) {
	if level < gl.level && level < PANIC {
		return
	}

	m := msg
	if len(args) > 0 {
		m = fmt.Sprintf(m, args...)
	}

	// Show file name and line in debug and trace
	if gl.level < INFO {
		m = fmt.Sprintf("[%s] %s", line(2), m)
	}

	m = fmt.Sprintf(colors[level]+"[%s] %s %s\n", levels[level], gl.name, "\033[0m"+m)
	gl.logger.Output(2, m)

	switch level {
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
func (self *Log) Log(message string, lvl Level, args ...interface{}) {
	self.log(message, lvl, args...)
}

// Trace logs a formatted message with the TRACE level.
func (self *Log) Trace(message string, args ...interface{}) {
	self.log(message, TRACE, args...)
}

// Debug logs a formatted message with the DEBUG level.
func (self *Log) Debug(message string, args ...interface{}) {
	self.log(message, DEBUG, args...)
}

// Info logs a formatted message with the INFO level.
func (self *Log) Info(message string, args ...interface{}) {
	self.log(message, INFO, args...)
}

// Notice logs a formatted message with the NOTICE level.
func (self *Log) Notice(message string, args ...interface{}) {
	self.log(message, NOTICE, args...)
}

// Warn logs a formatted message with the WARN level.
func (self *Log) Warn(message string, args ...interface{}) {
	self.log(message, WARN, args...)
}

// Error logs a formatted message with the ERROR level.
func (self *Log) Error(message string, args ...interface{}) {
	self.log(message, ERROR, args...)
}

// Panic logs a formatted message with the PANIC level and calls panic.
func (self *Log) Panic(message string, args ...interface{}) {
	self.log(message, PANIC, args...)
}

// Alert logs a formatted message with the ALERT level.
func (self *Log) Alert(message string, args ...interface{}) {
	self.log(message, ALERT, args...)
}

// Fatal logs a formatted message with the FATAL level and exits the application with an error.
func (self *Log) Fatal(message string, args ...interface{}) {
	self.log(message, FATAL, args...)
}

// StdLog returns the underlying stdlib logger.
func (gl *Log) StdLog() *stdlog.Logger {
	return gl.logger
}

// Set File creates or appends to the file path passed
func (gl *Log) SetFile(filePath string) error {
	var err error

	if gl.file == nil {
		// use specified log file
		gl.file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("Failed to create %s: %v\n", filePath, err)
		}

		// assign it to the standard logger
		stdlog.SetOutput(gl.file)

		gl.path = filePath
		gl.Trace("Opened log file")
	}

	return nil
}

// Close Log File closes the current logfile
func (gl *Log) CloseFile() error {
	if gl.file != nil {
		gl.Trace("Closing log file")
		return gl.file.Close()
	}

	return nil
}

// SignalProcessor is a goroutine that handles OS signals
func (gl *Log) SignalProcessor(c <-chan os.Signal) {
	for {
		sig := <-c
		fmt.Println("Got signal:", sig)
		switch sig {
		case syscall.SIGHUP:
			// Rotate logs if configured
			if gl.file != nil {
				gl.Info("Recieved SIGHUP, cycling logfile")
				gl.CloseFile()
				gl.SetFile(gl.path)
			} else {
				gl.Info("Ignoring SIGHUP, logfile not configured")
			}
		case syscall.SIGTERM:
		case syscall.SIGQUIT:
		case syscall.SIGINT:
			// Initiate shutdown
			gl.Info("Received signal, shutting down in 15 seconds")
			go func() {
				time.Sleep(15 * time.Second)
				gl.Info("Clean shutdown timed out, forcing exit")
				os.Exit(0)
			}()
		}
	}
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

func line(calldepth int) string {
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}

	return fmt.Sprintf("%s:%d", path.Base(file), line)
}

//SetLevel sets the level of default Logger
func SetLevel(lvl Level) {
	log.SetLevel(lvl)
}

// GetLogger retrieves a log instance.
func GetLogger() *Log {
	return log
}
