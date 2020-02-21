package state

import "testing"

//tests the test stringer is correct
func TestState_String(t *testing.T) {
	//define some states to check
	expectedStateStrings := []string{"UNKNOWN STATE: 0", "NOT_STARTED",
		"WAITING", "PRECOMPUTING", "STANDBY", "REALTIME", "ERROR", "CRASH",
		"UNKNOWN STATE: 8"}

	//check if states give the correct return
	for st := State(0); st <= NUM_STATES; st++ {
		if st.String() != expectedStateStrings[st] {
			t.Errorf("State %d did not string correctly: expected: %s,"+
				"recieved: %s", uint8(st), expectedStateStrings[st], st.String())
		}
	}
}
