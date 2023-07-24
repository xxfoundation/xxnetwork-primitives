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

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

// Tests that ID.Marshal returns the correct marshaled byte slice of the ID and
// that the returned data is a copy.
func TestID_Marshal(t *testing.T) {
	prng := rand.New(rand.NewSource(498651))
	id := NewRandomTestID(prng, Group, t)
	expected := id.Bytes()

	// Check for the correct values
	idBytes := id.Marshal()
	if !bytes.Equal(expected, idBytes) {
		t.Errorf("Marshal returned unexpected bytes."+
			"\nexpected: %+v\nreceived: %+v", expected, idBytes)
	}

	// Test if the returned bytes are copies
	if &id[0] == &idBytes[0] {
		t.Errorf("Marshal did not return a copy of the ID data."+
			"\nID pointer:    %+v\nbytes pointer: %+v", &id[0], &idBytes[0])
	}
}

// Tests that Unmarshal creates a new ID with the expected data.
func TestUnmarshal(t *testing.T) {
	expected := rngBytes(ArrIDLen, 42, t)

	// Unmarshal the bytes into an ID
	id, err := Unmarshal(expected)
	if err != nil {
		t.Errorf("Unmarshal produced an error: %+v", err)
	}

	// Make sure the ID contents are correct
	if !bytes.Equal(expected, id[:]) {
		t.Errorf("Unmarshal produced an ID with the incorrect bytes."+
			"\nexpected: %v\nreceived: %v", expected, id[:])
	}
}

// Error path: Tests that Unmarshal produces an error when the given data length
// is not equal to the length of an ID.
func TestUnmarshal_DataLengthError(t *testing.T) {
	invalidIdBytes := rngBytes(ArrIDLen+10, 42, t)
	expectedErr := fmt.Errorf("could not marshal byte slice to ID: "+
		"length of data must be %d, length received was %d",
		ArrIDLen, len(invalidIdBytes))

	// Unmarshal the bytes into an ID
	_, err := Unmarshal(invalidIdBytes)
	if err == nil {
		t.Errorf("Unmarshal did not product an expected error."+
			"\nexpected: %v\nreceived: %v", expectedErr, err)
	}
}

// Test that an ID that is marshaled and unmarshalled matches the original.
func TestID_Marshal_Unmarshal(t *testing.T) {
	originalID := NewRandomTestID(rand.New(rand.NewSource(89)), Node, t)

	idBytes := originalID.Marshal()

	id, err := Unmarshal(idBytes)
	if err != nil {
		t.Errorf("Unmarshal produced an error: %+v", err)
	}

	if !originalID.Equal(id) {
		t.Errorf("Original ID does not match marshaled/unmarshalled ID."+
			"\nexpected: %s\nreceived: %s", originalID, id)
	}
}

// Tests that the byte slice returned by ID.Bytes matches the data in the
// original ID and the data is a copy of the values and not the reference.
func TestID_Bytes(t *testing.T) {
	expected := rngBytes(ArrIDLen, 42, t)
	id := NewIdFromBytes(expected, t)

	// Check for the correct values
	idBytes := id.Bytes()
	if !bytes.Equal(expected, idBytes) {
		t.Errorf("Bytes returned unexpected bytes."+
			"\nexpected: %v\nreceived: %v", expected, idBytes)
	}

	// Test if the returned bytes are copies
	if &id[0] == &idBytes[0] {
		t.Errorf("Bytes did not return a copy of the ID data."+
			"\nID pointer:    %+v\nbytes pointer: %+v", &id[0], &idBytes[0])
	}
}

// Tests that ID.Bytes panics when the ID is nil.
func TestID_Bytes_NilError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Bytes failed to panic when the ID is nil.")
		}
	}()

	var id *ID
	_ = id.Bytes()
}

