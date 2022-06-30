////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                           //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file                                                               //
////////////////////////////////////////////////////////////////////////////////

package exponential

import (
	"math/rand"
	"reflect"
	"testing"
)

// Tests that NewMovingAvg returns a new MovingAvg with the expected
// values.
func TestNewMovingAvg(t *testing.T) {
	expected := &MovingAvg{
		cutoff: 0.23,
		aN:     0.5,
		s:      4,
		e:      1000,
	}
	p := MovingAvgParams{
		Cutoff:          expected.cutoff,
		InitialAverage:  expected.aN,
		SmoothingFactor: expected.s,
		NumberOfEvents:  expected.e,
	}

	ea := NewMovingAvg(p)

	if !reflect.DeepEqual(expected, ea) {
		t.Errorf("Received unexpected MovingAvg."+
			"\nexpected: %+v\nreceived: %+v", expected, ea)
	}
}

// Tests that MovingAvg.Intake does not return an error.
// NOTE: This is not a full or accurate test of MovingAvg.Intake. Not sure if
//  there is a good test for it, but if there is, you should add it.
func TestMovingAvg_Intake(t *testing.T) {
	ea := NewMovingAvg(DefaultMovingAvgParams())

	for i := 0; i < int(ea.e); i++ {
		err := ea.Intake(BoolToFloat(i%2 == 0))
		if err != nil {
			t.Errorf("Error on instake #%d: %+v", i, err)
			break
		}
	}
}

// Tests that MovingAvg.IsOverCutoff returns false when the cutoff has not
// been reach and true when it has been reached
func TestMovingAvg_IsOverCutoff(t *testing.T) {
	ea := NewMovingAvg(DefaultMovingAvgParams())
	prng := rand.New(rand.NewSource(42))

	if ea.IsOverCutoff() {
		t.Errorf("IsOverCutoff reported that the cutoff has been reached "+
			"when it has not.\ncutoff:  %f\naverage: %f", ea.cutoff, ea.aN)
	}

	var err error
	for err == nil {
		err = ea.Intake(BoolToFloat(prng.Uint64()%2 == 0))
	}

	if !ea.IsOverCutoff() {
		t.Errorf("IsOverCutoff reported that the cutoff has not been reached "+
			"when it has.\ncutoff:  %f\naverage: %f", ea.cutoff, ea.aN)
	}

}

// Tests both cases for BoolToFloat.
func Test_BoolToFloat(t *testing.T) {
	if BoolToFloat(true) != 1 {
		t.Errorf("Received incorrect float for boolean %t."+
			"\nexpected: %f\nreceived: %f", true, float32(1), BoolToFloat(true))
	}

	if BoolToFloat(false) != 0 {
		t.Errorf("Received incorrect float for boolean %t."+
			"\nexpected: %f\nreceived: %f", false, float32(0), BoolToFloat(false))
	}
}
