////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package current

import (
	"testing"

	"github.com/pkg/errors"

	"gitlab.com/elixxir/primitives/states"
)

// Consistency test of Activity.String.
func TestActivity_String(t *testing.T) {
	expected := []string{"NOT_STARTED", "WAITING", "PRECOMPUTING", "STANDBY",
		"REALTIME", "COMPLETED", "ERROR", "CRASH", "UNKNOWN STATE: 8"}

	for st := NOT_STARTED; st <= NUM_STATES; st++ {
		if st.String() != expected[st] {
			t.Errorf("Incorrect string for Activity %d."+
				"\nexpected: %s\nreceived: %s", st, expected[st], st.String())
		}
	}
}

// Test proper happy path of Activity.ConvertToRoundState.
func TestActivity_ConvertToRoundState(t *testing.T) {
	tests := []struct {
		activity Activity
		state    states.Round
		err      error
	}{
		{NOT_STARTED, 99, errors.Errorf("unable to convert activity %s (%d) "+
			"to valid state", NOT_STARTED, NOT_STARTED)},
		{WAITING, states.PENDING, nil},
		{PRECOMPUTING, states.PRECOMPUTING, nil},
		{STANDBY, states.STANDBY, nil},
		{REALTIME, states.REALTIME, nil},
		{COMPLETED, states.COMPLETED, nil},
		{ERROR, states.FAILED, nil},
		{CRASH, 99, errors.Errorf("unable to convert activity %s (%d) to "+
			"a valid state", CRASH, CRASH)},
	}

	for i, tt := range tests {
		state, err := tt.activity.ConvertToRoundState()
		if err != nil && tt.err == nil {
			if err.Error() != tt.err.Error() {
				t.Errorf(
					"Unexpected error for %s (%d).\nexpected: %s\nreceived: %s",
					tt.activity, i, tt.err, err)
			}
		} else if state != tt.state {
			t.Errorf("Unexpected conversation of %s (%d)."+
				"\nexpected: %s\nreceived: %s", tt.activity, i, tt.state, state)
		}
	}
}
