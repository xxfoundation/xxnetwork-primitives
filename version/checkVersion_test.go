////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package version

import (
	"testing"
)

// Tests that String() returns a correctly formatted version string.
func TestVersion_String(t *testing.T) {
	// Test values
	testVersion := Version{5, 10, "test"}
	expectedVersion := "5.10.test"

	// Get the string
	testString := testVersion.String()

	if testString != expectedVersion {
		t.Errorf("String() did not return the correct string"+
			"\n\texpected: %s\n\treceived: %s", expectedVersion, testString)
	}
}

// Tests that ParseVersion() creates the correct Version for the given string.
func TestParseVersion_HappyPath(t *testing.T) {
	// Test values
	expectedVersion := Version{1, 2, "3456"}

	// Parse the test version
	testVersion, err := ParseVersion(expectedVersion.String())

	// Make sure no errors occurred
	if err != nil {
		t.Errorf("ParseVersion() produced an unexpected error: %+v", err)
	}

	// Check the major value
	if testVersion.major != expectedVersion.major {
		t.Errorf("ParseVersion() produced a version with an incorrect major "+
			"value\n\texpected: %v\n\treceived: %v",
			expectedVersion.major, testVersion.major)
	}

	// Check the minor value
	if testVersion.minor != expectedVersion.minor {
		t.Errorf("ParseVersion() produced a version with an incorrect minor "+
			"value\n\texpected: %v\n\treceived: %v",
			expectedVersion.minor, testVersion.minor)
	}

	// Check the patch
	if testVersion.patch != expectedVersion.patch {
		t.Errorf("ParseVersion() produced a version with an incorrect patch"+
			"\n\texpected: %v\n\treceived: %v",
			expectedVersion.patch, testVersion.patch)
	}
}

// Tests that ParseVersion() returns an error for various incorrectly formatted
// version strings.
func TestParseVersion_Error(t *testing.T) {
	// Test values
	testStrings := []string{
		"",
		"0",
		"0.",
		"0.0",
		"0.0.",
		"a.0.0",
		"0.a.0",
		"a.a.a",
		"0.0.0.",
		"0.0.0.0",
		".",
		"..",
		"...",
	}

	// Check that all the test strings produce errors
	for _, testString := range testStrings {
		_, err := ParseVersion(testString)
		if err == nil {
			t.Errorf("ParseVersion() did not produce an error for the string "+
				"%#v.", testString)
		}
	}
}

// Tests that IsCompatible() correctly determine that the given versions are
// compatible with our version.
func Test_IsCompatible_HappyPath(t *testing.T) {
	ourVersion := Version{1, 7, "51"}
	theirVersions := []Version{
		{1, 7, "51"},
		{1, 7, "0"},
		{1, 12, "test"},
		{1, 00011, ""},
	}

	for _, theirVersion := range theirVersions {
		if !IsCompatible(ourVersion, theirVersion) {
			t.Errorf("IsCompatible() incorectly determined %+v and %+v to "+
				"not be compatible", ourVersion, theirVersion)
		}
	}
}

// Tests that IsCompatible() correctly determine that the given versions are
// incompatible with our version.
func Test_IsCompatible_Failure(t *testing.T) {
	ourVersion := Version{5, 4, "51"}
	theirVersions := []Version{
		{5, 3, "51"},
		{5, 0, "0"},
		{4, 4, "51"},
		{0, 4, "51"},
		{0, 1, ""},
		{3, 4, "51"},
	}

	for _, theirVersion := range theirVersions {
		if IsCompatible(ourVersion, theirVersion) {
			t.Errorf("IsCompatible() incorectly determined %+v and %+v to be "+
				"compatible", ourVersion, theirVersion)
		}
	}
}
