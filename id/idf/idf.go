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
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/utils"
)

// The length of the salt byte array
const saltLen = 32

const (
	ioReadErr    = "failed to read IDF %s: %v"
	unmarshalErr = "failed to unmarshal IDF JSON: %v"
	marshalErr   = "failed to JSON marshal the IDF: %v"
	writeErr     = "failed to write IDF: %v"
	saltSizeErr  = "length of salt is %d != %d expected"
)

// IdFile describes the information and structure for saving an ID and Salt to
// a JSON file that is both human readable and used for processing.
type IdFile struct {
	ID      id.ID             `json:"id"`
	Type    string            `json:"type"`
	Salt    [saltLen]byte     `json:"salt"`
	IdBytes [id.ArrIDLen]byte `json:"idBytes"`
}

// newIDF creates a a new IdFile with the given 32-byte salt and id.ID. An error
// is returned if the salt is not of the correct length.
func newIDF(salt []byte, genID id.ID) (IdFile, error) {
	// Check that the salt is of the correct length
	if len(salt) != saltLen {
		return IdFile{}, errors.Errorf(saltSizeErr, len(salt), saltLen)
	}

	// Create new IdFile with the ID
	newIDF := IdFile{
		ID:      genID,
		Type:    genID.GetType().String(),
		Salt:    [saltLen]byte{},
		IdBytes: genID,
	}

	// Copy salt into the IdFile
	copy(newIDF.Salt[:], salt)

	return newIDF, nil
}

// UnloadIDF unmarshal the JSON encoded IdFile at the given file path and
// returns its 32-byte salt and id.ID.
func UnloadIDF(path string) ([]byte, id.ID, error) {
	// Read the contents from the file
	jsonBytes, err := utils.ReadFile(path)
	if err != nil {
		return nil, id.ID{}, errors.Errorf(ioReadErr, path, err)
	}

	return unloadIDF(jsonBytes)
}

// unloadIDF unmarshalls the JSON data into an IdFile and returns the id.ID and
// a 32-byte salt.
func unloadIDF(data []byte) ([]byte, id.ID, error) {
	// Create empty IdFile and unmarshal the data into it
	var idf *IdFile
	if err := json.Unmarshal(data, &idf); err != nil {
		return nil, id.ID{}, errors.Errorf(unmarshalErr, err)
	}

	return idf.Salt[:], idf.ID, nil
}

// LoadIDF saves a JSON encoded IdFile containing the salt and generated ID.
func LoadIDF(filePath string, salt []byte, genID id.ID) error {
	// Marshal the IDF into JSON bytes
	idfJSON, err := loadIDF(salt, genID)
	if err != nil {
		return err
	}

	// Create new ID file
	err = utils.WriteFile(filePath, idfJSON, utils.FilePerms, utils.DirPerms)
	if err != nil {
		return errors.Errorf(writeErr, err)
	}

	return err
}

// loadIDF creates an IdFile from the salt and id.ID and returns it JSON
// encoded.
func loadIDF(salt []byte, genID id.ID) ([]byte, error) {
	// Generate new IdFile object
	idf, err := newIDF(salt, genID)
	if err != nil {
		return nil, err
	}

	// Marshal the IDF into JSON bytes
	idfJSON, err := json.Marshal(idf)
	if err != nil {
		return nil, errors.Errorf(marshalErr, err)
	}

	return idfJSON, err
}
