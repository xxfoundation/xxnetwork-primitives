////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package format

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/primitives/id"
)

const (
	KeyFPLen       = 32
	MacLen         = 32
	RecipientIDLen = 33

	MinimumPrimeSize = 2*MacLen + RecipientIDLen

	AssociatedDataSize = KeyFPLen + MacLen + RecipientIDLen
)

/*                               Message Structure (not to scale)
+--------------------------------------------------------------------------------------+
|                                       Message							   			   |
|                                  2*primeSize bits					    			   |
+------------------------------------------+-------------------------------------------+
|                  payloadA                |                 payloadB                  |
|               primeSize bits             |              primeSize bits               |
+---------+----------+---------------------+---------+-------+-----------+-------------+
| grpBitA |  keyFP   |      Contents1      | grpBitB |  mac  | Contents2 | recipientID |
|  1 bit  | 255 bits |       *below*       | 1 bit   | 255 b |  *below*  |  264 bits   |
+ --------+----------+---------------------+---------+-------+-----------+-------------+
|			                     Raw Contents	                     	 |
|                        2*primeSize - recipientID bits                  |
+------------------------------------------------------------------------+


size - size in bits of the data which is stored

Size Contents1 = primeSize - grpBitASize - keyFPSize - sizeSize
Size Contents2 = primeSize - grpBitBSize - macSize- recipientIDSize - timestampSize

the size of the data in the two contents fields is stored within the "size" field

/////Adherence to the group/////////////////////////////////////////////////////
The first bits of keyFingerprint and MAC are enforced to be 0, thus ensuring
PayloadA and PayloadB are within the group
*/

// Message structure stores all the data serially. Subsequent fields point to
// subsections of the serialised data.
type Message struct {
	data []byte

	//Note: These are mapped to locations in the data object
	payloadA []byte
	payloadB []byte

	keyFP       []byte
	contents1   []byte
	mac         []byte
	contents2   []byte
	recipientID []byte

	rawContents []byte
}

// NewMessage creates a new empty message based upon the size of the encryption
// primes. All subcomponents point to locations in the internal data buffer.
// Panics if the prime size to too small.
func NewMessage(numPrimeBytes int) Message {
	if numPrimeBytes < MinimumPrimeSize {
		jww.FATAL.Panicf("cannot make message based off of a prime of "+
			"size %v, minnimum size is %v to small prime", numPrimeBytes,
			MinimumPrimeSize)
	}

	data := make([]byte, 2*numPrimeBytes)

	return Message{
		data: data,

		payloadA: data[:numPrimeBytes],
		payloadB: data[numPrimeBytes:],

		keyFP:     data[0:KeyFPLen],
		contents1: data[KeyFPLen:numPrimeBytes],

		mac:         data[numPrimeBytes : numPrimeBytes+MacLen],
		contents2:   data[numPrimeBytes+MacLen : 2*numPrimeBytes-RecipientIDLen],
		recipientID: data[2*numPrimeBytes-RecipientIDLen : 2*numPrimeBytes],

		rawContents: data[:2*numPrimeBytes-RecipientIDLen],
	}
}

// Marshal marshals the message into a byte slice.
func (m *Message) Marshal() []byte {
	return copyByteSlice(m.data)
}

// Unmarshal unmarshalls a byte slice into a new Message.
func Unmarshal(b []byte) Message {
	m := NewMessage(len(b) / 2)
	copy(m.data, b)
	return m
}

// Returns a copy of the message
func (m Message) Copy() Message {
	m2 := NewMessage(len(m.data) / 2)
	copy(m2.data, m.data)
	return m2
}

//returns the size of the prime used
func (m Message) GetPrimeByteLen() int {
	return len(m.data) / 2
}

// Returns the maximum size of the contents
func (m Message) ContentsSize() int {
	return len(m.data) - AssociatedDataSize
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

// Get contents returns the exact contents of the message. This size of the
// return is based on the size of the contents actually stored
func (m Message) GetContents() []byte {

	c := make([]byte, len(m.contents1)+len(m.contents2))

	copy(c[:len(m.contents1)], m.contents1)
	copy(c[len(m.contents1):], m.contents2)

	return c
}

// sets the contents of the message. This will panic if the payload is greater
// than the maximum size. This overwrites any storage already in the message but
// will not clear bits beyond the size of the passed contents.
// If the passed contents is larger than the maximum contents size this will
// panic
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

// Get raw contents returns the exact contents of the message. This field
// crosses over the group barrier and the setter of this is responsible for
// ensuring the underlying payloads are within the group.
func (m Message) GetRawContents() []byte {
	return copyByteSlice(m.rawContents)
}

// Get size raw contents returns the exact contents of the message.
func (m Message) GetRawContentsSize() int {
	return len(m.rawContents)
}

// sets raw  contents of the message. This field crosses over the group barrier
// and the setter of this is responsible for ensuring the underlying payloads
// are within the group. This will panic if the payload is greater
// than the maximum size. This overwrites any storage already in the message.
// If the passed contents is larger than the maximum contents size this will
// panic
func (m Message) SetRawContents(c []byte) {
	if len(c) != len(m.rawContents) {
		jww.ERROR.Panicf("Set raw contents too large at %v bytes, "+
			"must be %v bytes or less", len(c),
			len(m.contents1)+len(m.contents2))
	}

	copy(m.rawContents, c)
}

// Gets the recipientID
func (m Message) GetRecipientID() *id.ID {
	rid := id.ID{}
	copy(rid[:], m.recipientID)
	return &rid
}

// Sets the recipientID
func (m Message) SetRecipientID(rid *id.ID) {
	copy(m.recipientID, rid[:])
}

// Gets the Key Fingerprint
func (m Message) GetKeyFP() Fingerprint {
	fp := Fingerprint{}
	copy(fp[:], m.keyFP)
	return fp
}

// Sets the Key Fingerprint. Checks that the first bit of the Key Fingerprint is
// 0, otherwise it panics.
func (m Message) SetKeyFP(fp Fingerprint) {
	if len(fp) != len(m.keyFP) {
		jww.ERROR.Panicf("key fingerprint not the correct size;"+
			"Expected: %v, Recieved: %v",
			len(m.keyFP), len(fp))
	}

	if fp[0]>>7 != 0 {
		jww.ERROR.Panicf("key fingerprint's first bit is not zero")
	}

	copy(m.keyFP, fp[:])
}

// Gets the MAC
func (m Message) GetMac() []byte {
	return copyByteSlice(m.mac)
}

// Sets the MAC. Checks that the first bit of the MAC is
// 0, otherwise it panics.
func (m Message) SetMac(mac []byte) {
	if len(mac) != len(m.mac) {
		jww.ERROR.Panicf("mac is not the correct size;"+
			"Expected: %v, Recieved: %v",
			len(m.mac), len(mac))
	}

	if mac[0]>>7 != 0 {
		jww.ERROR.Panicf("mac's first bit is not zero")
	}

	copy(m.mac, mac)
}

//helper function to copy a byte slice
func copyByteSlice(s []byte) []byte {
	c := make([]byte, len(s))
	copy(c, s)
	return c
}
