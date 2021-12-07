////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                           //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file                                                               //
////////////////////////////////////////////////////////////////////////////////

package id

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

// Tests that ID.Marshal returns the correct marshaled byte slice of the ID and
// that the returned data is a copy.
func TestID_Marshal(t *testing.T) {
	expected := rngBytes(ArrIDLen, 42, t)
	testID := NewIdFromBytes(expected, t)

	// Check for the correct values
	idBytes := testID.Marshal()
	if !bytes.Equal(expected, idBytes) {
		t.Errorf("Marshal returned unexpected bytes."+
			"\nexpected: %+v\nreceived: %+v", expected, idBytes)
	}

	// Test if the returned bytes are copies
	if &testID[0] == &idBytes[0] {
		t.Errorf("Marshal did not return a copy of the ID data."+
			"\nID pointer:    %+v\nbytes pointer: %+v", &testID[0], &idBytes[0])
	}
}

// Tests that Unmarshal creates a new ID with the expected data.
func TestUnmarshal(t *testing.T) {
	expected := rngBytes(ArrIDLen, 42, t)

	// Unmarshal the bytes into an ID
	testID, err := Unmarshal(expected)
	if err != nil {
		t.Errorf("Unmarshal produced an error: %+v", err)
	}

	// Make sure the ID contents are correct
	if !bytes.Equal(expected, testID[:]) {
		t.Errorf("Unmarshal produced an ID with the incorrect bytes."+
			"\nexpected: %v\nreceived: %v", expected, testID[:])
	}
}

// Tests that Unmarshal produces an error when the given data length is not
// equal to the length of an ID.
func TestUnmarshal_DataLengthError(t *testing.T) {
	invalidIdBytes := rngBytes(ArrIDLen+10, 42, t)
	expectedErr := fmt.Sprintf(unmarshalLenErr, len(invalidIdBytes), ArrIDLen)

	// Unmarshal the bytes into an ID
	_, err := Unmarshal(invalidIdBytes)
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Unmarshal did not product an expected error."+
			"\nexpected: %v\nreceived: %v", expectedErr, err)
	}
}

// Test that an ID that is marshaled and unmarshalled matches the original.
func TestID_Marshal_Unmarshal(t *testing.T) {
	originalID := NewIdFromBytes(rngBytes(ArrIDLen, 42, t), t)

	idBytes := originalID.Marshal()

	testID, err := Unmarshal(idBytes)
	if err != nil {
		t.Errorf("Unmarshal produced an error: %+v", err)
	}

	if originalID != testID {
		t.Errorf("Original ID does not match marshaled/unmarshalled ID."+
			"\nexpected: %s\nreceived: %s", originalID, testID)
	}
}

// Tests that the byte slice returned by ID.Bytes matches the data in the
// original ID and the data is a copy of the values and not the reference.
func TestID_Bytes(t *testing.T) {
	expected := rngBytes(ArrIDLen, 42, t)
	testID := NewIdFromBytes(expected, t)

	// Check for the correct values
	idBytes := testID.Bytes()
	if !bytes.Equal(expected, idBytes) {
		t.Errorf("Bytes returned unexpected bytes."+
			"\nexpected: %+v\nreceived: %+v", expected, idBytes)
	}

	// Test if the returned bytes are copies
	if &testID[0] == &idBytes[0] {
		t.Errorf("Bytes did not return a copy of the ID data."+
			"\nID pointer:    %+v\nbytes pointer: %+v", &testID[0], &idBytes[0])
	}
}

