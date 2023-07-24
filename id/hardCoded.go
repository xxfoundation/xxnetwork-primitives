////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package id

// Stores the global hard coded IDs. The last byte should be set to the correct
// ID type.
//
// Note: When adding or removing a hard coded ID, make sure to update
// GetHardCodedIDs with the changes.

var (
	// Permissioning is the ID for the permissioning server (data is the string
	// "Permissioning").
	Permissioning = ID{80, 101, 114, 109, 105, 115, 115, 105, 111, 110, 105,
		110, 103, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		byte(Generic)}

	// Authorizer is the ID for the authorizer (data is the string "authorizer").
	Authorizer = ID{97, 117, 116, 104, 111, 114, 105, 122, 101, 114, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(Generic)}

	// ClientRegistration is the ID for client registration (ID data is the
	// string "client-registration")
	ClientRegistration = ID{99, 108, 105, 101, 110, 116, 45, 114, 101, 103, 105,
		115, 116, 114, 97, 116, 105, 111, 110, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, byte(Generic)}

	// NotificationBot is the ID for the notification bot (data is the string
	// "notification-bot").
	NotificationBot = ID{110, 111, 116, 105, 102, 105, 99, 97, 116, 105, 111,
		110, 45, 98, 111, 116, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		byte(Generic)}

	// TempGateway is the ID for a temporary gateway (data is the string "tmp").
	TempGateway = ID{116, 109, 112, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(Gateway)}

	// ZeroUser is the ID for a user with the ID data set to all zeroes.
	ZeroUser = ID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(User)}

	// DummyUser is the ID for a dummy user (data is the string "dummy").
	DummyUser = ID{100, 117, 109, 109, 121, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(User)}

	// UDB is the ID for user discovery (data is in the range of dummy IDs).
	UDB = ID{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(User)}
)

// GetHardCodedIDs returns an array of all the hard coded IDs.
func GetHardCodedIDs() []*ID {
	return []*ID{
		&Permissioning,
		&Authorizer,
		&NotificationBot,
		&TempGateway,
		&ZeroUser,
		&DummyUser,
		&UDB,
		&ClientRegistration,
	}
}

// CollidesWithHardCodedID searches if the given ID collides with any hard coded
// IDs. If it collides, then the function returns true. Otherwise, it returns
// false.
func CollidesWithHardCodedID(testID *ID) bool {
	for _, hardCodedID := range GetHardCodedIDs() {
		if testID.Equal(hardCodedID) {
			return true
		}
	}

	return false
}
