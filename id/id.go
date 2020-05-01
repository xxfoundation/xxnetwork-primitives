////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

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
	"testing"
)

const (
	// Length of the ID data
	dataLen = 32

	// Length of the full ID array
	ArrIDLen = dataLen + 1
)

// ID is an array holding the generic identifier. The first 32 bytes hold the ID
// data and the last byte holds the ID type.
type ID [ArrIDLen]byte

// Marshal takes an ID and copies the data into the wire format.
func (i *ID) Marshal() []byte {
	return i.Bytes()
}

// Unmarshal takes in the ID wire format and copies the data into an ID.
func Unmarshal(data []byte) (*ID, error) {
	// Return an error if the length of data is incorrect
	if len(data) != ArrIDLen {
		return nil, errors.Errorf("Could not marshal byte slice to ID: "+
			"length of data must be %d, length received was %d",
			ArrIDLen, len(data))
	}

	newID := new(ID)
	copy(newID[:], data)

	return newID, nil
}

// Bytes returns a copy of an ID as a byte slice. Note that Bytes() is used by
// Marshal() and any changes made here will affect how Marshal() functions.
func (i *ID) Bytes() []byte {
	newBytes := make([]byte, ArrIDLen)
	copy(newBytes, i[:])

	return newBytes
}

// Cmp determines whether two IDs are the same. Returns true if they are equal
// and false otherwise.
func (i *ID) Cmp(y *ID) bool {
	return *i == *y
}

// DeepCopy creates a new copy of an ID.
func (i *ID) DeepCopy() *ID {
	newID := new(ID)
	copy(newID[:], i[:])

	return newID
}

// String converts an ID to a string via base64 encoding.
func (i *ID) String() string {
	return base64.StdEncoding.EncodeToString(i.Bytes())
}

// GetType returns the ID's type. It is the last byte of the array.
func (i *ID) GetType() Type {
	return Type(i[ArrIDLen-1])
}

// SetType changes the ID type by setting the last byte to the specified type.
func (i *ID) SetType(idType Type) {
	i[ArrIDLen-1] = byte(idType)
}

// NewIdFromBytes creates a new ID from a copy of the data. It is similar to
// Unmarshal() but does not do any error checking. If the data is longer than
// ArrIDLen, then it is truncated. If it is shorter, then the remaining bytes
// are filled with zeroes. This function is for testing purposes only.
func NewIdFromBytes(data []byte, t *testing.T) *ID {
	// Ensure that this function is only run in testing environments
	if t == nil {
		panic("NewIdFromBytes() can only be used for testing.")
	}

	newID := new(ID)
	copy(newID[:], data[:])

	return newID
}

// NewIdFromString creates a new ID from the given string and type. If the
// string is longer than dataLen, then it is truncated. If it is shorter, then
// the remaining bytes are filled with zeroes. This function is for testing
// purposes only.
func NewIdFromString(idString string, idType Type, t *testing.T) *ID {
	// Ensure that this function is only run in testing environments
	if t == nil {
		panic("NewIdFromString() can only be used for testing.")
	}

	// Convert the string to bytes and create new ID from it
	newID := NewIdFromBytes([]byte(idString), t)

	// Set the ID type
	newID.SetType(idType)

	return newID
}

// NewIdFromUInt converts the specified uint64 into bytes and returns a new ID
// based off it with the specified ID type. The remaining space of the array is
// filled with zeros. This function is for testing purposes only.
func NewIdFromUInt(idUInt uint64, idType Type, t *testing.T) *ID {
	// Ensure that this function is only run in testing environments
	if t == nil {
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
func NewIdFromUInts(idUInts [4]uint64, idType Type, t *testing.T) *ID {
	// Ensure that this function is only run in testing environments
	if t == nil {
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