// Consistency test of ID.String.
func TestID_String(t *testing.T) {
	expectedIDs := []string{
		"U4x/lrFkvxuXu59LtHLon1sUhPJSCcnZND6SugndnVLf",
		"KdkEjm+OfQuK4AyZGAqh+XPQaLfRhsO5d2NT1EIScyJX",
		"hlrqczHlHjtl41v0oJfVQzSGffJYAzLv+IWV1btmrbEA",
		"5Nsfgi0t3FtfpMNMQ04N3fy4gWHxgSC916JH4LValO34",
		"uZSFZGuvWGxS68Jz8dpFR1GIDr/1Tp5S9wM5ER+Pqtv/",
		"l1lbRvyllYut7FGP3HFqUD90MzWATAymSimDLxhEo8Up",
		"dOOQ7qkfhLvAjyEdSAOiqpfiSQDBSfdTQyrdJpF8GmCh",
		"Shj2WG8IStR89H94/IpaVLCFTvqLRuVCJ2gXPnyQ5ErL",
		"J1Wbsn0MwguQZZeQuvySrbgo4zpLZGOBVKkPcPXEu48i",
		"fdYg7K6RCDzryum37o3Gt3QUwNfLYVHYce61Re75SWKM",
	}

	for i, expected := range expectedIDs {
		testID := NewIdFromBytes(rngBytes(ArrIDLen, int64(i+42), t), t)

		if testID.String() != expected {
			t.Errorf("String did not output the expected value for ID %d."+
				"\nexpected: %s\nreceived: %s", i, expected, testID.String())
		}
	}
}

// Tests that ID.GetType returns the correct type for each ID type.
func TestID_GetType(t *testing.T) {
	testTypes := []Type{Generic, Gateway, Node, User, Group, NumTypes, 7}
	var testIDs []ID
	for i, idType := range testTypes {
		idDataBytes := rngBytes(dataLen, int64(i+42), t)
		testIDs = append(testIDs,
			NewIdFromBytes(append(idDataBytes, byte(idType)), t))
	}

	for i, testID := range testIDs {
		if testTypes[i] != testID.GetType() {
			t.Errorf("GetType returned the incorrect type (%d)."+
				"\nexpected: %s\nreceived: %s", i, testTypes[i], testID.GetType())
		}
	}
}

// Tests that ID.SetType sets the type of the ID correctly by checking if the
// ID's type changed after calling SetType.
func TestID_SetType(t *testing.T) {
	for i := 0; i < 10; i++ {
		testID := NewIdFromBytes(rngBytes(ArrIDLen, int64(i+42), t), t)
		newTypeID := testID.SetType(testID.GetType() + 1)

		if testID.GetType() == newTypeID.GetType() {
			t.Errorf("SetType failed to change the ID type."+
				"\noriginal ID type: %s\nnew ID type: %s",
				testID.GetType(), newTypeID.GetType())
		}
	}
}

// Tests that ID.Cmp returns the expected comparisons.
func TestID_Cmp(t *testing.T) {
	idA := NewIdFromUInt(1, Generic, t)
	idB := NewIdFromUInt(1, User, t)
	idC := NewIdFromUInt(2, Generic, t)

	testValues := []struct {
		expected int
		a, b     ID
	}{
		{-1, idA, idB},
		{+1, idB, idA},
		{-1, idB, idC},
		{+1, idC, idB},
		{+0, idA, idA},
		{+0, ID{}, ID{}},
	}

	for i, v := range testValues {
		if v.a.Cmp(v.b) != v.expected {
			t.Errorf("Cmp failed to return the expected comparison (%d)."+
				"\nexpected: %d\nreceived: %d", i, v.expected, v.a.Cmp(v.b))
		}
	}

}

