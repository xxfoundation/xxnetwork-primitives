////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package id

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests that ID.Marshal returns the correct byte slice of an ID and that it is
// a copy. This test ensures that ID.Marshal uses ID.Bytes.
func TestID_Marshal(t *testing.T) {
	// Test values
	expectedBytes := newRandomBytes(ArrIDLen, t)
	testID := NewIdFromBytes(expectedBytes, t)

	// Check for the correct values
	testVal := testID.Marshal()
	if !bytes.Equal(expectedBytes, testVal) {
		t.Errorf("Marshal returned the incorrect byte slice of the ID"+
			"\n\texpected: %+v\n\treceived: %+v", expectedBytes, testVal)
	}

	// Test if the returned bytes are copies
	if &testID[0] == &testVal[0] {
		t.Errorf("Marshal did not return a copy when it should have."+
			"\n\texpected: any value except %+v\n\treceived: %+v",
			&testID[0], &testVal[0])
	}
}

// Tests that Unmarshal creates a new ID with the correct contents and does not
// return an error.
func TestUnmarshal(t *testing.T) {
	// Test values
	expectedBytes := newRandomBytes(ArrIDLen, t)

	// Unmarshal the bytes into an ID
	newID, err := Unmarshal(expectedBytes)

	// Make sure no error occurred
	if err != nil {
		t.Errorf("Unmarshal produced an unexpected error."+
			"\n\texpected: %v\n\treceived: %v", nil, err)
	}

	// Make sure the ID contents are correct
	if !bytes.Equal(expectedBytes, newID[:]) {
		t.Errorf("Unmarshal produced an ID with the incorrect bytes."+
			"\n\texpected: %v\n\treceived: %v", expectedBytes, newID[:])
	}
}

// Tests that Unmarshal produces an error when the given data length is not
// equal to the length of an ID and that the ID returned is nil.
func TestUnmarshal_DataLengthError(t *testing.T) {
	// Test values
	expectedBytes := newRandomBytes(ArrIDLen+10, t)
	expectedError := fmt.Errorf("could not marshal byte slice to ID: "+
		"length of data must be %d, length received was %d",
		ArrIDLen, len(expectedBytes))

	// Unmarshal the bytes into an ID
	newID, err := Unmarshal(expectedBytes)

	// Make sure an error occurs
	if err == nil {
		t.Errorf("Unmarshal did not product an expected error."+
			"\n\texpected: %v\n\treceived: %v", expectedError, err)
	}

	// Make sure the returned ID is nil
	if newID != nil {
		t.Errorf("Unmarshal produced a non-nil ID on error."+
			"\n\texpected: %v\n\treceived: %v", nil, newID)
	}
}

// Tests that Bytes returns the correct byte slice of an ID and that it is
// a copy.
func TestID_Bytes(t *testing.T) {
	// Test values
	expectedBytes := newRandomBytes(ArrIDLen, t)
	testID := NewIdFromBytes(expectedBytes, t)

	// Check for the correct values
	testVal := testID.Bytes()
	if !bytes.Equal(expectedBytes, testVal) {
		t.Errorf("Bytes returned the incorrect byte slice of the ID"+
			"\n\texpected: %+v\n\treceived: %+v", expectedBytes, testVal)
	}

	// Test if the returned bytes are copies
	if &testID[0] == &testVal[0] {
		t.Errorf("Bytes did not return a copy when it should have."+
			"\n\texpected: any value except %+v\n\treceived: %+v",
			&testID[0], &testVal[0])
	}
}

// Tests that ID.Bytes panics when the ID is nil.
func TestID_Bytes_NilError(t *testing.T) {
	var id *ID

	defer func() {
		if r := recover(); r == nil {
			t.Error("Bytes failed to panic when the ID is nil.")
		}
	}()

	_ = id.Bytes()
}

