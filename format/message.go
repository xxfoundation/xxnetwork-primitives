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

// Defines message structure.  Based the "Basic MaryPoppins Structure" doc
// Defining ranges in slices in go is inclusive for the beginning but
// exclusive for the end, so the END consts are one more then the final
// index.
const (
	TOTAL_LEN int = 256

	//Byte used to ensure the highest bit of a serialization is zero
	ZEROER byte = 0x7F
)

// Structure which contains a message payload and the recipient payload in an
// easily accessible format
type MaryPoppins struct {
	*Message
	*AssociatedData
}

// Wrap the sender ID in its type
func (m MaryPoppins) GetSender() *id.User {
	result := new(id.User).SetBytes(m.senderID[:])
	return result
}

// Get the payload from a message
func (m MaryPoppins) GetPayload() []byte {
	return m.data
}

// Wrap the recipient ID in its type
func (m MaryPoppins) GetRecipient() *id.User {
	result := new(id.User).SetBytes(m.recipientID[:])
	return result
}

// Makes a new message for a certain sender and recipient
func NewMessage(sender, recipient *id.User, text []byte) (*MaryPoppins, error) {

	//build the recipient payload
	recipientPayload, err := NewAssociatedData(recipient)

	if err != nil {
		err = errors.New(fmt.Sprintf(
			"Unable to build message due to recipient error: %s",
			err.Error()))
		return nil, err
	}

	//Build the message Payloads
	messagePayload, err := NewMessagePayload(sender, text)

	message := MaryPoppins{messagePayload, recipientPayload}

	return &message, err
}

func (m MaryPoppins) SerializeMessage() MessageSerial {
	return MessageSerial{m.Message.SerializePayload(),
		m.AssociatedData.SerializeRecipient()}
}

func DeserializeMessage(ms MessageSerial) MaryPoppins {
	return MaryPoppins{DeserializeMessagePayload(ms.MessagePayload),
		DeserializeAssociatedData(ms.RecipientPayload)}
}
