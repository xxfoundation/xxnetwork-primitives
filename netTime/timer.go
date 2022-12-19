////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package netTime

import "time"

// The Timer type represents a single event. When the Timer expires, the current
// time will be sent on C, unless the Timer was created by AfterFunc. A Timer
// must be created with NewTimer or AfterFunc.
type Timer struct {
	*time.Timer
}

// NewTimer creates a new [Timer] that will send the current time on its channel
// after at least duration d. This wraps a [time.Timer] to provide a timer
// accurate to the set time source.
func NewTimer(d time.Duration) *Timer {
	return &Timer{time.NewTimer(d + getOffset())}
}

// Reset changes the timer to expire after duration d. It returns true if the
// timer had been active, false if the timer had expired or been stopped. Refer
// to [time/Timer.Reset] for more information.
func (t *Timer) Reset(d time.Duration) bool {
	return t.Timer.Reset(d + getOffset())
}

// After waits for the duration to elapse and then sends the current time on the
// returned channel. It is equivalent to NewTimer(d).C. Refer to [time.After]
// for more information.
func After(d time.Duration) <-chan time.Time {
	return time.After(d + getOffset())
}

// AfterFunc waits for the duration to elapse and then calls f in its own
// goroutine. It returns a [Timer] that can be used to cancel the call using its
// Stop method. Refer to [time.AfterFunc] for more information.
func AfterFunc(d time.Duration, f func()) *Timer {
	return &Timer{time.AfterFunc(d+getOffset(), f)}
}
