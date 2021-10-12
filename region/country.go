////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                           //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file                                                               //
////////////////////////////////////////////////////////////////////////////////

package region

// GetCountryBin return the bin for the given country alpha-2 code.
func GetCountryBin(countryCode string) (GeoBin, bool) {
	bin, exists := countryBins[countryCode]
	return bin, exists
}

// GetCountryList returns a list of all country alpha-2 codes.
func GetCountryList() []string {
	list := make([]string, 0, len(countryBins))

	for code := range countryBins {
		list = append(list, code)
	}

	return list
}

// GetCountryBins returns a copy of the countryBins map.
func GetCountryBins() map[string]GeoBin {
	// Create the target map
	targetMap := make(map[string]GeoBin)

	// Copy from the original map to the target map
	for key, value := range countryBins {
		targetMap[key] = value
	}
	return targetMap
}

// CountryLen returns the number of countries in the countryBins list.
func CountryLen() int {
	return len(countryBins)
}

// countryBins maps country alpha-2 codes to their regional bin.
var countryBins = map[string]GeoBin{
	"AQ": NorthAmerica,
	"MX": NorthAmerica,
}
