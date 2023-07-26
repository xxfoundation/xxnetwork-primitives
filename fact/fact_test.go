////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package fact

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// Tests that NewFact returns a correctly formatted Fact.
func TestNewFact(t *testing.T) {
	tests := []struct {
		ft       FactType
		fact     string
		expected Fact
	}{
		{Username, "myUsername", Fact{"myUsername", Username}},
		{Email, "email@example.com", Fact{"email@example.com", Email}},
		{Phone, "8005559486US", Fact{"8005559486US", Phone}},
		{Nickname, "myNickname", Fact{"myNickname", Nickname}},
	}

	for i, tt := range tests {
		fact, err := NewFact(tt.ft, tt.fact)
		if err != nil {
			t.Errorf("Failed to make new fact (%d): %+v", i, err)
		} else if !reflect.DeepEqual(tt.expected, fact) {
			t.Errorf("Unexpected new Fact (%d).\nexpected: %s\nreceived: %s",
				i, tt.expected, fact)
		}
	}
}

// Error path: Tests that NewFact returns error when a fact exceeds the
// maxFactLen.
func TestNewFact_ExceedMaxFactError(t *testing.T) {
	_, err := NewFact(Email,
		"devinputvalidation_devinputvalidation_devinputvalidation@elixxir.io")
	if err == nil {
		t.Fatal("Expected error when the fact is longer than the maximum " +
			"character length.")
	}

}

// Error path: Tests that NewFact returns error when the fact is not valid.
func TestNewFact_InvalidFactError(t *testing.T) {
	_, err := NewFact(Nickname, "hi")
	if err == nil {
		t.Fatal("Expected error when the fact is invalid.")
	}
}

// Tests that a Fact marshalled by Fact.Stringify and unmarshalled by
// UnstringifyFact matches the original.
func TestFact_Stringify_UnstringifyFact(t *testing.T) {
	facts := []Fact{
		{"myUsername", Username},
		{"email@example.com", Email},
		{"8005559486US", Phone},
		{"myNickname", Nickname},
	}

	for i, expected := range facts {
		factString := expected.Stringify()
		fact, err := UnstringifyFact(factString)
		if err != nil {
			t.Errorf(
				"Failed to unstringify fact %s (%d): %+v", expected, i, err)
		} else if !reflect.DeepEqual(expected, fact) {
			t.Errorf("Unexpected unstringified Fact %s (%d)."+
				"\nexpected: %s\nreceived: %s",
				factString, i, expected, fact)
		}
	}
}

// Consistency test of Fact.Stringify.
func TestFact_Stringify(t *testing.T) {
	tests := []struct {
		fact     Fact
		expected string
	}{
		{Fact{"myUsername", Username}, "UmyUsername"},
		{Fact{"email@example.com", Email}, "Eemail@example.com"},
		{Fact{"8005559486US", Phone}, "P8005559486US"},
		{Fact{"myNickname", Nickname}, "NmyNickname"},
	}

	for i, tt := range tests {
		factString := tt.fact.Stringify()

		if factString != tt.expected {
			t.Errorf("Unexpected strified Fact %s (%d)."+
				"\nexpected: %s\nreceived: %s",
				tt.fact, i, tt.expected, factString)
		}
	}
}

// Consistency test of UnstringifyFact
func TestUnstringifyFact(t *testing.T) {
	tests := []struct {
		factString string
		expected   Fact
	}{
		{"UmyUsername", Fact{"myUsername", Username}},
		{"Eemail@example.com", Fact{"email@example.com", Email}},
		{"P8005559486US", Fact{"8005559486US", Phone}},
		{"NmyNickname", Fact{"myNickname", Nickname}},
	}

	for i, tt := range tests {
		fact, err := UnstringifyFact(tt.factString)
		if err != nil {
			t.Errorf(
				"Failed to unstringify fact %s (%d): %+v", tt.factString, i, err)
		} else if !reflect.DeepEqual(tt.expected, fact) {
			t.Errorf("Unexpected unstringified Fact %s (%d)."+
				"\nexpected: %s\nreceived: %s",
				tt.factString, i, tt.expected, fact)
		}
	}
}

