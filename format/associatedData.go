////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"gitlab.com/elixxir/primitives/id"
)

const (
	// Length, start index, and end index of the Associated Data serial
	associatedDataLen   = 113 // 904 bits
	associatedDataStart = contentsEnd
	associatedDataEnd   = associatedDataStart + associatedDataLen

	// Length, start index, and end index of recipientID
	recipientIDLen   = 32 // 256 bits
	recipientIDStart = 0
	recipientIDEnd   = recipientIDStart + recipientIDLen

	// Length, start index, and end index of keyFP
	keyFPLen   = 32 // 256 bits
	keyFPStart = recipientIDEnd
	keyFPEnd   = keyFPStart + keyFPLen

	// Length, start index, and end index of timestamp
	timestampLen   = 16 // 128 bits
	timestampStart = keyFPEnd
	timestampEnd   = timestampStart + timestampLen

	// Length, start index, and end index of mac
	macLen   = 16 // 128 bits
	macStart = timestampEnd
	macEnd   = macStart + macLen

	// Length, start index, and end index of grpByte
	grpByteLen   = 1 // 8 bits
	grpByteStart = macEnd
	grpByteEnd   = grpByteStart + grpByteLen
)

// Structure for the associated data section of the message points to a
// subsection of the serialised Message structure.
type AssociatedData struct {
	serial      []byte // points to region in master
	recipientID []byte
	keyFP       []byte // key fingerprint
	timestamp   []byte
	mac         []byte // message authentication code
	grpByte     []byte // zero value byte ensures payloadB can be in the group
}

// Length of the key fingerprint
const KeyFPLen = 32

// Array form for storing a fingerprint
type Fingerprint [KeyFPLen]byte

// NewAssociatedData creates a new AssociatedData for a message and points
// recipientID, keyFP, timestamp, mac, and grpByte to serial. If the new serial
// is not exactly the same length as serial, then it panics.
func NewAssociatedData(newSerial []byte) *AssociatedData {
	newAD := &AssociatedData{}

	if len(newSerial) == associatedDataLen {
		newAD.serial = newSerial
	} else {
		panic("new serial not the same size as Associated Data serial")
	}

	newAD.recipientID = newAD.serial[recipientIDStart:recipientIDEnd]
	newAD.keyFP = newAD.serial[keyFPStart:keyFPEnd]
	newAD.timestamp = newAD.serial[timestampStart:timestampEnd]
	newAD.mac = newAD.serial[macStart:macEnd]
	newAD.grpByte = newAD.serial[grpByteStart:grpByteEnd]
	newAD.grpByte[0] = 0

	return newAD
}

// Get returns the AssociatedData's serialised data. The caller can read or
// write the data within this slice, but cannot change the slice header in the
// actual structure.
func (a *AssociatedData) Get() []byte {
	return a.serial
}

// Set sets the entire content of associated data. The number of bytes copied is
// returned. If the specified byte array is not exactly the same size as serial,
// then it panics.
func (a *AssociatedData) Set(newSerial []byte) int {
	if len(newSerial) == associatedDataLen {
		return copy(a.serial, newSerial)
	} else {
		panic("new serial not the same size as AssociatedData serial")
	}
}

// GetRecipientID returns the recipientID. The caller can read or write the data
// within this slice, but cannot change the slice header in the actual
// structure.
func (a *AssociatedData) GetRecipientID() []byte {
	return a.recipientID
}

// SetRecipientID sets the recipientID. The number of bytes copied is returned.
// If the specified byte array is not exactly the same size as recipientID, then
// it panics.
func (a *AssociatedData) SetRecipientID(newRecipientID []byte) int {
	if len(newRecipientID) == recipientIDLen {
		return copy(a.recipientID, newRecipientID)
	} else {
		panic("new recipientID not the same size as AssociatedData newRecipientID")
	}
}

// GetRecipient returns the recipientID as a user ID.
func (a *AssociatedData) GetRecipient() *id.User {
	return id.NewUserFromBytes(a.recipientID)
}

// SetRecipient sets the value of recipientID from a user ID. The number of
// bytes copied is returned.
func (a *AssociatedData) SetRecipient(newRecipientID *id.User) int {
	return copy(a.recipientID, newRecipientID.Bytes())
}

// GetKeyFP returns the keyFP as a Fingerprint.
func (a *AssociatedData) GetKeyFP() (fp Fingerprint) {
	copy(fp[:], a.keyFP)
	return fp
}

// SetKeyFP sets the keyFP from a Fingerprint. The number of bytes copied is
// returned.
func (a *AssociatedData) SetKeyFP(fp Fingerprint) int {
	return copy(a.keyFP, fp[:])
}

// GetTimestamp returns the timestamp. The caller can read or write the data
// within this slice, but cannot change the slice header in the actual
// structure.
func (a *AssociatedData) GetTimestamp() []byte {
	return a.timestamp
}

// SetTimestamp sets the timestamp. The number of bytes copied is returned. If
// the specified byte array is not exactly the same size as timestamp, then it
// panics.
func (a *AssociatedData) SetTimestamp(newTimestamp []byte) int {
	if len(newTimestamp) == timestampLen {
		return copy(a.timestamp, newTimestamp)
	} else {
		panic("new timestamp not the same size as AssociatedData timestamp")
	}
}

// GetMAC returns the mac. The caller can read or write the data within this
// slice, but cannot change the slice header in the actual structure.
func (a *AssociatedData) GetMAC() []byte {
	return a.mac
}

// SetMac sets the mac. The number of bytes copied is returned. If the specified
// byte array is not exactly the same size as mac, then it panics.
func (a *AssociatedData) SetMAC(newMAC []byte) int {
	if len(newMAC) == macLen {
		return copy(a.mac, newMAC)
	} else {
		panic("new timestamp not the same size as AssociatedData timestamp")
	}
}

// DeepCopy creates a copy of AssociatedData.
func (a *AssociatedData) DeepCopy() *AssociatedData {
	newCopy := NewAssociatedData(nil)
	copy(newCopy.serial[:], a.serial)

	return newCopy
}

// NewFingerprint creates a new fingerprint from a byte slice.
func NewFingerprint(data []byte) *Fingerprint {
	fp := &Fingerprint{}
	copy(fp[:], data[:])
	return fp
}
