////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package id

import (
	"encoding/binary"
	"strconv"
)

// Round is the round ID for each round run in cMix.
type Round uint64

// Marshal serialises the Round ID into a byte slice.
func (rid Round) Marshal() []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(rid))
	return b
}

// UnmarshalRound deserializes the byte slice into a Round ID.
func UnmarshalRound(b []byte) Round {
	return Round(binary.LittleEndian.Uint64(b))
}

// String returns the string representation of the Round ID. This functions
// adheres to the fmt.Stringer interface.
func (rid Round) String() string {
	return strconv.FormatUint(uint64(rid), 10)
}
