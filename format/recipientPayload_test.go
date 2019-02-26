////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"bytes"
	"encoding/hex"
	"gitlab.com/elixxir/primitives/id"
	"testing"
)

func TestRecipientPayload(t *testing.T) {
	numRecpts := 5

	rids := []uint64{10, 5, 1000, 3, 0}

	// Set the last byte of each item
	initVectBytes := []byte{5, 34, 89, 77, 10}
	initVects := [][]byte{}
	for i := range initVectBytes {
		initVects = append(initVects, make([]byte, RIV_LEN))
		initVects[i][len(initVects[i])-1] = initVectBytes[i]
	}

	emptyBytes := []byte{22, 40, 53, 17, 14}
	emptys := [][]byte{}
	for i := range emptyBytes {
		emptys = append(emptys, make([]byte, REMPTY_LEN))
		emptys[i][len(emptys[i])-1] = emptyBytes[i]
	}

	micBytes := []byte{54, 52, 43, 27, 12}
	mics := [][]byte{}
	for i := range micBytes {
		mics = append(mics, make([]byte, RMIC_LEN))
		mics[i][len(mics[i])-1] = micBytes[i]
	}

	recipients := make([]*AssociatedData, numRecpts)

	var err error

	for i := 0; i < numRecpts; i++ {
		recipients[i], err = NewAssociatedData(id.NewUserFromUint(rids[i],
			t))

		if err != nil && rids[i] != 0 {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"returned the following unexpected error: %s ", i, err.Error())
			continue
		} else if err == nil && rids[i] == 0 {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"didnt fails as expected ", i)
			continue
		} else if rids[i] == 0 {
			continue
		}

		e := hex.EncodeToString
		if !bytes.Equal(recipients[i].GetRecipientID(),
			id.NewUserFromUint(rids[i], t).Bytes()) {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient ID did not match;\n  Expected: %v, "+
				"Recieved: %v ", i, rids[i],
				e(recipients[i].GetRecipientID()))
		}

		copy(recipients[i].recipientInitVect, initVects[i])

		if !bytes.Equal(recipients[i].GetRecipientInitVect(), initVects[i]) {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Initialization Vectors did not match;\n  Expected: %v, "+
				"Recieved: %v ", i, e(initVects[i]),
				e(recipients[i].GetRecipientInitVect()))
		}

		copy(recipients[i].recipientMIC, mics[i])

		if !bytes.Equal(recipients[i].GetRecipientMIC(), mics[i]) {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"MICs did not match;\n  Expected: %v, "+
				"Recieved: %v ", i, e(mics[i]),
				e(recipients[i].GetRecipientMIC()))
		}

		// Make sure that things are still accessible after serialization/deserialization
		serial := recipients[i].SerializeRecipient()

		deserial := DeserializeAssociatedData(serial)

		if !bytes.Equal(deserial.GetRecipientID(), recipients[i].GetRecipientID()) {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient ID did not match post serialization;\n  Expected"+
				": %v, Recieved: %v ", i,
				e(recipients[i].GetRecipientID()),
				e(deserial.GetRecipientID()))
		}

		if !bytes.Equal(deserial.GetRecipientInitVect(), recipients[i].GetRecipientInitVect()) {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient InitVect did not match post serialization;\n"+
				"  Expected: %v, Recieved: %v ", i,
				e(recipients[i].GetRecipientInitVect()),
				e(deserial.GetRecipientInitVect()))
		}

		if !bytes.Equal(deserial.GetRecipientMIC(), recipients[i].GetRecipientMIC()) {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient MIC did not match post serialization;\n"+
				"  Expected: %v, Recieved: %v ", i,
				e(recipients[i].GetRecipientMIC()),
				e(deserial.GetRecipientMIC()))
		}

		dcopy := deserial.DeepCopy()

		if !bytes.Equal(deserial.GetRecipientID(), dcopy.GetRecipientID()) {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient ID did not match post deep copy;\n  Expected"+
				": %v, Recieved: %v ", i,
				e(deserial.GetRecipientID()),
				e(dcopy.GetRecipientID()))
		}

		if !bytes.Equal(deserial.GetRecipientInitVect(), dcopy.GetRecipientInitVect()) {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient InitVect did not match post deep copy;\n"+
				"  Expected: %v, Recieved: %v ", i,
				e(deserial.GetRecipientInitVect()),
				e(deserial.GetRecipientInitVect()))
		}

		if !bytes.Equal(deserial.GetRecipientMIC(), dcopy.GetRecipientMIC()) {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient MIC did not match post deep copy;\n"+
				"  Expected: %v, Recieved: %v ", i,
				e(deserial.GetRecipientMIC()),
				e(deserial.GetRecipientMIC()))
		}
	}
}
