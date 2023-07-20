////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package fact

import (
	"fmt"

	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
)

type FactType uint8

const (
	Username FactType = 0
	Email    FactType = 1
	Phone    FactType = 2
	Nickname FactType = 3
)

// String returns the string representation of the FactType. This functions
// adheres to the fmt.Stringer interface.
func (t FactType) String() string {
	switch t {
	case Username:
		return "Username"
	case Email:
		return "Email"
	case Phone:
		return "Phone"
	case Nickname:
		return "Nickname"
	default:
		return fmt.Sprintf("Unknown Fact FactType: %d", t)
	}
}

// Stringify marshals the FactType into a portable string.
func (t FactType) Stringify() string {
	switch t {
	case Username:
		return "U"
	case Email:
		return "E"
	case Phone:
		return "P"
	case Nickname:
		return "N"
	}
	jww.FATAL.Panicf("Unknown Fact FactType: %d", t)
	return "error"
}

// UnstringifyFactType unmarshalls the stringified FactType.
func UnstringifyFactType(s string) (FactType, error) {
	switch s {
	case "U":
		return Username, nil
	case "E":
		return Email, nil
	case "P":
		return Phone, nil
	case "N":
		return Nickname, nil
	}
	return 3, errors.Errorf("Unknown Fact FactType: %s", s)
}

// IsValid determines if the FactType is one of the defined types.
func (t FactType) IsValid() bool {
	return t == Username || t == Email || t == Phone || t == Nickname
}
