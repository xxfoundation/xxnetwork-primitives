////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package id

// Stored here are global hard coded IDs. The last byte should be set to the
// correct ID type.

// Note: When adding or removing a hard coded ID, make sure to update
// GetHardCodedIDs() with the changes.

// ID for permissioning (ID data is the string "Permissioning")
var Permissioning = ID{80, 101, 114, 109, 105, 115, 115, 105, 111, 110, 105,
	110, 103, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(Generic)}

// ID for authorizer (ID data is the string "authorizer")
var Authorizer = ID{97, 117, 116, 104, 111, 114, 105, 122, 101, 114, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(Generic)}

// ID for authorizer (ID data is the string "client-registration")
var ClientRegistration = ID{99, 108, 105, 101, 110, 116, 45, 114, 101, 103,
	105, 115, 116, 114, 97, 116, 105, 111, 110, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, byte(Generic)}

// ID for notification bot (ID data is the string "notification-bot")
var NotificationBot = ID{110, 111, 116, 105, 102, 105, 99, 97, 116, 105, 111,
	110, 45, 98, 111, 116, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(Generic)}

// ID for a temporary gateway (ID data is the string "tmp")
var TempGateway = ID{116, 109, 112, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(Gateway)}

// ID for a user with the ID data all zeroes
var ZeroUser = ID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(User)}

// ID for a dummy user (ID data is the string "dummy")
var DummyUser = ID{100, 117, 109, 109, 121, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(User)}

// ID for UDB (ID data is in the range of dummy IDs)
var UDB = ID{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(User)}

// GetHardCodedIDs returns an array of all the hard coded IDs.
func GetHardCodedIDs() (ids []*ID) {
	ids = append(ids, &Permissioning)
	ids = append(ids, &Authorizer)
	ids = append(ids, &NotificationBot)
	ids = append(ids, &TempGateway)
	ids = append(ids, &ZeroUser)
	ids = append(ids, &DummyUser)
	ids = append(ids, &UDB)
	ids = append(ids, &ClientRegistration)

	return
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
