////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package id

import (
	"testing"
)

// Tests that Type.String returns the correct string for each Type.
func TestType_String(t *testing.T) {
	testValues := map[Type]string{
		Generic:  "generic",
		Gateway:  "gateway",
		Node:     "node",
		User:     "user",
		Group:    "group",
		NumTypes: "5",
	}

	for idType, expected := range testValues {
		if expected != idType.String() {
			t.Errorf("String returned incorrect string for type."+
				"\nexpected: %s\nreceived: %s", expected, idType.String())
		}
	}
}

// Tests that Type.String returns an error when given a Type that has not been
// defined.
func TestType_String_Error(t *testing.T) {
	expectedError := "UNKNOWN ID TYPE: 6"
	testType := Type(6)

	// Test stringer error
	testVal := testType.String()
	if expectedError != testVal {
		t.Errorf("String did not return an error when it should have."+
			"\nexpected: %s\nreceived: %s", expectedError, testVal)
	}
}
