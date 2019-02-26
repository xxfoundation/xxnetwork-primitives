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
	// Length and position of sender ID
	// Because the first bit of all user IDs is zero, this will always be in the
	// cyclic group
	MP_SID_LEN   int = id.UserLen
	MP_SID_START int = 0
	MP_SID_END   int = MP_SID_START + MP_SID_LEN

	// Length and Position of message payload
	// Includes both padding and payload
	MP_PAYLOAD_LEN   int = TOTAL_LEN - MP_SID_LEN
	MP_PAYLOAD_START int = MP_SID_END
	MP_PAYLOAD_END   int = MP_PAYLOAD_START + MP_PAYLOAD_LEN
)

type Payload struct {
	// This array holds all of the message payload
	payloadSerial [TOTAL_LEN]byte
	// All other slices point to their respective parts of the array. So, the
	// message is always serialized and ready to go, and no copies are required
	payload  []byte
	senderID []byte
}

// Makes a new message for a certain sender.
// Only takes the first DATA_LEN bytes for the payload.
// Split into multiple messages elsewhere or use less space if the message is
// too long to fit.
// Will return an error if the message was too long to fit in one payload
// Make sure to populate the initialization vector and the MIC later
func NewPayload() *Payload {
	result := Payload{payloadSerial: [TOTAL_LEN]byte{}}
	result.payload = result.payloadSerial[MP_PAYLOAD_START:MP_PAYLOAD_END]
	result.senderID = result.payloadSerial[MP_SID_START:MP_SID_END]

	return &result
}

// Get the sender ID's slice
// This allows reading and writing the correct section of memory, but
// doesn't allow changing the slice header in the structure itself
func (p *Payload) GetSenderID() []byte {
	return p.senderID
}

// Wrap the sender ID in its type
func (p *Payload) GetSender() *id.User {
	result := new(id.User).SetBytes(p.senderID[:])
	return result
}

// Returns number of bytes copied
func (p *Payload) SetSenderID(newId []byte) int {
	return copy(p.senderID, newId)
}

func (p *Payload) SetSender(newId *id.User) {
	copy(p.senderID, newId.Bytes())
}

// This function returns a pointer to the payload payload
// This ensures that while the data can be edited, it cant be reallocated
func (p *Payload) GetPayload() []byte {
	return p.payload
}

// Returns number of bytes copied
func (p *Payload) SetPayload(payload []byte) int {
	return copy(p.payload, payload)
}

// Returns the serialized message payload
func (p *Payload) SerializePayload() []byte {
	return p.payloadSerial[:]
}

// Slices a serialized payload in the correct spots
func DeserializePayload(pSerial []byte) *Payload {
	result := NewPayload()
	copy(result.payloadSerial[:], pSerial)
	return result
}

func (p *Payload) DeepCopy() *Payload {
	return DeserializePayload(p.payloadSerial[:])
}