// Tests that ID.Equal properly reports when two IDs with the same values are
// equal and different values as not equal.
func TestID_Equal(t *testing.T) {
	prng := rand.New(rand.NewSource(8459612))

	idA := NewRandomTestID(prng, Gateway, t)
	idB, idC := NewIdFromBytes(idA.Bytes(), t), idA.DeepCopy()
	idC.SetType(Node)
	tests := []struct {
		a, b  *ID
		equal bool
	}{
		{idA, idA, true},
		{idA, idB, true},
		{idA, idC, false},
		{idA, NewRandomTestID(prng, Gateway, t), false},
		{NewRandomTestID(prng, Node, t), NewRandomTestID(prng, Node, t), false},
	}

	for i, tt := range tests {
		equal := tt.a.Equal(tt.b)
		if equal != tt.equal {
			t.Errorf("Incorrect result for %s == %s (%d)."+
				"\nexpected: %t\nreceived: %t", tt.a, tt.b, i, tt.equal, equal)
		}
	}
}

// Tests that ID.Equal panics when the ID is nil.
func TestID_Equal_NilError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Equal failed to panic when the ID is nil.")
		}
	}()

	var id *ID
	_ = id.Equal(nil)
}

// Tests that ID.Compare properly reports the comparison of different IDs.
func TestID_Compare(t *testing.T) {
	tests := []struct {
		a, b    *ID
		compare int
	}{
		{NewIdFromBytes([]byte(""), t), NewIdFromBytes([]byte(""), t), 0},
		{NewIdFromBytes([]byte("a"), t), NewIdFromBytes([]byte(""), t), 1},
		{NewIdFromBytes([]byte(""), t), NewIdFromBytes([]byte("a"), t), -1},
		{NewIdFromBytes([]byte("abc"), t), NewIdFromBytes([]byte("abc"), t), 0},
		{NewIdFromBytes([]byte("abd"), t), NewIdFromBytes([]byte("abc"), t), 1},
		{NewIdFromBytes([]byte("abc"), t), NewIdFromBytes([]byte("abd"), t), -1},
		{NewIdFromBytes([]byte("ab"), t), NewIdFromBytes([]byte("abc"), t), -1},
		{NewIdFromBytes([]byte("abc"), t), NewIdFromBytes([]byte("ab"), t), 1},
		{NewIdFromBytes([]byte("x"), t), NewIdFromBytes([]byte("ab"), t), 1},
		{NewIdFromBytes([]byte("ab"), t), NewIdFromBytes([]byte("x"), t), -1},
		{NewIdFromBytes([]byte("x"), t), NewIdFromBytes([]byte("a"), t), 1},
		{NewIdFromBytes([]byte("b"), t), NewIdFromBytes([]byte("x"), t), -1},
	}

	for i, tt := range tests {
		compare := tt.a.Compare(tt.b)
		if compare != tt.compare {
			t.Errorf("Incorrect result for comparing %s and %s (%d)."+
				"\nexpected: %d\nreceived: %d",
				tt.a, tt.b, i, tt.compare, compare)
		}
	}
}

// Tests that ID.Compare panics when the ID is nil.
func TestID_Compare_NilError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Compare failed to panic when the ID is nil.")
		}
	}()

	var id *ID
	_ = id.Compare(nil)
}

// Tests that ID.Less properly reports when an ID is less than another.
func TestID_Less(t *testing.T) {
	tests := []struct {
		a, b *ID
		less bool
	}{
		{NewIdFromBytes([]byte(""), t), NewIdFromBytes([]byte(""), t), false},
		{NewIdFromBytes([]byte("a"), t), NewIdFromBytes([]byte(""), t), false},
		{NewIdFromBytes([]byte(""), t), NewIdFromBytes([]byte("a"), t), true},
		{NewIdFromBytes([]byte("abc"), t), NewIdFromBytes([]byte("abc"), t), false},
		{NewIdFromBytes([]byte("abd"), t), NewIdFromBytes([]byte("abc"), t), false},
		{NewIdFromBytes([]byte("abc"), t), NewIdFromBytes([]byte("abd"), t), true},
		{NewIdFromBytes([]byte("ab"), t), NewIdFromBytes([]byte("abc"), t), true},
		{NewIdFromBytes([]byte("abc"), t), NewIdFromBytes([]byte("ab"), t), false},
		{NewIdFromBytes([]byte("x"), t), NewIdFromBytes([]byte("ab"), t), false},
		{NewIdFromBytes([]byte("ab"), t), NewIdFromBytes([]byte("x"), t), true},
		{NewIdFromBytes([]byte("x"), t), NewIdFromBytes([]byte("a"), t), false},
		{NewIdFromBytes([]byte("b"), t), NewIdFromBytes([]byte("x"), t), true},
	}

	for i, tt := range tests {
		less := tt.a.Less(tt.b)
		if less != tt.less {
			t.Errorf("Incorrect result for %s < %s (%d)."+
				"\nexpected: %t\nreceived: %t", tt.a, tt.b, i, tt.less, less)
		}
	}
}

