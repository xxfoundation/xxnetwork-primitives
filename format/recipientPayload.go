////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"errors"
	"fmt"
	"gitlab.com/elixxir/primitives/userid"
)

const (
	// Length and Position of the Recipient Initialization Vector
	RIV_LEN   uint64 = 9
	RIV_START uint64 = 0
	RIV_END   uint64 = RIV_LEN

	// Length and Position of the Recipient ID
	RID_LEN   uint64 = userid.UserIDLen
	RID_START uint64 = REMPTY_END
	RID_END   uint64 = RID_START + RID_LEN

	// Length and Position of the Recipient MIC
	RMIC_LEN   uint64 = 8
	RMIC_START uint64 = RID_END
	RMIC_END   uint64 = RMIC_START + RMIC_LEN

	// Length of unused region in recipient payload
	REMPTY_LEN   uint64 = TOTAL_LEN - RIV_LEN - RMIC_LEN - RID_LEN
	REMPTY_START uint64 = RIV_END
	REMPTY_END   uint64 = REMPTY_START + REMPTY_LEN
)

// Structure containing the components of the recipient payload
type RecipientPayload struct {
	recipientSerial   [TOTAL_LEN]byte
	recipientInitVect []byte
	recipientID       []byte
	recipientMIC      []byte
}

//Builds a recipient payload object
func NewRecipientPayload(ID *userid.UserID) (*RecipientPayload, error) {
	if ID == nil || *ID == *userid.ZeroID {
		return nil, errors.New(fmt.Sprintf(
			"Cannot build Recipient Payload; Invalid Recipient ID: %q",
			ID))
	}
	result := RecipientPayload{recipientSerial: [TOTAL_LEN]byte{}}
	result.recipientID = result.recipientSerial[RID_START:RID_END]
	copy(result.recipientID, ID.Bytes())
	result.recipientInitVect = result.recipientSerial[RIV_START:RIV_END]
	result.recipientMIC = result.recipientSerial[RMIC_START:RMIC_END]

	return &result, nil
}

// This function returns the recipient ID slice
// The caller can read or write the data within this slice, but can't change
// the slice header in the actual structure
func (r *RecipientPayload) GetRecipientID() []byte {
	return r.recipientID
}

// Get the recipient initialization vector
// The caller can read or write the data within this slice, but can't change
// the slice header in the actual structure
func (r *RecipientPayload) GetRecipientInitVect() []byte {
	return r.recipientInitVect
}

// Get the recipient MIC
// The caller can read or write the data within this slice, but can't change
// the slice header in the actual structure
func (r *RecipientPayload) GetRecipientMIC() []byte {
	return r.recipientMIC
}

// Returns the serialized recipient payload, without copying
func (r *RecipientPayload) SerializeRecipient() []byte {
	return r.recipientSerial[:]
}

// Slices a serialized recipient ID into its constituent fields
func DeserializeRecipient(rSerial []byte) *RecipientPayload {
	var rBytes [TOTAL_LEN]byte
	copy(rBytes[:], rSerial)
	return &RecipientPayload{
		rBytes,
		rBytes[RIV_START:RIV_END],
		rBytes[RID_START:RID_END],
		rBytes[RMIC_START:RMIC_END],
	}

}

// Creates a deep copy of the recipient, used for sending multiple messages
func (r *RecipientPayload) DeepCopy() *RecipientPayload {
	return DeserializeRecipient(r.recipientSerial[:])
}
