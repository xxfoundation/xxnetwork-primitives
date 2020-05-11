////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package id

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestNewIdListFromStrings(t *testing.T) {
	// Topology list will contain the list of strings to be passed to NewNodeListFromStrings
	var topologyList []string
	// ExpectedNodes will contain the constructed nodes, to be compared one-by-one to
	// the output of NewNodeListFromStrings
	var expectedNodes []*Node

	// Construct a topology list
	for i := 0; i < 10; i++ {
		// construct an id
		idBytes := make([]byte, NodeIdLen)
		rand.Read(idBytes)
		expectedId := NewNodeFromBytes(idBytes)

		// Append to the slices
		expectedNodes = append(expectedNodes, expectedId)
		topologyList = append(topologyList, expectedId.String())
	}

	// Pass topologyList into NewNodeListFromStrings
	receivedNodes, err := NewNodeListFromStrings(topologyList)
	if err != nil {
		t.Errorf("Failed to create node list: %+v", err)
	}

	// Iterate through the list, comparing receivedNodes to expectedNodes every iteration
	for index, receivedNode := range receivedNodes {
		expectedNode := expectedNodes[index]

		// Check the outputted list to the expected values
		if !bytes.Equal(receivedNode.Bytes(), expectedNode.Bytes()) {
			t.Errorf("Node of index %d was not converted correctly. "+
				"\n\treceived: %v\n\texpected: %v", index, receivedNode.Bytes(), expectedNode.Bytes())
		}

	}

}

// Error path: construct a list with a bad topology
func TestNewIdListFromStrings_Error(t *testing.T) {
	// Topology list will contain the list of strings to be passed to NewNodeListFromStrings
	var topologyList []string
	// ExpectedNodes will contain the constructed nodes, to be compared one-by-one to
	// the output of NewNodeListFromStrings
	var expectedNodes []*Node

	// Construct a topology list
	for i := 0; i < 10; i++ {
		// construct id
		idBytes := make([]byte, NodeIdLen)
		rand.Read(idBytes)
		expectedId := NewNodeFromBytes(idBytes)
		expectedNodes = append(expectedNodes, expectedId)

		// Inject a bad string into the list of strings at an arbitrary point
		if i == rand.Int()%9 {
			topologyList = append(topologyList, expectedId.String()+"badString")
		}
	}

	// Attempt to convert the topologyList
	_, err := NewNodeListFromStrings(topologyList)
	if err != nil {
		return
	}

	t.Errorf("Expected error case, should not successfully create a list of nodes due to a bad" +
		"topology")
}
