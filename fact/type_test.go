////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package fact

import (
	"testing"
)

// Consistency test of FactType.String.
func TestFactType_String(t *testing.T) {
	tests := map[FactType]string{
		Username:      "Username",
		Email:         "Email",
		Phone:         "Phone",
		Nickname:      "Nickname",
		FactType(200): "Unknown Fact FactType: 200",
	}

	for ft, expected := range tests {
		str := ft.String()
		if expected != str {
			t.Errorf("Unexpected FactType string.\nexpected: %q\nreceived: %q",
				expected, str)
		}
	}
}

// Tests that a FactType marshalled by FactType.Stringify and unmarshalled by
// UnstringifyFactType matches the original.
func TestFactType_Stringify_UnstringifyFactType(t *testing.T) {
	factTypes := []FactType{
		Username,
		Email,
		Phone,
		Nickname,
	}

	for _, expected := range factTypes {
		str := expected.Stringify()

		ft, err := UnstringifyFactType(str)
		if err != nil {
			t.Fatalf("Failed to unstringify fact type %q: %+v", str, err)
		}
		if expected != ft {
			t.Errorf("Unexpected unstringified FactType."+
				"\nexpected: %s\nreceived: %s", expected, str)
		}
	}
}

// Panic path: Tests that FactType.Stringify panics for an invalid FactType
func TestFactType_Stringify_InvalidFactTypePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Failed to panic for invalid FactType")
		}
	}()

	FactType(99).Stringify()
}

// Error path: Tests that FactType.UnstringifyFactType returns an error for an
// invalid FactType.
func TestFactType_Unstringify_UnknownFactTypeError(t *testing.T) {
	_, err := UnstringifyFactType("invalid")
	if err == nil {
		t.Errorf("Failed to get error for invalid FactType.")
	}
}

func TestFactType_IsValid(t *testing.T) {
	tests := map[FactType]bool{
		Username: true,
		Email:    true,
		Phone:    true,
		Nickname: true,
		99:       false,
	}

	for ft, expected := range tests {
		if ft.IsValid() != expected {
			t.Errorf("Unexpected IsValid result for %s."+
				"\nexpected: %t\nreceived: %t", ft, expected, ft.IsValid())
		}
	}
}