// Test that DeepCopy returns a copy with the same contents as the original
// and where the pointers are different.
func TestID_DeepCopy(t *testing.T) {
	// Test values
	expectedID := NewRandomTestID(rand.New(rand.NewSource(2125)), Gateway, t)

	// Test if the contents are equal
	testVal := expectedID.DeepCopy()
	if !reflect.DeepEqual(expectedID, testVal) {
		t.Errorf("DeepCopy returned a copy with the wrong contents."+
			"\nexpected: %+v\nreceived: %+v", expectedID, testVal)
	}

	// Test if the returned bytes are copies
	if &expectedID[0] == &testVal[0] {
		t.Errorf("DeepCopy did not return a copy when it should have."+
			"\nexpected: any value except %+v\nreceived: %+v",
			&expectedID[0], &testVal[0])
	}
}

// Tests that ID.DeepCopy panics when the ID is nil.
func TestID_DeepCopy_NilError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("DeepCopy failed to panic when the ID is nil.")
		}
	}()

	var id *ID
	_ = id.DeepCopy()
}

// Consistency test of ID.String.
func TestID_String(t *testing.T) {
	prng := rand.New(rand.NewSource(968541))
	expectedIDs := []string{
		"7cDhuqmCtFaDidi2WBsoMjjzRly6uO4DgR0PbwG8Od8A",
		"9EP/FtTKdaIGW2zBz0dX5/h2jEiF3UtmJoebiZt1oXEB",
		"3Inxc3Kxl/qPTNTkr24WslhIizm//zMKV8+/1Rr0wiIA",
		"JF3qunQHJsa0ZUxiTvNO6xrH+fB9ZiWjESlvHjSnX58D",
		"lERWkU9BKJQu1ZSrIftW+7X+7Zaxbry0f1qckIk8L5AA",
		"LibKwFbfLVcIbJFgmEhNFxhrNOmsJgEDTx9dTfIuL+oE",
		"BQo5PteLs/vejskCcFboH/rAJJrYm/CUkIVKw9WLA0AB",
		"7F+ouuu1drGB7cH1fK+6p8EG9Kps/iyRX1YU5V9PBk8E",
		"0EDjLA2F+PSbhsMkvK57x0S+u1JPAiFAGOIyu3M5wVwA",
		"p823UKuHN9s0Q+3eIkndKDJ3GcHDrgaBZV7xQJbcR/AD",
	}

	for i, expected := range expectedIDs {
		id := NewRandomTestID(prng, Type(prng.Intn(int(NumTypes))), t)

		if id.String() != expected {
			t.Errorf("String did not output the expected value for ID %d."+
				"\nexpected: %s\nreceived: %s", i, expected, id.String())
		}
	}
}

// Tests that ID.String panics when the ID is nil.
func TestID_String_NilError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("String failed to panic when the ID is nil.")
		}
	}()

	var id *ID
	_ = id.String()
}

