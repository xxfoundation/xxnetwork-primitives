////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"testing"
	"gitlab.com/elixxir/primitives/userid"
	"bytes"
	"encoding/hex"
)

func TestMessagePayload(t *testing.T) {
	tests := 3

	testStrings := make([][]byte, tests)

	testStrings[0] = testText[0 : DATA_LEN/2]
	testStrings[1] = testText[0:DATA_LEN]

	testStrings[2] = testText[0 : 2*DATA_LEN]

	expectedSlices := make([][][]byte, tests)

	expectedSlices[0] = make([][]byte, 1)

	expectedSlices[0][0] = []byte(testStrings[0])

	expectedSlices[1] = make([][]byte, 2)

	expectedSlices[1][0] = ([]byte(testStrings[1]))[0:DATA_LEN]

	expectedSlices[2] = make([][]byte, 3)

	expectedSlices[2][0] = ([]byte(testStrings[2]))[0:DATA_LEN]
	expectedSlices[2][1] = ([]byte(testStrings[2]))[DATA_LEN : 2*DATA_LEN]
	expectedSlices[2][2] = ([]byte(testStrings[2]))[2*DATA_LEN:]

	e := hex.EncodeToString

	for i := uint64(0); i < uint64(tests); i++ {
		pldSlc := NewPayload(userid.NewUserIDFromUint(i+1, t), testStrings[i])

		for indx, pld := range pldSlc {
			if *userid.NewUserIDFromUint(i+1, t) != *pld.GetSender() {
				t.Errorf("Test of Payload failed on test %v:%v, sID did not "+
					"match;\n  Expected: %v, Received: %v", i, indx, i,
					e(pld.GetSender().Bytes()))
			}

			expct := expectedSlices[i][indx]

			if !bytes.Equal(pld.data, expct) {
				t.Errorf("Test of Payload failed on test %v:%v, "+
					"bytes did not "+
					"match;\n Value Expected: %v, Value Received: %v", i, indx,
					string(expct), string(pld.data))
			}

			pld.GetPayloadMIC()[PMIC_LEN-1] = uint8(i)
			pld.GetPayloadInitVect()[PMIC_LEN-1] = uint8(i * 5)

			serial := pld.serializePayload()
			deserial := deserializePayload(serial)

			if !bytes.Equal(deserial.GetPayloadInitVect(),
				pld.GetPayloadInitVect()) {
				t.Errorf("Test of Payload failed on "+
					"test %v: %v, Init Vect did not match post serialization;\n"+
					"  Expected: %v, Recieved: %v ", i, indx,
					e(pld.GetPayloadInitVect()),
					e(deserial.GetPayloadInitVect()))
			}

			if !bytes.Equal(deserial.GetSenderID(), pld.GetSenderID()) {
				t.Errorf("Test of Payload failed on test %v:%v, "+
					"Sender ID did not match post serialization;\n"+
					"  Expected: %v, Recieved: %v ", i, indx,
					e(pld.GetSenderID()),
					e(deserial.GetSenderID()))
			}

			// Use Contains instead of Equal here because the deserialized
			// data string will include trailing zeroes. Maybe GetData should
			// trim trailing zeroes?
			if !bytes.Contains(deserial.GetData(), pld.GetData()) {
				t.Errorf("Test of Payload failed on test %v:%v, "+
					"Data did not match post serialization;\n"+
					"  Expected: %v, Recieved: %v ", i, indx,
					e(pld.GetData()),
					e(deserial.GetData()))
			}

			if !bytes.Equal(deserial.GetPayloadMIC(), pld.GetPayloadMIC()) {
				t.Errorf("Test of Payload failed on test %v:%v, "+
					"Payload MIC did not match post serialization;\n"+
					"  Expected: %v, Recieved: %v ", i, indx,
					e(pld.GetPayloadMIC()),
					e(deserial.GetPayloadMIC()))
			}
		}

	}

}
