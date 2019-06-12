////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"bytes"
	"encoding/base64"
	"gitlab.com/elixxir/primitives/id"
	"math/rand"
	"reflect"
	"testing"
)

// Tests that NewAssociatedData() properly sets AssociatedData's serial and all
// other fields.
func TestNewAssociatedData(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, associatedDataLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(randSlice)

	// Test fields
	if !bytes.Equal(ad.serial, randSlice) {
		t.Errorf("NewAssociatedData() did not properly set "+
			"AssociatedData's serial"+
			"\n\treceived: %v\n\texpected: %v",
			ad.serial, randSlice)
	}

	if !bytes.Equal(ad.recipientID, randSlice[recipientIDStart:recipientIDEnd]) {
		t.Errorf("NewAssociatedData() did not properly set "+
			"AssociatedData's recipientID"+
			"\n\treceived: %v\n\texpected: %v",
			ad.recipientID, randSlice[recipientIDStart:recipientIDEnd])
	}

	if !bytes.Equal(ad.keyFP, randSlice[keyFPStart:keyFPEnd]) {
		t.Errorf("NewAssociatedData() did not properly set "+
			"AssociatedData's keyFP"+
			"\n\treceived: %v\n\texpected: %v",
			ad.keyFP, randSlice[keyFPStart:keyFPEnd])
	}

	if !bytes.Equal(ad.timestamp, randSlice[timestampStart:timestampEnd]) {
		t.Errorf("NewAssociatedData() did not properly set "+
			"AssociatedData's timestamp"+
			"\n\treceived: %v\n\texpected: %v",
			ad.timestamp, randSlice[timestampStart:timestampEnd])
	}

	if !bytes.Equal(ad.mac, randSlice[macStart:macEnd]) {
		t.Errorf("NewAssociatedData() did not properly set "+
			"AssociatedData's mac"+
			"\n\treceived: %v\n\texpected: %v",
			ad.mac, randSlice[macStart:macEnd])
	}
}

// Tests that NewAssociatedData() panics when the new serial is not the same
// length as serial.
func TestNewAssociatedData_Panic(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, associatedDataLen-5)
	rand.Read(randSlice)

	// Defer to an error when NewAssociatedData() does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("NewAssociatedData() did not panic when expected")
		}
	}()

	// Create new AssociatedData
	NewAssociatedData(randSlice)
}

// Tests that NewAssociatedData() creates all the fields with the correct
// lengths.
func TestNewAssociatedData_Length(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, associatedDataLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(randSlice)

	// Test lengths
	if len(ad.serial) != associatedDataLen {
		t.Errorf("NewAssociatedData() did not create "+
			"AssociatedData's serial with the correct length"+
			"\n\treceived: %d\n\texpected: %d",
			len(ad.serial), associatedDataLen)
	}

	if len(ad.recipientID) != recipientIDLen {
		t.Errorf("NewAssociatedData() did not create "+
			"AssociatedData's recipientID with the correct length"+
			"\n\treceived: %d\n\texpected: %d",
			len(ad.recipientID), recipientIDLen)
	}

	if len(ad.keyFP) != keyFPLen {
		t.Errorf("NewAssociatedData() did not create "+
			"AssociatedData's keyFP with the correct length"+
			"\n\treceived: %d\n\texpected: %d",
			len(ad.keyFP), keyFPLen)
	}

	if len(ad.timestamp) != timestampLen {
		t.Errorf("NewAssociatedData() did not create "+
			"AssociatedData's timestamp with the correct length"+
			"\n\treceived: %d\n\texpected: %d",
			len(ad.timestamp), timestampLen)
	}

	if len(ad.mac) != macLen {
		t.Errorf("NewAssociatedData() did not create "+
			"AssociatedData's mac with the correct length"+
			"\n\treceived: %d\n\texpected: %d",
			len(ad.mac), macLen)
	}

	// Check that the summation of the length of all fields equals the length
	// of serial
	sum := len(ad.recipientID) + len(ad.keyFP) + len(ad.timestamp) + len(ad.mac)

	if sum != len(ad.serial) {
		t.Errorf("The sum of the lengths of all fields does not equal "+
			"the length of the serial"+
			"\n\treceived: %d\n\texpected: %d",
			sum, len(ad.serial))
	}
}

