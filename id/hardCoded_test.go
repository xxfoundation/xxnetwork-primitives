////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package id

import (
	"testing"
)

// Tests that GetHardCodedIDs() returns all the hard coded IDs in the order that
// they were added.
func TestGetHardCodedIDs(t *testing.T) {
	expectedIDs := []*ID{&Permissioning, &NotificationBot, &TempGateway,
		&ZeroUser, &DummyUser, &UDB}

	for i, testID := range GetHardCodedIDs() {
		if !expectedIDs[i].Cmp(testID) {
			t.Errorf("GetHardCodedIDs() did not return the expected ID at "+
				"index %d.\n\texepcted: %v\n\trecieved: %v",
				i, expectedIDs[i], testID)
		}
	}
}
