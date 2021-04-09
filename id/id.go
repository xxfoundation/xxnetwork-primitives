////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

// Contains the generic ID type, which is a byte array that represents an entity
// ID. The first bytes in the array contain the actual ID data while the last
// byte contains the ID type, which is either generic, gateway, node, or user.
// IDs can be hard coded or generated using a cryptographic function found in
// crypto.
package id

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
	"testing"
)

// ID data sizes
const (
	dataLen  = 32          // Length of the ID data in bytes
	ArrIDLen = dataLen + 1 // Length of the full ID array in bytes
)

// Error strings
const (
	unmarshalLenErr   = "ID Unmarshal: length of data %d != %d expected"
	readerErr         = "NewRandomID: failed to generate random bytes: %+v"
	fromBytesTestErr  = "NewIdFromBytes() can only be used for testing."
	fromStringTestErr = "NewIdFromString() can only be used for testing."
	fromUintTestErr   = "NewIdFromUInt() can only be used for testing."
	fromUintsTestErr  = "NewIdFromUInts() can only be used for testing."
)

// ID is a fixed-length array containing data that services as an identifier for
// entities. The first 32 bytes hold the ID data while the last byte holds the
// type, which describes the type of entity the ID belongs to.
type ID [ArrIDLen]byte

// Marshal returns the ID bytes in wire format.
func (id ID) Marshal() []byte {
	return id.Bytes()
}

// Unmarshal unmarshalls the ID wire format into an ID object.
func Unmarshal(buff []byte) (ID, error) {
	// Return an error if the provided data is not the correct length
	if len(buff) != ArrIDLen {
		return ID{}, errors.Errorf(unmarshalLenErr, len(buff), ArrIDLen)
	}

	return makeID(buff), nil
}

// Bytes returns a copy of an ID as a byte slice. Note that Bytes() is used by
// Marshal() and any changes made here will affect how Marshal() functions.
func (id ID) Bytes() []byte {
	return id[:]
}

// String converts an ID to a string via base64 encoding.
func (id ID) String() string {
	return base64.StdEncoding.EncodeToString(id.Bytes())
}

// GetType returns the ID's type (last byte of the array).
func (id ID) GetType() Type {
	return Type(id[ArrIDLen-1])
}

// SetType returns a copy of the ID with the specified ID type set.
func (id ID) SetType(t Type) ID {
	newID := id
	newID[ArrIDLen-1] = byte(t)
	return newID
}

// NewRandomID generates a random ID with the specified type.
func NewRandomID(r io.Reader, t Type) (ID, error) {
	// Generate random bytes
	idBytes := make([]byte, ArrIDLen)
	if _, err := r.Read(idBytes); err != nil {
		return ID{}, errors.Errorf(readerErr, err)
	}

	// Create ID from bytes and set the type
	id := makeID(idBytes).SetType(t)

	return id, nil
}

// NewIdFromBytes creates a new ID from the supplied byte slice. It is similar
// to Unmarshal() but does not do any error checking. If the data is longer than
// ArrIDLen, then it is truncated. If it is shorter, then the remaining bytes
// are filled with zeroes. This function is for testing purposes only.
func NewIdFromBytes(buff []byte, x interface{}) ID {
	// Ensure that this function is only run in testing environments
	switch x.(type) {
	case *testing.T, *testing.M, *testing.B, *testing.PB:
	default:
		panic(fromBytesTestErr)
	}

	return makeID(buff)
}

// NewIdFromString creates a new ID from the given string and type. If the
// string is longer than dataLen, then it is truncated. If it is shorter, then
// the remaining bytes are filled with zeroes. This function is for testing
// purposes only.
func NewIdFromString(idString string, t Type, x interface{}) ID {
	// Ensure that this function is only run in testing environments
	switch x.(type) {
	case *testing.T, *testing.M, *testing.B, *testing.PB:
		break
	default:
		panic(fromStringTestErr)
	}

	// Create a new ID from the string and set the type
	return NewIdFromBytes([]byte(idString), x).SetType(t)
}

// NewIdFromUInt converts the specified uint64 into bytes and returns a new ID
// based off it with the specified ID type. The remaining space of the array is
// filled with zeros. This function is for testing purposes only.
func NewIdFromUInt(idUInt uint64, t Type, x interface{}) ID {
	// Ensure that this function is only run in testing environments
	switch x.(type) {
	case *testing.T, *testing.M, *testing.B:
		break
	default:
		panic(fromUintTestErr)
	}
	// Convert the uints to bytes
	var id ID
	binary.BigEndian.PutUint64(id[:], idUInt)

	// Set ID type and return
	return id.SetType(t)
}

// NewIdFromUInt converts the specified uint64 array into bytes and returns a
// new ID based off it with the specified ID type. Unlike NewIdFromUInt(), the
// four uint64s provided fill the entire ID array. This function is for testing
// purposes only.
func NewIdFromUInts(idUInts [4]uint64, t Type, x interface{}) ID {
	// Ensure that this function is only run in testing environments
	switch x.(type) {
	case *testing.T, *testing.M, *testing.B:
		break
	default:
		panic(fromUintsTestErr)
	}

	// Convert the uints to bytes
	var id ID
	for i, idUint := range idUInts {
		binary.BigEndian.PutUint64(id[i*8:], idUint)
	}

	// Set ID type and return
	return id.SetType(t)
}

// makeID creates a new ID from the buffer.
func makeID(buff []byte) ID {
	var id ID
	copy(id[:], buff)
	return id
}