// Tests that ID.GetType returns the correct type for each ID type.
func TestID_GetType(t *testing.T) {
	prng := rand.New(rand.NewSource(334832))
	testTypes := []Type{Generic, Gateway, Node, User, Group, NumTypes, 7}
	testIDs := make([]*ID, len(testTypes))
	for i, idType := range testTypes {
		testIDs[i] = NewRandomTestID(prng, idType, t)
	}

	for i, id := range testIDs {
		if testTypes[i] != id.GetType() {
			t.Errorf("GetType returned the incorrect type (%d)."+
				"\nexpected: %s\nreceived: %s", i, testTypes[i], id.GetType())
		}
	}
}

// Tests that ID.GetType panics when the ID is nil.
func TestID_GetType_NilError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("GetType failed to panic when the ID is nil.")
		}
	}()

	var id *ID
	_ = id.GetType()
}

// Tests that ID.SetType sets the type of the ID correctly by checking if the
// ID's type changed after calling SetType.
func TestID_SetType(t *testing.T) {
	prng := rand.New(rand.NewSource(334832))
	testTypes := []Type{Generic, Gateway, Node, User, Group, NumTypes, 7}
	for i, idType := range testTypes {
		id := NewRandomTestID(prng, Type(prng.Intn(int(NumTypes))), t)
		id.SetType(idType)

		if idType != id.GetType() {
			t.Errorf("Incorrect type for ID %s (%d)."+
				"\nexpected: %s\nreceived: %s", id, i, idType, id.GetType())
		}
	}
}

// Tests that ID.SetType panics when the ID is nil.
func TestID_SetType_NilError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("SetType failed to panic when the ID is nil.")
		}
	}()

	var id *ID
	id.SetType(Generic)
}

// Tests that an ID can be JSON marshaled and unmarshalled.
func TestID_MarshalJSON_UnmarshalJSON(t *testing.T) {
	id := NewRandomTestID(rand.New(rand.NewSource(49056)), Node, t)

	jsonData, err := json.Marshal(id)
	if err != nil {
		t.Errorf("json.Marshal returned an error: %+v", err)
	}

	newID := &ID{}
	err = json.Unmarshal(jsonData, newID)
	if err != nil {
		t.Errorf("json.Unmarshal returned an error: %+v", err)
	}

	if !id.Equal(newID) {
		t.Errorf("Failed the JSON marshal and unmarshal ID."+
			"\noriginal ID: %s\nreceived ID: %s", id, newID)
	}
}

// Tests that an ID can be JSON marshaled and unmarshalled when it is the key in
// a map. This tests ID.MarshalText.
func TestID_TextMarshal(t *testing.T) {
	id := NewRandomTestID(rand.New(rand.NewSource(42534)), User, t)

	testMap := make(map[ID]int)
	testMap[*id] = 8675309

	jsonData, err := json.Marshal(testMap)
	require.NoError(t, err)

	newMap := make(map[ID]int)
	err = json.Unmarshal(jsonData, &newMap)
	require.NoError(t, err)

	require.Equal(t, testMap[*id], newMap[*id])
}

// Tests that an ID can be text marshaled and unmarshalled.
func TestID_MarshalText_UnmarshalText(t *testing.T) {
	id := NewRandomTestID(rand.New(rand.NewSource(6156)), Node, t)

	text, err := id.MarshalText()
	if err != nil {
		t.Errorf("Failed to text marshal: %+v", err)
	}

	newID := &ID{}
	err = newID.UnmarshalText(text)
	if err != nil {
		t.Errorf("Failed to text unmarshal: %+v", err)
	}

	if !id.Equal(newID) {
		t.Errorf("Text marshalled and unmarshalled ID does not match original."+
			"\nexpected: %s\nreceived: %s", id, newID)
	}
}

// Error path: Tests that ID.UnmarshalText returns an error when the bytes are
// not a valid base 64 string.
func TestID_UnmarshalText_InvalidBase64Error(t *testing.T) {
	expectedErr := base64.CorruptInputError(7)

	id := &ID{}
	err := id.UnmarshalText([]byte("invalid bytes"))
	if err == nil || !errors.Is(err, expectedErr) {
		t.Errorf("Failed to receive expected error.\nexpected: %v\nreceived: %+v",
			expectedErr, err)
	}
}

