package notifications

import "testing"

// Unit test the Provider.String function.
func TestType_String(t *testing.T) {
	//define some activities to check
	expectedTypeStringer := []string{"APNS", "FCM", "HUAWEI"}

	//check if states give the correct return
	for st := APNS; st <= HUAWEI; st++ {
		if st.String() != expectedTypeStringer[st] {
			t.Errorf("Provider %d did not string correctly"+
				"\nExpected: %s,"+
				"\nReceived: %s", uint8(st), expectedTypeStringer[st], st.String())
		}
	}
}
