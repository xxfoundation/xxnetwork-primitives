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
	// Length and position of the Recipient ID
	// The first bit of all user IDs should be zero
	AD_RID_LEN   int = id.UserLen
	AD_RID_START int = 0
	AD_RID_END   int = AD_RID_START + AD_RID_LEN

	// Length and position of the key fingerprint
	AD_KEYFP_LEN int = 32
	AD_KEYFP_START int = AD_RID_END
	AD_KEYFP_END int = AD_KEYFP_START + AD_KEYFP_LEN

	// Length and Position of the Recipient MAC
	AD_MAC_LEN   int = 32
	AD_MAC_START int = AD_KEYFP_END
	AD_MAC_END   int = AD_MAC_START + AD_MAC_LEN

	// Length of unused region in recipient payload
	// TODO @mario Should the empty data go at the end or in the middle
	// somewhere? Should this be PKCS padding instead?
	AD_EMPTY_LEN   int = TOTAL_LEN - AD_RID_LEN - AD_KEYFP_LEN - AD_MAC_LEN
	AD_EMPTY_START int = AD_RID_END
	AD_EMPTY_END   int = AD_EMPTY_START + AD_EMPTY_LEN
)

// Structure containing the components of the recipient payload
type AssociatedData struct {
	associatedDataSerial [TOTAL_LEN]byte
	recipientID          []byte
	keyFingerprint       []byte
	mac                  []byte
}

// Initializes an Associated data with the correct slices
func NewAssociatedData() (*AssociatedData) {
	result := AssociatedData{associatedDataSerial: [TOTAL_LEN]byte{}}
	result.recipientID = result.associatedDataSerial[AD_RID_START:AD_RID_END]
	result.keyFingerprint = result.associatedDataSerial[AD_KEYFP_START:AD_KEYFP_END]
	result.mac = result.associatedDataSerial[AD_MAC_START:AD_MAC_END]

	return &result
}

// This function returns the recipient ID slice
// The caller can read or write the data within this slice, but can't change
// the slice header in the actual structure
func (r *AssociatedData) GetRecipientID() []byte {
	return r.recipientID
}

func (r *AssociatedData) GetRecipient() *id.User {
	return new(id.User).SetBytes(r.recipientID)
}

// Returns number of bytes copied
func (r *AssociatedData) SetRecipientID(newID []byte) int {
	return copy(r.recipientID, newID)
}

func (r *AssociatedData) SetRecipient(newID *id.User) {
	copy(r.recipientID, newID.Bytes())
}

// Get the key fingerprint
// The caller can read or write the data within this slice, but can't change
// the slice header in the actual structure
func (r *AssociatedData) GetKeyFingerprint() []byte {
	return r.keyFingerprint
}

// Returns number of bytes copied
func (r *AssociatedData) SetKeyFingerprint(newKeyFP []byte) int {
	return copy(r.keyFingerprint, newKeyFP)
}

// Get the MAC for the associated data
// The caller can read or write the data within this slice, but can't change
// the slice header in the actual structure
func (r *AssociatedData) GetMAC() []byte {
	return r.mac
}

// Returns number of bytes copied
func (r *AssociatedData) SetMAC(newMAC []byte) int {
	return copy(r.mac, newMAC)
}

// Returns the serialized recipient payload, without copying
func (r *AssociatedData) SerializeAssociatedData() []byte {
	return r.associatedDataSerial[:]
}

// Slices a serialized recipient ID into its constituent fields
func DeserializeAssociatedData(rSerial []byte) *AssociatedData {
	result := NewAssociatedData()
	copy(result.associatedDataSerial[:], rSerial)
	return result
}

// Creates a deep copy of the recipient, used for sending multiple messages
func (r *AssociatedData) DeepCopy() *AssociatedData {
	return DeserializeAssociatedData(r.associatedDataSerial[:])
}
