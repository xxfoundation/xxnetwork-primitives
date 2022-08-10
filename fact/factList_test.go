///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package fact

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFactList_StringifyUnstringify(t *testing.T) {
	expected := FactList{}
	expected = append(expected, Fact{
		Fact: "vivian@elixxir.io",
		T:    Email,
	})
	expected = append(expected, Fact{
		Fact: "(270) 301-5797US",
		T:    Phone,
	})

	FlString := expected.Stringify()
	// Manually check and verify that the string version is as expected
	t.Log(FlString)

	actual, _, err := UnstringifyFactList(FlString)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Error("fact lists weren't equal")
	}
}

// Tests that a FactList can be JSON marshalled and unmarshalled.
func TestFactList_JSON(t *testing.T) {
	fl := FactList{
		{"devUsername", Username},
		{"devinputvalidation@elixxir.io", Email},
		{"6502530000US", Phone},
		{"name", Nickname},
	}

	out, err := json.Marshal(fl)
	if err != nil {
		t.Errorf("Failed to marshal FactList: %+v", err)
	}

	var newFactList FactList
	err = json.Unmarshal(out, &newFactList)
	if err != nil {
		t.Errorf("Failed to unmarshal FactList: %+v", err)
	}

	if !reflect.DeepEqual(fl, newFactList) {
		t.Errorf("Marshalled and unmarshalled FactList does not match original."+
			"\nexpected: %+v\nreceived: %+v", fl, newFactList)
	}
}