// Tests that no two sub fields overlap and that the start of serial is the same
// as the start of recipientID and that the end of serial is the same as the end
// of grpByte.
func TestNewAssociatedData_Overlap(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, associatedDataLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(randSlice)

	// Check the fields
	if reflect.ValueOf(ad.serial[:0]).Pointer() !=
		reflect.ValueOf(ad.recipientID[:0]).Pointer() {
		t.Errorf("The start of serial is not the same pointer as the "+
			"start of recipientID"+
			"\n\tstart of serial:      %d\n\tstart of recipientID: %d",
			reflect.ValueOf(ad.serial[:0]).Pointer(),
			reflect.ValueOf(ad.recipientID[:0]).Pointer())
	}

	if reflect.ValueOf(ad.recipientID[:recipientIDLen-1]).Pointer() >=
		reflect.ValueOf(ad.keyFP[:0]).Pointer() {
		t.Errorf("The end of recipientID overlaps with the start of keyFP"+
			"\n\tend of recipientID: %d\n\tstart of keyFP:     %d",
			reflect.ValueOf(ad.recipientID[recipientIDLen-1:]).Pointer(),
			reflect.ValueOf(ad.keyFP[:0]).Pointer())
	}

	if reflect.ValueOf(ad.keyFP[:keyFPLen-1]).Pointer() >=
		reflect.ValueOf(ad.timestamp[:0]).Pointer() {
		t.Errorf("The end of keyFP overlaps with the start of timestamp"+
			"\n\tend of keyFP:       %d\n\tstart of timestamp: %d",
			reflect.ValueOf(ad.keyFP[keyFPLen-1:]).Pointer(),
			reflect.ValueOf(ad.timestamp[:0]).Pointer())
	}

	if reflect.ValueOf(ad.timestamp[:timestampLen-1]).Pointer() >=
		reflect.ValueOf(ad.mac[:0]).Pointer() {
		t.Errorf("The end of timestamp overlaps with the start of mac"+
			"\n\tend of timestamp: %d\n\tstart of mac:     %d",
			reflect.ValueOf(ad.timestamp[timestampLen-1:]).Pointer(),
			reflect.ValueOf(ad.mac[:0]).Pointer())
	}

	if reflect.ValueOf(ad.serial[associatedDataLen-1:]).Pointer() !=
		reflect.ValueOf(ad.mac[macLen-1:]).Pointer() {
		t.Errorf("The end of serial is not the same pointer as the "+
			"end of mac"+
			"\n\tend of serial: %d\n\tend of mac:    %d",
			reflect.ValueOf(ad.serial[associatedDataLen-1:]).Pointer(),
			reflect.ValueOf(ad.mac[macLen-1:]).Pointer())
	}
}

// Tests that when values are set for each field that they are reflected in
// serial.
func TestAssociatedData_Values(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randAD := make([]byte, associatedDataLen)
	rand.Read(randAD)
	randRecipientID := make([]byte, recipientIDLen)
	rand.Read(randRecipientID)
	randKeyFP := make([]byte, keyFPLen)
	rand.Read(randKeyFP)
	fp := NewFingerprint(randKeyFP)
	randTimestamp := make([]byte, timestampLen)
	rand.Read(randTimestamp)
	randMAC := make([]byte, macLen)
	rand.Read(randMAC)

	// Create new AssociatedData and set fields
	ad := NewAssociatedData(randAD)
	ad.SetRecipientID(randRecipientID)
	ad.SetKeyFP(*fp)
	ad.SetTimestamp(randTimestamp)
	ad.SetMAC(randMAC)

	// Check if the values set to each field are reflected in serial
	if !bytes.Equal(ad.GetRecipientID(), ad.serial[recipientIDStart:recipientIDEnd]) {
		t.Errorf("recipientID is not properly mapped to serial"+
			"\n\treceived: %v\n\texpected: %v",
			ad.GetRecipientID(), ad.serial[recipientIDStart:recipientIDEnd])
	}

	fp2 := NewFingerprint(ad.serial[keyFPStart:keyFPEnd])
	if !reflect.DeepEqual(ad.GetKeyFP(), *fp2) {
		t.Errorf("keyFP is not properly mapped to serial"+
			"\n\treceived: %v\n\texpected: %v",
			ad.keyFP, *fp2)
	}

	if !bytes.Equal(ad.GetTimestamp(), ad.serial[timestampStart:timestampEnd]) {
		t.Errorf("timestamp is not properly mapped to serial"+
			"\n\treceived: %v\n\texpected: %v",
			ad.GetTimestamp(), ad.serial[timestampStart:timestampEnd])
	}

	if !bytes.Equal(ad.GetMAC(), ad.serial[macStart:macEnd]) {
		t.Errorf("mac is not properly mapped to serial"+
			"\n\treceived: %v\n\texpected: %v",
			ad.GetMAC(), ad.serial[macStart:macEnd])
	}
}

