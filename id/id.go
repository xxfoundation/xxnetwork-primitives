////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Package id contains the generic ID type, which is a byte array that
// represents an entity ID. The first bytes in the array contain the actual ID
// data while the last byte contains the ID type, which is either generic,
// gateway, node, or user. IDs can be hard coded or generated using a
// cryptographic function found in crypto.
package id

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"regexp"
	"testing"

	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
)

const (
	// Length of the ID data.
	dataLen = 32

	// ArrIDLen is the length of the full ID array.
	ArrIDLen = dataLen + 1

	// Contains the regular expression to search for an alphanumeric string.
	alphanumeric string = "^[a-zA-Z0-9]+$"
)

// regexAlphanumeric is the regex for searching for an alphanumeric string.
var regexAlphanumeric = regexp.MustCompile(alphanumeric)

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

// Bytes returns a copy of an ID as a byte slice. Note that Bytes is used by
// Marshal and any changes made here will affect how Marshal functions.
func (id *ID) Bytes() []byte {
	if id == nil {
		jww.FATAL.Panicf("%+v", errors.New(
			"Failed to get bytes of ID: ID is nil."))
	}

	newBytes := make([]byte, ArrIDLen)
	copy(newBytes, id[:])

	return newBytes
}

// Cmp determines whether two IDs are equal. Returns true if they are equal and
// false otherwise.
//
// Deprecated: Use ID.Equal instead.
func (id *ID) Cmp(y *ID) bool {
	return id.Equal(y)
}

// Equal determines whether two IDs are equal. Returns true if they are equal
// and false otherwise.
func (id *ID) Equal(y *ID) bool {
	if id == nil || y == nil {
		jww.FATAL.Panicf("%+v", errors.Errorf("Failed to compare IDs: one or both IDs are nil."))
	}

	return *id == *y
}

// Compare returns an integer comparing the two IDs lexicographically.
// The result will be 0 if id == y, -1 if id < y, and +1 if id > y.
func (id *ID) Compare(y *ID) int {
	if id == nil || y == nil {
		jww.FATAL.Panicf("%+v", errors.New(
			"Failed to compare IDs: one or both IDs are nil."))
	}

	return bytes.Compare(id[:], y[:])
}

// Less returns true if id is less than y.
func (id *ID) Less(y *ID) bool {
	return id.Compare(y) == -1
}

// DeepCopy creates a copy of an ID.
func (id *ID) DeepCopy() *ID {
	if id == nil {
		jww.FATAL.Panicf("%+v", errors.New(
			"Failed to create a copy of ID: ID is nil."))
	}

	return copyID(id[:])
}

// String converts an ID to a string via base64 encoding.
func (id *ID) String() string {
	if id == nil {
		jww.FATAL.Panicf("%+v", errors.New(
			"Failed to create string of ID: ID is nil."))
	}

	return base64.StdEncoding.EncodeToString(id[:])
}

// GetType returns the ID's type. It is the last byte of the array.
func (id *ID) GetType() Type {
	if id == nil {
		jww.FATAL.Panicf("%+v", errors.New(""+
			""+
			"Failed to get ID type: ID is nil."))
	}

	return Type(id[ArrIDLen-1])
}

// SetType changes the ID type by setting the last byte to the specified type.
func (id *ID) SetType(idType Type) {
	if id == nil {
		jww.FATAL.Panicf("%+v", errors.New(
			"Failed to set ID type: ID is nil."))
	}

	id[ArrIDLen-1] = byte(idType)
}

// MarshalJSON marshals the [ID] into valid JSON. This function adheres to the
// [json.Marshaler] interface.
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.Marshal())
}

// UnmarshalJSON unmarshalls the JSON into the [ID]. This function adheres to
// the [json.Unmarshaler] interface.
func (id *ID) UnmarshalJSON(b []byte) error {
	var buff []byte
	if err := json.Unmarshal(b, &buff); err != nil {
		return err
	}

	newID, err := Unmarshal(buff)
	if err != nil {
		return err
	}

	*id = *newID

	return nil
}

// MarshalText marshals the [ID] into base 64 encoded text. This function
// adheres to the [encoding.TextMarshaler] interface. This allows for the JSON
// marshalling of non-referenced IDs in maps (e.g., map[ID]int).
func (id ID) MarshalText() (text []byte, err error) {
	return []byte(base64.RawStdEncoding.EncodeToString(id[:])), nil
}

// UnmarshalText unmarshalls the text into an [ID]. This function adheres to the
// [encoding.TextUnmarshaler] interface. This allows for the JSON unmarshalling
// of non-referenced IDs in maps (e.g., map[ID]int).
func (id *ID) UnmarshalText(text []byte) error {
	idBytes, err := base64.RawStdEncoding.DecodeString(string(text))
	if err != nil {
		return err
	}

	newID, err := Unmarshal(idBytes)
	if err != nil {
		return err
	}

	copy(id[:], newID[:])
	return nil
}

// HexEncode encodes the ID without the 33rd type byte.
func (id *ID) HexEncode() string {
	return "0x" + hex.EncodeToString(id.Bytes()[:32])
}

