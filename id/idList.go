////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package id

import "github.com/pkg/errors"

// idList.go handles operations that create a list of id.ID objects.

// NewIDListFromBytes creates a list of id.ID from a list of byte slices returns
// it. An error is returned if any IDs fail to unmarshal.
func NewIDListFromBytes(topology [][]byte) ([]*ID, error) {
	list := make([]*ID, len(topology))

	for i, idBytes := range topology {
		id, err := Unmarshal(idBytes)
		if err != nil {
			return nil, errors.Errorf(
				"unable to unmarshal ID at index %d: %+v", i, err)
		}

		list[i] = id
	}

	return list, nil
}
