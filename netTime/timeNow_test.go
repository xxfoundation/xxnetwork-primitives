////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package netTime

import (
	"testing"
	"time"
)

// Tests that Now() returns time.Now() if it is unset.
func TestNow(t *testing.T) {
	expectedTime := time.Now().Round(time.Millisecond)
	receivedTime := Now().Round(time.Millisecond)

	if !expectedTime.Equal(receivedTime) {
		t.Errorf("Returned incorrect time.\nexpected: %s\nreceived: %s",
			expectedTime, receivedTime)
	}
}

// Tests that setting Now to a custom function results in the expected time.
func TestNow_Set(t *testing.T) {
	expectedTime := time.Date(1955, 11, 5, 12, 0, 0, 0, time.UTC)
	Now = func() time.Time { return expectedTime }

	if !Now().Equal(expectedTime) {
		t.Errorf("Now returned incorrect time.\nexpected: %s\nreceived: %s",
			expectedTime, Now())
	}
}

// Test that Since returns the expected duration.
func TestSince(t *testing.T) {
	expectedDuration := 24 * time.Hour
	testTime := time.Date(1955, 11, 5, 12, 0, 0, 0, time.UTC)
	Now = func() time.Time { return testTime }

	timeSince := Since(testTime.Add(-expectedDuration))
	if expectedDuration != timeSince {
		t.Errorf("Since returned incorrect duration."+
			"\nexpected: %s\nreceived: %s", expectedDuration, timeSince)
	}
}

// Test that Until returns the expected duration.
func TestUntil(t *testing.T) {
	expectedDuration := 24 * time.Hour
	testTime := time.Date(1955, 11, 5, 12, 0, 0, 0, time.UTC)
	Now = func() time.Time { return testTime }

	timeUntil := Until(testTime.Add(expectedDuration))
	if expectedDuration != timeUntil {
		t.Errorf("Until returned incorrect duration."+
			"\nexpected: %s\nreceived: %s", expectedDuration, timeUntil)
	}
}
