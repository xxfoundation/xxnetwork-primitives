////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"bytes"
	"encoding/hex"
	"gitlab.com/elixxir/primitives/userid"
	"testing"
)

func TestMessagePayload(t *testing.T) {
	tests := 3

	testStrings := [][]byte{
		testText[0 : DATA_LEN/2],
		testText[0:DATA_LEN],
		testText[0 : 2*DATA_LEN],
	}

	expectedSlices := make([][]byte, tests)

	expectedSlices[0] = []byte(testStrings[0])
	expectedSlices[1] = []byte(testStrings[1])[0:DATA_LEN]
	expectedSlices[2] = []byte(testStrings[2])[0:DATA_LEN]

	expectedErrors := []bool{false, false, true}

	e := hex.EncodeToString

	for i := uint64(0); i < uint64(tests); i++ {
		pld, err := NewMessagePayload(userid.NewUserIDFromUint(i+1, t),
			testStrings[i])

		if (err != nil) != expectedErrors[i] {
			t.Errorf("Didn't expect error result on test %v", i)
		}

		if *userid.NewUserIDFromUint(i+1, t) != *pld.GetSender() {
			t.Errorf("Test of Payload failed on test %v, sID did not "+
				"match;\n  Expected: %v, Received: %v", i, i,
				e(pld.GetSender().Bytes()))
		}

		expct := expectedSlices[i]

		if !bytes.Contains(pld.data, expct) {
			t.Errorf("Test of Payload failed on test %v, "+
				"bytes did not "+
				"match;\n Value Expected: %v, " +
				"Value Received: %v; Lengths: %v, %v", i,
				string(expct), string(pld.data), len(expct), len(pld.data))
		}

		pld.GetPayloadMIC()[MMIC_LEN-1] = uint8(i)
		pld.GetMessagePayloadInitVect()[MMIC_LEN-1] = uint8(i * 5)

		serial := pld.SerializePayload()
		deserial := DeserializeMessagePayload(serial)

		if !bytes.Equal(deserial.GetMessagePayloadInitVect(),
			pld.GetMessagePayloadInitVect()) {
			t.Errorf("Test of Payload failed on "+
				"test %v: Init Vect did not match post serialization;\n"+
				"  Expected: %v, Recieved: %v ", i,
				e(pld.GetMessagePayloadInitVect()),
				e(deserial.GetMessagePayloadInitVect()))
		}

		if !bytes.Equal(deserial.GetSenderID(), pld.GetSenderID()) {
			t.Errorf("Test of Payload failed on test %v, "+
				"Sender ID did not match post serialization;\n"+
				"  Expected: %v, Recieved: %v ", i,
				e(pld.GetSenderID()),
				e(deserial.GetSenderID()))
		}

		if !bytes.Equal(deserial.GetData(), pld.GetData()) {
			t.Errorf("Test of Payload failed on test %v, "+
				"Data did not match post serialization;\n"+
				"  Expected: %v, Recieved: %v ", i,
				e(pld.GetData()),
				e(deserial.GetData()))
		}

		if !bytes.Equal(deserial.GetPayloadMIC(), pld.GetPayloadMIC()) {
			t.Errorf("Test of Payload failed on test %v, "+
				"Payload MIC did not match post serialization;\n"+
				"  Expected: %v, Recieved: %v ", i,
				e(pld.GetPayloadMIC()),
				e(deserial.GetPayloadMIC()))
		}
	}

}
