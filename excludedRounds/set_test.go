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

func Test_NewSet(t *testing.T) {
	xr := NewSet()
	rid1 := id.Round(400)
	if xr.Has(rid1) {
		t.Errorf("NewSet excluded rounds set should not have anything in it")
	}
	xr.Insert(rid1)
	if !xr.Has(rid1) {
		t.Errorf("Should have found inserted round in excluded round set")
	}
}

func TestSet_Remove(t *testing.T) {
	xr := NewSet()
	rid1 := id.Round(400)
	xr.Insert(rid1)
	if !xr.Has(rid1) {
		t.Errorf("Should have found inserted round in excluded round set")
	}

	xr.Remove(rid1)
	if xr.Has(rid1) {
		t.Errorf("Remove should have removed round %d from set", rid1)
	}

}

func TestSet_Union(t *testing.T) {
	xr := NewSet()
	for i := 0; i < 10; i++ {
		xr.Insert(id.Round(i))
	}

	// Construct a 2nd set
	xr2 := NewSet()
	for i := 0; i < 10; i++ {
		xr2.Insert(id.Round(100 + i))
	}

	// Union the two sets
	union := xr2.Union(xr.xr)

	// Ensure union is the two sets combined
	for i := 0; i < 10; i++ {
		if !union.Has(id.Round(i)) {
			t.Errorf("Union should have placed round %d from set 1", id.Round(i))
		}

		if !union.Has(id.Round(100 + i)) {
			t.Errorf("Union should have placed round %d from set 2", id.Round(100+i))
		}
	}

}

func TestSet_Len(t *testing.T) {
	xr := NewSet()
	expected := 10
	for i := 1; i < expected; i++ { // Index at one to avoid off by one
		xr.Insert(id.Round(i))
	}

	if xr.Len() != expected {
		t.Errorf("Len did not produce expected length."+
			"\nExpected: %d"+
			"\nReceived: %d", expected, xr.Len())
	}
}
