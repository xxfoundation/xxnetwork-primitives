package state

import "testing"

//tests the test stringer is correct
func TestState_String(t *testing.T) {
	//define some states to check
	expectedStateStrings := []string{"NOT_STARTED",
		"WAITING", "PRECOMPUTING", "STANDBY", "REALTIME", "ERROR", "CRASH",
		"UNKNOWN STATE: 7"}

	//check if states give the correct return
	for st := NOT_STARTED; st <= NUM_STATES; st++ {
		if st.String() != expectedStateStrings[st] {
			t.Errorf("State %d did not string correctly: expected: %s,"+
				"recieved: %s", uint8(st), expectedStateStrings[st], st.String())
		}
	}
}
