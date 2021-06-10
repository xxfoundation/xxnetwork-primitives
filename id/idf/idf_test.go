////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package idf

import (
	"bytes"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/utils"
	"os"
	"reflect"
	"testing"
)

// Random test values
var randomIdfJson = "{\"id\":\"Uv38ByGCZU8WP18PmmIdcpVmx00QA3xNe7sEB9HixkkC\",\"type\":\"node\",\"salt\":[133,90,216,104,29,13,134,209,233,30,0,22,121,57,203,102,148,210,196,34,172,210,8,160,7,41,57,72,127,105,153,235],\"idBytes\":[82,253,252,7,33,130,101,79,22,63,95,15,154,98,29,114,149,102,199,77,16,3,124,77,123,187,4,7,209,226,198,73,2]}"
var randomIDBytes = [id.ArrIDLen]byte{82, 253, 252, 7, 33, 130, 101, 79, 22,
	63, 95, 15, 154, 98, 29, 114, 149, 102, 199, 77, 16, 3, 124, 77, 123, 187,
	4, 7, 209, 226, 198, 73, 2}
var randomSaltBytes = [saltLen]byte{133, 90, 216, 104, 29, 13, 134, 209, 233,
	30, 0, 22, 121, 57, 203, 102, 148, 210, 196, 34, 172, 210, 8, 160, 7, 41,
	57, 72, 127, 105, 153, 235}
var randomType = "node"
var randomHexNodeId = "0x52fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c649"
var randomIDF = IdFile{
	ID:        "Uv38ByGCZU8WP18PmmIdcpVmx00QA3xNe7sEB9HixkkC",
	Type:      randomType,
	IdBytes:   randomIDBytes,
	Salt:      randomSaltBytes,
}

// Tests that newIdfFromJSON() creates the correct IdFile object form the given
// JSON bytes.
func TestNewIdfFromJSON(t *testing.T) {
	// Expected values
	expectedIDF := &randomIDF
	testIdfJSON := []byte(randomIdfJson)

	// Create the new IDF from the known JSON
	newIDF, err := newIdfFromJSON(testIdfJSON)

	// Check that no error occurred
	if err != nil {
		t.Fatalf("newIdfFromJSON() produced an unexpected error:\n%+v", err)
	}

	// Check that the new IDF matches the expected IDF
	if !reflect.DeepEqual(expectedIDF, newIDF) {
		t.Errorf("newIdfFromJSON() produced an incorrect IDF."+
			"\n\texpected: %#v\n\treceived: %#v", expectedIDF, newIDF)
	}
}

// Tests that newIdfFromJSON() creates the correct IdFile object form the given
// JSON bytes.
func TestNewIdfFromJSON_JsonError(t *testing.T) {
	// Expected values
	expectedError := "Failed to unmarshal IDF JSON: ..."
	testIdfJSON := []byte("invalidJSON")

	// Create the new IDF from the known JSON
	newIDF, err := newIdfFromJSON(testIdfJSON)

	// Check that the expected error occurred
	if err == nil {
		t.Fatalf("newIdfFromJSON() did not produce the expected error."+
			"\n\texpected: %v\n\treceived: %v", expectedError, err)
	}

	// Check that the new IDF is nil
	if newIDF != nil {
		t.Errorf("newIdfFromJSON() returned a non-nil IDF on error."+
			"\n\texpected: %#v\n\treceived: %#v", nil, newIDF)
	}
}

// Tests that newIDF() creates the correct IdFile object.
func TestNewIDF(t *testing.T) {
	// Expected values
	expectedID := id.NewIdFromBytes(randomIDBytes[:], t)
	expectedSalt := randomSaltBytes
	expectedIDF := &randomIDF

	// Create the new IDF from the expected values
	newIDF, err := newIDF(expectedSalt[:], expectedID)

	// Check that no error occurred
	if err != nil {
		t.Fatalf("newIDF() produced an unexpected error:\n%v", err)
	}

	// Check that the new IDF matches the expected IDF
	if !reflect.DeepEqual(expectedIDF, newIDF) {
		t.Errorf("newIDF() produced an incorrect IDF."+
			"\n\texpected: %#v\n\treceived: %#v", expectedIDF, newIDF)
	}
}

// Tests that newIDF() returns and error when the provided salt is of the
// incorrect length.
func TestNewIDF_SaltLengthError(t *testing.T) {
	// Expected values
	expectedID := id.NewIdFromBytes(randomIDBytes[:], t)
	expectedSalt := []byte{1, 2, 3}
	expectedError := "Salt length must be 32, length received was %d"

	// Create the new IDF from the expected values
	newIDF, err := newIDF(expectedSalt, expectedID)

	// Check that the expected error occurred
	if err == nil {
		t.Fatalf("newIDF() did not produce the expected error."+
			"\n\texpected: %v\n\treceived: %v", expectedError, err)
	}

	// Check that the new IDF is nil
	if newIDF != nil {
		t.Errorf("newIDF() returned a non-nil IDF on error."+
			"\n\texpected: %#v\n\treceived: %#v", nil, newIDF)
	}
}

