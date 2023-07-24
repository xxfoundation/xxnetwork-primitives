////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package fact

import (
	"encoding/json"
	"reflect"
	"testing"
)

// Tests that a FactList marshalled by FactList.Stringify and unmarshalled by
// UnstringifyFactList matches the original.
func TestFactList_Stringify_UnstringifyFactList(t *testing.T) {
	expected := FactList{
		Fact{"vivian@elixxir.io", Email},
		Fact{"(270) 301-5797US", Phone},
		Fact{"invalidFact", Phone},
	}

	flString := expected.Stringify()
	factList, _, err := UnstringifyFactList(flString)
	if err != nil {
		t.Fatalf("Failed to unstringify %q: %+v", flString, err)
	}

	expected = expected[:2]
	if !reflect.DeepEqual(factList, expected) {
		t.Errorf("Unexpected unstringified FactList."+
			"\nexpected: %v\nreceived: %v", expected, factList)
	}
}

// Tests that a nil FactList marshalled by FactList.Stringify and unmarshalled
// by UnstringifyFactList matches the original.
func TestUnstringifyFactList_NilFactList(t *testing.T) {
	var expected FactList

	flString := expected.Stringify()
	factList, _, err := UnstringifyFactList(flString)
	if err != nil {
		t.Fatalf("Failed to unstringify %q: %+v", flString, err)
	}

	if !reflect.DeepEqual(factList, expected) {
		t.Errorf("Unexpected unstringified FactList."+
			"\nexpected: %v\nreceived: %v", expected, factList)
	}
}

// Error path: Tests that UnstringifyFactList returns an error for a malformed
// stringified FactList.
func Test_UnstringifyFactList_MissingFactBreakError(t *testing.T) {
	_, _, err := UnstringifyFactList("hi")
	if err == nil {
		t.Errorf("Expected error for invalid stringified list.")
	}
}

// Tests that a FactList JSON marshalled and unmarshalled matches the original.
func TestFactList_JsonMarshalUnmarshal(t *testing.T) {
	expected := FactList{
		{"devUsername", Username},
		{"devinputvalidation@elixxir.io", Email},
		{"6502530000US", Phone},
		{"name", Nickname},
	}

	data, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Failed to JSON marshal FactList: %+v", err)
	}

	var factList FactList
	err = json.Unmarshal(data, &factList)
	if err != nil {
		t.Errorf("Failed to JSON unmarshal FactList: %+v", err)
	}

	if !reflect.DeepEqual(expected, factList) {
		t.Errorf("Marshalled and unmarshalled FactList does not match original."+
			"\nexpected: %+v\nreceived: %+v", expected, factList)
	}
}
