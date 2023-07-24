////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Package idf contains the structure of the ID File. It holds a generated ID
// and a salt (a 256-bit random number) in a JSON file. This file is used by
// different entities to save their ID and salt to file. The file path is
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
	saltSizeErr  = "salt length must be %d; received length of %d"
	ioReadErr    = "failed to read IDF from path %q"
	unmarshalErr = "failed to JSON unmarshal the IDF"
	marshalErr   = "failed to JSON marshal the IDF"
	ioWriteErr   = "failed to write IDF to path %q"
)

// IdFile is used to save IDs and salts to file in a human-readable form.
type IdFile struct {
	ID        *id.ID            `json:"id"`
	Type      string            `json:"type"`
	Salt      [saltLen]byte     `json:"salt"`
	IdBytes   [id.ArrIDLen]byte `json:"idBytes"`
	HexNodeID string            `json:"hexNodeID"`
}

// newIDF creates a new IdFile with the given 32-byte salt and id.ID. An error
// is returned if the salt is not of the correct length.
func newIDF(salt []byte, genID *id.ID) (IdFile, error) {
	// Check that the salt is of the correct length
	if len(salt) != saltLen {
		return IdFile{}, errors.Errorf(saltSizeErr, saltLen, len(salt))
	}

	idf := IdFile{
		ID:        genID,
		Type:      genID.GetType().String(),
		Salt:      [saltLen]byte{},
		IdBytes:   *genID,
		HexNodeID: genID.HexEncode(),
	}
	copy(idf.Salt[:], salt)

	return idf, nil
}

// UnloadIDF unmarshal the JSON encoded IdFile at the given file path and
// returns its 32-byte salt and id.ID.
func UnloadIDF(path string) ([]byte, *id.ID, error) {
	// Read the contents from the file
	jsonBytes, err := utils.ReadFile(path)
	if err != nil {
		return nil, nil, errors.Wrapf(err, ioReadErr, path)
	}

	var idf IdFile
	if err = json.Unmarshal(jsonBytes, &idf); err != nil {
		return nil, nil, errors.Wrap(err, unmarshalErr)
	}

	return idf.Salt[:], idf.ID, err
}

// LoadIDF creates an IdFile object with the provided values and
// marshals it into JSON bytes ready to be written to a file.
func LoadIDF(filePath string, salt []byte, genID *id.ID) error {
	// Generate new IdFile object
	idf, err := newIDF(salt, genID)
	if err != nil {
		return err
	}

	// Marshal the IDF into JSON bytes
	idfJSON, err := json.Marshal(idf)
	if err != nil {
		return errors.Wrap(err, marshalErr)
	}

	// Create new ID file
	err = utils.WriteFile(filePath, idfJSON, utils.FilePerms, utils.DirPerms)
	return errors.Wrapf(err, ioWriteErr, filePath)
}
