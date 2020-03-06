package log

import (
	"fmt"
	stdlog "log"
	"os"

	"github.com/rs/zerolog"
)

// Wraps a zerolog.Logger so its interoperable with Go's standard "log" package

var _ StdLogger = &stdlog.Logger{}

type StdLogger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

func NewStd(log zerolog.Logger) StdLogger {
	return &stdLogger{log}
}

type stdLogger struct {
	log zerolog.Logger
}

func (s *stdLogger) Fatal(v ...interface{}) {
	s.log.Fatal().Msg(fmt.Sprint(v...))
	os.Exit(1)
}

func (s *stdLogger) Fatalf(format string, v ...interface{}) {
	s.log.Fatal().Msg(fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (s *stdLogger) Print(v ...interface{}) {
	s.log.Info().Msg(fmt.Sprint(v...))
}

func (s *stdLogger) Println(v ...interface{}) {
	s.log.Info().Msg(fmt.Sprintln(v...))
}

func (s *stdLogger) Printf(format string, v ...interface{}) {
	s.log.Info().Msg(fmt.Sprintf(format, v...))
}
