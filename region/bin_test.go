////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package region

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
)

// Unit test of GeoBin.String.
func TestGeoBin_String(t *testing.T) {
	testValues := []struct {
		bin      GeoBin
		expected string
	}{
		{NorthAmerica, "NorthAmerica"},
		{SouthAndCentralAmerica, "SouthAndCentralAmerica"},
		{WesternEurope, "WesternEurope"},
		{CentralEurope, "CentralEurope"},
		{EasternEurope, "EasternEurope"},
		{MiddleEast, "MiddleEast"},
		{NorthernAfrica, "NorthernAfrica"},
		{SouthernAfrica, "SouthernAfrica"},
		{Russia, "Russia"},
		{EasternAsia, "EasternAsia"},
		{WesternAsia, "WesternAsia"},
		{Oceania, "Oceania"},

		{Oceania + 1, "INVALID BIN " + strconv.Itoa(int(Oceania+1))},
	}

	for i, val := range testValues {
		if val.bin.String() != val.expected {
			t.Errorf("String did not return the expected string (%d)."+
				"\nexpected: %s\nreceived: %s", i, val.expected, val.bin)
		}
	}
}

// Unit test of GetRegion.
func TestGetRegion(t *testing.T) {
	testValues := []struct {
		region   string
		expected GeoBin
	}{
		{"NorthAmerica", NorthAmerica},
		{"SouthAndCentralAmerica", SouthAndCentralAmerica},
		{"WesternEurope", WesternEurope},
		{"CentralEurope", CentralEurope},
		{"EasternEurope", EasternEurope},
		{"MiddleEast", MiddleEast},
		{"NorthernAfrica", NorthernAfrica},
		{"SouthernAfrica", SouthernAfrica},
		{"Russia", Russia},
		{"EasternAsia", EasternAsia},
		{"WesternAsia", WesternAsia},
		{"Oceania", Oceania},
	}

	for i, val := range testValues {
		bin, err := GetRegion(val.region)
		if err != nil {
			t.Errorf("GetRegion returned an error (%d): %+v", i, err)
		}

		if bin != val.expected {
			t.Errorf("GetRegion did not return the expected value (%d)."+
				"\nexpected: %d\nreceived: %d", i, val.expected, bin)
		}
	}
}

// Error path: tests that GetRegion returns an error for an invalid
func TestGetRegion_InvalidRegionError(t *testing.T) {
	region := "INVALID REGION"
	expectedErr := fmt.Sprintf(invalidRegionErr, region)

	bin, err := GetRegion(region)
	if err == nil || err.Error() != expectedErr {
		t.Errorf("GetRegion did not return the expected error."+
			"\nexpected: %s\nreceived: %+v", expectedErr, err)
	}

	if bin != math.MaxUint8 {
		t.Errorf("GetRegion did not return the expected value."+
			"\nexpected: %d\nreceived: %d", math.MaxUint8, bin)
	}
}

// Tests that a GeoBin can be JSON marshalled and unmarshalled.
func TestGeoBin_JsonMarshalUnmarshal(t *testing.T) {
	bin := NorthAmerica

	data, err := json.Marshal(bin)
	if err != nil {
		t.Errorf("Failed to JSON marshal GeoBin: %+v", err)
	}

	var newBin GeoBin
	err = json.Unmarshal(data, &newBin)
	if err != nil {
		t.Errorf("Failed to JSON unmarshal GeoBin: %+v", err)
	}

	if bin != newBin {
		t.Errorf("JSON marshalled and unmarshalled GeoBin does not match original."+
			"\nexpected: %s\nreceived: %s", bin, newBin)
	}
}

// Unit test of GeoBin.Bytes.
func TestGeoBin_Bytes(t *testing.T) {
	bin := NorthAmerica
	expected := []byte{byte(bin)}

	if !bytes.Equal(expected, bin.Bytes()) {
		t.Errorf("Bytes did not return the expected value."+
			"\nexpected: %+v\nreceived: %+v", expected, bin.Bytes())
	}
}

// Error path: tests that GeoBin.UnmarshalJSON returns an error for invalid JSON
// data
func TestGeoBin_UnmarshalJSON_InvalidJsonError(t *testing.T) {
	expectedErr := strings.Split(jsonUnmarshalErr, "%")[0]

	var bin GeoBin
	err := bin.UnmarshalJSON([]byte("}{"))
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("UnmarshalJSON did not return the expected error."+
			"\nexpected: %s\nreceived: %+v", expectedErr, err)
	}
}

// Error path: tests that GeoBin.UnmarshalJSON returns an error when the JSON
// unmarshalled string does not match a region.
func TestGeoBin_UnmarshalJSON_InvalidRegionError(t *testing.T) {
	region := "INVALID REGION"
	expectedErr := fmt.Sprintf(invalidRegionErr, region)

	var bin GeoBin
	err := bin.UnmarshalJSON([]byte("\"" + region + "\""))
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("UnmarshalJSON did not return the expected error."+
			"\nexpected: %s\nreceived: %+v", expectedErr, err)
	}
}
