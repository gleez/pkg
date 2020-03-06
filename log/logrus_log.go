package log

/*
 * `zlog` is a simple implementation of `grpclog.LoggerV2` interface using `zerolog`.
 * Use this to log the internal actions of a gRPC server or client.
 *
 * Example:
 *
 *  logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
 *  logger = logger.With().Str("component", "client-grpc").Logger()
 *
 *  grpclog.SetLoggerV2(log.NewLogrus(logger))
 *
 */
import (
	"fmt"

	"github.com/rs/zerolog"
)

type LogrusLogger interface {
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnln(v ...interface{})
	Warning(v ...interface{})
	Warningf(format string, v ...interface{})
	Warningln(v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Trace(v ...interface{})
	Tracef(format string, v ...interface{})
	Traceln(v ...interface{})
	V(verbosity int) bool
	SetLevel(lvl zerolog.Level)
	GetLevel() zerolog.Level
	WithField(key string, value interface{}) zerolog.Context
	WithFields(fields map[string]interface{}) zerolog.Context
}

func NewLogrus(log zerolog.Logger) LogrusLogger {
	return &logrusLogger{log: log}
}

type logrusLogger struct {
	log zerolog.Logger
}

func (l *logrusLogger) Panic(args ...interface{}) {
	l.log.Panic().Msg(fmt.Sprint(args...))
}

func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	l.log.Panic().Msg(fmt.Sprintf(format, args...))
}

func (l *logrusLogger) Panicln(args ...interface{}) {
	l.Panic(args...)
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	l.log.Fatal().Msg(fmt.Sprint(args...))
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatal().Msg(fmt.Sprintf(format, args...))
}

func (l *logrusLogger) Fatalln(args ...interface{}) {
	l.Fatal(args...)
}

func (l *logrusLogger) Error(args ...interface{}) {
	l.log.Error().Msg(fmt.Sprint(args...))
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.log.Error().Msg(fmt.Sprintf(format, args...))
}

func (l *logrusLogger) Errorln(args ...interface{}) {
	l.Error(args...)
}

func (l *logrusLogger) Warning(args ...interface{}) {
	l.log.Warn().Msg(fmt.Sprint(args...))
}

func (l *logrusLogger) Warningf(format string, args ...interface{}) {
	l.log.Warn().Msg(fmt.Sprintf(format, args...))
}

func (l *logrusLogger) Warningln(args ...interface{}) {
	l.Warning(args...)
}

func (l *logrusLogger) Warn(args ...interface{}) {
	l.log.Warn().Msg(fmt.Sprint(args...))
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.log.Warn().Msg(fmt.Sprintf(format, args...))
}

func (l *logrusLogger) Warnln(args ...interface{}) {
	l.Warn(args...)
}

func (l *logrusLogger) Info(args ...interface{}) {
	l.log.Info().Msg(fmt.Sprint(args...))
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.log.Info().Msg(fmt.Sprintf(format, args...))
}

func (l *logrusLogger) Infoln(args ...interface{}) {
	l.Info(args...)
}

func (l *logrusLogger) Debug(args ...interface{}) {
	l.log.Debug().Msg(fmt.Sprint(args...))
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.log.Debug().Msg(fmt.Sprintf(format, args...))
}

func (l *logrusLogger) Debugln(args ...interface{}) {
	l.Debug(args...)
}

func (l *logrusLogger) Print(args ...interface{}) {
	l.log.Info().Msg(fmt.Sprint(args...))
}

func (l *logrusLogger) Printf(format string, args ...interface{}) {
	l.log.Info().Msg(fmt.Sprintf(format, args...))
}

func (l *logrusLogger) Println(args ...interface{}) {
	l.Print(args...)
}

func (l *logrusLogger) Trace(args ...interface{}) {
	l.log.Trace().Msg(fmt.Sprint(args...))
}

func (l *logrusLogger) Tracef(format string, args ...interface{}) {
	l.log.Trace().Msg(fmt.Sprintf(format, args...))
}

func (l *logrusLogger) Traceln(args ...interface{}) {
	l.Trace(args...)
}

func (l *logrusLogger) V(verbosity int) bool {
	// verbosity values:
	// 0 = info
	// 1 = warning
	// 2 = error
	// 3 = fatal
	switch l.log.GetLevel() {
	case zerolog.PanicLevel:
		return verbosity > 3
	case zerolog.FatalLevel:
		return verbosity == 3
	case zerolog.ErrorLevel:
		return verbosity == 2
	case zerolog.WarnLevel:
		return verbosity == 1
	case zerolog.InfoLevel:
		return verbosity == 0
	case zerolog.DebugLevel:
		return true
	case zerolog.TraceLevel:
		return true
	default:
		return false
	}
}

// Level creates a child logger with the minimum accepted level set to level.
func (l *logrusLogger) SetLevel(lvl zerolog.Level) {
	_ = l.log.Level(lvl)
}

// GetLevel returns the current Level of l.
func (l *logrusLogger) GetLevel() zerolog.Level {
	return l.log.GetLevel()
}

// WithField allocates a new entry and adds a field to it.
// Debug, Print, Info, Warn, Error, Fatal or Panic must be then applied to
// this new returned entry.
// If you want multiple fields, use `WithFields`.
func (l *logrusLogger) WithField(key string, value interface{}) zerolog.Context {
	return l.log.With().Interface(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (l *logrusLogger) WithFields(fields map[string]interface{}) zerolog.Context {
	return l.log.With().Fields(fields)
}
