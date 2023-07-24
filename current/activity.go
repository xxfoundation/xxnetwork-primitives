////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package current

import (
	"strconv"

	"github.com/pkg/errors"

	"gitlab.com/elixxir/primitives/states"
)

// This holds the enum for the activity of the server. It is in primitives so
// that other repos such as registration/permissioning can access it.

// Activity describes the activity a server has be doing.
type Activity uint32

// List of Activities.
const (
	NOT_STARTED = Activity(iota)
	WAITING
	PRECOMPUTING
	STANDBY
	REALTIME
	COMPLETED
	ERROR
	CRASH
	NUM_STATES
)

// String returns the string representation of the Activity. This functions
// adheres to the fmt.Stringer interface.
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
		return "UNKNOWN ACTIVITY: " + strconv.FormatUint(uint64(a), 10)
	}
}

// ConvertToRoundState converts an Activity to a valid round state or returns an
// error if it is invalid.
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
			"unable to convert activity %s (%d) to a valid state", a, a)
	}
}
