////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

// The package idf contains the structure of the ID File. It holds a generated
// ID and a salt (a 256-bit random number) in a JSON file. This file is used
// by different entities to save their ID and salt to file. The file path is
// usually stored and referenced from each entity's configuration files.
package idf

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/utils"
)

// The length of the salt byte array
const saltLen = 32

// IdFile structure matches the JSON structure used to save IDs and salts. The
// ID type is also saved as a string to make the file easy to read; it is never
// used or processed.
type IdFile struct {
	Salt [saltLen]byte     `json:"salt"`
	ID   [id.ArrIDLen]byte `json:"id"`
	Type string            `json:"type"`
}

// UnloadIDF reads the contents of the IDF at the given path and returns the
// salt and ID stored in it. It does so by unmarshalling the JSON in the file
// into an IdFile object. The ID bytes from the object are unmarshalled into an
// ID object and it is returned along with the salt.
//
// Errors are returned when there is a failure to read the IDF, unmarshall the
// JSON, or unmarshalling the ID.
func UnloadIDF(filePath string) ([]byte, *id.ID, error) {
	// Read the contents from the file
	jsonBytes, err := utils.ReadFile(filePath)
	if err != nil {
		return nil, nil, errors.Errorf("Could not read IDF file %s: %v",
			filePath, err)
	}

	// Create empty IdFile and unmarshal the data into it
	idf, err := newIdfFromJSON(jsonBytes)
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal ID bytes into ID
	newID, err := id.Unmarshal(idf.ID[:])

	return idf.Salt[:], newID, err
}

// Create empty IdFile, unmarshal the data into it, and return it.
func newIdfFromJSON(jsonBytes []byte) (*IdFile, error) {
	// Create new and empty IdFile object
	var idf *IdFile

	// Unmarshal JSON bytes into IdFile
	err := json.Unmarshal(jsonBytes, &idf)
	if err != nil {
		return nil, errors.Errorf("Failed to unmarshal IDF JSON: %v", err)
	}

	return idf, nil
}

// MarshalIdfToJSON creates an IdFile object with the provided values and
// marshals it into JSON bytes ready to be written to a file.
func LoadIDF(filePath string, salt []byte, genID *id.ID) error {
	// Generate new IdFile object
	idf, err := newIDF(salt, genID)
	if err != nil {
		return errors.Errorf("Failed to create new IDF: %v", err)
	}

	// Marshal the IDF into JSON bytes
	idfJSON, err := json.Marshal(idf)
	if err != nil {
		return errors.Errorf("Failed to marshal the IDF: %v", err)
	}

	// Create new ID file
	err = writeIDF(filePath, idfJSON)

	return err
}

// newIDF creates a pointer to a new IdFile object using the given ID and salt.
// The salt and marshaled ID are copied into the IdFile and the type is set from
// the ID. An error is returned if the salt is of the incorrect length.
func newIDF(salt []byte, genID *id.ID) (*IdFile, error) {
	// Check that the salt is of the correct length
	if len(salt) != saltLen {
		return nil, errors.Errorf("Salt length must be %d, length "+
			"received was %d", saltLen, len(salt))
	}

	// Create the new, empty IDF
	newIDF := &IdFile{}

	// Copy salt byte slice into IDF salt array
	copy(newIDF.Salt[:], salt)

	// Copy marshaled ID byte slice into IDF ID array
	copy(newIDF.ID[:], genID.Marshal())

	// Set the IDF type
	newIDF.Type = genID.GetType().String()

	return newIDF, nil
}

// writeIDF creates an ID file (IDF) at the given path with the given JSON data.
// Errors are returned if an error occurs making directories or files.
func writeIDF(filePath string, jsonData []byte) error {
	// Create new ID file
	err := utils.WriteFile(filePath, jsonData, utils.FilePerms, utils.DirPerms)
	if err != nil {
		return errors.Errorf("Failed to create IDF: %v", err)
	}

	return nil
}
