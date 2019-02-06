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

func TestRecipientPayload(t *testing.T) {
	numRecpts := 5

	rids := []uint64{10, 5, 1000, 3, 0}

	initVects := []*cyclic.Int{cyclic.NewInt(5), cyclic.NewInt(34),
		cyclic.NewInt(89), cyclic.NewInt(77), cyclic.NewInt(10)}

	emptys := []*cyclic.Int{cyclic.NewInt(22), cyclic.NewInt(40),
		cyclic.NewInt(53), cyclic.NewInt(17), cyclic.NewInt(14)}

	mics := []*cyclic.Int{cyclic.NewInt(54), cyclic.NewInt(52),
		cyclic.NewInt(43), cyclic.NewInt(27), cyclic.NewInt(12)}

	recipients := make([]Recipient, numRecpts)

	var err error

	for i := 0; i < numRecpts; i++ {
		recipients[i], err = NewRecipientPayload(id.NewUserIDFromUint(rids[i], t))

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

		if recipients[i].GetRecipientID().Uint64() != rids[i] {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient ID did not match;\n  Expected: %v, "+
				"Recieved: %v ", i, rids[i],
				recipients[i].GetRecipientID().Text(10))
		}

		recipients[i].recipientInitVect.Set(initVects[i])

		if recipients[i].GetRecipientInitVect().Cmp(initVects[i]) != 0 {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Initialization Vectors did not match;\n  Expected: %v, "+
				"Recieved: %v ", i, initVects[i].Text(16),
				recipients[i].GetRecipientInitVect().Text(16))
		}

		recipients[i].recipientEmpty.Set(emptys[i])

		if recipients[i].GetRecipientEmpty().Cmp(emptys[i]) != 0 {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Empty Regions did not match;\n  Expected: %v, "+
				"Recieved: %v ", i, emptys[i].Text(16),
				recipients[i].GetRecipientEmpty().Text(16))
		}

		recipients[i].recipientMIC.Set(mics[i])

		if recipients[i].GetRecipientMIC().Cmp(mics[i]) != 0 {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"MICs did not match;\n  Expected: %v, "+
				"Recieved: %v ", i, mics[i].Text(16),
				recipients[i].GetRecipientMIC().Text(16))
		}

		serial := recipients[i].SerializeRecipient()

		deserial := DeserializeRecipient(serial)

		if deserial.GetRecipientID().Cmp(recipients[i].GetRecipientID()) != 0 {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient ID did not match post serialization;\n  Expected"+
				": %v, Recieved: %v ", i,
				recipients[i].GetRecipientID().Text(10),
				deserial.GetRecipientID().Text(10))
		}

		if deserial.GetRecipientInitVect().Cmp(recipients[i].GetRecipientInitVect()) != 0 {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient InitVect did not match post serialization;\n"+
				"  Expected: %v, Recieved: %v ", i,
				recipients[i].GetRecipientInitVect().Text(16),
				deserial.GetRecipientInitVect().Text(16))
		}

		if deserial.GetRecipientEmpty().Cmp(recipients[i].GetRecipientEmpty()) != 0 {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient Empty did not match post serialization;\n"+
				"  Expected: %v, Recieved: %v ", i,
				recipients[i].GetRecipientEmpty().Text(16),
				deserial.GetRecipientEmpty().Text(16))
		}

		if deserial.GetRecipientMIC().Cmp(recipients[i].GetRecipientMIC()) != 0 {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient MIC did not match post serialization;\n"+
				"  Expected: %v, Recieved: %v ", i,
				recipients[i].GetRecipientMIC().Text(16),
				deserial.GetRecipientMIC().Text(16))
		}

		dcopy := deserial.DeepCopy()

		if deserial.GetRecipientID().Cmp(dcopy.GetRecipientID()) != 0 {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient ID did not match post deep copy;\n  Expected"+
				": %v, Recieved: %v ", i,
				deserial.GetRecipientID().Text(10),
				dcopy.GetRecipientID().Text(10))
		}

		if deserial.GetRecipientInitVect().Cmp(dcopy.GetRecipientInitVect()) != 0 {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient InitVect did not match post deep copy;\n"+
				"  Expected: %v, Recieved: %v ", i,
				deserial.GetRecipientInitVect().Text(16),
				deserial.GetRecipientInitVect().Text(16))
		}

		if deserial.GetRecipientEmpty().Cmp(dcopy.GetRecipientEmpty()) != 0 {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient Empty did not match post deep copy;\n"+
				"  Expected: %v, Recieved: %v ", i,
				deserial.GetRecipientEmpty().Text(16),
				deserial.GetRecipientEmpty().Text(16))
		}

		if deserial.GetRecipientMIC().Cmp(dcopy.GetRecipientMIC()) != 0 {
			t.Errorf("Test of Recipient Payload failed on test %v, "+
				"Recipient MIC did not match post deep copy;\n"+
				"  Expected: %v, Recieved: %v ", i,
				deserial.GetRecipientMIC().Text(16),
				deserial.GetRecipientMIC().Text(16))
		}
	}
}
