////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package excludedRounds

import (
	"testing"

	"gitlab.com/xx_network/primitives/id"
)

func TestSet(t *testing.T) {
	s := NewSet()
	if s.Len() != 1 {
		t.Errorf("Unexpected length.\nexpected: %d\nreceived: %d", 1, s.Len())
	}
	rid1 := id.Round(400)
	if s.Has(rid1) {
		t.Errorf("NewSet excluded rounds set should not have anything in it")
	}
	if !s.Insert(rid1) {
		t.Errorf("Insert failed.")
	}
	if s.Insert(rid1) {
		t.Errorf("Insert did not fail for already inserted item.")
	}
	if !s.Has(rid1) {
		t.Errorf("Should have found inserted round in excluded round set")
	}
	if s.Len() != 2 {
		t.Errorf("Unexpected length.\nexpected: %d\nreceived: %d", 2, s.Len())
	}
	s.Remove(rid1)
	if s.Has(rid1) {
		t.Errorf("Should not have found round in excluded round set")
	}
	if s.Len() != 1 {
		t.Errorf("Unexpected length.\nexpected: %d\nreceived: %d", 1, s.Len())
	}
}
