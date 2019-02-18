////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"errors"
	"gitlab.com/elixxir/primitives/userid"
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
	// This array holds all of the message data
	payloadSerial [TOTAL_LEN]byte
	// All other slices point to their respective parts of the array. So, the
	// message is always serialized and ready to go, and no copies are required
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
	result := Payload{payloadSerial: [TOTAL_LEN]byte{}}
	result.data = result.payloadSerial[DATA_START:DATA_END]
	result.payloadMIC = result.payloadSerial[PMIC_START:PMIC_END]
	result.senderID = result.payloadSerial[SID_START:SID_END]
	copy(result.senderID, sender.Bytes())
	result.payloadInitVect = result.payloadSerial[PIV_START:PIV_END]

	copyLen := copy(result.data, text)
	var err error
	if copyLen != len(text) {
		err = errors.New("Couldn't fit text in one payload")
	}

	return &result, err
}

// This function returns a pointer to the Payload Initialization Vector
// This ensures that while the data can be edited, it cant be reallocated
func (p *Payload) GetPayloadInitVect() []byte {
	return p.payloadInitVect
}

// This function returns a pointer to the Sender ID
// This ensures that while the data can be edited, it cant be reallocated
func (p *Payload) GetSenderID() []byte {
	return p.senderID
}

func (p *Payload) GetSender() *userid.UserID {
	result := new(userid.UserID).SetBytes(p.senderID[:])
	return result
}

// This function returns a pointer to the data payload
// This ensures that while the data can be edited, it cant be reallocated
func (p *Payload) GetData() []byte {
	return p.data
}

// This function returns a pointer to the payload MIC
// This ensures that while the data can be edited, it cant be reallocated
func (p *Payload) GetPayloadMIC() []byte {
	return p.payloadMIC
}

// Returns the serialized message payload
// TODO Does it make sense to make this an internal method?
func (p *Payload) SerializePayload() []byte {
	// It's actually unnecessary to ensure that the highest bit of the
	// serialized message is zero here if the initialization vector was
	// correctly generated, but just in case, we set the first bit to zero
	// to ensure that the payload fits in the cyclic group.
	p.payloadSerial[0] = p.payloadSerial[0] & ZEROER

	return p.payloadSerial[:]
}

// Slices a serialized payload in the correct spots
func DeserializePayload(pSerial []byte) *Payload {
	var pBytes [TOTAL_LEN]byte
	copy(pBytes[:], pSerial)

	return &Payload{
		pBytes,
		pBytes[PIV_START:PIV_END],
		pBytes[SID_START:SID_END],
		pBytes[DATA_START:DATA_END],
		pBytes[PMIC_START:PMIC_END],
	}
}

func (p *Payload) DeepCopy() *Payload {
	return DeserializePayload(p.payloadSerial[:])
}