// Error path: Tests that ID.UnmarshalText returns an error when the bytes are
// not a valid ID.
func TestID_UnmarshalText_InvalidIDError(t *testing.T) {
	data := []byte("InvalidID")
	expectedErr := errors.Errorf("Failed to unmarshal ID: length of data "+
		"must be %d, length received is %d", ArrIDLen, len(data))

	id := &ID{}
	err := id.UnmarshalText([]byte(base64.StdEncoding.EncodeToString(data)))
	if err == nil || !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Errorf("Failed to receive expected error.\nexpected: %v\nreceived: %+v",
			expectedErr, err)
	}
}

// Error path: supplied data is invalid JSON.
func TestID_UnmarshalJSON_JsonUnmarshalError(t *testing.T) {
	expectedErr := "invalid character"
	id := &ID{}
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
	id := &ID{}
	err := id.UnmarshalJSON([]byte("\"\""))
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("UnmarshalJSON failed to return the expected error for "+
			"invalid ID.\nexpected: %s\nreceived: %+v", expectedErr, err)
	}
}

// Consistency test of ID.HexEncode.
func TestID_HexEncode(t *testing.T) {
	prng := rand.New(rand.NewSource(6516))
	expectedIDs := []string{
		"0xae5843e77c9721df5ba243cff923d518d00ced010b0105c145dfae308b900d60",
		"0x0ca36fad72ac9a0032e84d8980377841f66ec5a61311adc85f588100d7b40ef5",
		"0x72567246a346927f4b672dbf2ef5f0a310f1cd1e81b625c37762860b010f6ec7",
		"0x18ec3e065d83db8c8eaee3b3434a1e551dbeda262b5babbd15ef0067f6aba57b",
		"0x5b1ebe5390064991e9f9f50395e98b89336e2901decf89cfc87c4615f3aa847e",
		"0xdeed3aa17f2fec2a8f400c713d5833873bfb39175bb126fee6317a43c2801b6e",
		"0x7946861001bbfb4c50c300cb13bf912427ea5d0a8dc9ee58d975697782cdc0d1",
		"0x75497ed22411fa8467f73d140e83b866d0e01d7052a8f40fad30ed77b643c88e",
		"0x34369f5ba85469f0fdf42e23df6eeca302917b85af084deaab7b27b8b3927ea9",
		"0xc8eb3ca60b689be3588bf254627086a2f678a30d79b9d24c813202894cda8113",
	}

	for i, expected := range expectedIDs {
		id := NewRandomTestID(prng, Type(prng.Intn(int(NumTypes))), t)

		hexString := id.HexEncode()
		if hexString != expected {
			t.Errorf("Unexpected hex string for ID %s (%d)."+
				"\nexpected: %s\nreceived: %s", id, i, expected, hexString)
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
		id, err := NewRandomID(prng, Type(prng.Intn(int(NumTypes))))
		if err != nil {
			t.Errorf("NewRandomID returned an error (%d): %+v", i, err)
		}
		if id.String() != expected {
			t.Errorf("NewRandomID did not generate the expected ID (%d)."+
				"\nexpected: %s\nreceived: %s", i, expected, id.String())
		}
	}
}

// Tests that NewRandomID returns unique IDs.
func TestNewRandomID_Unique(t *testing.T) {
	prng := rand.New(rand.NewSource(42))
	ids := map[*ID]struct{}{}

	for i := 0; i < 100; i++ {
		id, err := NewRandomID(prng, Type(prng.Intn(int(NumTypes))))
		if err != nil {
			t.Errorf("NewRandomID returned an error (%d): %+v", i, err)
		}

		if _, exists := ids[id]; exists {
			t.Errorf(
				"NewRandomID did not generate a unique ID (%d).\nID: %s", i, id)
		} else {
			ids[id] = struct{}{}
		}
	}
}

// This tests uses a custom PRNG which generates an ID with a base 64 encoding
// starting with a special character. NewRandomID should force another call
// to prng.Read, and this call should return an ID with encoding starting with
// an alphanumeric character. This test fails if it hangs forever (PRNG error)
// or the ID has an encoding beginning with a special character (NewRandomID
// error).
func TestNewRandomID_SpecialCharacter(t *testing.T) {
	prng := newAlphanumericPRNG()
	id, err := NewRandomID(prng, 0)
	if err != nil {
		t.Errorf("NewRandomID returned an error: %+v", err)
	}

	if !regexAlphanumeric.MatchString(string(id.String()[0])) {
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

// Tests that NewRandomTestID returns the expected IDs for a given PRNG.
func TestNewRandomTestID_Consistency(t *testing.T) {
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
		id := NewRandomTestID(prng, Type(prng.Intn(int(NumTypes))), t)

		if id.String() != expected {
			t.Errorf("NewRandomTestID did not generate the expected ID (%d)."+
				"\nexpected: %s\nreceived: %s", i, expected, id.String())
		}
	}
}

// Tests that NewRandomTestID panics when given a nil testing object.
func TestNewRandomTestID_TestError(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("NewRandomTestID did not panic when it received a nil " +
				"testing object when it should have.")
		}
	}()

	// Call the function with nil testing object
	_ = NewRandomTestID(nil, Generic, nil)
}

