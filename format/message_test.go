////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package format

import (
	"math/rand"
	"testing"
)

func TestContentsSize(t *testing.T) {

	m := NewMessage(256)

	prng := rand.New(rand.NewSource(42))

	for i := 0; i < 1000; i++ {
		size := uint16(prng.Uint64() % (1<<14 - 1))

		m.setContentsSize(size)

		gotSize := m.getContentsSize()

		if size != gotSize {
			t.Errorf("Reconstructed size not correct; "+
				"intial: %v, reconstructed: %v", size, gotSize)
		}

	}

}

/*
// Tests that NewAssociatedData() properly sets AssociatedData's serial and all
// other fields.
func TestNewMessage(t *testing.T) {
	// Create new Message
	msg := NewMessage()

	// Test fields
	if !bytes.Equal(msg.master[:], make([]byte, TotalLen)) {
		t.Errorf("NewMessage() did not properly create Message's master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.master, make([]byte, TotalLen))
	}

	if msg.master[TotalLen-1] != 0 {
		t.Errorf("NewMessage() did not set the last byte to zero"+
			"\n\treceived: %v\n\texpected: %v",
			msg.master[TotalLen-1], 0)
	}
}

// Tests that NewMessage() creates all the fields with the correct lengths.
func TestNewMessage_Length(t *testing.T) {
	// Create new Message
	msg := NewMessage()

	// Test lengths
	if len(msg.Contents.serial) != ContentsLen {
		t.Errorf("NewMessage() did not create Message's Contents "+
			"with the correct length"+
			"\n\treceived: %d\n\texpected: %d",
			len(msg.Contents.serial), ContentsLen)
	}

	if len(msg.AssociatedData.serial) != AssociatedDataLen {
		t.Errorf("NewMessage() did not create Message's AssociatedData "+
			"with the correct length"+
			"\n\treceived: %d\n\texpected: %d",
			len(msg.AssociatedData.serial), AssociatedDataLen)
	}

	if len(msg.payloadA) != PayloadLen {
		t.Errorf("NewMessage() did not create Message's payloadA "+
			"with the correct length"+
			"\n\treceived: %d\n\texpected: %d",
			len(msg.payloadA), PayloadLen)
	}

	if len(msg.payloadB) != PayloadLen {
		t.Errorf("NewMessage() did not create Message's payloadB "+
			"with the correct length"+
			"\n\treceived: %d\n\texpected: %d",
			len(msg.payloadB), PayloadLen)
	}

	if len(msg.grpByte) != GrpByteLen {
		t.Errorf("NewMessage() did not create Message's grpByte "+
			"with the correct length"+
			"\n\treceived: %d\n\texpected: %d",
			len(msg.grpByte), GrpByteLen)
	}

	// Check that Content and AssociatedData serials fit in master
	sum := len(msg.Contents.serial) + len(msg.AssociatedData.serial) + len(msg.grpByte)

	if sum != len(msg.master) {
		t.Errorf("The sum of the lengths of Content, AssociatedData, and "+
			"grpByte does not equal the length of master"+
			"\n\treceived: %d\n\texpected: %d",
			sum, len(msg.master))
	}

	// Check that payloadA and payloadB serials fit in master
	sum = len(msg.payloadA) + len(msg.payloadB)

	if sum != len(msg.master) {
		t.Errorf("The sum of the lengths of the payloads does not equal "+
			"the length of master"+
			"\n\treceived: %d\n\texpected: %d",
			sum, len(msg.master))
	}
}

// Tests that no two sub fields overlap (except when they are supposed to),
// that the start of payloadA, Contents, and master is the same, and that the
// end of payloadB, grpByte, and master are the same.
func TestNewMessage_Overlap(t *testing.T) {
	// Create new Message
	msg := NewMessage()

	// Check the fields
	if reflect.ValueOf(msg.master[:0]).Pointer() !=
		reflect.ValueOf(msg.payloadA[:0]).Pointer() {
		t.Errorf("The start of master is not the same pointer as the "+
			"start of payloadA"+
			"\n\tstart of master:   %d\n\tstart of payloadA: %d",
			reflect.ValueOf(msg.master[:0]).Pointer(),
			reflect.ValueOf(msg.payloadA[:0]).Pointer())
	}

	if reflect.ValueOf(msg.master[:0]).Pointer() !=
		reflect.ValueOf(msg.Contents.serial[:0]).Pointer() {
		t.Errorf("The start of master is not the same pointer as the "+
			"start of Contents"+
			"\n\tstart of master:   %d\n\tstart of payloadA: %d",
			reflect.ValueOf(msg.master[:0]).Pointer(),
			reflect.ValueOf(msg.Contents.serial[:0]).Pointer())
	}

	if reflect.ValueOf(msg.payloadA[:PayloadLen-1]).Pointer() >=
		reflect.ValueOf(msg.payloadB[:0]).Pointer() {
		t.Errorf("The end of payloadA overlaps with the start of payloadB"+
			"\n\tend of payloadA:   %d\n\tstart of payloadB: %d",
			reflect.ValueOf(msg.payloadA[:PayloadLen-1]).Pointer(),
			reflect.ValueOf(msg.payloadB[:0]).Pointer())
	}

	if reflect.ValueOf(msg.Contents.serial[:ContentsLen-1]).Pointer() >=
		reflect.ValueOf(msg.AssociatedData.serial[:0]).Pointer() {
		t.Errorf("The end of Contents overlaps with the start of AssociatedData"+
			"\n\tend of Contents:         %d\n\tstart of AssociatedData: %d",
			reflect.ValueOf(msg.Contents.serial[:ContentsLen-1]).Pointer(),
			reflect.ValueOf(msg.AssociatedData.serial[:0]).Pointer())
	}

	if reflect.ValueOf(msg.AssociatedData.serial[:AssociatedDataLen-1]).Pointer() >=
		reflect.ValueOf(msg.grpByte[:0]).Pointer() {
		t.Errorf("The end of AssociatedData overlaps with the start of grpByte"+
			"\n\tend of AssociatedData:   %d\n\tstart of grpByte:        %d",
			reflect.ValueOf(msg.AssociatedData.serial[:AssociatedDataLen-1]).Pointer(),
			reflect.ValueOf(msg.grpByte[:0]).Pointer())
	}

	if reflect.ValueOf(msg.master[TotalLen-1:]).Pointer() !=
		reflect.ValueOf(msg.payloadB[PayloadLen-1:]).Pointer() {
		t.Errorf("The end of master is not the same pointer as the "+
			"end of payloadB"+
			"\n\tend of master:   %d\n\tend of payloadB: %d",
			reflect.ValueOf(msg.master[TotalLen-1:]).Pointer(),
			reflect.ValueOf(msg.payloadB[PayloadLen-1:]).Pointer())
	}

	if reflect.ValueOf(msg.master[TotalLen-1:]).Pointer() !=
		reflect.ValueOf(msg.grpByte[GrpByteLen-1:]).Pointer() {
		t.Errorf("The end of master is not the same pointer as the "+
			"end of payloadB"+
			"\n\tend of master:  %d\n\tend of grpByte: %d",
			reflect.ValueOf(msg.master[TotalLen-1:]).Pointer(),
			reflect.ValueOf(msg.grpByte[GrpByteLen-1:]).Pointer())
	}
}

// Tests that when values are set for each field that they are reflected in
// master.
func TestMessage_Values(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randPayloadA := make([]byte, PayloadLen)
	rand.Read(randPayloadA)
	randPayloadB := make([]byte, PayloadLen)
	rand.Read(randPayloadB)
	randContents := make([]byte, ContentsLen)
	rand.Read(randContents)
	randAssociatedData := make([]byte, AssociatedDataLen)
	rand.Read(randAssociatedData)

	// Create new Message and set payload fields
	msg := NewMessage()
	msg.SetPayloadA(randPayloadA)
	msg.SetPayloadB(randPayloadB)

	// Check if the values set to each field are reflected in master
	if !bytes.Equal(msg.GetPayloadA(), msg.master[payloadAStart:payloadAEnd]) {
		t.Errorf("payloadA is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetPayloadA(), msg.master[payloadAStart:payloadAEnd])
	}

	if !bytes.Equal(msg.GetPayloadB(), msg.master[payloadBStart:payloadBEnd]) {
		t.Errorf("payloadB is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetPayloadB(), msg.master[payloadBStart:payloadBEnd])
	}

	if !bytes.Equal(msg.Contents.Get(), msg.master[contentsStart:contentsEnd]) {
		t.Errorf("Contents is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.Contents.Get(), msg.master[contentsStart:contentsEnd])
	}

	if msg.Contents.GetPosition() != invalidPosition {
		t.Errorf("Contents position is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.Contents.GetPosition(), invalidPosition)
	}

	if !bytes.Equal(msg.AssociatedData.Get(), msg.master[associatedDataStart:associatedDataEnd]) {
		t.Errorf("AssociatedData is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.AssociatedData.Get(), msg.master[associatedDataStart:associatedDataEnd])
	}

	if !bytes.Equal(msg.AssociatedData.GetRecipientID(), msg.master[associatedDataStart:associatedDataStart+RecipientIDLen]) {
		t.Errorf("AssociatedData recipientID is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.AssociatedData.GetRecipientID(), msg.master[associatedDataStart:associatedDataStart+RecipientIDLen])
	}

	fp := msg.AssociatedData.GetKeyFP()

	if !bytes.Equal(fp[:], msg.master[associatedDataStart+RecipientIDLen:associatedDataStart+RecipientIDLen+KeyFPLen]) {
		t.Errorf("AssociatedData keyFP is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			fp[:], msg.master[associatedDataStart+RecipientIDLen:associatedDataStart+RecipientIDLen+KeyFPLen])
	}

	if !bytes.Equal(msg.AssociatedData.GetTimestamp(), msg.master[associatedDataStart+RecipientIDLen+KeyFPLen:associatedDataStart+RecipientIDLen+KeyFPLen+TimestampLen]) {
		t.Errorf("AssociatedData timestamp is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.AssociatedData.GetTimestamp(), msg.master[associatedDataStart+RecipientIDLen+KeyFPLen:associatedDataStart+RecipientIDLen+KeyFPLen+TimestampLen])
	}

	if !bytes.Equal(msg.AssociatedData.GetMAC(), msg.master[associatedDataStart+RecipientIDLen+KeyFPLen+TimestampLen:associatedDataStart+RecipientIDLen+KeyFPLen+TimestampLen+MacLen]) {
		t.Errorf("AssociatedData MAC is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.AssociatedData.GetMAC(), msg.master[associatedDataStart+RecipientIDLen+KeyFPLen+TimestampLen:associatedDataStart+RecipientIDLen+KeyFPLen+TimestampLen+MacLen])
	}

	// Create new Message and set Contents and AssociatedData
	msg.Contents.Set(randContents)
	msg.AssociatedData.Set(randAssociatedData)

	// Check if the values set to each field are reflected in master
	if !bytes.Equal(msg.GetPayloadA(), msg.master[payloadAStart:payloadAEnd]) {
		t.Errorf("payloadA is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetPayloadA(), msg.master[payloadAStart:payloadAEnd])
	}

	if !bytes.Equal(msg.GetPayloadB(), msg.master[payloadBStart:payloadBEnd]) {
		t.Errorf("payloadB is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetPayloadB(), msg.master[payloadBStart:payloadBEnd])
	}

	if !bytes.Equal(msg.Contents.Get(), msg.master[contentsStart:contentsEnd]) {
		t.Errorf("Contents is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.Contents.Get(), msg.master[contentsStart:contentsEnd])
	}

	if !bytes.Equal(msg.AssociatedData.Get(), msg.master[associatedDataStart:associatedDataEnd]) {
		t.Errorf("AssociatedData is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.AssociatedData.Get(), msg.master[associatedDataStart:associatedDataEnd])
	}

	if !bytes.Equal(msg.AssociatedData.GetRecipientID(), msg.master[associatedDataStart:associatedDataStart+RecipientIDLen]) {
		t.Errorf("AssociatedData recipientID is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.AssociatedData.GetRecipientID(), msg.master[associatedDataStart:associatedDataStart+RecipientIDLen])
	}

	fp = msg.AssociatedData.GetKeyFP()

	if !bytes.Equal(fp[:], msg.master[associatedDataStart+RecipientIDLen:associatedDataStart+RecipientIDLen+KeyFPLen]) {
		t.Errorf("AssociatedData keyFP is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			fp[:], msg.master[associatedDataStart+RecipientIDLen:associatedDataStart+RecipientIDLen+KeyFPLen])
	}

	if !bytes.Equal(msg.AssociatedData.GetTimestamp(), msg.master[associatedDataStart+RecipientIDLen+KeyFPLen:associatedDataStart+RecipientIDLen+KeyFPLen+TimestampLen]) {
		t.Errorf("AssociatedData timestamp is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.AssociatedData.GetTimestamp(), msg.master[associatedDataStart+RecipientIDLen+KeyFPLen:associatedDataStart+RecipientIDLen+KeyFPLen+TimestampLen])
	}

	if !bytes.Equal(msg.AssociatedData.GetMAC(), msg.master[associatedDataStart+RecipientIDLen+KeyFPLen+TimestampLen:associatedDataStart+RecipientIDLen+KeyFPLen+TimestampLen+MacLen]) {
		t.Errorf("AssociatedData MAC is not properly mapped to master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.AssociatedData.GetMAC(), msg.master[associatedDataStart+RecipientIDLen+KeyFPLen+TimestampLen:associatedDataStart+RecipientIDLen+KeyFPLen+TimestampLen+MacLen])
	}
}

// Tests that Message is always constructed the same way.
func TestMessage_Consistency(t *testing.T) {
	// Define the master to check against (64-bit encoded)
	master, _ := base64.StdEncoding.DecodeString("U4x/lrFkvxuXu59LtHLon1s" +
		"UhPJSCcnZND6SugndnVLf15tNdkKbYXoMn58NO6VbDMDWFEyIhTWEGsvgcJsHWAg/Yd" +
		"N1vAK0HfT5GSnhj9qeb4LlTnSOgeeeS71v40zcuoQ+6NY+jE/+HOvqVG2PrBPdGqwEz" +
		"i6ih3xVec+ix44bC6+uiBuCp1EQikLtPJA8qkNGWnhiBhaXiu0M48bE8657w+BJW1cS" +
		"/v2+DBAoh+EA2s0tiF9pLLYH2gChHBxwceeWotwtwlpbdLLhKXBeJz8FySMmgo4rBW4" +
		"4F2WOEGFJiUf980RBDtTBFgI/qONXa2/tJ/+JdLrAyv2a0FaSsTYZ5ziWTf3Hno1TQ3" +
		"NmHP1m10/sHhuJSRq3I25LdSFikM8r60LDyicyhWDxqsBnzqbov0bUqytGgEAsX7KCD" +
		"ohdMmDx3peCg9Sgmjb5bCCUF0bj7U2mRqmui0+ntPw6ILr6GnXtMnqGuLDDmvHP0rO1" +
		"EhnqeVM6v0SNLEedMmB1M5BZFMjMHPCdo54Okp0CSry8sWk5e7c05+8KbgHxhU3rX+Q" +
		"k/vesIQiR9ZdeKSqiuKoEfGHNszNz6+csJ6CYwCGX2ua3MsNR32aPh04snxzgnKhgF+" +
		"fiF0gwP/QcGyPhHEjtF1OdaF928qeYvGTeDl2yhksq08Js5jgjQnZaE9Y=")

	// Generate random byte slice
	rand.Seed(42)
	randPayloadA := make([]byte, PayloadLen)
	rand.Read(randPayloadA)
	randPayloadB := make([]byte, PayloadLen)
	rand.Read(randPayloadB)

	// Create new Message and set payload fields
	msg := NewMessage()
	msg.SetPayloadA(randPayloadA)
	msg.SetPayloadB(randPayloadB)

	if !bytes.Equal(msg.GetMaster(), master) {
		t.Errorf("Message's master does not match the hardcoded master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetMaster(), master)
	}
}

// Tests that GetMaster() returns the correct bytes set to Message's master.
func TestMessage_GetMaster(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, TotalLen)
	rand.Read(randSlice)

	// Create new Message
	msg := NewMessage()
	copy(msg.master[:], randSlice)

	if !bytes.Equal(msg.GetMaster(), randSlice) {
		t.Errorf("GetMaster() did not return the correct data from "+
			"Message's master"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetMaster(), randSlice)
	}
}

// Tests that GetPayloadA() returns the correct bytes set to Message's payloadA.
func TestMessage_GetPayloadA(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, TotalLen)
	rand.Read(randSlice)

	// Create new Message
	msg := NewMessage()
	copy(msg.master[:], randSlice)

	if !bytes.Equal(msg.GetPayloadA(), randSlice[payloadAStart:payloadAEnd]) {
		t.Errorf("GetPayloadA() did not return the correct data from "+
			"Message's payloadA"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetPayloadA(), randSlice[payloadAStart:payloadAEnd])
	}
}

// Tests that SetPayloadA() sets the correct bytes to Message's payloadA and
// copies the correct number of bytes.
func TestMessage_SetPayloadA(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, PayloadLen)
	rand.Read(randSlice)

	// Create new Message and set payloadA
	msg := NewMessage()
	msg.SetPayloadA(randSlice)

	if !bytes.Equal(msg.GetPayloadA(), randSlice) {
		t.Errorf("SetPayloadA() did not properly set Message's payloadA"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetPayloadA(), randSlice)
	}
}

// Tests that SetPayloadA() panics when the new payload is not the same length
// as payloadA.
func TestMessage_SetPayloadA_Panic(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, PayloadLen+5)
	rand.Read(randSlice)

	// Defer to an error when SetPayloadA() does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetPayloadA() did not panic when expected")
		}
	}()

	// Create new Message and set payloadA
	msg := NewMessage()
	msg.SetPayloadA(randSlice)
}

// Tests that GetPayloadB() returns the correct bytes set to Message's payloadB.
func TestMessage_GetPayloadB(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, TotalLen)
	rand.Read(randSlice)

	// Create new Message
	msg := NewMessage()
	copy(msg.master[:], randSlice)

	if !bytes.Equal(msg.GetPayloadB(), randSlice[payloadBStart:payloadBEnd]) {
		t.Errorf("GetPayloadB() did not return the correct data from "+
			"Message's payloadB"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetPayloadB(), randSlice[payloadBStart:payloadBEnd])
	}
}

// Tests that SetPayloadB() sets the correct bytes to Message's payloadB and
// copies the correct number of bytes.
func TestMessage_SetPayloadB(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, PayloadLen)
	rand.Read(randSlice)

	// Create new Message and set payloadB
	msg := NewMessage()
	msg.SetPayloadB(randSlice)

	if !bytes.Equal(msg.GetPayloadB(), randSlice) {
		t.Errorf("SetPayloadB() did not properly set Message's payloadB"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetPayloadB(), randSlice)
	}
}

// Tests that SetPayloadB() panics when the new payload is not the same length
// as payloadB.
func TestMessage_SetPayloadB_Panic(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, PayloadLen+5)
	rand.Read(randSlice)

	// Defer to an error when SetPayloadB() does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetPayloadB() did not panic when expected")
		}
	}()

	// Create new Message and set payloadB
	msg := NewMessage()
	msg.SetPayloadB(randSlice)
}

// Tests that GetPayloadBForEncryption() returns the correct bytes set to
// Message's payloadB. Also checks that the first and last byte are swapped and
// that the first byte is zero.
func TestMessage_GetPayloadBForEncryption(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, TotalLen)
	rand.Read(randSlice)
	randSlice[TotalLen-1] = 0

	// Create new Message
	msg := NewMessage()
	copy(msg.master[:], randSlice)

	if !bytes.Equal(msg.GetPayloadBForEncryption()[1:PayloadLen-1], randSlice[payloadBStart+1:payloadBEnd-1]) {
		t.Errorf("GetPayloadBForEncryption() did not return the correct data from "+
			"Message's payloadB"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetPayloadBForEncryption()[1:PayloadLen-1], randSlice[payloadBStart+1:payloadBEnd-1])
	}

	if msg.GetPayloadBForEncryption()[0] != 0 {
		t.Errorf("GetPayloadBForEncryption() did not set the first byte to zero"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetPayloadBForEncryption()[0], 0)
	}

	if msg.GetPayloadBForEncryption()[PayloadLen-1] != randSlice[payloadBStart] {
		t.Errorf("GetPayloadBForEncryption() did not correctly swap the "+
			"first and last byte"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetPayloadBForEncryption()[PayloadLen-1], randSlice[payloadBStart])
	}
}

// Tests that SetDecryptedPayloadB() sets the correct bytes to Message's
// payloadB and copies the correct number of bytes. Also checks that the first
// and last bytes are swapped and that the first byte is zero
func TestMessage_SetDecryptedPayloadB(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, PayloadLen)
	rand.Read(randSlice)

	// Create new Message and set payloadB
	msg := NewMessage()
	msg.SetDecryptedPayloadB(randSlice)

	if !bytes.Equal(msg.GetPayloadB()[1:PayloadLen-1], randSlice[1:PayloadLen-1]) {
		t.Errorf("SetDecryptedPayloadB() did not return the correct data from "+
			"Message's payloadB"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetPayloadB()[1:PayloadLen-1], randSlice[1:PayloadLen-1])
	}

	if msg.GetPayloadB()[PayloadLen-1] != 0 {
		t.Errorf("SetDecryptedPayloadB() did not set the last byte to zero"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetPayloadB()[PayloadLen-1], 0)
	}

	if msg.GetPayloadB()[0] != randSlice[PayloadLen-1] {
		t.Errorf("SetDecryptedPayloadB() did not correctly swap the "+
			"first and last byte"+
			"\n\treceived: %v\n\texpected: %v",
			msg.GetPayloadB()[0], randSlice[PayloadLen-1])
	}
}

// Tests that SetDecryptedPayloadB() panics when the new payload is not the same
// length as payloadB.
func TestMessage_SetDecryptedPayloadB_Panic(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, PayloadLen+5)
	rand.Read(randSlice)

	// Defer to an error when SetPayloadB() does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetDecryptedPayloadB() did not panic when expected")
		}
	}()

	// Create new Message and set payloadB
	msg := NewMessage()
	msg.SetDecryptedPayloadB(randSlice)
}*/
