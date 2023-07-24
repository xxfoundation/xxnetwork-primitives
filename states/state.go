////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package states

import "strconv"

// This holds the enum for the states of a round. It is in primitives so
// other repos such as registration/permissioning, gateway, and client can
// access it

// Round describes the state of the round.
type Round uint32

// List of round states.
const (
	PENDING = Round(iota)
	PRECOMPUTING
	STANDBY
	QUEUED
	REALTIME
	COMPLETED
	FAILED
	NUM_STATES
)

// String returns the string representation of the Round state. This functions
// adheres to the fmt.Stringer interface.
func (r Round) String() string {
	switch r {
	case PENDING:
		return "PENDING"
	case PRECOMPUTING:
		return "PRECOMPUTING"
	case STANDBY:
		return "STANDBY"
	case QUEUED:
		return "QUEUED"
	case REALTIME:
		return "REALTIME"
	case COMPLETED:
		return "COMPLETED"
	case FAILED:
		return "FAILED"
	default:
		return "UNKNOWN STATE: " + strconv.FormatUint(uint64(r), 10)
	}
}
