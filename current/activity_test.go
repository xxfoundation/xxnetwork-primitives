////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package current

import "testing"

//tests the test stringer is correct
func TestActivity_String(t *testing.T) {
	//define some activities to check
	expectedActivityStringer := []string{"NOT_STARTED",
		"WAITING", "PRECOMPUTING", "STANDBY", "REALTIME", "COMPLETED", "ERROR", "CRASH",
		"UNKNOWN STATE: 8"}

	//check if states give the correct return
	for st := NOT_STARTED; st <= NUM_STATES; st++ {
		if st.String() != expectedActivityStringer[st] {
			t.Errorf("Activity %d did not string correctly: expected: %s,"+
				"recieved: %s", uint8(st), expectedActivityStringer[st], st.String())
		}
	}
}
