////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package id

// Stored here are global hard coded IDs. The last byte should be set to the
// correct ID type.

// ID for permissioning (ID data is the string "Permissioning")
var Permissioning = ID{80, 101, 114, 109, 105, 115, 115, 105, 111, 110, 105,
	110, 103, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(Generic)}

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
