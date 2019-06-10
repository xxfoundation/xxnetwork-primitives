////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

const (
	// Length of the entire message serial
	TotalLen = 512 // 4096 bits

	// Length, start index, and end index of the payloads
	subPayloadLen = 256 // 2048 bits
	payloadAStart = 0
	payloadAEnd   = payloadAStart + subPayloadLen
	payloadBStart = payloadAEnd
	payloadBEnd   = payloadBStart + subPayloadLen
)

// Structure for the message stores all the data serially. Subsequent fields
// point to subsections of the serialised data so that the message is always
// serialized, is ready to go, and no copies are required.
type Message struct {
	master         [TotalLen]byte // serialised message data
	Contents                      // points to the contents of the message
	AssociatedData                // points to the associate data of the message
	payloadA       []byte         // points to the first half of the message
	payloadB       []byte         // points to the second half of the message
}

/*
+----------------------------------------------------------------------------------------+
|                                         Message                                        |
|                                        4096 bits                                       |
+----------------------------------------------------------------------------------------+
|                  payloadA                  |                 payloadB                  |
|                 2048 bits                  |                2048 bits                  |
+------------------------------------+-------+-------------------------------------------+
|              Contents              |                   AssociatedData                  |
|              3192 bits             |                      904 bits                     |
+------------------------------------+---------------------------------------------------+
|     padding     |       data       | recipientID | keyFP | timestamp |  mac  | grpByte |
|   88–3192 bits  |    0–3104 bits   |    256 b    | 256 b |   128 b   | 256 b |   8 b   |
+-----------------+------------------+-------------+-------+-----------+-------+---------+
*/

// NewMessage creates a new empty message. It points the contents, associated
// data, payload A, and payload B, to their respective parts of master.
func NewMessage() *Message {
	newMsg := &Message{master: [TotalLen]byte{}}

	newMsg.Contents.serial = newMsg.master[contentsStart:contentsEnd]
	newMsg.AssociatedData.serial = newMsg.master[associatedDataStart:associatedDataEnd]
	newMsg.payloadA = newMsg.master[payloadAStart:payloadAEnd]
	newMsg.payloadB = newMsg.master[payloadBStart:payloadBEnd]

	newMsg.Contents = *NewContents(newMsg.Contents.serial)
	newMsg.AssociatedData = *NewAssociatedData(newMsg.AssociatedData.serial)

	return newMsg
}

// GetMaster returns the entire serialised message.
func (m *Message) GetMaster() []byte {
	return m.master[:]
}

// GetPayloadA returns payload A, which is the first half of the message.
func (m *Message) GetPayloadA() []byte {
	return m.payloadA
}

// SetPayloadA copies the passed byte slice into payloadA. The number of bytes
// copied is returned. If the specified byte array is not exactly the same size
// as payloadA, then it panics.
func (m *Message) SetPayloadA(payload []byte) int {
	if len(payload) == subPayloadLen {
		return copy(m.payloadA, payload)
	} else {
		panic("new payload not the same size as PayloadA")
	}
}

// GetPayloadB returns payload B, which is the last half of the message.
func (m *Message) GetPayloadB() []byte {
	return m.payloadB
}

// SetPayloadB copies the passed byte slice into payloadB. The number of bytes
// copied is returned. If the specified byte array is not exactly the same size
// as payloadB, then it panics.
func (m *Message) SetPayloadB(payload []byte) int {
	if len(payload) == subPayloadLen {
		return copy(m.payloadB, payload)
	} else {
		panic("new payload not the same size as PayloadB")
	}
}

// GetPayloadBForEncryption ensures payload B is in the group for encrypting.
// Specifically, it moves the first byte to the end and sets the first byte to
// zero.
func (m *Message) GetPayloadBForEncryption() []byte {
	payloadCopy := make([]byte, subPayloadLen)
	copy(payloadCopy, m.payloadB)
	payloadCopy[subPayloadLen-1] = payloadCopy[0]
	payloadCopy[0] = 0

	return payloadCopy
}

// SetDecryptedPayloadB is used when receiving a decrypted payload B to ensure
// all data is put back in the right order. Specifically, it moves the last byte
// to the front and sets the last byte to zero. The number of bytes copied is
// returned.
func (m *Message) SetDecryptedPayloadB(newPayload []byte) int {
	size := copy(m.payloadB, newPayload)
	m.payloadB[0] = m.payloadB[subPayloadLen-1]
	m.payloadB[subPayloadLen-1] = 0

	return size
}
