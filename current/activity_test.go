////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package current

import (
	"gitlab.com/elixxir/primitives/states"
	"testing"
)

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

// Test proper happy path conversions
func TestActivity_ConvertToRoundState(t *testing.T) {
	activity := NOT_STARTED
	state, err := activity.ConvertToRoundState()
	if err == nil {
		t.Errorf("Expected error when converting %+v", activity)
	}
	activity = WAITING
	state, err = activity.ConvertToRoundState()
	if err != nil {
		t.Errorf("Invalid conversion: %+v", err)
	}
	if state != states.PENDING {
		t.Errorf("Attempted to convert %+v. Expected %+v, got %+v",
			activity, states.PENDING, state)
	}
	activity = PRECOMPUTING
	state, err = activity.ConvertToRoundState()
	if err != nil {
		t.Errorf("Invalid conversion: %+v", err)
	}
	if state != states.PRECOMPUTING {
		t.Errorf("Attempted to convert %+v. Expected %+v, got %+v",
			activity, states.PRECOMPUTING, state)
	}
	activity = STANDBY
	state, err = activity.ConvertToRoundState()
	if err != nil {
		t.Errorf("Invalid conversion: %+v", err)
	}
	if state != states.STANDBY {
		t.Errorf("Attempted to convert %+v. Expected %+v, got %+v",
			activity, states.STANDBY, state)
	}
	activity = REALTIME
	state, err = activity.ConvertToRoundState()
	if err != nil {
		t.Errorf("Invalid conversion: %+v", err)
	}
	if state != states.REALTIME {
		t.Errorf("Expected %+v, got %+v", states.REALTIME, state)
	}
	activity = COMPLETED
	state, err = activity.ConvertToRoundState()
	if err != nil {
		t.Errorf("Invalid conversion: %+v", err)
	}
	if state != states.COMPLETED {
		t.Errorf("Attempted to convert %+v. Expected %+v, got %+v",
			activity, states.COMPLETED, state)
	}
	activity = ERROR
	state, err = activity.ConvertToRoundState()
	if err != nil {
		t.Errorf("Invalid conversion: %+v", err)
	}
	if state != states.FAILED {
		t.Errorf("Attempted to convert %+v. Expected %+v, got %+v",
			activity, states.FAILED, state)
	}
	activity = CRASH
	state, err = activity.ConvertToRoundState()
	if err == nil {
		t.Errorf("Expected error when converting %+v", activity)
	}
}
