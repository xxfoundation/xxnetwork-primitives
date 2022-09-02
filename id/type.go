////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package id

import (
	"fmt"
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

// String is a stringer to get the name of the ID type.
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
		return fmt.Sprintf("UNKNOWN ID TYPE: %d", t)
	}
}
