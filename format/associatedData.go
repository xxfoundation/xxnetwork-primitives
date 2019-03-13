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
	// First byte to ensure that the associated data are always within the
	// cyclic group
	AD_FIRST_LEN   int = 1
	AD_FIRST_START int = 0
	AD_FIRST_END   int = AD_FIRST_START + AD_FIRST_LEN

	// Length and position of the Recipient ID
	AD_RID_LEN   int = id.UserLen
	AD_RID_START int = AD_FIRST_END
	AD_RID_END   int = AD_RID_START + AD_RID_LEN

	// Length and position of the key fingerprint
	AD_KEYFP_LEN   int = 32
	AD_KEYFP_START int = AD_RID_END
	AD_KEYFP_END   int = AD_KEYFP_START + AD_KEYFP_LEN

	// Length and position of the encrypted timestamp
	// 128 bits, seconds+nanoseconds
	// Encrypt as one AES block
	AD_TIMESTAMP_LEN   int = 16
	AD_TIMESTAMP_START int = AD_KEYFP_END
	AD_TIMESTAMP_END   int = AD_TIMESTAMP_START + AD_TIMESTAMP_LEN

	// Length and Position of the MAC
	AD_MAC_LEN   int = 32
	AD_MAC_START int = AD_TIMESTAMP_END
	AD_MAC_END   int = AD_MAC_START + AD_MAC_LEN

	// TODO Delete this when the third phase has been removed
	// Length and position of the recipient MIC
	AD_RMIC_LEN   int = 32
	AD_RMIC_START int = AD_MAC_END
	AD_RMIC_END   int = AD_RMIC_START + AD_RMIC_LEN

	// Length of unused region in recipient payloadData
	// TODO @mario Should the empty data go at the end or in the middle
	// somewhere? Should this be PKCS padding instead?
	AD_EMPTY_LEN   int = TOTAL_LEN - AD_RID_LEN - AD_KEYFP_LEN - AD_MAC_LEN - AD_FIRST_LEN - AD_TIMESTAMP_LEN - AD_RMIC_LEN
	AD_EMPTY_START int = AD_RMIC_END
	AD_EMPTY_END   int = AD_EMPTY_START + AD_EMPTY_LEN
)

// Structure containing the components of the recipient payloadData
type AssociatedData struct {
	associatedDataSerial [TOTAL_LEN]byte
	recipientID          []byte
	keyFingerprint       [AD_KEYFP_LEN]byte
	timestamp            []byte
	mac                  []byte
	rmic                 []byte
}

// Initializes an Associated data with the correct slices
func NewAssociatedData() *AssociatedData {
	result := AssociatedData{
		associatedDataSerial: [TOTAL_LEN]byte{},
		keyFingerprint: [AD_KEYFP_LEN]byte{},
	}
	result.recipientID = result.associatedDataSerial[AD_RID_START:AD_RID_END]

	result.timestamp = result.associatedDataSerial[AD_TIMESTAMP_START:AD_TIMESTAMP_END]
	result.mac = result.associatedDataSerial[AD_MAC_START:AD_MAC_END]
	result.rmic = result.associatedDataSerial[AD_RMIC_START:AD_RMIC_END]

	ensureGroup(result.associatedDataSerial[AD_FIRST_START:AD_FIRST_END])

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
	return r.keyFingerprint[:]
}

// Returns number of bytes copied
func (r *AssociatedData) SetKeyFingerprint(newKeyFP []byte) int {
	return copy(r.keyFingerprint[:], newKeyFP)
}

// Get the MAC for the message
// The caller can read or write the data within this slice, but can't change
// the slice header in the actual structure
func (r *AssociatedData) GetMAC() []byte {
	return r.mac
}

// Returns number of bytes copied
func (r *AssociatedData) SetMAC(newMAC []byte) int {
	return copy(r.mac, newMAC)
}

// Get the MIC for the recipient ID
func (r *AssociatedData) GetRecipientMIC() []byte {
	return r.rmic
}

// Returns number of bytes copied
func (r *AssociatedData) SetRecipientMIC(newRecipientMIC []byte) int {
	return copy(r.rmic, newRecipientMIC)
}

// Get the message's timestamp
func (r *AssociatedData) GetTimestamp() []byte {
	return r.timestamp
}

func (r *AssociatedData) SetTimestamp(newTimestamp []byte) int {
	return copy(r.timestamp, newTimestamp)
}

// Returns the serialized Associated Data, without copying
func (r *AssociatedData) SerializeAssociatedData() []byte {
	return r.associatedDataSerial[:]
}

// Slices a serialized Associated Data into its constituent fields
func DeserializeAssociatedData(rSerial []byte) *AssociatedData {
	result := NewAssociatedData()
	copy(result.associatedDataSerial[:], rSerial)
	return result
}

// Creates a deep copy of the Associated Data, used for sending multiple messages
func (r *AssociatedData) DeepCopy() *AssociatedData {
	return DeserializeAssociatedData(r.associatedDataSerial[:])
}