// Error path: Tests all error paths of UnstringifyFact.
func TestUnstringifyFact_Error(t *testing.T) {
	longFact := strings.Repeat("A", maxFactLen+1)
	tests := []struct {
		factString  string
		expectedErr string
	}{
		{"", "stringified facts must at least have a type at the start"},
		{longFact, fmt.Sprintf("Fact (%s) exceeds maximum character limit for "+
			"a fact (%d characters)", longFact, maxFactLen)},
		{"P", "stringified facts must be at least 1 character long"},
		{"QA", `Failed to unstringify fact type for "QA"`},
	}

	for i, tt := range tests {
		_, err := UnstringifyFact(tt.factString)
		if err == nil || !strings.Contains(err.Error(), tt.expectedErr) {
			t.Errorf("Unexpected error when Unstringifying fact %q (%d)."+
				"\nexpected: %s\nreceived: %+v",
				tt.factString, i, tt.expectedErr, err)
		}
	}
}

// Consistency test of Fact.Normalized.
func TestFact_Normalized(t *testing.T) {
	tests := []struct {
		fact     Fact
		expected string
	}{
		{Fact{"myUsername", Username}, "MYUSERNAME"},
		{Fact{"email@example.com", Email}, "EMAIL@EXAMPLE.COM"},
		{Fact{"8005559486US", Phone}, "8005559486US"},
		{Fact{"myNickname", Nickname}, "MYNICKNAME"},
	}

	for i, tt := range tests {
		normal := tt.fact.Normalized()
		if normal != tt.expected {
			t.Errorf("Unexpected new normalized Fact %v (%d)."+
				"\nexpected: %q\nreceived: %q", tt.fact, i, tt.expected, normal)
		}
	}
}

// Tests that ValidateFact correctly validates various facts.
func TestValidateFact(t *testing.T) {
	facts := []Fact{
		{"myUsername", Username},
		{"email@example.com", Email},
		{"8005559486US", Phone},
		{"myNickname", Nickname},
	}

	for i, fact := range facts {
		err := ValidateFact(fact)
		if err != nil {
			t.Errorf(
				"Failed to validate fact %s (%d): %+v", fact, i, err)
		}
	}
}

// Error path: Tests that ValidateFact does not validate invalid facts
func TestValidateFact_InvalidFactsError(t *testing.T) {
	facts := []Fact{
		{"test@gmail@gmail.com", Email},
		{"US8005559486", Phone},
		{"020 8743 8000135UK", Phone},
		{"me", Nickname},
		{"me", 99},
	}

	for i, fact := range facts {
		err := ValidateFact(fact)
		if err == nil {
			t.Errorf("Did not error on invalid fact %s (%d)", fact, i)
		}
	}
}

// Error path: Tests all error paths of validateNumber.
func Test_validateNumber_Error(t *testing.T) {
	tests := []struct {
		number, countryCode string
		expectedErr         string
	}{
		{"5", "", "Number or input are of length 0"},
		{"", "US", "Number or input are of length 0"},
		// {"020 8743 8000135", "UK", `Could not parse number "020 8743 8000135"`},
		{"8005559486", "UK", `Could not parse number "8005559486"`},
		{"+343511234567", "ES", `Could not validate number "+343511234567"`},
	}

	for i, tt := range tests {
		err := validateNumber(tt.number, tt.countryCode)
		if err == nil || !strings.Contains(err.Error(), tt.expectedErr) {
			t.Errorf("Unexpected error when validating number %q with country "+
				"code %q (%d).\nexpected: %s\nreceived: %+v",
				tt.number, tt.countryCode, i, tt.expectedErr, err)
		}
	}
}

// Tests that a Fact JSON marshalled and unmarshalled matches the original.
func TestFact_JsonMarshalUnmarshal(t *testing.T) {
	facts := []Fact{
		{"myUsername", Username},
		{"email@example.com", Email},
		{"8005559486US", Phone},
		{"myNickname", Nickname},
	}

	for i, expected := range facts {
		data, err := json.Marshal(expected)
		if err != nil {
			t.Errorf("Failed to JSON marshal %s (%d): %+v", expected, i, err)
		}

		var fact Fact
		if err = json.Unmarshal(data, &fact); err != nil {
			t.Errorf("Failed to JSON unmarshal %s (%d): %+v", expected, i, err)
		}

		if !reflect.DeepEqual(expected, fact) {
			t.Errorf("Unexpected unmarshalled fact (%d)."+
				"\nexpected: %+v\nreceived: %+v", i, expected, fact)
		}
	}
}
