package state

import "fmt"

//this holds the state enum for the state of the server.  It will be used for

// type which holds states so they have have an associated stringer
type State uint8

// List of states server can be in
const (
	NOT_STARTED = State(iota)
	WAITING
	PRECOMPUTING
	STANDBY
	REALTIME
	ERROR
	CRASH
)

const NUM_STATES = CRASH + 1

// Stringer to get the name of the state, primarily for for error prints
func (s State) String() string {
	switch s {
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
