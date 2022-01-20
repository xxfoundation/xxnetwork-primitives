///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package fact

import "testing"

// Tests whether good usernames are considered valid.
func TestIsValidUsername_GoodUsernames(t *testing.T) {
	// Construct a list of good username
	goodUsernames := []string{
		"abcdefghijklmnopqrstuvwxyzABCDEF",
		"GHIJKLMNOPQRSTUVWXYZ0123456789",
		"john_doe",
		"daMan",
		"Mr.George",
		"josh-@+#b",
		"A........5",
	}

	// Test whether every good username is valid
	for _, goodUsername := range goodUsernames {
		err := validateUsername(goodUsername)
		if err != nil { // If invalid, fail test
			t.Errorf("IsValidUsername failed with username %q: %v", goodUsername, err)
		}
	}

}

// Tests whether invalid usernames are considered invalid.
func TestIsValidUsername_BadUsernames(t *testing.T) {
	// Construct a list of bad usernames
	badUsernames := []string{
		"",
		"  ",
		"pie",
		"123456789012345678901234567890123",
		"аdМіnіѕТrаТоr",
		"ÁdmïNIstrátör",
		"𝔞𝔡𝔪𝔦𝔫",
		"á̵͖͔͇̟͙̜͙̀͑̒̀̕ď̶̦̣̲m̴̡̬̺̯̩͂i̶͚͍̝̋ͅn̶̤̙̩͎̠͍̙̱̹̏",
		"ﬁnished",
		"GHIJKLMNOPQRSTUVWXYZ0123456789_-",
		"!@#$%^*?",
		"josh!!!!!",
	}

	// Test if every bad username is invalid
	for _, badUsername := range badUsernames {
		err := validateUsername(badUsername)
		if err == nil { // If considered valid, fail test
			t.Errorf("IsValidUsername did not fail with username %q", badUsername)
		}
	}

}

// Unit test for input validation of emails
// NOTE: Tests below might fail on goland due to bad internet connections
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
