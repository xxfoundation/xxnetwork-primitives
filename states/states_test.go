////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package states

import "testing"

//tests the test stringer is correct
func TestActivity_String(t *testing.T) {
	//define some states to check
	expectedStateStringer := []string{"PENDING", "PRECOMPUTING", "STANDBY",
		"REALTIME", "COMPLETED", "FAILED", "UNKNOWN STATE: 6"}

	//check if states give the correct return
	for st := PENDING; st <= NUM_STATES; st++ {
		if st.String() != expectedStateStringer[st] {
			t.Errorf("Round %d did not string correctly: expected: %s,"+
				"recieved: %s", uint8(st), expectedStateStringer[st], st.String())
		}
	}
}
