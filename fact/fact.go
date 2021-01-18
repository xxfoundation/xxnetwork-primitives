///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package fact

import (
	"github.com/badoux/checkmail"
	"github.com/nyaruka/phonenumbers"
	"github.com/pkg/errors"
)

type Fact struct {
	Fact string
	T    FactType
}

// NewFact checks if the inputted information is a valid fact on the
// fact type. If so, it returns a new fact object. If not, it returns a
// validation error.
func NewFact(ft FactType, fact string) (Fact, error) {

	f := Fact{
		Fact: fact,
		T:    ft,
	}
	if err := ValidateFact(f); err != nil {
		return Fact{}, err
	}

	return f, nil
}

// marshal is for transmission for UDB, not a part of the fact interface
func (f Fact) Stringify() string {
	return f.T.Stringify() + f.Fact
}

func UnstringifyFact(s string) (Fact, error) {
	if len(s) < 1 {
		return Fact{}, errors.New("stringified facts must at least have a type at the start")
	}
	T := s[:1]
	fact := s[1:]
	ft, err := UnstringifyFactType(T)
	if err != nil {
		return Fact{}, err
	}

	return NewFact(ft, fact)
}

// Take the fact passed in and checks the input to see if it
//  valid based on the type of fact it is
func ValidateFact(fact Fact) error {
	switch fact.T {
	case Username:
		return nil
	case Phone:
		// Extract specific information for validating a number
		number, code := extractNumberInfo(fact.Fact)
		err := validateNumber(number, code)
		if err != nil {
			return err
		}
		return nil
	case Email:
		// Check input of email inputted
		err := validateEmail(fact.Fact)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.Errorf("Unknown fact type: %v", fact.T)

	}

}

// Numbers are assumed to have the 2 letter country code appended
// to the fact, with the rest of the information being a phone number
// Example: 6502530000US is a valid US number with the country code
// that would be the fact information for a phone number
func extractNumberInfo(fact string) (number, countryCode string) {
	factLen := len(fact)
	number = fact[:factLen-2]
	countryCode = fact[factLen-2:]
	return
}

// Validate the email input and check if the host is contact-able
func validateEmail(email string) error {
	// Check that the input is validly formatted
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return errors.Errorf("Could not validate format for email [%s]: %v", email, err)
	}

	return nil
}

// Checks if the number and country code passed in is parse-able
// and is a valid phone number with that information
func validateNumber(number, countryCode string) error {
	num, err := phonenumbers.Parse(number, countryCode)
	if err != nil {
		return errors.Errorf("Could not parse number [%s]: %v", number, err)
	}
	if !phonenumbers.IsValidNumber(num) {
		return errors.Errorf("Could not validate number [%s]: %v", number, err)
	}
	return nil
}