// Tests that NewRandomID returns the expected IDs for a given PRNG.
func TestNewRandomID_Consistency(t *testing.T) {
	prng := rand.New(rand.NewSource(42))
	expectedIDs := []string{
		"G5e7n0u0cuifWxSE8lIJydk0PpK6Cd2dUt/Xm012QpsA",
		"egzA1hRMiIU1hBrL4HCbB1gIP2HTdbwCtB30+Rkp4Y8C",
		"nm+C5b1v40zcuoQ+6NY+jE/+HOvqVG2PrBPdGqwEzi4A",
		"h3xVec+iG4KnURCKQu08kDyqQ0ZaeGIGFpeK7QzjxsQA",
		"rv79vgwQKIfhANrNLYhfaSy2B9oAoRwccHHnlqLcLcIA",
		"W3SyySMmgo4rBW44F2WOEGFJiUf980RBDtTBFgI/qOME",
		"a2/tJ//QVpKxNhnnOJZN/ceejVNDc2Yc/WbXT+weG4kD",
		"YpDPK+tCw8onMoVg8arAZ86m6L9G1KsrRoBALF+ygg4A",
		"XTKgmjb5bCCUF0bj7U2mRqmui0+ntPw6ILr6GnXtMnoE",
		"uLDDmup5Uzq/RI0sR50yYHUzkFkUyMwc8J2jng6SnQIE",
	}

	for i, expected := range expectedIDs {
		testID, err := NewRandomID(prng, Type(prng.Intn(int(NumTypes))))
		if err != nil {
			t.Errorf("NewRandomID returned an error (%d): %+v", i, err)
		}
		if testID.String() != expected {
			t.Errorf("NewRandomID did not generate the expected ID."+
				"\nexpected: %s\nreceived: %s", expected, testID)
		}
	}
}

// Tests that NewRandomID returns unique IDs.
func TestNewRandomID_Unique(t *testing.T) {
	prng := rand.New(rand.NewSource(42))
	ids := map[ID]struct{}{}

	for i := 0; i < 100; i++ {
		testID, err := NewRandomID(prng, Type(prng.Intn(int(NumTypes))))
		if err != nil {
			t.Errorf("NewRandomID returned an error (%d): %+v", i, err)
		}

		if _, exists := ids[testID]; exists {
			t.Errorf("NewRandomID did not generate a unique ID (%d).\nID: %s",
				i, testID)
		} else {
			ids[testID] = struct{}{}
		}
	}
}

// This tests uses a custom PRNG which generates an ID with a base64 encoding
// starting with a special character. NewRandomID should force another call
// to prng.Read, and this call should return an ID with encoding starting with an
// alphanumeric character. This test fails if it hangs forever (PRNG error) or
// the ID has an encoding beginning with a special character (NewRandomID error).
func TestNewRandomID_SpecialCharacter(t *testing.T) {
	prng := newAlphanumericPRNG()
	testId, err := NewRandomID(prng, 0)
	if err != nil {
		t.Errorf("NewRandomID() returned an error: %+v", err)
	}

	if !regexAlphanumeric.MatchString(string(testId.String()[0])) {
		t.Errorf("Should not have an ID starting with a special character")
	}

}

// Tests that NewRandomID returns an error when the io reader encounters an
// error.
func TestNewRandomID_ReaderError(t *testing.T) {
	_, err := NewRandomID(strings.NewReader(""), Generic)
	if err == nil {
		t.Error("NewRandomID failed to return an error when the reader " +
			"failed.")
	}
}

// Tests that NewIdFromBytes creates a new ID with the correct contents.
func TestNewIdFromBytes(t *testing.T) {
	expectedBytes := rngBytes(ArrIDLen, 42, t)
	testID := NewIdFromBytes(expectedBytes, t)

	if !bytes.Equal(expectedBytes, testID[:]) {
		t.Errorf("NewIdFromBytes produced an ID with the incorrect bytes."+
			"\nexpected: %v\nreceived: %v", expectedBytes, testID[:])
	}
}

// Tests that an ID can be JSON marshaled and unmarshalled.
func TestID_MarshalJSON_UnmarshalJSON(t *testing.T) {
	testID := NewIdFromBytes(rngBytes(ArrIDLen, 42, t), t)

	jsonData, err := json.Marshal(testID)
	if err != nil {
		t.Errorf("json.Marshal returned an error: %+v", err)
	}
	t.Logf("%s", jsonData)

	newID := &ID{}
	err = json.Unmarshal(jsonData, newID)
	if err != nil {
		t.Errorf("json.Unmarshal returned an error: %+v", err)
	}

	if testID != *newID {
		t.Errorf("Failed the JSON marshal and unmarshal ID."+
			"\noriginal ID: %s\nreceived ID: %s", testID, *newID)
	}
}

