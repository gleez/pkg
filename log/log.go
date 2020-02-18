// Package log provides a global logger for zerolog.
package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// Logger is the global logger.
var Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

//InitLog initializes global logging settings
// level - debug, info, warn, error, fatal, panic, none
func SetupLogging(level string, isConsole, isCaller bool) {
	//set log level
	if level == "none" || level == "disabled" {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	} else {
		lvl, err := zerolog.ParseLevel(level)
		if err != nil {
			fmt.Printf("invalid log level: %s\n", err)
			panic("log initialization failed")
		}
		zerolog.SetGlobalLevel(lvl)
	}

	//set log format
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.DurationFieldInteger = true

	//create logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	if isConsole {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	if isCaller {
		logger = logger.With().Caller().Logger()
	}

	Logger = logger
}

// Service logs using info - source, msg, error, params
func Service(ctx context.Context, source, msg string, err error, params map[string]interface{}) {

	if params == nil {
		params = make(map[string]interface{})
	}

	params["svc"] = source

	if source == "account" {
		params["svc"] = "user"
	}

	if val, ok := params["took"].(time.Duration); ok {
		params["duration"] = val
		params["latency"] = val.String()

		//remove took
		delete(params, "took")
	}

	// 	if user, ok := ctx.Value("_authuser").(*models.AuthUser); ok {
	// 		params["id"] = user.ID
	// 		params["username"] = user.Username
	// 		params["tenant_id"] = user.TenantID
	// 	}

	if err != nil {
		params["error"] = err
		Error().Fields(params).Msg(msg)
		return
	}

	Info().Fields(params).Msg(msg)

}

// Output duplicates the global logger and sets w as its output.
func Output(w io.Writer) zerolog.Logger {
	return Logger.Output(w)
}

// With creates a child logger with the field added to its context.
func With() zerolog.Context {
	return Logger.With()
}

// Level creates a child logger with the minimum accepted level set to level.
func Level(level zerolog.Level) zerolog.Logger {
	return Logger.Level(level)
}

// Sample returns a logger with the s sampler.
func Sample(s zerolog.Sampler) zerolog.Logger {
	return Logger.Sample(s)
}

// Hook returns a logger with the h Hook.
func Hook(h zerolog.Hook) zerolog.Logger {
	return Logger.Hook(h)
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return Logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return Logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *zerolog.Event {
	return Logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return Logger.Error()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return Logger.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *zerolog.Event {
	return Logger.Panic()
}

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
func WithLevel(level zerolog.Level) *zerolog.Event {
	return Logger.WithLevel(level)
}

// Log starts a new message with no level. Setting zerolog.GlobalLevel to
// zerolog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func Log() *zerolog.Event {
	return Logger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	Logger.Print(v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	Logger.Printf(format, v...)
}

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
