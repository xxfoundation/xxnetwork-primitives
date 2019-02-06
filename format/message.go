////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"errors"
	"fmt"
	// TODO Should this dependency remain? Since we're mostly using cyclic ints
	// as a placeholder, it would surely make more sense to use byte arrays,
	// and each subrange would just be a slice into that array.
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/elixxir/primitives/userid"
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
	Payload   *cyclic.Int
	Recipient *cyclic.Int
}

// Structure which contains a message payload and the recipient payload in an
// easily accessible format
type Message struct {
	Payload
	Recipient
}

//Returns a serialized sender ID for the message interface
func (m Message) GetSender() *id.UserID {
	result := new(id.UserID).SetBytes(m.senderID.LeftpadBytes(SID_LEN))
	return result
}

//Returns the payload for the message interface
func (m Message) GetPayload() []byte {
	return m.data.Bytes()
}

//Returns a serialized recipient id for the message interface
// FIXME Two copies for this isn't great
func (m Message) GetRecipient() *id.UserID {
	result := new(id.UserID).SetBytes(m.recipientID.LeftpadBytes(RID_LEN))
	return result
}

// Makes a new message for a certain sender and recipient
func NewMessage(sender, recipient *id.UserID, text []byte) ([]Message, error) {

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
	return MessageSerial{m.Payload.SerializePayload(),
		m.Recipient.SerializeRecipient()}
}

func DeserializeMessage(ms MessageSerial) Message {
	return Message{DeserializePayload(ms.Payload),
		DeserializeRecipient(ms.Recipient)}
}
