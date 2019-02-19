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
	// Length and Position of the Message Payload Initialization Vector
	MIV_LEN   uint64 = 9
	MIV_START uint64 = 0
	MIV_END   uint64 = MIV_LEN

	// Length and Position of message payload
	DATA_LEN   uint64 = TOTAL_LEN - SID_LEN - MIV_LEN - MMIC_LEN
	DATA_START uint64 = MIV_END
	DATA_END   uint64 = DATA_START + DATA_LEN

	SID_LEN   uint64 = userid.UserIDLen
	SID_START uint64 = DATA_END
	SID_END   uint64 = SID_START + SID_LEN

	// Length and Position of the Message Payload MIC
	MMIC_LEN   uint64 = 8
	MMIC_START uint64 = SID_END
	MMIC_END   uint64 = MMIC_START + MMIC_LEN
)

type MessagePayload struct {
	// This array holds all of the message data
	payloadSerial [TOTAL_LEN]byte
	// All other slices point to their respective parts of the array. So, the
	// message is always serialized and ready to go, and no copies are required
	messageInitVect []byte
	senderID        []byte
	data            []byte
	messageMIC      []byte
}

// Makes a new message for a certain sender.
// Only takes the first DATA_LEN bytes for the payload.
// Split into multiple messages elsewhere or use less space if the message is
// too long to fit.
// Will return an error if the message was too long to fit in one payload
// Make sure to populate the initialization vector and the MIC later
func NewMessagePayload(sender *userid.UserID, text []byte) (*MessagePayload, error) {
	result := MessagePayload{payloadSerial: [TOTAL_LEN]byte{}}
	result.data = result.payloadSerial[DATA_START:DATA_END]
	result.messageMIC = result.payloadSerial[MMIC_START:MMIC_END]
	result.senderID = result.payloadSerial[SID_START:SID_END]
	copy(result.senderID, sender.Bytes())
	result.messageInitVect = result.payloadSerial[MIV_START:MIV_END]

	copyLen := copy(result.data, text)
	var err error
	if copyLen != len(text) {
		err = errors.New("Couldn't fit text in one payload")
	}

	return &result, err
}

// Get the initialization vector's slice
// This allows reading and writing the correct section of memory, but
// doesn't allow changing the slice header in the structure itself
func (p *MessagePayload) GetMessageInitVect() []byte {
	return p.messageInitVect
}

// Get the sender ID's slice
// This allows reading and writing the correct section of memory, but
// doesn't allow changing the slice header in the structure itself
func (p *MessagePayload) GetSenderID() []byte {
	return p.senderID
}

// Wrap the sender ID in its type
func (p *MessagePayload) GetSender() *userid.UserID {
	result := new(userid.UserID).SetBytes(p.senderID[:])
	return result
}

// This function returns a pointer to the data payload
// This ensures that while the data can be edited, it cant be reallocated
func (p *MessagePayload) GetData() []byte {
	return p.data
}

// This function returns a pointer to the payload MIC
// This ensures that while the data can be edited, it cant be reallocated
func (p *MessagePayload) GetPayloadMIC() []byte {
	return p.messageMIC
}

// Returns the serialized message payload
// TODO Does it make sense to make this an internal method?
func (p *MessagePayload) SerializePayload() []byte {
	// It's actually unnecessary to ensure that the highest bit of the
	// serialized message is zero here if the initialization vector was
	// correctly generated, but just in case, we set the first bit to zero
	// to ensure that the payload fits in the cyclic group.
	p.payloadSerial[0] = p.payloadSerial[0] & ZEROER

	return p.payloadSerial[:]
}

// Slices a serialized payload in the correct spots
func DeserializeMessagePayload(pSerial []byte) *MessagePayload {
	var pBytes [TOTAL_LEN]byte
	copy(pBytes[:], pSerial)

	return &MessagePayload{
		pBytes,
		pBytes[MIV_START:MIV_END],
		pBytes[SID_START:SID_END],
		pBytes[DATA_START:DATA_END],
		pBytes[MMIC_START:MMIC_END],
	}
}

func (p *MessagePayload) DeepCopy() *MessagePayload {
	return DeserializeMessagePayload(p.payloadSerial[:])
}