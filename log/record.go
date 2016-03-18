/*
	Copyright (c) 2016 Gleez Technologies, Sandeep Sangamreddi, contributors
	The use of this source code is governed by a MIT style license found in the LICENSE file
*/

package log

import (
	"time"
)

// Record holds information about log
type Record struct {
	Time time.Time // time when the log produced
	Msg  string    // content of the log
	Line string    // in which file and line that the log produced
	Lvl  Level     // level of the log
}

// NewRecord creates a record according to the arguments provided and returns it
func NewRecord(time time.Time, msg, line string, lvl Level) *Record {
	return &Record{
		Time: time,
		Msg:  msg,
		Line: line,
		Lvl:  lvl,
	}
}
