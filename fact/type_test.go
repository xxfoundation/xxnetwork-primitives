///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package fact

import (
	"testing"
)

func TestFactType_String(t *testing.T) {
	// FactTypes and expected strings for them
	FTs := []FactType{Username, Email, Phone, Nickname, FactType(200)}
	Strs := []string{"Username", "Email", "Phone", "Nickname", "Unknown Fact FactType: 200"}
	for i, ft := range FTs {
		if FactType.String(ft) != Strs[i] {
			t.Errorf("Got unexpected string for FactType.\n\tGot: %s\n\tExpected: %s", FactType.String(ft), Strs[i])
		}
	}
}

func TestFactType_Stringify(t *testing.T) {
	// FactTypes and expected strings for them
	FTs := []FactType{Username, Email, Phone, Nickname}
	Strs := []string{"U", "E", "P", "N"}
	for i, ft := range FTs {
		if FactType.Stringify(ft) != Strs[i] {
			t.Errorf("Got unexpected string for FactType.\n\tGot: %s\n\tExpected: %s", FactType.Stringify(ft), Strs[i])
		}
	}
}

func TestFactType_Stringify_Error(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Stringify() did not panic when the FactType passed in is unknown.")
		}

	}()

	// Test that an unknown FactType results in a Panic
	FactType.Stringify(FactType(200))
}

func TestFactType_Unstringify(t *testing.T) {
	// FactTypes and expected strings for them
	FTs := []FactType{Username, Email, Phone, Nickname}
	Strs := []string{"U", "E", "P", "N"}
	for i, ft := range FTs {
		gotft, err := UnstringifyFactType(Strs[i])
		if err != nil {
			t.Error(err)
		}
		if gotft != ft {
			t.Errorf("Got unexpected string for FactType.\n\tGot: %s\n\tExpected: %s", FactType.Stringify(ft), Strs[i])
		}
	}

	_, err := UnstringifyFactType("x")
	if err == nil {
		t.Errorf("UnstringifyFactType did not return an error on an invalid type")
	}
}

func TestFactType_IsValid(t *testing.T) {
	if !FactType.IsValid(Username) ||
		!FactType.IsValid(Email) ||
		!FactType.IsValid(Phone) {

		t.Errorf("FactType.IsValid did not report a FactType as valid")
	}

	if FactType.IsValid(FactType(200)) {
		t.Errorf("FactType.IsValid reported a non-valid FactType value as valid")
	}
}
