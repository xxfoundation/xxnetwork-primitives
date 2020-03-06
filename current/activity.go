////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package current

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/primitives/states"
)

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
	COMPLETED
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
	case COMPLETED:
		return "COMPLETED"
	case ERROR:
		return "ERROR"
	case CRASH:
		return "CRASH"
	default:
		return fmt.Sprintf("UNKNOWN STATE: %d", a)
	}
}

// Converts an Activity to a valid Round state, or returns an error if invalid
func (a Activity) Convert() (states.Round, error) {
	if a <= PRECOMPUTING {
		return states.PRECOMPUTING, nil
	} else if a == STANDBY {
		return states.STANDBY, nil
	} else if a == REALTIME {
		return states.REALTIME, nil
	} else if a == COMPLETED {
		return states.COMPLETED, nil
	} else if a > COMPLETED {
		return states.FAILED, nil
	} else {
		return states.Round(0), errors.Errorf(
			"unable to convert activity %+v to valid state", a)
	}
}