// Tests that writeIDF() write the IDF data to file correctly by reading the
// file back and comparing the contents to the original JSON.
func TestWriteIDF(t *testing.T) {
	// Expected values
	expectedJSON := []byte(randomIdfJson)
	filePath := "test_ID.json"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(filePath)
		if err != nil {
			t.Fatalf("Error deleting test IDF file %#v:\n%v", filePath, err)
		}
	}()

	// Write the IDF
	err := writeIDF(filePath, expectedJSON)

	// Check that no error occurred
	if err != nil {
		t.Fatalf("writeIDF() produced an unexpected error:\n%v", err)
	}

	// Read the contents from the file
	jsonBytes, err := utils.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error reading test file %#v:\n%v", filePath, err)
	}

	// Check that the IDF file contents match the expected IDF bytes
	if !bytes.Equal(expectedJSON, jsonBytes) {
		t.Errorf("writeIDF() wrote the IDF incorrectly."+
			"\n\texpected: %#v\n\treceived: %#v", expectedJSON, jsonBytes)
	}
}

// Tests that writeIDF() returns an error and does not create a new file when
// provided a bad path.
func TestWriteIDF_BadPathError(t *testing.T) {
	// Expected values
	expectedJSON := []byte(randomIdfJson)
	filePath := "~a/test_ID.json"
	expectedError := "Failed to create IDF: ..."

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(filePath)
		if err != nil {
			t.Fatalf("Error deleting test IDF file %#v:\n%v", filePath, err)
		}
	}()

	// Write the IDF
	err := writeIDF(filePath, expectedJSON)

	// Check that the expected error occurred
	if err == nil {
		t.Fatalf("writeIDF() did not produce the expected error."+
			"\n\texpected: %v\n\treceived: %v", expectedError, err)
	}

	// Check that the IDF was not created
	if utils.Exists(filePath) {
		t.Errorf("writeIDF() created the test file %#v when given a bad path.",
			filePath)
	}
}

// Tests that UnloadIDF() returns the expected salt and IdBytes.
func TestUnloadIDF(t *testing.T) {
	// Test values
	filePath := "test_ID.json"
	testJSON := []byte(randomIdfJson)

	// Expected values
	expectedSalt := randomSaltBytes[:]
	expectedID := id.NewIdFromBytes(randomIDBytes[:], t)

	// Create IDF at path
	err := utils.WriteFile(filePath, testJSON, utils.FilePerms, utils.DirPerms)
	if err != nil {
		t.Fatalf("Error creating test IDF %#v:\n%v", filePath, err)
	}

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(filePath)
		if err != nil {
			t.Fatalf("Error deleting test IDF file %#v:\n%v", filePath, err)
		}
	}()

	// Load the IDF from file
	newSalt, newID, err := UnloadIDF(filePath)

	// Check that no error occurred
	if err != nil {
		t.Fatalf("UnloadIDF() produced an unexpected error:\n%v", err)
	}

	// Check if returned salt is correct
	if !bytes.Equal(expectedSalt, newSalt) {
		t.Errorf("UnloadIDF() returned incorrect salt."+
			"\n\texpected: %v\n\treceived: %v", expectedSalt, newSalt)
	}

	// Check if returned IdBytes is correct
	if !expectedID.Cmp(newID) {
		t.Errorf("UnloadIDF() returned incorrect IdBytes."+
			"\n\texpected: %v\n\treceived: %v",
			expectedID.Bytes(), newID.Bytes())
	}
}

// Tests that UnloadIDF() returns an error when provided an invalid path.
func TestUnloadIDF_FilePathError(t *testing.T) {
	// Test values
	filePath := "~a/test_ID.json"
	expectedError := "Could not read IDF file " + filePath + ": ..."

	// Load the IDF from file
	newSalt, newID, err := UnloadIDF(filePath)

	// Check that the expected error occurred
	if err == nil {
		t.Errorf("UnloadIDF() did not produce the expected error."+
			"\n\texpected: %v\n\treceived: %v", expectedError, err)
	}

	// Check that the returned salt is nil
	if newSalt != nil {
		t.Errorf("UnloadIDF() returned non-nil salt on error."+
			"\n\texpected: %v\n\treceived: %v", nil, newSalt)
	}

	// Check that the returned IdBytes is nil
	if newID != nil {
		t.Errorf("UnloadIDF() returned non-nil IdBytes on error."+
			"\n\texpected: %v\n\treceived: %v", nil, newID)
	}
}

