////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package format

import (
	"encoding/binary"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/primitives/id"
)

const (
	GrpByteLen     = 1
	RecipientIDLen = 33
	KeyFPLen       = 32
	TimestampLen   = 16
	MacLen         = 32

	// Length of the entire message serial
	AssociatedDataLen = RecipientIDLen + KeyFPLen + TimestampLen + MacLen
)

/*                               Message Structure (not to scale)
+----------------------------------------------------------------------------------------+
|                                         Message                                        |
|                                        4096 bits                                       |
+----------------------------------------------------------------------------------------+
|                  payloadA                  |                 payloadB                  |
|                 2048 bits                  |                2048 bits                  |
+------------------------------------+-------+---------------------------------+---------+
|              Contents              |             AssociatedData              | grpByte |
|              3192 bits             |                896 bits                 | 8 bits  |
+------------------------------------+-----------------------------------------+         |
|     padding     |       data       | recipientID | keyFP | timestamp |  mac  |         |
|   88–3192 bits  |    0–3104 bits   |   264 bits  | 256 b |  128 bits | 256 b |         |
+-----------------+------------------+-------------+-------+-----------+-------+---------+
*/

// Message structure stores all the data serially. Subsequent fields point to
// subsections of the serialised data.
type Message struct {
	data []byte

	//Note: These are mapped to locations in the data object
	payloadA []byte
	payloadB []byte

	groupByteA  []byte
	contents1   []byte
	groupByteB  []byte
	contents2   []byte
	recipientID []byte
	keyFP       []byte
	timestamp   []byte
	mac         []byte

	associatedData []byte
}

// NewMessage creates a new empty message. It points the contents, associated
// data, payload A, and payload B, to their respective parts of master.
func NewMessage(numPrimeBytes int) Message {

	if numPrimeBytes < 2*(AssociatedDataLen+GrpByteLen) {
		panic("cannot make message based off of to small prime")
	}

	data := make([]byte, 2*numPrimeBytes)

	adStart := 2*numPrimeBytes - AssociatedDataLen - GrpByteLen

	return Message{
		data: data,

		payloadA: data[:numPrimeBytes],
		payloadB: data[numPrimeBytes:],

		groupByteA:  data[0:1],
		contents1:   data[1:numPrimeBytes],
		groupByteB:  data[numPrimeBytes : numPrimeBytes+1],
		contents2:   data[numPrimeBytes+1 : adStart],
		recipientID: data[adStart : adStart+RecipientIDLen],
		keyFP:       data[adStart+RecipientIDLen : adStart+RecipientIDLen+KeyFPLen],
		timestamp:   data[adStart+RecipientIDLen+KeyFPLen : adStart+RecipientIDLen+KeyFPLen+MacLen],

		associatedData: data[adStart : adStart+RecipientIDLen+KeyFPLen+MacLen],
	}
}

// GetMaster returns the entire serialised message.
func (m Message) GetData() []byte {
	return copyByteSlice(m.data)
}

// GetPayloadA returns payload A, which is the first half of the message.
func (m Message) GetPayloadA() []byte {
	return copyByteSlice(m.payloadA)
}

// SetPayloadA copies the passed byte slice into payload A. If the specified
// byte slice is not exactly the same size as payload A, then it panics.
func (m Message) SetPayloadA(payload []byte) {
	if len(payload) != len(m.payloadA) {
		jww.ERROR.Panicf("new payload not the same size as PayloadA;"+
			"Expected: %v, Recieved: %v",
			len(m.payloadA), len(payload))
	}

	copy(m.payloadA, payload)
}

// GetPayloadB returns payload B, which is the last half of the message.
func (m Message) GetPayloadB() []byte {
	return copyByteSlice(m.payloadB)
}

// SetPayloadB copies the passed byte slice into payload B. If the specified
// byte slice is not exactly the same size as payload B, then it panics.
func (m Message) SetPayloadB(payload []byte) {
	if len(payload) != len(m.payloadB) {
		jww.ERROR.Panicf("new payload not the same size as PayloadB;"+
			"Expected: %v, Recieved: %v",
			len(m.payloadB), len(payload))
	}

	copy(m.payloadB, payload)
}

func (m Message) GetContents() []byte {
	size := int(m.getContentsSize())
	c := make([]byte, size)

	if size <= len(m.contents1) {
		copy(c, m.contents1[:size])
	} else {
		copy(c[:len(m.contents1)], m.contents1)
		copy(c[len(m.contents1):size], m.contents2[:size-len(m.contents1)])
	}

	return c
}

func (m Message) SetContents(c []byte) {
	if len(c) > len(m.contents1)+len(m.contents2) {
		jww.ERROR.Panicf("contents too large at %v bytes, must be "+
			"%v bytes or less", len(c), len(m.contents1)+len(m.contents2))
	}

	if len(c) <= len(m.contents1) {
		copy(m.contents1, c)
	} else {
		copy(m.contents1, c[:len(m.contents1)])
		copy(m.contents2, c[len(m.contents1):])
	}
}

func (m Message) GetRecipientID() *id.ID {
	rid := id.ID{}
	copy(rid[:], m.recipientID)
	return &rid
}

func (m Message) SetRecipientID(rid *id.ID) {
	copy(m.recipientID, rid[:])
}

func (m Message) GetKeyFP() []byte {
	return copyByteSlice(m.keyFP)
}

func (m Message) SetKeyFP(fp []byte) {
	if len(fp) != len(m.keyFP) {
		jww.ERROR.Panicf("key fingerprint not the correct size;"+
			"Expected: %v, Recieved: %v",
			len(m.keyFP), len(fp))
	}

	copy(m.keyFP, fp)
}

func (m Message) GetTimestamp() []byte {
	return copyByteSlice(m.timestamp)
}

func (m Message) SetTimestamp(ts []byte) {
	if len(ts) != len(m.timestamp) {
		jww.ERROR.Panicf("timestamp not the correct size;"+
			"Expected: %v, Recieved: %v",
			len(m.timestamp), len(ts))
	}

	copy(m.timestamp, ts)
}

func (m Message) GetMac() []byte {
	return copyByteSlice(m.mac)
}

func (m Message) SetMac(mac []byte) {
	if len(mac) != len(m.mac) {
		jww.ERROR.Panicf("timestamp not the correct size;"+
			"Expected: %v, Recieved: %v",
			len(m.mac), len(mac))
	}

	copy(m.mac, mac)
}

/*private functions*/
func (m Message) setContentsSize(s uint16) {
	sEnc := make([]byte, 2)
	binary.BigEndian.PutUint16(sEnc, s)

	m.groupByteB[0] = sEnc[1] & 0x7f
	m.groupByteA[0] = (sEnc[0]&0x3f)<<1 | (sEnc[1]&0x80)>>7
}

func (m Message) getContentsSize() uint16 {
	sEnc := make([]byte, 2)

	sEnc[1] = m.groupByteB[0] | (m.groupByteA[0]&0x1)<<7
	sEnc[0] = (m.groupByteA[0] & 0x7E) >> 1

	return binary.BigEndian.Uint16(sEnc)

}

func copyByteSlice(s []byte) []byte {
	c := make([]byte, len(s))
	copy(c, s)
	return c
}