// Error path: supplied data is invalid JSON.
func TestID_UnmarshalJSON_JsonUnmarshalError(t *testing.T) {
	expectedErr := "invalid character"
	id := ID{}
	err := id.UnmarshalJSON([]byte("invalid JSON"))
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("UnmarshalJSON failed to return the expected error for "+
			"invalid JSON.\nexpected: %s\nreceived: %+v", expectedErr, err)
	}
}

// Error path: supplied data is valid JSON but an invalid ID.
func TestID_UnmarshalJSON_IdUnmarshalError(t *testing.T) {
	expectedErr := fmt.Sprintf(unmarshalLenErr, 0, ArrIDLen)
	id := ID{}
	err := id.UnmarshalJSON([]byte("\"\""))
	t.Log(err)
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("UnmarshalJSON failed to return the expected error for "+
			"invalid ID.\nexpected: %s\nreceived: %+v", expectedErr, err)
	}
}

// Tests that NewIdFromBytes panics when given a nil testing object.
func TestNewIdFromBytes_TestError(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("NewIdFromBytes did not panic when it received a nil " +
				"testing object when it should have.")
		}
	}()

	// Call the function with nil testing object
	_ = NewIdFromBytes(rngBytes(ArrIDLen, 42, t), nil)
}

// Tests that NewIdFromString creates a new ID from string correctly. The new ID
// is created from a string that is 32 bytes long so that no truncation or
// padding is required. The test checks that the original string is still
// present in the data.
func TestNewIdFromString(t *testing.T) {
	// Test values
	expectedIdString := "Test ID string of correct length"
	expectedType := Generic
	var expectedID ID
	copy(expectedID[:], append([]byte(expectedIdString), byte(expectedType)))

	// Create the ID and check its contents
	testID := NewIdFromString(expectedIdString, expectedType, t)

	// Check if the new ID matches the expected ID
	if expectedID != testID {
		t.Errorf("NewIdFromString produced an ID with the incorrect data."+
			"\nexpected: %v\nreceived: %v", expectedID[:], testID[:])
	}

	// Check if the original string is still in the first 32 bytes
	newIdString := string(testID.Bytes()[:ArrIDLen-1])
	if expectedIdString != newIdString {
		t.Errorf("NewIdFromString did not correctly convert the original "+
			"string to bytes.\nexpected string: %#v\nreceived string: %#v"+
			"\nexpected bytes: %v\nreceived bytes: %v",
			expectedIdString, newIdString,
			[]byte(expectedIdString), testID.Bytes()[:ArrIDLen-1])
	}
}

// Tests that NewIdFromString panics when given a nil testing object.
func TestNewIdFromString_TestError(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("NewIdFromString did not panic when it received a nil " +
				"testing object when it should have.")
		}
	}()

	// Call the function with nil testing object
	_ = NewIdFromString("test", Generic, nil)
}

// Tests that NewIdFromUInt creates a new ID with the correct contents by
// converting the ID back into a uint and comparing it to the original.
func TestNewIdFromUInt(t *testing.T) {
	// Expected values
	expectedUint := rand.Uint64()

	// Create the ID and check its contents
	newID := NewIdFromUInt(expectedUint, Generic, t)
	idUint := binary.BigEndian.Uint64(newID[:ArrIDLen-1])

	if expectedUint != idUint {
		t.Errorf("NewIdFromUInt produced an ID from uint incorrectly."+
			"\nexpected: %v\nreceived: %v", expectedUint, idUint)
	}
}

// Tests that NewIdFromUInt panics when given a nil testing object.
func TestNewIdFromUInt_TestError(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("NewIdFromUInt did not panic when it received a nil " +
				"testing object when it should have.")
		}
	}()

	// Call the function with nil testing object
	_ = NewIdFromUInt(rand.Uint64(), Generic, nil)
}

