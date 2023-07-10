////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Package netTime provides a custom time function that should provide the
// current accurate time used by the network from a custom time service.
package netTime

import (
	"sync/atomic"
	"time"
)

// NowFunc is defined as the function interface for [time.Now].
type NowFunc func() time.Time

// Now returns the current accurate time. The function must be set an accurate
// time service that returns the current time with an accuracy of +/- 300 ms.
var Now NowFunc = time.Now

// offset is an internal variable which will be applied to the result of every
// call to Now. This is set using the SetOffset call.
var offset = int64(0)

// TimeSource is an interface which matches a time service that may be used
// to set [Now].
type TimeSource interface {
	NowMs() int64
}

// SetTimeSource sets [Now] to a custom source. All calls to [Now] will use this
// [TimeSource]. Note that this is in-memory, so any restart will require that
// this function be recalled.
func SetTimeSource(nowFunc TimeSource) {
	Now = func() time.Time {
		// Get the current time using the passed in time source
		currentTime := nowFunc.NowMs() * int64(time.Millisecond)

		// Parse the time into Golang's library (time.Time)
		parsedTime := time.Unix(0, currentTime)

		// Add the offset to the time retrieved via the time source
		return parsedTime.Add(getOffset())
	}
}

// SetOffset sets the internal offset variable atomically. All calls to [Now]
// will have this offset added to the result. Negative offsets are accepted and
// will reduce the result of the call to [Now].
func SetOffset(timeToOffset time.Duration) {
	atomic.StoreInt64(&offset, timeToOffset.Nanoseconds())
}

// getOffset returns the offset duration. This function is thread safe.
func getOffset() time.Duration {
	return time.Duration(atomic.LoadInt64(&offset)) * time.Nanosecond
}

// Since returns the time elapsed since t. It is shorthand for:
//
//	netTime.Now().Sub(t).
func Since(t time.Time) time.Duration {
	return Now().Sub(t)
}

// Until returns the duration until t. It is shorthand for:
//
//	t.Sub(netTime.Now()).
func Until(t time.Time) time.Duration {
	return t.Sub(Now())
}