// Tests that AssociatedData is always constructed the same way.
func TestAssociatedData_Consistency(t *testing.T) {
	// Define the serial to check against (64-bit encoded)
	serial, _ := base64.StdEncoding.DecodeString("E90arATOLqKHfFV5z6LHjhs" +
		"Lr66IG4KnURCKQu08kDyqQ0ZaeGIGFpeK7QzjxsTzrnvD4ElbVxL+/b4MECiH4QDazS" +
		"2IX2kstgfaAKEcHHBx55ai3C3CWlt0suEpcF4nPwXJIyaCjisFbjgXZY4QYQ==")

	// Generate random byte slice
	rand.Seed(42)
	randAD := make([]byte, associatedDataLen)
	rand.Read(randAD)
	randRecipientID := make([]byte, recipientIDLen)
	rand.Read(randRecipientID)
	randKeyFP := make([]byte, keyFPLen)
	rand.Read(randKeyFP)
	fp := NewFingerprint(randKeyFP)
	randTimestamp := make([]byte, timestampLen)
	rand.Read(randTimestamp)
	randMAC := make([]byte, macLen)
	rand.Read(randMAC)

	// Create new AssociatedData and set fields
	ad := NewAssociatedData(randAD)
	ad.SetRecipientID(randRecipientID)
	ad.SetKeyFP(*fp)
	ad.SetTimestamp(randTimestamp)
	ad.SetMAC(randMAC)

	if !bytes.Equal(ad.Get(), serial) {
		t.Errorf("AssociatedData's serial does not match the hardcoded serial"+
			"\n\treceived: %v\n\texpected: %v",
			ad.Get(), serial)
	}
}

// Tests that Get() returns the correct bytes set to AssociatedData's serial.
func TestAssociatedData_Get(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, associatedDataLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(randSlice)

	if !bytes.Equal(ad.Get(), randSlice) {
		t.Errorf("Get() did not return the correct data from "+
			"AssociatedData's serial"+
			"\n\treceived: %v\n\texpected: %v",
			ad.Get(), randSlice)
	}
}

// Tests that Set() sets the correct bytes to AssociatedData's serial and copies
// the correct number of bytes.
func TestAssociatedData_Set(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, associatedDataLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(make([]byte, associatedDataLen))

	// Set AssociatedData's serial
	size := ad.Set(randSlice)

	if !bytes.Equal(ad.Get(), randSlice) {
		t.Errorf("Set() did not properly set AssociatedData's serial"+
			"\n\treceived: %v\n\texpected: %v",
			ad.Get(), randSlice)
	}

	if size != associatedDataLen {
		t.Errorf("Set() did not copy the correct number of bytes into "+
			"AssociatedData's serial"+
			"\n\treceived: %v\n\texpected: %v",
			size, associatedDataLen)
	}
}

// Tests that Set() panics when the new serial is not the same length as serial.
func TestAssociatedData_Set_Panic(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	serialRand := make([]byte, associatedDataLen-5)
	rand.Read(serialRand)

	// Create new AssociatedData
	ad := NewAssociatedData(make([]byte, associatedDataLen))

	// Defer to an error when Set() does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Set() did not panic when expected")
		}
	}()

	// Set AssociatedData's serial
	ad.Set(serialRand)
}

// Tests that GetRecipientID() returns the correct bytes set to AssociatedData's
// recipientID.
func TestAssociatedData_GetRecipientID(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, associatedDataLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(randSlice)

	if !bytes.Equal(ad.GetRecipientID(), randSlice[recipientIDStart:recipientIDEnd]) {
		t.Errorf("GetRecipientID() did not return the correct data from "+
			"AssociatedData's recipientID"+
			"\n\treceived: %v\n\texpected: %v",
			ad.GetRecipientID(), randSlice[recipientIDStart:recipientIDEnd])
	}
}

// Tests that SetRecipientID() sets the correct bytes to AssociatedData's
// recipientID and copies the correct number of bytes.
func TestAssociatedData_SetRecipientID(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, recipientIDLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(make([]byte, associatedDataLen))

	// Set AssociatedData's recipientID
	size := ad.SetRecipientID(randSlice)

	if !bytes.Equal(ad.GetRecipientID(), randSlice) {
		t.Errorf("SetRecipientID() did not properly set AssociatedData's "+
			"recipientID"+
			"\n\treceived: %v\n\texpected: %v",
			ad.GetRecipientID(), randSlice)
	}

	if size != recipientIDLen {
		t.Errorf("SetRecipientID() did not copy the correct number of "+
			"bytes into AssociatedData's recipientID"+
			"\n\treceived: %v\n\texpected: %v",
			size, recipientIDLen)
	}
}