// Tests that NewRandomTestID returns an error when the io reader encounters an
// error.
func TestNewRandomTestID_ReaderError(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("NewRandomTestID failed to return an error when the " +
				"reader failed.")
		}
	}()

	_ = NewRandomTestID(strings.NewReader(""), Generic, t)
}

// Tests that NewIdFromBytes creates a new ID with the correct contents.
func TestNewIdFromBytes(t *testing.T) {
	expectedBytes := rngBytes(ArrIDLen, 42, t)
	id := NewIdFromBytes(expectedBytes, t)

	if !bytes.Equal(expectedBytes, id[:]) {
		t.Errorf("NewIdFromBytes produced an ID with the incorrect bytes."+
			"\nexpected: %v\nreceived: %v", expectedBytes, id[:])
	}
}

// Tests that NewIdFromBytes creates a new ID from bytes with a length smaller
// than 33. The resulting ID should have the bytes and the rest should be 0.
func TestNewIdFromBytes_Underflow(t *testing.T) {
	prng := rand.New(rand.NewSource(65474))

	// Expected values
	expectedBytes := make([]byte, ArrIDLen/2)
	prng.Read(expectedBytes)
	expectedArr := [ArrIDLen]byte{}
	copy(expectedArr[:], expectedBytes)

	// Create the ID and check its contents
	id := NewIdFromBytes(expectedBytes, t)

	if !bytes.Equal(expectedArr[:], id[:]) {
		t.Errorf("NewIdFromBytes produced an ID with the incorrect bytes."+
			"\nexpected: %v\nreceived: %v", expectedArr, id[:])
	}
}