// NewRandomID generates a random ID using the passed in io.Reader r
// and sets the ID to Type t. If the base64 string of the generated
// ID does not begin with an alphanumeric character, then another ID
// is generated until the encoding begins with an alphanumeric character.
func NewRandomID(r io.Reader, t Type) (*ID, error) {
	for {
		// Generate random bytes
		idBytes := make([]byte, ArrIDLen)
		if _, err := r.Read(idBytes); err != nil {
			return nil, errors.Errorf(
				"failed to generate random bytes for new ID: %+v", err)
		}

		// Create ID from bytes
		id := copyID(idBytes)

		// Set new ID type
		id.SetType(t)

		// Avoid the first character being a special character
		base64Id := id.String()
		if regexAlphanumeric.MatchString(string(base64Id[0])) {
			return id, nil
		}
	}
}

// NewRandomTestID generates a random ID using the passed in io.Reader r and
// sets the ID to Type t. If the base64 string of the generated ID does not
// begin with an alphanumeric character, then another ID is generated until the
// encoding begins with an alphanumeric character.
//
// This function is intended for testing purposes.
func NewRandomTestID(r io.Reader, t Type, x interface{}) *ID {
	// Ensure that this function is only run in testing environments
	switch x.(type) {
	case *testing.T, *testing.M, *testing.B:
		break
	default:
		jww.FATAL.Panicf("NewRandomTestID can only be used for testing.")
	}

	id, err := NewRandomID(r, t)
	if err != nil {
		jww.FATAL.Panicf(
			"failed to generate random bytes for new ID: %+v", err)
	}

	return id
}

// NewIdFromBytes creates a new ID from the supplied byte slice. It is similar
// to Unmarshal but does not do any error checking. If the data is longer than
// ArrIDLen, then it is truncated. If it is shorter, then the remaining bytes
// are filled with zeroes. This function is for testing purposes only.
func NewIdFromBytes(data []byte, x interface{}) *ID {
	// Ensure that this function is only run in testing environments
	switch x.(type) {
	case *testing.T, *testing.M, *testing.B:
		break
	default:
		jww.FATAL.Panicf("NewIdFromBytes can only be used for testing.")
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
		jww.FATAL.Panicf("NewIdFromString can only be used for testing.")
	}

	// Convert the string to bytes and create new ID from it
	newID := NewIdFromBytes([]byte(idString), x)

	// Set the ID type
	newID.SetType(idType)

	return newID
}

var (
	// nonBase64Regex matches any character that are not base 64 compatible
	// except strings.
	nonBase64Regex = regexp.MustCompile(`[^a-zA-Z0-9+/\s]+`)

	// whitespaceRegex matches one or more whitespace characters.
	whitespaceRegex = regexp.MustCompile(`\s+`)
)

// NewIdFromBase64String creates a new ID that when base 64 encoded, looks like
// the passed in base64String. This function is for testing purposes only.
//
// If the string is longer than the data portion of the base 64 string, then it
// is truncated. If it is shorter, then the remaining bytes are filled with
// zeroes. The string is made to be base 64 compatible by replacing one or more
// consecutive white spaces with a plus "+" and stripping all other invalid
// characters (any character that is not an upper- and lower-case alphabet
// character (A–Z, a–z), numeral (0–9), or the plus "+" and slash "/" symbols).
func NewIdFromBase64String(base64String string, idType Type, x interface{}) *ID {
	// Ensure that this function is only run in testing environments
	switch x.(type) {
	case *testing.T, *testing.M, *testing.B:
		break
	default:
		jww.FATAL.Panicf("NewIdFromBase64String can only be used for testing.")
	}

	// Convert the string to bytes and create new ID from it
	newID := NewIdFromBytes([]byte{}, x)

	// Set the ID type
	newID.SetType(idType)

	// Escape ID string to be base 64 compatible by replacing all strings with +
	// and stripping all other invalid characters
	base64String = nonBase64Regex.ReplaceAllString(base64String, "")
	base64String = whitespaceRegex.ReplaceAllString(base64String, "+")

	b64Str := base64.StdEncoding.EncodeToString(newID.Marshal())

	// Trim the string if it is over the max length
	if len(base64String) > len(b64Str)-1 {
		base64String = base64String[:len(b64Str)-1]
	}

	// Concatenate the string with the rest of the generated bytes and type
	base64String = base64String + b64Str[len(base64String):]

	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		jww.FATAL.Panicf("Failed to decode string: %+v", err)
	}

	return NewIdFromBytes(data, x)
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
		jww.FATAL.Panicf("NewIdFromUInt can only be used for testing.")
	}

	// Create the new ID
	newID := new(ID)

	// Convert the uints to bytes
	binary.BigEndian.PutUint64(newID[:], idUInt)

	// Set the ID's type
	newID.SetType(idType)

	return newID
}

// NewIdFromUInts converts the specified uint64 array into bytes and returns a
// new ID based off it with the specified ID type. Unlike NewIdFromUInt, the
// four uint64s provided fill the entire ID array. This function is for testing
// purposes only.
func NewIdFromUInts(idUInts [4]uint64, idType Type, x interface{}) *ID {
	// Ensure that this function is only run in testing environments
	switch x.(type) {
	case *testing.T, *testing.M, *testing.B:
		break
	default:
		jww.FATAL.Panicf("NewIdFromUInts can only be used for testing.")
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
