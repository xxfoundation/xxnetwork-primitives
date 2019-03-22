////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package id

import (
	"encoding/base32"
	"encoding/binary"
	"golang.org/x/crypto/blake2b"
	"testing"
)

// Length of IDs in bytes
// 256 bits
const UserLen = 32

// Most string types in most languages (with C excepted) support 0 as a
// character in a string, for Unicode support. So it's possible to use normal
// strings as an immutable container for bytes in all the languages we care
// about supporting.
// However, when marshaling strings into protobufs, you'll get errors when
// the string isn't a valid UTF-8 string. So, the alternative underlying type
// that you can use as a map key in Go is an array, and that's what the package
// should use.
// TODO Should we enforce any aspects of the user ID generation here?
// (e.g. the first bit being zero to enforce everything being in the cyclic group)
type User [UserLen]byte

// Use this if you don't want to have to populate user ids for this manually
var ZeroID *User

func init() {
	// A zero ID should have all its bytes set to zero
	ZeroID = new(User).SetBytes(make([]byte, UserLen))
}

// Length of registration code in raw bytes
// Must be a multiple of 5 bytes to work with base 32
// 8 character long reg codes when base-32 encoded currently with length of 5
const RegCodeLen = 5

// This is a stopgap to be able to register fake users for fake demos.
// Replace ASAP!
func (u *User) RegistrationCode() string {
	return base32.StdEncoding.EncodeToString(userHash(u))
}

// userHash generates a hash of the UID to be used as a registration code for
// demos
// TODO Should we use the full-length hash? Should we even be doing registration
// like this?
func userHash(uid *User) []byte {
	h, _ := blake2b.New256(nil)
	h.Write(uid[:])
	huid := h.Sum(nil)
	huid = huid[len(huid)-RegCodeLen:]
	return huid
}

const sizeofUint64 = 8

// Only tests should use this method for compatibility with the old user ID
// structure, as a utility method to easily create user IDs with the correct
// length. So this func takes a testing.T.
func NewUserFromUint(newId uint64, t *testing.T) *User {
	// TODO Uncomment these lines to cause failure where this method's used in
	// the real codebase. Then, replace those occurrences with better code.
	//t.Log("Warning: Creating a new user ID from uint. " +
	//	"You should create user IDs some other way.")
	var result User
	binary.BigEndian.PutUint64(result[UserLen-sizeofUint64:], newId)
	return &result
}

// Since user IDs are 256 bits long, you need four uint64s to be able to set
// all the bits with uints. All the uints are big-endian, and are put in the
// ID in big-endian order above that.
func (u *User) SetUints(uints *[4]uint64) *User {
	for i := range uints {
		binary.BigEndian.PutUint64(u[i*8:], uints[i])
	}
	return u
}

// Returns a user ID set to the contents of the byte slice if the byte slice
// has the correct length
// Otherwise, returns a user ID that's all zeroes
func (u *User) SetBytes(data []byte) *User {
	if len(data) != UserLen {
		// Return a user ID with all zeroes which should get rejected somewhere
		// along the line due to cryptographic properties that the system provides
		return new(User)
	} else {
		copy(u[:], data)
		return u
	}
}

// Utility function to convert a user ID to a byte slice
func (u *User) Bytes() []byte {
	return u[:]
}

// Utility function to determine whether two user IDs are equal
func Equal(lhs *User, rhs *User) bool {
	return *lhs == *rhs
}

func (u *User) DeepCopy() *User {
	if u == nil {
		return nil
	}
	var nu User
	copy(nu[:], (*u)[:])
	return &nu
}
