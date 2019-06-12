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
	associatedDataLen   = 112 // 896 bits
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
	macLen   = 32 // 256 bits
	macStart = timestampEnd
	macEnd   = macStart + macLen
)

// Structure for the associated data section of the message points to a
// subsection of the serialised Message structure.
type AssociatedData struct {
	serial      []byte // points to region in master
	recipientID []byte
	keyFP       []byte // key fingerprint
	timestamp   []byte
	mac         []byte // message authentication code
}

// Array form for storing a fingerprint
type Fingerprint [keyFPLen]byte

// NewAssociatedData creates a new AssociatedData for a message and points
// recipientID, keyFP, timestamp, mac, and grpByte to serial. If the new serial
// is not exactly the same length as serial, then it panics.
func NewAssociatedData(newSerial []byte) *AssociatedData {
	if len(newSerial) != associatedDataLen {
		panic("new serial not the same size as Associated Data serial")
	}

	newAD := &AssociatedData{}
	newAD.serial = newSerial
	newAD.recipientID = newAD.serial[recipientIDStart:recipientIDEnd]
	newAD.keyFP = newAD.serial[keyFPStart:keyFPEnd]
	newAD.timestamp = newAD.serial[timestampStart:timestampEnd]
	newAD.mac = newAD.serial[macStart:macEnd]

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
	if len(newSerial) != associatedDataLen {
		panic("new serial not the same size as AssociatedData serial")
	}

	return copy(a.serial, newSerial)
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
	if len(newRecipientID) != recipientIDLen {
		panic("new recipientID not the same size as AssociatedData newRecipientID")
	}

	return copy(a.recipientID, newRecipientID)
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
	if len(newTimestamp) != timestampLen {
		panic("new timestamp not the same size as AssociatedData timestamp")
	}

	return copy(a.timestamp, newTimestamp)
}

// GetMAC returns the mac. The caller can read or write the data within this
// slice, but cannot change the slice header in the actual structure.
func (a *AssociatedData) GetMAC() []byte {
	return a.mac
}

// SetMac sets the mac. The number of bytes copied is returned. If the specified
// byte array is not exactly the same size as mac, then it panics.
func (a *AssociatedData) SetMAC(newMAC []byte) int {
	if len(newMAC) != macLen {
		panic("new timestamp not the same size as AssociatedData timestamp")
	}

	return copy(a.mac, newMAC)
}

// DeepCopy creates a copy of AssociatedData.
func (a *AssociatedData) DeepCopy() *AssociatedData {
	newCopy := NewAssociatedData(make([]byte, associatedDataLen))
	copy(newCopy.serial[:], a.serial)

	return newCopy
}

// NewFingerprint creates a new fingerprint from a byte slice. If the specified
// data iis not exactly the same size as keyFP, then it panics.
func NewFingerprint(data []byte) *Fingerprint {
	if len(data) != keyFPLen {
		panic("data is not smaller than or equal to AssociatedData keyFP")
	}

	fp := &Fingerprint{}
	copy(fp[:], data[:])
	return fp
}