// Tests that UnloadIDF() returns an error when the provided IDF contains
// invalid JSON.
func TestUnloadIDF_InvalidJsonError(t *testing.T) {
	// Test values
	filePath := "test_ID.json"
	testJSON := []byte("invalidJSON")

	// Expected values
	expectedError := "Failed to unmarshal IDF JSON: ..."

	// Create IDF at path
	err := utils.WriteFile(filePath, testJSON, utils.FilePerms, utils.DirPerms)
	if err != nil {
		t.Fatalf("Error creating test IDF %#v:\n%v", filePath, err)
	}

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(filePath)
		if err != nil {
			t.Fatalf("Error deleting test IDF file %#v:\n%v", filePath, err)
		}
	}()

	// Load the IDF from file
	newSalt, newID, err := UnloadIDF(filePath)

	// Check that the expected error occurred
	if err == nil {
		t.Errorf("UnloadIDF() did not produce the expected error."+
			"\n\texpected: %v\n\treceived: %v", expectedError, err)
	}

	// Check that the returned salt is nil
	if newSalt != nil {
		t.Errorf("UnloadIDF() returned non-nil salt on error."+
			"\n\texpected: %v\n\treceived: %v", nil, newSalt)
	}

	// Check that the returned IdBytes is nil
	if newID != nil {
		t.Errorf("UnloadIDF() returned non-nil IdBytes on error."+
			"\n\texpected: %v\n\treceived: %v", nil, newID)
	}
}

// Tests that LoadIDF() writes the correct data to file.
func TestLoadIDF(t *testing.T) {
	// Test values
	filePath := "test_ID.json"
	testSalt := randomSaltBytes[:]
	testID := id.NewIdFromBytes(randomIDBytes[:], t)

	// Expected values
	expectedIdfJSON := []byte(randomIdfJson)

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(filePath)
		if err != nil {
			t.Fatalf("Error deleting test IDF file %#v:\n%v", filePath, err)
		}
	}()

	// Load IDF into a file
	err := LoadIDF(filePath, testSalt, testID)

	// Check that no error occurred
	if err != nil {
		t.Fatalf("LoadIDF() produced an unexpected error:\n%v", err)
	}

	// Get NDF contents
	testIdfJSON, err := utils.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error reading test IDF file %#v:\n%v", filePath, err)
	}

	t.Logf("received: %s", string(testIdfJSON))

	// Check if returned IDF JSON is correct
	if !bytes.Equal(expectedIdfJSON, testIdfJSON) {
		t.Errorf("LoadIDF() created incorrect IDF."+
			"\n\texpected: %v\n\treceived: %v", expectedIdfJSON, testIdfJSON)
	}
}

// Tests that LoadIDF() returns an error when given a salt with incorrect
// length.
func TestLoadIDF_IdfError(t *testing.T) {
	// Test values
	filePath := "test_ID.json"
	testSalt := []byte{1, 2, 3}
	testID := id.NewIdFromBytes(randomIDBytes[:], t)

	// Expected values
	expectedError := "Failed to create new IDF: ..."

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(filePath)
		if err != nil {
			t.Fatalf("Error deleting test IDF file %#v:\n%v", filePath, err)
		}
	}()

	// Load IDF into a file
	err := LoadIDF(filePath, testSalt, testID)

	// Check that the expected error occurred
	if err == nil {
		t.Errorf("LoadIDF() did not produce the expected error."+
			"\n\texpected: %v\n\treceived: %v", expectedError, err)
	}

	// Check that the IDF was not created
	if utils.Exists(filePath) {
		t.Errorf("LoadIDF() created the test file %#v when given a bad path.",
			filePath)
	}
}

func TestIDF_LoadUnload(t *testing.T) {
	// Test values
	filePath := "test_ID.json"

	// Expected values
	expectedSalt := randomSaltBytes[:]
	expectedID := id.NewIdFromBytes(randomIDBytes[:], t)

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(filePath)
		if err != nil {
			t.Fatalf("Error deleting test IDF file %#v:\n%v", filePath, err)
		}
	}()

	// Load IDF into a file
	err := LoadIDF(filePath, expectedSalt, expectedID)

	// Check that no error occurred
	if err != nil {
		t.Fatalf("LoadIDF() produced an unexpected error:\n%v", err)
	}

	// Unload the IDF
	newSalt, newID, err := UnloadIDF(filePath)

	// Check that no error occurred
	if err != nil {
		t.Fatalf("UnloadIDF() produced an unexpected error:\n%v", err)
	}

	// Check if returned salt is correct
	if !bytes.Equal(expectedSalt, newSalt) {
		t.Errorf("UnloadIDF() returned incorrect salt."+
			"\n\texpected: %v\n\treceived: %v", expectedSalt, newSalt)
	}

	// Check if returned IdBytes is correct
	if !expectedID.Cmp(newID) {
		t.Errorf("UnloadIDF() returned incorrect IdBytes."+
			"\n\texpected: %v\n\treceived: %v",
			expectedID.Bytes(), newID.Bytes())
	}
}
