package log

import (
	"fmt"
	"sort"

	"github.com/rs/zerolog"
)

// EventReceiver is a sentinel EventReceiver uses zerolog to log events
type EventReceiver struct{}

// Event receives a simple notification when various events occur.
func (n *EventReceiver) Event(eventName string) {

}

// EventKv receives a notification when various events occur along with
// optional key/value data.
func (n *EventReceiver) EventKv(eventName string, kvs map[string]string) {

	k := writeMapToDict(kvs)

	Clogger.Info().Str("svc", "dbr").
		Str("method", eventName).
		Dict("kvs", k).
		Msg("EventEv")
}

// EventErr receives a notification of an error if one occurs.
func (n *EventReceiver) EventErr(eventName string, err error) error {

	Clogger.Error().Str("svc", "dbr").
		Str("method", eventName).
		Err(err).
		Msg("EventErr")

	return err
}

// EventErrKv receives a notification of an error if one occurs along with
// optional key/value data.
func (n *EventReceiver) EventErrKv(eventName string, err error, kvs map[string]string) error {

	k := writeMapToDict(kvs)

	Clogger.Error().Str("svc", "dbr").
		Str("method", eventName).
		Dict("kvs", k).
		Err(err).
		Msg("EventErrkv")

	return err
}

// Timing receives the time an event took to happen.
func (n *EventReceiver) Timing(eventName string, nanos int64) {

	Clogger.Info().Str("svc", "dbr").
		Int64("took", nanos).
		Str("method", eventName).
		Str("latency", writeNanoseconds(nanos)).
		Msg("Timing")
}

// TimingKv receives the time an event took to happen along with optional key/value data.
func (n *EventReceiver) TimingKv(eventName string, nanos int64, kvs map[string]string) {

	k := writeMapToDict(kvs)

	Clogger.Info().Str("svc", "dbr").
		Int64("took", nanos).
		Str("method", eventName).
		Str("latency", writeNanoseconds(nanos)).
		Dict("kvs", k).
		Msg("TimingKv")
}

func writeMapToDict(kvs map[string]string) *zerolog.Event {
	if kvs == nil {
		return nil
	}

	keys := make([]string, 0, len(kvs))
	for k := range kvs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := zerolog.Dict()
	for _, k := range keys {
		out.Str(k, kvs[k])
	}

	return out
}

func writeNanoseconds(nanos int64) string {
	switch {
	case nanos > 2000000:
		return fmt.Sprintf("%dms", nanos/1000000)
	case nanos > 2000:
		return fmt.Sprintf("%dμs", nanos/1000)
	default:
		return fmt.Sprintf("%dns", nanos)
	}

	return ""
}
