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
	"bytes"
)

// Defines message structure.  Based the "Basic Message Structure" doc
// Defining rangings in slices in go is inclusive for the beginning but
// exclusive for the end, so the END consts are one more then the final
// index.
const (
	TOTAL_LEN uint64 = 256

	//Byte used to ensure the highest bit of a serilization is zero
	ZEROER byte = 0x7F
)

//TODO: generate ranges programmatic

// Holds the payloads once they have been serialized
type MessageSerial struct {
	Payload   []byte
	Recipient []byte
}

// Structure which contains a message payload and the recipient payload in an
// easily accessible format
type Message struct {
	Payload
	Recipient
}

//Returns a serialized sender ID for the message interface
func (m Message) GetSender() *userid.UserID {
	result := new(userid.UserID).SetBytes(m.senderID[:])
	return result
}

//Returns the payload for the message interface
// Automatically trims trailing zeroes
func (m Message) GetPayload() []byte {
	return bytes.TrimRight(m.data, "\x00")
}

//Returns a serialized recipient id for the message interface
// FIXME Two copies for this isn't great
func (m Message) GetRecipient() *userid.UserID {
	result := new(userid.UserID).SetBytes(m.recipientID[:])
	return result
}

// Makes a new message for a certain sender and recipient
func NewMessage(sender, recipient *userid.UserID, text []byte) ([]Message, error) {

	//build the recipient payload
	recipientPayload, err := NewRecipientPayload(recipient)

	if err != nil {
		err = errors.New(fmt.Sprintf(
			"Unable to build message due to recipient error: %s",
			err.Error()))
		return nil, err
	}

	//Build the message Payloads
	messagePayload := NewPayload(sender, text)

	messageList := make([]Message, len(messagePayload))

	for indx, pld := range messagePayload {
		messageList[indx] = Message{pld, recipientPayload.DeepCopy()}
	}

	return messageList, nil
}

func (m Message) SerializeMessage() MessageSerial {
	return MessageSerial{m.Payload.serializePayload(),
		m.Recipient.serializeRecipient()}
}

func DeserializeMessage(ms MessageSerial) Message {
	return Message{deserializePayload(ms.Payload),
		deserializeRecipient(ms.Recipient)}
}
