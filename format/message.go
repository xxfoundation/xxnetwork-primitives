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
	SizeLen        = 2
	KeyFPLen       = 32
	MacLen         = 32
	RecipientIDLen = 33
	TimestampLen   = 16

	MinimumPrimeSize = 2*MacLen + RecipientIDLen + TimestampLen

	AssociatedDataSize = KeyFPLen + SizeLen + MacLen + RecipientIDLen + TimestampLen
)

/*                               Message Structure (not to scale)
+-----------------------------------------------------------------------------------------------------------+
|                                                Message                                                    |
|                                            2*primeSize bits                                               |
+---------------------------------------------------+-------------------------------------------------------+
|                  payloadA                         |                 payloadB                              |
|               primeSize bits                      |              primeSize bits                           |
+---------+----------+---------+--------------------+---------+-------+-----------+-------------+-----------+
| grpBitA |  keyFP   |  size   |     Contents1      | grpBitB |  mac  | Contents2 | recipientID | timestamp |
|  1 bit  | 255 bits | 16 bits |      *below*       | 1 bit   | 255 b |  *below*  |  264 bits   |   128 b   |
+ --------+----------+---------+--------------------+---------+-------+-----------+-------------+-----------+

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
	size        []byte
	contents1   []byte
	mac         []byte
	contents2   []byte
	recipientID []byte
	timestamp   []byte
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
		size:      data[KeyFPLen : KeyFPLen+SizeLen],
		contents1: data[KeyFPLen+SizeLen:],

		mac:         data[numPrimeBytes : numPrimeBytes+MacLen],
		contents2:   data[numPrimeBytes+MacLen : 2*numPrimeBytes-RecipientIDLen-TimestampLen],
		recipientID: data[2*numPrimeBytes-RecipientIDLen-TimestampLen : 2*numPrimeBytes-TimestampLen],
		timestamp:   data[2*numPrimeBytes-TimestampLen:],
	}
}

// Returns a copy of the message
func (m Message) Copy() Message {
	m2 := NewMessage(len(m.data) / 2)
	copy(m2.data, m.data)
	return m2
}

// Returns the maximum size of the contents
func (m Message) ContentsSize() int {
	return len(m.data) - AssociatedDataSize
}

// Returns the underlying data buffer of the message
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

// Get contents returns the exact contents of the message. This size of the
// return is based on the size of the contents actually stored
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

	m.setContentsSize(uint16(len(c)))
}

// Returns the maximum size of the contents
func (m Message) GetSecretPayloadSize() int {
	return SizeLen + len(m.contents1) + len(m.contents2)
}

// Gets the entire payload which needs to be end to end encrypted for the
// communication to be secret. Contains the size, Contents1 and Contents2
func (m Message) GetSecretPayload() []byte {
	sp := make([]byte, SizeLen+len(m.contents1)+len(m.contents2))
	copy(sp[:SizeLen], m.size)
	copy(sp[SizeLen:SizeLen+len(m.contents1)], m.contents1)
	copy(sp[SizeLen+len(m.contents1):SizeLen+len(m.contents1)+len(m.contents2)], m.contents2)
	return sp
}

// Sets the entire payload which needs to be end to end encrypted for the
// communication to be secret. Sets the size, Contents1 and Contents2. Must be
// the exact size of the SecretPayload.
func (m Message) SetSecretPayload(sp []byte) {
	if len(sp) != m.GetSecretPayloadSize() {
		jww.ERROR.Panicf("secretPayload wrong size at %v bytes, must be "+
			"%v bytes", len(sp), m.GetSecretPayloadSize())
	}

	copy(m.size, sp[:SizeLen])
	copy(m.contents1, sp[SizeLen:SizeLen+len(m.contents1)])
	copy(m.contents2, sp[SizeLen+len(m.contents1):SizeLen+len(m.contents1)+len(m.contents2)])
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

// Gets the timestamp
func (m Message) GetTimestamp() []byte {
	return copyByteSlice(m.timestamp)
}

// Sets the timestamp. Panics if the passed data is not the exact right size
func (m Message) SetTimestamp(ts []byte) {
	if len(ts) != len(m.timestamp) {
		jww.ERROR.Panicf("timestamp not the correct size;"+
			"Expected: %v, Recieved: %v",
			len(m.timestamp), len(ts))
	}

	copy(m.timestamp, ts)
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

/*private functions*/
// Sets the size of the contents
func (m Message) setContentsSize(s uint16) {
	binary.BigEndian.PutUint16(m.size, s)
}

//gets the size of the contents
func (m Message) getContentsSize() uint16 {
	return binary.BigEndian.Uint16(m.size)
}

//helper function to copy a byte slice
func copyByteSlice(s []byte) []byte {
	c := make([]byte, len(s))
	copy(c, s)
	return c
}
