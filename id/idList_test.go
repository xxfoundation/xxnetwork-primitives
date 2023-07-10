////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package id

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestNewIDListFromBytes(t *testing.T) {
	// Topology list will contain the list of strings to be passed to
	// ewIDListFromBytes
	var topologyList [][]byte

	// ExpectedNodes will contain the constructed IDs, to be compared one-by-one
	// to the output of NewIDListFromBytes
	var expectedIDs []*ID

	// Construct a topology list
	for i := 0; i < 10; i++ {
		// Construct an ID
		idBytes := newRandomBytes(ArrIDLen, t)
		expectedId := NewIdFromBytes(idBytes, t)

		// Append to the slices
		expectedIDs = append(expectedIDs, expectedId)
		topologyList = append(topologyList, expectedId.Bytes())
	}

	// Pass topologyList into NewIDListFromBytes
	receivedIDs, err := NewIDListFromBytes(topologyList)
	if err != nil {
		t.Errorf("Failed to create ID list: %+v", err)
	}

	// Iterate through the list and comparing receivedIDs to expectedIDs every
	// iteration
	for index, receivedID := range receivedIDs {
		expectedNode := expectedIDs[index]

		// Check the outputted list to the expected values
		if !bytes.Equal(receivedID.Bytes(), expectedNode.Bytes()) {
			t.Errorf("ID of index %d was not converted correctly."+
				"\nreceived: %v\nexpected: %v", index, receivedID.Bytes(),
				expectedNode.Bytes())
		}

	}

}

// Error path: construct a list with a bad topology
func TestNewIDListFromBytes_Error(t *testing.T) {
	// Topology list will contain the list of strings to be passed to
	// NewIDListFromBytes
	var topologyList [][]byte

	// ExpectedNodes will contain the constructed nodes, to be compared
	// one-by-one to the output of NewIDListFromBytes
	var expectedNodes []*ID

	// Construct a topology list
	for i := 0; i < 10; i++ {
		// Construct ID
		idBytes := newRandomBytes(ArrIDLen, t)
		expectedId := NewIdFromBytes(idBytes, t)
		expectedNodes = append(expectedNodes, expectedId)

		// Inject a bad byte slices into the list of slices at an arbitrary point
		if i == rand.Int()%5 {
			expectedIdBytes := []byte{1, 2, 3}
			topologyList = append(topologyList, expectedIdBytes)
		}
	}

	// Attempt to convert the topologyList
	_, err := NewIDListFromBytes(topologyList)
	if err != nil {
		return
	}

	t.Errorf("Expected error case, should not successfully create a list" +
		"of nodes due to a bad topology")
}
