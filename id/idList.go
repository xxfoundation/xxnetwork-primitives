////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package id

import "github.com/pkg/errors"

// idList.go handles operations that create a list of id.ID objects.

// NewIDListFromBytes creates a list of IDs from a list of byte slices. On
// success, it returns a new list. On failure, it returns a nil list and an
// error.
func NewIDListFromBytes(topology [][]byte) ([]*ID, error) {
	list := make([]*ID, len(topology))

	for index, id := range topology {
		newId, err := Unmarshal(id)
		if err != nil {
			return nil, errors.Errorf("Unable to marshal ID for index %d: %+v",
				index, err)
		}

		list[index] = newId
	}

	return list, nil
}
