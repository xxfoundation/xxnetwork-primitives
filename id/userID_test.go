////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package id

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

// Proves that the registration code we get for a certain user is what we get
func TestUserID_RegistrationCode(t *testing.T) {
	expected := "RUHPS2MI" // reg code for user 1
	var id User
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
// You shouldn't use NewUserFromUint in production code! Use NewUserFromUints instead
// and have the first three uints be zero.
// I wrote NewUserFromUint for compatibility with a lot of our tests that
// populated user IDs from a single uint64 and don't care too much about
// propagating the whole thing.
func TestNewUserFromUint(t *testing.T) {
	// This particular method for new-ing a user ID is only able to fill out
	// the bytes on the little end
	intId := uint64(rand.Int63())
	id := NewUserFromUint(intId, t)
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
		t.Error("A byte that NewUserFromUint set wasn't identical to the" +
			" uint64 input")
	}
}

// Proves that setting the bytes populates the user ID with all the same bytes
func TestUserID_SetBytes(t *testing.T) {
	idBytes := make([]byte, UserLen)
	rand.Read(idBytes)
	id := NewUserFromBytes(idBytes)
	if !bytes.Equal(id[:], idBytes) {
		t.Error("NewNodeFromBytes didn't set all the bytes correctly")
	}
}

// Proves that providing invalid input (wrong length) to NewNodeFromBytes() gives an
// invalid result
func TestUserID_SetBytes_Error(t *testing.T) {
	var idBytes []byte
	id := NewUserFromBytes(idBytes)
	if !id.Cmp(ZeroID) {
		t.Error("Got a non-zero ID out of setting the bytes, but shouldn't have")
	}
	if id != nil {
	}
}

// Proves that NewUserFromUints populates data all over the user ID as expected
func TestUserID_SetUints(t *testing.T) {
	uints := [4]uint64{798264, 350789, 34076, 154268}
	id := NewUserFromUints(&uints)
	for i := 0; i < len(uints); i++ {
		if binary.BigEndian.Uint64(id[i*8:]) != uints[i] {
			t.Errorf("Uint64 differed at index %v", i)
		}
	}
}

// Proves that Bytes converts a user ID to a mostly identical byte slice for
// easy compatibility with methods that take a byte slice
func TestUserID_Bytes(t *testing.T) {
	idBytes := make([]byte, UserLen)
	rand.Read(idBytes)
	id := NewUserFromBytes(idBytes)
	if !bytes.Equal(idBytes, id.Bytes()) {
		t.Error("Surprisingly, " +
			"the Bytes() method didn't return an equivalent byteslice")
	}
}

// Tests that Bytes() correctly makes a new copy of the bytes.
func TestUserID_Bytes_Copy(t *testing.T) {
	idBytes := make([]byte, UserLen)
	rand.Read(idBytes)
	id := NewUserFromBytes(idBytes)

	userBytes := id.Bytes()

	// Modify the original
	for j := 0; j < UserLen; j++ {
		id[j] = ^id[j]
	}

	if !bytes.Equal(userBytes, idBytes) {
		t.Errorf("Bytes() returned incorrect byte slice of User ID"+
			"\n\treceived: %v\n\texpected: %v", userBytes, idBytes)
	}
}

// Proves that equal returns true when two IDs are equal and returns false when
// they aren't
func TestUser_Cmp(t *testing.T) {
	id1 := User{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	id2 := User{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	id3 := User{1, 2, 3, 4, 5, 6, 7, 8, 9, 11}

	if !id1.Cmp(&id2) {
		t.Error("ID 1 and 2 should have been equal, but weren't")
	}
	if id3.Cmp(&id1) {
		t.Error("ID 1 and 3 shouldn't have been equal, but were")
	}
}

//Test that deep copy returns an exact copy and that changing one does nto change the other
func TestUser_DeepCopy(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		var original User
		rng.Read(original[:])
		deepcopy := (&original).DeepCopy()

		if !reflect.DeepEqual(original, *deepcopy) {
			t.Errorf("User.DeepCopy: On Attempt %v copy does not equal origonal %v %v", i, original, deepcopy)
		}

		for j := 0; j < UserLen; j++ {
			original[j] = ^original[j]
		}

		if reflect.DeepEqual(original, *deepcopy) {
			t.Errorf("User.DeepCopy: On Attempt %v copy still linked to origonal %v %v", i, original, deepcopy)
		}
	}
}
