///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package fact

import (
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