// Tests that Cmp returns the correct value when comparing equal and unequal
// IDs.
func TestID_Cmp(t *testing.T) {
	// Test values
	randomBytes1 := newRandomBytes(ArrIDLen, t)
	randomBytes2 := newRandomBytes(ArrIDLen, t)
	randomBytes3 := make([]byte, ArrIDLen)
	copy(randomBytes3, randomBytes2)
	randomBytes3[ArrIDLen-1] = ^randomBytes3[ArrIDLen-1]
	testID1 := NewIdFromBytes(randomBytes1, t)
	testID2 := NewIdFromBytes(randomBytes1, t)
	testID3 := NewIdFromBytes(randomBytes2, t)
	testID4 := NewIdFromBytes(randomBytes3, t)

	// Compare two equal IDs
	testVal := testID1.Cmp(testID2)
	if !testVal {
		t.Errorf("Cmp incorrectly determined the two IDs are not equal."+
			"\n\texpected: %+v\n\treceived: %+v", true, testVal)
	}

	// Compare two unequal IDs
	testVal = testID1.Cmp(testID3)
	if testVal {
		t.Errorf("Cmp incorrectly determined the two IDs are equal."+
			"\n\texpected: %+v\n\treceived: %+v", false, testVal)
	}

	// Compare two almost equal IDs
	testVal = testID3.Cmp(testID4)
	if testVal {
		t.Errorf("Cmp incorrectly determined the two IDs are equal."+
			"\n\texpected: %+v\n\treceived: %+v", false, testVal)
	}
}

// Tests that ID.Cmp panics when both IDs are nil.
func TestID_Cmp_NilError(t *testing.T) {
	var idA, idB *ID

	defer func() {
		if r := recover(); r == nil {
			t.Error("Cmp failed to panic when both IDs are nil.")
		}
	}()

	_ = idA.Cmp(idB)
}

// Tests that ID.Equal returns the correct value when comparing equal and
// unequal IDs.
func TestID_Equal(t *testing.T) {
	tests := []struct {
		x, y  *ID
		equal bool
	}{
		{NewRandomTestID(rand.New(rand.NewSource(42)), User, t),
			NewRandomTestID(rand.New(rand.NewSource(42)), User, t), true},
		{NewRandomTestID(rand.New(rand.NewSource(42)), User, t),
			NewRandomTestID(rand.New(rand.NewSource(32)), User, t), false},
		{NewIdFromBytes([]byte{1, 2}, t), NewIdFromBytes([]byte{1, 3}, t), false},
	}

	for i, tt := range tests {
		equal := tt.x.Equal(tt.y)
		if equal != tt.equal {
			t.Errorf("Failed to compare IDs %s and %s (%d)."+
				"\nexpected: %t\nreceived: %t", tt.x, tt.y, i, tt.equal, equal)
		}
	}
}

// Tests that ID.Equal panics when either and both IDs are nil.
func TestID_Equal_NilError(t *testing.T) {
	var idA, idB *ID

	defer func() {
		if r := recover(); r == nil {
			t.Error("Equal failed to panic when both IDs are nil.")
		}
	}()

	idA.Equal(idB)
}

// Unit test of ID.Compare.
func TestID_Compare(t *testing.T) {
	tests := []struct {
		x, y    *ID
		compare int
	}{
		{NewRandomTestID(rand.New(rand.NewSource(42)), User, t),
			NewRandomTestID(rand.New(rand.NewSource(42)), User, t), 0},
		{NewIdFromBytes([]byte{1, 2}, t), NewIdFromBytes([]byte{1, 3}, t), -1},
		{NewIdFromBytes([]byte{9, 9}, t), NewIdFromBytes([]byte{1, 3}, t), 1},
	}

	for i, tt := range tests {
		compare := tt.x.Compare(tt.y)
		if compare != tt.compare {
			t.Errorf("Failed to compare IDs %s and %s (%d)."+
				"\nexpected: %d\nreceived: %d",
				tt.x, tt.y, i, tt.compare, compare)
		}
	}
}

// Tests that ID.Compare panics when either and both IDs are nil.
func TestID_Compare_NilError(t *testing.T) {
	var idA, idB *ID

	defer func() {
		if r := recover(); r == nil {
			t.Error("Compare failed to panic when both IDs are nil.")
		}
	}()

	idA.Compare(idB)
}

