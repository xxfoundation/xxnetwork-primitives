////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"errors"
	"fmt"
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/elixxir/primitives/userid"
)

const (
	// Length and Position of the Recipeint Initialization Vector
	RIV_LEN   uint64 = 9
	RIV_START uint64 = 0
	RIV_END   uint64 = RIV_LEN

	// Length and Position of the Recipient ID
	RID_LEN   uint64 = id.UserIDLen
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
	recipientInitVect *cyclic.Int
	recipientEmpty    *cyclic.Int
	recipientID       *cyclic.Int
	recipientMIC      *cyclic.Int
}

//Builds a recipient payload object
func NewRecipientPayload(ID *id.UserID) (Recipient, error) {
	if *ID == *id.ZeroID {
		return Recipient{}, errors.New(fmt.Sprintf(
			"Cannot build Recipient Payload; Invalid Recipient ID: %q",
			ID))
	}

	//TODO: initialize the components at their max lengths
	return Recipient{
		cyclic.NewInt(0),
		cyclic.NewInt(0),
		cyclic.NewIntFromBytes(ID[:]),
		cyclic.NewInt(0),
	}, nil
}

// This function returns a pointer to the recipient ID
// This ensures that while the data can be edited, it cant be reallocated
func (r Recipient) GetRecipientID() *cyclic.Int {
	return r.recipientID
}

// This function returns a pointer to the unused component of the recipient ID
// This ensures that while the data can be edited, it cant be reallocated
func (r Recipient) GetRecipientEmpty() *cyclic.Int {
	return r.recipientEmpty
}

// This function returns a pointer to the recipient Initialization Vector
// This ensures that while the data can be edited, it cant be reallocated
func (r Recipient) GetRecipientInitVect() *cyclic.Int {
	return r.recipientInitVect
}

// This function returns a pointer to the recipient Initialization MIC
// This ensures that while the data can be edited, it cant be reallocated
func (r Recipient) GetRecipientMIC() *cyclic.Int {
	return r.recipientMIC
}

// Returns the serialized recipient payload
// Returns as a cyclic int because it is expected that the message will be
// immediately encrypted via cyclic int multiplication
func (r Recipient) SerializeRecipient() *cyclic.Int {
	rbytes := make([]byte, TOTAL_LEN)

	//Copy the Recipient Initialization Vector into the serialization
	copy(rbytes[RIV_START:RIV_END], r.recipientInitVect.LeftpadBytes(RIV_LEN))

	//Copy the empty region into the serialization
	copy(rbytes[REMPTY_START:REMPTY_END], r.recipientEmpty.LeftpadBytes(
		REMPTY_LEN))

	//Copy the recipient ID into the serialization
	copy(rbytes[RID_START:RID_END], r.recipientID.LeftpadBytes(RID_LEN))

	//Copy the Recipient MIC into the serialization
	copy(rbytes[RMIC_START:RMIC_END], r.recipientMIC.LeftpadBytes(RMIC_LEN))

	//Make sure the highest bit of the serialization is zero
	rbytes[0] = rbytes[0] & ZEROER

	return cyclic.NewIntFromBytes(rbytes)
}

//Returns a Deserialized recipient id
func DeserializeRecipient(rSerial *cyclic.Int) Recipient {
	rbytes := rSerial.LeftpadBytes(TOTAL_LEN)

	return Recipient{
		cyclic.NewIntFromBytes(rbytes[RIV_START:RIV_END]),
		cyclic.NewIntFromBytes(rbytes[REMPTY_START:REMPTY_END]),
		cyclic.NewIntFromBytes(rbytes[RID_START:RID_END]),
		cyclic.NewIntFromBytes(rbytes[RMIC_START:RMIC_END]),
	}

}

// Creates a deep copy of the recipient, used for sending multiple messages
func (r Recipient) DeepCopy() Recipient {
	return Recipient{
		cyclic.NewInt(0).Set(r.GetRecipientInitVect()),
		cyclic.NewInt(0).Set(r.GetRecipientEmpty()),
		cyclic.NewInt(0).Set(r.GetRecipientID()),
		cyclic.NewInt(0).Set(r.GetRecipientMIC()),
	}
}