// Tests that NewIdFromBytes creates a new ID from bytes with a length larger
// than 33. The resulting ID should the original bytes truncated to 33 bytes.
func TestNewIdFromBytes_Overflow(t *testing.T) {
	prng := rand.New(rand.NewSource(22445))

	// Expected values
	expectedBytes := make([]byte, ArrIDLen*2)
	prng.Read(expectedBytes)
	expectedArr := [ArrIDLen]byte{}
	copy(expectedArr[:], expectedBytes)

	// Create the ID and check its contents
	id := NewIdFromBytes(expectedBytes, t)

	if !bytes.Equal(expectedArr[:], id[:]) {
		t.Errorf("NewIdFromBytes produced an ID with the incorrect bytes."+
			"\nexpected: %v\nreceived: %v", expectedArr, id[:])
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
	id := NewIdFromString(expectedIdString, expectedType, t)

	// Check if the new ID matches the expected ID
	if !expectedID.Equal(id) {
		t.Errorf("NewIdFromString produced an ID with the incorrect data."+
			"\nexpected: %v\nreceived: %v", expectedID[:], id[:])
	}

	// Check if the original string is still in the first 32 bytes
	idString := string(id.Bytes()[:ArrIDLen-1])
	if expectedIdString != idString {
		t.Errorf("NewIdFromString did not correctly convert the original "+
			"string to bytes.\nexpected string: %q\nreceived string: %q"+
			"\nexpected bytes: %v\nreceived bytes: %v", expectedIdString,
			idString, []byte(expectedIdString), id.Bytes()[:ArrIDLen-1])
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

// Tests that NewIdFromBase64String creates an ID, that when base 64 encoded,
// looks similar to the passed in string.
func TestNewIdFromBase64String(t *testing.T) {
	tests := []struct{ base64String, expected string }{
		{"Test 1",
			"Test+1AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE"},
		{"[Test  2]",
			"Test+2AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE"},
		{"$#%[T$%%est $#$ 3]$$%",
			"Test+3AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE"},
		{"$#%[T$%%est $#$ 4+/]$$%",
			"Test+4+/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE"},
		{"Test 55555555555555555555555555555555555555",
			"Test+55555555555555555555555555555555555555E"},
		{"Test 66666666666666666666666666666666666666666",
			"Test+66666666666666666666666666666666666666E"},
		{"", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE"},
	}

	for i, tt := range tests {
		id := NewIdFromBase64String(tt.base64String, Group, t)

		b64 := base64.StdEncoding.EncodeToString(id.Marshal())
		if tt.expected != b64 {
			t.Errorf("Incorrect base 64 encoding for string %q (%d)."+
				"\nexpected: %s\nreceived: %s",
				tt.base64String, i, tt.expected, b64)
		}
	}
}

// Tests that NewIdFromBase64String panics when given a nil testing object.
func TestNewIdFromBase64String_TestError(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("NewIdFromBase64String did not panic when it received a " +
				"nil testing object when it should have.")
		}
	}()

	// Call the function with nil testing object
	_ = NewIdFromBase64String("", Generic, nil)
}

// Tests that NewIdFromUInt creates a new ID with the correct contents by
// converting the ID back into a uint and comparing it to the original.
func TestNewIdFromUInt(t *testing.T) {
	// Expected values
	expectedUint := rand.Uint64()

	// Create the ID and check its contents
	id := NewIdFromUInt(expectedUint, Generic, t)
	idUint := binary.BigEndian.Uint64(id[:ArrIDLen-1])

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
	expectedUints := [4]uint64{
		rand.Uint64(), rand.Uint64(), rand.Uint64(), rand.Uint64()}

	// Create the ID and check its contents
	id := NewIdFromUInts(expectedUints, Generic, t)
	idUints := [4]uint64{}
	for i := range idUints {
		idUints[i] = binary.BigEndian.Uint64(id[i*8 : (i+1)*8])
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

// alphaNumericPRNG is a custom PRNG which adheres to the io.Reader interface.
// This is used for testing special characters.
type alphaNumericPRNG struct {
	counter uint
}

// newAlphanumericPRNG generates a copy of the custom prng alphaNumericPRNG.
func newAlphanumericPRNG() *alphaNumericPRNG {
	return &alphaNumericPRNG{counter: 0}
}

// Hardcoded byte array that generates a base64 string starting with a special
// character.
//
// Expected encoding is "/ABBb6YWlbkgLmg2Ohx4f0eE4K7Zx4VkGE4THx58gR8A".
// Any other encoding output indicates that something is wrong (likely a
// dependency).
var hardCodedSpecialCharacter = []byte{252, 0, 65, 111, 166, 22, 149, 185, 32,
	46, 104, 54, 58, 28, 120, 127, 71, 132, 224, 174, 217, 199, 133, 100, 24,
	78, 19, 31, 30, 124, 129, 31, 189}

// Hardcoded byte array that generates a base64 string starting with an
// alphanumeric character.
//
// Expected encoding is "6iMc/s5V6MSzD6+DDQthcfA53w7wY988cenRkjxNwIcD".
// Any other encoding output indicates that something is wrong (likely a
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
