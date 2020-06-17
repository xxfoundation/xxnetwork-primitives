////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

// Package format contains the message payload format. The structure can
// seen from this diagram. The entire message is 4096 bits with two payloads each 2048 bits.
// Underneath these payloads is the contents of the message and the associated data,
// which includes all necessary metadata for the message.
// The message structure is unique in that all the data is contained in a master byte
// slice with all sub structures pointing to different ranges within the master slice.
// This enables the message structure to always be quickly serialized.

/*                               Message Structure (not to scale)
+----------------------------------------------------------------------------------------+
|                                         Message                                        |
|                                        4096 bits                                       |
+----------------------------------------------------------------------------------------+
|                  payloadA                  |                 payloadB                  |
|                 2048 bits                  |                2048 bits                  |
+------------------------------------+-------+---------------------------------+---------+
|              Contents              |             AssociatedData              | grpByte |
|              3192 bits             |                896 bits                 | 8 bits  |
+------------------------------------+-----------------------------------------+         |
|     padding     |       data       | recipientID | keyFP | timestamp |  mac  |         |
|   88–3192 bits  |    0–3104 bits   |   256 bits  | 256 b |  128 bits | 256 b |         |
+-----------------+------------------+-------------+-------+-----------+-------+---------+
*/

package format

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/primitives/id"
)

const (
	// Length, start index, and end index of the Associated Data serial
	AssociatedDataLen   = 112 // 896 bits
	associatedDataStart = contentsEnd
	associatedDataEnd   = associatedDataStart + AssociatedDataLen

	// Length, start index, and end index of recipientID
	RecipientIDLen   = 32 // 256 bits
	recipientIDStart = 0
	recipientIDEnd   = recipientIDStart + RecipientIDLen

	// Length, start index, and end index of keyFP
	KeyFPLen   = 32 // 256 bits
	keyFPStart = recipientIDEnd
	keyFPEnd   = keyFPStart + KeyFPLen

	// Length, start index, and end index of timestamp
	TimestampLen   = 16 // 128 bits
	timestampStart = keyFPEnd
	timestampEnd   = timestampStart + TimestampLen

	// Length, start index, and end index of mac
	MacLen   = 32 // 256 bits
	macStart = timestampEnd
	macEnd   = macStart + MacLen
)

// AssociatedData is the structure for the associated data section of the
// message.
type AssociatedData struct {
	serial      []byte // points to region in master
	recipientID []byte
	keyFP       []byte // key fingerprint
	timestamp   []byte
	mac         []byte // message authentication code
}

// Fingerprint is the array form for storing a fingerprint
type Fingerprint [KeyFPLen]byte

// NewAssociatedData creates a new AssociatedData for a message and points
// recipientID, keyFP, timestamp, and mac to serial. If the new serial is not
// exactly the same length as serial, then it panics.
func NewAssociatedData(newSerial []byte) *AssociatedData {
	if len(newSerial) != AssociatedDataLen {
		jww.ERROR.Panicf("new serial not the same size as "+
			"AssociatedData serial; Expected: %v, Recieved: %v",
			AssociatedDataLen, len(newSerial))
	}

	newAD := &AssociatedData{
		serial:      newSerial[:],
		recipientID: newSerial[recipientIDStart:recipientIDEnd],
		keyFP:       newSerial[keyFPStart:keyFPEnd],
		timestamp:   newSerial[timestampStart:timestampEnd],
		mac:         newSerial[macStart:macEnd],
	}

	return newAD
}

// Get returns the AssociatedData's serialised data. The caller can read or
// write the data within this slice, but cannot change the slice header in the
// actual structure.
func (a *AssociatedData) Get() []byte {
	return a.serial
}

// Set sets the entire content of associated data. If the specified byte array
// is not exactly the same size as serial, then it panics.
func (a *AssociatedData) Set(newSerial []byte) {
	if len(newSerial) != AssociatedDataLen {
		jww.ERROR.Panicf("new serial not the same size as "+
			"AssociatedData serial; Expected: %v, Recieved: %v",
			AssociatedDataLen, len(newSerial))
	}

	copy(a.serial, newSerial)
}

