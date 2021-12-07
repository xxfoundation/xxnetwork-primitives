////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                           //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file                                                               //
////////////////////////////////////////////////////////////////////////////////

package idf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/utils"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Happy path.
func Test_newIDF(t *testing.T) {
	prng := rand.New(rand.NewSource(42))
	expectedIDF := newRandomIDF(prng, t)

	idf, err := newIDF(expectedIDF.Salt[:], expectedIDF.ID)
	if err != nil {
		t.Errorf("newIDF returned an error: %+v", err)
	}

	if !reflect.DeepEqual(expectedIDF, idf) {
		t.Errorf("newIDF did not produce the expected IDF."+
			"\nexpected: %v\nreceived: %v", expectedIDF, idf)
	}
}

// Error path: length of salt is incorrect.
func Test_newIDF_SaltLengthError(t *testing.T) {
	expectedErr := fmt.Sprintf(saltSizeErr, saltLen*2, saltLen)

	_, err := newIDF(make([]byte, saltLen*2), id.ID{})
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("newIDF did not return the expected error."+
			"\nexpected: %s\nreceived: %v", expectedErr, err)
	}
}

// Happy path.
func Test_unloadIDF(t *testing.T) {
	prng := rand.New(rand.NewSource(42))
	expectedIDF := newRandomIDF(prng, t)

	data, err := json.Marshal(expectedIDF)
	if err != nil {
		t.Errorf("Failed to JSON marshal IDF: %+v", err)
	}

	testSalt, testID, err := unloadIDF(data)
	if err != nil {
		t.Errorf("unloadIDF returned an error : %+v", err)
	}

	if testID != expectedIDF.ID {
		t.Errorf("unloadIDF returned unexpected ID."+
			"\nexpected: %s\nreceived: %s", expectedIDF.ID, testID)
	}

	if !bytes.Equal(testSalt, expectedIDF.Salt[:]) {
		t.Errorf("unloadIDF returned unexpected ID."+
			"\nexpected: %v\nreceived: %v", expectedIDF.Salt, testSalt)
	}
}

// Consistency test.
func Test_unloadIDF_Consistency(t *testing.T) {
	expectedID := id.ID{82, 253, 252, 7, 33, 130, 101, 79, 22, 63, 95, 15, 154,
		98, 29, 114, 149, 102, 199, 77, 16, 3, 124, 77, 123, 187, 4, 7, 209,
		226, 198, 73, 2}
	expectedSalt := [saltLen]byte{133, 90, 216, 104, 29, 13, 134, 209, 233, 30,
		0, 22, 121, 57, 203, 102, 148, 210, 196, 34, 172, 210, 8, 160, 7, 41,
		57, 72, 127, 105, 153, 235}
	jsonData := "{\"id\":\"Uv38ByGCZU8WP18PmmIdcpVmx00QA3xNe7sEB9HixkkC\",\"" +
		"type\":\"node\",\"salt\":[133,90,216,104,29,13,134,209,233,30,0,22," +
		"121,57,203,102,148,210,196,34,172,210,8,160,7,41,57,72,127,105,153," +
		"235],\"idBytes\":[82,253,252,7,33,130,101,79,22,63,95,15,154,98,29," +
		"114,149,102,199,77,16,3,124,77,123,187,4,7,209,226,198,73,2]}"

	testSalt, testID, err := unloadIDF([]byte(jsonData))
	if err != nil {
		t.Errorf("unloadIDF returned an error : %+v", err)
	}

	if testID != expectedID {
		t.Errorf("unloadIDF returned unexpected ID."+
			"\nexpected: %s\nreceived: %s", expectedID, testID)
	}

	if !bytes.Equal(testSalt, expectedSalt[:]) {
		t.Errorf("unloadIDF returned unexpected ID."+
			"\nexpected: %v\nreceived: %v", expectedSalt, testSalt)
	}
}

// Error path: invalid JSON data.
func Test_unloadIDF_InvalidJsonErr(t *testing.T) {
	expectedErr := strings.Split(unmarshalErr, "%")[0]

	_, _, err := unloadIDF([]byte("invalid JSON"))
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("unloadIDF did not return the expected error."+
			"\nexpected: %s\nreceived: %v", expectedErr, err)
	}
}

