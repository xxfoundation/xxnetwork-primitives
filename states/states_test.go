////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package states

import "testing"

// Consistency test of Round.String.
func TestRound_String(t *testing.T) {
	expected := []string{"PENDING", "PRECOMPUTING", "STANDBY", "QUEUED",
		"REALTIME", "COMPLETED", "FAILED", "UNKNOWN STATE: 7"}

	for st := PENDING; st <= NUM_STATES; st++ {
		if st.String() != expected[st] {
			t.Errorf("Incorrect string for Round state %d."+
				"\nexpected: %s\nreceived: %s", st, expected[st], st.String())
		}
	}
}
