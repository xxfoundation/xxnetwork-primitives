package states

import "testing"

//tests the test stringer is correct
func TestActivity_String(t *testing.T) {
	//define some states to check
	expectedStateStringer := []string{"PRECOMPUTING", "STANDBY", "REALTIME",
		"COMPLETED", "FAILED", "UNKNOWN STATE: 5"}

	//check if states give the correct return
	for st := PRECOMPUTING; st <= NUM_STATES; st++ {
		if st.String() != expectedStateStringer[st] {
			t.Errorf("State %d did not string correctly: expected: %s,"+
				"recieved: %s", uint8(st), expectedStateStringer[st], st.String())
		}
	}
}
