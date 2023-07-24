////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

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

// String returns the ID Type in a human-readable form for use in logging and
// debugging. This functions adheres to the fmt.Stringer interface.
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
		return "UNKNOWN ID TYPE: " + strconv.Itoa(int(t))
	}
}
