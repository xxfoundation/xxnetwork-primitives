///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package excludedRounds

import (
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

func TestSet(t *testing.T) {
	xr := New()
	rid1 := id.Round(400)
	if xr.Has(rid1) {
		t.Errorf("New excluded rounds set should not have anything in it")
	}
	xr.Insert(rid1)
	if !xr.Has(rid1) {
		t.Errorf("Should have found inserted round in excluded round set")
	}
}