// Test that DeepCopy returns a copy with the same contents as the original
// and where the pointers are different.
func TestID_DeepCopy(t *testing.T) {
	// Test values
	expectedID := NewIdFromBytes(newRandomBytes(ArrIDLen, t), t)

	// Test if the contents are equal
	testVal := expectedID.DeepCopy()
	if !reflect.DeepEqual(expectedID, testVal) {
		t.Errorf("DeepCopy returned a copy with the wrong contents."+
			"\n\texpected: %+v\n\treceived: %+v", expectedID, testVal)
	}

	// Test if the returned bytes are copies
	if &expectedID[0] == &testVal[0] {
		t.Errorf("DeepCopy did not return a copy when it should have."+
			"\n\texpected: any value except %+v\n\treceived: %+v",
			&expectedID[0], &testVal[0])
	}
}

// Tests that ID.DeepCopy panics when the ID is nil.
func TestID_DeepCopy_NilError(t *testing.T) {
	var id *ID

	defer func() {
		if r := recover(); r == nil {
			t.Error("DeepCopy failed to panic when the ID is nil.")
		}
	}()

	_ = id.DeepCopy()
}

// Tests that the base64 encoded string returned by String can be decoded into
// the original ID.
func TestID_String(t *testing.T) {
	// Expected values
	expectedBytes := newRandomBytes(ArrIDLen, t)
	expectedID := NewIdFromBytes(expectedBytes, t)

	// Encode into string
	stringID := expectedID.String()

	// Decode the string and check
	newID, err := base64.StdEncoding.DecodeString(stringID)
	if err != nil {
		t.Fatalf("Failed to decode string returned by String:\n%v", err)
	}

	if !bytes.Equal(expectedBytes, newID) {
		t.Errorf("String did not encode the string correctly."+
			"The decoded strings differ.\n\texpected: %v\n\treceived: %v",
			expectedBytes, newID)
	}
}

// Tests that GetType returns the correct type for each ID type.
func TestID_GetType(t *testing.T) {
	// Test values
	testTypes := []Type{Generic, Gateway, Node, User, Group, NumTypes, 7}
	randomBytes := [][]byte{
		newRandomBytes(ArrIDLen-1, t), newRandomBytes(ArrIDLen-1, t),
		newRandomBytes(ArrIDLen-1, t), newRandomBytes(ArrIDLen-1, t),
		newRandomBytes(ArrIDLen-1, t), newRandomBytes(ArrIDLen-1, t),
	}
	testIDs := []*ID{
		NewIdFromBytes(append(randomBytes[0], byte(testTypes[0])), t),
		NewIdFromBytes(append(randomBytes[1], byte(testTypes[1])), t),
		NewIdFromBytes(append(randomBytes[2], byte(testTypes[2])), t),
		NewIdFromBytes(append(randomBytes[3], byte(testTypes[3])), t),
		NewIdFromBytes(append(randomBytes[4], byte(testTypes[4])), t),
		NewIdFromBytes(append(randomBytes[5], byte(testTypes[5])), t),
	}

	for i, testID := range testIDs {
		testVal := testID.GetType()
		if testTypes[i] != testVal {
			t.Errorf("GetType returned the incorrect type."+
				"\n\texpected: %v\n\treceived: %v", testTypes[i], testVal)
		}
	}
}

// Tests that ID.GetType panics when the ID is nil.
func TestID_GetType_NilError(t *testing.T) {
	var id *ID

	defer func() {
		if r := recover(); r == nil {
			t.Error("GetType failed to panic when the ID is nil.")
		}
	}()

	_ = id.GetType()
}

