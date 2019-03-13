////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

// Defines message structure.  Based on the "Basic Message Structure" doc
// Defining ranges in slices in go is inclusive for the beginning but
// exclusive for the end, so the END consts are one more then the final
// index.
const TOTAL_LEN int = 256

// Structure which contains a message payload and associated data in an
// easily accessible format
type Message struct {
	*Payload
	*AssociatedData
}

// Makes a new message
// TODO Should this allow population of any fields?
func NewMessage() *Message {
	return &Message{
		Payload:        NewPayload(),
		AssociatedData: NewAssociatedData(),
	}
}
