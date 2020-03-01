////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package current

import "fmt"

//this holds the enum for the activity of the server. It is in primitives so
//other repos such as registration/permissioning can access it

// type which holds activities so they have have an associated stringer
type Activity uint32

// List of Activities server can be in
const (
	NOT_STARTED = Activity(iota)
	WAITING
	PRECOMPUTING
	STANDBY
	REALTIME
	ERROR
	CRASH
)

const NUM_STATES = CRASH + 1

// Stringer to get the name of the activity, primarily for for error prints
func (a Activity) String() string {
	switch a {
	case NOT_STARTED:
		return "NOT_STARTED"
	case WAITING:
		return "WAITING"
	case PRECOMPUTING:
		return "PRECOMPUTING"
	case STANDBY:
		return "STANDBY"
	case REALTIME:
		return "REALTIME"
	case ERROR:
		return "ERROR"
	case CRASH:
		return "CRASH"
	default:
		return fmt.Sprintf("UNKNOWN STATE: %d", s)
	}
}
