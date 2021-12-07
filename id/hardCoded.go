////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                           //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file                                                               //
////////////////////////////////////////////////////////////////////////////////

package id

// Stores the global hard coded IDs. The last byte should be set to the correct
// ID type.
// Note: When adding or removing a hard coded ID, make sure to update
// GetHardCodedIDs with the changes.

// Permissioning is the ID for the permissioning server (data is the string
// "Permissioning").
var Permissioning = ID{
	80, 101, 114, 109, 105, 115, 115, 105, 111, 110, 105, 110, 103, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(Generic),
}

// NotificationBot is the ID for the notification bot (data is the string
// "notification-bot").
var NotificationBot = ID{
	110, 111, 116, 105, 102, 105, 99, 97, 116, 105, 111, 110, 45, 98, 111, 116,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(Generic),
}

// TempGateway is the ID for a temporary gateway (data is the string "tmp").
var TempGateway = ID{
	116, 109, 112, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, byte(Gateway),
}

// ZeroUser is the ID for a user with the ID data set to all zeroes.
var ZeroUser = ID{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, byte(User),
}

// DummyUser is the ID for a dummy user (data is the string "dummy").
var DummyUser = ID{
	100, 117, 109, 109, 121, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(User),
}

// UDB is the ID for user discovery (data is in the range of dummy IDs).
var UDB = ID{
	0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, byte(User),
}

// GetHardCodedIDs returns an array of all the hard coded IDs.
func GetHardCodedIDs() (ids []ID) {
	return []ID{
		Permissioning,
		NotificationBot,
		TempGateway, ZeroUser,
		DummyUser,
		UDB,
	}
}

// CollidesWithHardCodedID searches if the given ID collides with any hard coded
// IDs. If it collides, then the function returns true. Otherwise, it returns
// false.
func CollidesWithHardCodedID(testID ID) bool {
	for _, hardCodedID := range GetHardCodedIDs() {
		if testID == hardCodedID {
			return true
		}
	}

	return false
}
