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
	TOTAL_LEN        = 256
	MP_SID_LEN   int = id.UserLen
	MP_SID_START int = 0
	MP_SID_END   int = MP_SID_START + MP_SID_LEN

	// Length and Position of message payload
	// Includes both padding and payload
	MP_PAYLOAD_LEN   int = TOTAL_LEN - MP_SID_LEN
	MP_PAYLOAD_START int = MP_SID_END
	MP_PAYLOAD_END   int = MP_PAYLOAD_START + MP_PAYLOAD_LEN
)

// The payload data must be variable, because for E2E, padding will
// be added, which has a minimum length of 11 bytes. This means that
// the E2E encrypt function can't accept a byte slice of more than
// 256-11 bytes. The senderID and payload data can be set on payload
// struct as normal, and then GetPayload can be used to get a byte
// slice containing only the actual data, instead of TOTAL_LEN bytes.
// SetPayload can be used to set the encrypted data into the object,
// since after padding the size will always be TOTAL_LEN.
// On decryption, the function SetSplitPayload can be used to split
// a serialized byte slice into the senderID and payload data
type Payload struct {
	// This array holds all of the message payload
	payloadSerial [TOTAL_LEN]byte
	// All other slices point to their respective parts of the array. So, the
	// message is always serialized and ready to go, and no copies are required
	senderID    []byte
	payloadData []byte
	// Actual size of payloadData before encryption / after decryption
	size int
}

// Makes a new message for a certain sender.
// Only takes the first DATA_LEN bytes for the payload.
// Split into multiple messages elsewhere or use less space if the message is
// too long to fit.
// Will return an error if the message was too long to fit in one payload
// Make sure to populate the initialization vector and the MIC later
func NewPayload() *Payload {
	result := Payload{payloadSerial: [TOTAL_LEN]byte{}, size: 0}
	result.payloadData = result.payloadSerial[MP_PAYLOAD_START:MP_PAYLOAD_END]
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
	result := id.NewUserFromBytes(p.senderID[:])
	return result
}

// Returns number of bytes copied
func (p *Payload) SetSenderID(newId []byte) int {
	return copy(p.senderID, newId)
}

func (p *Payload) SetSender(newId *id.User) {
	copy(p.senderID, newId.Bytes())
}

// This function returns a pointer to the payload data
// This ensures that while the data can be edited, it cant be reallocated
func (p *Payload) GetPayloadData() []byte {
	return p.payloadData[:p.size]
}

// Returns number of bytes copied
func (p *Payload) SetPayloadData(payload []byte) int {
	p.size = copy(p.payloadData, payload)
	return p.size
}

// Returns the actual existing serialized payload
func (p *Payload) GetPayload() []byte {
	return p.payloadSerial[:MP_SID_LEN+p.size]
}

// Set serialized payload
// Returns number of bytes copied
func (p *Payload) SetPayload(pSerial []byte) int {
	return copy(p.payloadSerial[:], pSerial)
}

// Split decrypted payload into actual fields
// Assume first 32 bytes are always senderID
// and rest is payload data
// Returns number of bytes copied
func (p *Payload) SetSplitPayload(pSerial []byte) int {
	ret1 := p.SetSenderID(pSerial[:MP_SID_LEN])
	ret2 := p.SetPayloadData(pSerial[MP_SID_LEN:])
	return ret1 + ret2
}

// Returns the full serialized payload
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
