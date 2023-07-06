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

// Tests that Now returns time.Now if it is unset.
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

// Tests that SetOffset modifies the offset variable.
func TestSetOffset(t *testing.T) {
	expected := int64(25)
	newOffset := time.Duration(expected)
	SetOffset(newOffset)

	if offset != expected {
		t.Fatalf("SetOffset failed to set the offset variable."+
			"\nexpected: %d\nreceived: %d", expected, offset)
	}
}

// Tests that calling Now with an offset returns a time.Time with that offset
// applied. This test handles a positive offset. There exists a test for a
// negative offset in this file.
func TestNow_SetOffset_Positive(t *testing.T) {
	// Set the time source to return to a hard coded value for easy testing
	testTime := time.Date(1955, 11, 5, 12, 5, 0, 0, time.UTC)
	mockSource := &mockTimeSource{returnTime: testTime.UnixMilli()}
	SetTimeSource(mockSource)

	// Set an offset (positive value for this test)
	newOffset := 5 * time.Second
	SetOffset(newOffset)

	// Retrieve the time derived from the time source
	received := Now().UnixNano()

	// The expected value should be the hardcoded time added to the offset
	expected := testTime.Add(newOffset).UnixNano()

	// Ensure expected value matches received value
	if received != expected {
		t.Fatalf("Now did not return a time adjusted for the offset."+
			"\nexpected: %d\nreceived: %d", expected, received)
	}
}

// Tests that calling Now with an offset returns a time.Time with that offset
// applied. This test handles a negative offset. There exists a test for a
// positive offset in this file.
func TestNow_SetOffset_Negative(t *testing.T) {
	// Set the time source to return to a hard coded value for easy testing
	testTime := time.Date(1955, 11, 5, 12, 5, 0, 0, time.UTC)
	mockSource := &mockTimeSource{returnTime: testTime.UnixMilli()}
	SetTimeSource(mockSource)

	// Set an offset (negative value for this test)
	newOffset := -5 * time.Second
	SetOffset(newOffset)
	received := Now().UnixNano()

	// The expected value should be the hardcoded time added to the offset
	expected := testTime.Add(newOffset).UnixNano()

	// Ensure expected value matches received value
	if received != expected {
		t.Fatalf("Now did not return a time adjusted for the offset."+
			"\nexpected: %d\nreceived: %d", expected, received)
	}
}
