////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
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
func (a Activity) ConvertToRoundState() (states.Round, error) {
	switch a {
	case WAITING:
		return states.PENDING, nil
	case PRECOMPUTING:
		return states.PRECOMPUTING, nil
	case STANDBY:
		return states.STANDBY, nil
	case REALTIME:
		return states.REALTIME, nil
	case COMPLETED:
		return states.COMPLETED, nil
	case ERROR:
		return states.FAILED, nil
	default:
		// Unsupported conversion. Return an arbitrary round and error
		return states.Round(99), errors.Errorf(
			"unable to convert activity %+v to valid state", a)
	}
}
