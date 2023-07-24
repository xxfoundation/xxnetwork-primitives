////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package id

import (
	"testing"
)

// Tests that GetHardCodedIDs returns all the hard coded IDs in the order that
// they were added.
func TestGetHardCodedIDs(t *testing.T) {
	expectedIDs := []*ID{&Permissioning, &Authorizer, &NotificationBot,
		&TempGateway, &ZeroUser, &DummyUser, &UDB, &ClientRegistration}

	for i, testID := range GetHardCodedIDs() {
		if !expectedIDs[i].Equal(testID) {
			t.Errorf("GetHardCodedIDs did not return the expected ID (%d)."+
				"\nexepcted: %s\nrecieved: %s", i, expectedIDs[i], testID)
		}
	}
}

// Tests that CollidesWithHardCodedID returns false when none of the test IDs
// collide with the hard coded IDs.
func TestCollidesWithHardCodedID_HappyPath(t *testing.T) {
	testIDs := []*ID{
		NewIdFromString("Test1", Generic, t),
		NewIdFromString("Test2", Gateway, t),
		NewIdFromString("Test3", Node, t),
		NewIdFromString("Test4", User, t),
		NewIdFromString("Test4", Group, t),
	}

	for _, testID := range testIDs {
		if CollidesWithHardCodedID(testID) {
			t.Errorf("CollidesWithHardCodedID found collision when none "+
				"should exist.\ncolliding ID: %v", testID)
		}
	}
}

// Tests that CollidesWithHardCodedID returns true when checking if hard coded
// IDs collide.
func TestCollidesWithHardCodedID_True(t *testing.T) {
	testIDs := []*ID{&Permissioning, &Authorizer, &NotificationBot, &TempGateway,
		&ZeroUser, &DummyUser, &UDB, &ClientRegistration}

	for _, testID := range testIDs {
		if !CollidesWithHardCodedID(testID) {
			t.Errorf("CollidesWithHardCodedID did not find a collision when "+
				"one should exist.\ncolliding ID: %v", testID)
		}
	}
}
