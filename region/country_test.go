////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package region

import (
	"testing"
)

// Unit test of GetCountryBin.
func TestGetCountryBin(t *testing.T) {
	for country, bin := range countryBins {
		foundBin, exists := GetCountryBin(country)
		if !exists {
			t.Errorf("GetCountryBin did not find %s in map.", country)
		}

		if bin != foundBin {
			t.Errorf("GetCountryBin did not return the expected bin"+
				"\nexpected: %s\nreceived: %s", bin, foundBin)
		}
	}
}

// Unit test of GetCountryList.
func TestGetCountryList(t *testing.T) {
	list := GetCountryList()
	countryBinsCopy := make(map[string]GeoBin, len(countryBins))
	for key, val := range countryBins {
		countryBinsCopy[key] = val
	}

	for i, code := range list {
		if _, exists := countryBinsCopy[code]; !exists {
			t.Errorf("Country code %q not found in map (%d).", code, i)
		}
		delete(countryBinsCopy, code)
	}

	if len(countryBinsCopy) != 0 {
		t.Errorf("Map contains %d entires not found in list.", len(countryBinsCopy))
	}
}

// Unit test of GetCountryBins.
func TestGetCountryBins(t *testing.T) {
	bins := GetCountryBins()

	for k, v := range bins {
		if val, exists := bins[k]; !exists || v != val {
			t.Errorf("Country code %q with %s not found in map.", k, v)
		}
	}
}

// Unit test of CountryLen.
func TestCountryLen(t *testing.T) {
	if len(countryBins) != CountryLen() {
		t.Errorf("CountryLen did not return the expected length."+
			"\nexpected: %d\nreceived: %d", len(countryBins), CountryLen())
	}
}
