////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package fact

import (
	"strings"

	"github.com/badoux/checkmail"
	"github.com/pkg/errors"
	"github.com/ttacon/libphonenumber"
)

const (
	// The maximum character length of a fact.
	maxFactLen = 64

	// The minimum character length of a nickname.
	minNicknameLen = 3
)

// Fact represents a piece of user-identifying information. This structure can
// be JSON marshalled and unmarshalled.
//
// JSON example:
//
//	{
//	  "Fact": "john@example.com",
//	  "T": 1
//	}
type Fact struct {
	Fact string   `json:"Fact"`
	T    FactType `json:"T"`
}

// NewFact checks if the inputted information is a valid fact on the
// fact type. If so, it returns a new fact object. If not, it returns a
// validation error.
func NewFact(ft FactType, fact string) (Fact, error) {
	if len(fact) > maxFactLen {
		return Fact{}, errors.Errorf("Fact (%s) exceeds maximum character limit"+
			"for a fact (%d characters)", fact, maxFactLen)
	}

	f := Fact{
		Fact: fact,
		T:    ft,
	}
	if err := ValidateFact(f); err != nil {
		return Fact{}, err
	}

	return f, nil
}

// Stringify marshals the Fact for transmission for UDB. It is not a part of the
// fact interface.
func (f Fact) Stringify() string {
	return f.T.Stringify() + f.Fact
}

// UnstringifyFact unmarshalls the stringified fact into a Fact.
func UnstringifyFact(s string) (Fact, error) {
	if len(s) < 1 {
		return Fact{}, errors.New("stringified facts must at least " +
			"have a type at the start")
	}

	if len(s) > maxFactLen {
		return Fact{}, errors.Errorf("Fact (%s) exceeds maximum character limit"+
			"for a fact (%d characters)", s, maxFactLen)
	}

	T := s[:1]
	fact := s[1:]
	if len(fact) == 0 {
		return Fact{}, errors.New(
			"stringified facts must be at least 1 character long")
	}
	ft, err := UnstringifyFactType(T)
	if err != nil {
		return Fact{}, errors.WithMessagef(err,
			"Failed to unstringify fact type for %q", s)
	}

	return NewFact(ft, fact)
}

// Normalized returns the fact in all uppercase letters.
func (f Fact) Normalized() string {
	return strings.ToUpper(f.Fact)
}

// ValidateFact checks the fact to see if it valid based on its type.
func ValidateFact(fact Fact) error {
	switch fact.T {
	case Username:
		return nil
	case Phone:
		// Extract specific information for validating a number
		// TODO: removes phone validation entirely. It is not used right now anyhow
		number, code := extractNumberInfo(fact.Fact)
		return validateNumber(number, code)
	case Email:
		// Check input of email inputted
		return validateEmail(fact.Fact)
	case Nickname:
		return validateNickname(fact.Fact)
	default:
		return errors.Errorf("Unknown fact type: %v", fact.T)
	}
}

// Numbers are assumed to have the 2-letter country code appended
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
	if err := checkmail.ValidateFormat(email); err != nil {
		return errors.Errorf(
			"Could not validate format for email [%s]: %v", email, err)
	}

	return nil
}

// Checks if the number and country code passed in is parse-able
// and is a valid phone number with that information
func validateNumber(number, countryCode string) error {
	errCh := make(chan error)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				errCh <- errors.Errorf("Crash occured on phone "+
					"validation of: number: %s, country code: %s", number,
					countryCode)
			}
		}()

		if len(number) == 0 || len(countryCode) == 0 {
			errCh <- errors.New("Number or input are of length 0")
		}
		num, err := libphonenumber.Parse(number, countryCode)
		if err != nil || num == nil {
			errCh <- errors.Errorf("Could not parse number [%s]: %v", number, err)
		}
		if !libphonenumber.IsValidNumber(num) {
			errCh <- errors.Errorf("Could not validate number [%s]: %v", number, err)
		}
		errCh <- nil
	}()
	return <-errCh
}

func validateNickname(nickname string) error {
	if len(nickname) < minNicknameLen {
		return errors.Errorf("Could not validate nickname %s: "+
			"too short (< %d characters)", nickname, minNicknameLen)
	}
	return nil
}
