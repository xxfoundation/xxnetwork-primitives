////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"errors"
	"fmt"
	"gitlab.com/elixxir/primitives/id"
)

// Defines message structure.  Based the "Basic Message Structure" doc
// Defining ranges in slices in go is inclusive for the beginning but
// exclusive for the end, so the END consts are one more then the final
// index.
const (
	TOTAL_LEN uint64 = 256

	//Byte used to ensure the highest bit of a serialization is zero
	ZEROER byte = 0x7F
)

//TODO: generate ranges programmatic

// Holds the payloads once they have been serialized
type MessageSerial struct {
	MessagePayload   []byte
	RecipientPayload []byte
}

// Structure which contains a message payload and the recipient payload in an
// easily accessible format
type Message struct {
	*MessagePayload
	*RecipientPayload
}

// Wrap the sender ID in its type
func (m Message) GetSender() *id.User {
	result := new(id.User).SetBytes(m.senderID[:])
	return result
}

// Get the payload from a message
func (m Message) GetPayload() []byte {
	return m.data
}

// Wrap the recipient ID in its type
func (m Message) GetRecipient() *id.User {
	result := new(id.User).SetBytes(m.recipientID[:])
	return result
}

// Makes a new message for a certain sender and recipient
func NewMessage(sender, recipient *id.User, text []byte) (*Message, error) {

	//build the recipient payload
	recipientPayload, err := NewRecipientPayload(recipient)

	if err != nil {
		err = errors.New(fmt.Sprintf(
			"Unable to build message due to recipient error: %s",
			err.Error()))
		return nil, err
	}

	//Build the message Payloads
	messagePayload, err := NewMessagePayload(sender, text)

	message := Message{messagePayload, recipientPayload}

	return &message, err
}

func (m Message) SerializeMessage() MessageSerial {
	return MessageSerial{m.MessagePayload.SerializePayload(),
		m.RecipientPayload.SerializeRecipient()}
}

func DeserializeMessage(ms MessageSerial) Message {
	return Message{DeserializeMessagePayload(ms.MessagePayload),
		DeserializeRecipient(ms.RecipientPayload)}
}
