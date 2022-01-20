///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package fact

import (
	"github.com/pkg/errors"
	"strings"
)

// maxFactCharacterLimit is the maximum character length of a fact.
const maxFactCharacterLimit = 64

type Fact struct {
	Fact string
	T    FactType
}

// NewFact checks if the inputted information is a valid fact on the
// fact type. If so, it returns a new fact object. If not, it returns a
// validation error.
func NewFact(ft FactType, fact string) (Fact, error) {

	if len(fact) > maxFactCharacterLimit {
		return Fact{}, errors.Errorf("Fact (%s) exceeds maximum character limit"+
			"for a fact (%d characters)", fact, maxFactCharacterLimit)
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

// marshal is for transmission for UDB, not a part of the fact interface
func (f Fact) Stringify() string {
	return f.T.Stringify() + f.Fact
}

func (f Fact) Normalized() string {
	return strings.ToUpper(f.Fact)
}

func UnstringifyFact(s string) (Fact, error) {
	if len(s) < 1 {
		return Fact{}, errors.New("stringified facts must at least " +
			"have a type at the start")
	}

	if len(s) > maxFactCharacterLimit {
		return Fact{}, errors.Errorf("Fact (%s) exceeds maximum character limit"+
			"for a fact (%d characters)", s, maxFactCharacterLimit)
	}

	T := s[:1]
	fact := s[1:]
	if len(fact) == 0 {
		return Fact{}, errors.New("stringified facts must be at " +
			"least 1 character long")
	}
	ft, err := UnstringifyFactType(T)
	if err != nil {
		return Fact{}, err
	}

	return NewFact(ft, fact)
}
