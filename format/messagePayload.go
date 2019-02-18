////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"gitlab.com/elixxir/primitives/userid"
	"errors"
)

const (
	// Length and Position of the Payload Initialization Vector
	PIV_LEN   uint64 = 9
	PIV_START uint64 = 0
	PIV_END   uint64 = PIV_LEN

	// Length and Position of message payload
	DATA_LEN   uint64 = TOTAL_LEN - SID_LEN - PIV_LEN - PMIC_LEN
	DATA_START uint64 = PIV_END
	DATA_END   uint64 = DATA_START + DATA_LEN

	SID_LEN   uint64 = userid.UserIDLen
	SID_START uint64 = DATA_END
	SID_END   uint64 = SID_START + SID_LEN

	// Length and Position of the Payload MIC
	PMIC_LEN   uint64 = 8
	PMIC_START uint64 = SID_END
	PMIC_END   uint64 = PMIC_START + PMIC_LEN
)

type Payload struct {
	payloadInitVect []byte
	senderID        []byte
	data            []byte
	payloadMIC      []byte
}

// Makes a new message for a certain sender.
// Only takes the first DATA_LEN bytes for the payload.
// Split into multiple messages elsewhere or use less space if the message is
// too long to fit.
// Will return an error if the message was too long to fit in one payload
// Make sure to populate the initialization vector and the MIC later
func NewPayload(sender *userid.UserID, text []byte) (*Payload, error) {
	var datum [DATA_LEN]byte
	copyLen := copy(datum[:], text)
	var err error
	if copyLen != len(text) {
		err = errors.New("Couldn't fit text in one payload")
	}

	return &Payload{
		make([]byte, PIV_LEN),
		sender[:],
		datum[:],
		make([]byte, PMIC_LEN)}, err
}

// This function returns a pointer to the Payload Initialization Vector
// This ensures that while the data can be edited, it cant be reallocated
func (p Payload) GetPayloadInitVect() []byte {
	return p.payloadInitVect
}

// This function returns a pointer to the Sender ID
// This ensures that while the data can be edited, it cant be reallocated
func (p Payload) GetSenderID() []byte {
	return p.senderID
}

func (p Payload) GetSender() *userid.UserID {
	result := new(userid.UserID).SetBytes(p.senderID[:])
	return result
}

// This function returns a pointer to the data payload
// This ensures that while the data can be edited, it cant be reallocated
func (p Payload) GetData() []byte {
	return p.data
}

// This function returns a pointer to the payload MIC
// This ensures that while the data can be edited, it cant be reallocated
func (p Payload) GetPayloadMIC() []byte {
	return p.payloadMIC
}

// Returns the serialized message payload
// TODO Does it make sense to make this an internal method?
func (p Payload) SerializePayload() []byte {
	pbytes := make([]byte, TOTAL_LEN)

	// Copy the Payload Initialization Vector into the serialization
	copy(pbytes[PIV_START:PIV_END], p.payloadInitVect[:])

	// Copy the Sender ID into the serialization
	copy(pbytes[SID_START:SID_END], p.senderID[:])

	// Copy the payload data into the serialization
	copy(pbytes[DATA_START:DATA_END], p.data[:])

	// Copy the payload MIC into the serialization
	copy(pbytes[PMIC_START:PMIC_END], p.payloadMIC[:])

	//Make sure the highest bit of the serialization is zero
	pbytes[0] = pbytes[0] & ZEROER

	return pbytes
}

//Returns a Deserialized Message Payload
func DeserializePayload(pSerial []byte) Payload {

	return Payload{
		pSerial[PIV_START:PIV_END],
		pSerial[SID_START:SID_END],
		pSerial[DATA_START:DATA_END],
		pSerial[PMIC_START:PMIC_END],
	}
}
