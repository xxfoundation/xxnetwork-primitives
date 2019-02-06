////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"gitlab.com/elixxir/crypto/cyclic"
	"testing"
	"gitlab.com/elixxir/primitives/userid"
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

	for i := uint64(0); i < uint64(tests); i++ {
		pldSlc := NewPayload(id.NewUserIDFromUint(i+1, t), testStrings[i])

		for indx, pld := range pldSlc {
			if *id.NewUserIDFromUint(i+1, t) != *pld.GetSender() {
				t.Errorf("Test of Payload failed on test %v:%v, sID did not "+
					"match;\n  Expected: %v, Received: %v", i, indx, i,
					pld.GetSender())
			}

			expct := cyclic.NewIntFromBytes(expectedSlices[i][indx])

			if pld.data.Cmp(expct) != 0 {
				t.Errorf("Test of Payload failed on test %v:%v, "+
					"bytes did not "+
					"match;\n Value Expected: %v, Value Received: %v", i, indx,
					string(expct.Bytes()), string(pld.data.Bytes()))
			}

			pld.GetPayloadMIC().SetUint64(uint64(i))
			pld.GetPayloadInitVect().SetUint64(uint64(i * 5))

			serial := pld.SerializePayload()
			deserial := DeserializePayload(serial)

			if deserial.GetPayloadInitVect().Cmp(pld.GetPayloadInitVect()) != 0 {
				t.Errorf("Test of Payload failed on "+
					"test %v: %v, Init Vect did not match post serialization;\n"+
					"  Expected: %v, Recieved: %v ", i, indx,
					pld.GetPayloadInitVect().Text(16),
					deserial.GetPayloadInitVect().Text(16))
			}

			if deserial.GetSenderID().Cmp(pld.GetSenderID()) != 0 {
				t.Errorf("Test of Payload failed on test %v:%v, "+
					"Sender ID did not match post serialization;\n"+
					"  Expected: %v, Recieved: %v ", i, indx,
					pld.GetSenderID().Text(10),
					deserial.GetSenderID().Text(10))
			}

			if deserial.GetData().Cmp(pld.GetData()) != 0 {
				t.Errorf("Test of Payload failed on test %v:%v, "+
					"Data did not match post serialization;\n"+
					"  Expected: %v, Recieved: %v ", i, indx,
					pld.GetData().Text(16),
					deserial.GetData().Text(16))
			}

			if deserial.GetPayloadMIC().Cmp(pld.GetPayloadMIC()) != 0 {
				t.Errorf("Test of Payload failed on test %v:%v, "+
					"Payload MIC did not match post serialization;\n"+
					"  Expected: %v, Recieved: %v ", i, indx,
					pld.GetPayloadMIC().Text(16),
					deserial.GetPayloadMIC().Text(16))
			}
		}

	}

}