// Tests that SetRecipientID() panics when the new recipient ID is not the same
// length as recipientID.
func TestAssociatedData_SetRecipientID_Panic(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	serialRand := make([]byte, recipientIDLen-5)
	rand.Read(serialRand)

	// Create new AssociatedData
	ad := NewAssociatedData(make([]byte, associatedDataLen))

	// Defer to an error when SetRecipientID() does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetRecipientID() did not panic when expected")
		}
	}()

	// Set AssociatedData's recipientID
	ad.SetRecipientID(serialRand)
}

// Tests that GetRecipient() returns AssociatedData's recipientID as a *id.User.
func TestAssociatedData_GetRecipient(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, associatedDataLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(randSlice)

	userID := id.NewUserFromBytes(randSlice[recipientIDStart:recipientIDEnd])

	if !reflect.DeepEqual(ad.GetRecipient(), userID) {
		t.Errorf("GetRecipient() did not return the correct data from "+
			"AssociatedData's recipientID as a *id.User"+
			"\n\treceived: %v\n\texpected: %v",
			ad.GetRecipientID(), userID)
	}
}

// Tests that SetRecipient() sets the correct bytes from a *id.User to
// AssociatedData's recipientID and copies the correct number of bytes.
func TestAssociatedData_SetRecipient(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, recipientIDLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(make([]byte, associatedDataLen))

	// Set AssociatedData's recipient
	userID := id.NewUserFromBytes(randSlice)
	size := ad.SetRecipient(userID)

	if !reflect.DeepEqual(ad.GetRecipient(), userID) {
		t.Errorf("SetRecipient() did not properly set AssociatedData's "+
			"recipientID from a *id.User"+
			"\n\treceived: %v\n\texpected: %v",
			ad.GetRecipient(), userID)
	}

	if size != recipientIDLen {
		t.Errorf("SetRecipient() did not copy the correct number of "+
			"bytes into AssociatedData's recipientID"+
			"\n\treceived: %v\n\texpected: %v",
			size, recipientIDLen)
	}
}

// Tests that NewFingerprint() creates a new Fingerprint with the correct data.
func TestNewFingerprint(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, keyFPLen)
	rand.Read(randSlice)

	// Create Fingerprint
	fp := NewFingerprint(randSlice)

	if !bytes.Equal(fp[:], randSlice) {
		t.Errorf("NewFingerprint() did not copy the correct data to a "+
			"new Fingerprint"+
			"\n\treceived: %v\n\texpected: %v",
			fp[:], randSlice)
	}
}

// Tests that NewFingerprint() panics when the new data is not the same size as
// keyFP.
func TestNewFingerprint_Panic(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, keyFPLen-5)
	rand.Read(randSlice)

	// Defer to an error when NewFingerprint() does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("NewFingerprint() did not panic when expected")
		}
	}()

	// Create Fingerprint
	NewFingerprint(randSlice)
}

// Tests that GetKeyFP() returns AssociatedData's keyFP as a Fingerprint
func TestAssociatedData_GetKeyFP(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, associatedDataLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(randSlice)

	fp := NewFingerprint(randSlice[keyFPStart:keyFPEnd])

	if !reflect.DeepEqual(ad.GetKeyFP(), *fp) {
		t.Errorf("GetKeyFP() did not return the correct data from "+
			"AssociatedData's keyFP as a Fingerprint"+
			"\n\treceived: %v\n\texpected: %v",
			ad.GetKeyFP(), *fp)
	}
}

// Tests that SetKeyFP() sets the correct bytes from a Fingerprint to
// AssociatedData's keyFP and copies the correct number of bytes.
func TestAssociatedData_SetKeyFP(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, recipientIDLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(make([]byte, associatedDataLen))

	// Set AssociatedData's keyFP
	fp := NewFingerprint(randSlice)
	size := ad.SetKeyFP(*fp)

	if !reflect.DeepEqual(ad.GetKeyFP(), *fp) {
		t.Errorf("SetKeyFP() did not properly set AssociatedData's "+
			"keyFP from a Fingerprint"+
			"\n\treceived: %v\n\texpected: %v",
			ad.GetKeyFP(), *fp)
	}

	if size != keyFPLen {
		t.Errorf("SetKeyFP() did not copy the correct number of "+
			"bytes into AssociatedData's keyFP"+
			"\n\treceived: %v\n\texpected: %v",
			size, keyFPLen)
	}
}

