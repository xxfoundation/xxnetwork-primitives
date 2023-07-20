////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package version

import (
	"reflect"
	"strings"
	"testing"
)

// Happy path.
func TestNew(t *testing.T) {
	expected := Version{1, 2, "3"}

	test := New(expected.major, expected.minor, expected.patch)

	if !reflect.DeepEqual(expected, test) {
		t.Errorf("New did not create the expected Version."+
			"\nexpected: %+v\nreceived: %+v", expected, test)
	}
}

// Happy path.
func TestVersion_Major(t *testing.T) {
	v := Version{1, 2, "3"}

	if v.Major() != v.major {
		t.Errorf("Major did not return the expected value."+
			"\nexpected: %d\nreceived: %d", v.major, v.Major())
	}
}

// Happy path.
func TestVersion_Minor(t *testing.T) {
	v := Version{1, 2, "3"}

	if v.Minor() != v.minor {
		t.Errorf("Minor did not return the expected value."+
			"\nexpected: %d\nreceived: %d", v.minor, v.Minor())
	}
}

// Happy path.
func TestVersion_Patch(t *testing.T) {
	v := Version{1, 2, "3"}

	if v.Patch() != v.patch {
		t.Errorf("Patch did not return the expected value."+
			"\nexpected: %s\nreceived: %s", v.patch, v.Patch())
	}
}

// Happy path.
func TestVersion_String(t *testing.T) {
	// Test values
	testValues := map[string]Version{
		"0.0.":           {},
		"0.0.0":          {0, 0, "0"},
		"5.10.test":      {5, 10, "test"},
		"5.10.test.test": {5, 10, "test.test"},
	}

	for expected, ver := range testValues {
		test := ver.String()
		if expected != test {
			t.Errorf("String did not return the expected string."+
				"\nexpected: %s\nreceived: %s", expected, test)
		}
	}
}

// Happy path.
func TestParseVersion(t *testing.T) {
	// Test values
	testValues := map[string]Version{
		"0.0.0":          {0, 0, "0"},
		"5.10.test":      {5, 10, "test"},
		"5.10.test.test": {5, 10, "test.test"},
	}

	for versionString, expected := range testValues {
		test, err := ParseVersion(versionString)
		if err != nil {
			t.Errorf("ParseVersion produced an unexpected error: %+v", err)
		}

		if expected != test {
			t.Errorf("ParseVersion did not return the expected Version."+
				"\nexpected: %+v\nreceived: %+v", expected, test)
		}
	}
}

// Error path: tests various invalid Version strings return the expected error.
func TestParseVersion_Error(t *testing.T) {
	// Test values
	testStrings := map[string]string{
		"invalid version.":         "3 parts",
		"":                         "3 parts",
		"0":                        "3 parts",
		"0.":                       "3 parts",
		"0.0":                      "3 parts",
		"0.0.":                     "patch cannot be empty",
		"a.0.0":                    "major version",
		"0.a.0":                    "minor version",
		"a.a.a":                    "major version",
		".":                        "3 parts",
		"..":                       "major version",
		"...":                      "major version",
		"18446744073709551615.0.0": "major",
	}

	// Check that all the test strings produce the expected errors
	for str, expectedErr := range testStrings {
		_, err := ParseVersion(str)
		if err == nil || !strings.Contains(err.Error(), expectedErr) {
			t.Errorf("ParseVersion did not produce the expected error for \"%s\"."+
				"\nexpected: %s\nreceived: %+v", str, expectedErr, err)
		}
	}
}

// Happy path.
func TestIsCompatible(t *testing.T) {
	// Test values
	testValues := []struct{ required, current Version }{
		{Version{1, 17, "51"}, Version{1, 17, "51"}},
		{Version{5, 42, "51"}, Version{5, 42, "0"}},
		{Version{0, 14, "51"}, Version{0, 15, ""}},
		{Version{9, 72, ""}, Version{9, 73, ""}},
	}

	for i, v := range testValues {
		if !IsCompatible(v.required, v.current) {
			t.Errorf("IsCompatible incorectly determined the current version"+
				"is not comptabile with the required version (%d)."+
				"\nrequired: %s\ncurrent:  %s", i, v.required, v.current)
		}
	}
}

// Error path: tests multiple cases of invalid Version objects.
func TestIsCompatible_Failure(t *testing.T) {
	// Test values
	testValues := []struct{ required, current Version }{
		{Version{5, 41, "patch5"}, Version{5, 40, "patch5"}},
		{Version{4, 15, "patch0"}, Version{5, 15, "patch0"}},
		{Version{0, 14, "patch9"}, Version{0, 13, ""}},
	}

	for i, v := range testValues {
		if IsCompatible(v.required, v.current) {
			t.Errorf("IsCompatible incorectly determined the current version"+
				"is comptabile with the required version (%d)."+
				"\nrequired: %s\ncurrent:  %s", i, v.required, v.current)
		}
	}
}

// Happy path: tests that the same Version objects are equal.
func TestEqual_Same(t *testing.T) {
	// Test values
	testValues := []struct{ a, b Version }{
		{Version{}, Version{}},
		{Version{1, 1, "1"}, Version{1, 1, "1"}},
		{Version{4, 15, "patch0"}, Version{4, 15, "patch0"}},
	}

	for i, v := range testValues {
		if !Equal(v.a, v.b) {
			t.Errorf("Equal determined the versions are not equal (%d)."+
				"\na: %s\nb: %s", i, v.a, v.b)
		}
	}
}

// Happy path: tests that differing Version objects are not equal.
func TestEqual_Different(t *testing.T) {
	// Test values
	testValues := []struct{ a, b Version }{
		{Version{major: 1}, Version{minor: 1}},
		{Version{1, 0, "1"}, Version{1, 2, "1"}},
		{Version{1, 1, "1"}, Version{1, 1, "2"}},
		{Version{1, 1, "1"}, Version{2, 1, "1"}},
	}

	for i, v := range testValues {
		if Equal(v.a, v.b) {
			t.Errorf("Equal determined the versions are equal (%d)."+
				"\na: %s\nb: %s", i, v.a, v.b)
		}
	}
}

// Happy path: tests that the same Version objects are equal.
func TestCmp(t *testing.T) {
	// Test values
	testValues := []struct {
		expected int
		a, b     Version
	}{
		{0, Version{}, Version{}},
		{0, Version{1, 1, ""}, Version{1, 1, ""}},
		{0, Version{4, 15, ""}, Version{4, 15, ""}},
		{1, Version{4, 15, ""}, Version{3, 15, ""}},
		{1, Version{4, 15, ""}, Version{4, 14, ""}},
		{-1, Version{4, 15, ""}, Version{5, 15, ""}},
		{-1, Version{4, 15, ""}, Version{4, 16, ""}},
	}

	for i, v := range testValues {
		test := Cmp(v.a, v.b)
		if v.expected != test {
			t.Errorf("Cmp did not return the expected value for %s and %s (%d)."+
				"\nexpected: %d\nreceived: %d", v.a, v.b, i, v.expected, test)
		}
	}
}
