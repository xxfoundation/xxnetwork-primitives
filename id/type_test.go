////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package id

import (
	"testing"
)

// Tests that String() returns the correct string for each Type.
func TestType_String(t *testing.T) {
	// Expected values
	expectedGeneric := "generic"
	expectedGateway := "gateway"
	expectedNode := "node"
	expectedUser := "user"
	expectedGroup := "group"
	expectedNumTypes := "5"

	// Test Generic stringer
	testVal := Generic.String()
	if expectedGeneric != testVal {
		t.Errorf("String() returned incorrect string for Generic type."+
			"\n\texpected: %s\n\treceived: %s", expectedGeneric, testVal)
	}

	// Test Gateway stringer
	testVal = Gateway.String()
	if expectedGateway != testVal {
		t.Errorf("String() returned incorrect string for Gateway type."+
			"\n\texpected: %s\n\treceived: %s", expectedGateway, testVal)
	}

	// Test Node stringer
	testVal = Node.String()
	if expectedNode != testVal {
		t.Errorf("String() returned incorrect string for Node type."+
			"\n\texpected: %s\n\treceived: %s", expectedNode, testVal)
	}

	// Test User stringer
	testVal = User.String()
	if expectedUser != testVal {
		t.Errorf("String() returned incorrect string for User type."+
			"\n\texpected: %s\n\treceived: %s", expectedUser, testVal)
	}

	// Test Group stringer
	testVal = Group.String()
	if expectedGroup != testVal {
		t.Errorf("String() returned incorrect string for Group type."+
			"\n\texpected: %s\n\treceived: %s", expectedGroup, testVal)
	}

	// Test NumTypes stringer
	testVal = NumTypes.String()
	if expectedNumTypes != testVal {
		t.Errorf("String() returned incorrect string for NumTypes type."+
			"\n\texpected: %s\n\treceived: %s", expectedNumTypes, testVal)
	}
}

// Tests that String() returns an error when given a Type that has not been
// defined.
func TestType_String_Error(t *testing.T) {
	// Expected/test values
	expectedError := "UNKNOWN ID TYPE: 6"
	testType := Type(6)

	// Test stringer error
	testVal := testType.String()
	if expectedError != testVal {
		t.Errorf("String() did not return an error when it should have."+
			"\n\texpected: %s\n\treceived: %s", expectedError, testVal)
	}
}