// Tests that NewIdFromUInts creates a new ID with the correct contents by
// converting the ID back into a slice of uints and comparing it to the
// original.
func TestNewIdFromUInts(t *testing.T) {
	// Expected values
	expectedUints := [4]uint64{rand.Uint64(), rand.Uint64(),
		rand.Uint64(), rand.Uint64()}

	// Create the ID and check its contents
	newID := NewIdFromUInts(expectedUints, Generic, t)
	idUints := [4]uint64{}
	for i := range idUints {
		idUints[i] = binary.BigEndian.Uint64(newID[i*8 : (i+1)*8])
	}

	if !reflect.DeepEqual(expectedUints, idUints) {
		t.Errorf("NewIdFromUInts produced an ID from uints incorrectly."+
			"\nexpected: %#v\nreceived: %#v", expectedUints, idUints)
	}
}

// Tests that NewIdFromUInts panics when given a nil testing object.
func TestNewIdFromUInts_TestError(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("NewIdFromUInts did not panic when it received a nil " +
				"testing object when it should have.")
		}
	}()

	// Call the function with nil testing object
	newUint64s := [4]uint64{rand.Uint64(), rand.Uint64(),
		rand.Uint64(), rand.Uint64()}
	_ = NewIdFromUInts(newUint64s, Generic, nil)
}

// Generates a byte slice of the specified length containing random numbers.
func rngBytes(length int, seed int64, t *testing.T) []byte {
	prng := rand.New(rand.NewSource(seed))

	// Create new byte slice of the correct size
	idBytes := make([]byte, length)

	// Create random bytes
	_, err := prng.Read(idBytes)
	if err != nil {
		t.Fatalf("Failed to generate random bytes: %+v", err)
	}

	return idBytes
}

// alphaNumericPRNG is a custom PRNG which adheres to the io.Reader interface. This is used for
// testing special characters.
type alphaNumericPRNG struct {
	counter uint
}

// newAlphanumericPRNG generates a copy of the custom prng alphaNumericPRNG.
func newAlphanumericPRNG() *alphaNumericPRNG {
	return &alphaNumericPRNG{counter: 0}
}

// Hardcoded byte array which generates a base64 string starting with a special character.
// Expected encoding is "/ABBb6YWlbkgLmg2Ohx4f0eE4K7Zx4VkGE4THx58gR8A",
// any other encoding output indicates that something is wrong (likely a dependency).
var hardCodedSpecialCharacter = []byte{252, 0, 65, 111, 166, 22, 149, 185, 32, 46, 104, 54, 58, 28, 120, 127, 71, 132, 224, 174, 217, 199, 133, 100, 24, 78, 19, 31, 30, 124, 129, 31, 189}

// Hardcoded byte array which generates a base64 string starting with an alphanumeric character.
// Expected encoding is "6iMc/s5V6MSzD6+DDQthcfA53w7wY988cenRkjxNwIcD",
// any other encoding output indicates that something is wrong (likely a dependency).
var hardcodedAlphaNumeric = []byte{234, 35, 28, 254, 206, 85, 232, 196, 179, 15, 175, 131, 13, 11, 97, 113, 240, 57, 223, 14, 240, 99, 223, 60, 113, 233, 209, 146, 60, 77, 192, 135, 19}

// Read will copy a value into byte slice p. On the first call to alphaNumericPRNG.Read()
// a hardcoded value with a base64 encoding starting with a special character will be returned.
// For any other call to alphaNumericPRNG.Read(), a hardcoded value with a base64 encoding starting with an
// alphanumeric character will be returned.
func (prng *alphaNumericPRNG) Read(p []byte) (n int, err error) {
	defer func() {
		prng.counter++
	}()

	if prng.counter == 0 {
		copy(p, hardCodedSpecialCharacter)
		return len(p), nil
	} else {
		copy(p, hardcodedAlphaNumeric)
		return len(p), nil
	}
}
