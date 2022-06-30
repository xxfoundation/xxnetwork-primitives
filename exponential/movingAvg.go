////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                           //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file                                                               //
////////////////////////////////////////////////////////////////////////////////

package exponential

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"sync"
)

// MovingAvg tracks the exponential moving average across a number of events and
// reports when it has surpassed the set cutoff.
type MovingAvg struct {
	// The maximum limit for the moving average aN allowed.
	cutoff float32

	// A(n), the current exponential moving average (initialize to A(0)).
	aN float32

	// Exponential smoothing factor; gives the most recent events more weight.
	// The greater the smoothing factor, the greater the influence of more
	// recent events.
	s float32

	// The number of events to average over.
	e uint32

	sync.Mutex
}

// NewMovingAvg creates a new MovingAvg with the given cutoff, initial average,
// smoothing factor, and number of events to average.
func NewMovingAvg(p MovingAvgParams) *MovingAvg {
	jww.TRACE.Printf("[MAVG] Tracking new exponential moving average: %+v", p)
	return &MovingAvg{
		cutoff: p.Cutoff,
		aN:     p.InitialAverage,
		s:      p.SmoothingFactor,
		e:      p.NumberOfEvents,
	}
}

// Intake takes in the current average and calculates the exponential average
// returning true if it is over the cutoff and false otherwise.
//
// The moving average is calculated by:
//  A(n) = a × (S/E) + A(n-1) × (1 − S/E)
// Where:
//  A(n) is the current exponential moving average
//  A(n-1) is the previous exponential moving average
//  a is the intake value
//  S is the smoothing factor
//  E is the number of events the average is over
func (m *MovingAvg) Intake(a float32) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	// Calculate exponential moving average
	k := m.s / (1 + float32(m.e))
	m.aN = (a * k) + (m.aN * (1 - k))

	jww.TRACE.Printf(
		"[MAVG] Intake %.4f: new moving average %.2f%% over %d events",
		a, m.aN*100, m.e)

	if m.aN > m.cutoff {
		return errors.Errorf("exponential average for the last %d events of "+
			"%.2f%% went over cutoff %.2f%%", m.e, m.aN*100, m.cutoff*100)
	}
	return nil
}

// IsOverCutoff returns true if the average has reached the cutoff and false if
// it has not
func (m *MovingAvg) IsOverCutoff() bool {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	return m.aN > m.cutoff
}

// BoolToFloat returns 1 if true and 0 if false.
func BoolToFloat(b bool) float32 {
	if b {
		return 1
	}
	return 0
}
