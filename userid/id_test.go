////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package userid

import (
	"testing"
	"math/rand"
	"encoding/binary"
	"bytes"
)

// Proves that the registration code we get for a certain user is what we get
func TestUserID_RegistrationCode(t *testing.T) {
	expected := "RUHPS2MI" // reg code for user 1
	var id UserID
	copy(id[len(id)-1:], []byte{0x01})
	actual := id.RegistrationCode()
	if actual != expected {
		t.Errorf("Registration code differed from expected. Got %v, "+
			"expected %v", actual, expected)
	}
}

// Proves that results from setting up a new user ID from one uint64 are as
// expected, i.e. the first 3 uints worth of space are zero and the last uint
// worth of space is filled.
// You shouldn't use NewUserIDFromUint in production code! Use SetUints instead
// and have the first three uints be zero.
// I wrote NewUserIDFromUint for compatibility with a lot of our tests that
// populated user IDs from a single uint64 and don't care too much about
// propagating the whole thing.
func TestNewUserIDFromUint(t *testing.T) {
	// This particular method for new-ing a user ID is only able to fill out
	// the bytes on the little end
	intId := uint64(rand.Int63())
	id := NewUserIDFromUint(intId, t)
	// The first 64*3 bits should be left at zero
	for i := 0; i < sizeofUint64*3; i++ {
		if id[i] != 0 {
			t.Error("A byte that should have been zero wasn't")
		}
	}
	// The last bits should be the same starting at the big end of the int ID
	intIdBigEndian := make([]byte, 64/8)
	binary.BigEndian.PutUint64(intIdBigEndian, intId)
	if !bytes.Equal(intIdBigEndian, id[sizeofUint64*3:]) {
		t.Error("A byte that NewUserIDFromUint set wasn't identical to the" +
			" uint64 input")
	}
}

// Proves that setting the bytes populates the user ID with all the same bytes
func TestUserID_SetBytes(t *testing.T) {
	idBytes := make([]byte, UserIDLen)
	rand.Read(idBytes)
	id := new(UserID).SetBytes(idBytes)
	if !bytes.Equal(id[:], idBytes) {
		t.Error("SetBytes didn't set all the bytes correctly")
	}
}

// Proves that providing invalid input (wrong length) to SetBytes() gives an
// invalid result
func TestUserID_SetBytes_Error(t *testing.T) {
	var idBytes []byte
	id := new(UserID).SetBytes(idBytes)
	if !Equal(id, ZeroID) {
		t.Error("Got a non-zero ID out of setting the bytes, but shouldn't have")
	}
	if id != nil {
	}
}

// Proves that SetUints populates data all over the user ID as expected
func TestUserID_SetUints(t *testing.T) {
	uints := [4]uint64{798264,350789,34076,154268}
	id := new(UserID).SetUints(&uints)
	for i := 0; i < len(uints); i++ {
		if binary.BigEndian.Uint64(id[i*8:]) != uints[i] {
			t.Errorf("Uint64 differed at index %v", i)
		}
	}
}

// Proves that Bytes converts a user ID to a mostly identical byte slice for
// easy compatibility with methods that take a byte slice
func TestUserID_Bytes(t *testing.T) {
	idBytes := make([]byte, UserIDLen)
	rand.Read(idBytes)
	id := new(UserID).SetBytes(idBytes)
	if !bytes.Equal(idBytes, id.Bytes()) {
		t.Error("Surprisingly, " +
			"the Bytes() method didn't return an equivalent byteslice")
	}
}

// Proves that equal returns true when two IDs are equal and returns false when
// they aren't
func TestEqual(t *testing.T) {
	id1 := UserID{1,2,3,4,5,6,7,8,9,10}
	id2 := UserID{1,2,3,4,5,6,7,8,9,10}
	id3 := UserID{1,2,3,4,5,6,7,8,9,11}

	if !Equal(&id1, &id2) {
		t.Error("ID 1 and 2 should have been equal, but weren't")
	}
	if Equal(&id3, &id1) {
		t.Error("ID 1 and 3 shouldn't have been equal, but were")
	}
}