// Tests that GetTimestamp() returns the correct bytes set to AssociatedData's
// timestamp.
func TestAssociatedData_GetTimestamp(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, associatedDataLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(randSlice)

	if !bytes.Equal(ad.GetTimestamp(), randSlice[timestampStart:timestampEnd]) {
		t.Errorf("GetTimestamp() did not return the correct data from "+
			"AssociatedData's timestamp"+
			"\n\treceived: %v\n\texpected: %v",
			ad.GetTimestamp(), randSlice[timestampStart:timestampEnd])
	}
}

// Tests that SetTimestamp() sets the correct bytes to AssociatedData's
// timestamp and copies the correct number of bytes.
func TestAssociatedData_SetTimestamp(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, timestampLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(make([]byte, associatedDataLen))

	// Set AssociatedData's timestamp
	size := ad.SetTimestamp(randSlice)

	if !bytes.Equal(ad.GetTimestamp(), randSlice) {
		t.Errorf("SetTimestamp() did not properly set AssociatedData's "+
			"timestamp"+
			"\n\treceived: %v\n\texpected: %v",
			ad.GetTimestamp(), randSlice)
	}

	if size != timestampLen {
		t.Errorf("SetTimestamp() did not copy the correct number of "+
			"bytes into AssociatedData's timestamp"+
			"\n\treceived: %v\n\texpected: %v",
			size, timestampLen)
	}
}

// Tests that SetTimestamp() panics when the new recipient ID is not the same
// length as timestamp.
func TestAssociatedData_SetTimestamp_Panic(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	serialRand := make([]byte, timestampLen-5)
	rand.Read(serialRand)

	// Create new AssociatedData
	ad := NewAssociatedData(make([]byte, associatedDataLen))

	// Defer to an error when SetTimestamp() does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetTimestamp() did not panic when expected")
		}
	}()

	// Set AssociatedData's timestamp
	ad.SetTimestamp(serialRand)
}

// Tests that GetMAC() returns the correct bytes set to AssociatedData's
// mac.
func TestAssociatedData_GetMAC(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, associatedDataLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(randSlice)

	if !bytes.Equal(ad.GetMAC(), randSlice[macStart:macEnd]) {
		t.Errorf("GetMAC() did not return the correct data from "+
			"AssociatedData's mac"+
			"\n\treceived: %v\n\texpected: %v",
			ad.GetMAC(), randSlice[macStart:macEnd])
	}
}

// Tests that SetMAC() sets the correct bytes to AssociatedData's
// mac and copies the correct number of bytes.
func TestAssociatedData_SetMAC(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, macLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(make([]byte, associatedDataLen))

	// Set AssociatedData's mac
	size := ad.SetMAC(randSlice)

	if !bytes.Equal(ad.GetMAC(), randSlice) {
		t.Errorf("SetMAC() did not properly set AssociatedData's "+
			"mac"+
			"\n\treceived: %v\n\texpected: %v",
			ad.GetMAC(), randSlice)
	}

	if size != macLen {
		t.Errorf("SetMAC() did not copy the correct number of "+
			"bytes into AssociatedData's mac"+
			"\n\treceived: %v\n\texpected: %v",
			size, macLen)
	}
}

// Tests that SetMAC() panics when the new recipient ID is not the same
// length as mac.
func TestAssociatedData_SetMAC_Panic(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	serialRand := make([]byte, macLen-5)
	rand.Read(serialRand)

	// Create new AssociatedData
	ad := NewAssociatedData(make([]byte, associatedDataLen))

	// Defer to an error when SetMAC() does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetMAC() did not panic when expected")
		}
	}()

	// Set AssociatedData's mac
	ad.SetMAC(serialRand)
}

// Tests that changes made to a copy by DeepCopy() does not reflect to the
// original contents.
func TestAssociatedData_DeepCopy(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, associatedDataLen)
	randSlice2 := make([]byte, associatedDataLen)
	rand.Read(randSlice)

	// Create new AssociatedData
	ad := NewAssociatedData(randSlice)

	// Create copy and change the serial
	adCopy := ad.DeepCopy()
	rand.Read(randSlice2)
	adCopy.Set(randSlice2)

	if bytes.Equal(ad.serial, adCopy.serial) {
		t.Errorf("DeepCopy() did not properly create a new copy of AssociatedData"+
			"\n\treceived: %v\n\texpected: %v",
			ad.serial, adCopy.serial)
	}
}
