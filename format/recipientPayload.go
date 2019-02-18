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
	// Length and Position of the Recipeint Initialization Vector
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
type Recipient struct {
	recipientSerial   [TOTAL_LEN]byte
	recipientInitVect []byte
	recipientID       []byte
	recipientMIC      []byte
}

//Builds a recipient payload object
func NewRecipientPayload(ID *userid.UserID) (*Recipient, error) {
	if ID == nil || *ID == *userid.ZeroID {
		return nil, errors.New(fmt.Sprintf(
			"Cannot build Recipient Payload; Invalid Recipient ID: %q",
			ID))
	}
	result := Recipient{recipientSerial: [TOTAL_LEN]byte{}}
	result.recipientID = result.recipientSerial[RID_START:RID_END]
	copy(result.recipientID, ID.Bytes())
	result.recipientInitVect = result.recipientSerial[RIV_START:RIV_END]
	result.recipientMIC = result.recipientSerial[RMIC_START:RMIC_END]

	return &result, nil
}

// This function returns a pointer to the recipient ID
// This ensures that while the data can be edited, it cant be reallocated
func (r *Recipient) GetRecipientID() []byte {
	return r.recipientID
}

// This function returns a pointer to the recipient Initialization Vector
// This ensures that while the data can be edited, it cant be reallocated
func (r *Recipient) GetRecipientInitVect() []byte {
	return r.recipientInitVect
}

// This function returns a pointer to the recipient Initialization MIC
// This ensures that while the data can be edited, it cant be reallocated
func (r *Recipient) GetRecipientMIC() []byte {
	return r.recipientMIC
}

// Returns the serialized recipient payload, without copying
func (r *Recipient) SerializeRecipient() []byte {
	return r.recipientSerial[:]
}

//Returns a Deserialized recipient id
func DeserializeRecipient(rSerial []byte) *Recipient {
	var rBytes [TOTAL_LEN]byte
	copy(rBytes[:], rSerial)
	return &Recipient{
		rBytes,
		rBytes[RIV_START:RIV_END],
		rBytes[RID_START:RID_END],
		rBytes[RMIC_START:RMIC_END],
	}

}

// Creates a deep copy of the recipient, used for sending multiple messages
func (r *Recipient) DeepCopy() *Recipient {
	return DeserializeRecipient(r.recipientSerial[:])
}
