////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package idf

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/pkg/errors"

	"gitlab.com/xx_network/primitives/id"
)

// Happy path.
func Test_newIDF(t *testing.T) {
	expectedIDF := newRandomIDF(rand.New(rand.NewSource(687468)), t)

	idf, err := newIDF(expectedIDF.Salt[:], expectedIDF.ID)
	if err != nil {
		t.Errorf("newIDF returned an error: %+v", err)
	}

	if !reflect.DeepEqual(expectedIDF, idf) {
		t.Errorf("Unexpected IDF.\nexpected: %v\nreceived: %v", expectedIDF, idf)
	}
}

// Error path: Tests that newIDF returns the expected error when the length of
// the salt is incorrect.
func Test_newIDF_SaltLengthError(t *testing.T) {
	salt := make([]byte, saltLen*2)
	expectedErr := errors.Errorf(saltSizeErr, saltLen, len(salt))

	_, err := newIDF(salt, &id.ID{})
	if err == nil || errors.Is(err, expectedErr) {
		t.Errorf(
			"Incorrect error.\nexpected: %s\nreceived: %v", expectedErr, err)
	}
}

// Tests that an IDF that is JSON marshalled and unmarshalled matches the
// original.
func Test_JsonMarshalUnmarshal_IDF(t *testing.T) {
	expectedIDF := newRandomIDF(rand.New(rand.NewSource(88465)), t)

	data, err := json.Marshal(expectedIDF)
	if err != nil {
		t.Fatalf("Failed to JSON marshal %T: %+v", expectedIDF, err)
	}

	var idf IdFile
	if err = json.Unmarshal(data, &idf); err != nil {
		t.Fatalf("Failed to JSON unmarshal %T: %+v", idf, err)
	}

	if !reflect.DeepEqual(expectedIDF, idf) {
		t.Errorf("Unexpected IDF.\nexpected: %v\nreceived: %v", expectedIDF, idf)
	}
}

// Tests that an ID and salt saved using LoadIDF is the same unloaded by
// UnloadIDF.
func TestUnloadIDF_LoadIDF(t *testing.T) {
	filePath := "tempIDF.json"

	// Delete the test file at the end
	defer func() {
		if err := os.RemoveAll(filePath); err != nil {
			t.Fatalf("Error deleting test IDF %q: %+v", filePath, err)
		}
	}()

	idf := newRandomIDF(rand.New(rand.NewSource(55412)), t)
	err := LoadIDF(filePath, idf.Salt[:], idf.ID)
	if err != nil {
		t.Fatalf("Failed to load IDF: %+v", err)
	}

	salt, newID, err := UnloadIDF(filePath)
	if err != nil {
		t.Fatalf("Failed to unload IDF: %+v", err)
	}

	if !bytes.Equal(idf.Salt[:], salt) {
		t.Errorf("Incorrect salt.\nexpected: %X\nreceived: %X", idf.Salt, salt)
	}

	if !idf.ID.Equal(newID) {
		t.Errorf("Incorrect ID.\nexpected: %s\nreceived: %s", idf.ID, newID)
	}
}

// Error path: Tests that UnloadIDF returns the expected error for invalid JSON.
func TestUnloadIDF_InvalidJsonErr(t *testing.T) {
	filePath := "tempIDF.json"

	// Delete the test file at the end
	defer func() {
		if err := os.RemoveAll(filePath); err != nil {
			t.Fatalf("Error deleting test IDF %q: %+v", filePath, err)
		}
	}()

	err := os.WriteFile(filePath, []byte("invalid JSON"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid JSON to %q: %+v", filePath, err)
	}

	expectedErr := errors.New(unmarshalErr)

	_, _, err = UnloadIDF(filePath)
	if err == nil || errors.Is(err, expectedErr) {
		t.Errorf("Incorrect error."+
			"\nexpected: %v\nreceived: %v", expectedErr, err)
	}
}

// Error path: Tests that UnloadIDF returns the expected error for an invalid
// file path.
func TestUnloadIDF_FilePathError(t *testing.T) {
	expectedErr := errors.Errorf(ioReadErr, "")

	// Load the IDF from file
	_, _, err := UnloadIDF("")

	// Check that the expected error occurred
	if err == nil || errors.Is(err, expectedErr) {
		t.Errorf("Incorrect error."+
			"\nexpected: %v\nreceived: %v", expectedErr, err)
	}
}

// Error path: Tests that LoadIDF returns the expected error when the length of
// the salt is incorrect.
func TestLoadIDF_SaltLengthError(t *testing.T) {
	salt := make([]byte, saltLen*2)
	expectedErr := errors.Errorf(saltSizeErr, saltLen, len(salt))

	err := LoadIDF("", salt, &id.ID{})
	if err == nil || errors.Is(err, expectedErr) {
		t.Errorf(
			"Incorrect error.\nexpected: %s\nreceived: %v", expectedErr, err)
	}
}

// Error path: Tests that LoadIDF returns the expected error for an invalid file
// path.
func TestLoadIDF_FilePathError(t *testing.T) {
	idf := newRandomIDF(rand.New(rand.NewSource(6541651)), t)
	expectedErr := errors.Errorf(ioWriteErr, "")

	err := LoadIDF("", idf.Salt[:], idf.ID)
	if err == nil || errors.Is(err, expectedErr) {
		t.Errorf(
			"Incorrect error.\nexpected: %s\nreceived: %v", expectedErr, err)
	}
}

// newRandomIDF creates an IDF with randomly generated ID and salt.
func newRandomIDF(prng *rand.Rand, t *testing.T) IdFile {
	newID := id.NewRandomTestID(prng, id.Node, t)
	var salt [saltLen]byte
	prng.Read(salt[:])
	idf, err := newIDF(salt[:], newID)
	if err != nil {
		t.Fatalf("Failed to create new IDF: %+v", err)
	}
	return idf
}
