package fact

import (
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/pkg/errors"
	"github.com/ttacon/libphonenumber"
	"regexp"
)

// Character limits for usernames.
const (
	minimumUsernameLength = 4
	maximumUsernameLength = 32
)

// usernameRegex is the regular expression for the enforcing the following characters only:
//  abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-+.@#
// Furthermore, the regex enforces usernames to begin and end with an alphanumeric character.
var usernameRegex = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_\\-+@.#]*[a-zA-Z0-9]$")

// Take the fact passed in and checks the input to see if it
//  valid based on the type of fact it is
func ValidateFact(fact Fact) error {
	switch fact.T {
	case Username:
		return validateUsername(fact.Fact)
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
	if err := checkmail.ValidateFormat(email); err != nil {
		return errors.Errorf("Could not validate format for email [%s]: %v", email, err)
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
	if len(nickname) < 3 {
		return errors.New(fmt.Sprintf("Could not validate nickname %s: too short (< 3 characters)", nickname))
	}
	return nil
}

func validateUsername(username string) error {
	// Check for acceptable length
	if len(username) < minimumUsernameLength || len(username) > maximumUsernameLength {
		return errors.Errorf("username length %d is not between %d and %d",
			len(username), minimumUsernameLength, maximumUsernameLength)
	}

	// Check is username contains allowed characters only
	if !usernameRegex.MatchString(username) {
		return errors.Errorf("username can only contain alphanumberics " +
			"and the symbols _, -, +, ., @, # and must start and end with an alphanumeric character")
	}

	return nil
}

