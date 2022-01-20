///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package fact

import (
	"reflect"
	"testing"
)

// Test NewFact() function returns a correctly formatted Fact
func TestNewFact(t *testing.T) {
	// Expected result
	e := Fact{
		Fact: "devinputvalidation@elixxir.io",
		T:    1,
	}

	g, err := NewFact(Email, "devinputvalidation@elixxir.io")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(e, g) {
		t.Errorf("The returned Fact did not match the expected Fact")
	}
}

// Test NewFact() returns error when a fact exceeds the maxFactCharacterLimit.
func TestNewFact_ExceedMaxFactError(t *testing.T) {
	// Expected error case
	_, err := NewFact(Email, "devinputvalidation_devinputvalidation_devinputvalidation@elixxir.io")
	if err == nil {
		t.Fatalf("NewFact expected to fail due to the fact exceeding maximum character length")
	}

}

// Test Stringify() creates a string of the Fact
// The output is verified to work in the test below
func TestFact_Stringify(t *testing.T) {
	f := Fact{
		Fact: "devinputvalidation@elixxir.io",
		T:    1,
	}

	expected := "Edevinputvalidation@elixxir.io"
	got := f.Stringify()
	t.Log(got)

	if got != expected {
		t.Errorf("Marshalled object from Got did not match Expected.\n\tGot: %v\n\tExpected: %v", got, expected)
	}
}

// Test the UnstringifyFact function creates a Fact from a string
// NOTE: this test does not pass, with error "Unknown Fact FactType: Etesting"
func TestFact_UnstringifyFact(t *testing.T) {
	// Expected fact from above test
	e := Fact{
		Fact: "devinputvalidation@elixxir.io",
		T:    Email,
	}

	// Stringify-ed Fact from above test
	m := "Edevinputvalidation@elixxir.io"
	f, err := UnstringifyFact(m)
	if err != nil {
		t.Error(err)
	}

	t.Log(f.Fact)
	t.Log(f.T)

	if !reflect.DeepEqual(e, f) {
		t.Errorf("The returned Fact did not match the expected Fact")
	}
}
