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

func (t FactType) IsValid() bool {
	return t == Username || t == Email || t == Phone || t == Nickname
}
