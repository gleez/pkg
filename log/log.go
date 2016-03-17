// Copyright (c) 2016 Gleez Technologies, Sandeep Sangamreddi, contributors
// The use of this source code is governed by a MIT style license found in the LICENSE file

package log

import (
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"
)

const format = "[%s] %s [%s:%d] %s\n" // level application [file:line]: message newline

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

// GoLog is a level logging facade for the Go stdlib log used as a default logger.
type GoLog struct {
	logger *stdlog.Logger
	level  Level
	color  bool
	name   string
	path   string
	file   *os.File
	signal bool
}

// NewGoLog constructs a new instance of GoLog.
func NewGoLog(logger *stdlog.Logger, level Level, name string, osSignal bool) *GoLog {
	glog := &GoLog{
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

// Set File creates or appends to the file path passed
func (gl *GoLog) SetFile(filePath string) error {
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
func (gl *GoLog) CloseFile() error {
	if gl.file != nil {
		gl.Trace("Closing log file")
		return gl.file.Close()
	}

	return nil
}

// SetLevel sets the logging level of the logger. Messages with the level lower than currently set
// are ignored. Messages with the level FATAL or PANIC are never ignored.
func (gl *GoLog) SetLevel(level Level) {
	gl.level = level
}

// GetLevel retrieves the currently set logger level.
func (gl *GoLog) GetLevel() Level {
	return gl.level
}

func (gl *GoLog) log(msg string, level Level, args ...interface{}) {
	if level < gl.level && level < PANIC {
		return
	}

	m := msg
	if len(args) > 0 {
		m = fmt.Sprintf(m, args...)
	}

	_, f, l, _ := runtime.Caller(2)

	//	if gl.color {
	//		m = fmt.Sprintf(colors[level]+format, levels[level], path.Base(f), l, "\033[0m"+m)
	//	} else {
	//		m = fmt.Sprintf(format, levels[level], path.Base(f), l, m)
	//	}

	output := fmt.Sprintf(format, levels[level], gl.name, path.Base(f), l, m)
	gl.logger.Output(2, output)

	switch level {
	case PANIC:
		panic(output)
	case FATAL:
		os.Exit(1)
	default:
		// ignore
	}
}

// Log logs a formatted message with a given level. If level is below the one currently set
// the message will be ignored.
func (self *GoLog) Log(message string, level Level, args ...interface{}) {
	self.log(message, level, args...)
}

// Trace logs a formatted message with the TRACE level.
func (self *GoLog) Trace(message string, args ...interface{}) {
	self.log(message, TRACE, args...)
}

// Debug logs a formatted message with the DEBUG level.
func (self *GoLog) Debug(message string, args ...interface{}) {
	self.log(message, DEBUG, args...)
}

// Info logs a formatted message with the INFO level.
func (self *GoLog) Info(message string, args ...interface{}) {
	self.log(message, INFO, args...)
}

// Notice logs a formatted message with the NOTICE level.
func (self *GoLog) Notice(message string, args ...interface{}) {
	self.log(message, NOTICE, args...)
}

// Warn logs a formatted message with the WARN level.
func (self *GoLog) Warn(message string, args ...interface{}) {
	self.log(message, WARN, args...)
}

// Error logs a formatted message with the ERROR level.
func (self *GoLog) Error(message string, args ...interface{}) {
	self.log(message, ERROR, args...)
}

// Panic logs a formatted message with the PANIC level and calls panic.
func (self *GoLog) Panic(message string, args ...interface{}) {
	self.log(message, PANIC, args...)
}

// Alert logs a formatted message with the ALERT level.
func (self *GoLog) Alert(message string, args ...interface{}) {
	self.log(message, ALERT, args...)
}

// Fatal logs a formatted message with the FATAL level and exits the application with an error.
func (self *GoLog) Fatal(message string, args ...interface{}) {
	self.log(message, FATAL, args...)
}

// StdLog returns the underlying stdlib logger.
func (gl *GoLog) StdLog() *stdlog.Logger {
	return gl.logger
}

// SignalProcessor is a goroutine that handles OS signals
func (gl *GoLog) SignalProcessor(c <-chan os.Signal) {
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
