///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package netTime

import (
	"testing"
	"time"
)

// Happy path: tests that Now() returns time.Now() if it is unset.
func TestNow(t *testing.T) {
	expectedTime := time.Now().Round(time.Millisecond)
	receivedTime := Now().Round(time.Millisecond)

	if !expectedTime.Equal(receivedTime) {
		t.Errorf("Returned incorrect time.\nexpected: %s\nreceived: %s",
			expectedTime, receivedTime)
	}
}

// Happy path: tests that setting Now works.
func TestNow_Set(t *testing.T) {
	expectedTime := time.Now()
	testNow := func() time.Time {
		return expectedTime
	}

	Now = testNow

	now := Now()
	if !Now().Equal(expectedTime) {
		t.Errorf("Returned incorrect time.\nexpected: %s\nreceived: %s",
			expectedTime, now)
	}
}
