package contact

import (
	"encoding/json"
	"gitlab.com/xx_network/primitives/id"
	"reflect"
	"testing"
)

// Test that GetID returns an id.ID in the Contact
func TestContact_GetID(t *testing.T) {
	C := Contact{ID: &id.DummyUser}
	if !reflect.DeepEqual(C.GetID(), id.DummyUser.Bytes()) {
		t.Error("Did not return the correct ID")
	}
}

// Test that GetDHPublicKey returns the cyclic int in Contact
func TestContact_GetDHPublicKey(t *testing.T) {
	// TODO: Find a test key
	C := Contact{DhPubKey: []byte{5}}
	dh := C.GetDHPublicKey()
	if !reflect.DeepEqual(dh, []byte{5}) {
		t.Error("Returned DH key did not match expected DH key")
	}
}

// Test that GetOwnershipProof returns the []byte in Contact
func TestContact_GetOwnershipProof(t *testing.T) {
	C := Contact{OwnershipProof: []byte{30, 40, 50}}
	if !reflect.DeepEqual([]byte{30, 40, 50}, C.GetOwnershipProof()) {
		t.Error("Returned proof key did not match expected proof")
	}
}

// Test that GetFactList returns a FactList and our Fact is in it
func TestContact_GetFactList(t *testing.T) {
	FL := new([]Fact)
	C := Contact{Facts: *FL}

	C.Facts = append(C.Facts, Fact{
		Fact: "testing",
		T:    Phone,
	})

	gFL := C.GetFactList()
	if !reflect.DeepEqual(C.Facts[0], gFL.Get(0)) {
		t.Error("Expected Fact and got Fact did not match")
	}
}

// Test that Marshal can complete without error, the output should
// be verified with the test below
func TestContact_Marshal(t *testing.T) {
	C := Contact{
		ID:             &id.DummyUser,
		DhPubKey:       []byte{5},
		OwnershipProof: []byte{30, 40, 50},
		Facts: []Fact{{
			Fact: "testing",
			T:    Email,
		}},
	}

	M, err := C.Marshal()
	if err != nil {
		t.Error(err)
	}
	t.Log(M)
}

// Test that Unmarshal successfully unmarshals the Contact marshalled
// above
func TestContact_Unmarshal(t *testing.T) {
	// Expected contact
	E := Contact{
		ID:             &id.DummyUser,
		DhPubKey:       []byte{5},
		OwnershipProof: []byte{30, 40, 50},
		Facts: []Fact{{
			Fact: "testing",
			T:    Email,
		}},
	}

	// Marshalled contact, gotten from above test
	M := []byte{123, 34, 73, 68, 34, 58, 91, 49, 48, 48, 44, 49, 49, 55,
		44, 49, 48, 57, 44, 49, 48, 57, 44, 49, 50, 49, 44, 48, 44, 48,
		44, 48, 44, 48, 44, 48, 44, 48, 44, 48, 44, 48, 44, 48, 44, 48,
		44, 48, 44, 48, 44, 48, 44, 48, 44, 48, 44, 48, 44, 48, 44, 48,
		44, 48, 44, 48, 44, 48, 44, 48, 44, 48, 44, 48, 44, 48, 44, 48,
		44, 48, 44, 51, 93, 44, 34, 68, 104, 80, 117, 98, 75, 101, 121,
		34, 58, 123, 125, 44, 34, 79, 119, 110, 101, 114, 115, 104, 105,
		112, 80, 114, 111, 111, 102, 34, 58, 34, 72, 105, 103, 121, 34,
		44, 34, 70, 97, 99, 116, 115, 34, 58, 91, 123, 34, 70, 97, 99,
		116, 34, 58, 34, 116, 101, 115, 116, 105, 110, 103, 34, 44, 34,
		84, 34, 58, 49, 125, 93, 125}

	C, err := Unmarshal(M)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(C, E) {
		t.Error("Expected and Unmarshaled contact are not equal")
	}
}

// Test that Unmarshalling a contact with a bad FactType Fact fails
func TestContact_UnmarshalBadFact(t *testing.T) {
	C := Contact{Facts: []Fact{{
		Fact: "testing",
		T:    200,
	}}}

	M, err := json.Marshal(&C)
	if err != nil {
		t.Error(err)
	}

	_, err = Unmarshal(M)
	if err == nil {
		t.Error("Unmarshalling Contact containing Fact with an invalid Fact type should've errored")
	}
}

// Test that StringifyFacts can complete without error,
// the output should be verified with the test below
func TestContact_StringifyFacts(t *testing.T) {
	C := Contact{Facts: []Fact{
		{
			Fact: "testing",
			T:    Phone,
		},
		{
			Fact: "othertest",
			T:    Email,
		},
	}}

	S := C.StringifyFacts()
	t.Log(S)
}

// Test that UnstringifyFacts successfully unstingify-ies
// the FactList stringified above
// NOTE: This test does not pass, for reason "Invalid fact string passed"
func TestUnstringifyFacts(t *testing.T) {
	E := Contact{Facts: []Fact{
		{
			Fact: "testing",
			T:    Phone,
		},
		{
			Fact: "othertest",
			T:    Email,
		},
	}}

	FL, S, err := UnstringifyFacts("Ptesting,Eothertest;")
	if err != nil {
		t.Error(err)
	}

	t.Log(S)

	if !reflect.DeepEqual(E.Facts, FL) {
		t.Error("Expected FactList and got FactList do not match")
	}
}
