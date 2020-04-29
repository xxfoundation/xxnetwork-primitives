////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package id

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestNewIDListFromBytes(t *testing.T) {
	// Topology list will contain the list of strings to be passed to
	// NewIDListFromBytes
	var topologyList [][]byte

	// ExpectedNodes will contain the constructed nodes, to be compared
	// one-by-one to the output of NewIDListFromBytes
	var expectedNodes []*ID

	// Construct a topology list
	for i := 0; i < 10; i++ {
		// construct an id
		idBytes := newRandomBytes(ArrIDLen, t)
		expectedId := NewIdFromBytes(idBytes, t)

		// Append to the slices
		expectedNodes = append(expectedNodes, expectedId)
		topologyList = append(topologyList, expectedId.Bytes())
	}

	// Pass topologyList into NewIDListFromBytes
	receivedNodes, err := NewIDListFromBytes(topologyList)
	if err != nil {
		t.Errorf("Failed to create node list: %+v", err)
	}

	// Iterate through the list, comparing receivedNodes to expectedNodes every
	// iteration
	for index, receivedNode := range receivedNodes {
		expectedNode := expectedNodes[index]

		// Check the outputted list to the expected values
		if !bytes.Equal(receivedNode.Bytes(), expectedNode.Bytes()) {
			t.Errorf("Node of index %d was not converted correctly. "+
				"\n\treceived: %v\n\texpected: %v", index, receivedNode.Bytes(),
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
		// construct id
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