// Happy path.
func TestUnloadIDF(t *testing.T) {
	prng := rand.New(rand.NewSource(42))
	filePath := fmt.Sprintf("%s-%d.json", "testIDF", time.Now().Unix())
	expectedIDF := newRandomIDF(prng, t)

	data, err := json.Marshal(expectedIDF)
	if err != nil {
		t.Errorf("Failed to JSON marshal IDF: %+v", err)
	}

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(filePath)
		if err != nil {
			t.Fatalf("Error deleting test IDF file %s: %+v", filePath, err)
		}
	}()

	// Create IDF at path
	err = utils.WriteFile(filePath, data, utils.FilePerms, utils.DirPerms)
	if err != nil {
		t.Fatalf("Error creating test IDF %s: %+v", filePath, err)
	}

	// Load the IDF from file
	newSalt, newID, err := UnloadIDF(filePath)

	// Check that no error occurred
	if err != nil {
		t.Fatalf("UnloadIDF produced an unexpected error: %+v", err)
	}

	// Check if returned salt is correct
	if !bytes.Equal(expectedIDF.Salt[:], newSalt) {
		t.Errorf("UnloadIDF returned incorrect salt."+
			"\nexpected: %v\nreceived: %v", expectedIDF.Salt, newSalt)
	}

	// Check if returned IdBytes is correct
	if expectedIDF.ID != newID {
		t.Errorf("UnloadIDF returned incorrect ID."+
			"\nexpected: %s\nreceived: %s", expectedIDF.ID, newID)
	}
}

// Error path: provided path does not exist.
func TestUnloadIDF_FilePathError(t *testing.T) {
	expectedErr := strings.Split(ioReadErr, "%")[0]

	// Load the IDF from file
	_, _, err := UnloadIDF("")

	// Check that the expected error occurred
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("UnloadIDF did not produce the expected error."+
			"\nexpected: %s\nreceived: %v", expectedErr, err)
	}
}

// Happy path.
func Test_loadIDF(t *testing.T) {
	prng := rand.New(rand.NewSource(42))
	idf := newRandomIDF(prng, t)
	expectedData, err := json.Marshal(idf)
	if err != nil {
		t.Errorf("Failed to JSON marshal IDF: %+v", err)
	}

	testData, err := loadIDF(idf.Salt[:], idf.ID)
	if err != nil {
		t.Errorf("loadIDF returned an error: %+v", err)
	}

	if !bytes.Equal(expectedData, testData) {
		t.Errorf("loadIDF returned unexpected JSON data."+
			"\nexpected: %s\nreceived: %s", expectedData, testData)
	}
}

// Consistency test.
func Test_loadIDF_Consistency(t *testing.T) {
	prng := rand.New(rand.NewSource(42))
	idf := newRandomIDF(prng, t)
	expectedData := []byte(
		"{\"id\":\"39ebTXZCm2F6DJ+fDTulWwzA1hRMiIU1hBrL4HCbB1gC\",\"type\":\"" +
			"node\",\"salt\":[83,140,127,150,177,100,191,27,151,187,159,75," +
			"180,114,232,159,91,20,132,242,82,9,201,217,52,62,146,186,9,221," +
			"157,82],\"idBytes\":[223,215,155,77,118,66,155,97,122,12,159," +
			"159,13,59,165,91,12,192,214,20,76,136,133,53,132,26,203,224," +
			"112,155,7,88,2],\"hexNodeID\":\"0xdfd79b4d76429b617a0c9f9f0d3" +
			"ba55b0cc0d6144c888535841acbe0709b0758\"}")

	testData, err := loadIDF(idf.Salt[:], idf.ID)
	if err != nil {
		t.Errorf("loadIDF returned an error: %+v", err)
	}

	if !bytes.Equal(expectedData, testData) {
		t.Errorf("loadIDF returned unexpected JSON data."+
			"\nexpected: %s\nreceived: %s", expectedData, testData)
	}
}

