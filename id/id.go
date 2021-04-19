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
	jww "github.com/spf13/jwalterweatherman"
	"io"
	"testing"
)

const (
	// Length of the ID data
	dataLen = 32

	// Length of the full ID array
	ArrIDLen = dataLen + 1
)

// ID is a fixed-length array containing data that services as an identifier for
// entities. The first 32 bytes hold the ID data while the last byte holds the
// type, which describes the type of entity the ID belongs to.
type ID [ArrIDLen]byte

// Marshal returns the ID bytes in wire format.
func (id *ID) Marshal() []byte {
	return id.Bytes()
}

// Unmarshal unmarshalls the ID wire format into an ID object.
func Unmarshal(data []byte) (*ID, error) {
	// Return an error if the provided data is not the correct length
	if len(data) != ArrIDLen {
		return nil, errors.Errorf("Failed to unmarshal ID: length of data "+
			"must be %d, length received is %d", ArrIDLen, len(data))
	}

	return copyID(data), nil
}

// Bytes returns a copy of an ID as a byte slice. Note that Bytes() is used by
// Marshal() and any changes made here will affect how Marshal() functions.
func (id *ID) Bytes() []byte {
	if id == nil {
		jww.FATAL.Panicf("Failed to get bytes of ID: ID is nil.")
	}

	newBytes := make([]byte, ArrIDLen)
	copy(newBytes, id[:])

	return newBytes
}

// Cmp determines whether two IDs are equal. Returns true if they are equal and
// false otherwise.
func (id *ID) Cmp(y *ID) bool {
	if id == nil || y == nil {
		jww.FATAL.Panicf("Failed to compare IDs: one or both IDs are nil.")
	}

	return *id == *y
}

// DeepCopy creates a copy of an ID.
func (id *ID) DeepCopy() *ID {
	if id == nil {
		jww.FATAL.Panicf("Failed to create a copy of ID: ID is nil.")
	}

	return copyID(id.Bytes())
}

// String converts an ID to a string via base64 encoding.
func (id *ID) String() string {
	return base64.StdEncoding.EncodeToString(id.Bytes())
}

// Uint64 returns the top 8 bytes of the ID as a uint64.
func (id *ID) Uint64() uint64 {
	return binary.BigEndian.Uint64(id[:8])
}

// GetType returns the ID's type. It is the last byte of the array.
func (id *ID) GetType() Type {
	if id == nil {
		jww.FATAL.Panicf("Failed to get ID type: ID is nil.")
	}

	return Type(id[ArrIDLen-1])
}

// SetType changes the ID type by setting the last byte to the specified type.
func (id *ID) SetType(idType Type) {
	if id == nil {
		jww.FATAL.Panicf("Failed to set ID type: ID is nil.")
	}

	id[ArrIDLen-1] = byte(idType)
}

func NewRandomID(r io.Reader, t Type) (*ID, error) {
	// Generate random bytes
	idBytes := make([]byte, ArrIDLen)
	if _, err := r.Read(idBytes); err != nil {
		return nil, errors.Errorf("failed to generate random bytes for new "+
			"ID: %+v", err)
	}

	// Create ID from bytes
	id := copyID(idBytes)

	// Set new ID type
	id.SetType(t)

	return id, nil
}

// NewIdFromBytes creates a new ID from the supplied byte slice. It is similar
// to Unmarshal() but does not do any error checking. If the data is longer than
// ArrIDLen, then it is truncated. If it is shorter, then the remaining bytes
// are filled with zeroes. This function is for testing purposes only.
func NewIdFromBytes(data []byte, x interface{}) *ID {
	// Ensure that this function is only run in testing environments
	switch x.(type) {
	case *testing.T, *testing.M, *testing.B:
		break
	default:
		panic("NewIdFromBytes() can only be used for testing.")
	}

	return copyID(data)
}

// NewIdFromString creates a new ID from the given string and type. If the
// string is longer than dataLen, then it is truncated. If it is shorter, then
// the remaining bytes are filled with zeroes. This function is for testing
// purposes only.
func NewIdFromString(idString string, idType Type, x interface{}) *ID {
	// Ensure that this function is only run in testing environments
	switch x.(type) {
	case *testing.T, *testing.M, *testing.B:
		break
	default:
		panic("NewIdFromString() can only be used for testing.")
	}

	// Convert the string to bytes and create new ID from it
	newID := NewIdFromBytes([]byte(idString), x)

	// Set the ID type
	newID.SetType(idType)

	return newID
}

// NewIdFromUInt converts the specified uint64 into bytes and returns a new ID
// based off it with the specified ID type. The remaining space of the array is
// filled with zeros. This function is for testing purposes only.
func NewIdFromUInt(idUInt uint64, idType Type, x interface{}) *ID {
	// Ensure that this function is only run in testing environments
	switch x.(type) {
	case *testing.T, *testing.M, *testing.B:
		break
	default:
		panic("NewIdFromUInt() can only be used for testing.")
	}

	// Create the new ID
	newID := new(ID)

	// Convert the uints to bytes
	binary.BigEndian.PutUint64(newID[:], idUInt)

	// Set the ID's type
	newID.SetType(idType)

	return newID
}

// NewIdFromUInt converts the specified uint64 array into bytes and returns a
// new ID based off it with the specified ID type. Unlike NewIdFromUInt(), the
// four uint64s provided fill the entire ID array. This function is for testing
// purposes only.
func NewIdFromUInts(idUInts [4]uint64, idType Type, x interface{}) *ID {
	// Ensure that this function is only run in testing environments
	switch x.(type) {
	case *testing.T, *testing.M, *testing.B:
		break
	default:
		panic("NewIdFromUInts() can only be used for testing.")
	}

	// Create the new ID
	newID := new(ID)

	// Convert the uints to bytes
	for i, idUint := range idUInts {
		binary.BigEndian.PutUint64(newID[i*8:], idUint)
	}

	// Set the ID's type
	newID.SetType(idType)

	return newID
}

// copyID copies the bytes into a new ID.
func copyID(buff []byte) *ID {
	newID := new(ID)
	copy(newID[:], buff)
	return newID
}
