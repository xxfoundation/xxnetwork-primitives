////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                           //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file                                                               //
////////////////////////////////////////////////////////////////////////////////

package id

import (
	"math/rand"
	"testing"
)

// Tests that random rounds marshalled via Round.Marshal and then unmarshalled
// via UnmarshalRound matches the original round.
func TestRound_Marshal_UnmarshalRound(t *testing.T) {
	prng := rand.New(rand.NewSource(42))

	for i := 0; i < 100; i++ {
		rid := Round(prng.Uint64())
		marshalledBytes := rid.Marshal()
		unmarshalledRid := UnmarshalRound(marshalledBytes)

		if rid != unmarshalledRid {
			t.Errorf("Marshalled and unmarshalled round ID does not match "+
				"original.\nexpected: %d\nreceived: %d", rid, unmarshalledRid)
		}
	}
}