// Tests that SetType sets the type of the ID correctly by checking if the
// ID's type changed after calling SetType.
func TestID_SetType(t *testing.T) {
	// Test values
	expectedType := Node
	testType := Generic
	testBytes := newRandomBytes(dataLen, t)
	testID := NewIdFromBytes(append(testBytes, byte(testType)), t)

	// Change the ID
	testID.SetType(expectedType)

	// Check the ID
	testVal := testID.GetType()
	if expectedType != testVal {
		t.Errorf("SetType did not set the ID type correctly."+
			"\n\texpected: %v\n\treceived: %v", expectedType, testVal)
	}
}

// Tests that ID.SetType panics when the ID is nil.
func TestID_SetType_NilError(t *testing.T) {
	var id *ID

	defer func() {
		if r := recover(); r == nil {
			t.Error("SetType failed to panic when the ID is nil.")
		}
	}()

	id.SetType(Generic)
}

// Tests that an ID can be JSON marshaled and unmarshalled.
func TestID_MarshalJSON_UnmarshalJSON(t *testing.T) {
	testID := NewIdFromBytes(rngBytes(ArrIDLen, 42, t), t)

	jsonData, err := json.Marshal(testID)
	if err != nil {
		t.Errorf("json.Marshal returned an error: %+v", err)
	}

	newID := &ID{}
	err = json.Unmarshal(jsonData, newID)
	if err != nil {
		t.Errorf("json.Unmarshal returned an error: %+v", err)
	}

	if *testID != *newID {
		t.Errorf("Failed the JSON marshal and unmarshal ID."+
			"\noriginal ID: %s\nreceived ID: %s", testID, newID)
	}
}

// Tests that an ID can be JSON marshaled and unmarshalled.
func TestID_TextMarshal_TextUnmarshal(t *testing.T) {
	testID := NewRandomTestID(rand.New(rand.NewSource(42534)), User, t)

	testMap := make(map[ID]int)
	testMap[*testID] = 8675309

	jsonData, err := json.Marshal(testMap)
	require.NoError(t, err)

	newMap := make(map[ID]int)
	err = json.Unmarshal(jsonData, &newMap)
	require.NoError(t, err)

	require.Equal(t, testMap[*testID], newMap[*testID])
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
	expectedErr := fmt.Sprintf("Failed to unmarshal ID: length of data "+
		"must be %d, length received is %d", ArrIDLen, 0)
	id := ID{}
	err := id.UnmarshalJSON([]byte("\"\""))
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("UnmarshalJSON failed to return the expected error for "+
			"invalid ID.\nexpected: %s\nreceived: %+v", expectedErr, err)
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
	ids := map[*ID]struct{}{}

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
		t.Errorf("NewRandomID returned an error: %+v", err)
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
	// Expected values
	expectedBytes := newRandomBytes(ArrIDLen, t)

	// Create the ID and check its contents
	newID := NewIdFromBytes(expectedBytes, t)

	if !bytes.Equal(expectedBytes, newID[:]) {
		t.Errorf("NewIdFromBytes produced an ID with the incorrect bytes."+
			"\n\texpected: %v\n\treceived: %v", expectedBytes, newID[:])
	}
}

// Tests that NewIdFromBytes creates a new ID from bytes with a length smaller
// than 33. The resulting ID should have the bytes and the rest should be 0.
func TestNewIdFromBytes_Underflow(t *testing.T) {
	// Expected values
	expectedBytes := newRandomBytes(ArrIDLen/2, t)
	expectedArr := [ArrIDLen]byte{}
	copy(expectedArr[:], expectedBytes)

	// Create the ID and check its contents
	newID := NewIdFromBytes(expectedBytes, t)

	if !bytes.Equal(expectedArr[:], newID[:]) {
		t.Errorf("NewIdFromBytes produced an ID with the incorrect bytes."+
			"\n\texpected: %v\n\treceived: %v", expectedArr, newID[:])
	}
}

