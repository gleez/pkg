package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	SetupLogging("debug", true, true)

	Debug().Msg("debug")
	Info().Msg("info")
	Warn().Msg("warn")
	Error().Msg("error")

	//these would fail the test
	// Log.Fatal().Msg("fatal")
	// Log.Panic().Msg("panic")
}
