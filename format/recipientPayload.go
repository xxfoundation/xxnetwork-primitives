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
	recipientInitVect []byte
	recipientEmpty    []byte
	recipientID       []byte
	recipientMIC      []byte
}

//Builds a recipient payload object
func NewRecipientPayload(ID *userid.UserID) (*Recipient, error) {
	if *ID == *userid.ZeroID {
		return &Recipient{}, errors.New(fmt.Sprintf(
			"Cannot build Recipient Payload; Invalid Recipient ID: %q",
			ID))
	}

	return &Recipient{
		make([]byte, RIV_LEN),
		make([]byte, REMPTY_LEN),
		ID.Bytes(),
		make([]byte, RMIC_LEN),
	}, nil
}

// This function returns a pointer to the recipient ID
// This ensures that while the data can be edited, it cant be reallocated
func (r Recipient) GetRecipientID() []byte {
	return r.recipientID
}

// This function returns a pointer to the unused component of the recipient ID
// This ensures that while the data can be edited, it cant be reallocated
func (r Recipient) GetRecipientEmpty() []byte {
	return r.recipientEmpty
}

// This function returns a pointer to the recipient Initialization Vector
// This ensures that while the data can be edited, it cant be reallocated
func (r Recipient) GetRecipientInitVect() []byte {
	return r.recipientInitVect
}

// This function returns a pointer to the recipient Initialization MIC
// This ensures that while the data can be edited, it cant be reallocated
func (r Recipient) GetRecipientMIC() []byte {
	return r.recipientMIC
}

// Returns the serialized recipient payload
func (r Recipient) SerializeRecipient() []byte {
	rbytes := make([]byte, TOTAL_LEN)

	//Copy the Recipient Initialization Vector into the serialization
	copy(rbytes[RIV_START:RIV_END], r.recipientInitVect[:])

	//Copy the empty region into the serialization
	copy(rbytes[REMPTY_START:REMPTY_END], r.recipientEmpty[:])

	//Copy the recipient ID into the serialization
	copy(rbytes[RID_START:RID_END], r.recipientID[:])

	//Copy the Recipient MIC into the serialization
	copy(rbytes[RMIC_START:RMIC_END], r.recipientMIC[:])

	//Make sure the highest bit of the serialization is zero
	rbytes[0] = rbytes[0] & ZEROER

	return rbytes
}

//Returns a Deserialized recipient id
func DeserializeRecipient(rSerial []byte) Recipient {
	return Recipient{
		rSerial[RIV_START:RIV_END],
		rSerial[REMPTY_START:REMPTY_END],
		rSerial[RID_START:RID_END],
		rSerial[RMIC_START:RMIC_END],
	}

}

// Creates a deep copy of the recipient, used for sending multiple messages
func (r Recipient) DeepCopy() Recipient {
	riv := make([]byte, RIV_LEN)
	rempty := make([]byte, REMPTY_LEN)
	rid := make([]byte, RID_LEN)
	rmic := make([]byte, RMIC_LEN)
	copy(riv, r.GetRecipientInitVect())
	copy(rempty, r.GetRecipientEmpty())
	copy(rid, r.GetRecipientID())
	copy(rmic, r.GetRecipientMIC())

	return Recipient{
		riv,
		rempty,
		rid,
		rmic,
	}
}