// Error path: salt is of the wrong length.
func Test_loadIDF_SaltLengthError(t *testing.T) {
	expectedErr := fmt.Sprintf(saltSizeErr, saltLen*2, saltLen)

	_, err := loadIDF(make([]byte, saltLen*2), id.ID{})
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("loadIDF did not return the expected error."+
			"\nexpected: %s\nreceived: %v", expectedErr, err)
	}
}

// Happy path.
func TestLoadIDF(t *testing.T) {
	prng := rand.New(rand.NewSource(42))
	filePath := fmt.Sprintf("%s-%d.json", "testIDF", time.Now().Unix())
	expectedIDF := newRandomIDF(prng, t)

	expectedData, err := json.Marshal(expectedIDF)
	if err != nil {
		t.Errorf("Failed to JSON marshal IDF: %+v", err)
	}

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(filePath)
		if err != nil {
			t.Fatalf("Error deleting test IDF file %s: %+v", filePath, err)
		}
	}()

	// Load IDF into a file
	err = LoadIDF(filePath, expectedIDF.Salt[:], expectedIDF.ID)
	if err != nil {
		t.Fatalf("LoadIDF produced an error: %+v", err)
	}

	// Get NDF contents
	testData, err := utils.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error reading test IDF file %s: %+v", filePath, err)
	}

	// Check if returned IDF JSON is correct
	if !bytes.Equal(expectedData, testData) {
		t.Errorf("LoadIDF created incorrect IDF."+
			"\nexpected: %s\nreceived: %s", expectedData, testData)
	}
}

// Error path: provided path does not exist.
func TestLoadIDF_SaltLengthError(t *testing.T) {
	expectedErr := fmt.Sprintf(saltSizeErr, saltLen*2, saltLen)

	// Load IDF into a file
	err := LoadIDF("", make([]byte, saltLen*2), id.ID{})

	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("LoadIDF did not produce the expected error."+
			"\nexpected: %s\nreceived: %v", expectedErr, err)
	}
}

// Error path: provided path does not exist.
func TestLoadIDF_FilePathError(t *testing.T) {
	prng := rand.New(rand.NewSource(42))
	expectedIDF := newRandomIDF(prng, t)
	expectedErr := strings.Split(writeErr, "%")[0]

	// Load IDF into a file
	err := LoadIDF("", expectedIDF.Salt[:], expectedIDF.ID)

	// Check that the expected error occurred
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("LoadIDF did not produce the expected error."+
			"\nexpected: %s\nreceived: %v", expectedErr, err)
	}
}

// Checks that an IdFile can be loaded and unloaded correctly.
func TestIDF_LoadUnload(t *testing.T) {
	prng := rand.New(rand.NewSource(42))
	expectedIDF := newRandomIDF(prng, t)
	filePath := fmt.Sprintf("%s-%d.json", "testIDF", time.Now().Unix())

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(filePath)
		if err != nil {
			t.Fatalf("Error deleting test IDF file %s: %+v", filePath, err)
		}
	}()

	// Load IDF into a file
	err := LoadIDF(filePath, expectedIDF.Salt[:], expectedIDF.ID)
	if err != nil {
		t.Fatalf("LoadIDF produced an error: %+v", err)
	}

	// Unload the IDF
	newSalt, newID, err := UnloadIDF(filePath)
	if err != nil {
		t.Fatalf("UnloadIDF produced an error: %+v", err)
	}

	// Check if returned salt is correct
	if !bytes.Equal(expectedIDF.Salt[:], newSalt) {
		t.Errorf("UnloadIDF returned incorrect salt."+
			"\nexpected: %v\nreceived: %v", expectedIDF.Salt, newSalt)
	}

	// Check if returned ID is correct
	if expectedIDF.ID != newID {
		t.Errorf("UnloadIDF returned incorrect ID."+
			"\nexpected: %s\nreceived: %s", expectedIDF.ID, newID)
	}
}

func newRandomIDF(prng *rand.Rand, t *testing.T) IdFile {
	expectedID := id.NewIdFromString("myID", id.Node, t)
	expectedSalt := [saltLen]byte{}
	prng.Read(expectedSalt[:])
	prng.Read(expectedID[:id.ArrIDLen-1])
	idf, err := newIDF(expectedSalt[:], expectedID)
	if err != nil {
		t.Errorf("Failed to create new IDF: %+v", err)
	}
	return idf
}
