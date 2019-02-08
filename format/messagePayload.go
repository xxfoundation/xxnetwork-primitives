////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
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
	payloadInitVect []byte
	senderID        []byte
	data            []byte
	payloadMIC      []byte
}

// Makes a new message for a certain sender.
// Splits the message into multiple if it is too long
// TODO Doing message splitting here didn't end up meeting the needs of the
// client. Maybe we should remove it from here to simplify things.
func NewPayload(sender *userid.UserID, text []byte) []Payload {
	// Split the payload into multiple sub-payloads if it is longer than the
	// maximum allowed
	var dataLst [][]byte

	for uint64(len(text)) > DATA_LEN {
		dataLst = append(dataLst, text[0:DATA_LEN])
		text = text[DATA_LEN:]
	}

	dataLst = append(dataLst, text)

	//Create a message payload for every sub-payload
	var payloadLst []Payload

	for _, datum := range dataLst {
		payload := Payload{
			make([]byte, PIV_LEN),
			sender[:],
			datum[:],
			make([]byte, PMIC_LEN)}
		payloadLst = append(payloadLst, payload)
	}

	return payloadLst
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
func (p Payload) serializePayload() []byte {
	pbytes := make([]byte, TOTAL_LEN)

	// Copy the Payload Initialization Vector into the serialization
	copy(pbytes[PIV_START:PIV_END], p.payloadInitVect[:])

	// Copy the Sender ID into the serialization
	copy(pbytes[SID_START:SID_END], p.senderID[:])

	// Copy the payload data into the serialization
	copy(pbytes[DATA_START:DATA_END], p.data[:])

	// Copy the payloac MIC into the serialization
	copy(pbytes[PMIC_START:PMIC_END], p.payloadMIC[:])

	//Make sure the highest bit of the serialization is zero
	pbytes[0] = pbytes[0] & ZEROER

	return pbytes
}

//Returns a Deserialized Message Payload
func deserializePayload(pSerial []byte) Payload {

	return Payload{
		pSerial[PIV_START:PIV_END],
		pSerial[SID_START:SID_END],
		pSerial[DATA_START:DATA_END],
		pSerial[RMIC_START:RMIC_END],
	}
}
