////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package id

import (
	"strconv"
)

// Type holds the numerical representation of the ID type.
type Type byte

// List of ID types
const (
	Generic = Type(iota)
	Gateway
	Node
	User
	Group
	NumTypes // Gives number of ID types
)

const noIDTypeErr = "unknown ID type: "

// String returns the name of the ID type. This functions satisfies the
// fmt.Stringer interface.
func (t Type) String() string {
	switch t {
	case Generic:
		return "generic"
	case Gateway:
		return "gateway"
	case Node:
		return "node"
	case User:
		return "user"
	case Group:
		return "group"
	case NumTypes:
		return strconv.Itoa(int(NumTypes))
	default:
		return noIDTypeErr + strconv.Itoa(int(t))
	}
}
