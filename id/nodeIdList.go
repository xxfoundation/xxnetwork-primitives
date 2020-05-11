////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package id

import "github.com/pkg/errors"

// NodeIdList.go handles operations that create a list of id.Node objects

// NewNodeListFromStrings creates a list of Node's from a list of strings
//  On success it returns a new list
//  On failure it returns a nil list and an error
func NewNodeListFromStrings(topology []string) ([]*Node, error) {
	list := make([]*Node, len(topology))
	for index, id := range topology {
		newId, err := NewNodeFromString(id)
		if err != nil {
			return nil, errors.Errorf("Unable to convert id for index %d: %+v", index, err)
		}

		list[index] = newId
	}

	return list, nil
}