// Tests that NewIdFromBytes creates a new ID from bytes with a length larger
// than 33. The resulting ID should the original bytes truncated to 33 bytes.
func TestNewIdFromBytes_Overflow(t *testing.T) {
	// Expected values
	expectedBytes := newRandomBytes(ArrIDLen*2, t)
	expectedArr := [ArrIDLen]byte{}
	copy(expectedArr[:], expectedBytes)

	// Create the ID and check its contents
	newID := NewIdFromBytes(expectedBytes, t)

	if !bytes.Equal(expectedArr[:], newID[:]) {
		t.Errorf("NewIdFromBytes produced an ID with the incorrect bytes."+
			"\n\texpected: %v\n\treceived: %v", expectedArr, newID[:])
	}
}

// Tests that NewIdFromBytes panics when given a nil testing object.
func TestNewIdFromBytes_TestError(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("NewIdFromBytes did not panic when it received a " +
				"nil testing object when it should have.")
		}
	}()

	// Call function with nil testing object
	_ = NewIdFromBytes(newRandomBytes(ArrIDLen, t), nil)
}

// Tests that NewIdFromString creates a new ID from string correctly. The new
// ID is created from a string that is 32 bytes long so that no truncation or
// padding is required. The test checks that the original string is still
// present in the data.
func TestNewIdFromString(t *testing.T) {
	// Test values
	expectedIdString := "Test ID string of correct length"
	expectedType := Generic
	expectedID := new(ID)
	copy(expectedID[:], append([]byte(expectedIdString), byte(expectedType)))

	// Create the ID and check its contents
	newID := NewIdFromString(expectedIdString, expectedType, t)

	// Check if the new ID matches the expected ID
	if !expectedID.Cmp(newID) {
		t.Errorf("NewIdFromString produced an ID with the incorrect data."+
			"\n\texpected: %v\n\treceived: %v", expectedID[:], newID[:])
	}

	// Check if the original string is still in the first 32 bytes
	newIdString := string(newID.Bytes()[:ArrIDLen-1])
	if expectedIdString != newIdString {
		t.Errorf("NewIdFromString did not correctly convert the original "+
			"string to bytes.\n\texpected string: %#v\n\treceived string: %#v"+
			"\n\texpected bytes: %v\n\treceived bytes: %v",
			expectedIdString, newIdString,
			[]byte(expectedIdString), newID.Bytes()[:ArrIDLen-1])
	}
}

// Tests that NewIdFromString panics when given a nil testing object.
func TestNewIdFromString_TestError(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("NewIdFromString did not panic when it received a " +
				"nil testing object when it should have.")
		}
	}()

	// Call function with nil testing object
	_ = NewIdFromString("test", Generic, nil)
}

// Tests that NewIdFromBase64String creates an ID, that when base 64 encoded,
// looks similar to the passed in string.
func TestNewIdFromBase64String(t *testing.T) {
	tests := []struct{ base64String, expected string }{
		{"Test 1", "Test+1AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE"},
		{"[Test  2]", "Test+2AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE"},
		{"$#%[T$%%est $#$ 3]$$%", "Test+3AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE"},
		{"$#%[T$%%est $#$ 4+/]$$%", "Test+4+/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE"},
		{"Test 55555555555555555555555555555555555555", "Test+55555555555555555555555555555555555555E"},
		{"Test 66666666666666666666666666666666666666666", "Test+66666666666666666666666666666666666666E"},
		{"", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE"},
	}

	for i, tt := range tests {
		newID := NewIdFromBase64String(tt.base64String, Group, t)

		b64 := base64.StdEncoding.EncodeToString(newID.Marshal())
		if tt.expected != b64 {
			t.Errorf("Incorrect base 64 encoding for string %q (%d)."+
				"\nexpected: %s\nreceived: %s",
				tt.base64String, i, tt.expected, b64)
		}
	}
	//
	//
	// // Test values
	// expectedIdString := "TestIDStringOfCorrectLength"
	// expectedType := Group
	// expectedID := new(ID)
	// copy(expectedID[:], append([]byte(expectedIdString), byte(expectedType)))
	//
	// strs := []string{
	// 	"test",
	// 	`"test"`,
	// 	`Question ?`,
	// 	`open   angle bracket <`,
	// 	`close angle bracket >`,
	// 	`slash /`,
	// 	`slash \`,
	// }
	// for i, str := range strs {
	// 	escaped := url.QueryEscape(str)
	// 	escaped = whitespaceRegex.ReplaceAllString(str, "+")
	// 	escaped = nonBase64Regex.ReplaceAllString(escaped, "")
	// 	fmt.Printf("%2d. %s\n    %s\n", i, str, escaped)
	// }
	//
	// // Create the ID and check its contents
	// newID := NewIdFromBase64String(expectedIdString, expectedType, t)
	//
	// // Check if the new ID matches the expected ID
	// if !expectedID.Cmp(newID) {
	// 	t.Errorf("NewIdFromString produced an ID with the incorrect data."+
	// 		"\n\texpected: %v\n\treceived: %v", expectedID[:], newID[:])
	// }
	//
	// // Check if the original string is still in the first 32 bytes
	// newIdString := string(newID.Bytes()[:ArrIDLen-1])
	// if expectedIdString != newIdString {
	// 	t.Errorf("NewIdFromString did not correctly convert the original "+
	// 		"string to bytes.\n\texpected string: %#v\n\treceived string: %#v"+
	// 		"\n\texpected bytes: %v\n\treceived bytes: %v",
	// 		expectedIdString, newIdString,
	// 		[]byte(expectedIdString), newID.Bytes()[:ArrIDLen-1])
	// }
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
			"\n\texpected: %v\n\treceived: %v", expectedUint, idUint)
	}
}

