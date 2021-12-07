////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                           //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file                                                               //
////////////////////////////////////////////////////////////////////////////////
package id

import (
	"strings"
	"testing"
)

func TestNewIDListFromBytes(t *testing.T) {
	// Construct a topology list and a list of matching IDs
	n := 10
	expectedIDs := make([]ID, n)
	topologyList := make([][]byte, n)
	for i := 0; i < 10; i++ {
		expectedId := NewIdFromBytes(rngBytes(ArrIDLen, int64(i+42), t), t)
		expectedIDs[i] = expectedId
		topologyList[i] = expectedId.Bytes()
	}

	// Pass topologyList into NewIDListFromBytes
	receivedIDs, err := NewIDListFromBytes(topologyList)
	if err != nil {
		t.Errorf("NewIDListFromBytes returned an error: %+v", err)
	}

	// Iterate through the list and comparing receivedIDs to expectedIDs every
	// iteration
	for i, receivedID := range receivedIDs {
		// Check the outputted list to the expected values
		if receivedID != expectedIDs[i] {
			t.Errorf("ID of index %d was not converted correctly."+
				"\nreceived: %s\nexpected: %s", i, receivedID, expectedIDs[i])
		}

	}

}

// Error path: construct a list with a bad topology.
func TestNewIDListFromBytes_Error(t *testing.T) {
	topologyList := [][]byte{
		NewIdFromBytes(rngBytes(ArrIDLen, 42, t), t).Bytes(),
		[]byte("invalid ID"),
		NewIdFromBytes(rngBytes(ArrIDLen, 42, t), t).Bytes(),
	}

	// Attempt to convert the topologyList
	_, err := NewIDListFromBytes(topologyList)
	if err == nil || !strings.Contains(err.Error(), "failed to unmarshal ID") {
		t.Errorf("NewIDListFromBytes did not return an error when an invalid " +
			"ID is in the list.")
	}

}