// GetRecipientID returns the recipientID. The caller can read or write the data
// within this slice, but cannot change the slice header in the actual
// structure.
func (a *AssociatedData) GetRecipientID() []byte {
	return a.recipientID
}

// SetRecipientID sets the recipientID. If the specified byte array is not
// exactly the same size as recipientID, then it panics.
func (a *AssociatedData) SetRecipientID(newRecipientID []byte) {
	if len(newRecipientID) != RecipientIDLen {
		jww.ERROR.Panicf("new recipientID not the same size as "+
			"AssociatedData newRecipientID; Expected: %v, Recieved: %v",
			RecipientIDLen, len(newRecipientID))
	}

	copy(a.recipientID, newRecipientID)
}

// GetRecipient returns the recipientID as a user ID.
func (a *AssociatedData) GetRecipient() (*id.ID, error) {
	tempBytes := make([]byte, id.ArrIDLen)
	copy(tempBytes, a.recipientID)

	newID, err := id.Unmarshal(tempBytes)
	if err != nil {
		return nil, err
	}

	newID.SetType(id.User)

	return newID, nil
}

// SetRecipient sets the value of recipientID from a user ID.
func (a *AssociatedData) SetRecipient(newRecipientID *id.ID) {
	copy(a.recipientID, newRecipientID.Bytes())
}

// GetKeyFP returns the keyFP as a Fingerprint.
func (a *AssociatedData) GetKeyFP() (fp Fingerprint) {
	copy(fp[:], a.keyFP)
	return fp
}

// SetKeyFP sets the keyFP from a Fingerprint.
func (a *AssociatedData) SetKeyFP(fp Fingerprint) {
	copy(a.keyFP, fp[:])
}

// GetTimestamp returns the timestamp. The caller can read or write the data
// within this slice, but cannot change the slice header in the actual
// structure.
func (a *AssociatedData) GetTimestamp() []byte {
	return a.timestamp
}

// SetTimestamp sets the timestamp. If the specified byte array is not exactly
// the same size as timestamp, then it panics.
func (a *AssociatedData) SetTimestamp(newTimestamp []byte) {
	if len(newTimestamp) != TimestampLen {
		jww.ERROR.Panicf("new timestamp not the same size as "+
			"AssociatedData timestamp; Expected: %v, Recieved: %v",
			TimestampLen, len(newTimestamp))
	}

	copy(a.timestamp, newTimestamp)
}

// GetMAC returns the mac. The caller can read or write the data within this
// slice, but cannot change the slice header in the actual structure.
func (a *AssociatedData) GetMAC() []byte {
	return a.mac
}

// SetMAC sets the mac. If the specified byte array is not exactly the same size
// as mac, then it panics.
func (a *AssociatedData) SetMAC(newMAC []byte) {
	if len(newMAC) != MacLen {
		jww.ERROR.Panicf("new MAC not the same size as "+
			"AssociatedData MAC; Expected: %v, Recieved: %v",
			MacLen, len(newMAC))
	}

	copy(a.mac, newMAC)
}

// DeepCopy creates a copy of AssociatedData.
func (a *AssociatedData) DeepCopy() *AssociatedData {
	newCopy := NewAssociatedData(make([]byte, AssociatedDataLen))
	copy(newCopy.serial[:], a.serial)

	return newCopy
}

// NewFingerprint creates a new fingerprint from a byte slice. If the specified
// data iis not exactly the same size as keyFP, then it panics.
func NewFingerprint(data []byte) *Fingerprint {
	if len(data) != KeyFPLen {
		jww.ERROR.Panicf("fingerprint not the same size as "+
			"AssociatedData fingerprint; Expected: %v, Recieved: %v",
			KeyFPLen, len(data))
	}

	fp := &Fingerprint{}
	copy(fp[:], data[:])
	return fp
}
