////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package exponential

import (
	"encoding/json"
	"reflect"
	"testing"
)

// Tests that DefaultMovingAvgParams returns a MovingAvgParams with all the
// default values.
func TestDefaultMovingAvgParams(t *testing.T) {
	expected := MovingAvgParams{
		Cutoff:          defaultCutoff,
		InitialAverage:  defaultInitialAverage,
		SmoothingFactor: defaultSmoothingFactor,
		NumberOfEvents:  defaultNumberOfEvents,
	}
	p := DefaultMovingAvgParams()

	if !reflect.DeepEqual(expected, p) {
		t.Errorf("Did not received expected default parameters."+
			"\nexpected: %+v\nreceived: %+v", expected, p)
	}
}

// Tests that MovingAvgParams can be JSON marshalled and unmarshalled.
func TestMovingAvgParams_JSONMarshalUnmarshal(t *testing.T) {
	p := DefaultMovingAvgParams()

	data, err := json.Marshal(p)
	if err != nil {
		t.Errorf("Failed to JSON marshal the MovingAvgParams: %+v", err)
	}

	var newParams MovingAvgParams
	err = json.Unmarshal(data, &newParams)
	if err != nil {
		t.Errorf("Failed to JSON unmarshal the MovingAvgParams: %+v", err)
	}

	if !reflect.DeepEqual(p, newParams) {
		t.Errorf("Marshalled and unmarshalled params do not match original."+
			"\nexpected: %+v\nreceived: %+v", p, newParams)
	}
}
