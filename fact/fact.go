package fact

import (
	"github.com/badoux/checkmail"
	"github.com/pkg/errors"
)

type Fact struct {
	Fact string
	T    FactType
}

func NewFact(ft FactType, fact string) (Fact, error) {
	//todo: filter the fact string
	err := ValidateFact(ft, fact)
	if err != nil {
		return Fact{}, err
	}
	return Fact{
		Fact: fact,
		T:    ft,
	}, nil
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

func ValidateFact(ft FactType, fact string) error {
	switch ft {
		case Phone:
			err := validateNumber(fact)
			if err != nil {
				return err
			}
			return nil
		case Email:
			err := validateEmail(fact)
			if err != nil {
				return err
			}
			return nil
		default:
			return errors.Errorf("Unknown fact type: %v", ft)

	}

}

//todo: we need more information passed in, namely a country code
// look up documentation here: https://docs.google.com/document/d/1_CdhcKaKXI-BBwjWVUsavmI-fGi46RSZoDeZTPx-SBQ/edit#
func validateEmail(fact string) error{
	// Check that the input is validly formatted
	err := checkmail.ValidateFormat(fact)
	if err != nil {
		return errors.Errorf("Could not validate format for email: %s", fact)
	}

	// Check that the domain is valid and reachable
	err = checkmail.ValidateHost(fact)
	if err != nil {
		return errors.Errorf("Could not validate host for email: %s", fact)
	}
	return nil
}

func validateNumber(fact string)  error  {
	// fixme: need standardized way to get country code. Either also passed,
	// or concat with fact, but needs to be parsed out
	// OR use Twilio directly here, however would need auth key at a low level somehow
	return nil
}