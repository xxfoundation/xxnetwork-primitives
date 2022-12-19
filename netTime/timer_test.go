package netTime

import (
	"testing"
	"time"
)

// Tests that NewTimer produces a timer that works and where C is accessible.
func TestNewTimer(t *testing.T) {
	d := 20 * time.Millisecond
	timer := NewTimer(20 * time.Millisecond)

	select {
	case <-timer.C:
	case <-time.After(d):
		t.Errorf("Timed out waiting for timer.")
	}
}
