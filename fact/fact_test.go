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

// Unit test for input validation of emails
// NOTE: Tests for here might fail on goland due to bad internet connections
//  it is likely this will pass remotely
func TestValidateFact_Email(t *testing.T) {
	// Valid Fact
	validFact := Fact{
		Fact: "devinputvalidation@elixxir.io",
		T:    Email,
	}

	// Happy path with valid email and host
	err := ValidateFact(validFact)
	if err != nil {
		t.Errorf("Unexpected error in happy path: %v", err)
	}

	// Invalid Fact Host
	invalidEmail := Fact{
		Fact: "test@gmail@gmail.com",
		T:    Email,
	}

	// Should not be able to verify user
	err = ValidateFact(invalidEmail)
	if err == nil {
		t.Errorf("Expected error in error path: should not be able to verify %s", invalidEmail.Fact)
	}
}

// Unit test for input validation of emails
func TestValidateFact_PhoneNumber(t *testing.T) {
	USCountryCode := "US"
	// UKCountryCode := "UK"
	// InvalidNumber := "020 8743 8000135"
	USNumber := "6502530000"

	// Valid Fact
	USFact := Fact{
		Fact: USNumber + USCountryCode,
		T:    Phone,
	}

	// Check US valid fact combination
	err := ValidateFact(USFact)
	if err != nil {
		t.Errorf("Unexpected error in happy path: %v", err)
	}

	// Phone number validation disabled
	// InvalidFact := Fact{
	// 	Fact: USNumber + UKCountryCode,
	// 	T:    Phone,
	// }
	//
	// // Invalid number and country code combination
	// err = ValidateFact(InvalidFact)
	// if err == nil {
	// 	t.Errorf("Expected error path: should not be able to validate US number with UK country code")
	// }
	//
	// InvalidFact = Fact{
	// 	Fact: InvalidNumber,
	// 	T:    Phone,
	// }
	// // Pass in an invalid number with a valid country code
	// err = ValidateFact(InvalidFact)
	// if err == nil {
	// 	t.Errorf("Expected error path: should not be able to validate US number with UK country code")
	// }
}