// Tests that NewIdFromUInt panics when given a nil testing object.
func TestNewIdFromUInt_TestError(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("NewIdFromUInt did not panic when it received a " +
				"nil testing object when it should have.")
		}
	}()

	// Call function with nil testing object
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
			"\n\texpected: %#v\n\treceived: %#v", expectedUints, idUints)
	}
}

// Tests that NewIdFromUInts panics when given a nil testing object.
func TestNewIdFromUInts_TestError(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("NewIdFromUInts did not panic when it received a " +
				"nil testing object when it should have.")
		}
	}()

	// Call function with nil testing object
	newUint64s := [4]uint64{rand.Uint64(), rand.Uint64(),
		rand.Uint64(), rand.Uint64()}
	_ = NewIdFromUInts(newUint64s, Generic, nil)
}

// Generates a byte slice of the specified length containing random numbers.
func newRandomBytes(length int, t *testing.T) []byte {
	// Create new byte slice of the correct size
	idBytes := make([]byte, length)

	// Create random bytes
	_, err := rand.Read(idBytes)
	if err != nil {
		t.Fatalf("Failed to generate random bytes: %v", err)
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

// Hardcoded byte array which generates a base64 string starting with a special
// character.
//
// Expected encoding is "/ABBb6YWlbkgLmg2Ohx4f0eE4K7Zx4VkGE4THx58gR8A",
// any other encoding output indicates that something is wrong (likely a
// dependency).
var hardCodedSpecialCharacter = []byte{252, 0, 65, 111, 166, 22, 149, 185, 32,
	46, 104, 54, 58, 28, 120, 127, 71, 132, 224, 174, 217, 199, 133, 100, 24,
	78, 19, 31, 30, 124, 129, 31, 189}

// Hardcoded byte array which generates a base64 string starting with an
// alphanumeric character.
//
// Expected encoding is "6iMc/s5V6MSzD6+DDQthcfA53w7wY988cenRkjxNwIcD",
// any other encoding output indicates that something is wrong (likely a
// dependency).
var hardcodedAlphaNumeric = []byte{234, 35, 28, 254, 206, 85, 232, 196, 179, 15,
	175, 131, 13, 11, 97, 113, 240, 57, 223, 14, 240, 99, 223, 60, 113, 233,
	209, 146, 60, 77, 192, 135, 19}

// Read will copy a value into byte slice p. On the first call to
// alphaNumericPRNG.Read a hardcoded value with a base 64 encoding starting with
// a special character will be returned.
//
// For any other call to alphaNumericPRNG.Read, a hardcoded value with a base 64
// encoding starting with an alphanumeric character will be returned.
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
