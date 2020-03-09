package states

import "fmt"

// this holds the enum for the states of a round. It is in primitives so
// other repos such as registration/permissioning, gateway, and client can
// access it

// type which holds activities so they have have an associated stringer
type Round uint32

// List of Activities server can be in
const (
	PENDING = Round(iota)
	PRECOMPUTING
	STANDBY
	REALTIME
	COMPLETED
	FAILED
)

const NUM_STATES = FAILED + 1

// Stringer to get the name of the activity, primarily for for error prints
func (s Round) String() string {
	switch s {
	case PENDING:
		return "PENDING"
	case PRECOMPUTING:
		return "PRECOMPUTING"
	case STANDBY:
		return "STANDBY"
	case REALTIME:
		return "REALTIME"
	case COMPLETED:
		return "COMPLETED"
	case FAILED:
		return "FAILED"
	default:
		return fmt.Sprintf("UNKNOWN STATE: %d", s)
	}
}
